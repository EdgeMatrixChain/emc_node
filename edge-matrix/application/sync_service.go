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

	applicationStore ApplicationStore
	blockchainStore  blockchainStore
	minerAgent       *miner.MinerHubAgent
	network          *network.Server
	stream           *grpc.GrpcStream // reference to the grpc stream

	//peersBlockNumMap map[peer.ID]uint64
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
	applicationStore ApplicationStore,
	blockchainStore blockchainStore,
	minerAgent *miner.MinerHubAgent,
) SyncAppPeerService {
	return &syncAppService{
		logger:           logger,
		network:          network,
		applicationStore: applicationStore,
		blockchainStore:  blockchainStore,
		minerAgent:       minerAgent,
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

// GetStatus is a gRPC endpoint to return the latest  application status
func (s *syncAppService) GetStatus(
	ctx context.Context,
	req *empty.Empty,
) (*proto.AppStatus, error) {
	application := s.applicationStore.GetEndpointApplication()
	return &proto.AppStatus{
		Name:        application.Name,
		StartupTime: application.StartupTime,
		Uptime:      application.Uptime,
		GuageMax:    application.GuageMax,
		GuageHeight: application.GuageHeight,
	}, nil
}

func toProtoData(data map[string][]byte) *proto.Data {
	return &proto.Data{Data: data}
}
