package keeper_test

// import (
// 	"strings"

// 	sdkmath "cosmossdk.io/math"
// 	cheqdapp "github.com/cheqd/cheqd-node/app"
// 	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
// 	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// )

// type KeeperTestSuite struct {
// 	app            *cheqdapp.TestApp
// 	ctx            sdk.Context
// 	resourceKeeper resourcekeeper.Keeper
// 	msgSvr         resourcetypes.MsgServer
// }

// func (suite *KeeperTestSuite) SetupTest() error {
// 	var err error
// 	suite.app, err = cheqdapp.Setup(false)
// 	if err != nil {
// 		return err
// 	}

// 	suite.ctx = suite.app.BaseApp.NewContext(false)
// 	suite.resourceKeeper = suite.app.ResourceKeeper

// 	// Set default params
// 	err = suite.resourceKeeper.SetParams(suite.ctx, *resourcetypes.DefaultFeeParams())
// 	if err != nil {
// 		return err
// 	}

// 	suite.msgSvr = resourcekeeper.NewMsgServerImpl(suite.app.ResourceKeeper, suite.app.DidKeeper)
// 	return nil
// }

// type TestCaseUpdateParams struct {
// 	name   string
// 	input  *resourcetypes.MsgUpdateParams
// 	expErr bool
// }

// var _ = DescribeTable("UpdateParams", func(testCase TestCaseUpdateParams) {
// 	keeperSuite := new(KeeperTestSuite)
// 	err := keeperSuite.SetupTest()

// 	Expect(err).To(BeNil())
// 	// ✅ Only set authority if the test case didn't define one
// 	if strings.TrimSpace(testCase.input.Authority) == "" {
// 		testCase.input.Authority = keeperSuite.resourceKeeper.GetAuthority()
// 	}
// 	// Call UpdateParams method
// 	_, err = keeperSuite.msgSvr.UpdateParams(keeperSuite.ctx, testCase.input)

// 	if testCase.expErr {
// 		Expect(err).NotTo(BeNil())
// 		// fmt.Println("error here....", testCase.name, err.Error())
// 		// fmt.Println("error message....", err.Error())
// 		// Expect(err.Error()).To(ContainSubstring(testCase.expErrMsg))
// 	} else {
// 		Expect(err).To(BeNil())

// 		// Verify params were updated correctly
// 		params, err := keeperSuite.resourceKeeper.GetParams(keeperSuite.ctx)
// 		Expect(err).To(BeNil())
// 		Expect(params).To(Equal(testCase.input.Params))
// 	}
// },
// 	Entry("valid params - all fields",
// 		TestCaseUpdateParams{
// 			name: "valid params - all fields",
// 			input: &resourcetypes.MsgUpdateParams{
// 				Params: resourcetypes.FeeParams{
// 					Image:      sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
// 					Json:       sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
// 					Default:    sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
// 					BurnFactor: sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
// 				},
// 			},
// 			expErr: false,
// 		}),
// 	Entry("invalid image amount 0",
// 		TestCaseUpdateParams{
// 			name: "invalid create_did amount 0",
// 			input: &resourcetypes.MsgUpdateParams{
// 				Params: resourcetypes.FeeParams{
// 					Image:      sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(0)},
// 					Json:       sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
// 					Default:    sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
// 					BurnFactor: sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
// 				},
// 			},
// 			expErr: true,
// 		}),
// 	Entry("invalid image denom",
// 		TestCaseUpdateParams{
// 			name: "invalid image denom",
// 			input: &resourcetypes.MsgUpdateParams{
// 				Params: resourcetypes.FeeParams{
// 					Image:      sdk.Coin{Denom: "", Amount: sdkmath.NewInt(10000000000)},
// 					Json:       sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
// 					Default:    sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
// 					BurnFactor: sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
// 				},
// 			},
// 			expErr: true,
// 		}),
// 	Entry("invalid json amount 0",
// 		TestCaseUpdateParams{
// 			name: "invalid json amount 0",
// 			input: &resourcetypes.MsgUpdateParams{
// 				Params: resourcetypes.FeeParams{
// 					Image:      sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
// 					Json:       sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(0)},
// 					Default:    sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
// 					BurnFactor: sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
// 				},
// 			},
// 			expErr: true,
// 		}),
// 	Entry("invalid json denom",
// 		TestCaseUpdateParams{
// 			name: "invalid json denom",
// 			input: &resourcetypes.MsgUpdateParams{
// 				Params: resourcetypes.FeeParams{
// 					Image:      sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
// 					Json:       sdk.Coin{Denom: "wrongdenom", Amount: sdkmath.NewInt(4000000000)},
// 					Default:    sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
// 					BurnFactor: sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
// 				},
// 			},
// 			expErr: true,
// 		}),
// 	Entry("invalid burn_factor 0",
// 		TestCaseUpdateParams{
// 			name: "invalid burn_factor 0",
// 			input: &resourcetypes.MsgUpdateParams{
// 				Params: resourcetypes.FeeParams{
// 					Image:      sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
// 					Json:       sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
// 					Default:    sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
// 					BurnFactor: sdkmath.LegacyMustNewDecFromStr("0"),
// 				},
// 			},
// 			expErr: true,
// 		}),
// 	Entry("invalid burn_factor negative",
// 		TestCaseUpdateParams{
// 			name: "invalid burn_factor negative",
// 			input: &resourcetypes.MsgUpdateParams{
// 				Params: resourcetypes.FeeParams{
// 					Image:      sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
// 					Json:       sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
// 					Default:    sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
// 					BurnFactor: sdkmath.LegacyMustNewDecFromStr("-0.1"),
// 				},
// 			},
// 			expErr: true,
// 		}),
// 	Entry("invalid burn_factor equal to 1",
// 		TestCaseUpdateParams{
// 			name: "invalid burn_factor equal to 1",
// 			input: &resourcetypes.MsgUpdateParams{
// 				Params: resourcetypes.FeeParams{
// 					Image:      sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
// 					Json:       sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
// 					Default:    sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
// 					BurnFactor: sdkmath.LegacyMustNewDecFromStr("1.0"),
// 				},
// 			},
// 			expErr: true,
// 		}),
// 	Entry("invalid burn_factor greater than 1",
// 		TestCaseUpdateParams{
// 			name: "invalid burn_factor greater than 1",
// 			input: &resourcetypes.MsgUpdateParams{
// 				Params: resourcetypes.FeeParams{
// 					Image:      sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
// 					Json:       sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
// 					Default:    sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
// 					BurnFactor: sdkmath.LegacyMustNewDecFromStr("1.1"),
// 				},
// 			},
// 			expErr: true,
// 		}),
// 	Entry("invalid authority",
// 		TestCaseUpdateParams{
// 			name: "invalid authority",
// 			input: &resourcetypes.MsgUpdateParams{
// 				Authority: "invalid",
// 				Params: resourcetypes.FeeParams{
// 					Image:      sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
// 					Json:       sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
// 					Default:    sdk.Coin{Denom: resourcetypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
// 					BurnFactor: sdkmath.LegacyMustNewDecFromStr("0.6"),
// 				},
// 			},
// 			expErr: true,
// 		}),
// )
