package telepool

import "github.com/emc-protocol/edge-matrix/types"

/* QUERY methods */
// Used to query the pool for specific state info.

// GetNonce returns the next nonce for the account
//
// -> Returns the value from the TxPool if the account is initialized in-memory
//
// -> Returns the value from the world state otherwise
func (p *TelegramPool) GetNonce(addr types.Address) uint64 {
	account := p.accounts.get(addr)
	if account == nil {
		stateRoot := p.store.Header().StateRoot
		stateNonce := p.store.GetNonce(stateRoot, addr)

		return stateNonce
	}

	return account.getNonce()
}

// GetCapacity returns the current number of slots
// occupied in the pool as well as the max limit
func (p *TelegramPool) GetCapacity() (uint64, uint64) {
	return p.gauge.read(), p.gauge.max
}

// GetPendingTx returns the transaction by hash in the TxPool (pending txn) [Thread-safe]
func (p *TelegramPool) GetPendingTele(txHash types.Hash) (*types.Telegram, bool) {
	tx, ok := p.index.get(txHash)
	if !ok {
		return nil, false
	}

	return tx, true
}

// GetTxs gets pending and queued transactions
func (p *TelegramPool) GetTeles(inclQueued bool) (
	allPromoted, allEnqueued map[types.Address][]*types.Telegram,
) {
	allPromoted, allEnqueued = p.accounts.allTxs(inclQueued)

	return
}
