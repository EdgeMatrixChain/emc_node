package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emc-protocol/edge-matrix/contracts"
	"github.com/emc-protocol/edge-matrix/rtc"
	"github.com/hashicorp/go-hclog"
	"github.com/umbracle/fastrlp"
	"math/big"

	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/emc-protocol/edge-matrix/helper/common"
	"github.com/emc-protocol/edge-matrix/helper/progress"
	"github.com/emc-protocol/edge-matrix/state/runtime"
	"github.com/emc-protocol/edge-matrix/types"
)

type edgeTelePoolStore interface {
	// AddTele adds a new telegram to the telegram pool
	AddTele(tx *types.Telegram) (string, error)

	// GetPendingTx gets the pending transaction from the transaction pool, if it's present
	GetPendingTele(txHash types.Hash) (*types.Telegram, bool)

	// GetNonce returns the next nonce for this address
	GetNonce(addr types.Address) uint64
}

type edgeRtcStore interface {
	SendMsg(msg *rtc.RtcMsg) error
	Sender(msg *rtc.RtcMsg) (types.Address, error)
}

type Account struct {
	Balance *big.Int
	Nonce   uint64
}

type ethStateStore interface {
	GetAccount(root types.Hash, addr types.Address) (*Account, error)
	GetStorage(root types.Hash, addr types.Address, slot types.Hash) ([]byte, error)
	GetForksInTime(blockNumber uint64) chain.ForksInTime
	GetCode(root types.Hash, addr types.Address) ([]byte, error)
}

type ethBlockchainStore interface {
	// Header returns the current msg of the chain (genesis if empty)
	Header() *types.Header

	// GetHeaderByNumber gets a msg using the provided number
	GetHeaderByNumber(uint64) (*types.Header, bool)

	// GetBlockByHash gets a block using the provided hash
	GetBlockByHash(hash types.Hash, full bool) (*types.Block, bool)

	// GetBlockByNumber returns a block using the provided number
	GetBlockByNumber(num uint64, full bool) (*types.Block, bool)

	// ReadTxLookup returns a block hash in which a given txn was mined
	ReadTxLookup(txnHash types.Hash) (types.Hash, bool)

	// GetReceiptsByHash returns the receipts for a block hash
	GetReceiptsByHash(hash types.Hash) ([]*types.Receipt, error)

	// GetAvgGasPrice returns the average gas price
	GetAvgGasPrice() *big.Int

	// ApplyTxn applies a transaction object to the blockchain
	ApplyTxn(header *types.Header, txn *types.Telegram) (*runtime.ExecutionResult, error)

	// GetSyncProgression retrieves the current sync progression, if any
	GetSyncProgression() *progress.Progression
}

// edgeStore provides access to the methods needed by edge endpoint
type edgeStore interface {
	edgeTelePoolStore
	edgeRtcStore
	ethStateStore
	ethBlockchainStore
}

// Edge is the edge jsonrpc endpoint
type Edge struct {
	logger        hclog.Logger
	store         edgeStore
	chainID       uint64
	filterManager *FilterManager
	priceLimit    uint64
}

var (
	ErrInsufficientFunds = errors.New("insufficient funds for execution")
)

func (e *Edge) NewWallet() (interface{}, error) {
	return nil, nil
}

// ChainId returns the chain id of the client
//
//nolint:stylecheck
func (e *Edge) ChainId() (interface{}, error) {
	return argUintPtr(e.chainID), nil
}

func (e *Edge) Syncing() (interface{}, error) {
	if syncProgression := e.store.GetSyncProgression(); syncProgression != nil {
		// Node is bulk syncing, return the status
		return progression{
			Type:          string(syncProgression.SyncType),
			StartingBlock: argUint64(syncProgression.StartingBlock),
			CurrentBlock:  argUint64(syncProgression.CurrentBlock),
			HighestBlock:  argUint64(syncProgression.HighestBlock),
		}, nil
	}

	// Node is not bulk syncing
	return false, nil
}

// GetBlockByNumber returns information about a block by block number
func (e *Edge) GetBlockByNumber(number BlockNumber, fullTx bool) (interface{}, error) {
	num, err := GetNumericBlockNumber(number, e.store)
	if err != nil {
		return nil, err
	}

	block, ok := e.store.GetBlockByNumber(num, true)
	if !ok {
		return nil, nil
	}

	return toBlock(block, fullTx), nil
}

// GetBlockByHash returns information about a block by hash
func (e *Edge) GetBlockByHash(hash types.Hash, fullTx bool) (interface{}, error) {
	block, ok := e.store.GetBlockByHash(hash, true)
	if !ok {
		return nil, nil
	}

	return toBlock(block, fullTx), nil
}

func (e *Edge) GetBlockTelegramCountByNumber(number BlockNumber) (interface{}, error) {
	num, err := GetNumericBlockNumber(number, e.store)
	if err != nil {
		return nil, err
	}

	block, ok := e.store.GetBlockByNumber(num, true)

	if !ok {
		return nil, nil
	}

	return len(block.Telegrams), nil
}

// BlockNumber returns current block number
func (e *Edge) BlockNumber() (interface{}, error) {
	h := e.store.Header()
	if h == nil {
		return nil, fmt.Errorf("msg has a nil value")
	}

	return argUintPtr(h.Number), nil
}

// SendRawTelegram sends a raw telegram
func (e *Edge) SendRawTelegram(buf argBytes) (interface{}, error) {
	tele := &types.Telegram{}
	if err := tele.UnmarshalRLP(buf); err != nil {
		return nil, err
	}
	e.logger.Debug(fmt.Sprintf("SendRawTelegram To: %s, Nonce:%d", tele.To.String(), tele.Nonce))
	tele.ComputeHash()

	teleResp, teleErr := e.store.AddTele(tele)
	if teleErr != nil {
		return nil, teleErr
	}
	resp := fmt.Sprintf(`{"telegram_hash":"%s","response":"%s"}`, tele.Hash.String(), teleResp)
	return resp, nil
}

func (e *Edge) SendRawMsg(buf argBytes) (interface{}, error) {
	msg := &rtc.RtcMsg{}
	if err := msg.UnmarshalRLP(buf); err != nil {
		return nil, err
	}
	if e.logger.IsDebug() {
		marshal, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}
		e.logger.Debug(fmt.Sprintf("SendRawMsg: %s", string(marshal)))
	}
	msg.ComputeHash()

	if err := e.store.SendMsg(msg); err != nil {
		return nil, err
	}

	return msg.Hash.String(), nil
}

// GetTelegramByHash returns a telegram by its hash.
// If the telegram is still pending -> return the telegram with some fields omitted
// If the telegram is sealed into a block -> return the whole telegram with all fields
func (e *Edge) GetTelegramByHash(hash types.Hash) (interface{}, error) {
	// findSealedTx is a helper method for checking the world state
	// for the transaction with the provided hash
	findSealedTx := func() *transaction {
		// Check the chain state for the transaction
		blockHash, ok := e.store.ReadTxLookup(hash)
		if !ok {
			// Block not found in storage
			return nil
		}

		block, ok := e.store.GetBlockByHash(blockHash, true)

		if !ok {
			// Block receipts not found in storage
			return nil
		}

		// Find the transaction within the block
		for idx, telegram := range block.Telegrams {
			if telegram.Hash == hash {
				return toTransaction(
					telegram,
					argUintPtr(block.Number()),
					argHashPtr(block.Hash()),
					&idx,
				)
			}
		}

		return nil
	}

	// findPendingTx is a helper method for checking the TxPool
	// for the pending transaction with the provided hash
	findPendingTx := func() *transaction {
		// Check the TxPool for the transaction if it's pending
		if pendingTx, pendingFound := e.store.GetPendingTele(hash); pendingFound {
			return toPendingTransaction(pendingTx)
		}

		// Transaction not found in the TxPool
		return nil
	}

	// 1. Check the chain state for the txn
	if resultTxn := findSealedTx(); resultTxn != nil {
		return resultTxn, nil
	}

	// 2. Check the TxPool for the txn
	if resultTxn := findPendingTx(); resultTxn != nil {
		return resultTxn, nil
	}

	// Transaction not found in state or TxPool
	e.logger.Warn(
		fmt.Sprintf("Transaction with hash [%s] not found", hash),
	)

	return nil, nil
}

// GetTelegramReceipt returns a telegram receipt by his hash
func (e *Edge) GetTelegramReceipt(hash types.Hash) (interface{}, error) {
	blockHash, ok := e.store.ReadTxLookup(hash)
	if !ok {
		// txn not found
		return nil, nil
	}

	block, ok := e.store.GetBlockByHash(blockHash, true)
	if !ok {
		// block not found
		e.logger.Warn(
			fmt.Sprintf("Block with hash [%s] not found", blockHash.String()),
		)

		return nil, nil
	}

	receipts, err := e.store.GetReceiptsByHash(blockHash)
	if err != nil {
		// block receipts not found
		e.logger.Warn(
			fmt.Sprintf("Receipts for block with hash [%s] not found", blockHash.String()),
		)

		return nil, nil
	}

	if len(receipts) == 0 {
		// Receipts not written yet on the db
		e.logger.Warn(
			fmt.Sprintf("No receipts found for block with hash [%s]", blockHash.String()),
		)

		return nil, nil
	}
	// find the transaction in the body
	indx := -1

	for i, txn := range block.Telegrams {
		if txn.Hash == hash {
			indx = i

			break
		}
	}

	if indx == -1 {
		// txn not found
		return nil, nil
	}

	txn := block.Telegrams[indx]
	raw := receipts[indx]

	logs := make([]*Log, len(raw.Logs))
	for indx, elem := range raw.Logs {
		logs[indx] = &Log{
			Address:     elem.Address,
			Topics:      elem.Topics,
			Data:        argBytes(elem.Data),
			BlockHash:   block.Hash(),
			BlockNumber: argUint64(block.Number()),
			TxHash:      txn.Hash,
			TxIndex:     argUint64(indx),
			LogIndex:    argUint64(indx),
			Removed:     false,
		}
	}

	res := &receipt{
		Root:               raw.Root,
		CumulativeGasUsed:  argUint64(raw.CumulativeGasUsed),
		LogsBloom:          raw.LogsBloom,
		Status:             argUint64(*raw.Status),
		TxHash:             txn.Hash,
		TxIndex:            argUint64(indx),
		BlockHash:          block.Hash(),
		BlockNumber:        argUint64(block.Number()),
		GasUsed:            argUint64(raw.GasUsed),
		ApplicationAddress: raw.ApplicationAddress,
		FromAddr:           txn.From,
		ToAddr:             txn.To,
		Logs:               logs,
	}

	return res, nil
}

// GetStorageAt returns the contract storage at the index position
func (e *Edge) GetStorageAt(
	address types.Address,
	index types.Hash,
	filter BlockNumberOrHash,
) (interface{}, error) {
	header, err := GetHeaderFromBlockNumberOrHash(filter, e.store)
	if err != nil {
		return nil, err
	}

	// Get the storage for the passed in location
	result, err := e.store.GetStorage(header.StateRoot, address, index)
	if err != nil {
		if errors.Is(err, ErrStateNotFound) {
			return argBytesPtr(types.ZeroHash[:]), nil
		}

		return nil, err
	}

	// TODO: GetStorage should return the values already parsed

	// Parse the RLP value
	p := &fastrlp.Parser{}

	v, err := p.Parse(result)
	if err != nil {
		return argBytesPtr(types.ZeroHash[:]), nil
	}

	data, err := v.Bytes()
	if err != nil {
		return argBytesPtr(types.ZeroHash[:]), nil
	}

	// Pad to return 32 bytes data
	return argBytesPtr(types.BytesToHash(data).Bytes()), nil
}

// GasPrice returns the average gas price based on the last x blocks
// taking into consideration operator defined price limit
func (e *Edge) GasPrice() (interface{}, error) {
	// Fetch average gas price in uint64
	avgGasPrice := e.store.GetAvgGasPrice().Uint64()

	// Return --price-limit flag defined value if it is greater than avgGasPrice
	return argUint64(common.Max(e.priceLimit, avgGasPrice)), nil
}

//// Call executes a smart contract call using the transaction object data
//func (e *Edge) Call(arg *txnArgs, filter BlockNumberOrHash) (interface{}, error) {
//	header, err := GetHeaderFromBlockNumberOrHash(filter, e.store)
//	if err != nil {
//		return nil, err
//	}
//
//	transaction, err := DecodeTxn(arg, e.store)
//	if err != nil {
//		return nil, err
//	}
//	// If the caller didn't supply the gas limit in the message, then we set it to maximum possible => block gas limit
//	if transaction.Gas == 0 {
//		transaction.Gas = header.GasLimit
//	}
//
//	// The return value of the execution is saved in the transition (returnValue field)
//	result, err := e.store.ApplyTxn(header, transaction)
//	if err != nil {
//		return nil, err
//	}
//
//	// Check if an EVM revert happened
//	if result.Reverted() {
//		return nil, constructErrorFromRevert(result)
//	}
//
//	if result.Failed() {
//		return nil, fmt.Errorf("unable to execute call: %w", result.Err)
//	}
//
//	return argBytesPtr(result.ReturnValue), nil
//}

// GetFilterLogs returns an array of logs for the specified filter
func (e *Edge) GetFilterLogs(id string) (interface{}, error) {
	logFilter, err := e.filterManager.GetLogFilterFromID(id)
	if err != nil {
		return nil, err
	}

	return e.filterManager.GetLogsForQuery(logFilter.query)
}

// GetLogs returns an array of logs matching the filter options
func (e *Edge) GetLogs(query *LogQuery) (interface{}, error) {
	return e.filterManager.GetLogsForQuery(query)
}

// GetBalance returns the account's balance at the referenced block.
func (e *Edge) GetBalance(address types.Address, filter BlockNumberOrHash) (interface{}, error) {
	header, err := GetHeaderFromBlockNumberOrHash(filter, e.store)
	if err != nil {
		return nil, err
	}

	// Extract the account balance
	acc, err := e.store.GetAccount(header.StateRoot, address)
	if errors.Is(err, ErrStateNotFound) {
		// Account not found, return an empty account
		return argUintPtr(0), nil
	} else if err != nil {
		return nil, err
	}

	return argBigPtr(acc.Balance), nil
}

// GetTelegramCount returns account nonce
func (e *Edge) GetTelegramCount(address types.Address, filter BlockNumberOrHash) (interface{}, error) {
	var (
		blockNumber BlockNumber
		header      *types.Header
		err         error
	)

	// The filter is empty, use the latest block by default
	if filter.BlockNumber == nil && filter.BlockHash == nil {
		filter.BlockNumber, _ = createBlockNumberPointer("latest")
	}

	if filter.BlockNumber == nil {
		header, err = GetHeaderFromBlockNumberOrHash(filter, e.store)
		if err != nil {
			return nil, fmt.Errorf("failed to get msg from block hash or block number: %w", err)
		}

		blockNumber = BlockNumber(header.Number)
	} else {
		blockNumber = *filter.BlockNumber
	}

	nonce, err := GetNextNonce(address, blockNumber, e.store)
	if err != nil {
		if errors.Is(err, ErrStateNotFound) {
			return argUintPtr(0), nil
		}

		return nil, err
	}

	return argUintPtr(nonce), nil
}

func (e *Edge) Sender(msg *RtcMsg) (string, error) {
	bigV := new(big.Int)
	bigV.SetString(msg.V, 0)
	bigR := new(big.Int)
	bigR.SetString(msg.R, 0)
	bigS := new(big.Int)
	bigS.SetString(msg.S, 0)
	rtcMsg := &rtc.RtcMsg{
		//Nonce:       msg.Nonce,
		Subject:     msg.Subject,
		Application: msg.Application,
		Content:     msg.Content,
		V:           bigV,
		R:           bigR,
		S:           bigS,
		Type:        msg.Type,
	}
	sender, err := e.store.Sender(rtcMsg)
	if err != nil {
		return "", err
	}
	return sender.String(), nil
}

func (e *Edge) SendMsg(msg *RtcMsg) (interface{}, error) {
	hash := types.StringToHash(msg.Subject)
	obj, err := e.GetTelegramByHash(hash)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, fmt.Errorf("failed to send msg to subject: %s is not exist", msg.Subject)
	}
	telegram := obj.(*transaction)
	if *telegram.To != contracts.EdgeRtcSubjectPrecompile {
		return nil, fmt.Errorf("failed to send msg to subject: %s is not a valid subject hash", msg.Subject)
	}

	bigV := new(big.Int)
	bigV.SetString(msg.V, 0)
	bigR := new(big.Int)
	bigR.SetString(msg.R, 0)
	bigS := new(big.Int)
	bigS.SetString(msg.S, 0)
	rtcMsg := &rtc.RtcMsg{
		To:          types.StringToAddress(msg.To),
		Subject:     msg.Subject,
		Application: msg.Application,
		Content:     msg.Content,
		V:           bigV,
		R:           bigR,
		S:           bigS,
		Type:        msg.Type,
	}

	err = e.store.SendMsg(rtcMsg)
	if err != nil {
		return nil, err
	}
	return "0x1", nil
}

// GetCode returns account code at given block number
func (e *Edge) GetCode(address types.Address, filter BlockNumberOrHash) (interface{}, error) {
	header, err := GetHeaderFromBlockNumberOrHash(filter, e.store)
	if err != nil {
		return nil, err
	}

	emptySlice := []byte{}
	code, err := e.store.GetCode(header.StateRoot, address)

	if errors.Is(err, ErrStateNotFound) {
		// If the account doesn't exist / is not initialized yet,
		// return the default value
		return "0x", nil
	} else if err != nil {
		return argBytesPtr(emptySlice), err
	}

	return argBytesPtr(code), nil
}

// NewFilter creates a filter object, based on filter options, to notify when the state changes (logs).
func (e *Edge) NewFilter(filter *LogQuery) (interface{}, error) {
	return e.filterManager.NewLogFilter(filter, nil), nil
}

// NewBlockFilter creates a filter in the node, to notify when a new block arrives
func (e *Edge) NewBlockFilter() (interface{}, error) {
	return e.filterManager.NewBlockFilter(nil), nil
}

// GetFilterChanges is a polling method for a filter, which returns an array of logs which occurred since last poll.
func (e *Edge) GetFilterChanges(id string) (interface{}, error) {
	return e.filterManager.GetFilterChanges(id)
}

// UninstallFilter uninstalls a filter with given ID
func (e *Edge) UninstallFilter(id string) (bool, error) {
	return e.filterManager.Uninstall(id), nil
}

// Unsubscribe uninstalls a filter in a websocket
func (e *Edge) Unsubscribe(id string) (bool, error) {
	return e.filterManager.Uninstall(id), nil
}
