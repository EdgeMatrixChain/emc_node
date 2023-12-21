package miner

import (
	"context"
	"github.com/emc-protocol/edge-matrix/miner/proto"
	"github.com/emc-protocol/edge-matrix/secrets"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p/core/host"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	setOpt    = "set"
	removeOpt = "remove"
)

type MinerService struct {
	proto.UnimplementedMinerServer
	logger         hclog.Logger
	host           host.Host
	secretsManager secrets.SecretsManager

	// agent for communicating with EMC Hub
	minerAgent *MinerHubAgent
}

func NewMinerService(logger hclog.Logger, minerAgent *MinerHubAgent, host host.Host, secretsManager secrets.SecretsManager) *MinerService {
	return &MinerService{
		logger:         logger,
		minerAgent:     minerAgent,
		host:           host,
		secretsManager: secretsManager,
	}
}

// GetMiner return miner's status from secretsManager and IC canister
func (s *MinerService) GetMiner() (*proto.MinerStatus, error) {
	// query node from IC canister
	nodeId, nodeIdentity, wallet, registered, nodeType, err := s.minerAgent.MyNode(s.host.ID().String())
	if err != nil {
		return nil, err
	}

	status := proto.MinerStatus{
		NetName:      "Arbitrum One",
		NodeId:       nodeId,
		NodeIdentity: nodeIdentity,
		Principal:    wallet,
		NodeType:     nodeType,
		Registered:   registered,
	}
	return &status, nil
}

func (s *MinerService) GetCurrentEPower(context.Context, *emptypb.Empty) (*proto.CurrentEPower, error) {
	round, power, err := s.minerAgent.MyCurrentEPower(s.host.ID().String())
	if err != nil {
		return nil, err
	}
	_, _, multiple, err := s.minerAgent.MyStack(s.host.ID().String())
	if err != nil {
		return nil, err
	}

	ePower := proto.CurrentEPower{
		Round:    round,
		Total:    power,
		Multiple: float32(multiple) / 10000.0,
	}
	return &ePower, nil
}

// PeersStatus implements the 'peers status' operator service
func (s *MinerService) GetMinerStatus(context.Context, *emptypb.Empty) (*proto.MinerStatus, error) {
	return s.GetMiner()
}

// Regiser set or remove a principal for miner
func (s *MinerService) MinerRegiser(ctx context.Context, req *proto.MinerRegisterRequest) (*proto.MinerRegisterResponse, error) {
	result := ""

	if req.Commit == setOpt {
		result = "register ok"

		switch NodeType(req.Type) {
		case NodeTypeRouter:
			err := s.minerAgent.RegisterRouterNode(s.host.ID().String(), req.Principal)
			if err != nil {
				result = err.Error()
			}
		case NodeTypeValidator:
			err := s.minerAgent.RegisterValidatorNode(s.host.ID().String(), req.Principal)
			if err != nil {
				result = err.Error()
			}
		case NodeTypeComputing:
			err := s.minerAgent.RegisterComputingNode(s.host.ID().String(), req.Principal)
			if err != nil {
				result = err.Error()
			}
		default:
		}

	} else if req.Commit == removeOpt {
		result = "unregister ok"

		switch NodeType(req.Type) {
		case NodeTypeRouter:
			err := s.minerAgent.UnRegisterRouterNode(s.host.ID().String())
			if err != nil {
				result = err.Error()
			}
		case NodeTypeValidator:
			err := s.minerAgent.UnregisterValidatorNode(s.host.ID().String())
			if err != nil {
				result = err.Error()
			}
		case NodeTypeComputing:
			err := s.minerAgent.UnRegisterComputingNode(s.host.ID().String())
			if err != nil {
				result = err.Error()
			}
		default:
		}

	}
	// TODO update minerFlag in application endpoint

	response := proto.MinerRegisterResponse{
		Message: result,
	}
	return &response, nil
}
