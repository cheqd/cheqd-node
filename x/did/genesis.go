package cheqd

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/keeper"
	"github.com/cheqd/cheqd-node/x/did/types"
)

// InitGenesis initializes the cheqd module's state from a provided genesis
// state.
func InitGenesis(ctx context.Context, k keeper.Keeper, genState *types.GenesisState) {
	// Set didocs
	for _, versionSet := range genState.VersionSets {
		for _, didDoc := range versionSet.DidDocs {
			err := k.SetDidDocVersion(&ctx, didDoc, false)
			if err != nil {
				panic(err)
			}
		}

		err := k.SetLatestDidDocVersion(&ctx, versionSet.DidDocs[0].DidDoc.Id, versionSet.LatestVersion)
		if err != nil {
			panic(err)
		}
	}

	// Set did namespace
	err := k.SetDidNamespace(&ctx, genState.DidNamespace)
	if err != nil {
		panic(err)
	}

	// Set fee params
	err = k.SetParams(ctx, *genState.FeeParams)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx context.Context, k keeper.Keeper) *types.GenesisState {
	didDocs, err := k.GetAllDidDocs(&ctx)
	if err != nil {
		panic(err)
	}
	feeParams, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	DidNameSpace, err := k.GetDidNamespace(&ctx)
	if err != nil {
		panic(err)
	}
	genesis := types.GenesisState{
		DidNamespace: DidNameSpace,
		VersionSets:  didDocs,
		FeeParams:    &feeParams,
	}

	return &genesis
}
