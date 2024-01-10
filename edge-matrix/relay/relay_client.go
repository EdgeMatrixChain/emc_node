package relay

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/emc-protocol/edge-matrix/application"
	emcNetwork "github.com/emc-protocol/edge-matrix/network"
	"github.com/emc-protocol/edge-matrix/network/common"
	"github.com/emc-protocol/edge-matrix/network/grpc"
	"github.com/emc-protocol/edge-matrix/relay/proto"
	"github.com/emc-protocol/edge-matrix/secrets"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/client"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	"github.com/multiformats/go-multiaddr"
	rawGrpc "google.golang.org/grpc"
	"math/big"
	"sync"
	"time"
)

const (
	// maxDiscoveryPeerReqCount is the max peer count that
	// can be requested from other peers

	// bootnodeAliveInterval is the interval at which
	// random bootnodes are dialed for their peer sets
	bootnodeAliveInterval = 15 * 60 * time.Second
)

// RelayConnInfo holds the connection information about the peer
type RelayConnInfo struct {
	Info peer.AddrInfo

	connDirections  map[network.Direction]bool
	protocolStreams map[string]*rawGrpc.ClientConn
}

type RelayClient struct {
	logger hclog.Logger // the logger

	subscription application.Subscription // reference to the application subscription
	closeCh      chan struct{}            // the channel used for closing the RelayClient

	host host.Host // the libp2p host reference

	protocols     map[string]Protocol // supported protocols
	protocolsLock sync.Mutex          // lock for the supported protocols map

	peers     map[peer.ID]*PeerConnInfo // map of all peer connections
	peersLock sync.Mutex                // lock for the peer map

	relayPeers     map[peer.ID]*RelayPeerInfo // map of all relay peer connections
	relayPeersLock sync.Mutex                 // lock for the relay peer map

	relaynodes *relaynodesWrapper // reference of all relaynodes for the node

	application *application.Application // reference of application
}

// RelayPeerInfo holds the relay information about the peer
type RelayPeerInfo struct {
	Info        *RelayConnInfo
	reservation *client.Reservation
}

func (s *RelayClient) GetBootnodes() []*peer.AddrInfo {
	return s.relaynodes.getRelaynodes()
}

func (s *RelayClient) GetRandomBootnode() *peer.AddrInfo {
	nonConnectedNodes := make([]*peer.AddrInfo, 0)

	for _, v := range s.relaynodes.getRelaynodes() {
		//if !s.hasPeer(v.ID) {
		nonConnectedNodes = append(nonConnectedNodes, v)
		//}
	}

	if len(nonConnectedNodes) > 0 {
		randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(nonConnectedNodes))))

		return nonConnectedNodes[randNum.Int64()]
	}

	return nil
}

// setupRelaynodes sets up the node's relayer node connections
func (s *RelayClient) setupRelaynodes(relaynodes []string) error {
	// Check the relaynode config is present
	if relaynodes == nil {
		return ErrNoRelaynodes
	}

	// Check if at least one relaynode is specified
	if len(relaynodes) < MinimumRelayNodes {
		return ErrMinRelaynodes
	}

	relaynodesArr := make([]*peer.AddrInfo, 0)
	relaynodesMap := make(map[peer.ID]*peer.AddrInfo)

	for _, rawAddr := range relaynodes {
		bootnode, err := common.StringToAddrInfo(rawAddr)
		if err != nil {
			return fmt.Errorf("failed to parse relaynode %s: %w", rawAddr, err)
		}

		if bootnode.ID == s.host.ID() {
			s.logger.Info("Omitting relaynode with same ID as host", "id", bootnode.ID)

			continue
		}

		relaynodesArr = append(relaynodesArr, bootnode)
		relaynodesMap[bootnode.ID] = bootnode
	}

	s.relaynodes = &relaynodesWrapper{
		relaynodeArr:       relaynodesArr,
		relaynodesMap:      relaynodesMap,
		relaynodeConnCount: 0,
	}

	return nil
}

// addRelaynodes add relayer nodes
func (s *RelayClient) addRelaynodes(relaynodes []string) error {
	// Check the relaynode config is present
	if relaynodes == nil {
		return ErrNoRelaynodes
	}

	// Check if at least one relaynode is specified
	if len(relaynodes) < 1 {
		return ErrNoRelaynodes
	}

	relaynodesArr := s.relaynodes.relaynodeArr
	relaynodesMap := s.relaynodes.relaynodesMap

	for _, rawAddr := range relaynodes {
		bootnode, err := common.StringToAddrInfo(rawAddr)
		if err != nil {
			return fmt.Errorf("failed to parse relaynode %s: %w", rawAddr, err)
		}

		if bootnode.ID == s.host.ID() {
			s.logger.Info("Omitting relaynode with same ID as host", "id", bootnode.ID)

			continue
		}

		if s.hasRelayPeer(bootnode.ID) {
			continue
		}

		relaynodesArr = append(relaynodesArr, bootnode)
		relaynodesMap[bootnode.ID] = bootnode
	}

	s.relaynodes = &relaynodesWrapper{
		relaynodeArr:       relaynodesArr,
		relaynodesMap:      relaynodesMap,
		relaynodeConnCount: 0,
	}

	return nil
}

// numPeers returns the number of connected reply peers [Thread safe]
func (s *RelayClient) numRelayPeers() int64 {
	s.relayPeersLock.Lock()
	defer s.relayPeersLock.Unlock()

	return int64(len(s.relayPeers))
}

// RelayPeers returns a copy of the networking server's relay peer connection info set.
// Only one (initial) connection (inbound OR outbound) per peer is contained [Thread safe]
func (s *RelayClient) RelayPeers() []*RelayConnInfo {
	s.relayPeersLock.Lock()
	defer s.relayPeersLock.Unlock()

	peers := make([]*RelayConnInfo, 0)
	for _, relayPeer := range s.relayPeers {
		if relayPeer.Info != nil {
			peers = append(peers, relayPeer.Info)
		}
	}

	return peers
}

// hasRelayPeer checks if the peer is present in the relay peers list [Thread safe]
func (s *RelayClient) hasRelayPeer(peerID peer.ID) bool {
	s.relayPeersLock.Lock()
	defer s.relayPeersLock.Unlock()

	_, ok := s.relayPeers[peerID]

	return ok
}

func (s *RelayClient) keepAliveMinimumRelayConnections() {
	for {
		select {
		case <-time.After(5 * time.Second):
		case <-s.closeCh:
			return
		}

		relayPeersSnapshot := make(map[peer.ID]*RelayPeerInfo, 0)
		for _, relayPeerInfo := range s.relayPeers {
			relayPeersSnapshot[relayPeerInfo.Info.Info.ID] = relayPeerInfo
		}
		for _, relayPeerInfo := range relayPeersSnapshot {
			if relayPeerInfo != nil {
				if relayPeerInfo.reservation.Expiration.Before(time.Now()) {
					// dial random unconnected relaynode
					s.connectToRandomRelayNodes()

					// disconncet expired relaynode
					s.removeRelayPeerInfo(relayPeerInfo.Info.Info.ID)
					s.disconnectFromPeer(relayPeerInfo.Info.Info.ID, "reconnect for reservation")
					s.RemoveFromPeerStore(&relayPeerInfo.Info.Info)

					// update alive status
					go s.keepAliveToBootnodes()
				}
			}
		}

		if s.numRelayPeers() < MinimumRelayConnections {
			// dial random unconnected relaynode
			s.connectToRandomRelayNodes()
			// update alive status
			go s.keepAliveToBootnodes()
		}
	}
}

func (s *RelayClient) connectToRandomRelayNodes() {
	if relayinfo := s.GetRandomRelaynode(); relayinfo != nil {
		s.logger.Info("connectToRandomRelayNodes", "relayinfo", relayinfo.String())

		resv, err := client.Reserve(context.Background(), s.host, *relayinfo)
		if err != nil {
			s.logger.Error(fmt.Sprintf("privateSrvHost failed to receive a relay reservation from %v. %v", relayinfo.Addrs[0], err))
			var re client.ReservationError
			if !errors.As(err, &re) {
				s.logger.Error("expected error to be of type %T", re)
			}
			s.logger.Error(re.Error())
		} else {
			s.addRelayPeerInfo(relayinfo, network.DirUnknown, resv)
			s.logger.Info(fmt.Sprintf("reservation: LimitData=%d, LimitDuration=%v, Expiration=%v, Addrs=%v", resv.LimitData, resv.LimitDuration, resv.Expiration, resv.Addrs))
		}
	}
}

// GetRandomRelaynode fetches a random relaynode that's currently
// NOT connected, if any
func (s *RelayClient) GetRandomRelaynode() *peer.AddrInfo {
	nonConnectedNodes := make([]*peer.AddrInfo, 0)

	for _, v := range s.relaynodes.getRelaynodes() {
		if !s.hasRelayPeer(v.ID) {
			nonConnectedNodes = append(nonConnectedNodes, v)
		}
	}

	if len(nonConnectedNodes) > 0 {
		randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(nonConnectedNodes))))

		return nonConnectedNodes[randNum.Int64()]
	}

	return nil
}

func (s *RelayClient) GetRelayPeerInfo() *RelayPeerInfo {
	for _, relayInfo := range s.relayPeers {
		return relayInfo
	}

	return nil
}

func (s *RelayClient) addRelayPeerInfo(relayInfo *peer.AddrInfo, direction network.Direction, resv *client.Reservation) bool {
	s.relayPeersLock.Lock()
	defer s.relayPeersLock.Unlock()

	relayPeerInfo, relayPeerExists := s.relayPeers[relayInfo.ID]
	if relayPeerExists && relayPeerInfo.Info != nil && relayPeerInfo.Info.connDirections[direction] {
		// Check if this peer already has an active connection status (saved info).
		// There is no need to do further processing
		return true
	}

	// Check if the connection info is already initialized
	if !relayPeerExists {
		// Create a new record for the connection info
		relayPeerInfo = &RelayPeerInfo{
			Info: &RelayConnInfo{
				Info:            *relayInfo,
				connDirections:  make(map[network.Direction]bool),
				protocolStreams: make(map[string]*rawGrpc.ClientConn),
			},
		}
	}

	// Save the connection info to the networking server
	relayPeerInfo.Info.connDirections[direction] = true
	relayPeerInfo.reservation = resv
	s.relayPeers[relayInfo.ID] = relayPeerInfo

	return false
}

// removeRelayPeerInfo removes (pops) relay peer connection info from the networking
// server's relay peer map. Returns nil if no peer was removed
func (s *RelayClient) removeRelayPeerInfo(peerID peer.ID) *RelayPeerInfo {
	s.relayPeersLock.Lock()
	defer s.relayPeersLock.Unlock()

	// Remove the peer from the peers map
	relayPeerInfo, ok := s.relayPeers[peerID]
	if !ok {
		// Peer is not present in the relay peers map
		return nil
	}

	s.logger.Warn("removeRelayPeerInfo", "ID", peerID)
	// Delete the peer from the relay peers map
	delete(s.relayPeers, peerID)

	return relayPeerInfo
}

// Start starts the networking relay reserve job
func (s *RelayClient) StartRelayReserv() error {
	s.logger.Info(" LibP2P relay reserve job running")

	go s.keepAliveMinimumRelayConnections()

	// watch for disconnected relay peers
	s.host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: func(net network.Network, conn network.Conn) {
			s.removeRelayPeerInfo(conn.RemotePeer())
		},
	})

	return nil
}

func (s *RelayClient) GetHost() host.Host {
	return s.host
}

// setupLibp2pKey is a helper method for setting up the networking private key
func setupLibp2pKey(secretsManager secrets.SecretsManager) (crypto.PrivKey, error) {
	var key crypto.PrivKey

	if secretsManager.HasSecret(secrets.NetworkKey) {
		// The key is present in the secrets manager, read it
		networkingKey, readErr := emcNetwork.ReadLibp2pKey(secretsManager)
		if readErr != nil {
			return nil, fmt.Errorf("unable to read networking private key from Secrets Manager, %w", readErr)
		}

		key = networkingKey
	} else {
		// The key is not present in the secrets manager, generate it
		libp2pKey, libp2pKeyEncoded, keyErr := emcNetwork.GenerateAndEncodeLibp2pKey()
		if keyErr != nil {
			return nil, fmt.Errorf("unable to generate networking private key for Secrets Manager, %w", keyErr)
		}

		// Write the networking private key to disk
		if setErr := secretsManager.SetSecret(secrets.NetworkKey, libp2pKeyEncoded); setErr != nil {
			return nil, fmt.Errorf("unable to store networking private key to Secrets Manager, %w", setErr)
		}

		key = libp2pKey
	}

	return key, nil
}

// setupAlive Sets up the live service for the node
func (s *RelayClient) StartAlive(subscription application.Subscription) error {
	// Set up a fresh routing table
	//keyID := kb.ConvertPeerID(s.host.ID())
	//
	//routingTable, err := kb.NewRoutingTable(
	//	defaultBucketSize,
	//	keyID,
	//	time.Minute,
	//	s.host.Peerstore(),
	//	10*time.Second,
	//	nil,
	//)
	//if err != nil {
	//	return err
	//}
	//
	//// Set the PeerAdded event handler
	//routingTable.PeerAdded = func(p peer.ID) {
	//	//info := s.host.Peerstore().PeerInfo(p)
	//	//s.addToDialQueue(&info, common.PriorityRandomDial)
	//}
	//
	//// Set the PeerRemoved event handler
	//routingTable.PeerRemoved = func(p peer.ID) {
	//	//s.dialQueue.DeleteTask(p)
	//}

	// Register the network notify bundle handlers
	s.host.Network().Notify(s.GetNotifyBundle())

	// Make sure the alive service has the bootnodes in its routing table,
	// and instantiates connections to them
	s.ConnectToBootnodes(s.relaynodes.getRelaynodes())

	// Start application event update process
	go s.startApplicationEventProcess(subscription)

	// Start the alive job
	go s.startAliveService()
	return nil
}

func (s *RelayClient) GetNotifyBundle() *network.NotifyBundle {
	return &network.NotifyBundle{
		ConnectedF: func(net network.Network, conn network.Conn) {
			peerID := conn.RemotePeer()
			s.logger.Info("Conn", "peer", peerID, "direction", conn.Stat().Direction)
			// Update the peer connection info
			if connectionExists := s.addPeerInfo(peerID, conn.Stat().Direction); connectionExists {
				// The peer connection information was already present in the networking
				// server, so no connection metrics should be updated further
				return
			}
		},
	}
}

// addPeerInfo updates the networking server's internal peer info table
// and returns a flag indicating if the same peer connection previously existed.
// In case the peer connection previously existed, this is a noop
func (s *RelayClient) addPeerInfo(id peer.ID, direction network.Direction) bool {
	s.peersLock.Lock()
	defer s.peersLock.Unlock()

	connectionInfo, connectionExists := s.peers[id]
	if connectionExists && connectionInfo.connDirections[direction] {
		// Check if this peer already has an active connection status (saved info).
		// There is no need to do further processing
		return true
	}

	// Check if the connection info is already initialized
	if !connectionExists {
		// Create a new record for the connection info
		connectionInfo = &PeerConnInfo{
			Info:            s.host.Peerstore().PeerInfo(id),
			connDirections:  make(map[network.Direction]bool),
			protocolStreams: make(map[string]*rawGrpc.ClientConn),
		}
	}

	// Save the connection info to the networking server
	connectionInfo.connDirections[direction] = true

	s.peers[id] = connectionInfo

	return false
}

// AddToPeerStore adds peer information to the node's peer store
func (s *RelayClient) AddToPeerStore(peerInfo *peer.AddrInfo) {
	s.host.Peerstore().AddAddr(peerInfo.ID, peerInfo.Addrs[0], peerstore.PermanentAddrTTL)

}

func (s *RelayClient) RemoveFromPeerStore(peerInfo *peer.AddrInfo) {
	s.host.Peerstore().RemovePeer(peerInfo.ID)
}

// addToTable adds the node to the peer store and the routing table
func (d *RelayClient) addToTable(node *peer.AddrInfo) error {
	// before we include peers on the routing table -> dial queue
	// we have to add them to the peer store so that they are
	// available to all the libp2p services
	d.AddToPeerStore(node)
	//d.logger.Debug("service-->addToTable", "node", node.String())
	//if _, err := d.routingTable.TryAddPeer(
	//	node.ID,
	//	false,
	//	false,
	//); err != nil {
	//	// Since the routing table addition failed,
	//	// the peer can be removed from the libp2p peer store
	//	// in the base networking server
	//	//d.logger.Debug("service-->RemoveFromPeerStore", "node", node.String())
	//	d.RemoveFromPeerStore(node)
	//
	//	return err
	//}

	return nil
}

// ConnectToBootnodes attempts to connect to the bootnodes
// and add them to the peer / routing table
func (d *RelayClient) ConnectToBootnodes(bootnodes []*peer.AddrInfo) {
	for _, nodeInfo := range bootnodes {
		if err := d.addToTable(nodeInfo); err != nil {
			d.logger.Error(
				"Failed to add new peer to routing table",
				"peer",
				nodeInfo.ID,
				"err",
				err,
			)
		}
	}
}

// NewProtoConnection opens up a new stream on the set protocol to the peer,
// and returns a reference to the connection
func (s *RelayClient) NewProtoConnection(protocol string, peerID peer.ID) (*rawGrpc.ClientConn, error) {
	s.protocolsLock.Lock()
	defer s.protocolsLock.Unlock()

	p, ok := s.protocols[protocol]
	if !ok {
		return nil, fmt.Errorf("protocol not found: %s", protocol)
	}

	stream, err := s.NewStream(protocol, peerID)
	if err != nil {
		return nil, err
	}

	return p.Client(stream), nil
}

func (s *RelayClient) SaveProtocolStream(
	protocol string,
	stream *rawGrpc.ClientConn,
	peerID peer.ID,
) {
	s.peersLock.Lock()
	defer s.peersLock.Unlock()

	connectionInfo, ok := s.peers[peerID]
	if !ok {
		s.logger.Warn(
			fmt.Sprintf(
				"Attempted to save protocol %s stream for non-existing peer %s",
				protocol,
				peerID,
			),
		)

		return
	}

	connectionInfo.addProtocolStream(protocol, stream)
}

func (s *RelayClient) registerAliveProtocol() {
	grpcStream := grpc.NewGrpcStream()
	s.RegisterProtocol(EdgeAliveProto, grpcStream)
}

func (s *RelayClient) NewStream(proto string, id peer.ID) (network.Stream, error) {
	return s.host.NewStream(context.Background(), id, protocol.ID(proto))
}

func (s *RelayClient) RegisterProtocol(id string, p Protocol) {
	s.protocolsLock.Lock()
	defer s.protocolsLock.Unlock()

	s.protocols[id] = p
}

// NewAliveClient returns a new or existing alive service client connection
func (s *RelayClient) NewAliveClient(peerID peer.ID) (proto.AliveClient, error) {
	conn, err := s.NewProtoConnection(EdgeAliveProto, peerID)
	if err != nil {
		return nil, fmt.Errorf("failed to open a stream, err %w", err)
	}

	s.SaveProtocolStream(EdgeAliveProto, conn, peerID)

	return proto.NewAliveClient(conn), nil
}

func (s *RelayClient) CloseProtocolStream(protocol string, peerID peer.ID) error {
	s.peersLock.Lock()
	defer s.peersLock.Unlock()

	connectionInfo, ok := s.peers[peerID]
	if !ok {
		return nil
	}

	return connectionInfo.removeProtocolStream(protocol)
}

// sayHello call Hello to bootnode
func (s *RelayClient) sayHello(
	peerID peer.ID,
) (bool, string, error) {

	if s.application == nil {
		return false, "", errors.New("no application info")
	}

	clt, clientErr := s.NewAliveClient(peerID)
	if clientErr != nil {
		return false, "", fmt.Errorf("unable to create new alive client connection, %w", clientErr)
	}
	s.logger.Info("Keep alive", "to", peerID.String())
	//get latest app status
	relay := ""
	relayPeerInfo := s.GetRelayPeerInfo()
	if relayPeerInfo != nil && relayPeerInfo.Info != nil {
		relay = fmt.Sprintf("%s/p2p/%s", relayPeerInfo.Info.Info.Addrs[0].String(), relayPeerInfo.Info.Info.ID.String())
	}

	resp, err := clt.Hello(
		context.Background(),
		&proto.AliveStatus{
			Name:         s.application.Name,
			StartupTime:  s.application.StartupTime,
			Uptime:       s.application.Uptime,
			Relay:        relay,
			AppOrigin:    s.application.AppOrigin,
			Mac:          s.application.Mac,
			CpuInfo:      s.application.CpuInfo,
			GpuInfo:      s.application.GpuInfo,
			MemInfo:      s.application.MemInfo,
			ModelHash:    s.application.ModelHash,
			AveragePower: s.application.AveragePower,
			Version:      s.application.Version,
		},
	)
	if err != nil {
		return false, "", err
	}

	if closeErr := s.CloseProtocolStream(EdgeAliveProto, peerID); closeErr != nil {
		return false, "", closeErr
	}

	return resp.Success, resp.Discovery, nil
}

// startAliveService starts the AliveService loop,
// in which random peers are dialed for their peer sets,
// and random bootnodes are dialed for their peer sets
func (s *RelayClient) startAliveService() {
	go s.keepAliveToBootnodes()
	bootnodeAliveTicker := time.NewTicker(bootnodeAliveInterval)

	defer func() {
		bootnodeAliveTicker.Stop()
	}()

	for {
		select {
		case <-s.closeCh:
			return
		case <-bootnodeAliveTicker.C:
			go s.keepAliveToBootnodes()
		}
	}
}

// keepAliveToBootnodes queries a random (unconnected) bootnode for new peers
// and adds them to the routing table
func (s *RelayClient) keepAliveToBootnodes() {
	s.logger.Info("keepAliveToBootnodes doing...")

	var (
		bootnode *peer.AddrInfo // the reference bootnode
	)

	// Try to find a suitable bootnode to use as a reference peer
	for bootnode == nil {
		// Get a random unconnected bootnode from the bootnode set
		bootnode = s.GetRandomBootnode()
		if bootnode == nil {
			return
		}

		for s.application == nil {
			time.Sleep(1 * time.Second)
		}
		success, discovery, err := s.sayHello(bootnode.ID)
		if err != nil {
			s.logger.Debug("Unable to execute bootnode peer alive call",
				"bootnode", bootnode.ID.String(),
				"err", err.Error(),
			)
			bootnode = nil
		}
		s.logger.Debug("keepAliveToBootnodes result", "success", success, "discovery", discovery)

		// add a new found relay node
		if discovery != "" {
			s.addRelaynodes([]string{discovery})
		}
	}
}

func (s *RelayClient) disconnectFromPeer(peer peer.ID, reason string) {
	if s.host.Network().Connectedness(peer) == network.Connected {
		s.logger.Info(fmt.Sprintf("Closing connection to peer [%s] for reason [%s]", peer.String(), reason))

		if closeErr := s.host.Network().ClosePeer(peer); closeErr != nil {
			s.logger.Error(fmt.Sprintf("Unable to gracefully close peer connection, %v", closeErr))
		}
	}
}

// startApplicationEventProcess starts application event subscription
func (m *RelayClient) startApplicationEventProcess(subscrption application.Subscription) {
	m.subscription = subscrption
	for {
		var event *application.Event

		select {
		case <-m.closeCh:
			return
		case event = <-m.subscription.GetEventCh():
		}

		if l := len(event.NewApp); l > 0 {
			latest := event.NewApp[l-1]

			m.application = latest

		}
	}
}

// NewRelayClient returns a new instance of the relay client
func NewRelayClient(logger hclog.Logger, config *emcNetwork.Config, relayOn bool) (*RelayClient, error) {
	logger = logger.Named("relay-client")
	key, err := setupLibp2pKey(config.SecretsManager)
	if err != nil {
		return nil, err
	}

	listenAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", config.Addr.IP.String(), config.Addr.Port))
	if err != nil {
		return nil, err
	}

	if config.NatAddr != nil {
		addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", config.NatAddr.String(), config.Addr.Port))

		if addr != nil {
			listenAddr = addr
		}
	}

	var edgeNodeHost host.Host
	if relayOn {
		edgeNodeHost, err = libp2p.New(
			libp2p.Security(noise.ID, noise.New),
			libp2p.ListenAddrs(listenAddr),
			libp2p.EnableRelay(),
			libp2p.Identity(key),
			libp2p.ForceReachabilityPrivate(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create libp2p privateNodeHost: %w", err)
		}
	} else {
		edgeNodeHost, err = libp2p.New(
			libp2p.Security(noise.ID, noise.New),
			libp2p.ListenAddrs(listenAddr),
			libp2p.Identity(key),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create libp2p edgeNodeHost: %w", err)
		}
	}

	clt := &RelayClient{
		logger:     logger,
		host:       edgeNodeHost,
		closeCh:    make(chan struct{}),
		protocols:  map[string]Protocol{},
		relayPeers: make(map[peer.ID]*RelayPeerInfo),
		peers:      make(map[peer.ID]*PeerConnInfo),
		relaynodes: &relaynodesWrapper{
			relaynodeArr:       make([]*peer.AddrInfo, 0),
			relaynodesMap:      make(map[peer.ID]*peer.AddrInfo),
			relaynodeConnCount: 0,
		},
	}

	clt.logger.Info("LibP2P Relay client running", "addr", edgeNodeHost.Addrs()[0].String()+"/p2p/"+edgeNodeHost.ID().String())

	relaynodes := config.Chain.Relaynodes
	if setupErr := clt.setupRelaynodes(relaynodes); setupErr != nil {
		return nil, fmt.Errorf("unable to parse relaynode data, %w", setupErr)
	}

	clt.registerAliveProtocol()
	return clt, nil
}
