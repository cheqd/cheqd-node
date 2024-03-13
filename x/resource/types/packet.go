package types

// the data in the ibc packet
// this is from a CosmWasm contract `IBCSend` msg
// This is in the `Binary` format, which is a Rust Vec<u8> type of a base64 encoded string of a JSON obj
type ResourceReqPacket struct {
	// Id of the resource
	ResourceId string `json:"resourceId,omitempty"`
	// Id of the collection
	CollectionId string `json:"collectionId,omitempty"`
}
