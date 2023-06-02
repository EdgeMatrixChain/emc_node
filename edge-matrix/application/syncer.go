package application

import (
	"fmt"
	"github.com/emc-protocol/edge-matrix/application/proof"
	"github.com/emc-protocol/edge-matrix/miner"
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

	validatorStore  ValidatorStore
	blockchainStore blockchainStore
	host            host.Host
	address         types.Address

	// agent for communicating with IC Miner Canister
	minerAgent *miner.MinerAgent

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
	validatorStore ValidatorStore,
	minerAgent *miner.MinerAgent,
) Syncer {
	return &syncer{
		logger:             logger.Named(syncerName),
		syncAppPeerClient:  syncAppPeerClient,
		syncAppPeerService: syncAppPeerService,
		newStatusCh:        make(chan struct{}),
		peerMap:            new(PeerMap),
		host:               host,
		address:            validatorStore.GetSignerAddress(),
		blockchainStore:    blockchainStore,
		validatorStore:     validatorStore,
		minerAgent:         minerAgent,
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
		validators := s.validatorStore.GetCurrentValidators()
		if validators.Includes(s.address) {
			// get latest block number
			header := s.blockchainStore.Header()
			if header != nil {
				blockNumber := header.Number
				// check latest proof number
				latestProofNum, ok := s.peersBlockNumMap[peerStatus.ID]
				if !ok {
					latestProofNum = 0
				}
				var blockNumberFixed uint64 = 0
				if (blockNumber - latestProofNum) > proof.DefaultProofBlockMinDuration {
					// send proof task to peer node
					blockNumberFixed = (blockNumber / proof.DefaultProofBlockRange) * proof.DefaultProofBlockRange
					s.peersBlockNumMap[peerStatus.ID] = blockNumberFixed // commet this line for disable check blocknum
					start := time.Now()

					//  get data from peer
					s.logger.Info(fmt.Sprintf("\n------------------------------------------\nGetPeerData: %s", peerStatus.ID.String()))
					dataMap, err := s.syncAppPeerClient.GetPeerData(peerStatus.ID, header.Hash.String(), time.Second*30)
					if err != nil {
						s.logger.Error(err.Error())
					}
					usedTime := time.Since(start).Milliseconds()

					// validate data
					//if s.logger.IsDebug() {
					//	s.logger.Debug("PeerData: {")
					//	for dataKey, bytes := range dataMap {
					//		s.logger.Debug(dataKey, hex.EncodeToString(bytes))
					//	}
					//	s.logger.Debug("}")
					//}

					var hashArray = make([]string, proof.DefaultHashProofCount)
					target := proof.DefaultHashProofTarget
					loops := proof.DefaultHashProofCount
					i := 0
					initSeed := header.Hash.String()
					for i < loops {
						seed := fmt.Sprintf("%s,%d", initSeed, i)
						hashArray[i] = seed
						i += 1
					}

					validateSuccess := 0
					validateStart := time.Now()
					for _, hash := range hashArray {
						isValidate := proof.ValidateHash(hash, target, dataMap[hash])
						if isValidate {
							validateSuccess += 1
						}
					}

					validateUsedTime := time.Since(validateStart).Milliseconds()
					rate := float32(validateSuccess) / float32(proof.DefaultHashProofCount)
					s.logger.Debug(fmt.Sprintf("used time for validate\t\t: %dms", validateUsedTime))
					s.logger.Info(fmt.Sprintf("validate success\t\t\t: %d/%d rate:%f nodeID:%s", validateSuccess, loops, rate, peerStatus.ID.String()))
					if rate >= 0.95 {
						// valid proof
						s.logger.Info("\n------------------------------------------\nSubmit proof to IC", "usedTime(ms)", usedTime, "blockNumber", blockNumberFixed, "NodeID", peerStatus.ID.String())
						// submit proof result to IC canister
						err := s.minerAgent.SubmitValidation(
							int64(blockNumberFixed),
							s.minerAgent.GetIdentity(),
							usedTime,
							peerStatus.ID.String(),
						)
						if err != nil {
							s.logger.Error("\n------------------------------------------\nSubmitValidation:", "err", err)
						}
					}
				} else {
					s.logger.Warn(fmt.Sprintf("\n------------------------------------------\ninvalid blockNum: %d, NodeId:%s", blockNumberFixed, peerStatus.ID.String()))
				}
			}
		}

		//s.putToPeerMap(peerStatus)
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
