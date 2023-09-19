package relay

import (
	"context"
	"errors"
	"github.com/emc-protocol/edge-matrix/application"
	appProto "github.com/emc-protocol/edge-matrix/application/proto"
	"github.com/emc-protocol/edge-matrix/network/common"
	"github.com/emc-protocol/edge-matrix/network/grpc"
	"github.com/emc-protocol/edge-matrix/relay/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"regexp"
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

	GetRandomBootnode() *peer.AddrInfo
}

// BOOTNODE QUERIES //
// AliveService is a service that finds other peers in the network
// and connects them to the current running node
type AliveService struct {
	proto.UnimplementedAliveServer
	pendingPeerConnections sync.Map // Map that keeps track of the pending status of peers; peerID -> bool

	baseServer networkingServer // The interface towards the base networking server
	logger     hclog.Logger     // The AliveService logger
	//routingTable *kb.RoutingTable // Kademlia 'k-bucket' routing table that contains connected nodes info

	syncAppPeerClient application.SyncAppPeerClient
	closeCh           chan struct{} // Channel used for stopping the AliveService
}

// NewAliveService creates a new instance of the alive service
func NewAliveService(
	server networkingServer,
	//routingTable *kb.RoutingTable,
	logger hclog.Logger,
	syncAppPeerClient application.SyncAppPeerClient,
) *AliveService {
	return &AliveService{
		logger:     logger.Named("AliveService"),
		baseServer: server,
		//routingTable:      routingTable,
		syncAppPeerClient: syncAppPeerClient,
		closeCh:           make(chan struct{}),
	}
}

// Close stops the discovery service
func (d *AliveService) Close() {
	close(d.closeCh)
}

// RoutingTableSize returns the size of the routing table
//func (d *AliveService) RoutingTableSize() int {
//	return d.routingTable.Size()
//}

// RoutingTablePeers fetches the peers from the routing table
//func (d *AliveService) RoutingTablePeers() []peer.ID {
//	return d.routingTable.ListPeers()
//}

func (d *AliveService) Hello(ctx context.Context, status *proto.AliveStatus) (*proto.AliveStatusResp, error) {
	// Extract the requesting peer ID from the gRPC context
	grpcContext, ok := ctx.(*grpc.Context)
	if !ok {
		return nil, errors.New("invalid type assertion")
	}

	from := grpcContext.PeerID
	addr := ""
	innerIp := false
	addrInfo := d.baseServer.GetPeerAddrInfo(from)
	if len(addrInfo.Addrs) > 0 {
		addr = addrInfo.Addrs[0].String()
		innerIp = isInnerIp(addrInfo.Addrs[0])
	}
	d.logger.Debug("-------->Alive status", "from", from, "name", status.Name, "app_origin", status.AppOrigin, "addr", addr, "relay", status.Relay)
	if !innerIp || status.Relay != "" {
		d.syncAppPeerClient.PublishApplicationStatus(&appProto.AppStatus{
			Name:         status.Name,
			NodeId:       from.String(),
			Uptime:       status.Uptime,
			StartupTime:  status.StartupTime,
			Relay:        status.Relay,
			Addr:         addr,
			AppOrigin:    status.AppOrigin,
			Mac:          status.Mac,
			CpuInfo:      status.CpuInfo,
			GpuInfo:      status.GpuInfo,
			MemInfo:      status.MemInfo,
			ModelHash:    status.ModelHash,
			AveragePower: status.AveragePower,
			Version:      status.Version,
		})
	}

	newRelayNode := d.baseServer.GetRandomBootnode()
	discovery := ""
	if newRelayNode != nil {
		discovery = common.AddrInfoToString(newRelayNode)
	}
	return &proto.AliveStatusResp{
		Success:   true,
		Discovery: discovery,
	}, nil
}

func isInnerIp(ma multiaddr.Multiaddr) (innerIp bool) {
	innerIp = false
	ip4Addr, err := ma.ValueForProtocol(multiaddr.P_IP4)
	if err != nil {
		ip4Addr = ""
	}

	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)\.(\d+)`)
	submatches := re.FindStringSubmatch(ip4Addr)
	if len(submatches) > 0 {
		// 127.0.0.1
		// 10.0.0.0/8
		// 172.16.0.0/12
		// 169.254.0.0/16
		// 192.168.0.0/16
		if submatches[0] == "127.0.0.1" ||
			submatches[1] == "10" ||
			(submatches[1] == "172" && submatches[2] >= "16" && submatches[2] <= "31") ||
			(submatches[1] == "169" && submatches[2] == "254") ||
			(submatches[1] == "192" && submatches[2] == "168") {
			innerIp = true
		}
	}
	return innerIp
}
