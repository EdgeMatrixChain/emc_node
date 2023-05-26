package register

import (
	"bytes"
	"fmt"
)

type MinerRegisterResult struct {
	Address      string `json:"-"`
	Commit       string `json:"-"`
	NodeType     string `json:"-"`
	ResultMessge string `json:"-"`
}

func (r *MinerRegisterResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[MINER REGISTER]\n")
	buffer.WriteString(r.Message())
	buffer.WriteString("\n")

	return buffer.String()
}

func (r *MinerRegisterResult) Message() string {
	if r.Commit == setOpt {
		return fmt.Sprintf(
			"Commit for the add/update the address [%s] to the miner set, %s",
			r.NodeType,
			r.Address,
			r.ResultMessge,
		)
	}

	return fmt.Sprintf(
		"Commit for the removal of miner at address [%s] from the miner set, %s",
		r.Address,
		r.ResultMessge,
	)
}

func (r *MinerRegisterResult) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"message": "%s"}`, r.Message())), nil
}
