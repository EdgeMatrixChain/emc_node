package idl_test

import (
	"math/big"

	"github.com/emc-protocol/edge-matrix/helper/ic/utils/idl"
)

func ExampleOpt() {
	test([]idl.Type{idl.NewOpt(new(idl.Nat))}, []interface{}{nil})
	test([]idl.Type{idl.NewOpt(new(idl.Nat))}, []interface{}{big.NewInt(1)})
	// Output:
	// 4449444c016e7d010000
	// 4449444c016e7d01000101
}
