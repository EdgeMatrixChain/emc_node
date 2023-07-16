package relay

import (
	"context"
	"errors"
	"github.com/emc-protocol/edge-matrix/application"
	appProto "github.com/emc-protocol/edge-matrix/application/proto"
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
type networkingServer interface {

	// GetPeerAddrInfo fetches the AddrInfo of a peer
	GetPeerAddrInfo(peerID peer.ID) peer.AddrInfo
}

// BOOTNODE QUERIES //
// AliveService is a service that finds other peers in the network
// and connects them to the current running node
type AliveService struct {
	proto.UnimplementedAliveServer
	pendingPeerConnections sync.Map // Map that keeps track of the pending status of peers; peerID -> bool

	baseServer   networkingServer // The interface towards the base networking server
	logger       hclog.Logger     // The AliveService logger
	routingTable *kb.RoutingTable // Kademlia 'k-bucket' routing table that contains connected nodes info

	syncAppPeerClient application.SyncAppPeerClient
	closeCh           chan struct{} // Channel used for stopping the AliveService
}

// NewAliveService creates a new instance of the alive service
func NewAliveService(
	server networkingServer,
	routingTable *kb.RoutingTable,
	logger hclog.Logger,
	syncAppPeerClient application.SyncAppPeerClient,
) *AliveService {
	return &AliveService{
		logger:            logger.Named("discovery"),
		baseServer:        server,
		routingTable:      routingTable,
		syncAppPeerClient: syncAppPeerClient,
		closeCh:           make(chan struct{}),
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
	addr := ""
	addrInfo := d.baseServer.GetPeerAddrInfo(from)
	if len(addrInfo.Addrs) > 0 {
		addr = addrInfo.Addrs[0].String()
	}
	d.logger.Debug("-------->Alive status", "from", from, "name", status.Name, "app_origin", status.AppOrigin, "addr", addr, "relay", status.Relay)
	d.syncAppPeerClient.PublishApplicationStatus(&appProto.AppStatus{
		Name:        status.Name,
		NodeId:      from.String(),
		Uptime:      status.Uptime,
		StartupTime: status.StartupTime,
		Relay:       status.Relay,
		Addr:        addr,
		AppOrigin:   status.AppOrigin,
		Mac:         status.Mac,
		CpuInfo:     status.CpuInfo,
		MemInfo:     status.MemInfo,
		ModelHash:   status.ModelHash,
	})

	return &proto.AliveStatusResp{
		Success: true,
	}, nil
}
