package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	cheqdsimapp "github.com/cheqd/cheqd-node/simapp"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

type HandlerTestSuite struct {
	suite.Suite

	app        *cheqdsimapp.SimApp
	ctx        sdk.Context
	govHandler govv1beta1.Handler
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app = cheqdsimapp.Setup(suite.T(), false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.govHandler = params.NewParamChangeProposalHandler(suite.app.ParamsKeeper)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func testProposal(changes ...proposal.ParamChange) *proposal.ParameterChangeProposal {
	return proposal.NewParameterChangeProposal("title", "description", changes)
}

func (suite *HandlerTestSuite) TestProposalHandler() {
	testCases := []struct {
		name     string
		proposal *proposal.ParameterChangeProposal
		onHandle func()
		expErr   bool
		errMsg   string
	}{
		{
			"all fields",
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "10000000000"}, "json": {"denom": "ncheq", "amount": "4000000000"}, "default": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "0.600000000000000000"}`,
			}),
			func() {
				expectedFeeParams := resourcetypes.FeeParams{
					MediaTypes: map[string]sdk.Coin{
						resourcetypes.DefaultKeyCreateResourceImage: {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(10000000000)},
						resourcetypes.DefaultKeyCreateResourceJson:  {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(4000000000)},
						resourcetypes.DefaultKeyCreateResource:      {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(2000000000)},
					},
					BurnFactor: sdk.MustNewDecFromStr("0.600000000000000000"),
				}

				feeParams := suite.app.ResourceKeeper.GetParams(suite.ctx)

				suite.Require().Equal(expectedFeeParams, feeParams)
			},
			false,
			"",
		},
		{
			"new media type added",
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "5000000000"}, "json": {"denom": "ncheq", "amount": "2000000000"}, "default": {"denom": "ncheq", "amount": "1000000000"}, "text/html": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "0.500000000000000000"}`,
			}),
			func() {
				expectedFeeParams := resourcetypes.FeeParams{
					MediaTypes: map[string]sdk.Coin{
						resourcetypes.DefaultKeyCreateResourceImage: {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(resourcetypes.DefaultCreateResourceImageFee)},
						resourcetypes.DefaultKeyCreateResourceJson:  {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(resourcetypes.DefaultCreateResourceJsonFee)},
						resourcetypes.DefaultKeyCreateResource:      {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(resourcetypes.DefaultCreateResourceDefaultFee)},
						"text/html":                                 {Denom: resourcetypes.BaseMinimalDenom, Amount: sdk.NewInt(2000000000)},
					},
					BurnFactor: sdk.MustNewDecFromStr(resourcetypes.DefaultBurnFactor),
				}

				feeParams := suite.app.ResourceKeeper.GetParams(suite.ctx)

				suite.Require().Equal(expectedFeeParams, feeParams)
			},
			false,
			"",
		},
		{
			"empty value",
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{}`,
			}),
			func() {},
			true,
			"",
		},
		{
			"omit fields",
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
			func() {},
			true,
			"",
		},
		{
			"invalid value: case `image` amount 0",
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "0"}, "json": {"denom": "ncheq", "amount": "4000000000"}, "default": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "0.600000000000000000"}`,
			}),
			func() {},
			true,
			"",
		},
		{
			"invalid value: case `json` amount 0",
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "10000000000"}, "json": {"denom": "ncheq", "amount": "0"}, "default": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "0.600000000000000000"}`,
			}),
			func() {},
			true,
			"",
		},
		{
			"invalid value: case `default` amount 0",
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "10000000000"}, "json": {"denom": "ncheq", "amount": "4000000000"}, "default": {"denom": "ncheq", "amount": "0"}}, "burn_factor": "0.600000000000000000"}`,
			}),
			func() {},
			true,
			"",
		},
		{
			"invalid value: case `burn_factor` -1",
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "10000000000"}, "json": {"denom": "ncheq", "amount": "4000000000"}, "default": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "-1"}`,
			}),
			func() {},
			true,
			"",
		},
		{
			"invalid value: case `burn_factor` 1.1",
			testProposal(proposal.ParamChange{
				Subspace: resourcetypes.ModuleName,
				Key:      string(resourcetypes.ParamStoreKeyFeeParams),
				Value:    `{"media_types": {"image": {"denom": "ncheq", "amount": "10000000000"}, "json": {"denom": "ncheq", "amount": "4000000000"}, "default": {"denom": "ncheq", "amount": "2000000000"}}, "burn_factor": "1.1"}`,
			}),
			func() {},
			true,
			"",
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			err := suite.govHandler(suite.ctx, tc.proposal)
			if tc.expErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tc.onHandle()
			}
		})
	}
}
