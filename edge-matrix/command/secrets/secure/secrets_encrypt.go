package secure

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/emc-protocol/edge-matrix/command"
)

const (
	// maxInitNum is the maximum value for "num" flag
	maxInitNum = 30
)

var (
	errInvalidNum = fmt.Errorf("num flag value should be between 1 and %d", maxInitNum)

	basicParams initParams
	initNumber  int
)

func GetCommand() *cobra.Command {
	secretsInitCmd := &cobra.Command{
		Use: "encrypt",
		Short: "Encrypt private keys for the Edge Matrix " +
			"to the encrypted Secrets Manager",
		PreRunE: runPreRun,
		Run:     runCommand,
	}

	setFlags(secretsInitCmd)

	return secretsInitCmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(
		&basicParams.dataDir,
		dataDirFlag,
		"",
		"the directory for the Edge Matrix data if the local FS is used",
	)

	cmd.Flags().BoolVar(
		&basicParams.generatesECDSA,
		ecdsaFlag,
		true,
		"the flag indicating whether new ECDSA key is created",
	)

	cmd.Flags().BoolVar(
		&basicParams.generatesNetwork,
		networkFlag,
		true,
		"the flag indicating whether new Network key is created",
	)

	cmd.Flags().BoolVar(
		&basicParams.generatesBLS,
		blsFlag,
		true,
		"the flag indicating whether new BLS key is created",
	)

	cmd.Flags().BoolVar(
		&basicParams.generatesICPIdentity,
		icpFlag,
		true,
		"the flag indicating whether new ICP identity key is created",
	)

	cmd.Flags().BoolVar(
		&basicParams.ensecureLocalStore,
		localStoreFlag,
		true,
		"the flag indicating should the secrets stored on the local storage be encrypted",
	)
	cmd.Flags().IntVar(
		&initNumber,
		numFlag,
		1,
		"the flag indicating how many secrets should be created, only for the local FS",
	)
}

func runPreRun(_ *cobra.Command, _ []string) error {
	if initNumber < 1 || initNumber > maxInitNum {
		return errInvalidNum
	}

	return basicParams.validateFlags()
}

func runCommand(cmd *cobra.Command, _ []string) {
	outputter := command.InitializeOutputter(cmd)
	defer outputter.WriteOutput()

	paramsList := getParamsList()
	results := make(command.Results, len(paramsList))

	secretsPass1 := ""
	fmt.Print("Input a new password: ")
	_, err := fmt.Scanln(&secretsPass1)
	if err != nil {
		outputter.SetError(err)
		return
	}

	secretsPass2 := ""
	fmt.Print("Re-type your password: ")
	_, err = fmt.Scanln(&secretsPass2)
	if err != nil {
		outputter.SetError(err)
		return
	}

	if secretsPass1 != secretsPass2 {
		outputter.SetError(errors.New("assword did not match!"))
		return
	}
	outputter.SetError(errors.New("encrypt complete!(this subcommand is a test command, it does not work yet!)"))
	return

	for i, params := range paramsList {
		if err := params.encryptSecrets(secretsPass2); err != nil {
			outputter.SetError(err)

			return
		}

		res, err := params.getResult()
		if err != nil {
			outputter.SetError(err)

			return
		}

		results[i] = res
	}

	outputter.SetCommandResult(results)
}

// getParamsList creates a list of initParams with num elements.
// This function basically copies the given initParams but updating dataDir by applying an index.
func getParamsList() []initParams {
	if initNumber == 1 {
		return []initParams{basicParams}
	}

	paramsList := make([]initParams, initNumber)
	for i := 1; i <= initNumber; i++ {
		paramsList[i-1] = initParams{
			dataDir:              fmt.Sprintf("%s%d", basicParams.dataDir, i),
			configPath:           basicParams.configPath,
			generatesECDSA:       basicParams.generatesECDSA,
			generatesBLS:         basicParams.generatesBLS,
			generatesNetwork:     basicParams.generatesNetwork,
			generatesICPIdentity: basicParams.generatesICPIdentity,
			ensecureLocalStore:   basicParams.ensecureLocalStore,
		}
	}

	return paramsList
}
