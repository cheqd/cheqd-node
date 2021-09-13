package cheqd

import (
	"github.com/cheqd/cheqd-node/x/cheqd/keeper"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the cred_def
	for _, elem := range genState.Cred_defList {
		k.SetCred_def(ctx, *elem)
	}

	// Set cred_def count
	k.SetCred_defCount(ctx, uint64(len(genState.Cred_defList)))

	// Set all the schema
	for _, elem := range genState.SchemaList {
		k.SetSchema(ctx, *elem)
	}

	// Set schema count
	k.SetSchemaCount(ctx, uint64(len(genState.SchemaList)))

	// Set all the attrib
	for _, elem := range genState.AttribList {
		k.SetAttrib(ctx, *elem)
	}

	// Set attrib count
	k.SetAttribCount(ctx, uint64(len(genState.AttribList)))

	// Set all the did
	for _, elem := range genState.DidList {
		k.SetDid(ctx, *elem)
	}

	// Set did count
	k.SetDidCount(ctx, uint64(len(genState.DidList)))

	// Set all the nym
	for _, elem := range genState.NymList {
		k.SetNym(ctx, *elem)
	}

	// Set nym count
	k.SetNymCount(ctx, uint64(len(genState.NymList)))

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all cred_def
	cred_defList := k.GetAllCred_def(ctx)
	for _, elem := range cred_defList {
		elem := elem
		genesis.Cred_defList = append(genesis.Cred_defList, &elem)
	}

	// Get all schema
	schemaList := k.GetAllSchema(ctx)
	for _, elem := range schemaList {
		elem := elem
		genesis.SchemaList = append(genesis.SchemaList, &elem)
	}

	// Get all attrib
	attribList := k.GetAllAttrib(ctx)
	for _, elem := range attribList {
		elem := elem
		genesis.AttribList = append(genesis.AttribList, &elem)
	}

	// Get all did
	didList := k.GetAllDid(ctx)
	for _, elem := range didList {
		elem := elem
		genesis.DidList = append(genesis.DidList, &elem)
	}

	// Get all nym
	nymList := k.GetAllNym(ctx)
	for _, elem := range nymList {
		elem := elem
		genesis.NymList = append(genesis.NymList, &elem)
	}

	return genesis
}
