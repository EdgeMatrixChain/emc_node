package status

import (
	"bytes"
	"fmt"

	"github.com/emc-protocol/edge-matrix/command/helper"
)

type MinerStatusResult struct {
	NetName   string `json:"net_name"`
	PeerID    string `json:"peer_id"`
	ICPubKey  string `json:"peer_ic_pubkey"`
	Principal string `json:"principal"`
}

func (r *MinerStatusResult) GetOutput() string {
	var buffer bytes.Buffer

	buffer.WriteString("\n[MINER STATUS]\n")
	buffer.WriteString(helper.FormatKV([]string{
		fmt.Sprintf("NetName |%s", r.NetName),
		fmt.Sprintf("PeerID |%s", r.PeerID),
		fmt.Sprintf("ICPubKey |%s", r.ICPubKey),
		fmt.Sprintf("Principal |%s", r.Principal),
	}))
	buffer.WriteString("\n")

	return buffer.String()
}
