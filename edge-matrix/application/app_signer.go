package application

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/helper/keccak"
	"github.com/emc-protocol/edge-matrix/secrets"
	"github.com/emc-protocol/edge-matrix/secrets/helper"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/umbracle/fastrlp"
	"math/big"
)

// Signer is a utility interface used to recover data from a rtc
type Signer interface {
	// Hash returns the hash of the transaction
	Hash(resp *EdgeResponse) types.Hash

	// Provider returns the provider of the edge call request
	Provider(resp *EdgeResponse) (types.Address, error)

	// SignMsg signs a transaction
	SignEdgeResp(resp *EdgeResponse, priv *ecdsa.PrivateKey) (*EdgeResponse, error)

	// CalculateV calculates the V value based on the type of signer used
	CalculateV(parity byte) []byte
}

var signerPool fastrlp.ArenaPool

// Magic numbers from Ethereum, used in v calculation
var (
	big27 = big.NewInt(27)
	big35 = big.NewInt(35)
)

// calcCallHash calculates the rtc hash (keccak256 hash of the RLP value)
func calcResponseHash(resp *EdgeResponse, chainID uint64) types.Hash {
	a := signerPool.Get()

	v := a.NewArray()

	if len(resp.RespString) < 1 {
		v.Set(a.NewNull())
	} else {
		v.Set(a.NewString(resp.RespString))
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

// NewEIP155Signer returns a new EIP155Signer object
func NewEIP155Signer(forks chain.ForksInTime, chainID uint64) *EIP155Signer {
	return &EIP155Signer{chainID: chainID, isHomestead: forks.Homestead}
}

type EIP155Signer struct {
	chainID     uint64
	isHomestead bool
}

// Hash is a wrapper function that calls calcCallHash with the EIP155Signer's chainID
func (e *EIP155Signer) Hash(resp *EdgeResponse) types.Hash {
	return calcResponseHash(resp, e.chainID)
}

// Provider returns the telegram provider
func (e *EIP155Signer) Provider(resp *EdgeResponse) (types.Address, error) {
	// Check if v value conforms to an earlier standard (before EIP155)
	bigV := big.NewInt(0)
	if resp.V != nil {
		bigV.SetBytes(resp.V.Bytes())
	}

	// Reverse the V calculation to find the original V in the range [0, 1]
	// v = CHAIN_ID * 2 + 35 + {0, 1}
	mulOperand := big.NewInt(0).Mul(big.NewInt(int64(e.chainID)), big.NewInt(2))
	bigV.Sub(bigV, mulOperand)
	bigV.Sub(bigV, big35)

	sig, err := encodeSignature(resp.R, resp.S, bigV, e.isHomestead)
	if err != nil {
		return types.Address{}, err
	}

	pub, err := crypto.Ecrecover(e.Hash(resp).Bytes(), sig)
	if err != nil {
		return types.Address{}, err
	}

	buf := crypto.Keccak256(pub[1:])[12:]

	return types.BytesToAddress(buf), nil
}

// SignMsg signs the transaction using the passed in private key
func (e *EIP155Signer) SignEdgeResp(
	resp *EdgeResponse,
	privateKey *ecdsa.PrivateKey,
) (*EdgeResponse, error) {
	resp = resp.Copy()

	h := e.Hash(resp)

	sig, err := crypto.Sign(privateKey, h[:])
	if err != nil {
		return nil, err
	}

	resp.R = new(big.Int).SetBytes(sig[:32])
	resp.S = new(big.Int).SetBytes(sig[32:64])
	resp.V = new(big.Int).SetBytes(e.CalculateV(sig[64]))

	return resp, nil
}

// calculateV returns the V value for app provider signatures. Based on EIP155
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

// getOrCreateECDSAKey loads ECDSA key or creates a new key
func GetOrCreateECDSAKey(manager secrets.SecretsManager) (*ecdsa.PrivateKey, error) {
	if !manager.HasSecret(secrets.ValidatorKey) {
		if _, err := helper.InitECDSAValidatorKey(manager); err != nil {
			return nil, err
		}
	}

	keyBytes, err := manager.GetSecret(secrets.ValidatorKey)
	if err != nil {
		return nil, err
	}

	return crypto.BytesToECDSAPrivateKey(keyBytes)
}
