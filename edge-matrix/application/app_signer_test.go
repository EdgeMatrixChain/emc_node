package application

import (
	"github.com/emc-protocol/edge-matrix/crypto"
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

func TestEIP155Signer_Provider(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name    string
		chainID *big.Int
	}{
		{
			"mainnet",
			big.NewInt(1),
		},
		{
			"testnet",
			big.NewInt(2),
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

			resp := &EdgeResponse{
				RespString: `{"code":0,"text":"abc"}`,
			}

			signer := NewEIP155Signer(chain.AllForksEnabled.At(0), testCase.chainID.Uint64())

			signedResp, signErr := signer.SignEdgeResp(resp, key)
			if signErr != nil {
				t.Fatalf("Unable to sign edge response")
			}

			recoveredSender, recoverErr := signer.Provider(signedResp)
			if recoverErr != nil {
				t.Fatalf("Unable to recover provider")
			}

			assert.Equal(t, recoveredSender.String(), crypto.PubKeyToAddress(&key.PublicKey).String())
		})
	}
}
