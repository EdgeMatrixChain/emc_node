package rtc

import (
	"sync"
)

type void struct{}

// Subscription is the blockchain subscription interface
type Subscription interface {
	GetEventCh() chan *Event
	GetEvent() *Event
	Close()
}

// FOR TESTING PURPOSES //

type MockSubscription struct {
	eventCh chan *Event
}

// subscription is the Blockchain event subscription object
type subscription struct {
	updateCh chan *Event // Channel for update information
	closeCh  chan void   // Channel for close signals
}

// GetEventCh creates a new event channel, and returns it
func (s *subscription) GetEventCh() chan *Event {
	eventCh := make(chan *Event)

	go func() {
		for {
			evnt := s.GetEvent()
			if evnt == nil {
				return
			}
			eventCh <- evnt
		}
	}()

	return eventCh
}

// GetEvent returns the event From the subscription (BLOCKING)
func (s *subscription) GetEvent() *Event {
	for {
		// Wait for an update
		select {
		case ev := <-s.updateCh:
			return ev
		case <-s.closeCh:
			return nil
		}
	}
}

// Close closes the subscription
func (s *subscription) Close() {
	close(s.closeCh)
}

type EventType int

const (
	EventNew EventType = iota // New head event
)

// Event is the blockchain event that gets passed to the listeners
type Event struct {
	//// New part of the chain (or a fork)
	//NewChain []*types.Header

	NewMsgs []*RtcMsg

	// Type is the type of event
	Type EventType

	// Source is the source that generated the blocks for the event
	// right now it can be either the Sealer or the Syncer
	Source string
}

// AddNewHeader appends a header to the event's NewMessage array
func (e *Event) AddNewRtcMsg(newMsg *RtcMsg) {

	if e.NewMsgs == nil {
		// Array doesn't exist yet, create it
		e.NewMsgs = []*RtcMsg{}
	}

	e.NewMsgs = append(e.NewMsgs, newMsg)
}

// SubscribeEvents returns a blockchain event subscription
func (r *Rtc) SubscribeRtcEvents() Subscription {
	return r.stream.subscribe()
}

// eventStream is the structure that contains the event list,
// as well as the update channel which it uses to notify of updates
type eventStream struct {
	sync.Mutex

	// channel to notify updates
	updateCh []chan *Event
}

// subscribe Creates a new blockchain event subscription
func (e *eventStream) subscribe() *subscription {
	return &subscription{
		updateCh: e.newUpdateCh(),
		closeCh:  make(chan void),
	}
}

// newUpdateCh returns the event update channel
func (e *eventStream) newUpdateCh() chan *Event {
	e.Lock()
	defer e.Unlock()

	ch := make(chan *Event, 1)

	if e.updateCh == nil {
		e.updateCh = make([]chan *Event, 0)
	}

	e.updateCh = append(e.updateCh, ch)

	return ch
}

// push adds a new Event, and notifies listeners
func (e *eventStream) push(event *Event) {
	e.Lock()
	defer e.Unlock()

	// Notify the listeners
	for _, update := range e.updateCh {
		select {
		case update <- event:
		default:
		}
	}
}
