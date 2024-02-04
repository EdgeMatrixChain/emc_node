package status

import (
	"context"

	"github.com/emc-protocol/edge-matrix/command"
	"github.com/emc-protocol/edge-matrix/command/helper"
	minerOp "github.com/emc-protocol/edge-matrix/miner/proto"
	"github.com/spf13/cobra"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Returns the current node status",
		Run:   runCommand,
	}
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	statusResponse, err := getMinerStatus(helper.GetGRPCAddress(cmd))
	if err != nil {
		outputter.SetError(err)

		return
	}

	outputter.SetCommandResult(&MinerStatusResult{
		NetName:      statusResponse.NetName,
		Principal:    statusResponse.Principal,
		NodeIdentity: statusResponse.NodeIdentity,
		NodeID:       statusResponse.NodeId,
		NodeType:     statusResponse.NodeType,
		Registered:   statusResponse.Registered == 1,
	})
}

func getMinerStatus(grpcAddress string) (*minerOp.MinerStatus, error) {
	client, err := helper.GetMinerClientConnection(
		grpcAddress,
	)
	if err != nil {
		return nil, err
	}

	return client.GetMinerStatus(context.Background(), &empty.Empty{})
}
