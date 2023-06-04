package proof

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPocQueue(t *testing.T) {
	q := NewPocQueue()

	info0 := &PocCpuData{
		NodeId:   "a",
		Seed:     "0x0a0b0c0d0e0f0102",
		BlockNum: 100,
	}
	q.AddTask(info0, PriorityRequestedPoc)
	assert.Equal(t, 1, q.heap.Len())

	info1 := &PocCpuData{
		NodeId:   "b",
		Seed:     "0x0a0b0c0d0e0f0102",
		BlockNum: 100,
	}
	q.AddTask(info1, PriorityRequestedPoc)
	assert.Equal(t, 2, q.heap.Len())

	assert.Equal(t, q.popTaskImpl().pocCpuData.NodeId, "a")
	assert.Equal(t, q.popTaskImpl().pocCpuData.NodeId, "b")
	assert.Equal(t, 0, q.heap.Len())

	assert.Nil(t, q.popTaskImpl())

	done := make(chan struct{})

	go func() {
		q.PopTask()
		done <- struct{}{}
	}()

	// we should not get any task now
	select {
	case <-done:
		t.Fatal("not expected")
	case <-time.After(1 * time.Second):
	}

	q.AddTask(info0, 1)

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}
