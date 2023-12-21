package helper

import (
	"crypto/ecdsa"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNodeId(t *testing.T) {

	//privateKeyhex: 	"0c3d062cd3c642735af6a3c1492d761d39a668a67617a457113eaf50860e9e3f"
	//address: 			"0x81e83Dc147B81Db5771D998A2C265cc710BE43a5"
	//ecdsaPrivateKey, err := crypto.BytesToECDSAPrivateKey([]byte("0c3d062cd3c642735af6a3c1492d761d39a668a67617a457113eaf50860e9e3f"))
	//if err != nil {
	//	return
	//}
	//t.Log("ecdsaPrivateKey: ", ecdsaPrivateKey)

	libp2pKey, _, err := network.GenerateAndEncodeLibp2pKey()
	if err != nil {
		return
	}

	/**Public key (address) = 0x5743e866dC9aefd26cCCECc23045bdab6c0e3b16
	  Node ID              = 16Uiu2HAmKtxgX8cBuab9gSju7aZDHtxFsB4NYKXvSrZ5Qd5cb4bU
	*/
	//id, err := peer.Decode("16Uiu2HAmKtxgX8cBuab9gSju7aZDHtxFsB4NYKXvSrZ5Qd5cb4bU")
	//if err != nil {
	//	return
	//}
	//t.Log(pubKey)

	id, err := peer.IDFromPublicKey(libp2pKey.GetPublic())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(id.String())
	publicKey, err := id.ExtractPublicKey()
	if err != nil {
		return
	}
	bytes, err := publicKey.Raw()
	if err != nil {
		return
	}
	t.Log(hex.EncodeToString(bytes), bytes, publicKey.Type().String())
	parsePubKey, err := secp256k1.ParsePubKey(bytes)
	if err != nil {
		return
	}

	address := crypto.PubKeyToAddress(parsePubKey.ToECDSA())
	t.Log(address)

}

func newTestECDSAKey(t *testing.T) (*ecdsa.PrivateKey, []byte) {
	t.Helper()

	testKey, testKeyEncoded, err := crypto.GenerateAndEncodeECDSAPrivateKey()
	assert.NoError(t, err, "failed to initialize ECDSA key")

	return testKey, testKeyEncoded
}

func Test_ecrecover(t *testing.T) {
	t.Parallel()

	testKey, _ := newTestECDSAKey(t)

	buf, err := crypto.MarshalECDSAPrivateKey(testKey)
	assert.NoError(t, err)
	t.Log("MarshalECDSAPrivateKey: ", hex.EncodeToString(buf))

	// public key on the secp256k1 elliptic curve
	// get Ethereum address of a public key
	signerAddress := crypto.PubKeyToAddress(&testKey.PublicKey)
	t.Log("signerAddress: ", signerAddress)

	rawMessage := crypto.Keccak256([]byte{0x1})
	t.Log("rawMessage: ", hex.EncodeToString(rawMessage))

	signature, err := crypto.Sign(
		testKey,
		rawMessage,
	)
	assert.NoError(t, err)
	t.Log("signature: ", hex.EncodeToString(signature))

	recoveredAddress, err := ecrecover(signature, rawMessage)
	assert.NoError(t, err)

	assert.Equal(
		t,
		signerAddress,
		recoveredAddress,
	)
}

func Test_ecrecover2(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name               string
		privateKeyHex      string
		checksummedAddress string
		shouldFail         bool
	}{
		// Generated with Ganache
		{
			"Valid address #1",
			"0c3d062cd3c642735af6a3c1492d761d39a668a67617a457113eaf50860e9e3f",
			"0x81e83Dc147B81Db5771D998A2C265cc710BE43a5",
			false,
		},
		{
			"Valid address #2",
			"71e6439122f6a44884132d54a978318d7218021a5d8f39fd24f440774d564d87",
			"0xCe1f32314aD63F18123b822a23c214DabAA9F7Cf",
			false,
		},
		{
			"Valid address #3",
			"c6435f6cb3a8f19111737b72944a0b4a7e52d8a6e95f1ebaa2881679f2087709",
			"0x47B7DAc4361062Dfc43d0EA6A2a4C3d27bBcCbdb",
			false,
		},
		{
			"Invalid key",
			"c6435f6cb3a8f19111737b72944a0b4a7e52d8a6e95f1ebaa2881679f",
			"",
			true,
		},
	}

	for _, testCase := range testTable {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {

			testKey, err := crypto.BytesToECDSAPrivateKey([]byte(testCase.privateKeyHex))
			if err != nil && !testCase.shouldFail {
				t.Fatalf("Unable to parse private key, %v", err)
			}

			if !testCase.shouldFail {
				address, err := crypto.GetAddressFromKey(testKey)
				if err != nil {
					t.Fatalf("unable to extract key, %v", err)
				}

				assert.Equal(t, testCase.checksummedAddress, address.String())

				buf, err := crypto.MarshalECDSAPrivateKey(testKey)
				assert.NoError(t, err)
				t.Log("MarshalECDSAPrivateKey: ", hex.EncodeToString(buf))

				// public key on the secp256k1 elliptic curve
				// get Ethereum address of a public key
				signerAddress := crypto.PubKeyToAddress(&testKey.PublicKey)
				t.Log("signerAddress: ", signerAddress)

				rawMessage := crypto.Keccak256([]byte{0x1})
				t.Log("rawMessage: ", hex.EncodeToString(rawMessage))

				sigBytes, err := crypto.Sign(
					testKey,
					rawMessage,
				)

				assert.NoError(t, err)
				t.Log("signature: ", hex.EncodeToString(sigBytes))

				R := new(big.Int).SetBytes(sigBytes[:32])
				S := new(big.Int).SetBytes(sigBytes[32:64])
				V := new(big.Int).SetBytes(CalculateV(sigBytes[64], 2))

				t.Log("signature.R: ", R)
				t.Log("signature.S: ", S)
				t.Log("signature.V: ", V)

				recoveredAddress, err := ecrecover(sigBytes, rawMessage)
				assert.NoError(t, err)

				assert.Equal(
					t,
					signerAddress,
					recoveredAddress,
				)
			} else {
				assert.Nil(t, testKey)
			}
		})
	}

}
