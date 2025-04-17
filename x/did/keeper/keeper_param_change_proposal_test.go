package keeper_test

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	cheqdapp "github.com/cheqd/cheqd-node/app"
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func init() {
	cheqdapp.SetConfig()
}

type KeeperTestSuite struct {
	app         *cheqdapp.TestApp
	ctx         sdk.Context
	didKeeper   didkeeper.Keeper
	queryClient didtypes.QueryClient
	msgSvr      didtypes.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() error {
	var err error
	suite.app, err = cheqdapp.Setup(false)
	if err != nil {
		return err
	}

	suite.ctx = suite.app.BaseApp.NewContext(false)
	suite.didKeeper = suite.app.DidKeeper

	// Set default params
	err = suite.didKeeper.SetParams(suite.ctx, *didtypes.DefaultFeeParams())
	if err != nil {
		return err
	}

	suite.msgSvr = didkeeper.NewMsgServerImpl(suite.didKeeper)
	// Setup query client
	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	didtypes.RegisterQueryServer(queryHelper, suite.didKeeper)
	suite.queryClient = didtypes.NewQueryClient(queryHelper)

	return nil
}

type TestCaseUpdateParams struct {
	name      string
	input     *didtypes.MsgUpdateParams
	expErr    bool
	expErrMsg string
}

var _ = DescribeTable("UpdateParams", func(testCase TestCaseUpdateParams) {
	keeperSuite := new(KeeperTestSuite)
	err := keeperSuite.SetupTest()

	Expect(err).To(BeNil())
	if strings.TrimSpace(testCase.input.Authority) == "" {
		testCase.input.Authority = keeperSuite.didKeeper.GetAuthority()
	}
	// Call UpdateParams method
	_, err = keeperSuite.msgSvr.UpdateParams(keeperSuite.ctx, testCase.input)

	if testCase.expErr {
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring(testCase.expErrMsg))
	} else {
		Expect(err).To(BeNil())

		// Verify params were updated correctly
		params, err := keeperSuite.didKeeper.GetParams(keeperSuite.ctx)
		Expect(err).To(BeNil())
		Expect(params).To(Equal(testCase.input.Params))
	}
},
	Entry("valid params - all fields",
		TestCaseUpdateParams{
			name: "valid params - all fields",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
				},
			},
			expErr:    false,
			expErrMsg: "",
		}),
	Entry("invalid create_did amount 0",
		TestCaseUpdateParams{
			name: "invalid create_did amount 0",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(0)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid create did tx fee:",
		}),
	Entry("invalid create_did denom",
		TestCaseUpdateParams{
			name: "invalid create_did denom",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: "wrongdenom", Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid create did tx fee:",
		}),
	Entry("invalid update_did amount 0",
		TestCaseUpdateParams{
			name: "invalid update_did amount 0",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(0)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid update did tx fee:",
		}),
	Entry("invalid update_did denom",
		TestCaseUpdateParams{
			name: "invalid update_did denom",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: "wrongdenom", Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid update did tx fee:",
		}),
	Entry("invalid deactivate_did amount 0",
		TestCaseUpdateParams{
			name: "invalid deactivate_did amount 0",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(0)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid deactivate did tx fee:",
		}),
	Entry("invalid deactivate_did denom",
		TestCaseUpdateParams{
			name: "invalid deactivate_did denom",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: "wrongdenom", Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid deactivate did tx fee:",
		}),
	Entry("invalid burn_factor 0",
		TestCaseUpdateParams{
			name: "invalid burn_factor 0",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("0"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid burn factor:",
		}),
	Entry("invalid burn_factor negative",
		TestCaseUpdateParams{
			name: "invalid burn_factor negative",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("-0.1"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid burn factor:",
		}),
	Entry("invalid burn_factor equal to 1",
		TestCaseUpdateParams{
			name: "invalid burn_factor equal to 1",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("1.0"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid burn factor:",
		}),
	Entry("invalid burn_factor greater than 1",
		TestCaseUpdateParams{
			name: "invalid burn_factor greater than 1",
			input: &didtypes.MsgUpdateParams{
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("1.1"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid burn factor:",
		}),
	Entry("invalid authority",
		TestCaseUpdateParams{
			name: "invalid authority",
			input: &didtypes.MsgUpdateParams{
				Authority: "invalid",
				Params: didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("0.6"),
				},
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		}),
)
