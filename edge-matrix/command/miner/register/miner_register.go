package register

import (
	"fmt"
	"github.com/emc-protocol/edge-matrix/command"
	"github.com/emc-protocol/edge-matrix/command/helper"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	minerSnapshotCmd := &cobra.Command{
		Use:     "register",
		Short:   "register a Principal to be added or removed to the IC miner contract",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(minerSnapshotCmd)

	helper.SetRequiredFlags(minerSnapshotCmd, params.getRequiredFlags())

	return minerSnapshotCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.addressRaw,
		addressFlag,
		"",
		"the address of the miner to be register for",
	)

	cmd.Flags().StringVar(
		&params.commit,
		nodeFlag,
		"",
		fmt.Sprintf(
			"requested node type to the miner's address. Possible values: [%s, %s, %s]",
			validatorNodeOpt,
			routeNodeOpt,
			edgeNodeOpt,
		),
	)

	cmd.Flags().StringVar(
		&params.commit,
		commitFlag,
		"",
		fmt.Sprintf(
			"requested change to the miner's address. Possible values: [%s, %s]",
			setOpt,
			removeOpt,
		),
	)

	cmd.MarkFlagsRequiredTogether(addressFlag, commitFlag)
}

func runPreRun(_ *cobra.Command, _ []string) error {
	if err := params.validateFlags(); err != nil {
		return err
	}

	return params.initRawParams()
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	if err := params.registerMinerAddress(helper.GetGRPCAddress(cmd)); err != nil {
		outputter.SetError(err)

		return
	}

	outputter.SetCommandResult(params.getResult())
}
