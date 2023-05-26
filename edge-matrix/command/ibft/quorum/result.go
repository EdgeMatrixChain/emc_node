package quorum

import (
	"bytes"
	"fmt"

	"github.com/emc-protocol/edge-matrix/command/helper"
	"github.com/emc-protocol/edge-matrix/helper/common"
)

type IBFTQuorumResult struct {
	Chain string            `json:"chain"`
	From  common.JSONNumber `json:"from"`
}

func (r *IBFTQuorumResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[NEW IBFT QUORUM START]\n")

	outputs := []string{
		fmt.Sprintf("Chain|%s", r.Chain),
		fmt.Sprintf("From|%d", r.From.Value),
	}

	buffer.WriteString(helper.FormatKV(outputs))
	buffer.WriteString("\n")

	return buffer.String()
}
