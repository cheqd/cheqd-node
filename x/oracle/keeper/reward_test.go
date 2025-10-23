package keeper_test

import (
	"fmt"
	"math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cheqd/cheqd-node/x/oracle/types"
)

// Test the reward giving mechanism
func (s *IntegrationTestSuite) TestRewardBallotWinners() {
	app, ctx := s.app, s.ctx

	// Add claim pools
	claims := []types.Claim{
		types.NewClaim(10, 0, 0, valAddr),
		types.NewClaim(20, 0, 0, valAddr2),
	}

	missCounters := []types.MissCounter{
		{ValidatorAddress: valAddr.String(), MissCounter: uint64(2)},
		{ValidatorAddress: valAddr2.String(), MissCounter: uint64(4)},
	}

	for _, mc := range missCounters {
		operator, _ := sdk.ValAddressFromBech32(mc.ValidatorAddress)
		app.OracleKeeper.SetMissCounter(ctx, operator, mc.MissCounter)
	}

	// Prepare reward pool
	givingAmt := sdk.NewCoins(sdk.NewInt64Coin(types.CheqdDenom, 30000000))
	err := app.BankKeeper.MintCoins(ctx, "oracle", givingAmt)
	s.Require().NoError(err)

	var voteTargets []string
	params := app.OracleKeeper.GetParams(ctx)
	for _, v := range params.AcceptList {
		voteTargets = append(voteTargets, v.SymbolDenom)
	}

	// Add extra voteTargets to increase maximum miss count
	for i := 1; i <= 3; i++ {
		voteTargets = append(voteTargets, fmt.Sprintf("%s%d", types.CheqdSymbol, i))
	}
	maximumMissCounts := uint64(len(voteTargets)) * (app.OracleKeeper.SlashWindow(ctx) / app.OracleKeeper.VotePeriod(ctx))

	val1ExpectedRewardFactor := fmt.Sprintf("%f", 1-math.Log(1)/
		math.Log(float64(maximumMissCounts-missCounters[0].MissCounter+1)))

	val2ExpectedRewardFactor := fmt.Sprintf("%f", 1-math.Log(float64(missCounters[1].MissCounter-missCounters[0].MissCounter+1))/
		math.Log(float64(maximumMissCounts-missCounters[0].MissCounter+1)))

	votePeriodsPerWindow := sdkmath.LegacyNewDec((int64)(app.OracleKeeper.RewardDistributionWindow(ctx))).
		QuoInt64((int64)(app.OracleKeeper.VotePeriod(ctx))).
		TruncateInt64()
	app.OracleKeeper.RewardBallotWinners(ctx, (int64)(app.OracleKeeper.VotePeriod(ctx)), (int64)(app.OracleKeeper.RewardDistributionWindow(ctx)), voteTargets, claims)
	outstandingRewardsDecVal1, err := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr)
	s.Require().NoError(err)
	outstandingRewardsVal1, _ := outstandingRewardsDecVal1.TruncateDecimal()
	outstandingRewardsDecVal2, err := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr2)
	s.Require().NoError(err)
	outstandingRewardsVal2, _ := outstandingRewardsDecVal2.TruncateDecimal()
	s.Require().Equal(sdkmath.LegacyNewDecFromInt(givingAmt.AmountOf(types.CheqdDenom)).Mul(sdkmath.LegacyMustNewDecFromStr(val1ExpectedRewardFactor).QuoInt64(int64(len(claims)))).QuoInt64(votePeriodsPerWindow).TruncateInt(),
		outstandingRewardsVal1.AmountOf(types.CheqdDenom))
	s.Require().Equal(sdkmath.LegacyNewDecFromInt(givingAmt.AmountOf(types.CheqdDenom)).Mul(sdkmath.LegacyMustNewDecFromStr(val2ExpectedRewardFactor).QuoInt64(int64(len(claims)))).QuoInt64(votePeriodsPerWindow).TruncateInt(),
		outstandingRewardsVal2.AmountOf(types.CheqdDenom))
}

func (s *IntegrationTestSuite) TestRewardBallotWinnersZeroMissCounters() {
	app, ctx := s.app, s.ctx

	// Add claim pools
	claims := []types.Claim{
		types.NewClaim(10, 0, 0, valAddr),
		types.NewClaim(20, 0, 0, valAddr2),
	}

	// Prepare reward pool
	givingAmt := sdk.NewCoins(sdk.NewInt64Coin(types.CheqdDenom, 30000000))
	err := app.BankKeeper.MintCoins(ctx, "oracle", givingAmt)
	s.Require().NoError(err)

	var voteTargets []string
	params := app.OracleKeeper.GetParams(ctx)
	for _, v := range params.AcceptList {
		voteTargets = append(voteTargets, v.SymbolDenom)
	}

	votePeriodsPerWindow := sdkmath.LegacyNewDec((int64)(app.OracleKeeper.RewardDistributionWindow(ctx))).
		QuoInt64((int64)(app.OracleKeeper.VotePeriod(ctx))).
		TruncateInt64()
	app.OracleKeeper.RewardBallotWinners(ctx, (int64)(app.OracleKeeper.VotePeriod(ctx)), (int64)(app.OracleKeeper.RewardDistributionWindow(ctx)), voteTargets, claims)
	outstandingRewardsDecVal1, err := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr)
	s.Require().NoError(err)
	outstandingRewardsVal1, _ := outstandingRewardsDecVal1.TruncateDecimal()
	outstandingRewardsDecVal2, err := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr2)
	s.Require().NoError(err)
	outstandingRewardsVal2, _ := outstandingRewardsDecVal2.TruncateDecimal()
	s.Require().Equal(sdkmath.LegacyNewDecFromInt(givingAmt.AmountOf(types.CheqdDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(2).TruncateInt(),
		outstandingRewardsVal1.AmountOf(types.CheqdDenom))
	s.Require().Equal(sdkmath.LegacyNewDecFromInt(givingAmt.AmountOf(types.CheqdDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(2).TruncateInt(),
		outstandingRewardsVal2.AmountOf(types.CheqdDenom))
}

func (s *IntegrationTestSuite) TestRewardBallotWinnersZeroVoteTargets() {
	app, ctx := s.app, s.ctx

	// Add claim pools
	claims := []types.Claim{
		types.NewClaim(10, 0, 0, valAddr),
		types.NewClaim(20, 0, 0, valAddr2),
	}

	app.OracleKeeper.RewardBallotWinners(ctx, (int64)(app.OracleKeeper.VotePeriod(ctx)), (int64)(app.OracleKeeper.RewardDistributionWindow(ctx)), []string{}, claims)
	outstandingRewardsDecVal1, err := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr)
	s.Require().NoError(err)
	outstandingRewardsVal1, _ := outstandingRewardsDecVal1.TruncateDecimal()
	outstandingRewardsDecVal2, err := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr2)
	s.Require().NoError(err)
	outstandingRewardsVal2, _ := outstandingRewardsDecVal2.TruncateDecimal()
	s.Require().Equal(sdkmath.LegacyZeroDec().TruncateInt(), outstandingRewardsVal1.AmountOf(types.CheqdDenom))
	s.Require().Equal(sdkmath.LegacyZeroDec().TruncateInt(), outstandingRewardsVal2.AmountOf(types.CheqdDenom))
}

func (s *IntegrationTestSuite) TestRewardBallotWinnersZeroClaims() {
	app, ctx := s.app, s.ctx

	var voteTargets []string
	params := app.OracleKeeper.GetParams(ctx)
	for _, v := range params.AcceptList {
		voteTargets = append(voteTargets, v.SymbolDenom)
	}

	app.OracleKeeper.RewardBallotWinners(ctx, (int64)(app.OracleKeeper.VotePeriod(ctx)), (int64)(app.OracleKeeper.RewardDistributionWindow(ctx)), voteTargets, []types.Claim{})
	outstandingRewardsDecVal1, err := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr)
	s.Require().NoError(err)
	outstandingRewardsVal1, _ := outstandingRewardsDecVal1.TruncateDecimal()
	outstandingRewardsDecVal2, err := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr2)
	s.Require().NoError(err)
	outstandingRewardsVal2, _ := outstandingRewardsDecVal2.TruncateDecimal()
	s.Require().Equal(sdkmath.LegacyZeroDec().TruncateInt(), outstandingRewardsVal1.AmountOf(types.CheqdDenom))
	s.Require().Equal(sdkmath.LegacyZeroDec().TruncateInt(), outstandingRewardsVal2.AmountOf(types.CheqdDenom))
}
