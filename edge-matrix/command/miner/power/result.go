package power

import (
	"bytes"
	"fmt"

	"github.com/emc-protocol/edge-matrix/command/helper"
)

type CurrentEPowerResult struct {
	Round uint64 `json:"round"`
	Total uint64 `json:"total"`
}

func (r *CurrentEPowerResult) GetOutput() string {
	var buffer bytes.Buffer
	buffer.WriteString("\n[MINER e-Power]\n")
	buffer.WriteString(helper.FormatKV([]string{
		fmt.Sprintf("Round |%d", r.Round),
		fmt.Sprintf("Total |%d E", r.Total),
	}))
	buffer.WriteString("\n")

	return buffer.String()
}
