package application

import (
	"reflect"
	"time"

	"github.com/valyala/fasthttp"
)

var headerContentTypeJson = []byte("application/json")

type FastHttpClient struct {
	client *fasthttp.Client
}

func NewFastHttpClient() *FastHttpClient {
	// You may read the timeouts from some config
	readTimeout := 60 * time.Second
	writeTimeout := 60 * time.Second
	tcpDialer := &fasthttp.TCPDialer{
		Concurrency:      1000,
		DNSCacheDuration: time.Hour,
	}
	hc := &FastHttpClient{
		client: &fasthttp.Client{
			ReadTimeout:                   readTimeout,
			WriteTimeout:                  writeTimeout,
			NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
			DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
			DisablePathNormalizing:        true,
			Dial:                          tcpDialer.Dial,
		},
	}
	return hc
}

func (f *FastHttpClient) sendGetRequest(url string) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := f.client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err == nil {
		return resp.Body(), nil
	} else {
		return nil, err
	}
}

func (f *FastHttpClient) sendPostJsonRequest(url string, reqEntityBytes []byte) ([]byte, error) {
	// per-request timeout
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw(reqEntityBytes)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := f.client.Do(req, resp)
	if err == nil {
		respBody := resp.Body()
		return respBody, nil
	} else {
		return nil, err
	}
}

func HttpConnError(err error) (string, bool) {
	errName := ""
	known := false
	if err == fasthttp.ErrTimeout {
		errName = "timeout"
		known = true
	} else if err == fasthttp.ErrNoFreeConns {
		errName = "conn_limit"
		known = true
	} else if err == fasthttp.ErrConnectionClosed {
		errName = "conn_close"
		known = true
	} else {
		errName = reflect.TypeOf(err).String()
		if errName == "*net.OpError" {
			// Write and Read errors are not so often and in fact they just mean timeout problems
			errName = "timeout"
			known = true
		}
	}
	return errName, known
}
