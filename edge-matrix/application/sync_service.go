package application

import (
	"context"
	"fmt"
	"github.com/emc-protocol/edge-matrix/application/proof"
	"github.com/emc-protocol/edge-matrix/application/proto"
	"github.com/emc-protocol/edge-matrix/miner"
	"github.com/emc-protocol/edge-matrix/network"
	"github.com/emc-protocol/edge-matrix/network/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p/core/peer"
	"time"
)

type syncAppService struct {
	proto.UnimplementedSyncAppServer
	logger hclog.Logger

	endpoint          *Endpoint
	syncAppPeerClient SyncAppPeerClient
	blockchainStore   blockchainStore
	minerAgent        *miner.MinerAgent
	network           *network.Server
	stream            *grpc.GrpcStream // reference to the grpc stream

	peersBlockNumMap map[peer.ID]uint64
}

type SyncAppPeerService interface {
	// Start starts server
	Start()
	// Close terminates running processes for SyncPeerService
	Close() error
}

func NewSyncAppPeerService(
	logger hclog.Logger,
	network *network.Server,
	endpoint *Endpoint,
	syncAppPeerClient SyncAppPeerClient,
	blockchainStore blockchainStore,
	minerAgent *miner.MinerAgent,
) SyncAppPeerService {
	return &syncAppService{
		logger:            logger,
		network:           network,
		endpoint:          endpoint,
		syncAppPeerClient: syncAppPeerClient,
		blockchainStore:   blockchainStore,
		minerAgent:        minerAgent,
	}
}

// Start starts syncPeerService
func (s *syncAppService) Start() {
	s.registerAppSyncerService()
}

// Close closes syncPeerService
func (s *syncAppService) Close() error {
	return s.stream.Close()
}

// setupGRPCServer setup GRPC server
func (s *syncAppService) registerAppSyncerService() {
	s.stream = grpc.NewGrpcStream()

	proto.RegisterSyncAppServer(s.stream.GrpcServer(), s)
	s.stream.Serve()
	s.network.RegisterProtocol(appSyncerProto, s.stream)
}

func (s *syncAppService) PostAppStatus(req *proto.PostPeerStatusRequest, stream proto.SyncApp_PostAppStatusServer) error {
	nodeIDString := req.GetNodeId()
	s.logger.Info(fmt.Sprintf("\n------------------------------------------\nreq.GetNodeId(): %s", nodeIDString))
	_, err := peer.Decode(nodeIDString)
	if err != nil {
		s.logger.Error("peer.Decode:", err.Error())
	}

	// get latest block number
	//result := "ok"
	//go func() {
	//	header := s.blockchainStore.Header()
	//	if header != nil {
	//		blockNumber := header.Number
	//		// check latest proof number
	//		latestProofNum, ok := s.peersBlockNumMap[PeerID]
	//		if !ok {
	//			latestProofNum = 0
	//		}
	//		var blockNumberFixed uint64 = 0
	//		if (blockNumber - latestProofNum) > proof.DefaultProofBlockMinDuration {
	//			// send proof task to peer node
	//			blockNumberFixed = (blockNumber / proof.DefaultProofBlockRange) * proof.DefaultProofBlockRange
	//			s.peersBlockNumMap[PeerID] = blockNumberFixed // commet this line for disable check blocknum
	//			start := time.Now()
	//
	//			//  get data from peer
	//			s.logger.Info(fmt.Sprintf("\n------------------------------------------\nGetPeerData: %s", nodeIDString))
	//			dataMap, err := s.syncAppPeerClient.GetPeerData(PeerID, header.Hash.String(), time.Second*60)
	//			if err != nil {
	//				s.logger.Error("GetPeerData", "PeerID", PeerID, "err", err.Error())
	//			} else {
	//				usedTime := time.Since(start).Milliseconds()
	//
	//				// validate data
	//				//if s.logger.IsDebug() {
	//				//	s.logger.Debug("PeerData: {")
	//				//	for dataKey, bytes := range dataMap {
	//				//		s.logger.Debug(dataKey, hex.EncodeToString(bytes))
	//				//	}
	//				//	s.logger.Debug("}")
	//				//}
	//
	//				var hashArray = make([]string, proof.DefaultHashProofCount)
	//				target := proof.DefaultHashProofTarget
	//				loops := proof.DefaultHashProofCount
	//				i := 0
	//				initSeed := header.Hash.String()
	//				for i < loops {
	//					seed := fmt.Sprintf("%s,%d", initSeed, i)
	//					hashArray[i] = seed
	//					i += 1
	//				}
	//
	//				validateSuccess := 0
	//				validateStart := time.Now()
	//				for _, hash := range hashArray {
	//					isValidate := proof.ValidateHash(hash, target, dataMap[hash])
	//					if isValidate {
	//						validateSuccess += 1
	//					}
	//				}
	//
	//				validateUsedTime := time.Since(validateStart).Milliseconds()
	//				rate := float32(validateSuccess) / float32(proof.DefaultHashProofCount)
	//				s.logger.Debug(fmt.Sprintf("used time for validate\t\t: %dms", validateUsedTime))
	//				result = fmt.Sprintf("validate success\t\t\t: %d/%d rate:%f nodeID:%s", validateSuccess, loops, rate, nodeIDString)
	//				s.logger.Info(result)
	//				if rate >= 0.95 {
	//					// valid proof
	//					s.logger.Info("\n------------------------------------------\nSubmit proof to IC", "usedTime(ms)", usedTime, "blockNumber", blockNumberFixed, "NodeID", nodeIDString)
	//					// submit proof result to IC canister
	//					err := s.minerAgent.SubmitValidation(
	//						int64(blockNumberFixed),
	//						s.minerAgent.GetIdentity(),
	//						usedTime,
	//						nodeIDString,
	//					)
	//					if err != nil {
	//						s.logger.Error("\n------------------------------------------\nSubmitValidation:", "err", err)
	//					}
	//				}
	//			}
	//		} else {
	//			s.logger.Warn(fmt.Sprintf("\n------------------------------------------\ninvalid blockNum: %d, NodeId:%s", blockNumberFixed, nodeIDString))
	//		}
	//	}
	//}()

	// if client closes stream, context.Canceled is given
	//if err := stream.Send(&proto.Result{Data: result}); err != nil {
	//	s.logger.Error("\n------------------------------------------\nstream.Send:", "err", err)
	//	return nil
	//}
	return nil
}

// GetData is a gRPC endpoint to return Data
func (s *syncAppService) GetData(
	req *proto.GetDataRequest,
	stream proto.SyncApp_GetDataServer,
) error {
	var data = make(map[string][]byte)
	target := proof.DefaultHashProofTarget
	loops := proof.DefaultHashProofCount
	i := 0
	for i < loops {
		seed := fmt.Sprintf("%s,%d", req.GetDataHash(), i)
		_, bytes, err := proof.ProofByCalcHash(seed, target, time.Second*5)
		if err != nil {
			break
		}
		data[seed] = bytes
		i += 1
	}

	// if client closes stream, context.Canceled is given
	if err := stream.Send(toProtoData(data)); err != nil {
		return nil
	}

	return nil
}

// GetStatus is a gRPC endpoint to return the latest block number as a node status
func (s *syncAppService) GetStatus(
	ctx context.Context,
	req *empty.Empty,
) (*proto.AppStatus, error) {
	return &proto.AppStatus{
		Name:        s.endpoint.application.Name,
		StartupTime: s.endpoint.application.StartupTime,
		Uptime:      s.endpoint.application.Uptime,
		GuageMax:    s.endpoint.application.GuageMax,
		GuageHeight: s.endpoint.application.GuageHeight,
	}, nil
}

func toProtoData(data map[string][]byte) *proto.Data {
	return &proto.Data{Data: data}
}
