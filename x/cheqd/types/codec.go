package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateCred_def{}, "cheqd/CreateCred_def", nil)
	cdc.RegisterConcrete(&MsgUpdateCred_def{}, "cheqd/UpdateCred_def", nil)
	cdc.RegisterConcrete(&MsgDeleteCred_def{}, "cheqd/DeleteCred_def", nil)

	cdc.RegisterConcrete(&MsgCreateSchema{}, "cheqd/CreateSchema", nil)
	cdc.RegisterConcrete(&MsgUpdateSchema{}, "cheqd/UpdateSchema", nil)
	cdc.RegisterConcrete(&MsgDeleteSchema{}, "cheqd/DeleteSchema", nil)

	cdc.RegisterConcrete(&MsgCreateAttrib{}, "cheqd/CreateAttrib", nil)
	cdc.RegisterConcrete(&MsgUpdateAttrib{}, "cheqd/UpdateAttrib", nil)
	cdc.RegisterConcrete(&MsgDeleteAttrib{}, "cheqd/DeleteAttrib", nil)

	cdc.RegisterConcrete(&MsgCreateDid{}, "cheqd/CreateDid", nil)
	cdc.RegisterConcrete(&MsgUpdateDid{}, "cheqd/UpdateDid", nil)
	cdc.RegisterConcrete(&MsgDeleteDid{}, "cheqd/DeleteDid", nil)

	cdc.RegisterConcrete(&MsgCreateNym{}, "cheqd/CreateNym", nil)
	cdc.RegisterConcrete(&MsgUpdateNym{}, "cheqd/UpdateNym", nil)
	cdc.RegisterConcrete(&MsgDeleteNym{}, "cheqd/DeleteNym", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCred_def{},
		&MsgUpdateCred_def{},
		&MsgDeleteCred_def{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateSchema{},
		&MsgUpdateSchema{},
		&MsgDeleteSchema{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAttrib{},
		&MsgUpdateAttrib{},
		&MsgDeleteAttrib{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateDid{},
		&MsgUpdateDid{},
		&MsgDeleteDid{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateNym{},
		&MsgUpdateNym{},
		&MsgDeleteNym{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
