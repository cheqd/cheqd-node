package resource

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	resourcev2 "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"

	"github.com/cosmos/cosmos-sdk/version"
)
a

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: resourcev2.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Resource",
					Use:       "specific-resource [collection-id] [resource-id]",
					Short:     "Query a specific resource",
					Example:   fmt.Sprintf("%s query resource pecific-resource c82f2b02-bdab-4dd7-b833-3e143745d612 wGHEXrZvJxR8vw5P3UWH1j", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "collection_id"},
						{ProtoField: "resource-id"},
					},
				},
				{
					RpcMethod: "CollectionMetadata",
					Use:       "collection-metadata [collection-id]",
					Short:     "Query metadata for an entire Collection",
					Example:   fmt.Sprintf("%s query resource collection-metadata c82f2b02-bdab-4dd7-b833-3e143745d612 wGHEXrZvJxR8vw5P3UWH1j", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "collection_id"},
					},
				},
				{
					RpcMethod: "ResourceMetadata",
					Use:       "metadata [collection-id] [resource-id]",
					Short:     "Query metadata for a specific resource",
					Example:   fmt.Sprintf("%s query resource metadata c82f2b02-bdab-4dd7-b833-3e143745d612 wGHEXrZvJxR8vw5P3UWH1j", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "collection_id"},
						{ProtoField: "resource-id"},
					},
				},
			},
			EnhanceCustomCommand: true, // Set to true if we have manual commands for the resource module
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: resourcev2.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
			},
		},
	}
}
