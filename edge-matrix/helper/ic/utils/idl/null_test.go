package idl_test

import "github.com/emc-protocol/edge-matrix/helper/ic/utils/idl"

func ExampleNull() {
	test([]idl.Type{new(idl.Null)}, []interface{}{nil})
	// Output:
	// 4449444c00017f
}
