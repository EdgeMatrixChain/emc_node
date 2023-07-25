package sd

import (
	"fmt"
	"github.com/brett-lempereur/ish"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/hashicorp/vault/sdk/helper/xor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProofByTxt2img(t *testing.T) {
	pocSD := NewPocSD("http://192.168.31.15:7860")
	hashString := "0xc09008b138b5ad15bebbd28539b6f3c62a1bcc75ee6a09c34ab6b27e96d05c19"
	bi, _ := pocSD.MakeSeedByHashString(hashString)
	sdModelHash, imageHash, md5sum, infoString, err := pocSD.ProofByTxt2img(hashString, bi)
	if err != nil {
		t.Error(err)
	}
	t.Log("sdModelHash:", sdModelHash, "imageHash", imageHash, "md5sum:", md5sum, "info", infoString)
	//assert.Equal(t, "e6415c4892", sdModelHash)

	pocSD1 := NewPocSD("http://192.168.31.14:7860")
	bi1, _ := pocSD1.MakeSeedByHashString(hashString)
	sdModelHash1, imageHash1, md5sum1, infoString, err := pocSD1.ProofByTxt2img(hashString, bi1)
	if err != nil {
		t.Error(err)
	}
	t.Log("sdModelHash:", sdModelHash1, "imageHash", imageHash1, "md5sum:", md5sum1, "info", infoString)
	//assert.Equal(t, "e6415c4892", sdModelHash1)

	assert.Equal(t, imageHash1, imageHash)
}

func TestProofByTxt2imgBySeedHash(t *testing.T) {
	pocSD := NewPocSD("http://36.155.7.141:7860")
	hashString := "0xc09008b138b5ad15bebbd28539b6f3c62a1bcc75ee6a09c34ab6b27e96d05c08"
	bi, _ := pocSD.MakeSeedByHashString(hashString)
	sdModelHash1, imageHash1, md5sum1, infoString, err := pocSD.ProofByTxt2img(hashString, bi)
	if err != nil {
		t.Error(err)
	}
	t.Log("hashString=", sdModelHash1, "imageHash=", imageHash1, "md5sum1=", md5sum1, "infoString=", infoString)
}

func TestMakeSeedNumByHash(t *testing.T) {
	pocSD := NewPocSD("http://127.0.0.1:7860")
	hashString := "0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c"
	bi, _ := pocSD.MakeSeedByHashString(hashString)
	t.Log(bi)
}

func TestDifferenceHash(t *testing.T) {
	hasher := ish.NewDifferenceHash(8, 8)
	test_case := []string{
		"/Users/dev/Downloads/a1.jpg",
		"/Users/dev/Downloads/a2.jpg",
		"/Users/dev/Downloads/a3.jpg",
		"/Users/dev/Downloads/a4.png",
	}
	for _, filename := range test_case {
		img, ft, err := ish.LoadFile(filename)
		if err != nil {
			fmt.Printf("%s: %s\n", filename, err.Error())
			t.Failed()
		}
		dh, err := hasher.Hash(img)
		if err != nil {
			fmt.Printf("%s <%s>: %s", filename, ft, err)
		} else {
			dhs := hex.EncodeToString(dh)
			fmt.Printf("%s <%s>: %s\n", filename, ft, dhs)
		}
	}
}

func TestDifferenceBit(t *testing.T) {
	hash1 := "3e3cbc6a43863c74"
	hash2 := "3e3cbc6a43963d74"
	decodeHex1, err := hex.DecodeHex(hash1)
	if err != nil {
		t.Error(err)
		return
	}
	decodeHex2, err := hex.DecodeHex(hash2)
	if err != nil {
		t.Error(err)
		return
	}
	xorBytes, err := xor.XORBytes(decodeHex1, decodeHex2)
	if err != nil {
		t.Error(err)
		return
	}
	toHex := hex.EncodeToString(xorBytes)
	t.Log(toHex)

	differenceBitCount, err := DifferenceBitCount(hash1, hash2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("differenceBitCount:", differenceBitCount)
}
