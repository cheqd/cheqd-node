package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/cheqd/cheqd-node/util/cli"
	"github.com/cheqd/cheqd-node/x/oracle/types"
)

// GetQueryCmd returns the CLI query commands for the x/oracle module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdQueryAggregatePrevote(),
		GetCmdQueryAggregateVote(),
		GetCmdQueryParams(),
		GetCmdQueryExchangeRates(),
		GetCmdQueryExchangeRate(),
		GetCmdQueryFeederDelegation(),
		GetCmdQueryMissCounter(),
		GetCmdQuerySlashWindow(),
		CmdQueryEMA(),
		CmdQueryWMA(),
		CmdQuerySMA(),
	)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current Oracle params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParams{})
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryAggregateVote implements the query aggregate prevote of the
// validator command.
func GetCmdQueryAggregateVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-votes [validator]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query outstanding oracle aggregate votes",
		Long: strings.TrimSpace(`
Query outstanding oracle aggregate vote.

$ cheqd-noded query oracle aggregate-votes

Or, you can filter with voter address

$ cheqd-noded query oracle aggregate-votes cheqdvaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			query := types.QueryAggregateVote{}

			if len(args) > 0 {
				validator, err := sdk.ValAddressFromBech32(args[0])
				if err != nil {
					return err
				}
				query.ValidatorAddr = validator.String()
			}

			res, err := queryClient.AggregateVote(cmd.Context(), &query)
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryAggregatePrevote implements the query aggregate prevote of the
// validator command.
func GetCmdQueryAggregatePrevote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-prevotes [validator]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query outstanding oracle aggregate prevotes",
		Long: strings.TrimSpace(`
Query outstanding oracle aggregate prevotes.

$ cheqd-noded query oracle aggregate-prevotes

Or, can filter with voter address

$ cheqd-noded query oracle aggregate-prevotes cheqdvaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			query := types.QueryAggregatePrevote{}

			if len(args) > 0 {
				validator, err := sdk.ValAddressFromBech32(args[0])
				if err != nil {
					return err
				}
				query.ValidatorAddr = validator.String()
			}

			res, err := queryClient.AggregatePrevote(cmd.Context(), &query)
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryExchangeRates implements the query rate command.
func GetCmdQueryExchangeRates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exchange-rates",
		Args:  cobra.NoArgs,
		Short: "Query the exchange rates",
		Long: strings.TrimSpace(`
Query the current exchange rates of assets based on USD.
You can find the current list of active denoms by running

$ cheqd-noded query oracle exchange-rates
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ExchangeRates(cmd.Context(), &types.QueryExchangeRates{})
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryExchangeRates implements the query rate command.
func GetCmdQueryExchangeRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exchange-rate [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the exchange rates",
		Long: strings.TrimSpace(`
Query the current exchange rates of an asset based on USD.

$ cheqd-noded query oracle exchange-rate ATOM
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ExchangeRates(
				cmd.Context(),
				&types.QueryExchangeRates{
					Denom: args[0],
				},
			)
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryFeederDelegation implements the query feeder delegation command.
func GetCmdQueryFeederDelegation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feeder-delegation [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the current delegate for a given validator address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if _, err = sdk.ValAddressFromBech32(args[0]); err != nil {
				return err
			}
			res, err := queryClient.FeederDelegation(cmd.Context(), &types.QueryFeederDelegation{
				ValidatorAddr: args[0],
			})
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryMissCounter implements the miss counter query command.
func GetCmdQueryMissCounter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miss-counter [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the current miss counter for a given validator address",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			if _, err = sdk.ValAddressFromBech32(args[0]); err != nil {
				return err
			}
			res, err := queryClient.MissCounter(cmd.Context(), &types.QueryMissCounter{
				ValidatorAddr: args[0],
			})
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQuerySlashWindow implements the slash window query command.
func GetCmdQuerySlashWindow() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slash-window",
		Short: "Query the current slash window progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.SlashWindow(cmd.Context(), &types.QuerySlashWindow{})
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryEMA() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ema [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the ema of the given denom",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			denom := strings.ToUpper(args[0]) // Convert denom to uppercase
			req := &types.QueryEMARequest{Denom: denom}
			res, err := queryClient.EMA(cmd.Context(), req)
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdQueryWMA() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wma [denom]",
		Short: "Query WMA price for a denom with specified strategy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			denom := args[0]
			strategy, _ := cmd.Flags().GetString("strategy")
			weightsStr, _ := cmd.Flags().GetString("weights")

			var weights []int64
			if weightsStr != "" {
				for _, w := range strings.Split(weightsStr, ",") {
					val, err := strconv.ParseInt(strings.TrimSpace(w), 10, 64)
					if err != nil {
						return err
					}
					weights = append(weights, val)
				}
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.WMA(cmd.Context(), &types.QueryWMARequest{
				Denom:         denom,
				Strategy:      strategy,
				CustomWeights: weights,
			})
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	cmd.Flags().String("strategy", "BALANCED", "WMA strategy: BALANCED | OLDEST | RECENT | CUSTOM")
	cmd.Flags().String("weights", "", "Custom weights (comma-separated, e.g. 10,9,8...)")

	return cmd
}

func CmdQuerySMA() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sma [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the SMA of the given denom",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			denom := strings.ToUpper(args[0]) // Convert denom to uppercase
			req := &types.QuerySMARequest{Denom: denom}
			res, err := queryClient.SMA(cmd.Context(), req)
			return cli.PrintOrErr(res, err, clientCtx)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
