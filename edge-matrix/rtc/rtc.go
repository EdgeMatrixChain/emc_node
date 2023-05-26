package rtc

import (
	"context"
	"errors"
	"fmt"
	"github.com/emc-protocol/edge-matrix/helper/keccak"
	"github.com/emc-protocol/edge-matrix/network"
	"github.com/emc-protocol/edge-matrix/rtc/proto"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/hashicorp/go-hclog"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/umbracle/fastrlp"
	"google.golang.org/protobuf/types/known/anypb"
	"math/big"
	"sync"
)

const (
	topicNameV1 = "rtc/0.1"
	rtcSlotSize = 32 * 1024  // 32kB
	rtcMaxSize  = 128 * 1024 // 128Kb

)

const (
	local  msgOrigin = iota // json-RPC/gRPC endpoints
	gossip                  // gossip protocol
)

const (
	SubjectMsg   RtcType = 0x0
	StateMsg     RtcType = 0x01
	SubscribeMsg RtcType = 0x02
)

// errors
var (
	ErrNegativeValue           = errors.New("negative value")
	ErrExtractSignature        = errors.New("cannot extract signature")
	ErrInvalidSender           = errors.New("invalid sender")
	ErrRtcPoolOverflow         = errors.New("rtc pool is full")
	ErrNonceTooLow             = errors.New("nonce too low")
	ErrInvalidAccountState     = errors.New("invalid account state")
	ErrAlreadyKnown            = errors.New("already known")
	ErrOversizedData           = errors.New("oversized data")
	ErrMaxEnqueuedLimitReached = errors.New("maximum number of enqueued rtc msg reached")
	ErrRejectFutureTx          = errors.New("rejected future rtc msg due to low slots")
)

type msgOrigin int
type RtcType byte

var marshalArenaPool fastrlp.ArenaPool

func (o msgOrigin) String() (s string) {
	switch o {
	case local:
		s = "local"
	case gossip:
		s = "gossip"
	}

	return
}

type RtcMsg struct {
	//Nonce       uint64
	Subject     string
	Application string
	Content     string

	V    *big.Int
	R    *big.Int
	S    *big.Int
	Hash types.Hash
	From types.Address

	To types.Address

	Type RtcType
}

type enqueueRequest struct {
	msg *RtcMsg
}

type promoteRequest struct {
	account *RtcMsg
}
type signer interface {
	Sender(msg *RtcMsg) (types.Address, error)
}

type Rtc struct {
	sync.RWMutex

	logger  hclog.Logger // The logger object
	signer  signer
	ctx     context.Context
	genesis types.Hash   // The hash of the genesis block
	stream  *eventStream // Event subscriptions

	writeLock sync.Mutex

	topics        map[string]*pubsub.Topic
	subscriptions map[string]*pubsub.Subscription

	// channels on which the pool's event loop
	// does dispatching/handling requests.
	enqueueReqCh chan enqueueRequest
	promoteReqCh chan promoteRequest
	//pruneCh      chan struct{}

	// shutdown channel
	shutdownCh chan struct{}

	// flag indicating if the current node is a sealer,
	// and should therefore gossip transactions
	//sealing uint32

	// networking stack
	topic   *network.Topic
	network *network.Server
}

// SetSigner sets the signer the rtc will use
// to validate a rtc msg's signature.
func (p *Rtc) SetSigner(s signer) {
	p.signer = s
}

func NewRtc(network *network.Server, logger hclog.Logger) (*Rtc, error) {
	rtc := &Rtc{
		logger:  logger.Named("rtc"),
		ctx:     context.Background(),
		network: network,
		stream:  &eventStream{},
		//topics:        make(map[string]*pubsub.Topic),
		//subscriptions: make(map[string]*pubsub.Subscription),
		//	main loop channels
		enqueueReqCh: make(chan enqueueRequest),
		promoteReqCh: make(chan promoteRequest),
		//pruneCh:      make(chan struct{}),
		shutdownCh: make(chan struct{}),
	}
	if network != nil {
		// subscribe to the gossip protocol
		topic, err := network.NewTopic(topicNameV1, &proto.RtcTelegram{})
		if err != nil {
			return nil, err
		}

		if subscribeErr := topic.Subscribe(rtc.addGossipMsg); subscribeErr != nil {
			return nil, fmt.Errorf("unable to subscribe to gossip topic, %w", subscribeErr)
		}

		rtc.topic = topic
	}

	return rtc, nil
}

//func (r *Rtc) join(subjectHash types.Hash) error {
//	r.Lock()
//	defer r.Unlock()
//
//	subjectTopic, err := r.ps.Join(subjectHash.String())
//	if err != nil {
//		return err
//	}
//	r.topics[subjectTopic.String()] = subjectTopic
//
//	sub, err := subjectTopic.Subscribe()
//	if err != nil {
//		return err
//	}
//	r.subscriptions[subjectTopic.String()] = sub
//
//	return nil
//}

// addGossipMsg handles receiving transactions
// gossiped by the network.
func (r *Rtc) addGossipMsg(obj interface{}, _ peer.ID) {

	raw, ok := obj.(*proto.RtcTelegram)
	if !ok {
		r.logger.Error("failed to cast gossiped message to telegram")

		return
	}

	// Verify that the gossiped rtc message is not empty
	if raw == nil || raw.Raw == nil {
		r.logger.Error("malformed gossip rtc telegram message received")

		return
	}

	// decode rtc telegram
	msg := new(RtcMsg)

	// decode telegram
	if err := msg.UnmarshalRLP(raw.Raw.Value); err != nil {
		r.logger.Error("failed to decode broadcast rtc telegram", "err", err)

		return
	}

	// add rtcMsg
	if err := r.addRtcMsg(gossip, msg); err != nil {
		if errors.Is(err, ErrAlreadyKnown) {
			r.logger.Debug("rejecting known rtc msg (gossip)", "hash", msg.Hash.String())

			return
		}

		r.logger.Error("failed to add broadcast rtc msg", "err", err, "From", msg.From, "Subject", msg.Subject)
	}
}

// dispatchEvent pushes a new event to the stream
func (r *Rtc) dispatchEvent(evnt *Event) {
	r.stream.push(evnt)
}

func (r *Rtc) AddRtcMsg(msg *RtcMsg) error {
	// validate incoming msg
	if err := r.validateRtcMsg(msg); err != nil {
		return err
	}

	// broadcast the RtcMsg only if a topic
	// subscription is present
	if r.topic != nil {
		msgTele := &proto.RtcTelegram{
			Raw: &anypb.Any{
				Value: msg.MarshalRLP(),
			},
		}

		if err := r.topic.Publish(msgTele); err != nil {
			r.logger.Error("failed to topic rtc message", "err", err)
		}
	}

	return nil
}

func (p *Rtc) Sender(msg *RtcMsg) (types.Address, error) {
	// Check if the rtcMsg is signed properly

	// Extract the sender
	from, signerErr := p.signer.Sender(msg)
	if signerErr != nil {
		return types.ZeroAddress, ErrExtractSignature
	}

	return from, nil
}

// validateTele ensures the rtcMsg conforms to specific
// constraints before publish the msg.
func (p *Rtc) validateRtcMsg(msg *RtcMsg) error {
	// Check the transaction size to overcome DOS Attacks
	if uint64(len(msg.MarshalRLP())) > rtcSlotSize {
		return ErrOversizedData
	}

	// Check if the rtcMsg is signed properly

	// Extract the sender
	from, signerErr := p.signer.Sender(msg)
	if signerErr != nil {
		return ErrExtractSignature
	}

	// testAddress
	//from := types.StringToAddress("0x68b95f67a32935e3ed85600F558b74E0d2747120")

	// If the from field is set, check that
	// it matches the signer
	if msg.From != types.ZeroAddress &&
		msg.From != from {
		return ErrInvalidSender
	}

	// If no address was set, update it
	if msg.From == types.ZeroAddress {
		msg.From = from
	}

	return nil
}

func (r *Rtc) addRtcMsg(origin msgOrigin, msg *RtcMsg) error {
	r.logger.Debug("add msg",
		"origin", origin.String(),
		"From", msg.From,
	)

	// validate incoming msg
	if err := r.validateRtcMsg(msg); err != nil {
		return err
	}

	// add to index
	//if ok := r.index.add(msg); !ok {
	//	return ErrAlreadyKnown
	//}

	// send request [BLOCKING]
	//r.enqueueReqCh <- enqueueRequest{msg: msg}

	event := &Event{}
	event.AddNewRtcMsg(msg)
	event.Type = EventNew
	event.Source = "rtc"
	r.dispatchEvent(event)
	return nil
}

func (r *Rtc) Start() {
	// set default value of txpool pending transactions gauge
	//r.updatePending(0)

	//	run the handler for high gauge level pruning
	//go func() {
	//	for {
	//		select {
	//		case <-r.shutdownCh:
	//			return
	//			//case <-r.pruneCh:
	//			//	r.pruneAccountsWithNonceHoles()
	//		}
	//
	//		//	handler is in cooldown to avoid successive calls
	//		//	which could be just no-ops
	//		time.Sleep(pruningCooldown)
	//	}
	//}()

	//	run the handler for the tx pipeline
	go func() {
		for {
			select {
			case <-r.shutdownCh:
				return
				//case req := <-r.enqueueReqCh:
				//	go r.handleEnqueueRequest(req)
				//case req := <-r.promoteReqCh:
				//	go r.handlePromoteRequest(req)
			}
		}
	}()
}

// Close shuts down the pool's main loop.
func (r *Rtc) Close() {
	r.shutdownCh <- struct{}{}
}

func (r *RtcMsg) Copy() *RtcMsg {
	tt := &RtcMsg{
		To:          r.To,
		Application: r.Application,
		From:        r.From,
		Subject:     r.Subject,
		Content:     r.Content,
	}
	*tt = *r

	if r.R != nil {
		tt.R = new(big.Int)
		tt.R = big.NewInt(0).SetBits(r.R.Bits())
	}

	if r.S != nil {
		tt.S = new(big.Int)
		tt.S = big.NewInt(0).SetBits(r.S.Bits())
	}

	return tt
}

func rtcTypeFromByte(b byte) (RtcType, error) {
	tt := RtcType(b)

	switch tt {
	case SubjectMsg, StateMsg, SubscribeMsg:
		return tt, nil
	default:
		return tt, fmt.Errorf("unknown rtc type: %d", b)
	}
}

func (t *RtcMsg) MarshalRLP() []byte {
	return t.MarshalRLPTo(nil)
}

func (t *RtcMsg) MarshalRLPTo(dst []byte) []byte {
	if t.Type != SubjectMsg {
		dst = append(dst, byte(t.Type))
	}

	return types.MarshalRLPTo(t.MarshalRLPWith, dst)
}

// MarshalRLPWith marshals the rtc msg to RLP with a specific fastrlp.Arena
func (t *RtcMsg) MarshalRLPWith(arena *fastrlp.Arena) *fastrlp.Value {
	vv := arena.NewArray()

	// Subject may be empty
	if len(t.Subject) > 0 {
		vv.Set(arena.NewString(t.Subject))
	} else {
		vv.Set(arena.NewNull())
	}

	// Application may be empty
	if len(t.Application) > 0 {
		vv.Set(arena.NewString(t.Application))
	} else {
		vv.Set(arena.NewNull())
	}

	// Content may be empty
	if len(t.Content) > 0 {
		vv.Set(arena.NewString(t.Content))
	} else {
		vv.Set(arena.NewNull())
	}

	if len(t.To) > 0 {
		vv.Set(arena.NewBytes(t.To.Bytes()))
	} else {
		vv.Set(arena.NewNull())
	}

	// signature values
	vv.Set(arena.NewBigInt(t.V))
	vv.Set(arena.NewBigInt(t.R))
	vv.Set(arena.NewBigInt(t.S))

	if t.Type == SubjectMsg {
		vv.Set(arena.NewBytes((t.From).Bytes()))
	}

	return vv
}

func (t *RtcMsg) UnmarshalRLP(input []byte) error {
	t.Type = SubjectMsg
	offset := 0

	if len(input) > 0 && input[0] <= types.RLPSingleByteUpperLimit {
		var err error
		if t.Type, err = rtcTypeFromByte(input[0]); err != nil {
			return err
		}

		offset = 1
	}

	return types.UnmarshalRlp(t.unmarshalRLPFrom, input[offset:])
}

// unmarshalRLPFrom unmarshals a rtc msg in RLP format
func (t *RtcMsg) unmarshalRLPFrom(p *fastrlp.Parser, v *fastrlp.Value) error {
	elems, err := v.GetElems()
	if err != nil {
		return err
	}

	if len(elems) < 7 {
		return fmt.Errorf("incorrect number of elements to decode rtcMsg, expected 7 but found %d", len(elems))
	}

	p.Hash(t.Hash[:0], v)

	//// nonce
	//if t.Nonce, err = elems[0].GetUint64(); err != nil {
	//	return err
	//}

	// Subject
	if t.Subject, err = elems[0].GetString(); err != nil {
		return err
	}

	// Application
	if t.Application, err = elems[1].GetString(); err != nil {
		return err
	}

	// Content
	if t.Content, err = elems[2].GetString(); err != nil {
		return err
	}

	// To
	if vv, err := elems[3].Bytes(); err == nil && len(vv) == types.AddressLength {
		// to address
		t.To = types.BytesToAddress(vv)
	}
	if err != nil {
		return err
	}

	// V
	t.V = new(big.Int)
	if err = elems[4].GetBigInt(t.V); err != nil {
		return err
	}

	// R
	t.R = new(big.Int)
	if err = elems[5].GetBigInt(t.R); err != nil {
		return err
	}

	// S
	t.S = new(big.Int)
	if err = elems[6].GetBigInt(t.S); err != nil {
		return err
	}

	// set From with default value
	t.From = types.ZeroAddress

	// From
	if len(elems) >= 8 {
		if vv, err := v.Get(7).Bytes(); err == nil && len(vv) == types.AddressLength {
			// address
			t.From = types.BytesToAddress(vv)
		}
	}
	return nil
}

// ComputeHash computes the hash of the rtcMsg
func (r *RtcMsg) ComputeHash() *RtcMsg {
	ar := marshalArenaPool.Get()
	hash := keccak.DefaultKeccakPool.Get()

	v := r.MarshalRLPWith(ar)
	hash.WriteRlp(r.Hash[:0], v)

	marshalArenaPool.Put(ar)
	keccak.DefaultKeccakPool.Put(hash)

	return r
}
