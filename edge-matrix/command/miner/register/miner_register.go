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
		Short:   "register node owner (ethereum address) to be added or removed to the EMC Hub",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(minerSnapshotCmd)

	helper.SetRequiredFlags(minerSnapshotCmd, params.getRequiredFlags())

	return minerSnapshotCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&params.principal,
		principalFlag,
		"",
		"the owner (ethereum address) of the node to be register for",
	)

	cmd.Flags().StringVar(
		&params.nodeType,
		nodeFlag,
		"",
		fmt.Sprintf(
			"requested node type to the node. Possible values: [%s, %s, %s]",
			validatorNodeOpt,
			routeNodeOpt,
			computingNodeOpt,
		),
	)

	cmd.Flags().StringVar(
		&params.commit,
		commitFlag,
		"",
		fmt.Sprintf(
			"requested change to the node's owner. Possible values: [%s, %s]",
			setOpt,
			removeOpt,
		),
	)

	//cmd.MarkFlagsRequiredTogether(principalFlag, commitFlag)
}

func runPreRun(_ *cobra.Command, _ []string) error {
	if err := params.validateFlags(); err != nil {
		return err
	}

	return nil
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
