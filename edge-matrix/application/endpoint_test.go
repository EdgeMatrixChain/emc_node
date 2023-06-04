package application

import (
	"bytes"
	"fmt"
	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/helper/rpc"
	"github.com/emc-protocol/edge-matrix/secrets"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"

	"github.com/libp2p/go-libp2p"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peerstore"
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

func TestLocal(t *testing.T) {
	m1, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10000")
	m2, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10001")
	srvHost := newHost(t, m1)
	clientHost := newHost(t, m2)
	defer clientHost.Close()

	srvHost.Peerstore().AddAddrs(clientHost.ID(), clientHost.Addrs(), peerstore.PermanentAddrTTL)
	clientHost.Peerstore().AddAddrs(srvHost.ID(), srvHost.Addrs(), peerstore.PermanentAddrTTL)

	key, _, err := crypto.GenerateAndEncodeECDSAPrivateKey()
	assert.NoError(t, err)
	endpoint, err := NewApplicationEndpoint(hclog.NewNullLogger(), key, srvHost, "ec-test", "http://127.0.0.1/", false, nil, nil, rpc.NewDefaultJsonRpcClient())
	if err != nil {
		return
	}
	endpoint.SetSigner(NewEIP155Signer(chain.AllForksEnabled.At(0), uint64(2)))

	defer endpoint.Close()

	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost, endpoint.getProtocolOption()))
	client := &http.Client{Transport: tr}

	buf := bytes.NewBufferString("Hector")
	res, err := client.Post(fmt.Sprintf("libp2p://%s/echo", endpoint.getID().String()), "text/plain", buf)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	respBuf, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp := &EdgeResponse{}
	err = resp.UnmarshalRLP(respBuf)
	if err != nil {
		t.Fatal(err)
	}

	if resp.RespString != "commited Hector!" {
		t.Errorf("expected Hi Hector! but got %s", resp.RespString)
	}

	t.Log(resp.RespString)
}
