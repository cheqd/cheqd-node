package helpers

import (
	"github.com/cheqd/cheqd-node/app/params"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	oracletypes "github.com/cheqd/cheqd-node/x/oracle/types"
	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	param "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
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
	govtypesv1.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	govtypesv1beta1.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	param.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	oracletypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	Codec = encodingConfig.Codec
	Registry = encodingConfig.InterfaceRegistry
}
