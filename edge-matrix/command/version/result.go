package version

import (
	"bytes"
	"fmt"
	"github.com/emc-protocol/edge-matrix/command/helper"
)

type VersionResult struct {
	Version string `json:"version"`
	Build   string `json:"build"`
}

func (r *VersionResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[VERSION INFO]\n")
	buffer.WriteString(helper.FormatKV([]string{
		fmt.Sprintf("Release version|%s\n", r.Version),
		fmt.Sprintf("Build version|%s\n", r.Build),
	}))

	return buffer.String()
}
