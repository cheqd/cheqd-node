package types

func NewDidDocWithMetadata(didDoc *DidDoc, metadata *Metadata) DidDocWithMetadata {
	return DidDocWithMetadata{DidDoc: didDoc, Metadata: metadata}
}

func (d *DidDocWithMetadata) ReplaceDids(old, new string) {
	d.DidDoc.ReplaceDids(old, new)
}
