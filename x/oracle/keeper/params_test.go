package keeper_test

// import (
// 	"cosmossdk.io/math"
// 	"github.com/cheqd/cheqd-node/x/oracle/types"
// )

// func (s *IntegrationTestSuite) TestVoteThreshold() {
// 	app, ctx := s.app, s.ctx

// 	voteDec := app.OracleKeeper.VoteThreshold(ctx)
// 	s.Require().Equal(math.LegacyMustNewDecFromStr("0.5"), voteDec)

// 	newVoteTreshold := math.LegacyMustNewDecFromStr("0.6")
// 	defaultParams := types.DefaultParams()
// 	defaultParams.VoteThreshold = newVoteTreshold
// 	app.OracleKeeper.SetParams(ctx, defaultParams)

// 	voteThresholdDec := app.OracleKeeper.VoteThreshold(ctx)
// 	s.Require().Equal(newVoteTreshold, voteThresholdDec)
// }
