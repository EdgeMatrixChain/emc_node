package idl_test

import "github.com/emc-protocol/edge-matrix/helper/ic/utils/idl"

func ExampleText() {
	test([]idl.Type{new(idl.Text)}, []interface{}{""})
	test([]idl.Type{new(idl.Text)}, []interface{}{"Motoko"})
	test([]idl.Type{new(idl.Text)}, []interface{}{"Hi ☃\n"})
	// Output:
	// 4449444c00017100
	// 4449444c000171064d6f746f6b6f
	// 4449444c00017107486920e298830a
}
