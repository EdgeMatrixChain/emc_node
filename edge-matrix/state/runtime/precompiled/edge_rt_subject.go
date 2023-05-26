package precompiled

import (
	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/emc-protocol/edge-matrix/state/runtime"
	"github.com/emc-protocol/edge-matrix/types"
)

type edgeRtcSubject struct{}

func (c *edgeRtcSubject) gas(input []byte, _ *chain.ForksInTime) uint64 {
	return 0
}

func (c *edgeRtcSubject) run(input []byte, caller types.Address, host runtime.Host) ([]byte, error) {
	//if len(input) < 1 {
	//	return abiBoolFalse, runtime.ErrInvalidInputData
	//}

	// check if caller is native token contract
	//if caller != contracts.NativeTokenContract {
	//	return abiBoolFalse, runtime.ErrUnauthorizedCaller
	//}
	//
	//from := types.BytesToAddress(input[0:32])
	//to := types.BytesToAddress(input[32:64])
	//amount := new(big.Int).SetBytes(input[64:96])
	//
	//// state changes
	//if err := host.Transfer(from, to, amount); err != nil {
	//	return abiBoolFalse, err
	//}

	return abiBoolTrue, nil
}
