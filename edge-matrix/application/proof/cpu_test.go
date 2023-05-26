package proof

import (
	"crypto/rand"
	"fmt"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"testing"
	"time"
)

func TestProofBlockNumber(t *testing.T) {
	nums := make([]int64, 100)
	i := int64(0)
	for i < 100 {
		nums[i] = 1892345 + i
		i += 1
	}
	for _, blockNumber := range nums {
		t.Log("blockNumber:", blockNumber)
		blockNumberFixed := (blockNumber / 30) * 30
		t.Log("blockNumberFixed:", blockNumberFixed)
	}

}

func TestProofByHash(t *testing.T) {
	var data = make(map[string]*[]byte)
	target := "0000"
	loops := 60
	i := 0
	start := time.Now()
	for i < loops {
		randBytes := make([]byte, 32)
		_, err := rand.Read(randBytes)
		if err != nil {
			return
		}
		seed := hex.EncodeToHex(randBytes)
		_, bytes, err := ProofByCalcHash(seed, target, time.Second*3)
		if err != nil {
			t.Log(fmt.Sprintf("err: %s", err.Error()))
			return
		}
		data[seed] = &bytes
		i += 1
	}
	t.Log(fmt.Sprintf("calc time			: %fs", time.Since(start).Seconds()))

	validateSuccess := 0
	validateStart := time.Now()
	for seed, bytes := range data {
		validateHash := ValidateHash(seed, target, *bytes)
		if validateHash {
			validateSuccess += 1
		}
	}
	t.Log(fmt.Sprintf("validate time		: %dms", time.Since(validateStart).Microseconds()))
	t.Log(fmt.Sprintf("validate success	: %d/%d", validateSuccess, loops))

}
