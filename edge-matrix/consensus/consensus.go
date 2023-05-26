package consensus

import (
	"context"
	"github.com/emc-protocol/edge-matrix/telepool"
	"github.com/emc-protocol/edge-matrix/validators"
	"log"

	"github.com/emc-protocol/edge-matrix/blockchain"
	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/emc-protocol/edge-matrix/helper/progress"
	"github.com/emc-protocol/edge-matrix/network"
	"github.com/emc-protocol/edge-matrix/secrets"
	"github.com/emc-protocol/edge-matrix/state"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

// Consensus is the public interface for consensus mechanism
// Each consensus mechanism must implement this interface in order to be valid
type Consensus interface {
	// VerifyHeader verifies the header is correct
	VerifyHeader(header *types.Header) error

	// ProcessHeaders updates the snapshot based on the verified headers
	ProcessHeaders(headers []*types.Header) error

	// GetBlockCreator retrieves the block creator (or signer) given the block header
	GetBlockCreator(header *types.Header) (types.Address, error)

	// PreCommitState a hook to be called before finalizing state transition on inserting block
	PreCommitState(header *types.Header, txn *state.Transition) error

	// GetSyncProgression retrieves the current sync progression, if any
	GetSyncProgression() *progress.Progression

	// GetBridgeProvider returns an instance of BridgeDataProvider
	GetBridgeProvider() BridgeDataProvider

	// Initialize initializes the consensus (e.g. setup data)
	Initialize() error

	// Start starts the consensus and servers
	Start() error

	// Close closes the connection
	Close() error

	// Get current validators
	GetCurrentValidators() validators.Validators

	// Get singer address
	GetSignerAddress() types.Address
}

// Config is the configuration for the consensus
type Config struct {
	// Logger to be used by the consensus
	Logger *log.Logger

	// Params are the params of the chain and the consensus
	Params *chain.Params

	// Config defines specific configuration parameters for the consensus
	Config map[string]interface{}

	// Path is the directory path for the consensus protocol to store information
	Path string
}

type Params struct {
	Context        context.Context
	Config         *Config
	TelePool       *telepool.TelegramPool
	Network        *network.Server
	Blockchain     *blockchain.Blockchain
	Executor       *state.Executor
	Grpc           *grpc.Server
	Logger         hclog.Logger
	SecretsManager secrets.SecretsManager
	BlockTime      uint64

	NumBlockConfirmations uint64
}

// Factory is the factory function to create a discovery consensus
type Factory func(*Params) (Consensus, error)

// BridgeDataProvider is an interface providing bridge related functions
type BridgeDataProvider interface {
	// GenerateExit proof generates proof of exit for given exit event
	GenerateExitProof(exitID, epoch, checkpointBlock uint64) (types.Proof, error)

	// GetStateSyncProof retrieves the StateSync proof
	GetStateSyncProof(stateSyncID uint64) (types.Proof, error)
}
