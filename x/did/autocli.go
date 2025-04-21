package cheqd

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	didv2 "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: didv2.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "DidDoc",
					Use:       "did-document [id]",
					Short:     "Query a DID Document by DID",
					Long:      "Fetch latest version of a DID Document for a given DID",
					Example:   "", //TODO: add the example,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "DidDocVersion",
					Use:       "did-version [id] [version-id]",
					Short:     "Query specific version of a DID Document",
					Long:      " Fetch specific version of a DID Document for a given DID",
					Example:   "", //TODO: add the example,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
						{ProtoField: "version"},
					},
				},
				{
					RpcMethod: "AllDidDocVersionsMetadata",
					Use:       "did-metadata [id]",
					Short:     "Query all versions metadata for a DID",
					Long:      "Fetch list of all versions of DID Documents for a given DID",
					Example:   "", //TODO: add the example,
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current did parameters",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: didv2.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Burn",
					Use:       "burn [amount] [flags]",
					Short:     "Burn tokens from an address",
					Long:      "",
					Example:   "",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "Mint",
					Skip:      true, // skipped because authority gated
				},
			},
		},
	}
}
