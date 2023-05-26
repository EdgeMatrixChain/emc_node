package agent_test

import (
	"github.com/emc-protocol/edge-matrix/helper/ic/agent"
	"testing"
)

func TestClientStatus(t *testing.T) {
	c := agent.NewClient("https://ic0.app")
	status, _ := c.Status()
	t.Log(status.Version)
	// Output:
	// 0.18.0
}
