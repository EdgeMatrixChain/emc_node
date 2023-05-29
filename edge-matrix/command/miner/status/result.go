package status

import (
	"bytes"
	"fmt"

	"github.com/emc-protocol/edge-matrix/command/helper"
)

type MinerStatusResult struct {
	NetName      string `json:"net_name"`
	NodeID       string `json:"node_id"`
	NodeIdentity string `json:"node_identity"`
	Principal    string `json:"principal"`
}

func (r *MinerStatusResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[MINER STATUS]\n")
	buffer.WriteString(helper.FormatKV([]string{
		fmt.Sprintf("NetName |%s", r.NetName),
		fmt.Sprintf("NodeID |%s", r.NodeID),
		fmt.Sprintf("NodeIdentity |%s", r.NodeIdentity),
		fmt.Sprintf("Principal |%s", r.Principal),
	}))
	buffer.WriteString("\n")

	return buffer.String()
}
