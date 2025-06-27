package cheqd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cosmossdk.io/core/appmodule"

	"github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/did/exported"
	"github.com/cheqd/cheqd-node/x/did/types"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	abci "github.com/cometbft/cometbft/abci/types"

	// "github.com/cheqd/cheqd-node/x/did/client/rest"
	"github.com/cheqd/cheqd-node/x/did/keeper"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModuleBasic  = AppModule{}
	_ module.HasABCIGenesis  = AppModule{}
	_ module.HasServices     = AppModule{}
	_ module.HasInvariants   = AppModule{}
	_ module.HasABCIEndBlock = AppModule{}

	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface for the cheqd module.
type AppModuleBasic struct {
	cdc codec.Codec
}

func NewAppModuleBasic(cdc codec.Codec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

// Name returns the cheqd module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers the module's interface types
func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns the cheqd module's default genesis state.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the cheqd module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		log.Fatal(err)
	}
}

// GetTxCmd returns the cheqd module's root tx command.
func (a AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// // GetQueryCmd returns the cheqd module's root query command.
// func (AppModuleBasic) GetQueryCmd() *cobra.Command {
// 	return cli.GetQueryCmd()
// }

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface for the cheqd module.
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
	// legacySubspace is used solely for migration of x/params managed parameters
	legacySubspace exported.Subspace
}

var _ appmodule.AppModule = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

func NewAppModule(cdc codec.Codec, keeper keeper.Keeper, subspace exported.Subspace) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		legacySubspace: subspace,
	}
}

// ConsensusVersion is a sequence number for state-breaking change of the
// module. It should be incremented on each consensus-breaking change
// introduced by the module. To avoid wrong/empty versions, the initial version
// should be set to 1.
func (am AppModule) ConsensusVersion() uint64 {
	return 6
}

// Name returns the cheqd module's name.
func (am AppModule) Name() string {
	return am.AppModuleBasic.Name()
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.keeper))
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))

	m := keeper.NewMigrator(am.keeper, am.legacySubspace)

	if err := cfg.RegisterMigration(types.ModuleName, 4, m.Migrate4to5); err != nil {
		panic(fmt.Sprintf("failed to migrate x/did from version 4 to 5: %v", err))
	}
	if err := cfg.RegisterMigration(types.ModuleName, 5, m.Migrate5to6); err != nil {
		panic(fmt.Sprintf("failed to migrate x/did from version 5 to 6: %v", err))
	}
}

// RegisterInvariants registers the cheqd module's invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the cheqd module's genesis initialization It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)
	InitGenesis(ctx, am.keeper, &genState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the cheqd module's exported genesis state as raw JSON bytes.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(genState)
}

// BeginBlock executes all ABCI BeginBlock logic respective to the cheqd module.
func (am AppModule) BeginBlock(_ context.Context) error {
	return nil
}

// EndBlock executes all ABCI EndBlock logic respective to the cheqd module. It
// returns no validator updates.
func (am AppModule) EndBlock(_ context.Context) ([]abci.ValidatorUpdate, error) {
	return []abci.ValidatorUpdate{}, nil
}
