package application

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/emc-protocol/edge-matrix/application/proto"
	"github.com/emc-protocol/edge-matrix/helper/rpc"
	"github.com/emc-protocol/edge-matrix/miner"
	"github.com/emc-protocol/edge-matrix/network"
	"github.com/emc-protocol/edge-matrix/network/event"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

const (
	SyncAppPeerClientLoggerName = "sync-app-peer-client"
	statusTopicName             = "appsyncer/status/0.1"
	defaultTimeoutForStatus     = 10 * time.Second
)

type syncAppPeerClient struct {
	logger     hclog.Logger // logger used for console logging
	network    Network      // reference to the network module
	host       host.Host
	minerAgent *miner.MinerAgent

	subscription           Subscription          // reference to the application subscription
	topic                  *network.Topic        // reference to the network topic
	id                     string                // node id
	peerStatusUpdateCh     chan *AppPeer         // peer status update channel
	peerConnectionUpdateCh chan *event.PeerEvent // peer connection update channel

	jsonRpcClient *rpc.JsonRpcClient
	privateKey    *ecdsa.PrivateKey
	endpoint      *Endpoint

	shouldEmitData bool // flag for emitting data in the topic
	closeCh        chan struct{}
	closed         *uint64 // ACTIVE == 0, CLOSED == non-zero.
}

// Start processes for SyncAppPeerClient
func (m *syncAppPeerClient) Start(subscription Subscription, topicSubFlag bool) error {
	// Mark client active.
	atomic.StoreUint64(m.closed, 0)

	// TODO remove
	//if err := m.startGossip(topicSubFlag); err != nil {
	//	return err
	//}
	//m.logger.Info("startGossip", "topicSubFlag", topicSubFlag)

	go m.startApplicationEventProcess(subscription)
	go m.startPeerEventProcess()

	return nil
}

// Close terminates running processes for SyncAppPeerClient
func (m *syncAppPeerClient) Close() {
	if atomic.SwapUint64(m.closed, 1) > 0 {
		// Already closed.
		return
	}

	if m.topic != nil {
		m.topic.Close()
	}

	if m.subscription != nil {
		m.subscription.Close()

		m.subscription = nil
	}

	if m.closeCh != nil {
		close(m.closeCh)
	}

	close(m.peerStatusUpdateCh)
}

// DisablePublishingPeerStatus disables publishing own status via gossip
func (m *syncAppPeerClient) DisablePublishingPeerStatus() {
	m.shouldEmitData = false
}

// EnablePublishingPeerStatus enables publishing own status via gossip
func (m *syncAppPeerClient) EnablePublishingPeerStatus() {
	m.shouldEmitData = true
}

// GetPeerStatus fetches peer status
func (m *syncAppPeerClient) GetPeerStatus(peerID peer.ID) (*AppPeer, error) {
	clt, err := m.newSyncPeerClient(peerID)
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), defaultTimeoutForStatus)
	defer cancel()

	status, err := clt.GetStatus(timeoutCtx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return &AppPeer{
		ID:           peerID,
		Starup_time:  status.StartupTime,
		Uptime:       status.Uptime,
		Guage_height: status.GuageHeight,
		Guage_max:    status.GuageMax,
		Distance:     m.network.GetPeerDistance(peerID),
	}, nil
}

// GetConnectedPeerStatuses fetches the statuses of all connecting peers
func (m *syncAppPeerClient) GetConnectedPeerStatuses() []*AppPeer {
	var (
		ps            = m.network.Peers()
		syncPeers     = make([]*AppPeer, 0, len(ps))
		syncPeersLock sync.Mutex
		wg            sync.WaitGroup
	)

	for _, p := range ps {
		p := p

		wg.Add(1)

		go func() {
			defer wg.Done()

			peerID := p.Info.ID

			status, err := m.GetPeerStatus(peerID)
			if err != nil {
				m.logger.Warn("failed to get status from a peer, skip", "id", peerID, "err", err)
			}

			syncPeersLock.Lock()

			syncPeers = append(syncPeers, status)

			syncPeersLock.Unlock()
		}()
	}

	wg.Wait()

	return syncPeers
}

// GetPeerStatusUpdateCh returns a channel of peer's status update
func (m *syncAppPeerClient) GetPeerStatusUpdateCh() <-chan *AppPeer {
	return m.peerStatusUpdateCh
}

// GetPeerConnectionUpdateEventCh returns peer's connection change event
func (m *syncAppPeerClient) GetPeerConnectionUpdateEventCh() <-chan *event.PeerEvent {
	return m.peerConnectionUpdateCh
}

// startGossip creates new topic and starts subscribing
func (m *syncAppPeerClient) startGossip(topicSubFlag bool) error {
	topic, err := m.network.NewTopic(statusTopicName, &proto.AppStatus{})
	if err != nil {
		return err
	}

	m.topic = topic

	if topicSubFlag {
		if err := topic.Subscribe(m.handleStatusUpdate); err != nil {
			return fmt.Errorf("unable to subscribe to gossip topic, %w", err)
		}
		m.logger.Info("subscribe to gossip topic=AppStatus")
	}

	return nil
}

// handleStatusUpdate is a handler of gossip
func (m *syncAppPeerClient) handleStatusUpdate(obj interface{}, from peer.ID) {
	status, ok := obj.(*proto.AppStatus)
	if !ok {
		m.logger.Error("failed to cast gossiped message to status")

		return
	}

	if !m.network.IsConnected(from) {
		if m.id != from.String() {
			m.logger.Debug("received status from non-connected peer, ignore", "id", from)
		}

		return
	}

	if atomic.LoadUint64(m.closed) > 0 {
		m.logger.Debug("received status from peer after client was closed, ignoring", "id", from)

		return
	}

	m.peerStatusUpdateCh <- &AppPeer{
		ID:           from,
		Name:         status.Name,
		Starup_time:  status.StartupTime,
		Uptime:       status.Uptime,
		Guage_height: status.GuageHeight,
		Guage_max:    status.GuageMax,
		Distance:     m.network.GetPeerDistance(from),
	}
}

// startNewBlockProcess starts application event subscription
func (m *syncAppPeerClient) startApplicationEventProcess(subscrption Subscription) {
	m.subscription = subscrption
	for {
		var event *Event

		select {
		case <-m.closeCh:
			return
		case event = <-m.subscription.GetEventCh():
		}

		if !m.shouldEmitData {
			continue
		}

		if l := len(event.NewApp); l > 0 {
			latest := event.NewApp[l-1]
			m.logger.Debug("event", "latest", latest)

		}
	}
}

// startPeerEventProcess starts subscribing peer connection change events and process them
func (m *syncAppPeerClient) startPeerEventProcess() {
	defer close(m.peerConnectionUpdateCh)

	peerEventCh, err := m.network.SubscribeCh(context.Background())
	if err != nil {
		m.logger.Error("failed to subscribe", "err", err)

		return
	}

	for {
		select {
		case <-m.closeCh:
			return

		case e := <-peerEventCh:
			if e != nil && (e.Type == event.PeerConnected || e.Type == event.PeerDisconnected) {
				m.peerConnectionUpdateCh <- e
			}
		}
	}
}

// CloseStream closes stream
func (m *syncAppPeerClient) CloseStream(peerID peer.ID) error {
	return m.network.CloseProtocolStream(appSyncerProto, peerID)
}

// GetPeerData returns bytes of data from given hash to peer
func (m *syncAppPeerClient) PostPeerStatusData(peerID peer.ID, nodeId string) (string, error) {
	toPeerId := peerID.String()
	clt, err := m.newSyncPeerClient(peerID)
	if err != nil {
		return "", fmt.Errorf("failed to create sync peer client to %s: %w", toPeerId, err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// TODO handle timeout
	data, err := clt.PostAppStatus(ctx, &proto.PostPeerStatusRequest{
		NodeId: nodeId,
	})
	if err != nil {
		cancel()

		return "", fmt.Errorf("failed to PostPeerStatusData: %w", err)
	}
	if err != nil {
		return "", err
	}
	recv, err := data.Recv()
	if err != nil {
		m.logger.Warn(err.Error())
		return "", err
	}
	m.logger.Info("PostAppStatus result:", recv.Data)
	return recv.Data, nil
}

// GetPeerData returns bytes of data from given hash to peer
func (m *syncAppPeerClient) GetPeerData(peerID peer.ID, hash string, timeout time.Duration) (map[string][]byte, error) {
	clt, err := m.newSyncPeerClient(peerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync peer client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// TODO handle timeout
	data, err := clt.GetData(ctx, &proto.GetDataRequest{
		DataHash: hash,
	})
	if err != nil {
		cancel()

		return nil, fmt.Errorf("failed to open GetData stream: %w", err)
	}
	recv, err := data.Recv()
	if err != nil {
		return nil, err
	}
	return recv.Data, nil
}

// newSyncPeerClient creates gRPC client
func (m *syncAppPeerClient) newSyncPeerClient(peerID peer.ID) (proto.SyncAppClient, error) {
	conn, err := m.network.NewProtoConnection(appSyncerProto, peerID)
	if err != nil {
		return nil, fmt.Errorf("failed to open a stream, err %w", err)
	}

	m.network.SaveProtocolStream(appSyncerProto, conn, peerID)

	return proto.NewSyncAppClient(conn), nil
}

func fromProto(protoData *proto.Data) (map[string][]byte, error) {
	return protoData.Data, nil
}

func dataStreamToChannel(stream proto.SyncApp_GetDataClient) (chan map[string][]byte, <-chan error) {
	dataCh := make(chan map[string][]byte)
	errorCh := make(chan error, 1)

	go func() {
		defer close(dataCh)

		for {
			protoData, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}

			if err != nil {
				errorCh <- err

				break
			}

			data, err := fromProto(protoData)
			if err != nil {
				errorCh <- err

				break
			}

			dataCh <- data
		}
	}()

	return dataCh, errorCh
}

type SyncAppPeerClient interface {
	// Start processes for SyncAppPeerClient
	Start(subscription Subscription, topicSubFlag bool) error
	// Close terminates running processes for SyncAppPeerClient
	Close()
	// GetPeerStatus fetches peer status
	GetPeerStatus(id peer.ID) (*AppPeer, error)
	// GetPeerData fetches peer data
	GetPeerData(peerID peer.ID, dataHash string, timeout time.Duration) (map[string][]byte, error)
	// GetConnectedPeerStatuses fetches the statuses of all connecting peers
	GetConnectedPeerStatuses() []*AppPeer
	// GetPeerStatusUpdateCh returns a channel of peer's status update
	GetPeerStatusUpdateCh() <-chan *AppPeer
	// GetPeerConnectionUpdateEventCh returns peer's connection change event
	GetPeerConnectionUpdateEventCh() <-chan *event.PeerEvent
	// CloseStream close a stream
	CloseStream(peerID peer.ID) error
	// DisablePublishingPeerStatus disables publishing status in syncer topic
	DisablePublishingPeerStatus()
	// EnablePublishingPeerStatus enables publishing status in syncer topic
	EnablePublishingPeerStatus()
}

func NewSyncAppPeerClient(
	logger hclog.Logger,
	network Network,
	minerAgent *miner.MinerAgent,
	host host.Host,
	jsonRpcClient *rpc.JsonRpcClient,
	privateKey *ecdsa.PrivateKey,
	endpoint *Endpoint,
) SyncAppPeerClient {
	return &syncAppPeerClient{
		logger:                 logger.Named(SyncAppPeerClientLoggerName),
		network:                network,
		id:                     network.AddrInfo().ID.String(),
		peerStatusUpdateCh:     make(chan *AppPeer, 1),
		peerConnectionUpdateCh: make(chan *event.PeerEvent, 1),
		shouldEmitData:         true,
		closeCh:                make(chan struct{}),
		closed:                 new(uint64),
		minerAgent:             minerAgent,
		host:                   host,
		jsonRpcClient:          jsonRpcClient,
		privateKey:             privateKey,
		endpoint:               endpoint,
	}
}
