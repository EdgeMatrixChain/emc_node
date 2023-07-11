package application

import (
	"context"
	"github.com/emc-protocol/edge-matrix/network"
	"github.com/emc-protocol/edge-matrix/network/event"
	"github.com/libp2p/go-libp2p/core/peer"
	rawGrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"math/big"
)

type Network interface {
	// AddrInfo returns Network Info
	AddrInfo() *peer.AddrInfo
	// RegisterProtocol registers gRPC service
	RegisterProtocol(string, network.Protocol)
	// Peers returns current connected peers
	Peers() []*network.PeerConnInfo
	// SubscribeCh returns a channel of peer event
	SubscribeCh(context.Context) (<-chan *event.PeerEvent, error)
	// GetPeerDistance returns the distance between the node and given peer
	GetPeerDistance(peer.ID) *big.Int
	// NewProtoConnection opens up a new stream on the set protocol to the peer,
	// and returns a reference to the connection
	NewProtoConnection(protocol string, peerID peer.ID) (*rawGrpc.ClientConn, error)
	// NewTopic Creates New Topic for gossip
	NewTopic(protoID string, obj proto.Message) (*network.Topic, error)
	// IsConnected returns the node is connecting to the peer associated with the given ID
	IsConnected(peerID peer.ID) bool
	// SaveProtocolStream saves stream
	SaveProtocolStream(protocol string, stream *rawGrpc.ClientConn, peerID peer.ID)
	// CloseProtocolStream closes stream
	CloseProtocolStream(protocol string, peerID peer.ID) error
}

type ApplicationStore interface {
	// ApplicationStore returns the application of endpoint
	GetEndpointApplication() *Application
	// UpdateApplicationPeer set/add application to applicationPeers map
	//UpdateApplicationPeer(app *Application)
}
