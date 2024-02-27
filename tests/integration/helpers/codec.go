package helpers

import (
	"github.com/cheqd/cheqd-node/app/params"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
)

var (
	Codec    codec.Codec
	Registry types.InterfaceRegistry
)

func init() {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	resourcetypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	didtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	Codec = encodingConfig.Codec
	Registry = encodingConfig.InterfaceRegistry
}
