package cheqd

import (
	"github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the cheqd module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the nym
	for _, elem := range genState.NymList {
		k.SetNym(ctx, *elem)
	}

	// Set nym count
	k.SetNymCount(ctx, uint64(len(genState.NymList)))

}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all nym
	nymList := k.GetAllNym(ctx)
	for _, elem := range nymList {
		elem := elem
		genesis.NymList = append(genesis.NymList, &elem)
	}

	return genesis
}
