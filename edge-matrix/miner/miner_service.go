package miner

import (
	"context"
	"crypto/ed25519"
	"github.com/emc-protocol/edge-matrix/crypto"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/helper/ic/agent"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/idl"
	"github.com/emc-protocol/edge-matrix/miner/proto"
	"github.com/emc-protocol/edge-matrix/secrets"
	"github.com/emc-protocol/edge-matrix/secrets/helper"
	"github.com/libp2p/go-libp2p/core/host"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MinerService struct {
	proto.UnimplementedMinerServer

	host           host.Host
	agent          *agent.Agent
	secretsManager secrets.SecretsManager
}

func NewMinerService(agent *agent.Agent, host host.Host, secretsManager secrets.SecretsManager) *MinerService {
	return &MinerService{
		agent:          agent,
		host:           host,
		secretsManager: secretsManager,
	}
}

// PeersStatus implements the 'peers status' operator service
func (s *MinerService) GetMiner() (*proto.MinerStatus, error) {
	if !s.secretsManager.HasSecret(secrets.ICPIdentityKey) {
		helper.InitICPIdentityKey(s.secretsManager)
	}
	icPrivKey, err := s.secretsManager.GetSecret(secrets.ICPIdentityKey)
	if err != nil {
		return nil, err
	}
	decodedPrivKey, err := crypto.BytesToEd25519PrivateKey(icPrivKey)
	decodedPubKey := make([]byte, ed25519.PublicKeySize)
	copy(decodedPubKey, decodedPrivKey[ed25519.PublicKeySize:])

	// TODO query status from IC canister

	status := proto.MinerStatus{
		NetName:  "IC",
		PeerId:   s.host.ID().String(),
		IcPubKey: hex.EncodeToString(decodedPubKey),
	}

	return &status, nil
}

// PeersStatus implements the 'peers status' operator service
func (s *MinerService) GetMinerStatus(context.Context, *emptypb.Empty) (*proto.MinerStatus, error) {
	return s.GetMiner()
}

// Regiser set or remove a address
func (s *MinerService) MinerRegiser(ctx context.Context, req *proto.MinerRegisterRequest) (*proto.MinerRegisterResponse, error) {
	// TODO call IC canister
	// principalId := req.Id
	// peerId :=s.Host.ID().String()
	canister := "xb3xh-uaaaa-aaaam-abi3a-cai"
	methodName := "greet"

	var argType []idl.Type
	argType = append(argType, new(idl.Text))

	var argValue []interface{}
	argValue = append(argValue, req.Id)

	result := ""
	arg, _ := idl.Encode(argType, argValue)
	_, resp, _, err := s.agent.Query(canister, methodName, arg)
	if err != nil {
		result = err.Error()
	}
	result = resp[0].(string)
	response := proto.MinerRegisterResponse{
		Message: result,
	}
	return &response, nil
}
