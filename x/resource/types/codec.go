package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// Sdk messages
	cdc.RegisterConcrete(&MsgCreateResource{}, "resource/CreateResource", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "resource/MsgUpdateParams", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// Sdk messages
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateResource{},
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
