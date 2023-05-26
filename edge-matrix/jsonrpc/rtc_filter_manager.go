package jsonrpc

import (
	"encoding/json"
	"fmt"
	"github.com/emc-protocol/edge-matrix/rtc"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"sync"
	"time"
)

// filterManagerStore provides methods required by FilterManager
type rtcFilterManagerStore interface {
	// SubscribeEvents subscribes for chain head events
	SubscribeRtcEvents() rtc.Subscription
}

// RtcFilterManager manages all running rtc filters
type RtcFilterManager struct {
	sync.RWMutex

	logger hclog.Logger

	timeout time.Duration

	store rtcFilterManagerStore
	//topic *network.Topic
	subscription rtc.Subscription
	//rtcStream    *rtcStream

	filters  map[string]filter
	timeouts timeHeapImpl

	updateCh chan struct{}
	closeCh  chan struct{}
}

func NewRtcFilterManager(logger hclog.Logger, store rtcFilterManagerStore) *RtcFilterManager {
	m := &RtcFilterManager{
		logger:   logger.Named("rtc-filter"),
		timeout:  defaultTimeout,
		store:    store,
		filters:  make(map[string]filter),
		timeouts: timeHeapImpl{},
		updateCh: make(chan struct{}),
		closeCh:  make(chan struct{}),
	}

	m.subscription = store.SubscribeRtcEvents()

	return m
}

// Run starts worker process to handle events
func (f *RtcFilterManager) Run() {
	// watch for new events in the blockchain
	watchCh := make(chan *rtc.Event)

	go func() {
		for {
			evnt := f.subscription.GetEvent()
			if evnt == nil {
				return
			}
			watchCh <- evnt
		}
	}()

	var timeoutCh <-chan time.Time

	for {
		// check for the next filter to be removed
		filterID, filterExpiresAt := f.nextTimeoutFilter()

		// set timer to remove filter
		if filterID != "" {
			timeoutCh = time.After(time.Until(filterExpiresAt))
		}

		select {
		case evnt := <-watchCh:
			// new blockchain event
			if err := f.dispatchEvent(evnt); err != nil {
				f.logger.Error("failed to dispatch event", "err", err)
			}

		case <-timeoutCh:
			// timeout for filter
			// if filter still exists
			if !f.Uninstall(filterID) {
				f.logger.Warn("failed to uninstall filter", "id", filterID)
			}

		case <-f.updateCh:
			// filters change, reset the loop to start the timeout timer

		case <-f.closeCh:
			// stop the filter manager
			return
		}
	}
}

// Close closed closeCh so that terminate worker
func (f *RtcFilterManager) Close() {
	close(f.closeCh)
}

// newRtcFilterBase initializes filterBase with unique ID
func newRtcFilterBase(ws wsConn) filterBase {
	return filterBase{
		id:        uuid.New().String(),
		ws:        ws,
		heapIndex: NoIndexInHeap,
	}
}

// rtcFilter is a filter to store subjects that meet the conditions in query
// logFilter is a filter to store logs that meet the conditions in query
type rtcFilter struct {
	filterBase
	sync.Mutex

	query *RtcQuery
	msgs  []*rtc.RtcMsg
}

// appendLog appends new log to logs
func (f *rtcFilter) appendLog(msg *rtc.RtcMsg) {
	f.Lock()
	defer f.Unlock()

	f.msgs = append(f.msgs, msg)
}

// takeRtcMsgUpdates returns all saved logs in filter and set new log slice
func (f *rtcFilter) takeRtcMsgUpdates() []*rtc.RtcMsg {
	f.Lock()
	defer f.Unlock()

	msgs := f.msgs
	f.msgs = []*rtc.RtcMsg{} // create brand-new slice so that prevent new msgs from being added to current msgs

	return msgs
}

// getUpdates returns stored logs in string
func (f *rtcFilter) getUpdates() (interface{}, error) {
	logs := f.takeRtcMsgUpdates()

	return logs, nil
}

// sendUpdates writes stored logs to web socket stream
func (f *rtcFilter) sendUpdates() error {
	updates := f.takeRtcMsgUpdates()

	for _, msg := range updates {
		res, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		if err := f.writeMessageToWs(string(res)); err != nil {
			return err
		}
	}

	return nil
}

// Exists checks the filter with given ID exists
func (f *RtcFilterManager) Exists(id string) bool {
	f.RLock()
	defer f.RUnlock()

	_, ok := f.filters[id]

	return ok
}

// GetFilterChanges returns the updates of the filter with given ID in string, and refreshes the timeout on the filter
func (f *RtcFilterManager) GetFilterChanges(id string) (interface{}, error) {
	filter, res, err := f.getFilterAndChanges(id)

	if err == nil && !filter.hasWSConn() {
		// Refresh the timeout on this filter
		f.Lock()
		f.refreshFilterTimeout(filter.getFilterBase())
		f.Unlock()
	}

	return res, err
}

// getFilterAndChanges returns the updates of the filter with given ID in string (read lock only)
func (f *RtcFilterManager) getFilterAndChanges(id string) (filter, interface{}, error) {
	f.RLock()
	defer f.RUnlock()

	filter, ok := f.filters[id]

	if !ok {
		return nil, nil, ErrFilterNotFound
	}

	// we cannot get updates from a ws filter with getFilterChanges
	if filter.hasWSConn() {
		return nil, nil, ErrWSFilterDoesNotSupportGetChanges
	}

	res, err := filter.getUpdates()
	if err != nil {
		return nil, nil, err
	}

	return filter, res, nil
}

// Uninstall removes the filter with given ID from list
func (f *RtcFilterManager) Uninstall(id string) bool {
	f.Lock()
	defer f.Unlock()

	return f.removeFilterByID(id)
}

// removeFilterByID removes the filter with given ID [NOT Thread Safe]
func (f *RtcFilterManager) removeFilterByID(id string) bool {
	// Make sure filter exists
	filter, ok := f.filters[id]
	if !ok {
		return false
	}

	delete(f.filters, id)

	if removed := f.timeouts.removeFilter(filter.getFilterBase()); removed {
		f.emitSignalToUpdateCh()
	}

	return true
}

// RemoveFilterByWs removes the filter with given WS [Thread safe]
func (f *RtcFilterManager) RemoveFilterByWs(ws wsConn) {
	f.Lock()
	defer f.Unlock()

	f.removeFilterByID(ws.GetFilterID())
}

// refreshFilterTimeout updates the timeout for a filter to the current time
func (f *RtcFilterManager) refreshFilterTimeout(filter *filterBase) {
	f.timeouts.removeFilter(filter)
	f.addFilterTimeout(filter)
}

// addFilterTimeout set timeout and add to heap
func (f *RtcFilterManager) addFilterTimeout(filter *filterBase) {
	filter.expiresAt = time.Now().Add(f.timeout)
	f.timeouts.addFilter(filter)
	f.emitSignalToUpdateCh()
}

// addFilter is an internal method to add given filter to list and heap
func (f *RtcFilterManager) addFilter(filter filter) string {
	f.Lock()
	defer f.Unlock()

	base := filter.getFilterBase()

	f.filters[base.id] = filter

	// Set timeout and add to heap if filter doesn't have web socket connection
	if !filter.hasWSConn() {
		f.addFilterTimeout(base)
	}

	return base.id
}

func (f *RtcFilterManager) emitSignalToUpdateCh() {
	select {
	// notify worker of new filter with timeout
	case f.updateCh <- struct{}{}:
	default:
	}
}

// nextTimeoutFilter returns the filter that will be expired next
// nextTimeoutFilter returns the only filter with timeout
func (f *RtcFilterManager) nextTimeoutFilter() (string, time.Time) {
	f.RLock()
	defer f.RUnlock()

	if len(f.timeouts) == 0 {
		return "", time.Time{}
	}

	// peek the first item
	base := f.timeouts[0]

	return base.id, base.expiresAt
}

// dispatchEvent is an event handler for new block event
func (f *RtcFilterManager) dispatchEvent(evnt *rtc.Event) error {
	// store new event in each filters
	f.processEvent(evnt)

	// send data to web socket stream
	if err := f.flushWsFilters(); err != nil {
		return err
	}

	return nil
}

// processEvent makes each filter append the new data that interests them
func (f *RtcFilterManager) processEvent(evnt *rtc.Event) {
	f.RLock()
	defer f.RUnlock()

	for _, newMsg := range evnt.NewMsgs {
		if processErr := f.appendRtcMsgsToFilters(newMsg); processErr != nil {
			f.logger.Error(fmt.Sprintf("Unable to process new RtcMsg, %v", processErr))
		}

	}
}

// NewRtcFilter adds new RtcFilter
func (f *RtcFilterManager) NewRtcFilter(rtcQuery *RtcQuery, ws wsConn) string {
	filter := &rtcFilter{
		filterBase: newRtcFilterBase(ws),
		query:      rtcQuery,
	}

	if filter.hasWSConn() {
		ws.SetFilterID(filter.id)
	}

	return f.addFilter(filter)
}

// appendLogsToFilters makes each LogFilters append logs in the msg
func (f *RtcFilterManager) appendRtcMsgsToFilters(msg *rtc.RtcMsg) error {
	// Get rtcFilters from filters
	rtcFilters := make([]*rtcFilter, 0)

	for _, ft := range f.filters {
		if rf, ok := ft.(*rtcFilter); ok {
			rtcFilters = append(rtcFilters, rf)
		}
	}

	if len(rtcFilters) == 0 {
		return nil
	}

	for _, ft := range rtcFilters {
		if ft.query.Match(msg) {
			ft.appendLog(&rtc.RtcMsg{
				To:          msg.To,
				Subject:     msg.Subject,
				Application: msg.Application,
				Content:     msg.Content,
				V:           msg.V,
				R:           msg.R,
				S:           msg.S,
				Hash:        msg.Hash,
				From:        msg.From,
				Type:        msg.Type,
			})
		}
	}
	return nil
}

// flushWsFilters make each filters with web socket connection write the updates to web socket stream
// flushWsFilters also removes the filters if flushWsFilters notices the connection is closed
func (f *RtcFilterManager) flushWsFilters() error {
	closedFilterIDs := make([]string, 0)

	f.RLock()

	for id, filter := range f.filters {
		if !filter.hasWSConn() {
			continue
		}

		if flushErr := filter.sendUpdates(); flushErr != nil {
			f.logger.Error(fmt.Sprintf("Unable to process flush, %v", flushErr))
			// mark as closed if the connection is closed
			//if errors.Is(flushErr, websocket.ErrCloseSent) || errors.Is(flushErr, net.ErrClosed) {
			closedFilterIDs = append(closedFilterIDs, id)

			f.logger.Warn(fmt.Sprintf("Subscription %s has been closed", id))

			continue
			//}

		}
	}

	f.RUnlock()

	// remove filters with closed web socket connections from FilterManager
	if len(closedFilterIDs) > 0 {
		f.Lock()
		for _, id := range closedFilterIDs {
			f.removeFilterByID(id)
		}
		f.Unlock()

		f.logger.Info(fmt.Sprintf("Removed %d filters due to closed connections", len(closedFilterIDs)))
	}

	return nil
}
