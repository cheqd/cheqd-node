package v4_test

import (
	"fmt"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"cosmossdk.io/core/header"
	"github.com/cheqd/cheqd-node/app/testsuite"
	v4 "github.com/cheqd/cheqd-node/app/upgrades/v4"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	"github.com/cometbft/cometbft/abci/types"

	"github.com/stretchr/testify/suite"
)

type UpgradeTestSuite struct {
	testsuite.UpgradeTestSuite
}

func TestUpgrade(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (suite *UpgradeTestSuite) TestV4MinorUpgradeHandler() {
	suite.Setup()

	const upgradeHeight = 10
	const upgradeName = v4.FeatureUpgradeName
	denom := oracletypes.DefaultParams().UsdcIbcDenom

	// -----------------------------
	// 1️⃣ Non-mainnet chain
	// -----------------------------
	ctx := suite.Ctx.WithBlockHeight(upgradeHeight - 1).WithChainID("testing-1").WithBlockTime(time.Now().UTC())
	ctx = ctx.WithHeaderInfo(header.Info{Height: upgradeHeight - 1})

	// Schedule the upgrade plan
	plan := upgradetypes.Plan{Name: upgradeName, Height: upgradeHeight}
	err := suite.App.UpgradeKeeper.ScheduleUpgrade(ctx, plan)
	suite.Require().NoError(err)

	// Advance to the upgrade height
	ctx = suite.Ctx.WithBlockHeight(upgradeHeight).WithHeaderInfo(header.Info{Height: upgradeHeight})
	// Trigger BeginBlocker — this runs the upgrade handler
	suite.Require().NotPanics(func() {
		fmt.Println(ctx.BlockHeight())
		r, _ := suite.App.UpgradeKeeper.GetUpgradePlan(ctx)
		fmt.Println(r)
		_, err = suite.App.PreBlocker(ctx, &types.RequestFinalizeBlock{})
	})

	// ✅ Check TWAP initialized
	rate, err := suite.App.FeeabsKeeper.GetTwapRate(ctx, denom)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.LegacyMustNewDecFromStr("1.0"), rate)
}
