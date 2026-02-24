package keeper_test

import (
	types1 "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/cheqd/cheqd-node/x/oracle/keeper"
	"github.com/cheqd/cheqd-node/x/oracle/types"
)

func (s *IntegrationTestSuite) TestMigrateProposal() {
	ctx := s.ctx
	cdc := s.app.AppCodec()
	storeKey := s.app.GetKey(govtypes.StoreKey)

	// create legacy prop and set it in store
	legacyMsg := types.MsgLegacyGovUpdateParams{
		Authority:   "cheqd10d07y265gmmuvt4z0w9aw880jnsr700j5ql9az",
		Title:       "title",
		Description: "desc",
		Keys: []string{
			"VotePeriod",
		},
		Changes: types.Params{
			VotePeriod: 5,
		},
	}
	bz, err := cdc.Marshal(&legacyMsg)
	s.Require().NoError(err)
	prop := govv1.Proposal{
		Id: 1,
		Messages: []*types1.Any{
			{
				TypeUrl:          "/cheqd.oracle.v2.MsgGovUpdateParams",
				Value:            bz,
				XXX_unrecognized: []byte{},
			},
		},
		Status: govv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
	}
	err = s.app.GovKeeper.SetProposal(ctx, prop)
	s.Require().NoError(err)

	// try to retrieve proposal and fail
	_, err = s.app.GovKeeper.Proposals.Get(ctx, prop.Id)
	s.Require().Error(err)

	// successfully retrieve proposal after migration
	err = keeper.MigrateProposals(ctx, storeKey, cdc)
	s.Require().NoError(err)

	_, err = s.app.GovKeeper.Proposals.Get(ctx, prop.Id)
	s.Require().NoError(err)
}
