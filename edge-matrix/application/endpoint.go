package application

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/emc-protocol/edge-matrix/application/proof"
	"github.com/hashicorp/go-hclog"
	gostream "github.com/libp2p/go-libp2p-gostream"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	// proto tag for p2phttp
	ProtoTagEcApp = "/em-app"
)

const (
	txSlotSize  = 32 * 1024  // 32kB
	txMaxSize   = 128 * 1024 // 128Kb
	topicNameV1 = "em_app/0.1"

	// maximum allowed number of times an account
	// was excluded from request handler
	maxAccountDemotions uint64 = 100

	// maximum allowed number of consecutive blocks that don't have the account's call
	pruningCooldown = 5000 * time.Millisecond

	// callPoolMetrics is a prefix used for callpool-related metrics
	callPoolMetrics = "callpool"
)

type Application struct {
	Name    string
	Tag     string
	Version string
	PeerID  peer.ID

	// app startup time
	StartupTime uint64
	// app uptime
	Uptime uint64
	// amount of slots currently occupying the app
	GuageHeight uint64
	// max limit
	GuageMax uint64
}

func (a *Application) Copy() *Application {
	newApp := &Application{
		Name:        a.Name,
		Tag:         a.Tag,
		Version:     a.Version,
		PeerID:      a.PeerID,
		StartupTime: a.StartupTime,
		Uptime:      a.Uptime,
		GuageHeight: a.GuageHeight,
		GuageMax:    a.GuageMax,
	}

	return newApp
}

type Endpoint struct {
	logger hclog.Logger

	// gauge for measuring app capacity
	gauge slotGauge

	name       string
	appUrl     string
	h          host.Host
	tag        string
	listener   net.Listener
	httpClient *FastHttpClient
	signer     Signer
	privateKey *ecdsa.PrivateKey

	stream *eventStream // Event subscriptions

	miner       bool
	application *Application
}

func (e *Endpoint) getID() peer.ID {
	return e.h.ID()
}

func (e *Endpoint) getProtocolOption() p2phttp.Option {
	return p2phttp.ProtocolOption(protocol.ID(e.tag))
}

func (e *Endpoint) Close() {
	e.listener.Close()
	e.h.Close()
}

// SetSigner sets the signer the endpint will use
// to validate a edge call response's signature.
func (e *Endpoint) SetSigner(s Signer) {
	e.signer = s
}

func NewApplicationEndpoint(logger hclog.Logger,
	privateKey *ecdsa.PrivateKey, srvHost host.Host, name string, appUrl string, miner bool) (*Endpoint, error) {
	endpoint := &Endpoint{
		logger: logger.Named("app_endpoint"),
		name:   name,
		appUrl: appUrl,
		h:      srvHost,
		tag:    ProtoTagEcApp,
		stream: &eventStream{},
		miner:  miner,
	}
	endpoint.httpClient = NewFastHttpClient()
	listener, err := gostream.Listen(srvHost, protocol.ID(ProtoTagEcApp))
	if err != nil {
		return nil, err
	}
	endpoint.listener = listener

	endpoint.privateKey = privateKey
	// Push the initial event to the stream
	endpoint.stream.push(&Event{})

	// Create an event and send it to the stream
	event := &Event{}
	endpoint.stream.push(event)

	// init application metric
	endpoint.application = &Application{
		Name:        name,
		PeerID:      srvHost.ID(),
		StartupTime: uint64(time.Now().UnixMilli()),
		Uptime:      0,
		GuageHeight: 0,
		GuageMax:    200,
	}

	// TODO check miner status
	if endpoint.miner {
		go func() {
			event := &Event{}
			event.AddNewApp(&Application{
				Name:        endpoint.application.Name,
				PeerID:      endpoint.application.PeerID,
				StartupTime: endpoint.application.StartupTime,
				Uptime:      uint64(time.Now().UnixMilli()) - endpoint.application.StartupTime,
				GuageHeight: endpoint.application.GuageHeight,
				GuageMax:    endpoint.application.GuageMax,
			})
			endpoint.stream.push(event)
			endpoint.logger.Info("endpoint.miner---->", "start push", event)

			ticker := time.NewTicker(proof.DefaultProofDuration)
			for {
				<-ticker.C
				event := &Event{}
				event.AddNewApp(&Application{
					Name:        endpoint.application.Name,
					PeerID:      endpoint.application.PeerID,
					StartupTime: endpoint.application.StartupTime,
					Uptime:      uint64(time.Now().UnixMilli()) - endpoint.application.StartupTime,
					GuageHeight: endpoint.application.GuageHeight,
					GuageMax:    endpoint.application.GuageMax,
				})
				endpoint.stream.push(event)
				endpoint.logger.Info("endpoint.miner---->", "push", event)
			}
			ticker.Stop()
		}()
	}

	go func() {
		http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			endpoint.logger.Debug(fmt.Sprintf("/api =>request: %s", string(body)))
			_, err = json.Marshal(body)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			var obj struct {
				Method  string          `json:"method"`
				Headers []string        `json:"headers"`
				Path    string          `json:"path"`
				Body    json.RawMessage `json:"body,omitempty"`
			}
			if err := json.Unmarshal(body, &obj); err != nil {
				http.Error(w, err.Error(), 500)
			}
			if obj.Method == "GET" {
				resp, err := endpoint.httpClient.sendGetRequest(endpoint.appUrl)
				if err != nil {
					resp = []byte("endpoint err: " + err.Error())
				}
				encodeString := base64.StdEncoding.EncodeToString(resp)
				edgeResp := &EdgeResponse{
					RespString: encodeString,
				}
				endpoint.logger.Debug(fmt.Sprintf("/api =>resp size: %d", len(edgeResp.RespString)))

				signedResp, err := endpoint.signer.SignEdgeResp(edgeResp, endpoint.privateKey)
				if err != nil {
					w.Write([]byte(err.Error()))
				}
				provider, err := endpoint.signer.Provider(signedResp)
				if err != nil {
					return
				}
				signedResp.From = provider
				signedResp.Hash = endpoint.signer.Hash(edgeResp)

				w.Write(signedResp.MarshalRLP())
			} else if obj.Method == "POST" {
				resp, err := endpoint.httpClient.sendPostJsonRequest(
					endpoint.appUrl+obj.Path, obj.Body)
				if err != nil {
					resp = []byte("endpoint err: " + err.Error())
				}
				encodeString := base64.StdEncoding.EncodeToString(resp)
				edgeResp := &EdgeResponse{
					RespString: encodeString,
				}
				endpoint.logger.Debug(fmt.Sprintf("/api =>resp size: %d", len(edgeResp.RespString)))
				endpoint.logger.Debug(fmt.Sprintf("/api =>SignEdgeResp with hash: %s", endpoint.signer.Hash(edgeResp).String()))
				signedResp, err := endpoint.signer.SignEdgeResp(edgeResp, endpoint.privateKey)
				if err != nil {
					w.Write([]byte(err.Error()))
				}
				provider, err := endpoint.signer.Provider(signedResp)
				if err != nil {
					return
				}
				signedResp.From = provider
				signedResp.Hash = endpoint.signer.Hash(edgeResp)

				w.Write(signedResp.MarshalRLP())
			}
		})

		http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			resp := fmt.Sprintf("recieved data: %s", body)
			encodeString := base64.StdEncoding.EncodeToString([]byte(resp))
			edgeResp := &EdgeResponse{
				RespString: encodeString,
			}
			endpoint.logger.Debug(fmt.Sprintf("/api =>resp size: %d", len(edgeResp.RespString)))

			signedResp, err := endpoint.signer.SignEdgeResp(edgeResp, endpoint.privateKey)
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			provider, err := endpoint.signer.Provider(signedResp)
			if err != nil {
				return
			}
			signedResp.From = provider
			signedResp.Hash = endpoint.signer.Hash(edgeResp)

			w.Write(signedResp.MarshalRLP())
		})

		http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			var infoObj struct {
				Name        string `json:"name"`
				PeerID      string `json:"peerId"`
				Uptime      uint64 `json:"uptime"`
				StartupTime uint64 `json:"startupTime"`
				Version     string `json:"version"`
				Tag         string `json:"tag"`
			}
			infoObj.PeerID = endpoint.application.PeerID.String()
			infoObj.Version = endpoint.application.Version
			infoObj.Tag = endpoint.application.Tag
			infoObj.Uptime = uint64(time.Now().UnixMilli()) - endpoint.application.StartupTime
			infoObj.StartupTime = endpoint.application.StartupTime
			infoObj.Name = endpoint.application.Name

			info := make([]byte, 0)
			info, err := json.Marshal(infoObj)
			if err != nil {
				info = []byte("endpoint err: " + err.Error())
			}
			resp := base64.StdEncoding.EncodeToString(info)
			edgeResp := &EdgeResponse{
				RespString: resp,
			}
			endpoint.logger.Debug(fmt.Sprintf("/api =>resp size: %d", len(edgeResp.RespString)))

			signedResp, err := endpoint.signer.SignEdgeResp(edgeResp, endpoint.privateKey)
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			provider, err := endpoint.signer.Provider(signedResp)
			if err != nil {
				return
			}
			signedResp.From = provider
			signedResp.Hash = endpoint.signer.Hash(edgeResp)

			w.Write(signedResp.MarshalRLP())
		})

		http.HandleFunc("/idl", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			idlData, err := os.ReadFile("idl.json")
			if nil != err {
				// TODO Fetch idl json text through GET #{appUrl}/idl
				idlData = []byte("{}")
			}

			resp := base64.StdEncoding.EncodeToString(idlData)
			edgeResp := &EdgeResponse{
				RespString: resp,
			}
			endpoint.logger.Debug(fmt.Sprintf("/api =>resp size: %d", len(edgeResp.RespString)))

			signedResp, err := endpoint.signer.SignEdgeResp(edgeResp, endpoint.privateKey)
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			provider, err := endpoint.signer.Provider(signedResp)
			if err != nil {
				return
			}
			signedResp.From = provider
			signedResp.Hash = endpoint.signer.Hash(edgeResp)

			w.Write(signedResp.MarshalRLP())
		})

		server := &http.Server{}
		server.Serve(listener)
	}()

	return endpoint, nil
}
