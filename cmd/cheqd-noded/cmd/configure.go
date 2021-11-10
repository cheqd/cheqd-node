package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
	"strconv"
)

const (
	flagSeedMode = "seed-mode"
	flagSeeds = "seeds"
	flagExternalAddress = "external-address"
	flagPersistentPeers = "persistent-peers"
	flagMinGasPrices    = "gas-prices"
)

// configureCmd returns configure cobra Command.
func configureCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Review and adjust most important node parameters",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			// Read config
			tmConfig, err := readTmConfig(clientCtx.HomeDir)
			if err != nil {
				return err
			}

			cosmosConfig, err := readCosmosConfig(clientCtx.HomeDir)
			if err != nil {
				return err
			}

			// Update seed mode
			seedModeStr, err := cmd.Flags().GetString(flagSeedMode)
			if err != nil {
				return err
			}

			if seedModeStr != "" {
				seedModeBool, err := strconv.ParseBool(seedModeStr)
				if err != nil {
					return errors.Wrap(err, "can't parse seed mode")
				}

				fmt.Println("Updating seed mode, new value: ", seedModeBool)
				tmConfig.P2P.SeedMode = seedModeBool
			}

			// Update seeds
			seeds, err := cmd.Flags().GetString(flagSeeds)
			if err != nil {
				return err
			}

			if seeds != "" {
				fmt.Println("Updating seeds, new value: ", seeds)
				tmConfig.P2P.Seeds = seeds
			}

			// Update external address
			externalAddress, err := cmd.Flags().GetString(flagExternalAddress)
			if err != nil {
				return err
			}

			if externalAddress != "" {
				fmt.Println("Updating external address, new value: ", externalAddress)
				tmConfig.P2P.ExternalAddress = externalAddress
			}

			// Update persistent peers
			persistentPeers, err := cmd.Flags().GetString(flagPersistentPeers)
			if err != nil {
				return err
			}

			if persistentPeers != "" {
				fmt.Println("Updating persistent peers, new value: ", persistentPeers)
				tmConfig.P2P.PersistentPeers = persistentPeers
			}

			// Update gas prices
			minGasPrices, err := cmd.Flags().GetString(flagMinGasPrices)
			if err != nil {
				return err
			}

			if minGasPrices != "" {
				fmt.Println("Updating gas prices, new value: ", minGasPrices)
				cosmosConfig.MinGasPrices = minGasPrices
			}

			// Validate config
			err = tmConfig.ValidateBasic()
			if err != nil {
				return err
			}

			err = cosmosConfig.ValidateBasic()
			if err != nil {
				return err
			}

			writeTmConfig(clientCtx.HomeDir, &tmConfig)
			writeCosmosConfig(clientCtx.HomeDir, &cosmosConfig)

			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	cmd.Flags().String(flagSeedMode, "",  "seed mode (true|false), updated only if set")
	cmd.Flags().String(flagSeeds, "",  "seeds, updated only if set")
	cmd.Flags().String(flagExternalAddress, "",  "external address, updated only if set")
	cmd.Flags().String(flagPersistentPeers, "",  "persistent peers, updated only if set")
	cmd.Flags().String(flagMinGasPrices, "",  "ming gas prices, updated only if set")

	return cmd
}
