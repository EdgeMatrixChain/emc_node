package types

import (
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/emc-protocol/edge-matrix/helper/keccak"
)

type TeleType byte

const (
	LegacyTx TeleType = 0x0
	StateTx  TeleType = 0x7f

	StateTransactionGasLimit = 1000000 // some arbitrary default gas limit for state transactions
)

func txTypeFromByte(b byte) (TeleType, error) {
	tt := TeleType(b)

	switch tt {
	case LegacyTx, StateTx:
		return tt, nil
	default:
		return tt, fmt.Errorf("unknown transaction type: %d", b)
	}
}

func (t TeleType) String() (s string) {
	switch t {
	case LegacyTx:
		return "LegacyTx"
	case StateTx:
		return "StateTx"
	}

	return
}

type Telegram struct {
	Nonce    uint64
	GasPrice *big.Int
	Gas      uint64
	To       *Address
	Value    *big.Int
	Input    []byte
	V        *big.Int
	R        *big.Int
	S        *big.Int

	Hash Hash
	From Address

	RespV    *big.Int
	RespR    *big.Int
	RespS    *big.Int
	RespHash Hash
	RespFrom Address

	Type TeleType

	// Cache
	size atomic.Value
}

// IsContractCreation checks if tx is contract creation
func (t *Telegram) IsContractCreation() bool {
	return t.To == nil
}

// ComputeHash computes the hash of the transaction
func (t *Telegram) ComputeHash() *Telegram {
	ar := marshalArenaPool.Get()
	hash := keccak.DefaultKeccakPool.Get()

	v := t.MarshalRLPWith(ar)
	hash.WriteRlp(t.Hash[:0], v)

	marshalArenaPool.Put(ar)
	keccak.DefaultKeccakPool.Put(hash)

	return t
}

func (t *Telegram) Copy() *Telegram {
	tt := new(Telegram)
	*tt = *t

	tt.GasPrice = new(big.Int)
	if t.GasPrice != nil {
		tt.GasPrice.Set(t.GasPrice)
	}

	tt.Value = new(big.Int)
	if t.Value != nil {
		tt.Value.Set(t.Value)
	}

	if t.R != nil {
		tt.R = new(big.Int)
		tt.R = big.NewInt(0).SetBits(t.R.Bits())
	}

	if t.S != nil {
		tt.S = new(big.Int)
		tt.S = big.NewInt(0).SetBits(t.S.Bits())
	}

	tt.Input = make([]byte, len(t.Input))
	copy(tt.Input[:], t.Input[:])

	return tt
}

// Cost returns gas * gasPrice + value
func (t *Telegram) Cost() *big.Int {
	total := new(big.Int).Mul(t.GasPrice, new(big.Int).SetUint64(t.Gas))
	total.Add(total, t.Value)

	return total
}

func (t *Telegram) Size() uint64 {
	if size := t.size.Load(); size != nil {
		sizeVal, ok := size.(uint64)
		if !ok {
			return 0
		}

		return sizeVal
	}

	size := uint64(len(t.MarshalRLP()))
	t.size.Store(size)

	return size
}

func (t *Telegram) ExceedsBlockGasLimit(blockGasLimit uint64) bool {
	return t.Gas > blockGasLimit
}

func (t *Telegram) IsUnderpriced(priceLimit uint64) bool {
	return t.GasPrice.Cmp(big.NewInt(0).SetUint64(priceLimit)) < 0
}
