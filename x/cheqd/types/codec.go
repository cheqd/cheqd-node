package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// Sdk messages
	cdc.RegisterConcrete(&MsgCreateDidDoc{}, "did/CreateDid", nil)
	cdc.RegisterConcrete(&MsgUpdateDidDoc{}, "did/UpdateDid", nil)
	cdc.RegisterConcrete(&MsgDeactivateDid{}, "did/DeleteDid", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Sdk messages
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateDidDoc{},
		&MsgUpdateDidDoc{},
	)

	// State value data
	registry.RegisterInterface("StateValueData", (*StateValueData)(nil))
	registry.RegisterImplementations((*StateValueData)(nil), &Did{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
