package genesis

import (
	"errors"
	"fmt"
	"time"

	"github.com/emc-protocol/edge-matrix/chain"
	"github.com/emc-protocol/edge-matrix/command"
	"github.com/emc-protocol/edge-matrix/command/helper"
	"github.com/emc-protocol/edge-matrix/consensus/ibft"
	"github.com/emc-protocol/edge-matrix/consensus/ibft/fork"
	"github.com/emc-protocol/edge-matrix/consensus/ibft/signer"
	"github.com/emc-protocol/edge-matrix/contracts/staking"
	stakingHelper "github.com/emc-protocol/edge-matrix/helper/staking"
	"github.com/emc-protocol/edge-matrix/server"
	"github.com/emc-protocol/edge-matrix/types"
	"github.com/emc-protocol/edge-matrix/validators"
)

const (
	dirFlag           = "dir"
	nameFlag          = "name"
	premineFlag       = "premine"
	chainIDFlag       = "chain-id"
	epochSizeFlag     = "epoch-size"
	epochRewardFlag   = "epoch-reward"
	blockGasLimitFlag = "block-gas-limit"
	posFlag           = "pos"
	minValidatorCount = "min-validator-count"
	maxValidatorCount = "max-validator-count"
)

// Legacy flags that need to be preserved for running clients
const (
	chainIDFlagLEGACY = "chainid"
)

var (
	params = &genesisParams{}
)

var (
	errValidatorsNotSpecified = errors.New("validator information not specified")
	errUnsupportedConsensus   = errors.New("specified consensusRaw not supported")
	errInvalidEpochSize       = errors.New("epoch size must be greater than 1")
)

type genesisParams struct {
	genesisPath         string
	name                string
	consensusRaw        string
	validatorPrefixPath string
	premine             []string
	bootnodes           []string
	ibftValidators      validators.Validators

	ibftValidatorsRaw []string

	chainID   uint64
	epochSize uint64

	blockGasLimit uint64
	isPos         bool

	minNumValidators uint64
	maxNumValidators uint64

	rawIBFTValidatorType string
	ibftValidatorType    validators.ValidatorType

	extraData []byte
	consensus server.ConsensusType

	consensusEngineConfig map[string]interface{}

	genesisConfig *chain.Chain

	// PolyBFT
	manifestPath           string
	smartContractsRootPath string
	validatorSetSize       int
	sprintSize             uint64
	blockTime              time.Duration
	bridgeJSONRPCAddr      string
	epochReward            uint64
}

func (p *genesisParams) validateFlags() error {
	// Check if the consensusRaw is supported
	if !server.ConsensusSupported(p.consensusRaw) {
		return errUnsupportedConsensus
	}

	// Check if validator information is set at all
	if p.isIBFTConsensus() &&
		!p.areValidatorsSetManually() &&
		!p.areValidatorsSetByPrefix() {
		return errValidatorsNotSpecified
	}

	// Check if the genesis file already exists
	if generateError := verifyGenesisExistence(p.genesisPath); generateError != nil {
		return errors.New(generateError.GetMessage())
	}

	// Check that the epoch size is correct
	if p.epochSize < 2 && (p.isIBFTConsensus()) {
		// Epoch size must be greater than 1, so new transactions have a chance to be added to a block.
		// Otherwise, every block would be an endblock (meaning it will not have any transactions).
		// Check is placed here to avoid additional parsing if epochSize < 2
		return errInvalidEpochSize
	}

	// Validate min and max validators number
	if err := command.ValidateMinMaxValidatorsNumber(p.minNumValidators, p.maxNumValidators); err != nil {
		return err
	}

	return nil
}

func (p *genesisParams) isIBFTConsensus() bool {
	return server.ConsensusType(p.consensusRaw) == server.IBFTConsensus
}

func (p *genesisParams) areValidatorsSetManually() bool {
	return len(p.ibftValidatorsRaw) != 0
}

func (p *genesisParams) areValidatorsSetByPrefix() bool {
	return p.validatorPrefixPath != ""
}

func (p *genesisParams) getRequiredFlags() []string {
	if p.isIBFTConsensus() {
		return []string{
			command.BootnodeFlag,
		}
	}

	return []string{}
}

func (p *genesisParams) initRawParams() error {
	p.consensus = server.ConsensusType(p.consensusRaw)

	if err := p.initIBFTValidatorType(); err != nil {
		return err
	}

	if err := p.initValidatorSet(); err != nil {
		return err
	}

	p.initIBFTExtraData()
	p.initConsensusEngineConfig()

	return nil
}

// setValidatorSetFromCli sets validator set from cli command
func (p *genesisParams) setValidatorSetFromCli() error {
	if len(p.ibftValidatorsRaw) == 0 {
		return nil
	}

	newValidators, err := validators.ParseValidators(p.ibftValidatorType, p.ibftValidatorsRaw)
	if err != nil {
		return err
	}

	if err = p.ibftValidators.Merge(newValidators); err != nil {
		return err
	}

	return nil
}

// setValidatorSetFromPrefixPath sets validator set from prefix path
//func (p *genesisParams) setValidatorSetFromPrefixPath() error {
//	if !p.areValidatorsSetByPrefix() {
//		return nil
//	}
//
//	validators, err := command.GetValidatorsFromPrefixPath(
//		p.validatorPrefixPath,
//		p.ibftValidatorType,
//	)
//
//	if err != nil {
//		return fmt.Errorf("failed to read from prefix: %w", err)
//	}
//
//	if err := p.ibftValidators.Merge(validators); err != nil {
//		return err
//	}
//
//	return nil
//}

func (p *genesisParams) initIBFTValidatorType() error {
	var err error
	if p.ibftValidatorType, err = validators.ParseValidatorType(p.rawIBFTValidatorType); err != nil {
		return err
	}

	return nil
}

func (p *genesisParams) initValidatorSet() error {
	p.ibftValidators = validators.NewValidatorSetFromType(p.ibftValidatorType)

	// Set validator set
	// Priority goes to cli command over prefix path
	//if err := p.setValidatorSetFromPrefixPath(); err != nil {
	//	return err
	//}

	if err := p.setValidatorSetFromCli(); err != nil {
		return err
	}

	// Validate if validator number exceeds max number
	if ok := p.isValidatorNumberValid(); !ok {
		return command.ErrValidatorNumberExceedsMax
	}

	return nil
}

func (p *genesisParams) isValidatorNumberValid() bool {
	return p.ibftValidators == nil || uint64(p.ibftValidators.Len()) <= p.maxNumValidators
}

func (p *genesisParams) initIBFTExtraData() {
	if p.consensus != server.IBFTConsensus {
		return
	}

	var committedSeal signer.Seals

	switch p.ibftValidatorType {
	case validators.ECDSAValidatorType:
		committedSeal = new(signer.SerializedSeal)
	case validators.BLSValidatorType:
		committedSeal = new(signer.AggregatedSeal)
	}

	ibftExtra := &signer.IstanbulExtra{
		Validators:     p.ibftValidators,
		ProposerSeal:   []byte{},
		CommittedSeals: committedSeal,
	}

	p.extraData = make([]byte, signer.IstanbulExtraVanity)
	p.extraData = ibftExtra.MarshalRLPTo(p.extraData)
}

func (p *genesisParams) initConsensusEngineConfig() {
	if p.consensus != server.IBFTConsensus {
		p.consensusEngineConfig = map[string]interface{}{
			p.consensusRaw: map[string]interface{}{},
		}

		return
	}

	if p.isPos {
		p.initIBFTEngineMap(fork.PoS)

		return
	}

	p.initIBFTEngineMap(fork.PoA)
}

func (p *genesisParams) initIBFTEngineMap(ibftType fork.IBFTType) {
	p.consensusEngineConfig = map[string]interface{}{
		string(server.IBFTConsensus): map[string]interface{}{
			fork.KeyType:          ibftType,
			fork.KeyValidatorType: p.ibftValidatorType,
			ibft.KeyEpochSize:     p.epochSize,
		},
	}
}

func (p *genesisParams) generateGenesis() error {
	if err := p.initGenesisConfig(); err != nil {
		return err
	}

	if err := helper.WriteGenesisConfigToDisk(
		p.genesisConfig,
		p.genesisPath,
	); err != nil {
		return err
	}

	return nil
}

func (p *genesisParams) initGenesisConfig() error {
	chainConfig := &chain.Chain{
		Name: p.name,
		Genesis: &chain.Genesis{
			GasLimit:   p.blockGasLimit,
			Difficulty: 1,
			Alloc:      map[types.Address]*chain.GenesisAccount{},
			ExtraData:  p.extraData,
			GasUsed:    command.DefaultGenesisGasUsed,
		},
		Params: &chain.Params{
			ChainID: int64(p.chainID),
			Forks:   chain.AllForksEnabled,
			Engine:  p.consensusEngineConfig,
		},
		Bootnodes: p.bootnodes,
	}

	// Predeploy staking smart contract if needed
	if p.shouldPredeployStakingSC() {
		stakingAccount, err := p.predeployStakingSC()
		if err != nil {
			return err
		}

		chainConfig.Genesis.Alloc[staking.AddrStakingContract] = stakingAccount
	}

	for _, premineRaw := range p.premine {
		premineInfo, err := parsePremineInfo(premineRaw)
		if err != nil {
			return err
		}

		chainConfig.Genesis.Alloc[premineInfo.address] = &chain.GenesisAccount{
			Balance: premineInfo.balance,
		}
	}

	p.genesisConfig = chainConfig

	return nil
}

func (p *genesisParams) shouldPredeployStakingSC() bool {
	// If the consensus selected is IBFT / Dev and the mechanism is Proof of Stake,
	// deploy the Staking SC
	return p.isPos && (p.consensus == server.IBFTConsensus)
}

func (p *genesisParams) predeployStakingSC() (*chain.GenesisAccount, error) {
	stakingAccount, predeployErr := stakingHelper.PredeployStakingSC(
		p.ibftValidators,
		stakingHelper.PredeployParams{
			MinValidatorCount: p.minNumValidators,
			MaxValidatorCount: p.maxNumValidators,
		})
	if predeployErr != nil {
		return nil, predeployErr
	}

	return stakingAccount, nil
}

func (p *genesisParams) getResult() command.CommandResult {
	return &GenesisResult{
		Message: fmt.Sprintf("Genesis written to %s\n", p.genesisPath),
	}
}
