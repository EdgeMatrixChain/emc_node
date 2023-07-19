package proof

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/types"
	"time"
)
import (
	_ "crypto/sha256"
)

const (
	DefaultProofBlockRange       = 30
	DefaultProofBlockMinDuration = 300
	DefaultProofDuration         = 15 * 60 * time.Second
	DefaultHashProofTarget       = "0000"
	DefaultHashProofCount        = 60
)

func ProofByCalcHash(seed string, target string, timeout time.Duration) (types.Hash, []byte, error) {
	hash, data, err := generateHash(seed, target, timeout)
	if err != nil {
		return types.ZeroHash, nil, err
	}
	return types.StringToHash("0x" + hash), data, nil
}

func generateHash(seed string, target string, timeout time.Duration) (string, []byte, error) {
	hash := make([]byte, 32)
	ctxt, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		//case <-time.After(1 * time.Second):
		//	fmt.Println("overslept")
		case <-ctxt.Done():
			fmt.Println(ctxt.Err()) // prints "context deadline exceeded"
			return "", nil, nil
		default:
			_, err := rand.Read(hash)
			if err != nil {
				return "", nil, err
			}
			dst := append([]byte(seed), hash...)
			keccak256Hash := crypto.Keccak256Hash(dst)
			hashStr := hex.EncodeToString(keccak256Hash.Bytes())
			if hashStr[0:len(target)] == target {
				return hashStr, hash, nil
			}
		}
	}
	return "", nil, nil
}

func ValidateHash(seed string, target string, bytes []byte) bool {
	dst := append([]byte(seed), bytes...)
	keccak256Hash := crypto.Keccak256Hash(dst)
	hashStr := hex.EncodeToString(keccak256Hash.Bytes())
	if hashStr[0:len(target)] == target {
		return true
	}
	return false
}
