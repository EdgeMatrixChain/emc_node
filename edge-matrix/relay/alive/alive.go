package alive

import (
	"context"
	"errors"
	"github.com/emc-protocol/edge-matrix/network/grpc"
	"github.com/emc-protocol/edge-matrix/relay/proto"
	"github.com/hashicorp/go-hclog"
	kb "github.com/libp2p/go-libp2p-kbucket"
	"github.com/libp2p/go-libp2p/core/peer"
	"sync"
)

const (
	EdgeAliveProto = "/alive/0.2"
)

// networkingServer defines the base communication interface between
// any networking server implementation and the AliveService

// AliveService is a service that finds other peers in the network
// and connects them to the current running node
type AliveService struct {
	proto.UnimplementedAliveServer
	pendingPeerConnections sync.Map // Map that keeps track of the pending status of peers; peerID -> bool

	//baseServer   networkingServer // The interface towards the base networking server
	logger       hclog.Logger     // The AliveService logger
	routingTable *kb.RoutingTable // Kademlia 'k-bucket' routing table that contains connected nodes info

	closeCh chan struct{} // Channel used for stopping the AliveService
}

// NewAliveService creates a new instance of the alive service
func NewAliveService(
	//server networkingServer,
	routingTable *kb.RoutingTable,
	logger hclog.Logger,
) *AliveService {
	return &AliveService{
		logger: logger.Named("discovery"),
		//baseServer:   server,
		routingTable: routingTable,
		closeCh:      make(chan struct{}),
	}
}

// Close stops the discovery service
func (d *AliveService) Close() {
	close(d.closeCh)
}

// RoutingTableSize returns the size of the routing table
func (d *AliveService) RoutingTableSize() int {
	return d.routingTable.Size()
}

// RoutingTablePeers fetches the peers from the routing table
func (d *AliveService) RoutingTablePeers() []peer.ID {
	return d.routingTable.ListPeers()
}

func (d *AliveService) Hello(ctx context.Context, status *proto.AliveStatus) (*proto.AliveStatusResp, error) {
	// Extract the requesting peer ID from the gRPC context
	grpcContext, ok := ctx.(*grpc.Context)
	if !ok {
		return nil, errors.New("invalid type assertion")
	}

	from := grpcContext.PeerID
	d.logger.Info("-------->Receive Hello", "from", from, "stats", status.String())

	//if info := d.baseServer.GetPeerInfo(id); len(info.Addrs) > 0 {
	//	filteredPeers = append(filteredPeers, common.AddrInfoToString(info))
	//}

	return &proto.AliveStatusResp{
		Success: true,
	}, nil
}
