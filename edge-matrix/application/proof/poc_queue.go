package proof

import (
	"container/heap"
	"sync"
)

type PocPriority uint64

const (
	PriorityRequestedPoc PocPriority = 1
	PriorityPushPoc      PocPriority = 10
)

// PocQueue is a queue that holds poc tasks for miner node, implemented as a min-heap
type PocQueue struct {
	sync.Mutex

	heap  dialQueueImpl
	tasks map[string]*PocTask

	updateCh chan struct{}
	closeCh  chan struct{}
}

// NewPocQueue creates a new PocQueue instance
func NewPocQueue() *PocQueue {
	return &PocQueue{
		heap:     dialQueueImpl{},
		tasks:    map[string]*PocTask{},
		updateCh: make(chan struct{}),
		closeCh:  make(chan struct{}),
	}
}

// Close closes the running PocQueue
func (d *PocQueue) Close() {
	close(d.closeCh)
}

// PopTask is a loop that handles update and close events [BLOCKING]
func (d *PocQueue) PopTask() *PocTask {
	for {
		task := d.popTaskImpl() // Blocking pop
		if task != nil {
			return task
		}

		select {
		case <-d.updateCh:
		case <-d.closeCh:
			return nil
		}
	}
}

// PopTask is a loop that handles update and close events [BLOCKING]
func (d *PocQueue) Len() int {
	return len(d.tasks)
}

// popTaskImpl is the implementation for task popping from the min-heap
func (d *PocQueue) popTaskImpl() *PocTask {
	d.Lock()
	defer d.Unlock()

	if len(d.heap) != 0 {
		// pop the first value and remove it from the heap
		tt := heap.Pop(&d.heap)

		task, ok := tt.(*PocTask)
		if !ok {
			return nil
		}

		return task
	}

	return nil
}

// DeleteTask deletes a task from the dial queue for the specified peer
func (d *PocQueue) DeleteTask(peer string) {
	d.Lock()
	defer d.Unlock()

	item, ok := d.tasks[peer]
	if ok {
		// negative index for popped element
		if item.index >= 0 {
			heap.Remove(&d.heap, item.index)
		}

		delete(d.tasks, peer)
	}
}

// AddTask adds a new task to the dial queue
func (d *PocQueue) AddTask(
	pocData *PocCpuData,
	priority PocPriority,
) {
	d.Lock()
	defer d.Unlock()

	task := &PocTask{
		pocCpuData: pocData,
		priority:   uint64(priority),
	}

	d.tasks[pocData.Validator] = task
	heap.Push(&d.heap, task)

	select {
	case d.updateCh <- struct{}{}:
	default:
	}
}
