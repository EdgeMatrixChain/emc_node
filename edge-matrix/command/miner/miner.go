package miner

import (
	"github.com/emc-protocol/edge-matrix/command/helper"
	"github.com/emc-protocol/edge-matrix/command/miner/power"
	"github.com/emc-protocol/edge-matrix/command/miner/register"
	"github.com/emc-protocol/edge-matrix/command/miner/status"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	minerCmd := &cobra.Command{
		Use:   "node",
		Short: "Top level Node command for interacting with the emc. Only accepts subcommands.",
	}

	helper.RegisterGRPCAddressFlag(minerCmd)

	registerSubcommands(minerCmd)

	return minerCmd
}

func registerSubcommands(baseCmd *cobra.Command) {
	baseCmd.AddCommand(
		// miner status
		status.GetCommand(),
		// miner register
		register.GetCommand(),
		// miner power
		power.GetCommand(),
	)
}
