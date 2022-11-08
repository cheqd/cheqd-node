package types

func NewDidDocWithMetadata(didDoc *DidDoc, metadata *Metadata) DidDocWithMetadata {
	return DidDocWithMetadata{DidDoc: didDoc, Metadata: metadata}
}
