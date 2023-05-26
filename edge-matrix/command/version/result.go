package version

import (
	"bytes"
	"fmt"
	"github.com/emc-protocol/edge-matrix/command/helper"
)

type VersionResult struct {
	Version string `json:"version"`
}

func (r *VersionResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[VERSION INFO]\n")
	buffer.WriteString(helper.FormatKV([]string{
		fmt.Sprintf("Release version|%s", r.Version),
	}))

	return buffer.String()
}
