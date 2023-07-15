package application

import (
	"fmt"
	"github.com/emc-protocol/edge-matrix/secrets"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// newHost illustrates how to build a libp2p host with secio using
// a randomly generated key-pair
func newHost(t *testing.T, listen multiaddr.Multiaddr) host.Host {
	h, err := libp2p.New(
		libp2p.ListenAddrs(listen),
	)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

type mockSecretManager struct {
	secrets.SecretsManager

	HasSecretFunc func(name string) bool
	GetSecretFunc func(name string) ([]byte, error)
}

func TestMultiAddrTest(t *testing.T) {
	ip_addr_string := "/ip4/21.229.33.23/tcp/51004"
	ma, err := multiaddr.NewMultiaddr(ip_addr_string)
	if err != nil {
		t.Error(err.Error())
	}
	addr, err := ma.ValueForProtocol(multiaddr.P_IP4)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(addr)

	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)\.(\d+)`)
	submatches := re.FindStringSubmatch(addr)
	fmt.Printf("%s.%s.%s.%s\n", submatches[1], submatches[2], submatches[3], submatches[4])

}

func TestRunPocSubmitSlice(t *testing.T) {
	batchSize := 500
	sliceSize := 100
	sliceBegin := 0
	sliceEnd := 0
	wg := sync.WaitGroup{}
	for sliceBegin < batchSize {
		sliceEnd += sliceSize
		if sliceEnd > batchSize {
			sliceEnd = batchSize
		}
		t.Log(sliceBegin, ":", sliceEnd)
		wg.Add(1)
		go func(begin, end int) {
			t.Log("do slice [", begin, ":", end, "]")
			time.Sleep(1 * time.Second)
			wg.Done()
		}(sliceBegin, sliceEnd)

		sliceBegin += sliceSize
	}
	wg.Wait()
	t.Log("done")

}

//func TestLocal(t *testing.T) {
//	m1, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10000")
//	m2, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10001")
//	srvHost := newHost(t, m1)
//	clientHost := newHost(t, m2)
//	defer clientHost.Close()
//
//	srvHost.Peerstore().AddAddrs(clientHost.ID(), clientHost.Addrs(), peerstore.PermanentAddrTTL)
//	clientHost.Peerstore().AddAddrs(srvHost.ID(), srvHost.Addrs(), peerstore.PermanentAddrTTL)
//
//	key, _, err := crypto.GenerateAndEncodeECDSAPrivateKey()
//	assert.NoError(t, err)
//	endpoint, err := NewApplicationEndpoint(hclog.NewNullLogger(), key, srvHost, "ec-test", "http://127.0.0.1/", false, nil, , rpc.NewDefaultJsonRpcClient())
//	if err != nil {
//		return
//	}
//	endpoint.SetSigner(NewEIP155Signer(chain.AllForksEnabled.At(0), uint64(2)))
//
//	defer endpoint.Close()
//
//	tr := &http.Transport{}
//	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost, endpoint.getProtocolOption()))
//	client := &http.Client{Transport: tr}
//
//	buf := bytes.NewBufferString("Hector")
//	res, err := client.Post(fmt.Sprintf("libp2p://%s/echo", endpoint.getID().String()), "text/plain", buf)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer res.Body.Close()
//	respBuf, err := io.ReadAll(res.Body)
//	if err != nil {
//		t.Fatal(err)
//	}
//	resp := &EdgeResponse{}
//	err = resp.UnmarshalRLP(respBuf)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if resp.RespString != "commited Hector!" {
//		t.Errorf("expected Hi Hector! but got %s", resp.RespString)
//	}
//
//	t.Log(resp.RespString)
//}
