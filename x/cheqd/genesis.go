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
	// Set all the credDef
	//for _, elem := range genState.CredDefList {
	//	k.SetCredDef(ctx, *elem)
	//}
	//
	//// Set credDef count
	//k.SetCredDefCount(ctx, uint64(len(genState.CredDefList)))
	//
	//// Set all the schema
	//for _, elem := range genState.SchemaList {
	//	k.SetSchema(ctx, *elem)
	//}
	//
	//// Set schema count
	//k.SetSchemaCount(ctx, uint64(len(genState.SchemaList)))
	//
	//// Set all the attrib
	//for _, elem := range genState.AttribList {
	//	k.SetAttrib(ctx, *elem)
	//}
	//
	//// Set attrib count
	//k.SetAttribCount(ctx, uint64(len(genState.AttribList)))
	//
	//// Set all the did
	//for _, elem := range genState.DidList {
	//	k.SetDid(ctx, *elem)
	//}
	//
	//// Set did count
	//k.SetDidCount(ctx, uint64(len(genState.DidList)))
	//
	//// Set all the nym
	//for _, elem := range genState.NymList {
	//	k.SetNym(ctx, *elem)
	//}
	//
	//// Set nym count
	//k.SetNymCount(ctx, uint64(len(genState.NymList)))

}

// ExportGenesis returns the cheqd module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all credDef
	//credDefList := k.GetAllCredDef(ctx)
	//for _, elem := range credDefList {
	//	elem := elem
	//	genesis.CredDefList = append(genesis.CredDefList, &elem)
	//}
	//
	//// Get all schema
	//schemaList := k.GetAllSchema(ctx)
	//for _, elem := range schemaList {
	//	elem := elem
	//	genesis.SchemaList = append(genesis.SchemaList, &elem)
	//}
	//
	//// Get all attrib
	//attribList := k.GetAllAttrib(ctx)
	//for _, elem := range attribList {
	//	elem := elem
	//	genesis.AttribList = append(genesis.AttribList, &elem)
	//}
	//
	//// Get all did
	//didList := k.GetAllDid(ctx)
	//for _, elem := range didList {
	//	elem := elem
	//	genesis.DidList = append(genesis.DidList, &elem)
	//}
	//
	//// Get all nym
	//nymList := k.GetAllNym(ctx)
	//for _, elem := range nymList {
	//	elem := elem
	//	genesis.NymList = append(genesis.NymList, &elem)
	//}

	return genesis
}
