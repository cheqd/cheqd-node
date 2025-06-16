package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cheqdapp "github.com/cheqd/cheqd-node/app"
	utils "github.com/cheqd/cheqd-node/util"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	oraclekeeper "github.com/cheqd/cheqd-node/x/oracle/keeper"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"

	"github.com/cheqd/cheqd-node/x/did/keeper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cheqd DID Module - MsgUpdateParams Suite")
}

type KeeperTestSuite struct {
	app         *cheqdapp.TestApp
	ctx         sdk.Context
	didKeeper   keeper.Keeper
	queryClient didtypes.QueryClient
	msgSvr      didtypes.MsgServer
}

func NewTestSuite(t GinkgoTInterface) *KeeperTestSuite {
	return &KeeperTestSuite{}
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
	suite.app.OracleKeeper.SetAverage(suite.ctx, oraclekeeper.KeyWMAWithStrategy(
		oracletypes.CheqdSymbol,
		string(oraclekeeper.WmaStrategyBalanced)), sdkmath.LegacyMustNewDecFromStr("0.016"))
	suite.msgSvr = keeper.NewMsgServerImpl(suite.didKeeper)

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

var _ = Describe("MsgServer MsgUpdateParams", Ordered, func() {
	var suite *KeeperTestSuite

	BeforeAll(func() {
		suite = NewTestSuite(GinkgoT())
	})

	DescribeTable("UpdateParams validation",
		func(tc TestCaseUpdateParams) {
			err := suite.SetupTest()
			Expect(err).ToNot(HaveOccurred())

			msgServer := keeper.NewMsgServerImpl(suite.app.DidKeeper)
			ctx := suite.ctx

			_, err = msgServer.UpdateParams(ctx, tc.input)

			if tc.expErr {
				Expect(err).To(HaveOccurred(), "expected error but got nil")
				if tc.expErrMsg != "" {
					Expect(err.Error()).To(ContainSubstring(tc.expErrMsg))
				}
			} else {
				Expect(err).ToNot(HaveOccurred(), "expected no error but got %s", err)
			}
		},

		Entry("invalid params - no overlapping range (usd + cheq)",
			TestCaseUpdateParams{
				name: "invalid params - no overlapping range (usd + cheq)",
				input: &didtypes.MsgUpdateParams{
					Authority: "cheqd10d07y265gmmuvt4z0w9aw880jnsr700j5ql9az",
					Params: didtypes.FeeParams{
						CreateDid: []didtypes.FeeRange{
							{
								Denom:     didtypes.BaseMinimalDenom,
								MinAmount: sdkmath.NewInt(10000000000),
								MaxAmount: utils.PtrInt(30000000000),
							},
							{
								Denom:     oracletypes.UsdDenom,
								MinAmount: sdkmath.NewInt(500_000_000_000_000_000), // 0.5 USD
								MaxAmount: utils.PtrInt(800_000_000_000_000_000),   // 0.8 USD
							},
						},
						UpdateDid: []didtypes.FeeRange{{
							Denom:     didtypes.BaseMinimalDenom,
							MinAmount: sdkmath.NewInt(25000000000),
							MaxAmount: nil,
						}},
						DeactivateDid: []didtypes.FeeRange{{
							Denom:     didtypes.BaseMinimalDenom,
							MinAmount: sdkmath.NewInt(10000000000),
							MaxAmount: utils.PtrInt(20000000000),
						}},
						BurnFactor: sdkmath.LegacyMustNewDecFromStr("0.6"),
					},
				},
				expErr:    true,
				expErrMsg: "no overlapping fee range found",
			},
		),

		Entry("invalid params - gap between cheq and usd range",
			TestCaseUpdateParams{
				name: "invalid params - gap between cheq and usd range",
				input: &didtypes.MsgUpdateParams{
					Authority: "cheqd10d07y265gmmuvt4z0w9aw880jnsr700j5ql9az",
					Params: didtypes.FeeParams{
						CreateDid: []didtypes.FeeRange{
							{
								Denom:     didtypes.BaseMinimalDenom,
								MinAmount: sdkmath.NewInt(5000000000),
								MaxAmount: utils.PtrInt(20000000000),
							},
							{
								Denom:     oracletypes.UsdDenom,
								MinAmount: sdkmath.NewInt(500_000_000_000_000_000),
								MaxAmount: utils.PtrInt(1_000_000_000_000_000_000),
							},
						},
						UpdateDid: []didtypes.FeeRange{{
							Denom:     didtypes.BaseMinimalDenom,
							MinAmount: sdkmath.NewInt(25000000000),
						}},
						DeactivateDid: []didtypes.FeeRange{{
							Denom:     didtypes.BaseMinimalDenom,
							MinAmount: sdkmath.NewInt(10000000000),
							MaxAmount: utils.PtrInt(20000000000),
						}},
						BurnFactor: sdkmath.LegacyMustNewDecFromStr("0.6"),
					},
				},
				expErr:    true,
				expErrMsg: "no overlapping fee range found",
			},
		),

		Entry("valid params with overlapping range",
			TestCaseUpdateParams{
				name: "valid params with overlapping range",
				input: &didtypes.MsgUpdateParams{
					Authority: "cheqd10d07y265gmmuvt4z0w9aw880jnsr700j5ql9az",
					Params: didtypes.FeeParams{
						CreateDid: []didtypes.FeeRange{
							{
								Denom:     didtypes.BaseMinimalDenom,
								MinAmount: sdkmath.NewInt(50000000000),
								MaxAmount: utils.PtrInt(100000000000),
							},
							{
								Denom:     oracletypes.UsdDenom,
								MinAmount: sdkmath.NewInt(1200000000000000000),
								MaxAmount: utils.PtrInt(2000000000000000000),
							},
						},
						UpdateDid: []didtypes.FeeRange{{
							Denom:     didtypes.BaseMinimalDenom,
							MinAmount: sdkmath.NewInt(25000000000),
							MaxAmount: nil,
						}},
						DeactivateDid: []didtypes.FeeRange{{
							Denom:     didtypes.BaseMinimalDenom,
							MinAmount: sdkmath.NewInt(10000000000),
							MaxAmount: utils.PtrInt(20000000000),
						}},
						BurnFactor: sdkmath.LegacyMustNewDecFromStr("0.5"),
					},
				},
				expErr: false,
			},
		),

		Entry("valid params with only single price",
			TestCaseUpdateParams{
				name: "valid params fixed prices",
				input: &didtypes.MsgUpdateParams{
					Authority: "cheqd10d07y265gmmuvt4z0w9aw880jnsr700j5ql9az",
					Params: didtypes.FeeParams{
						CreateDid: []didtypes.FeeRange{
							{
								Denom:     didtypes.BaseMinimalDenom,
								MinAmount: sdkmath.NewInt(50000000000),
								MaxAmount: utils.PtrInt(100000000000),
							},
						},
						UpdateDid: []didtypes.FeeRange{{
							Denom:     didtypes.BaseMinimalDenom,
							MinAmount: sdkmath.NewInt(25000000000),
							MaxAmount: nil,
						}},
						DeactivateDid: []didtypes.FeeRange{{
							Denom:     didtypes.BaseMinimalDenom,
							MinAmount: sdkmath.NewInt(10000000000),
							MaxAmount: utils.PtrInt(20000000000),
						}},
						BurnFactor: sdkmath.LegacyMustNewDecFromStr("0.5"),
					},
				},
				expErr: false,
			},
		),
	)
})
