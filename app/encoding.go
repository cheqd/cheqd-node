package app

import (
	"github.com/cheqd/cheqd-node/app/params"
	"github.com/cosmos/cosmos-sdk/std"
)

// MakeTestEncodingConfig creates an EncodingConfig for testing. This function
// should be used only in tests or when creating a new app instance (NewApp*()).
// App user shouldn't create new codecs - use the app.AppCodec instead.
// [DEPRECATED]
func MakeTestEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	return encodingConfig
}
