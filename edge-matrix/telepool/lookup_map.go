package telepool

import (
	"sync"

	"github.com/emc-protocol/edge-matrix/types"
)

// Lookup map used to find transactions present in the pool
type lookupMap struct {
	sync.RWMutex
	all map[types.Hash]*types.Telegram
}

// add inserts the given transaction into the map. Returns false
// if it already exists. [thread-safe]
func (m *lookupMap) add(msg *types.Telegram) bool {
	m.Lock()
	defer m.Unlock()

	if _, exists := m.all[msg.Hash]; exists {
		return false
	}

	m.all[msg.Hash] = msg

	return true
}

// remove removes the given transactions from the map. [thread-safe]
func (m *lookupMap) remove(msgs ...*types.Telegram) {
	m.Lock()
	defer m.Unlock()

	for _, msg := range msgs {
		delete(m.all, msg.Hash)
	}
}

// get returns the transaction associated with the given hash. [thread-safe]
func (m *lookupMap) get(hash types.Hash) (*types.Telegram, bool) {
	m.RLock()
	defer m.RUnlock()

	tx, ok := m.all[hash]
	if !ok {
		return nil, false
	}

	return tx, true
}
