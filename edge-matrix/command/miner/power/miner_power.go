package power

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
		Use:   "power",
		Short: "Returns the miner's e-power be generated today",
		Run:   runCommand,
	}
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	response, err := getCurrentEPower(helper.GetGRPCAddress(cmd))
	if err != nil {
		outputter.SetError(err)

		return
	}

	outputter.SetCommandResult(&CurrentEPowerResult{
		Round:    response.Round,
		Total:    response.Total,
		Multiple: response.Multiple,
	})
}

func getCurrentEPower(grpcAddress string) (*minerOp.CurrentEPower, error) {
	client, err := helper.GetMinerClientConnection(
		grpcAddress,
	)
	if err != nil {
		return nil, err
	}

	return client.GetCurrentEPower(context.Background(), &empty.Empty{})
}
