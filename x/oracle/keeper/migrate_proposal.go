package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1migrations "github.com/cosmos/cosmos-sdk/x/gov/migrations/v1"

	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// MigrateProposals migrates all legacy MsgUpgateGovParam proposals into non legacy param update versions.
func MigrateProposals(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	propStore := prefix.NewStore(store, v1migrations.ProposalsKeyPrefix)

	iter := propStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var prop govv1.Proposal
		err := cdc.Unmarshal(iter.Value(), &prop)
		// if error unmarshaling prop, convert to non legacy prop
		if err != nil {
			newProp, err := convertProposal(prop, cdc)
			if err != nil {
				return err
			}
			bz, err := cdc.Marshal(&newProp)
			if err != nil {
				return err
			}
			// Set new value on store.
			propStore.Set(iter.Key(), bz)
		}
	}

	return nil
}

func convertProposal(prop govv1.Proposal, cdc codec.BinaryCodec) (govv1.Proposal, error) {
	msgs := prop.Messages

	for _, msg := range msgs {
		var oldUpdateParamMsg oracletypes.MsgLegacyGovUpdateParams
		err := cdc.Unmarshal(msg.GetValue(), &oldUpdateParamMsg)
		if err != nil {
			return govv1.Proposal{}, err
		}

		newUpdateParamMsg := oracletypes.MsgGovUpdateParams{
			Authority:   oldUpdateParamMsg.Authority,
			Title:       oldUpdateParamMsg.Title,
			Description: oldUpdateParamMsg.Description,
			Plan: oracletypes.ParamUpdatePlan{
				Keys:    oldUpdateParamMsg.Keys,
				Height:  0, // placeholder value for height
				Changes: oldUpdateParamMsg.Changes,
			},
		}

		msg.Value, err = newUpdateParamMsg.Marshal()
		if err != nil {
			return govv1.Proposal{}, err
		}
	}

	prop.Messages = msgs
	return prop, nil
}
