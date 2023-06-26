package application

import (
	"github.com/emc-protocol/edge-matrix/application/proof"
	"sync"
)

// Lookup map used to find transactions present in the pool
type PocMap struct {
	sync.RWMutex
	all map[string]*proof.PocCpuRequest
}

// add inserts the given PocCpuRequest into the map. Returns false
// if it already exists. [thread-safe]
func (m *PocMap) add(msg *proof.PocCpuRequest) bool {
	m.Lock()
	defer m.Unlock()

	if _, exists := m.all[msg.NodeId]; exists {
		return false
	}

	m.all[msg.NodeId] = msg

	return true
}

// remove removes the given PocCpuRequests from the map. [thread-safe]
func (m *PocMap) remove(msgs ...*proof.PocCpuRequest) {
	m.Lock()
	defer m.Unlock()

	for _, msg := range msgs {
		delete(m.all, msg.NodeId)
	}
}

// get returns the PocCpuRequest associated with the given nodeId. [thread-safe]
func (m *PocMap) get(nodeId string) (*proof.PocCpuRequest, bool) {
	m.RLock()
	defer m.RUnlock()

	request, ok := m.all[nodeId]
	if !ok {
		return nil, false
	}

	return request, true
}
