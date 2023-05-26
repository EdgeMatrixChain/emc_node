package consensus

import (
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/emc-protocol/edge-matrix/types/buildroot"
)

// BuildBlockParams are parameters passed into the BuildBlock helper method
type BuildBlockParams struct {
	Header   *types.Header
	Teles    []*types.Telegram
	Receipts []*types.Receipt
}

// BuildBlock is a utility function that builds a block, based on the passed in header, transactions and receipts
func BuildBlock(params BuildBlockParams) *types.Block {
	teles := params.Teles
	header := params.Header

	if len(teles) == 0 {
		header.TeleRoot = types.EmptyRootHash
	} else {
		header.TeleRoot = buildroot.CalculateTelegramsRoot(teles)
	}

	if len(params.Receipts) == 0 {
		header.ReceiptsRoot = types.EmptyRootHash
	} else {
		header.ReceiptsRoot = buildroot.CalculateReceiptsRoot(params.Receipts)
	}

	// TODO: Compute uncles
	//header.Sha3Uncles = types.EmptyUncleHash
	header.ComputeHash()

	return &types.Block{
		Header:    header,
		Telegrams: teles,
	}
}
