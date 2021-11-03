package cheqd

import (
	"fmt"
	"github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the cheqd module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState v1.GenesisState) {
	for _, elem := range genState.DidList {
		did, err := elem.GetDid()
		if err != nil {
			panic(fmt.Sprintf("Cannot import geneses case: %s", err.Error()))
		}

		if err = k.SetDid(ctx, *did, elem.Metadata); err != nil {
			panic(fmt.Sprintf("Cannot set did case: %s", err.Error()))
		}
	}

	// Set nym count
	k.SetDidCount(ctx, uint64(len(genState.DidList)))

	k.SetDidNamespace(ctx, genState.DidNamespace)
}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *v1.GenesisState {
	genesis := v1.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all did
	didList := k.GetAllDid(ctx)
	for _, elem := range didList {
		elem := elem
		genesis.DidList = append(genesis.DidList, &elem)
	}

	genesis.DidNamespace = k.GetDidNamespace(ctx)

	return genesis
}
