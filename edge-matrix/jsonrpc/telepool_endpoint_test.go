package jsonrpc

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/emc-protocol/edge-matrix/types"

	"github.com/stretchr/testify/assert"
)

func TestContentEndpoint(t *testing.T) {
	t.Parallel()

	t.Run("returns empty ContentResponse if tx pool has no transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		txPoolEndpoint := &TelePool{mockStore}

		result, _ := txPoolEndpoint.Content()
		//nolint:forcetypeassert
		response := result.(ContentResponse)

		assert.True(t, mockStore.includeQueued)
		assert.Equal(t, 0, len(response.Pending))
		assert.Equal(t, 0, len(response.Queued))
	})

	//nolint:dupl
	t.Run("returns correct data for pending transaction", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		address1 := types.Address{0x1}
		testTx := newTestTransaction(2, address1)
		mockStore.pending[address1] = []*types.Telegram{testTx}
		txPoolEndpoint := &TelePool{mockStore}

		result, _ := txPoolEndpoint.Content()
		//nolint:forcetypeassert
		response := result.(ContentResponse)

		assert.Equal(t, 1, len(response.Pending))
		assert.Equal(t, 0, len(response.Queued))
		assert.Equal(t, 1, len(response.Pending[address1]))

		txData := response.Pending[address1][testTx.Nonce]
		assert.NotNil(t, txData)
		assert.Equal(t, testTx.Gas, uint64(txData.Gas))
		assert.Equal(t, *testTx.GasPrice, big.Int(txData.GasPrice))
		assert.Equal(t, testTx.To, txData.To)
		assert.Equal(t, testTx.From, txData.From)
		assert.Equal(t, *testTx.Value, big.Int(txData.Value))
		assert.Equal(t, testTx.Input, []byte(txData.Input))
		assert.Equal(t, nil, txData.BlockNumber)
		assert.Equal(t, nil, txData.TxIndex)
	})

	//nolint:dupl
	t.Run("returns correct data for queued transaction", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		address1 := types.Address{0x1}
		testTx := newTestTransaction(2, address1)
		mockStore.queued[address1] = []*types.Telegram{testTx}
		txPoolEndpoint := &TelePool{mockStore}

		result, _ := txPoolEndpoint.Content()
		//nolint:forcetypeassert
		response := result.(ContentResponse)

		assert.Equal(t, 0, len(response.Pending))
		assert.Equal(t, 1, len(response.Queued))
		assert.Equal(t, 1, len(response.Queued[address1]))

		txData := response.Queued[address1][testTx.Nonce]
		assert.NotNil(t, txData)
		assert.Equal(t, testTx.Gas, uint64(txData.Gas))
		assert.Equal(t, *testTx.GasPrice, big.Int(txData.GasPrice))
		assert.Equal(t, testTx.To, txData.To)
		assert.Equal(t, testTx.From, txData.From)
		assert.Equal(t, *testTx.Value, big.Int(txData.Value))
		assert.Equal(t, testTx.Input, []byte(txData.Input))
		assert.Equal(t, nil, txData.BlockNumber)
		assert.Equal(t, nil, txData.TxIndex)
	})

	t.Run("returns correct ContentResponse data for multiple transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		address1 := types.Address{0x1}
		testTx1 := newTestTransaction(2, address1)
		testTx2 := newTestTransaction(4, address1)
		testTx3 := newTestTransaction(11, address1)
		address2 := types.Address{0x2}
		testTx4 := newTestTransaction(7, address2)
		testTx5 := newTestTransaction(8, address2)
		mockStore.pending[address1] = []*types.Telegram{testTx1, testTx2}
		mockStore.pending[address2] = []*types.Telegram{testTx4}
		mockStore.queued[address1] = []*types.Telegram{testTx3}
		mockStore.queued[address2] = []*types.Telegram{testTx5}
		txPoolEndpoint := &TelePool{mockStore}

		result, _ := txPoolEndpoint.Content()
		//nolint:forcetypeassert
		response := result.(ContentResponse)

		assert.True(t, mockStore.includeQueued)
		assert.Equal(t, 2, len(response.Pending))
		assert.Equal(t, 2, len(response.Pending[address1]))
		assert.Equal(t, 1, len(response.Pending[address2]))
		assert.Equal(t, 2, len(response.Queued))
	})
}

func TestInspectEndpoint(t *testing.T) {
	t.Parallel()

	t.Run("returns empty InspectResponse if tx pool has no transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		mockStore.maxSlots = 1024
		txPoolEndpoint := &TelePool{mockStore}

		result, _ := txPoolEndpoint.Inspect()
		//nolint:forcetypeassert
		response := result.(InspectResponse)

		assert.True(t, mockStore.includeQueued)
		assert.Equal(t, 0, len(response.Pending))
		assert.Equal(t, 0, len(response.Queued))
		assert.Equal(t, uint64(0), response.CurrentCapacity)
		assert.Equal(t, mockStore.maxSlots, response.MaxCapacity)
	})

	t.Run("returns correct data for queued transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		mockStore.capacity = 1
		address1 := types.Address{0x1}
		testTx := newTestTransaction(2, address1)
		mockStore.queued[address1] = []*types.Telegram{testTx}
		txPoolEndpoint := &TelePool{mockStore}

		result, _ := txPoolEndpoint.Inspect()
		//nolint:forcetypeassert
		response := result.(InspectResponse)

		assert.Equal(t, 0, len(response.Pending))
		assert.Equal(t, 1, len(response.Queued))
		assert.Equal(t, uint64(1), response.CurrentCapacity)
		transactionInfo := response.Queued[testTx.From.String()]
		assert.NotNil(t, transactionInfo)
		assert.NotNil(t, transactionInfo[strconv.FormatUint(testTx.Nonce, 10)])
	})

	t.Run("returns correct data for pending transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		mockStore.capacity = 2
		address1 := types.Address{0x1}
		testTx := newTestTransaction(2, address1)
		testTx2 := newTestTransaction(3, address1)
		mockStore.pending[address1] = []*types.Telegram{testTx, testTx2}
		txPoolEndpoint := &TelePool{mockStore}

		result, _ := txPoolEndpoint.Inspect()
		//nolint:forcetypeassert
		response := result.(InspectResponse)

		assert.Equal(t, 1, len(response.Pending))
		assert.Equal(t, 0, len(response.Queued))
		assert.Equal(t, uint64(2), response.CurrentCapacity)
		transactionInfo := response.Pending[testTx.From.String()]
		assert.NotNil(t, transactionInfo)
		assert.NotNil(t, transactionInfo[strconv.FormatUint(testTx.Nonce, 10)])
		assert.NotNil(t, transactionInfo[strconv.FormatUint(testTx2.Nonce, 10)])
	})
}

func TestStatusEndpoint(t *testing.T) {
	t.Parallel()

	t.Run("returns empty StatusResponse if tx pool has no transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		txPoolEndpoint := &TelePool{mockStore}

		result, _ := txPoolEndpoint.Status()
		//nolint:forcetypeassert
		response := result.(StatusResponse)

		assert.Equal(t, uint64(0), response.Pending)
		assert.Equal(t, uint64(0), response.Queued)
	})

	t.Run("returns correct count of pending/queued transactions", func(t *testing.T) {
		t.Parallel()

		mockStore := newMockTxPoolStore()
		address1 := types.Address{0x1}
		testTx1 := newTestTransaction(2, address1)
		testTx2 := newTestTransaction(4, address1)
		testTx3 := newTestTransaction(11, address1)
		address2 := types.Address{0x2}
		testTx4 := newTestTransaction(7, address2)
		testTx5 := newTestTransaction(8, address2)
		mockStore.pending[address1] = []*types.Telegram{testTx1, testTx2}
		mockStore.pending[address2] = []*types.Telegram{testTx4}
		mockStore.queued[address1] = []*types.Telegram{testTx3}
		mockStore.queued[address2] = []*types.Telegram{testTx5}
		txPoolEndpoint := &TelePool{mockStore}

		result, _ := txPoolEndpoint.Status()
		//nolint:forcetypeassert
		response := result.(StatusResponse)

		assert.Equal(t, uint64(3), response.Pending)
		assert.Equal(t, uint64(2), response.Queued)
	})
}

type mockTxPoolStore struct {
	pending       map[types.Address][]*types.Telegram
	queued        map[types.Address][]*types.Telegram
	capacity      uint64
	maxSlots      uint64
	includeQueued bool
}

func newMockTxPoolStore() *mockTxPoolStore {
	return &mockTxPoolStore{
		pending: make(map[types.Address][]*types.Telegram),
		queued:  make(map[types.Address][]*types.Telegram),
	}
}

func (s *mockTxPoolStore) GetTeles(inclQueued bool) (map[types.Address][]*types.Telegram, map[types.Address][]*types.Telegram) {
	s.includeQueued = inclQueued

	return s.pending, s.queued
}

func (s *mockTxPoolStore) GetCapacity() (uint64, uint64) {
	return s.capacity, s.maxSlots
}

func newTestTransaction(nonce uint64, from types.Address) *types.Telegram {
	txn := &types.Telegram{
		Nonce:    nonce,
		GasPrice: big.NewInt(1),
		Gas:      nonce * 100,
		Value:    big.NewInt(200),
		Input:    []byte{0xff},
		From:     from,
		To:       &addr1,
		V:        big.NewInt(1),
		R:        big.NewInt(1),
		S:        big.NewInt(1),
	}

	txn.ComputeHash()

	return txn
}
