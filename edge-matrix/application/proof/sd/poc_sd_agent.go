package sd

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emc-protocol/edge-matrix/helper/rpc"
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

func (p *PocSD) ProofByTxt2img(prompt string, seed int64) (sdModelHash string, md5sum string, err error) {
	sdModelHash = ""
	md5sum = ""
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
    }`
	postJson := fmt.Sprintf(txt2imgReq, prompt, seed)
	bytes, err := p.httpClient.SendPostJsonRequest(apiUrl, []byte(postJson))
	response := &txt2imgResponse{}
	err = json.Unmarshal(bytes, response)
	if err != nil {
		return "", "", errors.New("Response json.Unmarshal error")
	}

	if response.Images == nil || len(response.Images) < 1 {
		return "", "", errors.New("Response Images error")
	}
	images := response.Images
	imageBas64 := images[0]
	sum := md5.Sum([]byte(imageBas64))
	md5sum = fmt.Sprintf("%x", sum)

	if len(response.Info) < 1 {
		return "", "", errors.New("Response Info error")
	}
	info := &txt2imgResponseInfo{}
	infoString := response.Info
	err = json.Unmarshal([]byte(infoString), info)
	if err != nil {
		return "", "", errors.New("Response json.Unmarshal data.Info  error")
	}
	sdModelHash = info.SdModelHash

	return
}
