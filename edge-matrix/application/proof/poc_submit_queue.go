package proof

import (
	"container/heap"
	"sync"
)

// PocQueue is a queue that holds poc tasks for miner node, implemented as a min-heap
type PocSubmitQueue struct {
	sync.Mutex

	heap  submitQueueImpl
	tasks map[string]*PocSubmitTask

	updateCh chan struct{}
	closeCh  chan struct{}
}

// NewPocSubmitQueue creates a new PocSubmitQueue instance
func NewPocSubmitQueue() *PocSubmitQueue {
	return &PocSubmitQueue{
		heap:     submitQueueImpl{},
		tasks:    map[string]*PocSubmitTask{},
		updateCh: make(chan struct{}),
		closeCh:  make(chan struct{}),
	}
}

// Close closes the running PocQueue
func (d *PocSubmitQueue) Close() {
	close(d.closeCh)
}

// PopTask is a loop that handles update and close events [BLOCKING]
func (d *PocSubmitQueue) PopTask() *PocSubmitTask {
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
func (d *PocSubmitQueue) Len() int {
	return len(d.tasks)
}

// popTaskImpl is the implementation for task popping from the min-heap
func (d *PocSubmitQueue) popTaskImpl() *PocSubmitTask {
	d.Lock()
	defer d.Unlock()

	if len(d.heap) != 0 {
		// pop the first value and remove it from the heap
		tt := heap.Pop(&d.heap)

		task, ok := tt.(*PocSubmitTask)
		if !ok {
			return nil
		}

		return task
	}

	return nil
}

// DeleteTask deletes a task from the dial queue for the specified peer
func (d *PocSubmitQueue) DeleteTask(peer string) {
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
func (d *PocSubmitQueue) AddTask(
	pocSubmitData *PocSubmitData,
	priority PocPriority,
) {
	d.Lock()
	defer d.Unlock()

	task := &PocSubmitTask{
		pocSubmitData: pocSubmitData,
		priority:      uint64(priority),
	}

	d.tasks[pocSubmitData.TargetNodeID] = task
	heap.Push(&d.heap, task)

	select {
	case d.updateCh <- struct{}{}:
	default:
	}
}
