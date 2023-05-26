package state

import (
	"math/big"
	"testing"

	"github.com/emc-protocol/edge-matrix/state/runtime"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func newTestTransition(preState map[types.Address]*PreState) *Transition {
	if preState == nil {
		preState = defaultPreState
	}

	return &Transition{
		logger: hclog.NewNullLogger(),
		state:  newTestTxn(preState),
	}
}

func TestTransfer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		preState    map[types.Address]*PreState
		from        types.Address
		to          types.Address
		amount      int64
		expectedErr error
	}{
		{
			name: "should succeed",
			preState: map[types.Address]*PreState{
				addr1: {
					Nonce:   0,
					Balance: 1000,
				},
				addr2: {
					Nonce:   0,
					Balance: 0,
				},
			},
			from:        addr1,
			to:          addr2,
			amount:      100,
			expectedErr: nil,
		},
		{
			name: "should fail by ErrNotEnoughFunds",
			preState: map[types.Address]*PreState{
				addr1: {
					Nonce:   0,
					Balance: 10,
					State:   map[types.Hash]types.Hash{},
				},
			},
			from:   addr1,
			to:     addr2,
			amount: 100,
			// should return ErrInsufficientBalance when state.transfer returns ErrNotEnoughFunds
			expectedErr: runtime.ErrInsufficientBalance,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			transition := newTestTransition(tt.preState)

			amount := big.NewInt(tt.amount)
			err := transition.Transfer(tt.from, tt.to, amount)

			assert.Equal(t, tt.expectedErr, err)
			if err == nil {
				// should move balance
				oldBalanceOfFrom := big.NewInt(int64(tt.preState[tt.from].Balance))
				oldBalanceOfTo := big.NewInt(int64(tt.preState[tt.to].Balance))
				newBalanceOfFrom := transition.GetBalance(tt.from)
				newBalanceOfTo := transition.GetBalance(tt.to)
				diffOfFrom := new(big.Int).Sub(newBalanceOfFrom, oldBalanceOfFrom)
				diffOfTo := new(big.Int).Sub(newBalanceOfTo, oldBalanceOfTo)

				assert.Zero(t, diffOfFrom.Cmp(new(big.Int).Mul(big.NewInt(-1), amount)))
				assert.Zero(t, diffOfTo.Cmp(amount))
			}
		})
	}
}
