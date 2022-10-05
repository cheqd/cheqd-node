package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	cheqdsimapp "github.com/cheqd/cheqd-node/simapp"
	cheqdtypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	// stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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
	}{
		{
			"all fields",
			testProposal(proposal.ParamChange{
				Subspace: cheqdtypes.ModuleName,
				Key:      cheqdtypes.FeeParamsKey,
				Value:    `{"create_did": {"denom": "ncheq", "amount": "10000000000"}, "update_did": {"denom": "ncheq", "amount": "4000000000"}, "deactivate_did": {"denom": "ncheq", "amount": "2000000000"}, "burn_factor": "0.600000000000000000"}`,
			}),
			func() {
				// TODO: Refactor to comply with v.0.46.1 ParamsKeeper Subspace API
				feeParams := suite.app.CheqdKeeper.GetParams(suite.ctx)
				expected := sdk.NewCoins(sdk.NewCoin(cheqdtypes.BaseMinimalDenom, sdk.NewInt(10000000000)))
				fmt.Println(expected)
				fmt.Println(feeParams.CreateDid)
				suite.Require().Equal(expected, feeParams.CreateDid)
			},
			false,
		},
		{
			"invalid type",
			testProposal(proposal.NewParamChange(cheqdtypes.ModuleName, string(cheqdtypes.FeeParamsKey), `{}`)),
			func() {},
			true,
		},
		{
			"omit empty fields",
			testProposal(proposal.ParamChange{
				Subspace: govtypes.ModuleName,
				Key:      string(govv1.ParamStoreKeyDepositParams),
				Value:    `{"min_deposit": [{"denom": "uatom","amount": "64000000"}], "max_deposit_period": "172800000000000"}`,
			}),
			func() {
				depositParams := suite.app.GovKeeper.GetDepositParams(suite.ctx)
				defaultPeriod := govv1.DefaultPeriod
				suite.Require().Equal(govv1.DepositParams{
					MinDeposit:       sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(64000000))),
					MaxDepositPeriod: &defaultPeriod,
				}, depositParams)
			},
			false,
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
