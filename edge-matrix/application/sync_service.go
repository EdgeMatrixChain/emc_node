package application

import (
	"context"
	"fmt"
	"github.com/emc-protocol/edge-matrix/application/proof"
	"github.com/emc-protocol/edge-matrix/application/proto"
	"github.com/emc-protocol/edge-matrix/network"
	"github.com/emc-protocol/edge-matrix/network/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"time"
)

type syncAppService struct {
	proto.UnimplementedSyncAppServer

	endpoint *Endpoint
	network  *network.Server
	stream   *grpc.GrpcStream // reference to the grpc stream
}

type SyncAppPeerService interface {
	// Start starts server
	Start()
	// Close terminates running processes for SyncPeerService
	Close() error
}

func NewSyncAppPeerService(
	network *network.Server,
	endpoint *Endpoint,
) SyncAppPeerService {
	return &syncAppService{
		network:  network,
		endpoint: endpoint,
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

// toProtoBlock converts type.Block -> proto.Block
func toProtoData(data map[string][]byte) *proto.Data {
	return &proto.Data{Data: data}
}
