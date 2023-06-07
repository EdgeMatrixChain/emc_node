package power

import (
	"bytes"
	"fmt"

	"github.com/emc-protocol/edge-matrix/command/helper"
)

type CurrentEPowerResult struct {
	Round    uint64  `json:"round"`
	Total    float32 `json:"total"`
	Multiple float32 `json:"multiple"`
}

func (r *CurrentEPowerResult) GetOutput() string {
	var buffer bytes.Buffer
	buffer.WriteString("\n[MINER e-Power]\n")
	buffer.WriteString(helper.FormatKV([]string{
		fmt.Sprintf("Round |%d", r.Round),
		fmt.Sprintf("Total |%.8f E", r.Total),
		fmt.Sprintf("Multiple |%.8f", r.Multiple),
	}))
	buffer.WriteString("\n")

	return buffer.String()
}
