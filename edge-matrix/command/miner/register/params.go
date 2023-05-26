package register

import (
	"context"
	"errors"
	"github.com/emc-protocol/edge-matrix/command"
	"github.com/emc-protocol/edge-matrix/command/helper"
	minerOp "github.com/emc-protocol/edge-matrix/miner/proto"
)

const (
	commitFlag  = "commit"
	addressFlag = "addr"
	nodeFlag    = "node"
)

const (
	setOpt    = "set"
	removeOpt = "remove"
)

const (
	validatorNodeOpt = "validator"
	routeNodeOpt     = "router"
	edgeNodeOpt      = "edge"
)

var (
	errInvalidCommitType    = errors.New("invalid register type")
	errInvalidNodeType      = errors.New("invalid node type")
	errInvalidAddressFormat = errors.New("invalid address format")
)

var (
	params = &registerParams{}
)

type registerParams struct {
	addressRaw string

	commit   string
	nodeType string
	address  string
	message  string
}

func (p *registerParams) getRequiredFlags() []string {
	return []string{
		commitFlag,
		addressFlag,
		nodeFlag,
	}
}

func (p *registerParams) validateFlags() error {
	if !isValidCommitType(p.commit) {
		return errInvalidCommitType
	}
	if !isValidNodeType(p.nodeType) {
		return errInvalidCommitType
	}

	return nil
}

func (p *registerParams) initRawParams() error {
	if err := p.initAddress(); err != nil {
		return err
	}

	return nil
}

func (p *registerParams) initAddress() error {
	p.address = p.addressRaw
	return nil
}

func isValidCommitType(commit string) bool {
	return commit == setOpt || commit == removeOpt
}

func isValidNodeType(node string) bool {
	return node == routeNodeOpt || node == edgeNodeOpt || node == validatorNodeOpt
}

func (p *registerParams) registerMinerAddress(grpcAddress string) error {
	minerClient, err := helper.GetMinerClientConnection(grpcAddress)
	if err != nil {
		return err
	}

	result, err := minerClient.MinerRegiser(
		context.Background(),
		p.getRegisterUpdate(),
	)
	if err != nil {
		p.message = err.Error()
	} else {
		p.message = result.Message
	}
	return nil
}

func (p *registerParams) getRegisterUpdate() *minerOp.MinerRegisterRequest {
	req := &minerOp.MinerRegisterRequest{
		Id: p.address,
	}
	return req
}

func (p *registerParams) getResult() command.CommandResult {
	return &MinerRegisterResult{
		Address:      p.address,
		Commit:       p.commit,
		NodeType:     p.nodeType,
		ResultMessge: p.message,
	}
}
