package cheqd

import (
	"github.com/cheqd/cheqd-node/x/did/keeper"
	"github.com/cheqd/cheqd-node/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the cheqd module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set didocs
	for _, versionSet := range genState.VersionSets {
		for _, didDoc := range versionSet.DidDocs {
			k.SetDidDocVersion(&ctx, didDoc, false)
		}

		k.SetLatestDidDocVersion(&ctx, versionSet.DidDocs[0].DidDoc.Id, versionSet.LatestVersion)
	}

	// Set did namespace
	k.SetDidNamespace(&ctx, genState.DidNamespace)
}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.GenesisState{
		DidNamespace: k.GetDidNamespace(&ctx),
		VersionSets:  k.GetAllDidDocs(&ctx),
	}

	return &genesis
}
