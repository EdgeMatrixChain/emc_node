package sd

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brett-lempereur/ish"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/helper/rpc"
	"github.com/hashicorp/vault/sdk/helper/xor"
	"image"
	"math/big"
	"math/bits"
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

type sdModel struct {
	Title     string `json:"title"`
	ModelName string `json:"model_name"`
	Hash      string `json:"hash"`
	Sha256    string `json:"sha256"`
	Filename  string `json:"filename"`
	Config    string `json:"config"`
}

type sdLoral struct {
	Name     string          `json:"name"`
	Alias    string          `json:"alias"`
	Path     string          `json:"path"`
	Metadata sdLoralMetadata `json:"metadata"`
}

type sdLoralMetadata struct {
	Ss_output_name  string `json:"ss_output_name"`
	Sshs_model_hash string `json:"sshs_model_hash"`
}

type txt2imgResponse struct {
	Images []string `json:"images"`
	Info   string   `json:"info"`
}

type GetAppNodeResponse struct {
	NodeId string `json:"data"`
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

func DifferenceBitCount(hash1, hash2 string) (int, error) {
	decodeHex1, err := hex.DecodeHex(hash1)
	if err != nil {
		return 0, err
	}
	decodeHex2, err := hex.DecodeHex(hash2)
	if err != nil {
		return 0, err
	}
	xorBytes, err := xor.XORBytes(decodeHex1, decodeHex2)
	if err != nil {
		return 0, err
	}

	toInt64, err := BytesToInt64(xorBytes)
	if err != nil {
		return 0, err
	}

	return bits.OnesCount64(toInt64), nil
}

func BytesToInt64(buf []byte) (uint64, error) {
	if buf == nil || len(buf) != 8 {
		return 0, errors.New("length for bytes must be 8")
	}
	return binary.BigEndian.Uint64(buf), nil
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

func (p *PocSD) ProofByTxt2imgWithModel(prompt string, seed int64, modelName string) (sdModelHash string, imageHash string, md5sum string, infoString string, err error) {
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
      "override_settings": {"sd_model_checkpoint":"%s"},
      "override_settings_restore_afterwards": true,
      "sampler_index": "Euler",
      "send_images": true,
      "save_images": false,
      "alwayson_scripts": {}
    }`
	postJson := fmt.Sprintf(txt2imgReq, prompt, seed, seed, modelName)
	jsonBytes, err := p.httpClient.SendPostJsonRequest(apiUrl, []byte(postJson))
	if err != nil {
		err = errors.New("ProofByTxt2imgWithModel error:" + err.Error())
		return
	}
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
	if err != nil {
		err = errors.New("ProofByTxt2img error:" + err.Error())
		return
	}
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

func (p *PocSD) BindAppNode(nodeId string) (err error) {
	err = nil
	apiUrl := p.sdPath + "/hubapi/v1/bindNode"
	bindReq := `{
      "nodeId":"%s"
    }`
	postJson := fmt.Sprintf(bindReq, nodeId)
	_, err = p.httpClient.SendPostJsonRequest(apiUrl, []byte(postJson))
	if err != nil {
		err = errors.New("BindNode error:" + err.Error())
		return
	}
	return
}

func (p *PocSD) GetAppNode() (err error, nodeId string) {
	err = nil
	nodeId = ""
	apiUrl := p.sdPath + "/hubapi/v1/getNode"
	jsonBytes, err := p.httpClient.SendGetRequest(apiUrl)
	if err != nil {
		err = errors.New("BindNode error:" + err.Error())
		return
	}
	response := &GetAppNodeResponse{}
	err = json.Unmarshal(jsonBytes, response)
	if err != nil {
		err = errors.New("GetAppNodeResponse json.Unmarshal error")
		return
	}
	nodeId = response.NodeId
	return
}

func (p *PocSD) SdLoras() (models []sdLoral, err error) {
	err = nil
	apiUrl := p.sdPath + "/sdapi/v1/loras"
	jsonBytes, err := p.httpClient.SendGetRequest(apiUrl)
	//log.Println(string(jsonBytes))
	err = json.Unmarshal(jsonBytes, &models)
	if err != nil {
		err = errors.New("Response json.Unmarshal error")
		return
	}

	return
}

func (p *PocSD) SdModels() (models []sdModel, err error) {
	err = nil
	apiUrl := p.sdPath + "/sdapi/v1/sd-models"
	jsonBytes, err := p.httpClient.SendGetRequest(apiUrl)
	//log.Println(string(jsonBytes))
	err = json.Unmarshal(jsonBytes, &models)
	if err != nil {
		err = errors.New("Response json.Unmarshal error")
		return
	}

	return
}
