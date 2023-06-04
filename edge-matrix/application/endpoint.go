package application

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/emc-protocol/edge-matrix/application/proof"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/helper/rpc"
	"github.com/emc-protocol/edge-matrix/miner"
	"github.com/emc-protocol/edge-matrix/types"
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
	"sync"
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
	sync.Mutex
	nextNonce        uint64
	nonceCacheEnable bool

	name       string
	appUrl     string
	h          host.Host
	tag        string
	listener   net.Listener
	httpClient *rpc.FastHttpClient
	signer     Signer
	privateKey *ecdsa.PrivateKey
	address    types.Address
	stream     *eventStream // Event subscriptions

	miner         bool
	application   *Application
	minerAgent    *miner.MinerAgent
	jsonRpcClient *rpc.JsonRpcClient

	peersBlockNumMap map[string]uint64
	pocQueue         *proof.PocQueue
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

func (e *Endpoint) GetNextNonce() (uint64, error) {
	e.Lock()
	defer e.Unlock()
	// TODO disable cache in 60s if no tasks in processing
	if !e.nonceCacheEnable {
		nonce, err := e.jsonRpcClient.GetNextNonce(e.address.String())
		if err != nil {
			e.logger.Error("runPoc --> unable to GetNextNonce, %v", err)
			return 0, err
		}
		e.nextNonce = nonce
		e.nonceCacheEnable = true
		return e.nextNonce, nil
	}
	return e.nextNonce, nil
}

func (e *Endpoint) IncreaseNonce() {
	e.Lock()
	defer e.Unlock()

	if e.nonceCacheEnable {
		e.nextNonce += 1
	}
}
func (e *Endpoint) DisableNonceCache() {
	e.Lock()
	defer e.Unlock()

	e.nonceCacheEnable = false
}

func (e *Endpoint) runPoc() {
	for {
		// TODO  get available slot from slotGauge
		tt := e.pocQueue.PopTask()

		if tt == nil {
			// The poc queue is closed,
			// no further poc tasks are incoming
			e.logger.Error("The poc queue is closed")
			return
		}
		go func() {
			pocData := tt.GetPocCpuDataInfo()

			// DO poc_cpu
			start := time.Now()
			//  get data from peer
			nonce, err := e.GetNextNonce()
			if err != nil {
				e.logger.Error("\"runPoc -->GetNextNonce", "err:", err.Error())
				return
			}
			e.logger.Info(fmt.Sprintf("Calling peer [%s] as validator [%s]", pocData.NodeId, e.getID().String()), "queue.len", e.pocQueue.Len(), "nonce", nonce)
			inputString := fmt.Sprintf("{\"peerId\": \"%s\",\"endpoint\": \"/poc_cpu\",\"Input\": {\"seed\": \"%s\"}}", pocData.NodeId, pocData.Seed)
			response, err := e.jsonRpcClient.SendRawTelegram(
				rpc.EdgeCallPrecompile,
				nonce,
				inputString,
				e.privateKey,
			)
			if err != nil {
				e.DisableNonceCache()
				e.logger.Warn("\"runPoc -->SendRawTelegram for poc_cpu", "err:", err.Error())
				return
			}
			e.IncreaseNonce()
			e.logger.Info("SendRawTelegram", "TelegramHash:", response.Result.TelegramHash)
			respBytes, err := base64.StdEncoding.DecodeString(response.Result.Response)
			if err != nil {
				e.logger.Warn("runPoc -->base64 decode err: ", err.Error())
				return
			}
			var dataMapJson []map[string]string
			err = json.Unmarshal(respBytes, &dataMapJson)
			if err != nil {
				e.logger.Warn("runPoc --> json.Unmarshal", "resp", string(respBytes), "err", err.Error())
				return
			}
			dataMap := make(map[string][]byte)
			for _, data := range dataMapJson {
				bytes, err := hex.DecodeString(data["v"])
				if err != nil {
					e.logger.Warn("runPoc --> hex.DecodeString(data[\"v\"]) err: ", err.Error())
					continue
				}
				dataMap[data["k"]] = bytes
			}
			usedTime := time.Since(start).Milliseconds()
			// validate data
			//if s.logger.IsDebug() {
			//	s.logger.Debug("PeerData: {")
			//	for dataKey, bytes := range dataMap {
			//		s.logger.Debug(dataKey, hex.EncodeToString(bytes))
			//	}
			//	s.logger.Debug("}")
			//}
			var hashArray = make([]string, proof.DefaultHashProofCount)
			target := proof.DefaultHashProofTarget
			loops := proof.DefaultHashProofCount
			i := 0
			initSeed := pocData.Seed
			for i < loops {
				seed := fmt.Sprintf("%s,%d", initSeed, i)
				hashArray[i] = seed
				i += 1
			}

			validateSuccess := 0
			validateStart := time.Now()
			for _, hash := range hashArray {
				isValidate := proof.ValidateHash(hash, target, dataMap[hash])
				if isValidate {
					validateSuccess += 1
				}
			}

			validateUsedTime := time.Since(validateStart).Milliseconds()
			rate := float32(validateSuccess) / float32(proof.DefaultHashProofCount)
			e.logger.Debug(fmt.Sprintf("used time for validate\t\t: %dms", validateUsedTime))
			result := fmt.Sprintf("validate success\t\t\t: %d/%d rate:%f nodeID:%s", validateSuccess, loops, rate, pocData.NodeId)
			e.logger.Info(result)
			if rate >= 0.95 {
				// valid proof
				e.logger.Info("\n------------------------------------------\nSubmit proof to IC", "usedTime(ms)", usedTime, "blockNumber", pocData.BlockNum, "NodeID", pocData.NodeId)
				// submit proof result to IC canister
				err := e.minerAgent.SubmitValidation(
					int64(pocData.BlockNum),
					e.minerAgent.GetIdentity(),
					usedTime,
					pocData.NodeId,
				)
				if err != nil {
					e.logger.Warn("\n------------------------------------------\nSubmitValidation:", "err", err)
					return
				}
			}
		}()

	}
}

func NewApplicationEndpoint(logger hclog.Logger,
	privateKey *ecdsa.PrivateKey,
	srvHost host.Host,
	name string,
	appUrl string,
	miner bool,
	blockchainStore blockchainStore,
	minerAgent *miner.MinerAgent,
	jsonRpcClient *rpc.JsonRpcClient) (*Endpoint, error) {
	endpoint := &Endpoint{
		logger:           logger.Named("app_endpoint"),
		name:             name,
		appUrl:           appUrl,
		h:                srvHost,
		tag:              ProtoTagEcApp,
		stream:           &eventStream{},
		miner:            miner,
		minerAgent:       minerAgent,
		peersBlockNumMap: make(map[string]uint64),
		jsonRpcClient:    jsonRpcClient,
		pocQueue:         proof.NewPocQueue(),
		nonceCacheEnable: false,
	}
	endpoint.httpClient = rpc.NewDefaultHttpClient()
	listener, err := gostream.Listen(srvHost, protocol.ID(ProtoTagEcApp))
	if err != nil {
		return nil, err
	}
	endpoint.listener = listener

	address, err := crypto.GetAddressFromKey(privateKey)
	if err != nil {
		endpoint.logger.Error("unable to extract key, error: %v", err.Error())
		return nil, err
	}
	endpoint.address = address
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
				endpoint.logger.Info("Application---->", "push", event.LatestApp())
			}
			ticker.Stop()
		}()
	}

	go func() {
		http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			endpoint.logger.Debug(fmt.Sprintf("/api =>request: %s", string(body)))
			_, err = json.Marshal(body)
			if err != nil {
				http.Error(w, err.Error(), 400)
			}
			var obj struct {
				Method  string          `json:"method"`
				Headers []string        `json:"headers"`
				Path    string          `json:"path"`
				Body    json.RawMessage `json:"body,omitempty"`
			}
			if err := json.Unmarshal(body, &obj); err != nil {
				http.Error(w, err.Error(), 400)
			}
			if obj.Method == "GET" {
				resp, err := endpoint.httpClient.SendGetRequest(endpoint.appUrl + obj.Path)
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
				resp, err := endpoint.httpClient.SendPostJsonRequest(
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
				http.Error(w, err.Error(), 400)
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

			writeResponse(w, info, endpoint)
		})

		http.HandleFunc("/idl", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			idlData, err := os.ReadFile("idl.json")
			if nil != err {
				// TODO Fetch idl json text through GET #{appUrl}/idl
				idlData = []byte("{}")
			}
			writeResponse(w, idlData, endpoint)
		})

		http.HandleFunc("/poc_cpu", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			pocResp := make([]byte, 0)
			body, err := io.ReadAll(r.Body)
			if err != nil {
				pocResp = []byte("endpoint err: " + err.Error())
				writeResponse(w, pocResp, endpoint)
				return
			}
			endpoint.logger.Info(fmt.Sprintf("/poc_cpu =>request: %s", string(body)))
			_, err = json.Marshal(body)
			if err != nil {
				pocResp = []byte("endpoint json.Marshal err: " + err.Error())
				writeResponse(w, pocResp, endpoint)
				return
			}
			var req struct {
				Seed string `json:"seed"`
			}
			if err := json.Unmarshal(body, &req); err != nil {
				pocResp = []byte("endpoint err: " + err.Error())
				writeResponse(w, pocResp, endpoint)
				return
			}
			var data = make(map[string][]byte)
			target := proof.DefaultHashProofTarget
			loops := proof.DefaultHashProofCount
			i := 0
			for i < loops {
				seed := fmt.Sprintf("%s,%d", req.Seed, i)
				_, bytes, err := proof.ProofByCalcHash(seed, target, time.Second*5)
				if err != nil {
					pocResp = []byte("endpoint err: " + err.Error())
					writeResponse(w, pocResp, endpoint)
					return
				}
				data[seed] = bytes
				i += 1
			}
			resp := "["
			dataIdx := 0
			for k, v := range data {
				resp += fmt.Sprintf("{\"k\":\"%s\",\"v\":\"%s\"}", k, hex.EncodeToString(v))
				dataIdx += 1
				if dataIdx < len(data) {
					resp += ","
				}
			}
			resp += "]"
			pocResp = []byte(resp)

			writeResponse(w, pocResp, endpoint)

		})

		http.HandleFunc("/poc_request", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			pocResp := make([]byte, 0)
			body, err := io.ReadAll(r.Body)
			if err != nil {
				pocResp = []byte("endpoint err: " + err.Error())
				writeResponse(w, pocResp, endpoint)
				return
			}
			endpoint.logger.Info(fmt.Sprintf("/poc_request =>request: %s", string(body)))
			_, err = json.Marshal(body)
			if err != nil {
				pocResp = []byte("endpoint json.Marshal body, err: " + err.Error())
				writeResponse(w, pocResp, endpoint)
				return
			}
			var obj struct {
				Node_id string `json:"node_id"`
			}
			if err := json.Unmarshal(body, &obj); err != nil {
				pocResp = []byte("endpoint json.Unmarshal obj, err: " + err.Error())
				writeResponse(w, pocResp, endpoint)
				return
			}

			header := blockchainStore.Header()
			if header != nil {
				blockNumber := header.Number
				// check latest proof number
				latestProofNum, ok := endpoint.peersBlockNumMap[obj.Node_id]
				if !ok {
					latestProofNum = 0
				}
				var blockNumberFixed uint64 = 0
				if (blockNumber - latestProofNum) > proof.DefaultProofBlockMinDuration {
					// send proof task to peer node
					blockNumberFixed = (blockNumber / proof.DefaultProofBlockRange) * proof.DefaultProofBlockRange

					// add poc request to queue
					endpoint.pocQueue.AddTask(
						&proof.PocCpuData{
							NodeId:   obj.Node_id,
							Seed:     header.Hash.String(),
							BlockNum: blockNumberFixed,
						},
						proof.PriorityRequestedPoc,
					)
					endpoint.peersBlockNumMap[obj.Node_id] = blockNumberFixed // commet this line for disable check blocknum
				} else {
					logger.Warn(fmt.Sprintf("\n\n------------------------------------------\ninvalid blockNum, blockNumber: %d, latestProofNum: %d, NodeId:%s\n------------------------------------------\n", blockNumber, latestProofNum, obj.Node_id))
					pocResp = []byte("endpoint err: invalid blockNum")
					writeResponse(w, pocResp, endpoint)
					return
				}
			}

			writeResponse(w, pocResp, endpoint)
		})
		server := &http.Server{}
		server.Serve(listener)
	}()

	go endpoint.runPoc()

	return endpoint, nil
}

func writeResponse(w http.ResponseWriter, info []byte, endpoint *Endpoint) {
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
}
