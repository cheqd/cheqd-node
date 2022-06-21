package resource

import (
	"fmt"

	"github.com/cheqd/cheqd-node/x/resource/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the cheqd module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, resource := range genState.ResourceList {
		if err := k.SetResource(&ctx, resource); err != nil {
			panic(fmt.Sprintf("Cannot set resource case: %s", err.Error()))
		}
	}

	// Set nym count
	k.SetResourceCount(&ctx, uint64(len(genState.ResourceList)))

	// k.SetResourceNamespace(ctx, genState.ResourceNamespace)
}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all resource
	//resourceList := k.GetAllResource(&ctx)
	//for _, elem := range resourceList {
	//	elem := elem
	//	genesis.ResourceList = append(genesis.ResourceList, &elem)
	//}
	//
	//genesis.ResourceNamespace = k.GetResourceNamespace(ctx)

	return genesis
}
