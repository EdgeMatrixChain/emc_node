package sd

// Lookup map used to find poc present in the pool
type pocLookupMap struct {
	//sync.RWMutex
	all map[string]*PocSdData
}

// add inserts the given transaction into the map. Returns false
// if it already exists. [thread-safe]
func (m *pocLookupMap) add(msg *PocSdData) bool {
	//m.Lock()
	//defer m.Unlock()

	if _, exists := m.all[msg.NodeId]; exists {
		return false
	}

	m.all[msg.NodeId] = msg

	return true
}

// remove removes the given pocData from the map. [thread-safe]
func (m *pocLookupMap) remove(msgs ...*PocSdData) {
	//m.Lock()
	//defer m.Unlock()

	for _, msg := range msgs {
		delete(m.all, msg.NodeId)
	}
}

// get returns the pocData associated with the given nodeId. [thread-safe]
func (m *pocLookupMap) get(nodeId string) (*PocSdData, bool) {
	//m.RLock()
	//defer m.RUnlock()

	tx, ok := m.all[nodeId]
	if !ok {
		return nil, false
	}

	return tx, true
}

func (m *pocLookupMap) len() int {
	//m.RLock()
	//defer m.RUnlock()

	return len(m.all)
}

func (m *pocLookupMap) getAll() map[string]*PocSdData {
	//m.RLock()
	//defer m.RUnlock()

	return m.all
}
