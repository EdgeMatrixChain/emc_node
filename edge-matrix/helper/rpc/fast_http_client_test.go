package rpc

import "testing"

func TestGet(t *testing.T) {
	httpClient := NewDefaultHttpClient()
	resp, err := httpClient.SendGetRequest("http://google.com")
	if err != nil {
		resp = []byte("endpoint err: " + err.Error())
	}
	t.Log("resp:", string(resp))
}
