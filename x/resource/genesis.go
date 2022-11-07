package resource

import (
	"fmt"

	"github.com/cheqd/cheqd-node/x/resource/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the resource module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, resource := range genState.Resources {
		if err := k.SetResource(&ctx, resource); err != nil {
			panic(fmt.Sprintf("Cannot set resource case: %s", err.Error()))
		}
	}
}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.GenesisState{}

	// Get all resource
	resourceList := k.GetAllResources(&ctx)
	for _, elem := range resourceList {
		genesis.Resources = append(genesis.Resources, &elem)
	}

	return &genesis
}
