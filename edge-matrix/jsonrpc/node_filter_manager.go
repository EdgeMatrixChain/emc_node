package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emc-protocol/edge-matrix/application"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-hclog"
	"net"
	"sync"
	"time"
)

// nodeFilterManagerStore provides methods required by NodeFilterManager
type nodeFilterManagerStore interface {
	// SubscribeAppEvents subscribes for chain head events
	SubscribeAppEvents() application.Subscription
}

// NodeFilterManager manages all running node filters
type NodeFilterManager struct {
	sync.RWMutex

	logger hclog.Logger

	timeout time.Duration

	store nodeFilterManagerStore
	//topic *network.Topic
	subscription application.Subscription

	filters  map[string]filter
	timeouts timeHeapImpl

	updateCh chan struct{}
	closeCh  chan struct{}
}

// nodeFilter is a filter to store node that meet the conditions in query
type nodeFilter struct {
	filterBase
	sync.Mutex

	query *NodeQuery
	msgs  []*application.Application
}

// appendLog appends new log to logs
func (f *nodeFilter) appendLog(msg *application.Application) {
	f.Lock()
	defer f.Unlock()

	f.msgs = append(f.msgs, msg)
}

// takeNodeMsgUpdates returns all online node msg  in filter
func (f *nodeFilter) takeNodeMsgUpdates() []*application.Application {
	f.Lock()
	defer f.Unlock()

	msgs := f.msgs
	f.msgs = []*application.Application{} // create brand-new slice so that prevent new msgs from being added to current msgs

	return msgs
}

// getUpdates returns stored node msg in string
func (f *nodeFilter) getUpdates() (interface{}, error) {
	logs := f.takeNodeMsgUpdates()

	return logs, nil
}

// sendUpdates writes stored logs to web socket stream
func (f *nodeFilter) sendUpdates() error {
	updates := f.takeNodeMsgUpdates()

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

func NewNodeFilterManager(logger hclog.Logger, store nodeFilterManagerStore) *NodeFilterManager {
	m := &NodeFilterManager{
		logger:   logger.Named("node-filter"),
		timeout:  defaultTimeout,
		store:    store,
		filters:  make(map[string]filter),
		timeouts: timeHeapImpl{},
		updateCh: make(chan struct{}),
		closeCh:  make(chan struct{}),
	}

	m.subscription = store.SubscribeAppEvents()

	return m
}

// Run starts worker process to handle events
func (f *NodeFilterManager) Run() {
	// watch for new events in the blockchain
	watchCh := make(chan *application.Event)

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
func (f *NodeFilterManager) Close() {
	close(f.closeCh)
}

// newNodeFilterBase initializes filterBase with unique ID
func newNodeFilterBase(ws wsConn) filterBase {
	return filterBase{
		id:        uuid.New().String(),
		ws:        ws,
		heapIndex: NoIndexInHeap,
	}
}

// Exists checks the filter with given ID exists
func (f *NodeFilterManager) Exists(id string) bool {
	f.RLock()
	defer f.RUnlock()

	_, ok := f.filters[id]

	return ok
}

// GetFilterChanges returns the updates of the filter with given ID in string, and refreshes the timeout on the filter
func (f *NodeFilterManager) GetFilterChanges(id string) (interface{}, error) {
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
func (f *NodeFilterManager) getFilterAndChanges(id string) (filter, interface{}, error) {
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
func (f *NodeFilterManager) Uninstall(id string) bool {
	f.Lock()
	defer f.Unlock()

	return f.removeFilterByID(id)
}

// removeFilterByID removes the filter with given ID [NOT Thread Safe]
func (f *NodeFilterManager) removeFilterByID(id string) bool {
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
func (f *NodeFilterManager) RemoveFilterByWs(ws wsConn) {
	f.Lock()
	defer f.Unlock()

	f.removeFilterByID(ws.GetFilterID())
}

// refreshFilterTimeout updates the timeout for a filter to the current time
func (f *NodeFilterManager) refreshFilterTimeout(filter *filterBase) {
	f.timeouts.removeFilter(filter)
	f.addFilterTimeout(filter)
}

// addFilterTimeout set timeout and add to heap
func (f *NodeFilterManager) addFilterTimeout(filter *filterBase) {
	filter.expiresAt = time.Now().Add(f.timeout)
	f.timeouts.addFilter(filter)
	f.emitSignalToUpdateCh()
}

// addFilter is an internal method to add given filter to list and heap
func (f *NodeFilterManager) addFilter(filter filter) string {
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

func (f *NodeFilterManager) emitSignalToUpdateCh() {
	select {
	// notify worker of new filter with timeout
	case f.updateCh <- struct{}{}:
	default:
	}
}

// nextTimeoutFilter returns the filter that will be expired next
// nextTimeoutFilter returns the only filter with timeout
func (f *NodeFilterManager) nextTimeoutFilter() (string, time.Time) {
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
func (f *NodeFilterManager) dispatchEvent(evnt *application.Event) error {
	// store new event in each filters
	f.processEvent(evnt)

	// send data to web socket stream
	if err := f.flushWsFilters(); err != nil {
		return err
	}

	return nil
}

// processEvent makes each filter append the new data that interests them
func (f *NodeFilterManager) processEvent(evnt *application.Event) {
	f.RLock()
	defer f.RUnlock()

	for _, newMsg := range evnt.NewApp {
		if processErr := f.appendNodeLogToFilters(newMsg); processErr != nil {
			f.logger.Error(fmt.Sprintf("Unable to process new RtcMsg, %v", processErr))
		}

	}
}

// NewRtcFilter adds new RtcFilter
func (f *NodeFilterManager) NewNodeFilter(nodeQuery *NodeQuery, ws wsConn) string {
	filter := &nodeFilter{
		filterBase: newNodeFilterBase(ws),
		query:      nodeQuery,
	}

	if filter.hasWSConn() {
		ws.SetFilterID(filter.id)
	}

	return f.addFilter(filter)
}

// appendNodeLogToFilters makes each NodeFilters append logs in the msg
func (f *NodeFilterManager) appendNodeLogToFilters(msg *application.Application) error {
	// Get nodeFilters from filters
	nodeFilters := make([]*nodeFilter, 0)

	for _, ft := range f.filters {
		if rf, ok := ft.(*nodeFilter); ok {
			nodeFilters = append(nodeFilters, rf)
		}
	}

	if len(nodeFilters) == 0 {
		return nil
	}

	for _, ft := range nodeFilters {
		if ft.query.Match(msg) {
			ft.appendLog(msg.Copy())
		}
	}
	return nil
}

// flushWsFilters make each filters with web socket connection write the updates to web socket stream
// flushWsFilters also removes the filters if flushWsFilters notices the connection is closed
func (f *NodeFilterManager) flushWsFilters() error {
	closedFilterIDs := make([]string, 0)

	f.RLock()

	for id, filter := range f.filters {
		if !filter.hasWSConn() {
			continue
		}

		if flushErr := filter.sendUpdates(); flushErr != nil {
			f.logger.Error(fmt.Sprintf("Unable to process flush, %v", flushErr))
			// mark as closed if the connection is closed
			if errors.Is(flushErr, websocket.ErrCloseSent) || errors.Is(flushErr, net.ErrClosed) {
				closedFilterIDs = append(closedFilterIDs, id)

				f.logger.Warn(fmt.Sprintf("Subscription %s has been closed", id))

				continue
			}
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
