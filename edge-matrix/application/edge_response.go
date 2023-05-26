package application

import (
	"fmt"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/umbracle/fastrlp"
	"math/big"
)

type EdgeResponse struct {
	RespString string
	V          *big.Int
	R          *big.Int
	S          *big.Int
	From       types.Address

	Hash types.Hash
}

func (r *EdgeResponse) Copy() *EdgeResponse {
	tt := new(EdgeResponse)
	*tt = *r

	tt.From = r.From

	if len(r.RespString) > 0 {
		tt.RespString = r.RespString
	}

	if r.R != nil {
		tt.R = new(big.Int)
		tt.R = big.NewInt(0).SetBits(r.R.Bits())
	}

	if r.S != nil {
		tt.S = new(big.Int)
		tt.S = big.NewInt(0).SetBits(r.S.Bits())
	}

	return tt
}

func (r *EdgeResponse) MarshalRLP() []byte {
	return r.MarshalRLPTo(nil)
}

func (r *EdgeResponse) MarshalRLPTo(dst []byte) []byte {
	return types.MarshalRLPTo(r.MarshalRLPWith, dst)
}

// MarshalRLPWith marshals the EdgeResponse to RLP with a specific fastrlp.Arena
func (r *EdgeResponse) MarshalRLPWith(arena *fastrlp.Arena) *fastrlp.Value {
	vv := arena.NewArray()

	// RespString may be empty
	if len(r.RespString) > 0 {
		vv.Set(arena.NewString(r.RespString))
	} else {
		vv.Set(arena.NewNull())
	}
	// signature values
	vv.Set(arena.NewBigInt(r.V))
	vv.Set(arena.NewBigInt(r.R))
	vv.Set(arena.NewBigInt(r.S))

	vv.Set(arena.NewBytes((r.From).Bytes()))
	vv.Set(arena.NewBytes((r.Hash).Bytes()))

	return vv
}

func (r *EdgeResponse) UnmarshalRLP(input []byte) error {
	return types.UnmarshalRlp(r.unmarshalRLPFrom, input[0:])
}

// unmarshalRLPFrom unmarshals a EdgeResponse in RLP format
func (r *EdgeResponse) unmarshalRLPFrom(_ *fastrlp.Parser, v *fastrlp.Value) error {
	elems, err := v.GetElems()
	if err != nil {
		return err
	}

	if len(elems) < 4 {
		return fmt.Errorf("incorrect number of elements to decode rtcMsg, expected 4 but found %d", len(elems))
	}

	// RespString
	if r.RespString, err = elems[0].GetString(); err != nil {
		return err
	}

	// V
	r.V = new(big.Int)
	if err = elems[1].GetBigInt(r.V); err != nil {
		return err
	}

	// R
	r.R = new(big.Int)
	if err = elems[2].GetBigInt(r.R); err != nil {
		return err
	}

	// S
	r.S = new(big.Int)
	if err = elems[3].GetBigInt(r.S); err != nil {
		return err
	}

	// set From with default value
	r.From = types.ZeroAddress

	// From
	if len(elems) >= 5 {
		if vv, err := v.Get(4).Bytes(); err == nil && len(vv) == types.AddressLength {
			// address
			r.From = types.BytesToAddress(vv)
		}
		if vv, err := v.Get(5).Bytes(); err == nil && len(vv) == types.HashLength {
			// address
			r.Hash = types.BytesToHash(vv)
		}
	}

	return nil
}
