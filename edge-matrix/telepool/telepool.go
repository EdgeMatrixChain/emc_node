package telepool

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/armon/go-metrics"
	"github.com/emc-protocol/edge-matrix/application"
	"github.com/emc-protocol/edge-matrix/blockchain"
	"github.com/emc-protocol/edge-matrix/contracts"
	"github.com/emc-protocol/edge-matrix/network"
	"github.com/emc-protocol/edge-matrix/telepool/proto"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/umbracle/fastrlp"
	"io"
	"math/big"
	"sync/atomic"
	"time"
)

// indicates origin of a transaction
type teleOrigin int

const (
	local  teleOrigin = iota // json-RPC/gRPC endpoints
	gossip                   // gossip protocol
)

// errors
var (
	ErrIntrinsicGas            = errors.New("intrinsic gas too low")
	ErrBlockLimitExceeded      = errors.New("exceeds block gas limit")
	ErrNegativeValue           = errors.New("negative value")
	ErrExtractSignature        = errors.New("cannot extract signature")
	ErrInvalidSender           = errors.New("invalid sender")
	ErrInvalidProvider         = errors.New("invalid provider")
	ErrTxPoolOverflow          = errors.New("txpool is full")
	ErrUnderpriced             = errors.New("transaction underpriced")
	ErrNonceTooLow             = errors.New("nonce too low")
	ErrNonceTooHigh            = errors.New("nonce too high")
	ErrInsufficientFunds       = errors.New("insufficient funds for gas * price + value")
	ErrInvalidAccountState     = errors.New("invalid account state")
	ErrAlreadyKnown            = errors.New("already known")
	ErrOversizedData           = errors.New("oversized data")
	ErrMaxEnqueuedLimitReached = errors.New("maximum number of enqueued transactions reached")
	ErrRejectFutureTx          = errors.New("rejected future tx due to low slots")
	ErrSmartContractRestricted = errors.New("smart contract deployment restricted")
)

func (o teleOrigin) String() (s string) {
	switch o {
	case local:
		s = "local"
	case gossip:
		s = "gossip"
	}

	return
}

const (
	txSlotSize  = 32 * 1024
	txMaxSize   = 512 * 1024 // 512k
	topicNameV1 = "tele/0.3"

	// maximum allowed number of times an account
	// was excluded from block building (ibft.writeTransactions)
	maxAccountDemotions uint64 = 10

	// maximum allowed number of consecutive blocks that don't have the account's transaction
	maxAccountSkips = uint64(10)
	pruningCooldown = 2000 * time.Millisecond

	// txPoolMetrics is a prefix used for txpool-related metrics
	txPoolMetrics = "telepool"
)

var marshalArenaPool fastrlp.ArenaPool

type enqueueRequest struct {
	tele *types.Telegram
}

// store interface defines State helper methods the TxPool should have access to
type store interface {
	Header() *types.Header
	GetNonce(root types.Hash, addr types.Address) uint64
	GetBalance(root types.Hash, addr types.Address) (*big.Int, error)
	GetBlockByHash(types.Hash, bool) (*types.Block, bool)
}

type signer interface {
	Sender(tele *types.Telegram) (types.Address, error)
	Provider(tele *types.Telegram) (types.Address, error)
}

type providerSigner interface {
	Provider(tele *types.Telegram) (types.Address, error)
}

// A promoteRequest is created each time some account
// is eligible for promotion. This request is signaled
// on 2 occasions:
//
// 1. When an enqueued transaction's nonce is
// not greater than the expected (account's nextNonce).
// == nextNonce - transaction is expected (addTele)
// < nextNonce - transaction was demoted (Demote)
//
// 2. When an account's nextNonce is updated (during ResetWithHeader)
// and the first enqueued transaction matches the new nonce.
type promoteRequest struct {
	account types.Address
}

type Config struct {
	MaxSlots           uint64
	MaxAccountEnqueued uint64
}

type TelegramPool struct {
	logger         hclog.Logger
	signer         signer
	providerSigner providerSigner
	store          store
	// map of all accounts registered by the pool
	accounts accountsMap

	// all the primaries sorted by max gas price
	executables *pricedQueue

	// lookup map keeping track of all
	// transactions present in the pool
	index lookupMap

	// networking stack
	topic       *network.Topic
	network     *network.Server
	edgeNetwork *network.Server

	appSyncer application.Syncer

	// gauge for measuring pool capacity
	gauge slotGauge

	// channels on which the pool's event loop
	// does dispatching/handling requests.
	enqueueReqCh chan enqueueRequest
	promoteReqCh chan promoteRequest
	pruneCh      chan struct{}

	// shutdown channel
	shutdownCh chan struct{}

	// flag indicating if the current node is a sealer,
	// and should therefore gossip transactions
	sealing uint32

	// Event manager for telepool events
	eventManager *eventManager

	// indicates which txpool operator commands should be implemented
	proto.UnimplementedTxnPoolOperatorServer

	// pending is the list of pending and ready transactions. This variable
	// is accessed with atomics
	pending int64
}

// NewTelegramPool returns a new pool for processing incoming telegram.
func NewTelegramPool(
	logger hclog.Logger,
	store store,
	network *network.Server,
	edgeNetwork *network.Server,
	config *Config,
	teleVesion string,
) (*TelegramPool, error) {
	pool := &TelegramPool{
		logger:      logger.Named("telepool"),
		store:       store,
		executables: newPricedQueue(),
		accounts:    accountsMap{maxEnqueuedLimit: config.MaxAccountEnqueued},
		index:       lookupMap{all: make(map[types.Hash]*types.Telegram)},
		gauge:       slotGauge{height: 0, max: config.MaxSlots},
		//	main loop channels
		enqueueReqCh: make(chan enqueueRequest),
		promoteReqCh: make(chan promoteRequest),
		pruneCh:      make(chan struct{}),
		shutdownCh:   make(chan struct{}),
		network:      network,
		edgeNetwork:  edgeNetwork,
	}

	// Attach the event manager
	pool.eventManager = newEventManager(pool.logger)

	if network != nil {
		// subscribe to the gossip protocol
		protoId := topicNameV1
		if teleVesion != "" {
			protoId = "tele/" + teleVesion
		}
		topic, err := network.NewTopic(protoId, &proto.Txn{})
		if err != nil {
			return nil, err
		}

		if subscribeErr := topic.Subscribe(pool.addGossipTele); subscribeErr != nil {
			return nil, fmt.Errorf("unable to subscribe to gossip topic, %w", subscribeErr)
		}

		pool.topic = topic
	}

	return pool, nil
}

// sealing returns the current set sealing flag
func (p *TelegramPool) getSealing() bool {
	return atomic.LoadUint32(&p.sealing) == 1
}

// addGossipTele handles receiving telegram
// gossiped by the network.
func (p *TelegramPool) addGossipTele(obj interface{}, _ peer.ID) {
	if !p.getSealing() {
		return
	}

	raw, ok := obj.(*proto.Txn)
	if !ok {
		p.logger.Error("failed to cast gossiped message to telegram")

		return
	}

	// Verify that the gossiped telegram message is not empty
	if raw == nil || raw.Raw == nil {
		p.logger.Error("malformed gossip telegram message received")

		return
	}

	tele := new(types.Telegram)

	// decode telegram
	if err := tele.UnmarshalRLP(raw.Raw.Value); err != nil {
		p.logger.Error("failed to decode broadcast telegram", "err", err)

		return
	}

	// add telegram
	if _, err := p.addTele(gossip, tele); err != nil {
		if errors.Is(err, ErrAlreadyKnown) {
			p.logger.Debug("rejecting known telegram (gossip)", "hash", tele.Hash.String())

			return
		}

		p.logger.Error("failed to add broadcast telegram", "err", err, "hash", tele.Hash.String())
	}
}

func (p *TelegramPool) SetAppSyncer(appSyncer application.Syncer) {
	p.appSyncer = appSyncer
}

// AddTele adds a new telegram to the pool (sent from json-RPC/gRPC endpoints)
// and broadcasts it to the network (if enabled).
func (p *TelegramPool) AddTele(tele *types.Telegram) (string, error) {
	resp := &application.EdgeResponse{}
	if tele.To != nil && *tele.To == contracts.EdgeCallPrecompile {
		input := tele.Input
		call := &application.EdgeCall{}
		if err := json.Unmarshal(input, &call); err != nil {
			return "", err
		}
		if call.Endpoint == "/poc_cpu_request" || call.Endpoint == "/poc_cpu_validate" {
			host := p.edgeNetwork.GetHost()

			relayAddr, addr := p.getAppPeerAddr(call.PeerId)
			p.logger.Info("edge call", "PeerId", call.PeerId, "Endpoint", call.Endpoint, addr, "Relay", relayAddr)
			if relayAddr != "" || addr != "" {
				clientHost, err2 := p.newTempHost()
				if err2 != nil {
					return "", err2
				}
				defer clientHost.Close()

				host = clientHost
				err := p.addAddrToHost(call.PeerId, host, addr, relayAddr)
				if err != nil {
					return "", err
				}
			}

			respBuf, callErr := application.Call(host, application.ProtoTagEcApp, call)
			if callErr != nil {
				return "", callErr
			}

			err := resp.UnmarshalRLP(respBuf)
			if err != nil {
				return "", err
			}
			tele.RespFrom = resp.From
			tele.RespR = resp.R
			tele.RespV = resp.V
			tele.RespS = resp.S
			tele.RespHash = resp.Hash
			if len(resp.RespString) > 0 {
				return resp.RespString, nil
			} else {
				return "", nil
			}
		}
	}

	if tele.RespV == nil {
		tele.RespV = big.NewInt(0)
	}
	if tele.RespR == nil {
		tele.RespR = big.NewInt(0)
	}
	if tele.RespS == nil {
		tele.RespS = big.NewInt(0)
	}
	tele.RespHash = types.ZeroHash
	tele.RespFrom = types.ZeroAddress

	respString, err := p.addTele(local, tele)
	if err != nil {
		p.logger.Error("failed to add telegram", "err", err)

		return "", err
	}

	// broadcast the transaction only if a topic
	// subscription is present
	if p.topic != nil {
		tx := &proto.Txn{
			Raw: &any.Any{
				Value: tele.MarshalRLP(),
			},
		}

		if err := p.topic.Publish(tx); err != nil {
			p.logger.Error("failed to topic tx", "err", err)
		}
	}

	return respString, nil
}

func (p *TelegramPool) addAddrToHost(peerId string, host host.Host, addr string, relayAddr string) error {
	if relayAddr != "" {
		targetRelayInfo, err := peer.AddrInfoFromString(fmt.Sprintf("%s/p2p-circuit/p2p/%s", relayAddr, peerId))
		if err != nil {
			return err
		}
		host.Peerstore().AddAddrs(targetRelayInfo.ID, targetRelayInfo.Addrs, peerstore.AddressTTL)
	} else if addr != "" {
		addrInfo, err := peer.AddrInfoFromString(fmt.Sprintf("%s/p2p/%s", addr, peerId))
		if err != nil {
			return err
		}
		host.Peerstore().AddAddrs(addrInfo.ID, addrInfo.Addrs, peerstore.RecentlyConnectedAddrTTL)
	}
	return nil
}

func (p *TelegramPool) newTempHost() (host.Host, error) {
	var r io.Reader
	r = rand.Reader
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}
	listen, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/10001")
	clientHost, err := libp2p.New(
		libp2p.ListenAddrs(listen),
		libp2p.Security(noise.ID, noise.New),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		return nil, err
	}
	return clientHost, nil
}

func (p *TelegramPool) getAppPeerAddr(peerId string) (relayAddr string, addr string) {
	if p.appSyncer != nil {
		appPeer := p.appSyncer.GetAppPeer(peerId)
		if appPeer != nil {
			relayAddr = appPeer.Relay
			addr = appPeer.Addr
			return
		}
	}
	return "", ""
}

// addTele is the main entry point to the pool
// for all new transactions. If the call is
// successful, an account is created for this address
// (only once) and an enqueueRequest is signaled.
func (p *TelegramPool) addTele(origin teleOrigin, tele *types.Telegram) (string, error) {
	// validate incoming tele
	if err := p.validateTele(tele); err != nil {
		return "", err
	}

	if p.gauge.highPressure() {
		p.signalPruning()

		//	only accept transactions with expected nonce
		if account := p.accounts.get(tele.From); account != nil &&
			tele.Nonce > account.getNonce() {
			return "", ErrRejectFutureTx
		}
	}

	// check for overflow
	if p.gauge.read()+slotsRequired(tele) > p.gauge.max {
		return "", ErrTxPoolOverflow
	}

	respString := ""
	// telegram for edge call
	if origin == local {
		resp := &application.EdgeResponse{}
		if tele.To != nil && *tele.To == contracts.EdgeCallPrecompile {
			input := tele.Input
			call := &application.EdgeCall{}
			if err := json.Unmarshal(input, &call); err != nil {
				return "", err
			}

			host := p.edgeNetwork.GetHost()

			relayAddr, addr := p.getAppPeerAddr(call.PeerId)
			p.logger.Info("edge call", "PeerId", call.PeerId, "Endpoint", call.Endpoint, addr, "Relay", relayAddr)
			if relayAddr != "" || addr != "" {
				clientHost, err2 := p.newTempHost()
				if err2 != nil {
					return "", err2
				}
				defer clientHost.Close()

				host = clientHost
				err := p.addAddrToHost(call.PeerId, host, addr, relayAddr)
				if err != nil {
					return "", err
				}
			}

			// TODO relpace Call to CallWithFrom
			//respBuf, callErr := application.CallWithFrom(p.edgeNetwork.GetHost(), application.ProtoTagEcApp, call, tele.From)
			respBuf, callErr := application.Call(host, application.ProtoTagEcApp, call)
			if callErr != nil {
				return "", callErr
			}

			err := resp.UnmarshalRLP(respBuf)
			if err != nil {
				return "", err
			}
			tele.RespFrom = resp.From
			tele.RespR = resp.R
			tele.RespV = resp.V
			tele.RespS = resp.S
			tele.RespHash = resp.Hash
			if len(resp.RespString) > 0 {
				respString = resp.RespString
			}
		}
	}
	tele.ComputeHash()

	// add to index
	if ok := p.index.add(tele); !ok {
		return "", ErrAlreadyKnown
	}

	// initialize account for this address once
	p.createAccountOnce(tele.From)

	// send request [BLOCKING]
	p.enqueueReqCh <- enqueueRequest{tele: tele}
	//p.eventManager.signalEvent(proto.EventType_ADDED, tele.Hash)

	return respString, nil
}

// validateTele ensures the telegram conforms to specific
// constraints before entering the pool.
func (p *TelegramPool) validateTele(tele *types.Telegram) error {
	// Check the transaction size to overcome DOS Attacks
	if uint64(len(tele.MarshalRLP())) > txMaxSize {
		return ErrOversizedData
	}

	// Check if the transaction is signed properly

	// Extract the sender
	from, signerErr := p.signer.Sender(tele)
	if signerErr != nil {
		return ErrExtractSignature
	}

	p.logger.Debug(fmt.Sprintf("validateTele from: %s", from.String()))

	// Extract the provider
	if tele.RespFrom != types.ZeroAddress {
		respFrom, signerErr := p.signer.Provider(tele)
		if signerErr != nil {
			return ErrExtractSignature
		}
		p.logger.Debug(fmt.Sprintf("validateTele RespFrom:%s, provider: %s", tele.RespFrom, respFrom.String()))
		if respFrom != tele.RespFrom {
			return ErrInvalidProvider

		}
	}
	// testAddress
	//from := types.StringToAddress("0x68b95f67a32935e3ed85600F558b74E0d2747120")

	// If the from field is set, check that
	// it matches the signer
	if tele.From != types.ZeroAddress &&
		tele.From != from {
		return ErrInvalidSender
	}

	// If no address was set, update it
	if tele.From == types.ZeroAddress {
		tele.From = from
	}

	// Grab the state root for the latest block
	stateRoot := p.store.Header().StateRoot

	// Check nonce ordering
	if p.store.GetNonce(stateRoot, tele.From) > tele.Nonce {
		return ErrNonceTooLow
	}

	// Check max nonce
	teleAcct := p.accounts.get(tele.From)
	if teleAcct != nil && tele.Nonce >= (p.store.GetNonce(stateRoot, tele.From)+teleAcct.maxEnqueued) {
		// don't signal promotion for
		// higher nonce txs
		return ErrNonceTooHigh
	}

	return nil
}

// createAccountOnce creates an account and
// ensures it is only initialized once.
func (p *TelegramPool) createAccountOnce(newAddr types.Address) *account {
	p.logger.Debug(fmt.Sprintf("createAccountOnce: %s", newAddr.String()))

	if p.accounts.exists(newAddr) {
		return nil
	}

	// fetch nonce from state
	stateRoot := p.store.Header().StateRoot
	stateNonce := p.store.GetNonce(stateRoot, newAddr)

	p.logger.Debug(fmt.Sprintf("accounts.initOnce: %s", newAddr.String()))

	// initialize the account
	return p.accounts.initOnce(newAddr, stateNonce)
}

func (p *TelegramPool) signalPruning() {
	select {
	case p.pruneCh <- struct{}{}:
	default: //	pruning handler is active or in cooldown
	}
}

func (p *TelegramPool) updatePending(i int64) {
	newPending := atomic.AddInt64(&p.pending, i)
	metrics.SetGauge([]string{txPoolMetrics, "pending_transactions"}, float32(newPending))
}

// Start runs the pool's main loop in the background.
// On each request received, the appropriate handler
// is invoked in a separate goroutine.
func (p *TelegramPool) Start() {
	// set default value of txpool pending transactions gauge
	p.updatePending(0)

	//	run the handler for high gauge level pruning
	go func() {
		for {
			select {
			case <-p.shutdownCh:
				return
			case <-p.pruneCh:
				p.pruneAccountsWithNonceHoles()
			}

			//	handler is in cooldown to avoid successive calls
			//	which could be just no-ops
			time.Sleep(pruningCooldown)
		}
	}()

	//	run the handler for the tx pipeline
	go func() {
		for {
			select {
			case <-p.shutdownCh:
				return
			case req := <-p.enqueueReqCh:
				go p.handleEnqueueRequest(req)
			case req := <-p.promoteReqCh:
				go p.handlePromoteRequest(req)
			}
		}
	}()
}

// Close shuts down the pool's main loop.
func (p *TelegramPool) Close() {
	p.eventManager.Close()
	p.shutdownCh <- struct{}{}
}

// SetSigner sets the signer the pool will use
// to validate a telegram's signature.
func (p *TelegramPool) SetSigner(s signer) {
	p.signer = s
}

// SetSealing sets the sealing flag
func (p *TelegramPool) SetSealing(sealing bool) {
	newValue := uint32(0)
	if sealing {
		newValue = 1
	}

	atomic.CompareAndSwapUint32(
		&p.sealing,
		p.sealing,
		newValue,
	)
}

func (p *TelegramPool) pruneAccountsWithNonceHoles() {
	p.accounts.Range(
		func(_, value interface{}) bool {
			account, _ := value.(*account)

			account.enqueued.lock(true)
			defer account.enqueued.unlock()

			firstTx := account.enqueued.peek()

			if firstTx == nil {
				return true
			}

			if firstTx.Nonce == account.getNonce() {
				return true
			}

			removed := account.enqueued.clear()

			p.index.remove(removed...)
			p.gauge.decrease(slotsRequired(removed...))

			return true
		},
	)
}

// handleEnqueueRequest attempts to enqueue the transaction
// contained in the given request to the associated account.
// If, afterwards, the account is eligible for promotion,
// a promoteRequest is signaled.
func (p *TelegramPool) handleEnqueueRequest(req enqueueRequest) {
	tele := req.tele
	addr := req.tele.From

	// fetch account
	account := p.accounts.get(addr)

	p.logger.Debug(fmt.Sprintf("handleEnqueueRequest From:%s", addr.String()))

	// enqueue telegram
	if err := account.enqueue(tele); err != nil {
		p.logger.Error("enqueue request", "err", err)

		p.index.remove(tele)

		return
	}

	p.logger.Debug("enqueue request", "hash", tele.Hash.String())

	p.gauge.increase(slotsRequired(tele))

	//p.eventManager.signalEvent(proto.EventType_ENQUEUED, tele.Hash)

	if tele.Nonce > account.getNonce() {
		// don't signal promotion for
		// higher nonce txs
		return
	}

	p.promoteReqCh <- promoteRequest{account: addr} // BLOCKING
}

// handlePromoteRequest handles moving promotable transactions
// of some account from enqueued to promoted. Can only be
// invoked by handleEnqueueRequest or resetAccount.
func (p *TelegramPool) handlePromoteRequest(req promoteRequest) {
	addr := req.account
	account := p.accounts.get(addr)

	// promote enqueued txs
	promoted, pruned := account.promote()
	p.logger.Debug("promote request", "promoted", promoted, "addr", addr.String())

	p.index.remove(pruned...)
	p.gauge.decrease(slotsRequired(pruned...))

	// update metrics
	p.updatePending(int64(len(promoted)))

	//p.eventManager.signalEvent(proto.EventType_PROMOTED, toHash(promoted...)...)
}

// Prepare generates all the transactions
// ready for execution. (primaries)
func (p *TelegramPool) Prepare() {
	// clear from previous round
	if p.executables.length() != 0 {
		p.executables.clear()
	}

	// fetch primary from each account
	primaries := p.accounts.getPrimaries()

	// push primaries to the executables queue
	for _, tx := range primaries {
		p.executables.push(tx)
	}
}

// Peek returns the best-price selected
// transaction ready for execution.
func (p *TelegramPool) Peek() *types.Telegram {
	// Popping the executables queue
	// does not remove the actual tx
	// from the pool.
	// The executables queue just provides
	// insight into which account has the
	// highest priced tx (head of promoted queue)
	return p.executables.pop()
}

// Pop removes the given transaction from the
// associated promoted queue (account).
// Will update executables with the next primary
// from that account (if any).
func (p *TelegramPool) Pop(tx *types.Telegram) {
	// fetch the associated account
	account := p.accounts.get(tx.From)

	account.promoted.lock(true)
	defer account.promoted.unlock()

	// pop the top most promoted tx
	account.promoted.pop()

	// successfully popping an account resets its demotions count to 0
	account.resetDemotions()

	// update state
	p.gauge.decrease(slotsRequired(tx))

	// update metrics
	p.updatePending(-1)

	// update executables
	if tx := account.promoted.peek(); tx != nil {
		p.executables.push(tx)
	}
}

// Drop clears the entire account associated with the given transaction
// and reverts its next (expected) nonce.
func (p *TelegramPool) Drop(tx *types.Telegram) {
	// fetch associated account
	account := p.accounts.get(tx.From)

	account.promoted.lock(true)
	account.enqueued.lock(true)

	// num of all txs dropped
	droppedCount := 0

	// pool resource cleanup
	clearAccountQueue := func(txs []*types.Telegram) {
		p.index.remove(txs...)
		p.gauge.decrease(slotsRequired(txs...))

		// increase counter
		droppedCount += len(txs)
	}

	defer func() {
		account.enqueued.unlock()
		account.promoted.unlock()
	}()

	// rollback nonce
	nextNonce := tx.Nonce
	account.setNonce(nextNonce)

	// drop promoted
	dropped := account.promoted.clear()
	clearAccountQueue(dropped)

	// update metrics
	p.updatePending(-1 * int64(len(dropped)))

	// drop enqueued
	dropped = account.enqueued.clear()
	clearAccountQueue(dropped)

	//p.eventManager.signalEvent(proto.EventType_DROPPED, tx.Hash)
	p.logger.Debug("dropped account txs",
		"num", droppedCount,
		"next_nonce", nextNonce,
		"address", tx.From.String(),
	)
}

// Demote excludes an account from being further processed during block building
// due to a recoverable error. If an account has been demoted too many times (maxAccountDemotions),
// it is Dropped instead.
func (p *TelegramPool) Demote(tx *types.Telegram) {
	account := p.accounts.get(tx.From)
	if account.Demotions() >= maxAccountDemotions {
		p.logger.Debug(
			"Demote: threshold reached - dropping account",
			"addr", tx.From.String(),
		)

		p.Drop(tx)

		// reset the demotions counter
		account.resetDemotions()

		return
	}

	account.incrementDemotions()

	p.eventManager.signalEvent(proto.EventType_DEMOTED, tx.Hash)
}

// ResetWithHeaders processes the transactions from the new
// headers to sync the pool with the new state.
func (p *TelegramPool) ResetWithHeaders(headers ...*types.Header) {
	// process the txs in the event
	// to make sure the pool is up-to-date
	p.processEvent(&blockchain.Event{
		NewChain: headers,
	})
}

// processEvent collects the latest nonces for each account containted
// in the received event. Resets all known accounts with the new nonce.
func (p *TelegramPool) processEvent(event *blockchain.Event) {
	// Grab the latest state root now that the block has been inserted
	stateRoot := p.store.Header().StateRoot
	stateNonces := make(map[types.Address]uint64)

	// discover latest (next) nonces for all accounts
	for _, header := range event.NewChain {
		block, ok := p.store.GetBlockByHash(header.Hash, true)
		if !ok {
			p.logger.Error("could not find block in store", "hash", header.Hash.String())

			continue
		}

		// remove mined txs from the lookup map
		p.index.remove(block.Telegrams...)

		// Extract latest nonces
		for _, tx := range block.Telegrams {
			var err error

			addr := tx.From
			if addr == types.ZeroAddress {
				// From field is not set, extract the signer
				if addr, err = p.signer.Sender(tx); err != nil {
					p.logger.Error(
						fmt.Sprintf("unable to extract signer for transaction, %v", err),
					)

					continue
				}
			}

			// skip already processed accounts
			if _, processed := stateNonces[addr]; processed {
				continue
			}

			// fetch latest nonce from the state
			latestNonce := p.store.GetNonce(stateRoot, addr)

			// update the result map
			stateNonces[addr] = latestNonce
		}
	}

	// reset accounts with the new state
	p.resetAccounts(stateNonces)

	if !p.getSealing() {
		// only non-validator cleanup inactive accounts
		p.updateAccountSkipsCounts(stateNonces)
	}
}

// updateAccountSkipsCounts update the accounts' skips,
// the number of the consecutive blocks that doesn't have the account's transactions
func (p *TelegramPool) updateAccountSkipsCounts(latestActiveAccounts map[types.Address]uint64) {
	p.accounts.Range(
		func(key, value interface{}) bool {
			address, _ := key.(types.Address)
			account, _ := value.(*account)

			if _, ok := latestActiveAccounts[address]; ok {
				account.resetSkips()

				return true
			}

			firstTx := account.getLowestTx()
			if firstTx == nil {
				// no need to increment anything,
				// account has no txs
				return true
			}

			account.incrementSkips()

			if account.skips < maxAccountSkips {
				return true
			}

			// account has been skipped too many times
			p.Drop(firstTx)

			account.resetSkips()

			return true
		},
	)
}

// resetAccounts updates existing accounts with the new nonce and prunes stale transactions.
func (p *TelegramPool) resetAccounts(stateNonces map[types.Address]uint64) {
	if len(stateNonces) == 0 {
		return
	}

	var (
		allPrunedPromoted []*types.Telegram
		allPrunedEnqueued []*types.Telegram
	)

	// clear all accounts of stale txs
	for addr, newNonce := range stateNonces {
		account := p.accounts.get(addr)

		if account == nil {
			// no updates for this account
			continue
		}

		prunedPromoted, prunedEnqueued := account.reset(newNonce, p.promoteReqCh)

		// append pruned
		allPrunedPromoted = append(allPrunedPromoted, prunedPromoted...)
		allPrunedEnqueued = append(allPrunedEnqueued, prunedEnqueued...)

		// new state for account -> demotions are reset to 0
		account.resetDemotions()
	}

	// pool cleanup callback
	cleanup := func(stale []*types.Telegram) {
		p.index.remove(stale...)
		p.gauge.decrease(slotsRequired(stale...))
	}

	// prune pool state
	if len(allPrunedPromoted) > 0 {
		cleanup(allPrunedPromoted)

		p.eventManager.signalEvent(
			proto.EventType_PRUNED_PROMOTED,
			toHash(allPrunedPromoted...)...,
		)

		p.updatePending(int64(-1 * len(allPrunedPromoted)))
	}

	if len(allPrunedEnqueued) > 0 {
		cleanup(allPrunedEnqueued)

		p.eventManager.signalEvent(
			proto.EventType_PRUNED_ENQUEUED,
			toHash(allPrunedEnqueued...)...,
		)
	}
}

// Length returns the total number of all promoted transactions.
func (p *TelegramPool) Length() uint64 {
	return p.accounts.promoted()
}

// toHash returns the hash(es) of given transaction(s)
func toHash(txs ...*types.Telegram) (hashes []types.Hash) {
	for _, tx := range txs {
		hashes = append(hashes, tx.Hash)
	}

	return
}
