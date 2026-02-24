package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cheqd/cheqd-node/util"
	"github.com/cheqd/cheqd-node/x/oracle/types"
)

// Simulation parameter constants
const (
	votePeriodKey               = "vote_period"
	voteThresholdKey            = "vote_threshold"
	rewardBandsKey              = "reward_bands"
	rewardDistributionWindowKey = "reward_distribution_window"
	slashFractionKey            = "slash_fraction"
	slashWindowKey              = "slash_window"
	minValidPerWindowKey        = "min_valid_per_window"
	historicStampPeriodKey      = "historic_stamp_period"
	medianStampPeriodKey        = "median_stamp_period"
	maximumPriceStampsKey       = "maximum_price_stamps"
	maximumMedianStampsKey      = "maximum_median_stamps"
)

// GenVotePeriod produces a randomized VotePeriod in the range of [5, 100]
func GenVotePeriod(r *rand.Rand) uint64 {
	return util.SafeIntToUint64(5 + r.Intn(100))
}

// GenVoteThreshold produces a randomized VoteThreshold in the range of [0.34, 0.67]
func GenVoteThreshold(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(34, 2).Add(math.LegacyNewDecWithPrec(int64(r.Intn(33)), 2))
}

// GenRewardBand produces a randomized RewardBand in the range of [0.000, 0.100]
func GenRewardBand(r *rand.Rand) math.LegacyDec {
	return math.LegacyZeroDec().Add(math.LegacyNewDecWithPrec(int64(r.Intn(100)), 3))
}

// GenRewardDistributionWindow produces a randomized RewardDistributionWindow in the range of [100, 100000]
func GenRewardDistributionWindow(r *rand.Rand) uint64 {
	return util.SafeIntToUint64(100 + r.Intn(100000))
}

// GenSlashFraction produces a randomized SlashFraction in the range of [0.000, 0.100]
func GenSlashFraction(r *rand.Rand) math.LegacyDec {
	return math.LegacyZeroDec().Add(math.LegacyNewDecWithPrec(int64(r.Intn(100)), 3))
}

// GenSlashWindow produces a randomized SlashWindow in the range of [100, 100000]
func GenSlashWindow(r *rand.Rand) uint64 {
	return util.SafeIntToUint64(100 + r.Intn(100000))
}

// GenMinValidPerWindow produces a randomized MinValidPerWindow in the range of [0, 0.500]
func GenMinValidPerWindow(r *rand.Rand) math.LegacyDec {
	return math.LegacyZeroDec().Add(math.LegacyNewDecWithPrec(int64(r.Intn(500)), 3))
}

// GenHistoricStampPeriod produces a randomized HistoricStampPeriod in the range of [100, 1000]
func GenHistoricStampPeriod(r *rand.Rand) uint64 {
	return util.SafeIntToUint64(100 + r.Intn(1000))
}

// GenMedianStampPeriod produces a randomized MedianStampPeriod in the range of [100, 1000]
func GenMedianStampPeriod(r *rand.Rand) uint64 {
	return util.SafeIntToUint64(10001 + r.Intn(100000))
}

// GenMaximumPriceStamps produces a randomized MaximumPriceStamps in the range of [10, 100]
func GenMaximumPriceStamps(r *rand.Rand) uint64 {
	return util.SafeIntToUint64(11 + r.Intn(100))
}

// GenMaximumMedianStamps produces a randomized MaximumMedianStamps in the range of [10, 100]
func GenMaximumMedianStamps(r *rand.Rand) uint64 {
	return util.SafeIntToUint64(11 + r.Intn(100))
}

// RandomizedGenState generates a random GenesisState for oracle
func RandomizedGenState(simState *module.SimulationState) {
	oracleGenesis := types.DefaultGenesisState()

	var votePeriod uint64
	simState.AppParams.GetOrGenerate(
		votePeriodKey, &votePeriod, simState.Rand,
		func(r *rand.Rand) { votePeriod = GenVotePeriod(r) },
	)

	var voteThreshold math.LegacyDec
	simState.AppParams.GetOrGenerate(
		voteThresholdKey, &voteThreshold, simState.Rand,
		func(r *rand.Rand) { voteThreshold = GenVoteThreshold(r) },
	)

	var rewardBands types.RewardBandList
	simState.AppParams.GetOrGenerate(
		rewardBandsKey, &rewardBands, simState.Rand,
		func(r *rand.Rand) {
			for _, denom := range oracleGenesis.Params.MandatoryList {
				rb := types.RewardBand{
					RewardBand:  GenRewardBand(r),
					SymbolDenom: denom.SymbolDenom,
				}
				rewardBands = append(rewardBands, rb)
			}
			for _, denom := range oracleGenesis.Params.AcceptList {
				rb := types.RewardBand{
					RewardBand:  GenRewardBand(r),
					SymbolDenom: denom.SymbolDenom,
				}
				rewardBands = append(rewardBands, rb)
			}
		},
	)

	var rewardDistributionWindow uint64
	simState.AppParams.GetOrGenerate(
		rewardDistributionWindowKey, &rewardDistributionWindow, simState.Rand,
		func(r *rand.Rand) { rewardDistributionWindow = GenRewardDistributionWindow(r) },
	)

	var slashFraction math.LegacyDec
	simState.AppParams.GetOrGenerate(
		slashFractionKey, &slashFraction, simState.Rand,
		func(r *rand.Rand) { slashFraction = GenSlashFraction(r) },
	)

	var slashWindow uint64
	simState.AppParams.GetOrGenerate(
		slashWindowKey, &slashWindow, simState.Rand,
		func(r *rand.Rand) { slashWindow = GenSlashWindow(r) },
	)

	var minValidPerWindow math.LegacyDec
	simState.AppParams.GetOrGenerate(
		minValidPerWindowKey, &minValidPerWindow, simState.Rand,
		func(r *rand.Rand) { minValidPerWindow = GenMinValidPerWindow(r) },
	)

	var historicStampPeriod uint64
	simState.AppParams.GetOrGenerate(
		historicStampPeriodKey, &historicStampPeriod, simState.Rand,
		func(r *rand.Rand) { historicStampPeriod = GenHistoricStampPeriod(r) },
	)

	var medianStampPeriod uint64
	simState.AppParams.GetOrGenerate(
		medianStampPeriodKey, &medianStampPeriod, simState.Rand,
		func(r *rand.Rand) { medianStampPeriod = GenMedianStampPeriod(r) },
	)

	var maximumPriceStamps uint64
	simState.AppParams.GetOrGenerate(
		maximumPriceStampsKey, &maximumPriceStamps, simState.Rand,
		func(r *rand.Rand) { maximumPriceStamps = GenMaximumPriceStamps(r) },
	)

	var maximumMedianStamps uint64
	simState.AppParams.GetOrGenerate(
		maximumMedianStampsKey, &maximumMedianStamps, simState.Rand,
		func(r *rand.Rand) { maximumMedianStamps = GenMaximumMedianStamps(r) },
	)

	oracleGenesis.Params = types.Params{
		VotePeriod:               votePeriod,
		VoteThreshold:            voteThreshold,
		RewardBands:              rewardBands,
		RewardDistributionWindow: rewardDistributionWindow,
		AcceptList: types.DenomList{
			{SymbolDenom: types.CheqdSymbol, BaseDenom: types.CheqdDenom},
		},
		SlashFraction:       slashFraction,
		SlashWindow:         slashWindow,
		MinValidPerWindow:   minValidPerWindow,
		HistoricStampPeriod: historicStampPeriod,
		MedianStampPeriod:   medianStampPeriod,
		MaximumPriceStamps:  historicStampPeriod,
		MaximumMedianStamps: historicStampPeriod,
	}

	bz, err := json.MarshalIndent(&oracleGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated oracle parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(oracleGenesis)
}
