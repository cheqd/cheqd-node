package testsuite

import (
	"cosmossdk.io/log"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cheqd/cheqd-node/app"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
)

type UpgradeTestSuite struct {
	suite.Suite

	App         *app.TestApp
	Ctx         sdk.Context
	QueryHelper *baseapp.QueryServiceTestHelper
	TestAccs    []sdk.AccAddress
}

// Setup sets up basic environment for suite (App, Ctx, and test accounts)
func (s *UpgradeTestSuite) Setup() {
	var err error
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(s.T())
	s.App, err = app.NewAppWithCustomOptions(false, app.SetupOptions{
		Logger:  logger.With("instance", "first"),
		DB:      db,
		AppOpts: simtestutil.NewAppOptionsWithFlagHome(s.T().TempDir()),
	})
	s.Require().NoError(err)
	s.Ctx = s.App.BaseApp.NewContext(false)
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}

	s.TestAccs = createRandomAccounts(3)
}

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}
