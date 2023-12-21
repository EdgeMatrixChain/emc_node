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

type GetAppNodeResponse struct {
	NodeId string `json:"data"`
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
	response := &GetAppNodeResponse{}
	err = json.Unmarshal(jsonBytes, response)
	if err != nil {
		err = errors.New("GetAppNodeResponse json.Unmarshal error")
		return
	}
	nodeId = response.NodeId
	return
}

func (p *AppAgent) GetAppOrigin() (err error, nodeId string) {
	err = nil
	nodeId = ""
	apiUrl := p.appPath + "/hubapi/v1/getOrigin"
	jsonBytes, err := p.httpClient.SendGetRequest(apiUrl)
	if err != nil {
		err = errors.New("GetAppOrigin error:" + err.Error())
		return
	}
	response := &GetAppNodeResponse{}
	err = json.Unmarshal(jsonBytes, response)
	if err != nil {
		err = errors.New("GetAppOriginResponse json.Unmarshal error")
		return
	}
	nodeId = response.NodeId
	return
}
