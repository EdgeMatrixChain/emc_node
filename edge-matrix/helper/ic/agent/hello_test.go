package agent

import (
	"testing"

	"github.com/emc-protocol/edge-matrix/helper/ic/utils/idl"
)

func TestHello_QueryRaw(t *testing.T) {
	//EXT canister
	//canisterID := "bzsui-sqaaa-aaaah-qce2a-cai"

	//PUNK canister
	// canisterID := "qfh5c-6aaaa-aaaah-qakeq-cai"

	//agent := New(true, "")
	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")

	canister := "xb3xh-uaaaa-aaaam-abi3a-cai"
	methodName := "greet"

	var argType []idl.Type
	argType = append(argType, new(idl.Text))

	var argValue []interface{}
	argValue = append(argValue, "world")

	arg, _ := idl.Encode(argType, argValue)
	_, result, str, err := agent.Query(canister, methodName, arg)
	if err != nil {
		panic(err)
	}
	t.Log(str, " ->", result[0])
}

//func TestAgent_UpdateRaw(t *testing.T) {
//	// canisterID := "gvbup-jyaaa-aaaah-qcdwa-cai"
//	// agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
//
//	// methodName := "transfer"
//	// var argType []idl.Type
//	// var argValue []interface{}
//	// p, _ := principal.Decode("aaaaa-aa")
//	// argType = append(argType, new(idl.Principal))
//	// argType = append(argType, new(idl.Nat))
//	// argValue = append(argValue, p)
//	// argValue = append(argValue, big.NewInt(10000000000))
//	var myresult uint64
//	canisterID := "d24m2-dqaaa-aaaah-aa4zq-cai"
//	ag := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
//	methodName := "total"
//	arg, _ := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
//	_, result, err := ag.Update(canisterID, methodName, arg, 30)
//	if err != nil {
//		panic(err)
//	}
//	utils.Decode(&myresult, result[0])
//
//	// arg, _ := idl.Encode(argType, argValue)
//	// _, result, err := agent.Update(canisterID, methodName, arg)
//
//	t.Log("errMsg:", err, "result:", myresult)
//}
//
//func TestAgent_GetCanisterModule(t *testing.T) {
//	canisterID := "bzsui-sqaaa-aaaah-qce2a-cai"
//	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
//	result, err := agent.GetCanisterModule(canisterID)
//	if err != nil {
//		t.Log("err:", err)
//	} else if result == nil {
//		t.Log("no module")
//	} else {
//		t.Log("hash:", hex.EncodeToString(result))
//	}
//}
//
//func TestAgent_GetCanisterControllers(t *testing.T) {
//	canisterID := "6b4pv-sqaaa-aaaah-qaava-cai"
//	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
//	result, err := agent.GetCanisterControllers(canisterID)
//	if err != nil {
//		t.Log("err:", err)
//	} else {
//		for _, i := range result {
//			t.Log("controller:", i.Encode())
//		}
//	}
//	t.Log(result)
//}
//
//func TestPrincipal(t *testing.T) {
//	pkBytes, _ := hex.DecodeString("833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
//	identity := identity.New(false, pkBytes)
//	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
//	t.Log(p.Encode(), len(identity.PubKeyBytes()))
//}
//
//func TestCbor(t *testing.T) {
//	canisterID, _ := principal.Decode("gvbup-jyaaa-aaaah-qcdwa-cai")
//	agent := New(true, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
//
//	req := Request{
//		Type:          "call",
//		Sender:        agent.Sender(),
//		IngressExpiry: uint64(agent.getExpiryDate().UnixNano()),
//		CanisterID:    canisterID,
//		MethodName:    "transfer",
//		Arguments:     []byte("i love vivian"),
//	}
//
//	envelope := Envelope{
//		req,
//		[]byte{},
//		[]byte{},
//	}
//
//	data, _ := cbor.Marshal(envelope)
//	resp := new(Envelope)
//	cbor.Unmarshal(data, resp)
//	t.Log("sender", resp.Content.Sender.Encode())
//	t.Log("type", resp.Content.Type)
//	t.Log("ingress expiryt", resp.Content.IngressExpiry)
//	t.Log("method", resp.Content.MethodName)
//	t.Log("arg", resp.Content.Arguments)
//	t.Log("canister", resp.Content.CanisterID.Encode())
//}
//
//func TestAgent_GetCanisterTime(t *testing.T) {
//	canisterID := "b65vx-3qaaa-aaaaa-7777q-cai"
//	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
//	result, err := agent.GetCanisterTime(canisterID)
//	if err != nil {
//		t.Log("err:", err)
//	} else {
//		t.Log("result:", result)
//
//	}
//}
//
//func TestAgent_GetCanisterCandid(t *testing.T) {
//	canisterID := "oeee4-qaaaa-aaaak-qaaeq-cai"
//	agent := New(false, "833fe62409237b9d62ec77587520911e9a759cec1d19755b7da901b96dca3d42")
//	arg, _ := idl.Encode([]idl.Type{new(idl.Null)}, []interface{}{nil})
//	methodName := "__get_candid_interface_tmp_hack"
//	_, result, err := agent.Update(canisterID, methodName, arg, 30)
//	if err != nil {
//		panic(err)
//	}
//	if err != nil {
//		t.Log("err:", err)
//	} else {
//		t.Log("result:", result)
//
//	}
//}
