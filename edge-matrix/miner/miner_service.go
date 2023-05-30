package miner

import (
	"context"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/identity"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/principal"
	"github.com/emc-protocol/edge-matrix/miner/proto"
	"github.com/emc-protocol/edge-matrix/secrets"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p/core/host"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MinerService struct {
	proto.UnimplementedMinerServer
	logger         hclog.Logger
	icHost         string
	host           host.Host
	secretsManager secrets.SecretsManager

	// agent for communicating with IC Miner Canister
	minerAgent *MinerAgent
}

func NewMinerService(logger hclog.Logger, minerAgent *MinerAgent, host host.Host, secretsManager secrets.SecretsManager) *MinerService {
	return &MinerService{
		logger:         logger,
		minerAgent:     minerAgent,
		host:           host,
		secretsManager: secretsManager,
	}
}

// GetMiner return miner's status from secretsManager and IC canister
func (s *MinerService) GetMiner() (*proto.MinerStatus, error) {
	identity := s.GetIdentity()
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	s.logger.Debug("MinerRegiser", "node_identity", p.Encode())

	// query node from IC canister
	wallet, ntype, err := s.minerAgent.MyNode(s.host.ID().String())
	if err != nil {
		return nil, err
	}
	nodeType := ""
	if ntype > -1 {
		switch NodeType(ntype) {
		case NodeTypeRouter:
			nodeType = "router"
		case NodeTypeValidator:
			nodeType = "validator"
		case NodeTypeComputing:
			nodeType = "computing"
		default:
		}
	}

	status := proto.MinerStatus{
		NetName:      "IC",
		NodeId:       s.host.ID().String(),
		NodeIdentity: p.Encode(),
		Principal:    wallet,
		NodeType:     nodeType,
	}
	return &status, nil
}

// PeersStatus implements the 'peers status' operator service
func (s *MinerService) GetMinerStatus(context.Context, *emptypb.Empty) (*proto.MinerStatus, error) {
	return s.GetMiner()
}

func (s *MinerService) GetIdentity() *identity.Identity {
	icPrivKey, err := s.secretsManager.GetSecret(secrets.ICPIdentityKey)
	if err != nil {
		return nil
	}
	decodedPrivKey, err := crypto.BytesToEd25519PrivateKey(icPrivKey)
	identity := identity.New(false, decodedPrivKey.Seed())
	return identity
}

// Regiser set or remove a principal for miner
func (s *MinerService) MinerRegiser(ctx context.Context, req *proto.MinerRegisterRequest) (*proto.MinerRegisterResponse, error) {
	identity := s.GetIdentity()
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	s.logger.Info("MinerRegiser", "node identity", p.Encode(), "NodeId", s.host.ID().String(), "Principal", req.Principal)

	result := "register ok"
	err := s.minerAgent.RegisterNode(NodeType(req.Type), s.host.ID().String(), req.Principal)
	if err != nil {
		result = err.Error()
	}
	// TODO update minerFlag in application endpoint

	response := proto.MinerRegisterResponse{
		Message: result,
	}
	return &response, nil
}
