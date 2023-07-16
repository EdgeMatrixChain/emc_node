package application

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/emc-protocol/edge-matrix/application/proof"
	"github.com/emc-protocol/edge-matrix/application/proof/helper"
	"github.com/emc-protocol/edge-matrix/application/proof/sd"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/principal"
	"github.com/emc-protocol/edge-matrix/helper/rpc"
	"github.com/emc-protocol/edge-matrix/miner"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/emc-protocol/edge-matrix/versioning"
	"github.com/hashicorp/go-hclog"
	gostream "github.com/libp2p/go-libp2p-gostream"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/orcaman/concurrent-map/v2"
	"io"
	"math/big"
	"math/rand"
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
	txSlotSize = 32 * 1024 // 32kB
)

const (
	DefaultBlockNumSyncDuration  = 2 * time.Second
	DefaultAppStatusSyncDuration = 5 * time.Second
	PocSubmitBatchSize           = 1
	PocSubmitSliceSize           = 1
)

type Endpoint struct {
	logger hclog.Logger

	// gauge for measuring app capacity
	gauge slotGauge
	sync.Mutex
	nextNonce        uint64
	nonceCacheEnable bool

	name       string
	appUrl     string
	appOrigin  string
	h          host.Host
	tag        string
	listener   net.Listener
	httpClient *rpc.FastHttpClient
	signer     Signer
	privateKey *ecdsa.PrivateKey
	address    types.Address
	stream     *eventStream // Event subscriptions

	application     *Application
	minerAgent      *miner.MinerAgent
	jsonRpcClient   *rpc.JsonRpcClient
	blockchainStore blockchainStore

	peersPocRequestMap cmap.ConcurrentMap[string, proof.PocCpuRequest]
	pocQueue           *proof.PocQueue
	pocSubmitQueue     *proof.PocSubmitQueue
	randomNum          int

	//applicationPeersMap cmap.ConcurrentMap[peer.ID, Application]

	// poc_cpu_validate flag
	pocCpuValidateFlag bool

	// poc_gpu_validate flag
	pocGpuValidateFlag bool

	latestBlockHeadHash string
	latestBlockNum      uint64

	isEdgeMode bool
	round      *sd.PocSDRound
	pastRound  *sd.PocSDRound
}

// SubscribeEvents returns a application event subscription
func (b *Endpoint) SubscribeEvents() Subscription {
	return b.stream.subscribe()
}

func (e *Endpoint) EnablePocCpuValidate(flag bool) {
	e.pocCpuValidateFlag = flag
}

func (e *Endpoint) EnablePocGpuValidate(flag bool) {
	e.pocGpuValidateFlag = flag
}

func (e *Endpoint) AddPocTask(
	pocData *proof.PocCpuData,
	priority proof.PocPriority,
) {
	e.pocQueue.AddTask(pocData, priority)
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
			e.logger.Error("doPocTask --> unable to GetNextNonce, %v", err)
			return 0, err
		}
		e.nextNonce = nonce
		e.nonceCacheEnable = true
		return e.nextNonce, nil
	}
	e.nextNonce += 1
	return e.nextNonce, nil
}

func (e *Endpoint) DisableNonceCache() {
	e.Lock()
	defer e.Unlock()

	e.nonceCacheEnable = false
}

func (e *Endpoint) runPocSubmit() {
	batchSize := PocSubmitBatchSize
	sliceSize := PocSubmitSliceSize

	e.logger.Info("runPocSubmit", "batchSize", batchSize)
	batchSubmitData := make([]*proof.PocSubmitData, batchSize)
	taskCount := 0
	for {
		tt := e.pocSubmitQueue.PopTask()

		if tt == nil {
			e.logger.Error("The poc submit queue is closed")
			continue
		}
		batchSubmitData[taskCount] = tt.GetPocSubmitData()
		e.logger.Info("runPocSubmit->batchSubmitData", "count", taskCount, "TargetNodeID", batchSubmitData[taskCount].TargetNodeID, "blockNum", batchSubmitData[taskCount].ValidationTicket, "validator", batchSubmitData[taskCount].Validator, "power", batchSubmitData[taskCount].Power)
		taskCount += 1

		if taskCount < batchSize {
			continue
		}
		e.logger.Info("runPocSubmit->pocSubmitQueue", "remain", e.pocSubmitQueue.Len())
		taskCount = 0
		vecValues := make([]interface{}, len(batchSubmitData))
		for i, pocSubmitData := range batchSubmitData {
			p, err := principal.Decode(pocSubmitData.Validator)
			if err != nil {
				e.logger.Error("principal.Decode", "err", err)
				continue
			}

			vecValues[i] = map[string]interface{}{
				"validationTicket": big.NewInt(pocSubmitData.ValidationTicket),
				"validator":        p,
				"power":            big.NewInt(pocSubmitData.Power),
				"targetNodeID":     pocSubmitData.TargetNodeID,
			}
		}

		// Do batch submit
		sliceBegin := 0
		sliceEnd := 0
		wg := sync.WaitGroup{}
		for sliceBegin < batchSize {
			sliceEnd += sliceSize
			if sliceEnd > batchSize {
				sliceEnd = batchSize
			}
			sliceVecValues := vecValues[sliceBegin:sliceEnd]
			wg.Add(1)
			go func(vec []interface{}) {
				e.submitToIc(vec)
				wg.Done()
			}(sliceVecValues)

			sliceBegin += sliceSize
		}
		wg.Wait()
	}
}

func (e *Endpoint) submitToIc(vecValues []interface{}) {
	e.logger.Info("vecValues", "len", len(vecValues))
	if vecValues == nil || len(vecValues) < 1 {
		return
	}
	// submit proof result to IC canister
	err := e.minerAgent.SubmitValidationVec(vecValues)
	if err != nil {
		e.logger.Warn("\n------------------------------------------\nSubmitValidation", "err", err)
	} else {
		e.logger.Info("\n------------------------------------------\nSubmitValidation", "success", len(vecValues))
	}
}

func (e *Endpoint) doPocTask() {
	for {
		// TODO  get available slot from slotGauge
		tt := e.pocQueue.PopTask()

		if tt == nil {
			// The poc queue is closed,
			// no further poc tasks are incoming
			e.logger.Error("The poc queue is closed")
			return
		}

		pocData := tt.GetPocCpuDataInfo()
		inputString := ""
		if e.pocGpuValidateFlag && e.appUrl != "" && e.appOrigin == proof.AppOriginSD {
			// Do poc by gpu
			pocSD := sd.NewPocSD(e.appUrl)
			prompt := pocData.Seed
			seedNum, _ := pocSD.MakeSeedByHashString(pocData.Seed)
			sdModelHash, md5sum, err := pocSD.ProofByTxt2img(prompt, seedNum)
			if err != nil {
				e.logger.Error("doPocTask -->ProofByTxt2img", "err", err.Error())
			}
			var obj struct {
				NodeId    string `json:"node_id"`
				ModelHash string `json:"model_hash"`
				SeedHash  string `json:"seed_hash"`
				Md5num    string `json:"md5num"`
			}
			obj.NodeId = e.h.ID().String()
			obj.ModelHash = sdModelHash
			obj.Md5num = md5sum
			obj.SeedHash = pocData.Seed

			jsonBuf, err := json.Marshal(obj)
			if err != nil {
				e.logger.Error("doPocTask -->Marshal() error", "err", err.Error())
			}
			inputData := string(jsonBuf)
			inputString = fmt.Sprintf("{\"peerId\": \"%s\",\"endpoint\": \"/poc_gpu_validate\",\"Input\": %s}", pocData.Validator, inputData)
		} else if e.pocCpuValidateFlag {
			// DO poc by cpu
			var data = make(map[string][]byte)
			target := proof.DefaultHashProofTarget
			loops := proof.DefaultHashProofCount
			i := 0
			for i < loops {
				seed := fmt.Sprintf("%s,%d", pocData.Seed, i)
				_, bytes, err := proof.ProofByCalcHash(seed, target, time.Second*5)
				if err != nil {
					e.logger.Error("doPocTask -->ProofByCalcHash", "err", err.Error())
				} else {
					if bytes != nil && len(bytes) > 0 {
						data[seed] = bytes
					} else {
						e.logger.Warn("doPocTask -->ProofByCalcHash failed", "seed", seed)
					}
				}
				i += 1
			}
			if len(data) < proof.DefaultHashProofCount-3 {
				e.logger.Warn("doPocTask -->validate data size too low", "size", len(data))
				continue
			}
			pocDataBuf := "["
			dataIdx := 0
			for k, v := range data {
				pocDataBuf += fmt.Sprintf("{\"k\":\"%s\",\"v\":\"%s\"}", k, hex.EncodeToString(v))
				dataIdx += 1
				if dataIdx < len(data) {
					pocDataBuf += ","
				}
			}
			pocDataBuf += "]"
			inputData := fmt.Sprintf("{\"node_id\" : \"%s\", \"poc_data\": %s}", e.h.ID().String(), pocDataBuf)
			inputString = fmt.Sprintf("{\"peerId\": \"%s\",\"endpoint\": \"/poc_cpu_validate\",\"Input\": %s}", pocData.Validator, inputData)
		}

		if inputString == "" {
			continue
		}

		redoMax := 1
		attemptCount := 0
		callFail := true
		teleResponse := rpc.TelegramResponse{}
		for attemptCount <= redoMax {
			nonce, err := e.GetNextNonce()
			if err != nil {
				e.logger.Error("\"doPocTask -->GetNextNonce", "err:", err.Error())
				attemptCount += 1
				continue
			}
			e.logger.Info(fmt.Sprintf("Calling peer [%s] as validator [%s]", pocData.Validator, e.getID().String()), "queue.len", e.pocQueue.Len(), "nonce", nonce, "data", inputString)
			response, err := e.jsonRpcClient.SendRawTelegram(
				rpc.EdgeCallPrecompile,
				nonce,
				inputString,
				e.privateKey,
			)
			if err != nil {
				e.DisableNonceCache()
				e.logger.Warn("\"doPocTask -->SendRawTelegram for poc", "nonce", nonce, "attemptCount", attemptCount, "err", err.Error())
				if attemptCount >= redoMax {
					break
				}
			} else {
				e.logger.Info("doPocTask -->SendRawTelegram for poc", "TelegramHash", response.Result.TelegramHash, "nonce", nonce, "attemptCount", attemptCount)
				teleResponse = *response
				callFail = false
				break
			}
			attemptCount += 1
		}
		if callFail {
			continue
		}
		respBytes, err := base64.StdEncoding.DecodeString(teleResponse.Result.Response)
		if err != nil {
			e.logger.Warn("doPocTask -->base64 decode err: ", err.Error())
			continue
		}
		var obj struct {
			Message string          `json:"message"`
			Err     json.RawMessage `json:"err"`
		}
		if err := json.Unmarshal(respBytes, &obj); err != nil {
			e.logger.Error("doPocTask -->json.Unmarsha", "err", err.Error())
			continue
		}
		e.logger.Info("doPocTask -->", "message", obj.Message)
		if obj.Err != nil && len(obj.Err) > 0 {
			e.logger.Warn("doPocTask -->", "response.Err", string(obj.Err))
		}
	}
}

//func (e *Endpoint) UpdateApplicationPeer(app *Application) {
//	if app == nil {
//		return
//	}
//	e.applicationPeersMap.Set(app.PeerID, *app)
//}

func (e *Endpoint) GetEndpointApplication() *Application {
	return e.application
}

func (e *Endpoint) doPocRequest() {
	if !e.pocGpuValidateFlag || e.appUrl == "" {
		return
	}
	// check miner status
	_, _, wallet, _, _, err := e.minerAgent.MyNode(e.h.ID().String())
	if err != nil {
		e.logger.Error("doPocRequest -->MyNode", "err", err.Error())
		return
	}
	if wallet == "" {
		e.logger.Info("doPocRequest -->wallet princial=nil")
		return
	}

	// Do test poc by gpu
	pocSD := sd.NewPocSD(e.appUrl)
	randBytes := make([]byte, 32)
	_, err = rand.Read(randBytes)
	if err != nil {
		return
	}
	prompt := hex.EncodeToHex(randBytes)
	seedNum, _ := pocSD.MakeSeedByHashString(prompt)
	sdModelHash, md5sum, err := pocSD.ProofByTxt2img(prompt, seedNum)
	if err != nil {
		e.logger.Error("doPocRequest", "err", err.Error())
		return
	}

	e.logger.Info("doPocRequest", "sdModelHash", sdModelHash)
	if sdModelHash == "" || md5sum == "" {
		e.logger.Error("doPocRequest sd call fail")
		return
	}
	// update application mode hash
	e.application.ModelHash = sdModelHash

	validators, err := e.minerAgent.ListValidatorsNodeId()
	if err != nil {
		e.logger.Error("endpoint.miner -->ListValidatorsNodeId", err.Error())
		return
	}
	//validators := []string{
	//	"16Uiu2HAmGpKZdnpaaYgKTZqagLVJcnMphdeqaHtKBaFFkb5MYRUy",
	//	"16Uiu2HAmTPfBgUkQ4V8qaBvTaJp54Cm32TWGvYZaxcuPxoaSbZAS",
	//	"16Uiu2HAmPfFVHNnYKdDQywJXnzbgM1MdAi6P1MsCkxN7Hr6VaiYa",
	//	"16Uiu2HAmKt7agigzA6oGDdMre4eCU7QER91vrW9M3aneiHEvGu1Y",
	//	"16Uiu2HAmEoDReK7pKygYYYFgJ8uuXS8oWsYFWiiEbCSF9HjYcih2",
	//	"16Uiu2HAkyPw8SEeDpErEwcEZ2QtXzPq5KQf4woWybsr7KN6VH7yX",
	//	"16Uiu2HAm7BqtmjH7JECa5Y4iNgiZXuet3HqXYZbeXNN9XiQwgSbf",
	//}
	if len(validators) < 1 {
		return
	}
	for _, validatorNodeID := range validators {
		// post status by sendTelegram
		redoMax := 1
		attemptCount := 0
		teleResponse := rpc.TelegramResponse{}
		sendOk := false
		for attemptCount <= redoMax {
			nonce, err := e.GetNextNonce()
			if err != nil {
				e.logger.Error("unable to GetNextNonce, %v", err)
				attemptCount += 1
				continue
			}
			inputString := fmt.Sprintf("{\"peerId\": \"%s\",\"endpoint\": \"/poc_cpu_request\",\"Input\": {\"node_id\": \"%s\"}}", validatorNodeID, e.h.ID().String())
			response, err := e.jsonRpcClient.SendRawTelegram(
				rpc.EdgeCallPrecompile,
				nonce,
				inputString,
				e.privateKey,
			)
			if err != nil {
				e.DisableNonceCache()
				e.logger.Warn("endpoint.miner -->SendRawTelegram for doPocRequest", "nonce", nonce, "attemptCount", attemptCount, "input", inputString, "err", err.Error())
				if attemptCount >= redoMax {
					break
				}
				attemptCount += 1
			} else {
				e.logger.Debug("endpoint.miner -->SendRawTelegram for doPocRequest", "TelegramHash", response.Result.TelegramHash, "nonce", nonce, "attemptCount", attemptCount, "input", inputString)
				sendOk = true
				teleResponse = *response
				break
			}
		}
		if !sendOk {
			e.logger.Error("endpoint.miner---->doPocRequest failed", "validatorNodeID", validatorNodeID)
			continue
		}
		e.logger.Debug("endpoint.miner -->SendRawTelegram", "TelegramHash:", teleResponse.Result.TelegramHash)
		decodeBytes, err := base64.StdEncoding.DecodeString(teleResponse.Result.Response)
		if err != nil {
			e.logger.Error("SendRawTelegram", "DecodeString err:", err.Error())
		} else {
			e.logger.Debug("endpoint.miner---->doPocRequest:", "validatorNodeID", validatorNodeID, "resp", string(decodeBytes))
			var obj struct {
				Validator string `json:"validator"`
				Seed      string `json:"seed"`
				Err       string `json:"err"`
			}
			if err := json.Unmarshal(decodeBytes, &obj); err != nil {
				e.logger.Error("endpoint.miner -->json.Unmarshal", "err", err.Error())
				continue
			}
			if obj.Err != "" {
				e.logger.Error("endpoint.miner -->Response", "err", obj.Err)
				continue
			}
			e.AddPocTask(&proof.PocCpuData{
				Validator: obj.Validator,
				Seed:      obj.Seed,
			}, proof.PriorityRequestedPoc)
			e.logger.Info("endpoint.miner -->AddPocTask", "Validator", validatorNodeID, "Seed", obj.Seed)
		}
	}
}
func NewApplicationEndpoint(
	logger hclog.Logger,
	privateKey *ecdsa.PrivateKey,
	srvHost host.Host,
	name string,
	appUrl string,
	appOrigin string,
	blockchainStore blockchainStore,
	minerAgent *miner.MinerAgent,
	jsonRpcClient *rpc.JsonRpcClient,
	isEdgeMode bool) (*Endpoint, error) {
	endpoint := &Endpoint{
		logger:              logger.Named("app_endpoint"),
		name:                name,
		appUrl:              appUrl,
		appOrigin:           appOrigin,
		h:                   srvHost,
		tag:                 ProtoTagEcApp,
		stream:              &eventStream{},
		minerAgent:          minerAgent,
		jsonRpcClient:       jsonRpcClient,
		pocQueue:            proof.NewPocQueue(),
		pocSubmitQueue:      proof.NewPocSubmitQueue(),
		nonceCacheEnable:    false,
		blockchainStore:     blockchainStore,
		peersPocRequestMap:  cmap.New[proof.PocCpuRequest](),
		latestBlockHeadHash: "",
		latestBlockNum:      0,
		isEdgeMode:          isEdgeMode,
		//peersPocRequestMap: PocMap{all: make(map[string]*proof.PocCpuRequest)},
	}
	rand.Seed(time.Now().Unix())
	endpoint.randomNum = rand.Intn(1000)
	endpoint.httpClient = rpc.NewDefaultHttpClient()
	listener, err := gostream.Listen(srvHost, ProtoTagEcApp)
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

	// init application metric
	mac, _ := helper.GetLocalMac()
	endpoint.application = &Application{
		Name:        name,
		PeerID:      srvHost.ID(),
		StartupTime: uint64(time.Now().UnixMilli()),
		Uptime:      0,
		AppOrigin:   appOrigin,
		GuageHeight: 0,
		GuageMax:    200,
		Mac:         mac,
		CpuInfo:     helper.GetCpuInfo(),
		MemInfo:     helper.GetMemInfo(),
		Version:     versioning.Version + " Build" + versioning.Build,
	}

	// TODO check miner status
	go func() {
		ticker := time.NewTicker(DefaultAppStatusSyncDuration)
		for {
			<-ticker.C
			event := &Event{}
			endpoint.application.Uptime = uint64(time.Now().UnixMilli()) - endpoint.application.StartupTime
			endpoint.application.MemInfo = helper.GetMemInfo()

			event.AddNewApp(endpoint.application)
			endpoint.stream.push(event)
			endpoint.logger.Info("endpoint----> status", "ModelHash", endpoint.application.ModelHash, "Mac", endpoint.application.Mac, "CpuInfo", endpoint.application.CpuInfo, "MemInfo", endpoint.application.MemInfo)
		}
		ticker.Stop()
	}()

	go func() {
		ticker := time.NewTicker(proof.DefaultProofDuration)
		//ticker := time.NewTicker(1 * 60 * time.Second) // for test
		for {
			<-ticker.C
			go endpoint.doPocRequest()
		}
		ticker.Stop()
	}()

	go func() {
		ticker := time.NewTicker(DefaultBlockNumSyncDuration)
		for {
			<-ticker.C
			head := endpoint.blockchainStore.Header()
			if head != nil {
				endpoint.latestBlockHeadHash = head.Hash.String()
				endpoint.latestBlockNum = head.Number

				if endpoint.pocGpuValidateFlag {
					currentRoundBlockNumber := (endpoint.latestBlockNum / uint64(proof.DefaultProofBlockMinDuration)) * uint64(proof.DefaultProofBlockMinDuration)
					if endpoint.round == nil {
						endpoint.round = sd.NewPocSDRound(endpoint.logger, currentRoundBlockNumber, endpoint.latestBlockHeadHash)
					}
					round := endpoint.round
					roundBlockNum, _ := round.GetRoundSeed()
					//endpoint.logger.Info("POC current round", "roundBlockNum", roundBlockNum)
					if roundBlockNum != currentRoundBlockNumber {
						endpoint.logger.Info("POC Round change", "newRoundBlockNum", currentRoundBlockNumber)
						endpoint.pastRound = round
						// change to new round
						endpoint.round = sd.NewPocSDRound(endpoint.logger, currentRoundBlockNumber, endpoint.latestBlockHeadHash)
						go func() {
							// complete past round
							validPocSdData, err := endpoint.pastRound.CompleteRound()
							if err != nil {
								endpoint.logger.Info("poc round", "err", err.Error())
							}
							for _, pocSdData := range validPocSdData {
								// valid proof
								endpoint.pocSubmitQueue.AddTask(&proof.PocSubmitData{
									ValidationTicket: int64(pocSdData.BlockNum),
									Validator:        endpoint.minerAgent.GetIdentity(),
									Power:            pocSdData.Power,
									TargetNodeID:     pocSdData.NodeId,
								}, proof.PriorityRequestedPoc)
							}
						}()
					}
				}
			}
		}
		ticker.Stop()
	}()

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
				// ai model hash string
				ModelHash string `json:"model_hash"`
				// mac addr
				Mac string `json:"mac"`
				// memory info
				MemInfo string `json:"mem_info"`
				// cpu info
				CpuInfo string `json:"cpu_info"`
			}
			infoObj.PeerID = endpoint.application.PeerID.String()
			infoObj.Version = endpoint.application.Version
			infoObj.Tag = endpoint.application.AppOrigin
			infoObj.Uptime = uint64(time.Now().UnixMilli()) - endpoint.application.StartupTime
			infoObj.StartupTime = endpoint.application.StartupTime
			infoObj.Name = endpoint.application.Name
			infoObj.CpuInfo = endpoint.application.CpuInfo
			infoObj.MemInfo = endpoint.application.MemInfo
			infoObj.Mac = endpoint.application.Mac
			infoObj.ModelHash = endpoint.application.ModelHash

			info := make([]byte, 0)
			info, err := json.Marshal(infoObj)
			if err != nil {
				info = []byte("endpoint err: " + err.Error())
			}

			writeResponse(w, info, endpoint)
		})

		http.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			resp := fmt.Sprintf("%s", time.Now().String())
			w.Write([]byte(resp))
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
		http.HandleFunc("/poc_cpu_validate", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			pocResp := []byte(fmt.Sprintf("{\"message\":\"validate failed\"}"))
			// Close poc_cpu_validate
			if !endpoint.pocCpuValidateFlag {
				writeResponse(w, pocResp, endpoint)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				pocResp = []byte("endpoint err: " + err.Error())
				writeResponse(w, pocResp, endpoint)
				return
			}
			endpoint.logger.Debug(fmt.Sprintf("/poc_cpu_validate =>request: %s", string(body)))

			var obj struct {
				Node_id  string          `json:"node_id"`
				Poc_data json.RawMessage `json:"poc_data"`
			}
			if err := json.Unmarshal(body, &obj); err != nil {
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "json.Unmarshal obj-> "+err.Error()))
				writeResponse(w, pocResp, endpoint)
				return
			}

			if obj.Node_id == "" || obj.Poc_data == nil || len(obj.Poc_data) < 1 {
				endpoint.logger.Warn("poc_cpu_validate --> invalid poc_data")
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "invalid request"))
				writeResponse(w, pocResp, endpoint)
				return
			}
			pocCpuRequest, ok := endpoint.peersPocRequestMap.Get(obj.Node_id)
			if !ok {
				endpoint.logger.Warn("poc_cpu_validate --> invalid request")
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "invalid request"))
				writeResponse(w, pocResp, endpoint)
				return
			}

			var dataMapJson []map[string]string
			err = json.Unmarshal(obj.Poc_data, &dataMapJson)
			if err != nil {
				endpoint.logger.Warn("poc_cpu_validate --> json.Unmarshal", "resp", string(body), "err", err.Error())
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "invalid request"))
				writeResponse(w, pocResp, endpoint)
				return
			}
			dataMap := make(map[string][]byte)
			for _, data := range dataMapJson {
				bytes, err := hex.DecodeString(data["v"])
				if err != nil {
					endpoint.logger.Warn("poc_cpu_validate --> hex.DecodeString(data[\"v\"]) err: ", err.Error())
					continue
				}
				dataMap[data["k"]] = bytes
			}

			usedTime := time.Since(pocCpuRequest.Start).Milliseconds()
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
			initSeed := pocCpuRequest.Seed
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
			endpoint.logger.Debug(fmt.Sprintf("used time for validate\t\t: %dms", validateUsedTime))
			if rate >= 0.95 {
				// valid proof
				endpoint.pocSubmitQueue.AddTask(&proof.PocSubmitData{
					ValidationTicket: int64(pocCpuRequest.BlockNum),
					Validator:        endpoint.minerAgent.GetIdentity(),
					Power:            usedTime,
					TargetNodeID:     obj.Node_id,
				}, proof.PriorityRequestedPoc)
				result := fmt.Sprintf("validate success\t\t\t: %d/%d rate:%f nodeID:%s", validateSuccess, loops, rate, obj.Node_id)
				endpoint.logger.Info(result)
				pocResp = []byte(fmt.Sprintf("{\"message\":\"SubmitValidation\"}"))
			}
			writeResponse(w, pocResp, endpoint)
		})

		http.HandleFunc("/poc_gpu_validate", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			pocResp := []byte(fmt.Sprintf("{\"message\":\"validate failed\"}"))
			// Close poc_cpu_validate
			if !endpoint.pocGpuValidateFlag {
				writeResponse(w, pocResp, endpoint)
				return
			}

			if endpoint.round == nil {
				writeResponse(w, pocResp, endpoint)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				pocResp = []byte("endpoint err: " + err.Error())
				writeResponse(w, pocResp, endpoint)
				return
			}
			endpoint.logger.Debug(fmt.Sprintf("/poc_gpu_validate =>request: %s", string(body)))

			var obj struct {
				NodeId    string `json:"node_id"`
				ModelHash string `json:"model_hash"`
				SeedHash  string `json:"seed_hash"`
				Md5num    string `json:"md5num"`
			}
			if err := json.Unmarshal(body, &obj); err != nil {
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "json.Unmarshal obj-> "+err.Error()))
				writeResponse(w, pocResp, endpoint)
				return
			}

			if obj.NodeId == "" || obj.SeedHash == "" || obj.Md5num == "" {
				endpoint.logger.Warn("poc_gpu_validate --> invalid poc_data")
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "invalid request"))
				writeResponse(w, pocResp, endpoint)
				return
			}
			pocCpuRequest, ok := endpoint.peersPocRequestMap.Get(obj.NodeId)
			if !ok {
				endpoint.logger.Warn("poc_gpu_validate --> invalid request")
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "invalid request"))
				writeResponse(w, pocResp, endpoint)
				return
			}

			usedTime := time.Since(pocCpuRequest.Start).Milliseconds()
			roundBlockNumber := (pocCpuRequest.BlockNum / uint64(proof.DefaultProofBlockMinDuration)) * uint64(proof.DefaultProofBlockMinDuration)
			pocData := &sd.PocSdData{
				NodeId:    obj.NodeId,
				ModelHash: obj.ModelHash,
				SeedHash:  obj.SeedHash,
				Md5num:    obj.Md5num,
				BlockNum:  roundBlockNumber,
				Power:     usedTime,
			}
			err = endpoint.round.AddPocData(pocData)
			if err != nil {
				endpoint.logger.Warn("poc_gpu_validate --> invalid request", "err", err.Error())
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", err.Error()))
				writeResponse(w, pocResp, endpoint)
				return
			}

			result := fmt.Sprintf("validate submit\t\t\t: NodeId:%s ModelHash:%s SeedHash:%s Md5num:%s Power:%d", obj.NodeId, obj.ModelHash, obj.SeedHash, obj.Md5num, usedTime)
			endpoint.logger.Info(result)
			pocResp = []byte(fmt.Sprintf("{\"message\":\"SubmitValidation\"}"))
			writeResponse(w, pocResp, endpoint)
		})

		http.HandleFunc("/poc_cpu_request", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			pocResp := []byte(fmt.Sprintf("{\"validator\":\"%s\",\"err\":\"%s\"}", endpoint.h.ID().String(), "invalid block time"))

			body, err := io.ReadAll(r.Body)
			if err != nil {
				pocResp = []byte(fmt.Sprintf("{\"validator\":\"%s\",\"seed\":\"%s\",\"err\":\"%s\"}", endpoint.h.ID().String(), "", err.Error()))
				writeResponse(w, pocResp, endpoint)
				return
			}

			endpoint.logger.Info(fmt.Sprintf("/doPocRequest =>request: %s", string(body)))

			if endpoint.latestBlockHeadHash == "" {
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "latestBlockHeadHash is nil"))
				writeResponse(w, pocResp, endpoint)
				return
			}
			isJSON := json.Valid(body)
			if !isJSON {
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "json.Valid body-> "+err.Error()))
				writeResponse(w, pocResp, endpoint)
				return
			}
			var obj struct {
				Node_id string `json:"node_id"`
			}
			if err := json.Unmarshal(body, &obj); err != nil {
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "json.Unmarshal obj-> "+err.Error()))
				writeResponse(w, pocResp, endpoint)
				return
			}
			if obj.Node_id == "" {
				pocResp = []byte(fmt.Sprintf("{\"err\":\"%s\"}", "Node_id is nil"))
				writeResponse(w, pocResp, endpoint)
				return
			}

			// check latest proof number
			blockNumber := endpoint.latestBlockNum
			//var blockNumberFixed uint64 = 0
			blockNumberFixed := (blockNumber / uint64(proof.DefaultProofBlockRange)) * uint64(proof.DefaultProofBlockRange)
			validateRequest := false
			var latestProofNum uint64 = 0

			latestPocCpuRequest, loaded := endpoint.peersPocRequestMap.Get(
				obj.Node_id)
			if !loaded {
				latestProofNum = 0
				validateRequest = true
			} else {
				latestProofNum = latestPocCpuRequest.BlockNum
				validateRequest = (blockNumber - latestProofNum) > uint64(proof.DefaultProofBlockMinDuration)
			}

			if validateRequest {
				seed := endpoint.latestBlockHeadHash
				if endpoint.pocGpuValidateFlag && endpoint.round != nil {
					_, seed = endpoint.round.GetRoundSeed()
				}
				start := time.Now()
				endpoint.peersPocRequestMap.Set(obj.Node_id,
					proof.PocCpuRequest{
						NodeId:   obj.Node_id,
						Seed:     seed,
						BlockNum: blockNumberFixed,
						Start:    start,
					})
				// send proof task to peer node
				pocResp = []byte(fmt.Sprintf("{\"validator\":\"%s\",\"seed\":\"%s\"}", endpoint.h.ID().String(), seed))
				writeResponse(w, pocResp, endpoint)
				return
			} else {
				logger.Warn(fmt.Sprintf("\n\n------------------------------------------\n"+
					"invalid blockNum, blockNumber: %d, latestProofNum: %d, NodeId:%s"+
					"\n------------------------------------------\n", blockNumber, latestProofNum, obj.Node_id))
				pocResp := []byte(fmt.Sprintf("{\"validator\":\"%s\",\"err\":\"%s\"}", endpoint.h.ID().String(), "block num too low"))

				writeResponse(w, pocResp, endpoint)
				return
			}
		})

		server := &http.Server{}
		server.Serve(listener)
	}()

	go endpoint.doPocTask()
	if !endpoint.isEdgeMode {
		go endpoint.runPocSubmit()
	}
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
