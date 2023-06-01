package application

import "testing"

func TestGet(t *testing.T) {
	httpClient := NewFastHttpClient()
	resp, err := httpClient.sendGetRequest("http://google.com")
	if err != nil {
		resp = []byte("endpoint err: " + err.Error())
	}
	t.Log("resp:", string(resp))
}
