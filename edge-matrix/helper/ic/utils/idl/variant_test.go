package idl_test

import (
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/idl"
	"testing"
)

func Test_ExampleVariant(t *testing.T) {
	result := map[string]idl.Type{
		"ok":  new(idl.Text),
		"err": new(idl.Text),
	}
	test_([]idl.Type{idl.NewVariant(result)}, []interface{}{idl.FieldValue{
		Name:  "ok",
		Value: "good",
	}})
	test_([]idl.Type{idl.NewVariant(result)}, []interface{}{idl.FieldValue{
		Name:  "err",
		Value: "uhoh",
	}})
	// Output:
	// 4449444c016b029cc20171e58eb4027101000004676f6f64
	// 4449444c016b029cc20171e58eb402710100010475686f68
}
