package types

import (
	"encoding/binary"
	"fmt"
	"sync/atomic"

	"github.com/emc-protocol/edge-matrix/helper/hex"
)

// Header represents a block header in the Ethereum blockchain.
type Header struct {
	ParentHash   Hash
	Miner        []byte
	subject      []byte
	StateRoot    Hash
	TeleRoot     Hash
	ReceiptsRoot Hash
	LogsBloom    Bloom
	Number       uint64
	Timestamp    uint64
	ExtraData    []byte
	Nonce        Nonce
	Hash         Hash

	GasLimit uint64
	GasUsed  uint64
}

func (h *Header) Equal(hh *Header) bool {
	return h.Hash == hh.Hash
}

func (h *Header) HasBody() bool {
	return h.TeleRoot != EmptyRootHash
}

func (h *Header) HasReceipts() bool {
	return h.ReceiptsRoot != EmptyRootHash
}

func (h *Header) SetNonce(i uint64) {
	binary.BigEndian.PutUint64(h.Nonce[:], i)
}

func (h *Header) IsGenesis() bool {
	return h.Hash != ZeroHash && h.Number == 0
}

type Nonce [8]byte

func (n Nonce) String() string {
	return hex.EncodeToHex(n[:])
}

// MarshalText implements encoding.TextMarshaler
func (n Nonce) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

func (h *Header) Copy() *Header {
	newHeader := &Header{
		ParentHash:   h.ParentHash,
		StateRoot:    h.StateRoot,
		TeleRoot:     h.TeleRoot,
		ReceiptsRoot: h.ReceiptsRoot,
		Hash:         h.Hash,
		LogsBloom:    h.LogsBloom,
		Nonce:        h.Nonce,
		Number:       h.Number,
		Timestamp:    h.Timestamp,
		GasLimit:     h.GasLimit,
		GasUsed:      h.GasUsed,
	}

	newHeader.Miner = make([]byte, len(h.Miner))
	copy(newHeader.Miner[:], h.Miner[:])

	newHeader.ExtraData = make([]byte, len(h.ExtraData))
	copy(newHeader.ExtraData[:], h.ExtraData[:])

	return newHeader
}

type Body struct {
	Telegrams []*Telegram
}

type FullBlock struct {
	Block    *Block
	Receipts []*Receipt
}

type Block struct {
	Header    *Header
	Telegrams []*Telegram
	Uncles    []*Header

	// Cache
	size atomic.Value // *uint64
}

func (b *Block) Hash() Hash {
	return b.Header.Hash
}

func (b *Block) Number() uint64 {
	return b.Header.Number
}

func (b *Block) ParentHash() Hash {
	return b.Header.ParentHash
}

func (b *Block) Body() *Body {
	return &Body{
		Telegrams: b.Telegrams,
	}
}

func (b *Block) Size() uint64 {
	sizePtr := b.size.Load()
	if sizePtr == nil {
		bytes := b.MarshalRLP()
		size := uint64(len(bytes))
		b.size.Store(&size)

		return size
	}

	sizeVal, ok := sizePtr.(*uint64)
	if !ok {
		return 0
	}

	return *sizeVal
}

func (b *Block) String() string {
	str := fmt.Sprintf(`Block(#%v):`, b.Number())

	return str
}

// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block) WithSeal(header *Header) *Block {
	cpy := *header

	return &Block{
		Header:    &cpy,
		Telegrams: b.Telegrams,
		Uncles:    b.Uncles,
	}
}
