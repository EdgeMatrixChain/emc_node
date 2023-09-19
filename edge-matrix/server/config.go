package server

import (
	"github.com/emc-protocol/edge-matrix/chain"
	"net"

	"github.com/hashicorp/go-hclog"

	"github.com/emc-protocol/edge-matrix/network"
	"github.com/emc-protocol/edge-matrix/secrets"
)

const DefaultGRPCPort int = 50000
const DefaultJSONRPCPort int = 50002

// Config is used to parametrize the minimal client
type Config struct {
	Chain *chain.Chain

	JSONRPC    *JSONRPC
	GRPCAddr   *net.TCPAddr
	LibP2PAddr *net.TCPAddr
	RelayAddr  *net.TCPAddr // the relay address

	PriceLimit         uint64
	MaxAccountEnqueued uint64
	MaxSlots           uint64
	BlockTime          uint64

	Telemetry   *Telemetry
	Network     *network.Config
	EdgeNetwork *network.Config

	DataDir     string
	RestoreFile *string

	Seal bool

	SecretsManager *secrets.SecretsManagerConfig

	LogLevel hclog.Level

	JSONLogFormat bool

	LogFilePath string

	RelayOn        bool
	RelayDiscovery bool

	NumBlockConfirmations uint64

	AppName     string
	AppUrl      string
	AppOrigin   string
	RunningMode string

	IcHost        string
	MinerCanister string
	EmcHost       string

	PocCpu bool
	PocGpu bool
}

// Telemetry holds the config details for metric services
type Telemetry struct {
	PrometheusAddr *net.TCPAddr
}

// JSONRPC holds the config details for the JSON-RPC server
type JSONRPC struct {
	JSONRPCAddr              *net.TCPAddr
	AccessControlAllowOrigin []string
	BatchLengthLimit         uint64
	BlockRangeLimit          uint64
}
