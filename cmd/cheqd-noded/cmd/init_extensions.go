package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

func extendInit(initCmd *cobra.Command) *cobra.Command {
	baseRunE := initCmd.RunE

	initCmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := baseRunE(cmd, args)
		if err != nil {
			return err
		}

		err = applyConfigDefaults(cmd)
		if err != nil {
			return err
		}

		return nil
	}

	return initCmd
}

func applyConfigDefaults(cmd *cobra.Command) error {
	clientCtx := client.GetClientContextFromCmd(cmd)

	err := applyTmDefaults(clientCtx.HomeDir)
	if err != nil {
		return err
	}

	err = applyCosmosDefaults(clientCtx.HomeDir)
	if err != nil {
		return err
	}

	return nil
}

func applyTmDefaults(homeDir string) error {
	tmConfig, err := readTmConfig(homeDir)
	if err != nil {
		return err
	}

	tmConfig.Consensus.CreateEmptyBlocks = false

	err = tmConfig.ValidateBasic()
	if err != nil {
		return err
	}

	writeTmConfig(homeDir, &tmConfig)

	return nil
}

func applyCosmosDefaults(homeDir string) error {
	cosmConfig, err := readCosmosConfig(homeDir)
	if err != nil {
		return err
	}

	cosmConfig.BaseConfig.MinGasPrices = "25ncheq"

	err = cosmConfig.ValidateBasic()
	if err != nil {
		return err
	}

	writeCosmosConfig(homeDir, &cosmConfig)

	return nil
}