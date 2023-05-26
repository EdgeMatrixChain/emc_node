package idl_test

import (
	"math/big"

	"github.com/emc-protocol/edge-matrix/helper/ic/utils/idl"
)

func ExampleVec() {
	test([]idl.Type{idl.NewVec(new(idl.Int))}, []interface{}{
		[]interface{}{big.NewInt(0), big.NewInt(1), big.NewInt(2), big.NewInt(3)},
	})
	// Output:
	// 4449444c016d7c01000400010203
}
