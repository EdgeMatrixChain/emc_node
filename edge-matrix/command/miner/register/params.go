package register

import (
	"context"
	"errors"
	"github.com/emc-protocol/edge-matrix/command"
	"github.com/emc-protocol/edge-matrix/command/helper"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/miner"
	minerOp "github.com/emc-protocol/edge-matrix/miner/proto"
	"github.com/emc-protocol/edge-matrix/types"
	"regexp"
)

const (
	commitFlag    = "commit"
	principalFlag = "principal"
	nodeFlag      = "node"
)

const (
	setOpt    = "set"
	removeOpt = "remove"
)

const (
	validatorNodeOpt = "validator"
	routeNodeOpt     = "router"
	computingNodeOpt = "computing"
)

var (
	errInvalidCommitType = errors.New("invalid commit type")
	errInvalidNodeType   = errors.New("invalid node type")
	errInvalidPrincipal  = errors.New("invalid ethereum address")
)

var (
	params = &registerParams{}
)

type registerParams struct {
	commit    string
	nodeType  string
	principal string
	message   string
}

func (p *registerParams) getRequiredFlags() []string {
	return []string{
		commitFlag,
	}
}

func (p *registerParams) validateFlags() error {
	if !isValidCommitType(p.commit) {
		return errInvalidCommitType
	}
	if !isValidNodeType(p.nodeType) {
		return errInvalidNodeType
	}
	if p.commit == setOpt {
		if !isValidPrincipal(p.principal) {
			return errInvalidPrincipal
		}
	}
	return nil
}

func isValidCommitType(commit string) bool {
	return commit == setOpt || commit == removeOpt
}

func isValidNodeType(node string) bool {
	//return node == computingNodeOpt
	return node == routeNodeOpt || node == computingNodeOpt || node == validatorNodeOpt
}

func isValidPrincipal(principal string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	isHex := re.MatchString(principal)
	if !isHex {
		return false
	}
	bytes := hex.DropHexPrefix([]byte(principal))
	return len(bytes) == 2*types.AddressLength
}

func (p *registerParams) registerMinerAddress(grpcAddress string) error {
	minerClient, err := helper.GetMinerClientConnection(grpcAddress)
	if err != nil {
		return err
	}

	if p.commit == setOpt {
		result, err := minerClient.MinerRegiser(
			context.Background(),
			p.getRegisterUpdate(),
		)
		if err != nil {
			p.message = err.Error()
		} else {
			p.message = result.Message
		}
	} else if p.commit == removeOpt {
		result, err := minerClient.MinerRegiser(
			context.Background(),
			p.getRegisterUpdate(),
		)
		if err != nil {
			p.message = err.Error()
		} else {
			p.message = result.Message
		}
	}
	return nil
}

func (p *registerParams) getRegisterUpdate() *minerOp.MinerRegisterRequest {
	nodeType := miner.NodeTypeComputing
	if p.nodeType == routeNodeOpt {
		nodeType = miner.NodeTypeRouter
	} else if p.nodeType == computingNodeOpt {
		nodeType = miner.NodeTypeComputing
	} else if p.nodeType == validatorNodeOpt {
		nodeType = miner.NodeTypeValidator
	}
	req := &minerOp.MinerRegisterRequest{
		Principal: p.principal,
		Type:      uint64(nodeType),
		Commit:    p.commit,
	}
	return req
}

func (p *registerParams) getResult() command.CommandResult {
	return &MinerRegisterResult{
		Principal:    p.principal,
		Commit:       p.commit,
		NodeType:     p.nodeType,
		ResultMessge: p.message,
	}
}
