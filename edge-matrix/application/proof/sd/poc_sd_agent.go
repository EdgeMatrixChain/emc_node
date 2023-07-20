package sd

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brett-lempereur/ish"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/helper/rpc"
	"image"
	"math/big"
)

type PocSD struct {
	sdPath     string
	httpClient *rpc.FastHttpClient
}

func NewPocSD(sdPath string) *PocSD {
	poc := &PocSD{
		httpClient: rpc.NewDefaultHttpClient(),
		sdPath:     sdPath,
	}
	return poc
}

type txt2imgResponse struct {
	Images []string `json:"images"`
	Info   string   `json:"info"`
}
type txt2imgResponseInfo struct {
	Prompt                string      `json:"prompt"`
	AllPrompts            []string    `json:"all_prompts"`
	NegativePrompt        string      `json:"negative_prompt"`
	AllNegativePrompts    []string    `json:"all_negative_prompts"`
	Seed                  int         `json:"seed"`
	AllSeeds              []int       `json:"all_seeds"`
	Subseed               int64       `json:"subseed"`
	AllSubseeds           []int64     `json:"all_subseeds"`
	SubseedStrength       float64     `json:"subseed_strength"`
	Width                 int         `json:"width"`
	Height                int         `json:"height"`
	SamplerName           string      `json:"sampler_name"`
	CfgScale              float64     `json:"cfg_scale"`
	Steps                 int         `json:"steps"`
	BatchSize             int         `json:"batch_size"`
	RestoreFaces          bool        `json:"restore_faces"`
	FaceRestorationModel  interface{} `json:"face_restoration_model"`
	SdModelHash           string      `json:"sd_model_hash"`
	SeedResizeFromW       int         `json:"seed_resize_from_w"`
	SeedResizeFromH       int         `json:"seed_resize_from_h"`
	DenoisingStrength     float64     `json:"denoising_strength"`
	ExtraGenerationParams struct {
	} `json:"extra_generation_params"`
	IndexOfFirstImage             int      `json:"index_of_first_image"`
	Infotexts                     []string `json:"infotexts"`
	Styles                        []string `json:"styles"`
	JobTimestamp                  string   `json:"job_timestamp"`
	ClipSkip                      int      `json:"clip_skip"`
	IsUsingInpaintingConditioning bool     `json:"is_using_inpainting_conditioning"`
}

func (p *PocSD) MakeSeedByHashString(hashString string) (seedNum int64, err error) {
	seedNum = -1
	err = nil
	if len(hashString) < 18 {
		err = errors.New("length too low")
		return
	}
	bi := big.NewInt(0)
	bi.SetString(hashString[2:16], 16)
	seedNum = bi.Int64()
	return
}

func (p *PocSD) ProofByTxt2img(prompt string, seed int64) (sdModelHash string, imageHash string, md5sum string, infoString string, err error) {
	sdModelHash = ""
	md5sum = ""
	imageHash = ""
	err = nil
	apiUrl := p.sdPath + "/sdapi/v1/txt2img"
	txt2imgReq := ` {
      "enable_hr": false,
      "denoising_strength": 0,
      "firstphase_width": 0,
      "firstphase_height": 0,
      "hr_scale": 2,
      "hr_upscaler": "",
      "hr_second_pass_steps": 0,
      "hr_resize_x": 0,
      "hr_resize_y": 0,
      "prompt": "%s",
      "styles": [
        ""
      ],
      "seed": %d,
      "subseed": %d,
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
      "sampler_index": "Euler",
      "send_images": true,
      "save_images": false,
      "alwayson_scripts": {}
    }`
	postJson := fmt.Sprintf(txt2imgReq, prompt, seed, seed)
	jsonBytes, err := p.httpClient.SendPostJsonRequest(apiUrl, []byte(postJson))
	//log.Println(string(jsonBytes))
	response := &txt2imgResponse{}
	err = json.Unmarshal(jsonBytes, response)
	if err != nil {
		err = errors.New("Response json.Unmarshal error")
		return
	}

	if len(response.Info) < 1 {
		err = errors.New("Response Info error")
		return
	}
	info := &txt2imgResponseInfo{}
	infoString = response.Info
	err = json.Unmarshal([]byte(infoString), info)
	if err != nil {
		err = errors.New("Response json.Unmarshal data.Info  error")
		return
	}
	sdModelHash = info.SdModelHash

	if response.Images == nil || len(response.Images) < 1 {
		err = errors.New("Response Images error")
		return
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(response.Images[0])
	if err != nil {
		err = errors.New("Base64 decode error")
		return
	}

	sum := md5.Sum(decodedBytes)
	md5sum = fmt.Sprintf("%x", sum)

	r := bytes.NewReader(decodedBytes)

	img, _, err := image.Decode(r)
	if err != nil {
		return
	}

	// generate Perceptual or Average hash
	hasher := ish.NewAverageHash(8, 8)
	dh, err := hasher.Hash(img)
	if err != nil {
		return
	} else {
		imageHash = hex.EncodeToString(dh)
	}

	return
}
