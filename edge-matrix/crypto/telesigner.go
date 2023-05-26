package crypto

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/bits"

	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/emc-protocol/edge-matrix/helper/keccak"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/umbracle/fastrlp"
)

// TxSigner is a utility interface used to recover data from a transaction
type TxSigner interface {
	// Hash returns the hash of the transaction
	Hash(tx *types.Telegram) types.Hash

	// Sender returns the sender of the transaction
	Sender(tx *types.Telegram) (types.Address, error)

	// SignTx signs a transaction
	SignTele(tx *types.Telegram, priv *ecdsa.PrivateKey) (*types.Telegram, error)

	// CalculateV calculates the V value based on the type of signer used
	CalculateV(parity byte) []byte
}

// NewSigner creates a new signer object (EIP155 or FrontierSigner)
func NewSigner(forks chain.ForksInTime, chainID uint64) TxSigner {
	var signer TxSigner

	if forks.EIP155 {
		signer = &EIP155Signer{chainID: chainID, isHomestead: forks.Homestead}
	} else {
		signer = &FrontierSigner{forks.Homestead}
	}

	return signer
}

type FrontierSigner struct {
	isHomestead bool
}

var signerPool fastrlp.ArenaPool

// calcTeleHash calculates the transaction hash (keccak256 hash of the RLP value)
func calcTeleHash(tele *types.Telegram, chainID uint64) types.Hash {
	a := signerPool.Get()

	v := a.NewArray()
	v.Set(a.NewUint(tele.Nonce))
	v.Set(a.NewBigInt(tele.GasPrice))
	v.Set(a.NewUint(tele.Gas))

	if tele.To == nil {
		v.Set(a.NewNull())
	} else {
		v.Set(a.NewCopyBytes((*tele.To).Bytes()))
	}

	v.Set(a.NewBigInt(tele.Value))
	v.Set(a.NewCopyBytes(tele.Input))

	// EIP155
	if chainID != 0 {
		v.Set(a.NewUint(chainID))
		v.Set(a.NewUint(0))
		v.Set(a.NewUint(0))
	}

	hash := keccak.Keccak256Rlp(nil, v)

	signerPool.Put(a)

	return types.BytesToHash(hash)
}

// Hash is a wrapper function for the calcTeleHash, with chainID 0
func (f *FrontierSigner) Hash(tx *types.Telegram) types.Hash {
	return calcTeleHash(tx, 0)
}

// Magic numbers from Ethereum, used in v calculation
var (
	big27 = big.NewInt(27)
	big35 = big.NewInt(35)
)

// Sender decodes the signature and returns the sender of the transaction
func (f *FrontierSigner) Sender(tx *types.Telegram) (types.Address, error) {
	refV := big.NewInt(0)
	if tx.V != nil {
		refV.SetBytes(tx.V.Bytes())
	}

	refV.Sub(refV, big27)

	sig, err := encodeSignature(tx.R, tx.S, refV, f.isHomestead)
	if err != nil {
		return types.Address{}, err
	}

	pub, err := Ecrecover(f.Hash(tx).Bytes(), sig)
	if err != nil {
		return types.Address{}, err
	}

	buf := Keccak256(pub[1:])[12:]

	return types.BytesToAddress(buf), nil
}

// SignTx signs the transaction using the passed in private key
func (f *FrontierSigner) SignTele(
	tx *types.Telegram,
	privateKey *ecdsa.PrivateKey,
) (*types.Telegram, error) {
	tx = tx.Copy()

	h := f.Hash(tx)

	sig, err := Sign(privateKey, h[:])
	if err != nil {
		return nil, err
	}

	tx.R = new(big.Int).SetBytes(sig[:32])
	tx.S = new(big.Int).SetBytes(sig[32:64])
	tx.V = new(big.Int).SetBytes(f.CalculateV(sig[64]))

	return tx, nil
}

// calculateV returns the V value for transactions pre EIP155
func (f *FrontierSigner) CalculateV(parity byte) []byte {
	reference := big.NewInt(int64(parity))
	reference.Add(reference, big27)

	return reference.Bytes()
}

// NewEIP155Signer returns a new EIP155Signer object
func NewEIP155Signer(forks chain.ForksInTime, chainID uint64) *EIP155Signer {
	return &EIP155Signer{chainID: chainID, isHomestead: forks.Homestead}
}

type EIP155Signer struct {
	chainID     uint64
	isHomestead bool
}

// Hash is a wrapper function that calls calcTeleHash with the EIP155Signer's chainID
func (e *EIP155Signer) Hash(tx *types.Telegram) types.Hash {
	return calcTeleHash(tx, e.chainID)
}

// Sender returns the transaction sender
func (e *EIP155Signer) Sender(tele *types.Telegram) (types.Address, error) {
	protected := true

	// Check if v value conforms to an earlier standard (before EIP155)
	bigV := big.NewInt(0)
	if tele.V != nil {
		bigV.SetBytes(tele.V.Bytes())
	}

	if vv := bigV.Uint64(); bits.Len(uint(vv)) <= 8 {
		protected = vv != 27 && vv != 28
	}

	if !protected {
		return (&FrontierSigner{}).Sender(tele)
	}

	// Reverse the V calculation to find the original V in the range [0, 1]
	// v = CHAIN_ID * 2 + 35 + {0, 1}
	mulOperand := big.NewInt(0).Mul(big.NewInt(int64(e.chainID)), big.NewInt(2))
	bigV.Sub(bigV, mulOperand)
	bigV.Sub(bigV, big35)

	sig, err := encodeSignature(tele.R, tele.S, bigV, e.isHomestead)
	if err != nil {
		return types.Address{}, err
	}

	pub, err := Ecrecover(e.Hash(tele).Bytes(), sig)
	if err != nil {
		return types.Address{}, err
	}

	buf := Keccak256(pub[1:])[12:]

	return types.BytesToAddress(buf), nil
}

// Provider returns the telegram provider
func (e *EIP155Signer) Provider(tele *types.Telegram) (types.Address, error) {

	// Check if v value conforms to an earlier standard (before EIP155)
	bigV := big.NewInt(0)
	if tele.RespV != nil {
		bigV.SetBytes(tele.RespV.Bytes())
	}

	// Reverse the V calculation to find the original V in the range [0, 1]
	// v = CHAIN_ID * 2 + 35 + {0, 1}
	mulOperand := big.NewInt(0).Mul(big.NewInt(int64(e.chainID)), big.NewInt(2))
	bigV.Sub(bigV, mulOperand)
	bigV.Sub(bigV, big35)

	sig, err := encodeSignature(tele.RespR, tele.RespS, bigV, e.isHomestead)
	if err != nil {
		return types.Address{}, err
	}

	pub, err := Ecrecover(tele.RespHash.Bytes(), sig)
	if err != nil {
		return types.Address{}, err
	}

	buf := Keccak256(pub[1:])[12:]

	return types.BytesToAddress(buf), nil
}

// SignTele signs the telegram using the passed in private key
func (e *EIP155Signer) SignTele(
	tx *types.Telegram,
	privateKey *ecdsa.PrivateKey,
) (*types.Telegram, error) {
	tx = tx.Copy()

	h := e.Hash(tx)

	sig, err := Sign(privateKey, h[:])
	if err != nil {
		return nil, err
	}

	tx.R = new(big.Int).SetBytes(sig[:32])
	tx.S = new(big.Int).SetBytes(sig[32:64])
	tx.V = new(big.Int).SetBytes(e.CalculateV(sig[64]))

	return tx, nil
}

// calculateV returns the V value for transaction signatures. Based on EIP155
func (e *EIP155Signer) CalculateV(parity byte) []byte {
	reference := big.NewInt(int64(parity))
	reference.Add(reference, big35)

	mulOperand := big.NewInt(0).Mul(big.NewInt(int64(e.chainID)), big.NewInt(2))

	reference.Add(reference, mulOperand)

	return reference.Bytes()
}

// encodeSignature generates a signature value based on the R, S and V value
func encodeSignature(R, S, V *big.Int, isHomestead bool) ([]byte, error) {
	if !ValidateSignatureValues(V, R, S, isHomestead) {
		return nil, fmt.Errorf("invalid txn signature")
	}

	sig := make([]byte, 65)
	copy(sig[32-len(R.Bytes()):32], R.Bytes())
	copy(sig[64-len(S.Bytes()):64], S.Bytes())
	sig[64] = byte(V.Int64()) // here is safe to convert it since ValidateSignatureValues will validate the v value

	return sig, nil
}
