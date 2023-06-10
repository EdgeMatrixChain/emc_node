package application

import (
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/emc-protocol/edge-matrix/validators"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"time"
)

const (
	appSyncerProto = "/appsyncer/0.1"
	syncerName     = "appsyncer"
)

type blockchainStore interface {
	// Header returns the current msg of the chain (genesis if empty)
	Header() *types.Header

	// GetHeaderByNumber gets a msg using the provided number
	GetHeaderByNumber(uint64) (*types.Header, bool)

	// GetBlockByHash gets a block using the provided hash
	GetBlockByHash(hash types.Hash, full bool) (*types.Block, bool)

	// GetBlockByNumber returns a block using the provided number
	GetBlockByNumber(num uint64, full bool) (*types.Block, bool)

	// ReadTxLookup returns a block hash in which a given txn was mined
	ReadTxLookup(txnHash types.Hash) (types.Hash, bool)

	// GetReceiptsByHash returns the receipts for a block hash
	GetReceiptsByHash(hash types.Hash) ([]*types.Receipt, error)
}

type syncer struct {
	logger hclog.Logger

	peerMap            *PeerMap
	syncAppPeerClient  SyncAppPeerClient
	syncAppPeerService SyncAppPeerService

	// Timeout for syncing a block
	blockTimeout time.Duration

	// Channel to notify Sync that a new status arrived
	newStatusCh chan struct{}

	blockchainStore blockchainStore
	host            host.Host

	peersBlockNumMap map[peer.ID]uint64
}

type ValidatorStore interface {
	// Get current validators
	GetCurrentValidators() validators.Validators

	// Get singer address
	GetSignerAddress() types.Address
}

type Syncer interface {
	// Start starts syncer processes
	Start(s Subscription, topicSubFlag bool) error
	// Close terminates syncer process
	Close() error
}

func NewSyncer(
	logger hclog.Logger,
	syncAppPeerClient SyncAppPeerClient,
	syncAppPeerService SyncAppPeerService,
	host host.Host,
	blockchainStore blockchainStore,
) Syncer {
	return &syncer{
		logger:             logger.Named(syncerName),
		syncAppPeerClient:  syncAppPeerClient,
		syncAppPeerService: syncAppPeerService,
		newStatusCh:        make(chan struct{}),
		peerMap:            new(PeerMap),
		host:               host,
		blockchainStore:    blockchainStore,
		peersBlockNumMap:   make(map[peer.ID]uint64),
	}
}

// initializePeerMap fetches peer statuses and initializes map
func (s *syncer) initializePeerMap() {
	peerStatuses := s.syncAppPeerClient.GetConnectedPeerStatuses()
	s.peerMap.Put(peerStatuses...)
}

// Close terminates goroutine processes
func (s *syncer) Close() error {
	close(s.newStatusCh)

	if err := s.syncAppPeerService.Close(); err != nil {
		return err
	}

	s.syncAppPeerClient.Close()

	return nil
}

func (s *syncer) Start(sub Subscription, topicSubFlag bool) error {
	if err := s.syncAppPeerClient.Start(sub, topicSubFlag); err != nil {
		return err
	}

	s.syncAppPeerService.Start()

	go s.startPeerStatusUpdateProcess()
	//go s.startPeerConnectionEventProcess()

	return nil

}

// startPeerStatusUpdateProcess subscribes peer status change event and updates peer map
func (s *syncer) startPeerStatusUpdateProcess() {
	for peerStatus := range s.syncAppPeerClient.GetPeerStatusUpdateCh() {
		s.logger.Info("AppPeerStatus updated ", "NodeID", peerStatus.ID.String())

		defer func() {
			err := s.syncAppPeerClient.CloseStream(peerStatus.ID)
			if err != nil {
				s.logger.Error("Failed to close stream: ", err)
			}
		}()
		// to get a proof result as a validator
		//validators := s.validatorStore.GetCurrentValidators()
		//if validators.Includes(s.address) {
		//
		//}
		s.putToPeerMap(peerStatus)
	}
}

// putToPeerMap puts given status to peer map
func (s *syncer) putToPeerMap(status *AppPeer) {
	s.peerMap.Put(status)
	s.notifyNewStatusEvent()
}

// removeFromPeerMap removes the peer from peer map
func (s *syncer) removeFromPeerMap(peerID peer.ID) {
	s.peerMap.Remove(peerID)
}

// notifyNewStatusEvent emits signal to newStatusCh
func (s *syncer) notifyNewStatusEvent() {
	select {
	case s.newStatusCh <- struct{}{}:
	default:
	}
}
