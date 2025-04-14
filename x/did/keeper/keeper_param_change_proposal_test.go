package keeper_test

import (
	"github.com/stretchr/testify/suite"

	cheqdapp "github.com/cheqd/cheqd-node/app"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	sdkmath "cosmossdk.io/math"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type HandlerTestSuite struct {
	suite.Suite

	app        *cheqdapp.TestApp
	ctx        sdk.Context
	govHandler govv1beta1.Handler
}

func (suite *HandlerTestSuite) SetupTest() error {
	var err error
	suite.app, err = cheqdapp.Setup(false)
	if err != nil {
		return err
	}
	suite.ctx = suite.app.BaseApp.NewContext(false)
	suite.govHandler = params.NewParamChangeProposalHandler(suite.app.ParamsKeeper)
	return nil
}

func testProposal(changes ...proposal.ParamChange) *proposal.ParameterChangeProposal {
	return proposal.NewParameterChangeProposal("title", "description", changes)
}

type TestCaseKeeperProposal struct {
	proposal *proposal.ParameterChangeProposal
	onHandle func(*HandlerTestSuite)
	expErr   bool
	errMsg   string
}

var _ = DescribeTable("Proposal Handler", func(testCase TestCaseKeeperProposal) {
	handlerSuite := new(HandlerTestSuite)
	err := handlerSuite.SetupTest()
	Expect(err).To(BeNil())

	err = handlerSuite.govHandler(handlerSuite.ctx, testCase.proposal)
	if testCase.expErr {
		Expect(err).NotTo(BeNil())
	} else {
		Expect(err).To(BeNil())
		testCase.onHandle(handlerSuite)
	}
},
	Entry("all fields",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: didtypes.ModuleName,
				Key:      string(didtypes.ParamsKey),
				Value:    `{"create_did": {"denom": "ncheq", "amount": "10000000000"}, "update_did": {"denom": "ncheq", "amount": "4000000000"}, "deactivate_did": {"denom": "ncheq", "amount": "2000000000"}, "burn_factor": "0.600000000000000000"}`,
			}),
			func(handlerSuite *HandlerTestSuite) {
				expectedFeeParams := didtypes.FeeParams{
					CreateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(10000000000)},
					UpdateDid:     sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(4000000000)},
					DeactivateDid: sdk.Coin{Denom: didtypes.BaseMinimalDenom, Amount: sdkmath.NewInt(2000000000)},
					BurnFactor:    sdkmath.LegacyMustNewDecFromStr("0.600000000000000000"),
				}

				feeParams, err := handlerSuite.app.DidKeeper.Params.Get(handlerSuite.ctx)
				Expect(err).To(BeNil())
				Expect(expectedFeeParams).To(Equal(feeParams))
			},
			false,
			"",
		}),
	Entry("empty value",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: didtypes.ModuleName,
				Key:      string(didtypes.ParamsKey),
				Value:    `{}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		},
	),
	Entry("omit fields",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: didtypes.ModuleName,
				Key:      string(didtypes.ParamsKey),
				Value:    `{"create_did": {"denom": "ncheq", "amount": "10000000000"}, "update_did": {"denom": "ncheq", "amount": "4000000000"}, "burn_factor": "0.600000000000000000"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		},
	),
	Entry("invalid value: case `create_did` amount 0",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: didtypes.ModuleName,
				Key:      string(didtypes.ParamsKey),
				Value:    `{"create_did": {"denom": "ncheq", "amount": "0"}, "update_did": {"denom": "ncheq", "amount": "4000000000"}, "deactivate_did": {"denom": "ncheq", "amount": "2000000000"}, "burn_factor": "0.600000000000000000"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		},
	),
	Entry("invalid value: case `update_did` amount 0",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: didtypes.ModuleName,
				Key:      string(didtypes.ParamsKey),
				Value:    `{"create_did": {"denom": "ncheq", "amount": "10000000000"}, "update_did": {"denom": "ncheq", "amount": "0"}, "deactivate_did": {"denom": "ncheq", "amount": "2000000000"}, "burn_factor": "0.600000000000000000"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		}),
	Entry("invalid value: case `deactivate_did` amount 0",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: didtypes.ModuleName,
				Key:      string(didtypes.ParamsKey),
				Value:    `{"create_did": {"denom": "ncheq", "amount": "10000000000"}, "update_did": {"denom": "ncheq", "amount": "4000000000"}, "deactivate_did": {"denom": "ncheq", "amount": "0"}, "burn_factor": "0.600000000000000000"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		},
	),
	Entry("invalid value: case `burn_factor` -1",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: didtypes.ModuleName,
				Key:      string(didtypes.ParamsKey),
				Value:    `{"create_did": {"denom": "ncheq", "amount": "10000000000"}, "update_did": {"denom": "ncheq", "amount": "4000000000"}, "deactivate_did": {"denom": "ncheq", "amount": "2000000000"}, "burn_factor": "-1"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		},
	),
	Entry("invalid value: case `burn_factor` 1.1",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: didtypes.ModuleName,
				Key:      string(didtypes.ParamsKey),
				Value:    `{"create_did": {"denom": "ncheq", "amount": "10000000000"}, "update_did": {"denom": "ncheq", "amount": "4000000000"}, "deactivate_did": {"denom": "ncheq", "amount": "2000000000"}, "burn_factor": "1.1"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		},
	),
)
