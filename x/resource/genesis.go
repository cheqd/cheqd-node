package resource

import (
	"context"
	"fmt"

	"github.com/cheqd/cheqd-node/x/resource/keeper"
	"github.com/cheqd/cheqd-node/x/resource/types"
)

// InitGenesis initializes the resource module's state from a provided genesis
// state.
func InitGenesis(ctx context.Context, k keeper.Keeper, genState *types.GenesisState) {
	for _, resource := range genState.Resources {
		if err := k.SetResource(ctx, resource); err != nil {
			panic(fmt.Sprintf("Cannot set resource case: %s", err.Error()))
		}
	}

	// set fee params
	if err := k.SetParams(ctx, *genState.FeeParams); err != nil {
		panic(err)
	}

	// set ibc port binding
	if err := k.SetPort(ctx, types.ResourcePortID); err != nil {
		panic(err)
	}

	// Bind Port claims the capability over the ResourcePortID
	if !k.IsBound(ctx, types.ResourcePortID) {
		err := k.BindPort(ctx, types.ResourcePortID)
		if err != nil {
			panic(fmt.Sprintf("could not claim port capability: %v", err))
		}
	}
}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx context.Context, k keeper.Keeper) *types.GenesisState {
	// Get all resource
	resourceList, err := k.GetAllResources(ctx)
	if err != nil {
		panic(fmt.Sprintf("Cannot get all resource: %s", err.Error()))
	}

	// get fee params
	feeParams, err := k.GetParams(ctx)
	if err != nil {
		panic(fmt.Sprintf("Cannot get fee params: %s", err.Error()))
	}

	return &types.GenesisState{
		Resources: resourceList,
		FeeParams: &feeParams,
	}
}
