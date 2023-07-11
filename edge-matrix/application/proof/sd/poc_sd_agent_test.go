package sd

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestProofByTxt2img(t *testing.T) {
	pocSD := NewPocSD("http://183.207.184.184:7772")
	sdModelHash, md5sum, err := pocSD.ProofByTxt2img("emc-3123123", 123123123)
	if err != nil {
		t.Error(err)
	}
	t.Log("sdModelHash:", sdModelHash, "md5sum:", md5sum)
	assert.Equal(t, "e6415c4892", sdModelHash)
	assert.Equal(t, "2ed218ca9d99263dae3ee9e7ab5a8d1c", md5sum)
}

func TestMakeSeedNumByHash(t *testing.T) {
	pocSD := NewPocSD("http://183.207.184.184:7772")
	hashString := "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c"
	bi, _ := pocSD.MakeSeedByHashString(hashString)
	t.Log(bi)
}
