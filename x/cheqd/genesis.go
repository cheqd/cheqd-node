package cheqd

import (
	"fmt"

	"github.com/canow-co/cheqd-node/x/cheqd/keeper"
	"github.com/canow-co/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the cheqd module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, elem := range genState.DidList {
		if err := k.SetDid(&ctx, elem); err != nil {
			panic(fmt.Sprintf("Cannot set did case: %s", err.Error()))
		}
	}

	// Set nym count
	k.SetDidCount(&ctx, uint64(len(genState.DidList)))

	k.SetDidNamespace(&ctx, genState.DidNamespace)
}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all did
	didList := k.GetAllDid(&ctx)
	for _, elem := range didList {
		elem := elem
		genesis.DidList = append(genesis.DidList, &elem)
	}

	genesis.DidNamespace = k.GetDidNamespace(&ctx)

	return genesis
}
