package crypto

import (
	"encoding/json"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/rtc"
	"github.com/emc-protocol/edge-matrix/types"
	"math/big"
	"testing"

	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/stretchr/testify/assert"
)

func TestFrontierKeyGen(t *testing.T) {
	key, err := crypto.GenerateECDSAKey()
	assert.NoError(t, err)

	t.Log(key.PublicKey)
}

func TestFrontierSigner(t *testing.T) {
	signer := &FrontierSigner{}

	//toAddress := types.StringToAddress("1")
	key, err := crypto.GenerateECDSAKey()
	assert.NoError(t, err)

	msg := &rtc.RtcMsg{
		Subject:     "0x1234",
		Application: "edge_rtc",
		Content:     "hello",
	}
	signedMsg, err := signer.SignRtc(msg, key)
	assert.NoError(t, err)

	from, err := signer.Sender(signedMsg)
	assert.NoError(t, err)
	assert.Equal(t, from, crypto.PubKeyToAddress(&key.PublicKey))
}

func TestEIP155Signer_Sender(t *testing.T) {
	t.Parallel()

	//toAddress := types.StringToAddress("1")

	testTable := []struct {
		name    string
		chainID *big.Int
	}{
		{
			"mainnet",
			big.NewInt(1),
		},
		{
			"expanse mainnet",
			big.NewInt(2),
		},
		{
			"ropsten",
			big.NewInt(3),
		},
		{
			"rinkeby",
			big.NewInt(4),
		},
		{
			"goerli",
			big.NewInt(5),
		},
		{
			"kovan",
			big.NewInt(42),
		},
		{
			"geth private",
			big.NewInt(1337),
		},
		{
			"mega large",
			big.NewInt(0).Exp(big.NewInt(2), big.NewInt(20), nil), // 2**20
		},
	}

	for _, testCase := range testTable {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			key, keyGenError := crypto.GenerateECDSAKey()
			if keyGenError != nil {
				t.Fatalf("Unable to generate key")
			}

			txn := &rtc.RtcMsg{
				Subject:     "0x1234",
				Application: "edge_rtc",
				Content:     "hello",
			}

			signer := NewEIP155Signer(chain.AllForksEnabled.At(0), testCase.chainID.Uint64())

			signedTx, signErr := signer.SignRtc(txn, key)
			if signErr != nil {
				t.Fatalf("Unable to sign transaction")
			}

			recoveredSender, recoverErr := signer.Sender(signedTx)
			if recoverErr != nil {
				t.Fatalf("Unable to recover sender")
			}

			assert.Equal(t, recoveredSender.String(), crypto.PubKeyToAddress(&key.PublicKey).String())
		})
	}
}

func TestEIP155Signer_BroadcastSubscribeRtcMsg(t *testing.T) {
	t.Parallel()

	//toAddress := types.BytesToAddress(hex.MustDecodeHex("0x68b95f67a32935e3ed85600F558b74E0d2747120"))

	testTable := []struct {
		name               string
		privateKeyHex      string
		checksummedAddress string
		shouldFail         bool
	}{
		// Generated with Ganache
		{
			"Valid address #1",
			"03b7dfc824b0cbcfe789ec0ce4571f3460befd0490e3d0d2aad8e3c07dbcce14",
			"0x0aF137aa3EcC7d10d926013ee34049AfA77382e6",
			false,
		},
	}

	for _, testCase := range testTable {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//key, keyGenError := GenerateECDSAKey()
			//if keyGenError != nil {
			//	t.Fatalf("Unable to generate key")
			//}

			privateKey, err := crypto.BytesToECDSAPrivateKey([]byte(testCase.privateKeyHex))
			if err != nil && !testCase.shouldFail {
				t.Fatalf("Unable to parse private key, %v", err)
			}
			address, err := crypto.GetAddressFromKey(privateKey)
			if err != nil {
				t.Fatalf("unable to extract key, %v", err)
			}

			assert.Equal(t, testCase.checksummedAddress, address.String())

			txn := &rtc.RtcMsg{
				Subject:     "0x28b9beef497e4fec6e80218f7e756f888fa347ed3eeb27abdec5b3e479f7f5c5",
				Application: "edge_chat",
				Content:     "helloこんにちは你好안녕하세요",
			}
			jsonTxn, err := json.Marshal(txn)
			if err != nil {
				t.Fatalf("Unable to marshal transaction")
			}
			t.Log(string(jsonTxn))
			t.Log("from addr: " + address.String())
			t.Log("publickKey: " + crypto.PubKeyToAddress(&privateKey.PublicKey).String())
			signer := NewEIP155Signer(chain.AllForksEnabled.At(0), big.NewInt(2).Uint64())

			signedTx, signErr := signer.SignRtc(txn, privateKey)
			if signErr != nil {
				t.Fatalf("Unable to sign transaction")
			}

			recoveredSender, recoverErr := signer.Sender(signedTx)
			if recoverErr != nil {
				t.Fatalf("Unable to recover sender")
			}
			t.Log("recoveredSender:" + recoveredSender.String())

			jsonSignedTxn, err := json.Marshal(signedTx)
			if err != nil {
				t.Fatalf("Unable to marshal jsonSignedTxn")
			}
			t.Log(string(jsonSignedTxn))
			bytes := signedTx.MarshalRLP()
			t.Log("MarshalRLP: " + hex.EncodeToHex(bytes))

			assert.Equal(t, recoveredSender.String(), crypto.PubKeyToAddress(&privateKey.PublicKey).String())
		})
	}
}

func TestEIP155Signer_RtcMsg(t *testing.T) {
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
			"03b7dfc824b0cbcfe789ec0ce4571f3460befd0490e3d0d2aad8e3c07dbcce14",
			"0x0aF137aa3EcC7d10d926013ee34049AfA77382e6",
			false,
		},
		//{
		//	"Valid address #2",
		//	"b22be9c19b61adc1d8e89a1dae0346ed274ac9fa239c06286910c29f9fee59d3",
		//	"0x57397Be2eDfc3AF7e3d9a3455aE80A58425Cb767",
		//	false,
		//},
		//{
		//	"Valid address #3",
		//	"c6435f6cb3a8f19111737b72944a0b4a7e52d8a6e95f1ebaa2881679f2087709",
		//	"0x47B7DAc4361062Dfc43d0EA6A2a4C3d27bBcCbdb",
		//	false,
		//},
	}

	for _, testCase := range testTable {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			privateKey, err := crypto.BytesToECDSAPrivateKey([]byte(testCase.privateKeyHex))
			if err != nil && !testCase.shouldFail {
				t.Fatalf("Unable to parse private key, %v", err)
			}
			key := crypto.MarshalPublicKey(&privateKey.PublicKey)
			t.Log("MarshalPublicKey: " + hex.EncodeToString(key))
			address, err := crypto.GetAddressFromKey(privateKey)
			if err != nil {
				t.Fatalf("unable to extract key, %v", err)
			}

			assert.Equal(t, testCase.checksummedAddress, address.String())

			txn := &rtc.RtcMsg{
				Subject:     "0x8eeb338239ada22d81ffb7adc995fe31a4d1dc2d701bc8a58fffe5b53e14281e",
				Application: "edge_chat",
				Content:     "hello",
				To:          types.ZeroAddress,
			}
			jsonTxn, err := json.Marshal(txn)
			if err != nil {
				t.Fatalf("Unable to marshal transaction")
			}
			t.Log(string(jsonTxn))
			t.Log("from addr: " + address.String())
			t.Log("address: " + crypto.PubKeyToAddress(&privateKey.PublicKey).String())
			signer := NewEIP155Signer(chain.AllForksEnabled.At(0), big.NewInt(2).Uint64())

			signedTx, signErr := signer.SignRtc(txn, privateKey)
			if signErr != nil {
				t.Fatalf("Unable to sign transaction")
			}

			recoveredSender, recoverErr := signer.Sender(signedTx)
			if recoverErr != nil {
				t.Fatalf("Unable to recover sender")
			}
			t.Log("recoveredSender:" + recoveredSender.String())

			jsonSignedTxn, err := json.Marshal(signedTx)
			if err != nil {
				t.Fatalf("Unable to marshal jsonSignedTxn")
			}
			t.Log(string(jsonSignedTxn))
			bytes := signedTx.MarshalRLP()
			t.Log("MarshalRLP: " + hex.EncodeToHex(bytes))

			assert.Equal(t, recoveredSender.String(), crypto.PubKeyToAddress(&privateKey.PublicKey).String())
		})
	}
}

func TestEIP155Signer_RtcMsgTo(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name               string
		privateKeyHex      string
		checksummedAddress string
		shouldFail         bool
	}{
		// Generated with Ganache
		{
			"Valid address #2",
			"b22be9c19b61adc1d8e89a1dae0346ed274ac9fa239c06286910c29f9fee59d3",
			"0x57397Be2eDfc3AF7e3d9a3455aE80A58425Cb767",
			false,
		},
	}

	for _, testCase := range testTable {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//key, keyGenError := GenerateECDSAKey()
			//if keyGenError != nil {
			//	t.Fatalf("Unable to generate key")
			//}

			privateKey, err := crypto.BytesToECDSAPrivateKey([]byte(testCase.privateKeyHex))
			if err != nil && !testCase.shouldFail {
				t.Fatalf("Unable to parse private key, %v", err)
			}
			address, err := crypto.GetAddressFromKey(privateKey)
			if err != nil {
				t.Fatalf("unable to extract key, %v", err)
			}

			assert.Equal(t, testCase.checksummedAddress, address.String())

			txn := &rtc.RtcMsg{
				Subject:     "0x8eeb338239ada22d81ffb7adc995fe31a4d1dc2d701bc8a58fffe5b53e14281e",
				Application: "edge_chat",
				Content:     "hello",
				To:          types.StringToAddress("0x0aF137aa3EcC7d10d926013ee34049AfA77382e6"),
			}
			jsonTxn, err := json.Marshal(txn)
			if err != nil {
				t.Fatalf("Unable to marshal transaction")
			}
			t.Log(string(jsonTxn))
			t.Log("from addr: " + address.String())
			t.Log("publickKey: " + crypto.PubKeyToAddress(&privateKey.PublicKey).String())
			signer := NewEIP155Signer(chain.AllForksEnabled.At(0), big.NewInt(2).Uint64())

			signedTx, signErr := signer.SignRtc(txn, privateKey)
			if signErr != nil {
				t.Fatalf("Unable to sign transaction")
			}

			recoveredSender, recoverErr := signer.Sender(signedTx)
			if recoverErr != nil {
				t.Fatalf("Unable to recover sender")
			}
			t.Log("recoveredSender:" + recoveredSender.String())

			jsonSignedTxn, err := json.Marshal(signedTx)
			if err != nil {
				t.Fatalf("Unable to marshal jsonSignedTxn")
			}
			t.Log(string(jsonSignedTxn))
			bytes := signedTx.MarshalRLP()
			t.Log("MarshalRLP: " + hex.EncodeToHex(bytes))

			assert.Equal(t, recoveredSender.String(), crypto.PubKeyToAddress(&privateKey.PublicKey).String())
		})
	}
}

func TestEIP155Signer_ChainIDMismatch(t *testing.T) {
	chainIDS := []uint64{1, 10, 100}
	//toAddress := types.StringToAddress("1")

	for _, chainIDTop := range chainIDS {
		key, keyGenError := crypto.GenerateECDSAKey()
		if keyGenError != nil {
			t.Fatalf("Unable to generate key")
		}

		txn := &rtc.RtcMsg{
			Subject:     "0x1234",
			Application: "edge_rtc",
			Content:     "hello",
		}

		signer := NewEIP155Signer(chain.AllForksEnabled.At(0), chainIDTop)

		signedTx, signErr := signer.SignRtc(txn, key)
		if signErr != nil {
			t.Fatalf("Unable to sign transaction")
		}

		for _, chainIDBottom := range chainIDS {
			signerBottom := NewEIP155Signer(chain.AllForksEnabled.At(0), chainIDBottom)

			recoveredSender, recoverErr := signerBottom.Sender(signedTx)
			if chainIDTop == chainIDBottom {
				// Addresses should match, no error should be present
				assert.NoError(t, recoverErr)

				assert.Equal(t, recoveredSender.String(), crypto.PubKeyToAddress(&key.PublicKey).String())
			} else {
				// There should be an error for mismatched chain IDs
				assert.Error(t, recoverErr)
			}
		}
	}
}
