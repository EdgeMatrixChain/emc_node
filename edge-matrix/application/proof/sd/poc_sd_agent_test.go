package sd

import (
	"fmt"
	"github.com/brett-lempereur/ish"
	"github.com/emc-protocol/edge-matrix/application/hub"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/helper/ic/agent"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/identity"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/principal"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/xor"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestPrivKey = "7c4b45acadbcd3fa0312ff99ab002ca0cfdab0587a01dbb02f13da282ebfc6a4"
)

func TestSdModels(t *testing.T) {
	pocSD := NewPocSD("http://36.155.7.132:7860")
	response, err := pocSD.SdModels()
	if err != nil {
		t.Error(err)
	}
	for _, model := range response {
		t.Log(model.ModelName, ":", model.Hash)
	}
}

func TestSdLoras(t *testing.T) {
	pocSD := NewPocSD("http://36.155.7.132:7860")
	response, err := pocSD.SdLoras()
	if err != nil {
		t.Error(err)
	}
	for _, loral := range response {
		t.Log(loral.Metadata.Ss_output_name, ":", loral.Metadata.Sshs_model_hash)
	}
}

func TestProofByTxt2imgWithModel(t *testing.T) {
	pocSD := NewPocSD("http://36.155.7.132:7860")
	hashString := "0xc09008b138b5ad15bebbd28539b6f3c62a1bcc75ee6a09c34ab6b27e96d05c19"
	bi, _ := pocSD.MakeSeedByHashString(hashString)
	sdModelHash, imageHash, md5sum, infoString, err := pocSD.ProofByTxt2imgWithModel(hashString, bi, "darkSushi25D25D_v20")
	if err != nil {
		t.Error(err)
	}
	t.Log("sdModelHash:", sdModelHash, "imageHash", imageHash, "md5sum:", md5sum, "info", infoString)
	//assert.Equal(t, "e6415c4892", sdModelHash)
}

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

func TestSDModelTest(t *testing.T) {
	pocSD := NewPocSD("http://36.155.7.132:7860")
	// fetch checkpoint models
	cpModels, err := pocSD.SdModels()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("cpModels:", cpModels)
	needReload := false
	localModels := make(map[string]string, 0)
	checkpointModels := make(map[string]string, 0)
	pocModels := make([]string, 0)

	for _, model := range cpModels {
		if model.Sha256 == "" {
			needReload = true
			// get modelHash
			hashString := "0xc09008b138b5ad15bebbd28539b6f3c62a1bcc75ee6a09c34ab6b27e96d05c19"
			bi, _ := pocSD.MakeSeedByHashString(hashString)
			_, _, _, _, err := pocSD.ProofByTxt2imgWithModel(hashString, bi, model.ModelName)
			if err != nil {
				t.Error(err)
				return
			}
		} else {
			localModels[model.Sha256] = model.ModelName
			checkpointModels[model.Sha256] = model.ModelName
		}
	}
	if needReload {
		response, err := pocSD.SdModels()
		if err != nil {
			t.Error(err)
			return
		}
		for _, model := range response {
			if model.Sha256 != "" {
				localModels[model.Sha256] = model.ModelName
				checkpointModels[model.Sha256] = model.ModelName
			}
		}
	}
	t.Log("checkpointModels:", checkpointModels)

	// fetch loras models
	loraModels, err := pocSD.SdLoras()
	if err != nil {
		t.Error(err)
		return
	}
	for _, loral := range loraModels {
		localModels[loral.Metadata.Sshs_model_hash] = loral.Metadata.Ss_output_name
	}

	// fetch white list from canister
	missedModelHash := ""
	icAgent := agent.NewWithHost("https://ic0.app", false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	hubAgent := hub.NewHubAgent(hclog.NewNullLogger(), icAgent)
	wlModels, err := hubAgent.ListModelsByeType("StableDiffusion")
	t.Log("wlModels:", wlModels)
	for _, modelHash := range wlModels {
		if _, ok := localModels[modelHash]; !ok {
			missedModelHash = missedModelHash + " " + modelHash
		} else {
			if modelName, ok := checkpointModels[modelHash]; ok {
				pocModels = append(pocModels, modelName)
			}
		}

	}
	t.Log("pocModels:", pocModels)
	if len(missedModelHash) > 0 {
		t.Error("missed modelHash:" + missedModelHash)
		return
	}
}
