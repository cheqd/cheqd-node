package simapp

// import (
// 	"github.com/cosmos/cosmos-sdk/std"
// 	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
// )

// // MakeTestEncodingConfig creates an EncodingConfig for testing. This function
// // should be used only in tests or when creating a new app instance (NewApp*()).
// // App user shouldn't create new codecs - use the app.AppCodec instead.
// // [DEPRECATED]
// func MakeTestEncodingConfig() moduletestutil.TestEncodingConfig {
// 	encodingConfig := moduletestutil.MakeTestEncodingConfig()
// 	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
// 	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
// 	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
// 	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
// 	return encodingConfig
// }
