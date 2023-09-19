package server

import (
	"errors"
	"github.com/emc-protocol/edge-matrix/chain"
	"net"

	"github.com/emc-protocol/edge-matrix/command/server/config"
	"github.com/emc-protocol/edge-matrix/network"
	"github.com/emc-protocol/edge-matrix/secrets"
	"github.com/emc-protocol/edge-matrix/server"
	"github.com/hashicorp/go-hclog"
	"github.com/multiformats/go-multiaddr"
)

const (
	configFlag                   = "config"
	genesisPathFlag              = "chain"
	dataDirFlag                  = "data-dir"
	libp2pAddressFlag            = "base-libp2p"
	edgeLibp2pAddressFlag        = "libp2p"
	relayLibp2pAddressFlag       = "relay-libp2p"
	prometheusAddressFlag        = "prometheus"
	natFlag                      = "nat"
	dnsFlag                      = "dns"
	sealFlag                     = "seal"
	maxPeersFlag                 = "max-peers"
	maxInboundPeersFlag          = "max-inbound-peers"
	maxOutboundPeersFlag         = "max-outbound-peers"
	priceLimitFlag               = "price-limit"
	jsonRPCBatchRequestLimitFlag = "json-rpc-batch-request-limit"
	jsonRPCBlockRangeLimitFlag   = "json-rpc-block-range-limit"
	maxSlotsFlag                 = "max-slots"
	maxEnqueuedFlag              = "max-enqueued"
	blockGasTargetFlag           = "block-gas-target"
	secretsConfigFlag            = "secrets-config"
	restoreFlag                  = "restore"
	blockTimeFlag                = "block-time"
	devIntervalFlag              = "dev-interval"
	devFlag                      = "dev"
	corsOriginFlag               = "access-control-allow-origins"
	logFileLocationFlag          = "log-to"

	numBlockConfirmationsFlag = "num-block-confirmations"

	relayOnFlag        = "relay-on"
	relayDiscoveryFlag = "relay-discovery"
	runningModeFlag    = "running-mode"
	appNameFlag        = "app-name"
	//appUrlFlag         = "app-url"
	appOriginFlag      = "app-origin"
	icHostFlag         = "ic-host"
	minerCanistertFlag = "miner-canister"

	pocCpuFlag = "poc-cpu"
	pocGpuFlag = "poc-gpu"
)

// Flags that are deprecated, but need to be preserved for
// backwards compatibility with existing scripts
const (
	ibftBaseTimeoutFlagLEGACY = "ibft-base-timeout"
)

const (
	unsetPeersValue = -1
)

var (
	params = &serverParams{
		rawConfig: &config.Config{
			Telemetry: &config.Telemetry{},
			Network:   &config.Network{},
			TelePool:  &config.TelePool{},
		},
	}
)

var (
	errInvalidNATAddress = errors.New("could not parse NAT IP address")
)

type serverParams struct {
	rawConfig  *config.Config
	configPath string

	libp2pAddress      *net.TCPAddr
	edgeLibp2pAddress  *net.TCPAddr
	relayLibp2pAddress *net.TCPAddr
	prometheusAddress  *net.TCPAddr
	natAddress         net.IP
	dnsAddress         multiaddr.Multiaddr
	grpcAddress        *net.TCPAddr
	jsonRPCAddress     *net.TCPAddr

	blockGasTarget uint64
	devInterval    uint64
	isDevMode      bool

	corsAllowedOrigins []string

	ibftBaseTimeoutLegacy uint64

	genesisConfig *chain.Chain
	secretsConfig *secrets.SecretsManagerConfig

	logFileLocation string
}

func (p *serverParams) isMaxPeersSet() bool {
	return p.rawConfig.Network.MaxPeers != unsetPeersValue
}

func (p *serverParams) isPeerRangeSet() bool {
	return p.rawConfig.Network.MaxInboundPeers != unsetPeersValue ||
		p.rawConfig.Network.MaxOutboundPeers != unsetPeersValue
}

func (p *serverParams) isSecretsConfigPathSet() bool {
	return p.rawConfig.SecretsConfigPath != ""
}

func (p *serverParams) isPrometheusAddressSet() bool {
	return p.rawConfig.Telemetry.PrometheusAddr != ""
}

func (p *serverParams) isNATAddressSet() bool {
	return p.rawConfig.Network.NatAddr != ""
}

func (p *serverParams) isDNSAddressSet() bool {
	return p.rawConfig.Network.DNSAddr != ""
}

func (p *serverParams) isLogFileLocationSet() bool {
	return p.rawConfig.LogFilePath != ""
}

//func (p *serverParams) isDevConsensus() bool {
//	return server.ConsensusType(p.genesisConfig.Params.GetEngine()) == server.DevConsensus
//}

func (p *serverParams) getRestoreFilePath() *string {
	if p.rawConfig.RestoreFile != "" {
		return &p.rawConfig.RestoreFile
	}

	return nil
}

func (p *serverParams) setRawGRPCAddress(grpcAddress string) {
	p.rawConfig.GRPCAddr = grpcAddress
}

func (p *serverParams) setRawJSONRPCAddress(jsonRPCAddress string) {
	p.rawConfig.JSONRPCAddr = jsonRPCAddress
}

func (p *serverParams) setJSONLogFormat(jsonLogFormat bool) {
	p.rawConfig.JSONLogFormat = jsonLogFormat
}

func (p *serverParams) generateConfig() *server.Config {
	return &server.Config{
		Chain: p.genesisConfig,
		JSONRPC: &server.JSONRPC{
			JSONRPCAddr:              p.jsonRPCAddress,
			AccessControlAllowOrigin: p.corsAllowedOrigins,
			BatchLengthLimit:         p.rawConfig.JSONRPCBatchRequestLimit,
			BlockRangeLimit:          p.rawConfig.JSONRPCBlockRangeLimit,
		},
		GRPCAddr:   p.grpcAddress,
		LibP2PAddr: p.libp2pAddress,
		Telemetry: &server.Telemetry{
			PrometheusAddr: p.prometheusAddress,
		},
		Network: &network.Config{
			NoDiscover:       p.rawConfig.Network.NoDiscover,
			Addr:             p.libp2pAddress,
			NatAddr:          p.natAddress,
			DNS:              p.dnsAddress,
			DataDir:          p.rawConfig.DataDir,
			MaxPeers:         p.rawConfig.Network.MaxPeers,
			MaxInboundPeers:  p.rawConfig.Network.MaxInboundPeers,
			MaxOutboundPeers: p.rawConfig.Network.MaxOutboundPeers,
			Chain:            p.genesisConfig,
		},
		EdgeNetwork: &network.Config{
			NoDiscover:       p.rawConfig.Network.NoDiscover,
			Addr:             p.edgeLibp2pAddress,
			NatAddr:          p.natAddress,
			DNS:              p.dnsAddress,
			DataDir:          p.rawConfig.DataDir,
			MaxPeers:         p.rawConfig.Network.MaxPeers,
			MaxInboundPeers:  p.rawConfig.Network.MaxInboundPeers,
			MaxOutboundPeers: p.rawConfig.Network.MaxOutboundPeers,
			Chain:            p.genesisConfig,
		},
		RelayAddr: p.relayLibp2pAddress,
		DataDir:   p.rawConfig.DataDir,
		Seal:      p.rawConfig.ShouldSeal,
		//PriceLimit:         p.rawConfig.TelePool.PriceLimit,
		MaxSlots:           p.rawConfig.TelePool.MaxSlots,
		MaxAccountEnqueued: p.rawConfig.TelePool.MaxAccountEnqueued,
		SecretsManager:     p.secretsConfig,
		RestoreFile:        p.getRestoreFilePath(),
		BlockTime:          p.rawConfig.BlockTime,
		LogLevel:           hclog.LevelFromString(p.rawConfig.LogLevel),
		JSONLogFormat:      p.rawConfig.JSONLogFormat,
		LogFilePath:        p.logFileLocation,

		RelayOn:               p.rawConfig.RelayOn,
		RelayDiscovery:        p.rawConfig.RelayDiscovery,
		NumBlockConfirmations: p.rawConfig.NumBlockConfirmations,

		RunningMode: p.rawConfig.RunningMode,
		AppName:     p.rawConfig.AppName,
		AppUrl:      "http://127.0.0.1:9527",
		AppOrigin:   p.rawConfig.AppOrigin,

		IcHost:        p.rawConfig.IcHost,
		MinerCanister: p.rawConfig.MinerCanister,
		EmcHost:       p.rawConfig.EmcHost,

		PocCpu: p.rawConfig.PocCpu,
		PocGpu: p.rawConfig.PocGpu,
	}
}
