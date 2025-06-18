package oracle

// import (
// 	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
// 	oraclev2 "github.com/cheqd/cheqd-node/api/v2/cheqd/oracle/v2"
// )

// func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
// 	return &autocliv1.ModuleOptions{
// 		Query: &autocliv1.ServiceCommandDescriptor{
// 			Service: oraclev2.Query_ServiceDesc.ServiceName,
// 			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
// 				{
// 					RpcMethod: "ExchangeRates",
// 					Use:       "exchange-rates [denom]",
// 					Short:     "Query the exchange rate for a specific denom",
// 					Long:      "Query the current exchange rate of a specific asset based on USD.",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
// 						{ProtoField: "denom", Optional: true},
// 					},
// 				},
// 				{
// 					RpcMethod: "ActiveExchangeRates",
// 					Use:       "active-exchange-rates",
// 					Short:     "Query all active exchange rate denoms",
// 				},
// 				{
// 					RpcMethod: "FeederDelegation",
// 					Use:       "feeder-delegation [validator_addr]",
// 					Short:     "Query the feeder delegation for a validator",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
// 						{ProtoField: "validator_addr"},
// 					},
// 				},
// 				{
// 					RpcMethod: "MissCounter",
// 					Use:       "miss-counter [validator_addr]",
// 					Short:     "Query the oracle miss counter for a validator",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
// 						{ProtoField: "validator_addr"},
// 					},
// 				},
// 				{
// 					RpcMethod: "SlashWindow",
// 					Use:       "slash-window",
// 					Short:     "Query the oracle slash window parameters",
// 				},
// 				{
// 					RpcMethod: "AggregatePrevote",
// 					Use:       "aggregate-prevote [validator_addr]",
// 					Short:     "Query the aggregate prevote for a validator",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
// 						{ProtoField: "validator_addr"},
// 					},
// 				},
// 				{
// 					RpcMethod: "AggregatePrevotes",
// 					Use:       "aggregate-prevotes",
// 					Short:     "Query all aggregate prevotes",
// 				},
// 				{
// 					RpcMethod: "AggregateVote",
// 					Use:       "aggregate-vote [validator_addr]",
// 					Short:     "Query the aggregate vote for a validator",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
// 						{ProtoField: "validator_addr"},
// 					},
// 				},
// 				{
// 					RpcMethod: "AggregateVotes",
// 					Use:       "aggregate-votes",
// 					Short:     "Query all aggregate votes",
// 				},
// 				{
// 					RpcMethod: "Params",
// 					Use:       "params",
// 					Short:     "Query the oracle module parameters",
// 				},
// 				{
// 					RpcMethod: "Medians",
// 					Use:       "medians [denom]",
// 					Short:     "Query median exchange rate(s)",
// 					Long:      "Query median exchange rate for all denoms or a specific one",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
// 						{ProtoField: "denom"}, // assuming denom is optional in the proto
// 					},
// 				},
// 				{
// 					RpcMethod: "MedianDeviations",
// 					Use:       "median-deviations [denom]",
// 					Short:     "Query median deviation(s)",
// 					Long:      "Query median deviation for all denoms or a specific one",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
// 						{ProtoField: "denom"}, // assuming denom is optional in the proto
// 					},
// 				},
// 				{
// 					RpcMethod: "ValidatorRewardSet",
// 					Use:       "validator-reward-set",
// 					Short:     "Query the list of validators eligible for oracle rewards",
// 				},
// 				{
// 					RpcMethod: "EMA",
// 					Use:       "ema [denom]",
// 					Short:     "Query the EMA price for a denom",
// 					Long:      "Returns the Exponential Moving Average (EMA) price for a given denom and strategy (e.g., FAST, SLOW, BALANCED).",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
// 						{ProtoField: "denom"},
// 					},
// 				},
// 				{
// 					RpcMethod: "WMA",
// 					Use:       "wma [denom] [strategy]",
// 					Short:     "Query the WMA price for a denom using a strategy",
// 					Long:      "Returns the Weighted Moving Average (WMA) price for a given denom and strategy (e.g., BALANCED, RECENT, CUSTOM).",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
// 						{ProtoField: "denom"},
// 						{ProtoField: "strategy"},
// 					},
// 					FlagOptions: map[string]*autocliv1.FlagOptions{
// 						"custom_weights": {
// 							Name: "custom_weights",
// 						},
// 					},
// 				},
// 				{
// 					RpcMethod: "SMA",
// 					Use:       "sma [denom] [strategy]",
// 					Short:     "Query the SMA price for a denom ",
// 					Long:      "Returns the Simple Moving Average (SMA) price for a given denom and strategy (e.g., FAST, SLOW, BALANCED).",
// 					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
// 						{ProtoField: "denom"},
// 					},
// 				},
// 			},
// 		},
// 	}
// }
