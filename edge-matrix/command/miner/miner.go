package miner

import (
	"github.com/emc-protocol/edge-matrix/command/helper"
	"github.com/emc-protocol/edge-matrix/command/miner/register"
	"github.com/emc-protocol/edge-matrix/command/miner/status"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	minerCmd := &cobra.Command{
		Use:   "miner",
		Short: "Top level Miner command for interacting with the Miner contracts. Only accepts subcommands.",
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
	)
}
