package relaylist

import (
	"context"

	"github.com/emc-protocol/edge-matrix/command"
	"github.com/emc-protocol/edge-matrix/command/helper"
	"github.com/emc-protocol/edge-matrix/server/proto"
	"github.com/spf13/cobra"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func GetCommand() *cobra.Command {
	peersListCmd := &cobra.Command{
		Use:   "relaylist",
		Short: "Returns the list of relay nodes, including the current connected node",
		Run:   runCommand,
	}

	return peersListCmd
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	peersList, err := getPeersRelayList(helper.GetGRPCAddress(cmd))
	if err != nil {
		outputter.SetError(err)

		return
	}

	outputter.SetCommandResult(
		newPeersListResult(peersList.Peers),
	)
}

func getPeersRelayList(grpcAddress string) (*proto.PeersListResponse, error) {
	client, err := helper.GetSystemClientConnection(grpcAddress)
	if err != nil {
		return nil, err
	}

	return client.PeersRelayList(context.Background(), &empty.Empty{})
}
