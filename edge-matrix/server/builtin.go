package server

import (
	"github.com/emc-protocol/edge-matrix/consensus"
	consensusIBFT "github.com/emc-protocol/edge-matrix/consensus/ibft"
	"github.com/emc-protocol/edge-matrix/secrets"
	"github.com/emc-protocol/edge-matrix/secrets/awsssm"
	"github.com/emc-protocol/edge-matrix/secrets/gcpssm"
	"github.com/emc-protocol/edge-matrix/secrets/hashicorpvault"
	"github.com/emc-protocol/edge-matrix/secrets/local"
)

//type GenesisFactoryHook func(config *chain.Chain, engineName string) func(*state.Transition) error

type ConsensusType string

const (
	//DevConsensus     ConsensusType = "dev"
	IBFTConsensus ConsensusType = "ibft"
)

var consensusBackends = map[ConsensusType]consensus.Factory{
	//	DevConsensus:     consensusDev.Factory,
	IBFTConsensus: consensusIBFT.Factory,
}

// secretsManagerBackends defines the SecretManager factories for different
// secret management solutions
var secretsManagerBackends = map[secrets.SecretsManagerType]secrets.SecretsManagerFactory{
	secrets.Local:          local.SecretsManagerFactory,
	secrets.HashicorpVault: hashicorpvault.SecretsManagerFactory,
	secrets.AWSSSM:         awsssm.SecretsManagerFactory,
	secrets.GCPSSM:         gcpssm.SecretsManagerFactory,
}

func ConsensusSupported(value string) bool {
	_, ok := consensusBackends[ConsensusType(value)]

	return ok
}
