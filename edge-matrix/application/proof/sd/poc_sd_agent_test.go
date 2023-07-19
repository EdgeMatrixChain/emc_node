package sd

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestProofByTxt2img(t *testing.T) {
	pocSD := NewPocSD("http://36.155.7.141:7860")
	sdModelHash, md5sum, infoString, err := pocSD.ProofByTxt2img("emc-3123123", 123123123)
	if err != nil {
		t.Error(err)
	}
	t.Log("sdModelHash:", sdModelHash, "md5sum:", md5sum, "info", infoString)
	assert.Equal(t, "e6415c4892", sdModelHash)

	pocSD1 := NewPocSD("http://36.155.7.138:7860")
	sdModelHash1, md5sum1, infoString, err := pocSD1.ProofByTxt2img("emc-3123123", 123123123)
	if err != nil {
		t.Error(err)
	}
	t.Log("sdModelHash:", sdModelHash1, "md5sum:", md5sum1, "info", infoString)
	assert.Equal(t, "e6415c4892", sdModelHash1)

	assert.Equal(t, md5sum1, md5sum)
}

func TestProofByTxt2imgBySeedHash(t *testing.T) {
	pocSD := NewPocSD("http://192.168.31.180:7862")
	hashString := "0xa8cc4bcd2513a7adc040c17158af5174abc9fc143d69d7f64407947d1ea9e638"
	bi, _ := pocSD.MakeSeedByHashString(hashString)
	sdModelHash1, md5sum1, infoString, err := pocSD.ProofByTxt2img(hashString, bi)
	if err != nil {
		t.Error(err)
	}
	t.Log("hashString", sdModelHash1, "md5sum1", md5sum1, "infoString", infoString)
}

func TestMakeSeedNumByHash(t *testing.T) {
	pocSD := NewPocSD("http://183.207.184.184:7772")
	hashString := "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c"
	bi, _ := pocSD.MakeSeedByHashString(hashString)
	t.Log(bi)
}
