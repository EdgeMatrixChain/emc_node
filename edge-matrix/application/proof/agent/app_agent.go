package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emc-protocol/edge-matrix/helper/rpc"
)

type AppAgent struct {
	appPath    string
	httpClient *rpc.FastHttpClient
}

func NewAppAgent(appPath string) *AppAgent {
	poc := &AppAgent{
		httpClient: rpc.NewDefaultHttpClient(),
		appPath:    appPath,
	}
	return poc
}

type GetDataResponse struct {
	Data string `json:"data"`
}

func (p *AppAgent) BindAppNode(nodeId string) (err error) {
	err = nil
	apiUrl := p.appPath + "/hubapi/v1/bindNode"
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

func (p *AppAgent) GetAppNode() (err error, nodeId string) {
	err = nil
	nodeId = ""
	apiUrl := p.appPath + "/hubapi/v1/getNode"
	jsonBytes, err := p.httpClient.SendGetRequest(apiUrl)
	if err != nil {
		err = errors.New("GetAppNode error:" + err.Error())
		return
	}
	response := &GetDataResponse{}
	err = json.Unmarshal(jsonBytes, response)
	if err != nil {
		err = errors.New("GetDataResponse json.Unmarshal error")
		return
	}
	nodeId = response.Data
	return
}

func (p *AppAgent) GetAppOrigin() (err error, appOrigin string) {
	err = nil
	appOrigin = ""
	apiUrl := p.appPath + "/hubapi/v1/getOrigin"
	jsonBytes, err := p.httpClient.SendGetRequest(apiUrl)
	if err != nil {
		err = errors.New("GetAppOrigin error:" + err.Error())
		return
	}
	response := &GetDataResponse{}
	err = json.Unmarshal(jsonBytes, response)
	if err != nil {
		err = errors.New("GetAppOriginResponse json.Unmarshal error")
		return
	}
	appOrigin = response.Data
	return
}

func (p *AppAgent) GetAppIdl() (err error, appOrigin string) {
	err = nil
	appOrigin = ""
	apiUrl := p.appPath + "/hubapi/v1/getIdl"
	jsonBytes, err := p.httpClient.SendGetRequest(apiUrl)
	if err != nil {
		err = errors.New("GetAppIdl error:" + err.Error())
		return
	}
	response := &GetDataResponse{}
	err = json.Unmarshal(jsonBytes, response)
	if err != nil {
		err = errors.New("GetAppOriginResponse json.Unmarshal error")
		return
	}
	appOrigin = response.Data
	return
}
