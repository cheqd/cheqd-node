package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// Sdk messages
	cdc.RegisterConcrete(&MsgCreateDidDoc{}, "did/CreateDidDoc", nil)
	cdc.RegisterConcrete(&MsgUpdateDidDoc{}, "did/UpdateDidDoc", nil)
	cdc.RegisterConcrete(&MsgDeactivateDidDoc{}, "did/DeleteDidDoc", nil)
	cdc.RegisterConcrete(&MsgBurn{}, "did/MsgBurn", nil)
	cdc.RegisterConcrete(&MsgMint{}, "did/MsgMint", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "did/MsgUpdateParams", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Sdk messages
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateDidDoc{},
		&MsgUpdateDidDoc{},
		&MsgDeactivateDidDoc{},
		&MsgBurn{},
		&MsgMint{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
