package resource

import (
	"fmt"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	resourcev2 "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"

	"github.com/cosmos/cosmos-sdk/version"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: resourcev2.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Resource",
					Use:       "resource [collection-id] [id]",
					Short:     "Fetch data/payload for a specific resource (with metadata)",
					Example:   fmt.Sprintf("%s query resource resource c82f2b02-bdab-4dd7-b833-3e143745d612 93f2573c-eca9-4098-96cb-a1ec676a29ed", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "collection_id"},
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "ResourceMetadata",
					Use:       "resource-metadata [collection-id] [id]",
					Short:     "Fetch only metadata for a specific resource",
					Example:   fmt.Sprintf("%s query resource resource-metadata c82f2b02-bdab-4dd7-b833-3e143745d612 93f2573c-eca9-4098-96cb-a1ec676a29ed", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "collection_id"},
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "CollectionResources",
					Use:       "collection-resources [collection-id]",
					Short:     "Fetch metadata for all resources in a collection",
					Example:   fmt.Sprintf("%s query resource collection-resources c82f2b02-bdab-4dd7-b833-3e143745d612", version.AppName),
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "collection_id"},
					},
				},
			},
			EnhanceCustomCommand: true, // Set to true if we have manual commands for the resource module
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: resourcev2.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "CreateResource",
					Use:       "create-resource",
					Short:     "Create a new resource",
					Long:      fmt.Sprintf("Create a new resource with the specified parameters. Resources are associated with DID document collections and have various metadata fields. Example:\n%s tx resource create-resource --collection-id=c82f2b02-bdab-4dd7-b833-3e143745d612 --id=93f2573c-eca9-4098-96cb-a1ec676a29ed --name=\"My Resource\" --resource-type=\"AnonCredsSchema\" --from mykey", version.AppName),
					FlagOptions: map[string]*autocliv1.FlagOptions{
						"payload.collection_id": {
							Name:  "collection-id",
							Usage: "Identifier of the DID Document collection the resource belongs to",
						},
						"payload.id": {
							Name:  "id",
							Usage: "Unique id of the resource (UUID format)",
						},
						"payload.name": {
							Name:  "name",
							Usage: "Human-readable name of the resource",
						},
						"payload.version": {
							Name:  "version",
							Usage: "Version of the resource (e.g., 1.0.0)",
						},
						"payload.resource_type": {
							Name:  "resource-type",
							Usage: "Type of the resource (e.g., AnonCredsSchema, StatusList2021)",
						},
						"payload.data": {
							Name:  "data",
							Usage: "Base64-encoded data representing the actual content to store",
						},
					},
				},
				{
					RpcMethod: "UpdateParams",
					Use:       "update-params",
					Short:     "Update the module parameters",
					Long:      "Update the resource module parameters. Can only be executed by the module authority (governance).",
					Skip:      true, // Skip because this is authority gated
				},
			},
			EnhanceCustomCommand: true,
		},
	}
}
