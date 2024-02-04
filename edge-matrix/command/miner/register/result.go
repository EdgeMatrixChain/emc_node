package register

import (
	"bytes"
	"fmt"
)

type MinerRegisterResult struct {
	Principal    string `json:"-"`
	Commit       string `json:"-"`
	NodeType     string `json:"-"`
	ResultMessge string `json:"-"`
}

func (r *MinerRegisterResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[NODE REGISTER]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *MinerRegisterResult) Message() string {
	if r.Commit == setOpt {
		return fmt.Sprintf(
			"Commit for the add/update node [%s] owner [%s] to the node pool\n%s \n",
			r.NodeType,
			r.Principal,
			r.ResultMessge,
		)
	}

	return fmt.Sprintf(
		"Commit for the removal node from the node poll\n%s \n",
		r.ResultMessge,
	)
}

func (r *MinerRegisterResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
