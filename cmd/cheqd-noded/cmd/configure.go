package cmd

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	cosmcfg "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
)

// configureCmd returns configure cobra Command.
func configureCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Adjust node parameters",
	}

	cmd.AddCommand(
		minGasPricesCmd(defaultNodeHome),
		p2pCmd(defaultNodeHome),
		rpcLaddrCmd(defaultNodeHome),
		createEmptyBlocksCmd(defaultNodeHome),
		fastsyncVersionCmd(defaultNodeHome))

	return cmd
}

// p2pCmd returns configure cobra Command.
func p2pCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "p2p",
		Short: "Adjust p2p parameters",
	}

	cmd.AddCommand(
		seedModeCmd(defaultNodeHome),
		seedsCmd(defaultNodeHome),
		externalAddressCmd(defaultNodeHome),
		persistentPeersCmd(defaultNodeHome),
		sendRateCmd(defaultNodeHome),
		recvRateCmd(defaultNodeHome),
		maxPacketMsgPayloadSizeCmd(defaultNodeHome))

	return cmd
}

// minGasPricesCmd returns configuration cobra Command.
func minGasPricesCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "min-gas-prices [value]",
		Short: "Update min-gas-prices value in app.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			return updateCosmConfig(clientCtx.HomeDir, func(config *cosmcfg.Config) {
				config.MinGasPrices = args[0]
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

// seedModeCmd returns configuration cobra Command.
func seedModeCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed-mode [value]",
		Short: "Update seed-mode value in config.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			value, err := strconv.ParseBool(args[0])
			if err != nil {
				return errors.Wrap(err, "can't parse seed mode")
			}

			return updateTmConfig(clientCtx.HomeDir, func(config *tmcfg.Config) {
				config.P2P.SeedMode = value
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

// seedsCmd returns configuration cobra Command.
func seedsCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seeds [value]",
		Short: "Update seeds value in config.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			return updateTmConfig(clientCtx.HomeDir, func(config *tmcfg.Config) {
				config.P2P.Seeds = args[0]
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

// externalAddressCmd returns configuration cobra Command.
func externalAddressCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "external-address [value]",
		Short: "Update external-address value in config.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			return updateTmConfig(clientCtx.HomeDir, func(config *tmcfg.Config) {
				config.P2P.ExternalAddress = args[0]
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

// persistentPeersCmd returns configuration cobra Command.
func persistentPeersCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "persistent-peers [value]",
		Short: "Update persistent-peers value in config.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			return updateTmConfig(clientCtx.HomeDir, func(config *tmcfg.Config) {
				config.P2P.PersistentPeers = args[0]
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

// sendRateCmd returns configuration cobra Command.
func sendRateCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-rate [value]",
		Short: "Update send-rate value in config.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			value, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return errors.Wrap(err, "can't parse send rate")
			}

			return updateTmConfig(clientCtx.HomeDir, func(config *tmcfg.Config) {
				config.P2P.SendRate = value
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

// recvRateCmd returns configuration cobra Command.
func recvRateCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recv-rate [value]",
		Short: "Update recv-rate value in config.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			value, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return errors.Wrap(err, "can't parse recv rate")
			}

			return updateTmConfig(clientCtx.HomeDir, func(config *tmcfg.Config) {
				config.P2P.RecvRate = value
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

// maxPacketMsgPayloadSizeCmd returns configuration cobra Command.
func maxPacketMsgPayloadSizeCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "max-packet-msg-payload-size [value]",
		Short: "Update max-packet-msg-payload-size value in config.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			value, err := strconv.ParseInt(args[0], 10, 32)
			if err != nil {
				return errors.Wrap(err, "can't parse max-packet-msg-payload-size")
			}

			return updateTmConfig(clientCtx.HomeDir, func(config *tmcfg.Config) {
				config.P2P.MaxPacketMsgPayloadSize = int(value)
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

// createEmptyBlocksCmd returns configuration cobra Command.
func createEmptyBlocksCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-empty-blocks [value]",
		Short: "Update create-empty-blocks value in config.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			value, err := strconv.ParseBool(args[0])
			if err != nil {
				return errors.Wrap(err, "can't parse create-empty-blocks")
			}

			return updateTmConfig(clientCtx.HomeDir, func(config *tmcfg.Config) {
				config.Consensus.CreateEmptyBlocks = value
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

// createEmptyBlocksCmd returns configuration cobra Command.
func rpcLaddrCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rpc-laddr [value]",
		Short: "Update rpc.laddr value in config.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			return updateTmConfig(clientCtx.HomeDir, func(config *tmcfg.Config) {
				config.RPC.ListenAddress = args[0]
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}

// fastsyncVersionCmd returns configuration cobra Command.
func fastsyncVersionCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fastsync-version [value]",
		Short: "Update fastsync.version value in config.toml",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			return updateTmConfig(clientCtx.HomeDir, func(config *tmcfg.Config) {
				config.FastSync.Version = args[0]
			})
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")

	return cmd
}
