package hub

import (
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"github.com/emc-protocol/edge-matrix/helper/ic/agent"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/identity"
	"github.com/emc-protocol/edge-matrix/helper/ic/utils/principal"
	"github.com/hashicorp/go-hclog"
	"testing"
)

const (
	TestPrivKey = "ed1cb741ef10f2e353c6c395f0b91270e762e55cf30d03ab6fa340ff306fb9d9"
)

func Test_Idendity(t *testing.T) {
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))
}

func Test_AddModel(t *testing.T) {
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	hubAgent := NewHubAgent(hclog.NewNullLogger(), icAgent)

	err = hubAgent.AddModel(
		"StableDiffussion", "hash0000004")
	if err != nil {
		t.Log(err)
	}
}

func Test_RemoveModel(t *testing.T) {
	icAgent := agent.NewWithHost(DEFAULT_IC_HOST, false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	hubAgent := NewHubAgent(hclog.NewNullLogger(), icAgent)

	err = hubAgent.RemoveModel(
		"hash0000002")
	if err != nil {
		t.Log(err)
	}
}

func Test_IsModelListed(t *testing.T) {
	icAgent := agent.NewWithHost("https://ic0.app", false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	hubAgent := NewHubAgent(hclog.NewNullLogger(), icAgent)

	isValidModel, err := hubAgent.IsModelListed("hash0000004")
	if err != nil {
		t.Log(err)
	}
	t.Log(isValidModel)
}

func Test_ListModels(t *testing.T) {
	icAgent := agent.NewWithHost("https://ic0.app", false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	hubAgent := NewHubAgent(hclog.NewNullLogger(), icAgent)

	models, err := hubAgent.ListModels()
	if err != nil {
		t.Log(err)
	}
	t.Log(len(models))
	for _, modelHash := range models {
		t.Log(modelHash)
	}
}

func Test_ListModelsByeType(t *testing.T) {
	icAgent := agent.NewWithHost("https://ic0.app", false, TestPrivKey)
	privKeyBytes, err := hex.DecodeHex(TestPrivKey)
	if err != nil {
		return
	}
	identity := identity.New(false, privKeyBytes)
	p := principal.NewSelfAuthenticating(identity.PubKeyBytes())
	t.Log("identity:", p.Encode(), len(identity.PubKeyBytes()))

	hubAgent := NewHubAgent(hclog.NewNullLogger(), icAgent)

	models, err := hubAgent.ListModelsByeType("openai")
	if err != nil {
		t.Log(err)
	}
	t.Log("len:", len(models))
	for _, modelHash := range models {
		t.Log(modelHash)
	}
}
