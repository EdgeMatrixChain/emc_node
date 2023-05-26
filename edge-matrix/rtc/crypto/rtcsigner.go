package crypto

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/rtc"
	"math/big"
	"math/bits"

	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/emc-protocol/edge-matrix/helper/keccak"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/umbracle/fastrlp"
)

// Signer is a utility interface used to recover data from a rtc
type Signer interface {
	// Hash returns the hash of the transaction
	Hash(rtc *rtc.RtcMsg) types.Hash

	// Sender returns the sender of the transaction
	Sender(rtc *rtc.RtcMsg) (types.Address, error)

	// SignMsg signs a transaction
	SignRtc(rtc *rtc.RtcMsg, priv *ecdsa.PrivateKey) (*rtc.RtcMsg, error)

	// CalculateV calculates the V value based on the type of signer used
	CalculateV(parity byte) []byte
}

// NewSigner creates a new signer object (FrontierSigner)
func NewRtcFrontierSigner(chainID uint64) Signer {
	var signer Signer

	signer = &FrontierSigner{true}

	return signer
}

type FrontierSigner struct {
	isHomestead bool
}

var signerPool fastrlp.ArenaPool

// calcMsgHash calculates the rtc hash (keccak256 hash of the RLP value)
func calcMsgHash(msg *rtc.RtcMsg, chainID uint64) types.Hash {
	a := signerPool.Get()

	v := a.NewArray()

	if len(msg.Subject) < 1 {
		v.Set(a.NewNull())
	} else {
		v.Set(a.NewString((msg.Subject)))
	}
	if len(msg.Application) < 1 {
		v.Set(a.NewNull())
	} else {
		v.Set(a.NewString((msg.Application)))
	}
	if len(msg.Content) < 1 {
		v.Set(a.NewNull())
	} else {
		v.Set(a.NewString((msg.Content)))
	}

	if len(msg.To) < 1 {
		v.Set(a.NewNull())
	} else {
		v.Set(a.NewBytes((msg.To).Bytes()))
	}

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

// Hash is a wrapper function for the calcMsgHash, with chainID 0
func (f *FrontierSigner) Hash(msg *rtc.RtcMsg) types.Hash {
	return calcMsgHash(msg, 0)
}

// Magic numbers from Ethereum, used in v calculation
var (
	big27 = big.NewInt(27)
	big35 = big.NewInt(35)
)

// Sender decodes the signature and returns the sender of the transaction
func (f *FrontierSigner) Sender(msg *rtc.RtcMsg) (types.Address, error) {
	refV := big.NewInt(0)
	if msg.V != nil {
		refV.SetBytes(msg.V.Bytes())
	}

	refV.Sub(refV, big27)

	sig, err := encodeSignature(msg.R, msg.S, refV, f.isHomestead)
	if err != nil {
		return types.Address{}, err
	}

	pub, err := crypto.Ecrecover(f.Hash(msg).Bytes(), sig)
	if err != nil {
		return types.Address{}, err
	}

	buf := crypto.Keccak256(pub[1:])[12:]

	return types.BytesToAddress(buf), nil
}

// SignRtc signs the transaction using the passed in private key
func (f *FrontierSigner) SignRtc(
	tx *rtc.RtcMsg,
	privateKey *ecdsa.PrivateKey,
) (*rtc.RtcMsg, error) {
	tx = tx.Copy()

	h := f.Hash(tx)

	sig, err := crypto.Sign(privateKey, h[:])
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

// Hash is a wrapper function that calls calcMsgHash with the EIP155Signer's chainID
func (e *EIP155Signer) Hash(msg *rtc.RtcMsg) types.Hash {
	return calcMsgHash(msg, e.chainID)
}

// Sender returns the transaction sender
func (e *EIP155Signer) Sender(msg *rtc.RtcMsg) (types.Address, error) {
	protected := true

	// Check if v value conforms to an earlier standard (before EIP155)
	bigV := big.NewInt(0)
	if msg.V != nil {
		bigV.SetBytes(msg.V.Bytes())
	}

	if vv := bigV.Uint64(); bits.Len(uint(vv)) <= 8 {
		protected = vv != 27 && vv != 28
	}

	if !protected {
		return (&FrontierSigner{}).Sender(msg)
	}

	// Reverse the V calculation to find the original V in the range [0, 1]
	// v = CHAIN_ID * 2 + 35 + {0, 1}
	mulOperand := big.NewInt(0).Mul(big.NewInt(int64(e.chainID)), big.NewInt(2))
	bigV.Sub(bigV, mulOperand)
	bigV.Sub(bigV, big35)

	sig, err := encodeSignature(msg.R, msg.S, bigV, e.isHomestead)
	if err != nil {
		return types.Address{}, err
	}

	pub, err := crypto.Ecrecover(e.Hash(msg).Bytes(), sig)
	if err != nil {
		return types.Address{}, err
	}

	buf := crypto.Keccak256(pub[1:])[12:]

	return types.BytesToAddress(buf), nil
}

// SignMsg signs the transaction using the passed in private key
func (e *EIP155Signer) SignRtc(
	msg *rtc.RtcMsg,
	privateKey *ecdsa.PrivateKey,
) (*rtc.RtcMsg, error) {
	msg = msg.Copy()

	h := e.Hash(msg)

	sig, err := crypto.Sign(privateKey, h[:])
	if err != nil {
		return nil, err
	}

	msg.R = new(big.Int).SetBytes(sig[:32])
	msg.S = new(big.Int).SetBytes(sig[32:64])
	msg.V = new(big.Int).SetBytes(e.CalculateV(sig[64]))

	return msg, nil
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
	if !crypto.ValidateSignatureValues(V, R, S, isHomestead) {
		return nil, fmt.Errorf("invalid rtcMsg signature")
	}

	sig := make([]byte, 65)
	copy(sig[32-len(R.Bytes()):32], R.Bytes())
	copy(sig[64-len(S.Bytes()):64], S.Bytes())
	sig[64] = byte(V.Int64()) // here is safe to convert it since ValidateSignatureValues will validate the v value

	return sig, nil
}
