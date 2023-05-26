package init

import (
	"fmt"
	"github.com/spf13/cobra"
)

var basicParams initParams

func GetCommand() *cobra.Command {
	secretsInitCmd := &cobra.Command{
		Use: "init",
		Short: "Initializes private keys for the Edge Matrix (Validator + Networking) " +
			"to the specified Secrets Manager",
		Run: runCommand,
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

	cmd.Flags().StringVar(
		&basicParams.configPath,
		configFlag,
		"",
		"the path to the SecretsManager config file, "+
			"if omitted, the local FS secrets manager is used",
	)

	// Don't accept data-dir and config flags because they are related to different secrets managers.
	// data-dir is about the local FS as secrets storage, config is about remote secrets manager.
	cmd.MarkFlagsMutuallyExclusive(dataDirFlag, configFlag)

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
		&basicParams.insecureLocalStore,
		localStoreFlag,
		true,
		"the flag indicating should the secrets stored on the local storage be encrypted",
	)
}

func runCommand(cmd *cobra.Command, _ []string) {
	if err := initSecrets(); err != nil {
		fmt.Errorf("init error: %s", err)
		return
	}
}

func initSecrets() error {
	if err := initSecretsManager(); err != nil {
		return err
	}

	return basicParams.initNetworkingKey()
}

func initSecretsManager() error {
	return basicParams.initLocalSecretsManager()
}
