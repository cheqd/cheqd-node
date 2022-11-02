package cheqd

import (
	"fmt"

	"github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the cheqd module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set didocs
	for _, elem := range genState.DidDocs {
		if err := k.SetDidDoc(&ctx, elem); err != nil {
			panic(fmt.Sprintf("Cannot set did case: %s", err.Error()))
		}
	}

	// Set did namespace
	k.SetDidNamespace(&ctx, genState.DidNamespace)
}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.GenesisState{
		DidNamespace: k.GetDidNamespace(&ctx),
	}

	// Add diddocs
	didDocs := k.GetAllDidDocs(&ctx)
	for _, elem := range didDocs {
		genesis.DidDocs = append(genesis.DidDocs, &elem)
	}

	return &genesis
}
