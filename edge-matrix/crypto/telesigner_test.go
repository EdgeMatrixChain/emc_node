package crypto

import (
	"encoding/json"
	"fmt"
	"github.com/emc-protocol/edge-matrix/contracts"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"math/big"
	"testing"

	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/stretchr/testify/assert"
)

func TestFrontierKeyGen(t *testing.T) {
	key, err := GenerateECDSAKey()
	assert.NoError(t, err)

	t.Log(key.PublicKey)
}

func TestFrontierSigner(t *testing.T) {
	signer := &FrontierSigner{}

	toAddress := types.StringToAddress("1")
	key, err := GenerateECDSAKey()
	assert.NoError(t, err)

	txn := &types.Telegram{
		To:       &toAddress,
		Value:    big.NewInt(10),
		GasPrice: big.NewInt(0),
	}
	signedTx, err := signer.SignTele(txn, key)
	assert.NoError(t, err)

	from, err := signer.Sender(signedTx)
	assert.NoError(t, err)
	assert.Equal(t, from, PubKeyToAddress(&key.PublicKey))
}

func TestEIP155Signer_Sender(t *testing.T) {
	t.Parallel()

	toAddress := types.StringToAddress("1")

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

			key, keyGenError := GenerateECDSAKey()
			if keyGenError != nil {
				t.Fatalf("Unable to generate key")
			}

			txn := &types.Telegram{
				To:       &toAddress,
				Value:    big.NewInt(1),
				GasPrice: big.NewInt(0),
			}

			signer := NewEIP155Signer(chain.AllForksEnabled.At(0), testCase.chainID.Uint64())

			signedTx, signErr := signer.SignTele(txn, key)
			if signErr != nil {
				t.Fatalf("Unable to sign transaction")
			}

			recoveredSender, recoverErr := signer.Sender(signedTx)
			if recoverErr != nil {
				t.Fatalf("Unable to recover sender")
			}

			assert.Equal(t, recoveredSender.String(), PubKeyToAddress(&key.PublicKey).String())
		})
	}
}

func TestEIP155Signer_TeleCreateRtcSubject(t *testing.T) {
	t.Parallel()
	nonce := 0
	chainId := 1
	testTable := []struct {
		name               string
		privateKeyHex      string
		checksummedAddress string
		shouldFail         bool
	}{
		// Generated with Ganache
		//{
		//	"Valid address #1",
		//	"46b73cb531acd7da2225f809be9572f981743ce862d8f2c3e3c8a80f1cd804db",
		//	"0x0DbDdfF0823A173F016Bb8C27c3a700F8E561B58",
		//	false,
		//},
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

			privateKey, err := BytesToECDSAPrivateKey([]byte(testCase.privateKeyHex))
			if err != nil && !testCase.shouldFail {
				t.Fatalf("Unable to parse private key, %v", err)
			}
			address, err := GetAddressFromKey(privateKey)
			if err != nil {
				t.Fatalf("unable to extract key, %v", err)
			}

			assert.Equal(t, testCase.checksummedAddress, address.String())

			tele := &types.Telegram{
				To:    &contracts.EdgeRtcSubjectPrecompile,
				Nonce: uint64(nonce),
				Input: []byte("edge-chat"),
			}
			jsonTxn, err := json.Marshal(tele)
			if err != nil {
				t.Fatalf("Unable to marshal transaction")
			}
			t.Log(string(jsonTxn))
			t.Log("from addr: " + address.String())
			t.Log("publickKey: " + PubKeyToAddress(&privateKey.PublicKey).String())
			signer := NewEIP155Signer(chain.AllForksEnabled.At(0), uint64(chainId))

			signedTx, signErr := signer.SignTele(tele, privateKey)
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

			assert.Equal(t, recoveredSender.String(), PubKeyToAddress(&privateKey.PublicKey).String())
		})
	}
}
func TestEIP155Signer_TeleCallPocRequest(t *testing.T) {
	t.Parallel()

	toAddress := contracts.EdgeCallPrecompile

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

			privateKey, err := BytesToECDSAPrivateKey([]byte(testCase.privateKeyHex))
			if err != nil && !testCase.shouldFail {
				t.Fatalf("Unable to parse private key, %v", err)
			}
			address, err := GetAddressFromKey(privateKey)
			if err != nil {
				t.Fatalf("unable to extract key, %v", err)
			}

			assert.Equal(t, testCase.checksummedAddress, address.String())

			inputFmt := `{"peerId": "%s","endpoint": "/poc_request","input": "{"node_id": "%s"}"}`

			tele := &types.Telegram{
				To:       &toAddress,
				Value:    big.NewInt(0),
				GasPrice: big.NewInt(0),
				Nonce:    uint64(7),
				Input:    []byte(fmt.Sprintf(inputFmt, "16Uiu2HAmKt7agigzA6oGDdMre4eCU7QER91vrW9M3aneiHEvGu1Y", "16Uiu2HAmQkbuGb3K3DmCyEDvKumSVCphVJCGPGHNoc4CobJbxfsC")),
			}
			jsonTxn, err := json.Marshal(tele)
			if err != nil {
				t.Fatalf("Unable to marshal transaction")
			}
			t.Log(string(jsonTxn))
			t.Log("from addr: " + address.String())
			t.Log("publickKey: " + PubKeyToAddress(&privateKey.PublicKey).String())
			signer := NewEIP155Signer(chain.AllForksEnabled.At(0), big.NewInt(1).Uint64())

			signedTx, signErr := signer.SignTele(tele, privateKey)
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

			assert.Equal(t, recoveredSender.String(), PubKeyToAddress(&privateKey.PublicKey).String())
		})
	}
}

func TestEIP155Signer_TeleCallInfo(t *testing.T) {
	t.Parallel()

	toAddress := contracts.EdgeCallPrecompile

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
		nonce := uint64(0)
		chainId := uint64(2)
		input := `{"peerId": "16Uiu2HAm8MbbU7Cge34Y17GXnMULjhyGHtMGUPXaGdepqUxn77M9","endpoint": "/info","input": ""}`

		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			//key, keyGenError := GenerateECDSAKey()
			//if keyGenError != nil {
			//	t.Fatalf("Unable to generate key")
			//}

			privateKey, err := BytesToECDSAPrivateKey([]byte(testCase.privateKeyHex))
			if err != nil && !testCase.shouldFail {
				t.Fatalf("Unable to parse private key, %v", err)
			}
			address, err := GetAddressFromKey(privateKey)
			if err != nil {
				t.Fatalf("unable to extract key, %v", err)
			}

			assert.Equal(t, testCase.checksummedAddress, address.String())

			tele := &types.Telegram{
				To:       &toAddress,
				Value:    big.NewInt(0),
				GasPrice: big.NewInt(0),
				Nonce:    nonce,
				Input:    []byte(input),
			}
			jsonTxn, err := json.Marshal(tele)
			if err != nil {
				t.Fatalf("Unable to marshal transaction")
			}
			t.Log(string(jsonTxn))
			t.Log("from addr: " + address.String())
			t.Log("publickKey: " + PubKeyToAddress(&privateKey.PublicKey).String())
			signer := NewEIP155Signer(chain.AllForksEnabled.At(0), chainId)

			signedTx, signErr := signer.SignTele(tele, privateKey)
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

			assert.Equal(t, recoveredSender.String(), PubKeyToAddress(&privateKey.PublicKey).String())
		})
	}
}

func TestEIP155Signer_TeleCallApi(t *testing.T) {
	t.Parallel()
	nonce := 20
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

			privateKey, err := BytesToECDSAPrivateKey([]byte(testCase.privateKeyHex))
			if err != nil && !testCase.shouldFail {
				t.Fatalf("Unable to parse private key, %v", err)
			}
			address, err := GetAddressFromKey(privateKey)
			if err != nil {
				t.Fatalf("unable to extract key, %v", err)
			}

			assert.Equal(t, testCase.checksummedAddress, address.String())

			input := `{"peerId":"16Uiu2HAm14xAsnJHDqnQNQ2Qqo1SapdRk9j8mBKY6mghVDP9B9u5","endpoint":"/api","Input":{"method": "POST","headers":[],"path":"/sdapi/v1/txt2img","body":{
      "enable_hr": false,
      "denoising_strength": 0,
      "firstphase_width": 0,
      "firstphase_height": 0,
      "hr_scale": 2,
      "hr_upscaler": "",
      "hr_second_pass_steps": 0,
      "hr_resize_x": 0,
      "hr_resize_y": 0,
      "prompt": "white cat and dog",
      "styles": [
        ""
      ],
      "seed": -1,
      "subseed": -1,
      "subseed_strength": 0,
      "seed_resize_from_h": -1,
      "seed_resize_from_w": -1,
      "sampler_name": "",
      "batch_size": 1,
      "n_iter": 1,
      "steps": 50,
      "cfg_scale": 7,
      "width": 512,
      "height": 512,
      "restore_faces": false,
      "tiling": false,
      "do_not_save_samples": false,
      "do_not_save_grid": false,
      "negative_prompt": "",
      "eta": 0,
      "s_churn": 0,
      "s_tmax": 0,
      "s_tmin": 0,
      "s_noise": 1,
      "override_settings": {},
      "override_settings_restore_afterwards": true,
      "script_args": [],
      "sampler_index": "Euler",
      "script_name": "",
      "send_images": true,
      "save_images": false,
      "alwayson_scripts": {}
    }}}`

			//			input := `{
			//  "peerId": "16Uiu2HAm14xAsnJHDqnQNQ2Qqo1SapdRk9j8mBKY6mghVDP9B9u5",
			//  "endpoint": "/api",
			//  "Input": {
			//    "method": "POST",
			//    "headers": [],
			//    "path": "/sdapi/v1/txt2img",
			//    "body": {
			//      "enable_hr": false,
			//      "denoising_strength": 0,
			//      "firstphase_width": 0,
			//      "firstphase_height": 0,
			//      "hr_scale": 2,
			//      "hr_upscaler": "",
			//      "hr_second_pass_steps": 0,
			//      "hr_resize_x": 0,
			//      "hr_resize_y": 0,
			//      "prompt": "white cat and dog",
			//      "styles": [
			//        ""
			//      ],
			//      "seed": -1,
			//      "subseed": -1,
			//      "subseed_strength": 0,
			//      "seed_resize_from_h": -1,
			//      "seed_resize_from_w": -1,
			//      "sampler_name": "",
			//      "batch_size": 1,
			//      "n_iter": 1,
			//      "steps": 50,
			//      "cfg_scale": 7,
			//      "width": 512,
			//      "height": 512,
			//      "restore_faces": false,
			//      "tiling": false,
			//      "do_not_save_samples": false,
			//      "do_not_save_grid": false,
			//      "negative_prompt": "",
			//      "eta": 0,
			//      "s_churn": 0,
			//      "s_tmax": 0,
			//      "s_tmin": 0,
			//      "s_noise": 1,
			//      "override_settings": {},
			//      "override_settings_restore_afterwards": true,
			//      "script_args": [],
			//      "sampler_index": "Euler",
			//      "script_name": "",
			//      "send_images": true,
			//      "save_images": false,
			//      "alwayson_scripts": {}
			//    }
			//  }
			//}`
			if err != nil {
				return
			}
			tele := &types.Telegram{
				To:    &contracts.EdgeCallPrecompile,
				Nonce: uint64(nonce),
				Input: []byte(input),
			}
			jsonTxn, err := json.Marshal(tele)
			if err != nil {
				t.Fatalf("Unable to marshal transaction")
			}
			t.Log(string(jsonTxn))
			t.Log("from addr: " + address.String())
			t.Log("publickKey: " + PubKeyToAddress(&privateKey.PublicKey).String())
			signer := NewEIP155Signer(chain.AllForksEnabled.At(0), big.NewInt(2).Uint64())

			signedTx, signErr := signer.SignTele(tele, privateKey)
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

			assert.Equal(t, recoveredSender.String(), PubKeyToAddress(&privateKey.PublicKey).String())

			ummarshalTx := types.Telegram{}
			err = ummarshalTx.UnmarshalRLP(bytes)
			if err != nil {
				t.Fatalf(err.Error())
				return
			}
			t.Log("ummarshalTx.Nonce: ", ummarshalTx.Nonce)
			t.Log("ummarshalTx.To: ", ummarshalTx.To)
		})
	}
}

func TestEIP155Signer_ChainIDMismatch(t *testing.T) {
	chainIDS := []uint64{1, 10, 100}
	toAddress := types.StringToAddress("1")

	for _, chainIDTop := range chainIDS {
		key, keyGenError := GenerateECDSAKey()
		if keyGenError != nil {
			t.Fatalf("Unable to generate key")
		}

		txn := &types.Telegram{
			To:       &toAddress,
			Value:    big.NewInt(1),
			GasPrice: big.NewInt(0),
		}

		signer := NewEIP155Signer(chain.AllForksEnabled.At(0), chainIDTop)

		signedTx, signErr := signer.SignTele(txn, key)
		if signErr != nil {
			t.Fatalf("Unable to sign transaction")
		}

		for _, chainIDBottom := range chainIDS {
			signerBottom := NewEIP155Signer(chain.AllForksEnabled.At(0), chainIDBottom)

			recoveredSender, recoverErr := signerBottom.Sender(signedTx)
			if chainIDTop == chainIDBottom {
				// Addresses should match, no error should be present
				assert.NoError(t, recoverErr)

				assert.Equal(t, recoveredSender.String(), PubKeyToAddress(&key.PublicKey).String())
			} else {
				// There should be an error for mismatched chain IDs
				assert.Error(t, recoverErr)
			}
		}
	}
}
