package resource

import (
	"fmt"

	"github.com/cheqd/cheqd-node/x/resource/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the resource module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState *types.GenesisState) {
	for _, resource := range genState.Resources {
		if err := k.SetResource(&ctx, resource); err != nil {
			panic(fmt.Sprintf("Cannot set resource case: %s", err.Error()))
		}
	}

	// set fee params
	k.SetParams(ctx, *genState.FeeParams)
}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	// Get all resource
	resourceList, err := k.GetAllResources(&ctx)
	if err != nil {
		panic(fmt.Sprintf("Cannot get all resource: %s", err.Error()))
	}

	// get fee params
	feeParams := k.GetParams(ctx)

	return &types.GenesisState{
		Resources: resourceList,
		FeeParams: &feeParams,
	}
}
