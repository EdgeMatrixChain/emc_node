package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/emc-protocol/edge-matrix/types"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
	"io"
	"net/http"
	"net/url"
)

type EdgeCall struct {
	PeerId   string          `json:"peerId"`
	Endpoint string          `json:"endpoint"`
	Input    json.RawMessage `json:"input"`
}

func (e *EdgeCall) Copy() *EdgeCall {
	tt := &EdgeCall{
		PeerId:   e.PeerId,
		Endpoint: e.Endpoint,
	}

	if len(e.Input) > 0 {
		tt.Input = make([]byte, len(e.Input))
		copy(tt.Input[:], e.Input)
	}

	return tt
}

func DecodeEdgeCallFromInterface(i interface{}) (*EdgeCall, error) {
	// once the rtc filter is decoded as map[string]interface we cannot use unmarshal json
	raw, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	call := &EdgeCall{}
	if err := json.Unmarshal(raw, &call); err != nil {
		return nil, err
	}

	return call, nil
}

func Call(clientHost host.Host, protoTag string, call *EdgeCall) ([]byte, error) {
	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost, p2phttp.ProtocolOption(protocol.ID(protoTag))))
	client := &http.Client{Transport: tr}

	if call.Input == nil {
		return nil, nil
	}
	raw, err := json.Marshal(call.Input)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(raw)
	res, err := client.Post(fmt.Sprintf("libp2p://%s%s", call.PeerId, call.Endpoint), "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	all, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return all, nil

}

func CallWithFrom(clientHost host.Host, protoTag string, call *EdgeCall, from types.Address) ([]byte, error) {
	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost, p2phttp.ProtocolOption(protocol.ID(protoTag))))
	client := &http.Client{Transport: tr}

	if call.Input == nil {
		return nil, nil
	}
	raw, err := json.Marshal(call.Input)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(raw)

	URL, err := url.Parse(fmt.Sprintf("libp2p://%s%s", call.PeerId, call.Endpoint))
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		URL:    URL,
		Method: "POST",
		Header: http.Header{"Content-Type": {"application/json"}, "Emc-From": {from.String()}, "Emc-Router": {clientHost.ID().String()}},
		Body:   io.NopCloser(buf),
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	all, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return all, nil

}
