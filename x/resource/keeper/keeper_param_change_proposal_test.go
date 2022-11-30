package keeper_test

import (
	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	cheqdsimapp "github.com/cheqd/cheqd-node/simapp"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type HandlerTestSuite struct {
	suite.Suite

	app        *cheqdsimapp.SimApp
	ctx        sdk.Context
	govHandler govv1beta1.Handler
}

func (suite *HandlerTestSuite) SetupTest() error {
	var err error
	suite.app, err = cheqdsimapp.Setup(false)
	if err != nil {
		return err
	}
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.govHandler = params.NewParamChangeProposalHandler(suite.app.ParamsKeeper)
	return nil
}

// func TestHandlerTestSuite(t *testing.T) {
// 	suite.Run(t, new(HandlerTestSuite))
// }

func testProposal(changes ...proposal.ParamChange) *proposal.ParameterChangeProposal {
	return proposal.NewParameterChangeProposal("title", "description", changes)
}

type TestCaseKeeperProposal struct {
	proposal *proposal.ParameterChangeProposal
	onHandle func(*HandlerTestSuite)
	expErr   bool
	errMsg   string
}

var _ = DescribeTable("Proposal Handler", func(testcase TestCaseKeeperProposal) {
	handlerSuite := new(HandlerTestSuite)
	err := handlerSuite.SetupTest()
	Expect(err).To(BeNil())

	err = handlerSuite.govHandler(handlerSuite.ctx, testcase.proposal)
	if testcase.expErr {
		Expect(err).NotTo(BeNil())
	} else {
		Expect(err).To(BeNil())
		testcase.onHandle(handlerSuite)
	}
},
	Entry("all fields",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "10000000000"}, "json": {"denom": "ncheq", "amount": "4000000000"}, "default": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "0.600000000000000000"}`,
			}),
			func(handlerSuite *HandlerTestSuite) {
				expectedFeeParams := resourcetypes.FeeParams{
					MediaTypes: map[string]sdk.Coin{
						resourcetypes.DefaultKeyCreateResourceImage: {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(10000000000)},
						resourcetypes.DefaultKeyCreateResourceJson:  {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(4000000000)},
						resourcetypes.DefaultKeyCreateResource:      {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(2000000000)},
					},
					BurnFactor: sdk.MustNewDecFromStr("0.600000000000000000"),
				}

				feeParams := handlerSuite.app.ResourceKeeper.GetParams(handlerSuite.ctx)

				Expect(expectedFeeParams).To(Equal(feeParams))
			},
			false,
			"",
		}),
	Entry("new media type added",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "5000000000"}, "json": {"denom": "ncheq", "amount": "2000000000"}, "default": {"denom": "ncheq", "amount": "1000000000"}, "text/html": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "0.500000000000000000"}`,
			}),
			func(handlerSuite *HandlerTestSuite) {
				expectedFeeParams := resourcetypes.FeeParams{
					MediaTypes: map[string]sdk.Coin{
						resourcetypes.DefaultKeyCreateResourceImage: {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(resourcetypes.DefaultCreateResourceImageFee)},
						resourcetypes.DefaultKeyCreateResourceJson:  {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(resourcetypes.DefaultCreateResourceJsonFee)},
						resourcetypes.DefaultKeyCreateResource:      {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(resourcetypes.DefaultCreateResourceDefaultFee)},
						"text/html":                                 {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(2000000000)},
					},
					BurnFactor: sdk.MustNewDecFromStr(resourcetypes.DefaultBurnFactor),
				}

				feeParams := handlerSuite.app.ResourceKeeper.GetParams(handlerSuite.ctx)

				Expect(expectedFeeParams).To(Equal(feeParams))
			},
			false,
			"",
		}),
	Entry("empty value",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		}),
	Entry("omit fields",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value: `
				{
					"media_types": {
						"image": {"denom": "ncheq", "amount": "10000000000"},
						"json": {"denom": "ncheq", "amount": "4000000000"}
					},
					"burn_factor": "0.600000000000000000"
				}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		}),
	Entry("invalid value: case `image` amount 0",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "0"}, "json": {"denom": "ncheq", "amount": "4000000000"}, "default": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "0.600000000000000000"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		}),
	Entry("invalid value: case `json` amount 0",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "10000000000"}, "json": {"denom": "ncheq", "amount": "0"}, "default": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "0.600000000000000000"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		}),
	Entry("invalid value: case `default` amount 0",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "10000000000"}, "json": {"denom": "ncheq", "amount": "4000000000"}, "default": {"denom": "ncheq", "amount": "0"}}, "burn_factor": "0.600000000000000000"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		}),
	Entry("invalid value: case `burn_factor` -1",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "10000000000"}, "json": {"denom": "ncheq", "amount": "4000000000"}, "default": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "-1"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		}),
	Entry("invalid value: case `burn_factor` 1.1",
		TestCaseKeeperProposal{
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "10000000000"}, "json": {"denom": "ncheq", "amount": "4000000000"}, "default": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "1.1"}`,
			}),
			func(*HandlerTestSuite) {},
			true,
			"",
		}),
)
