package types

type AssertionMethodJSONUnescaped struct {
	Id                 string                                `json:"id"`
	Type               string                                `json:"type"`
	Controller         string                                `json:"controller"`
	PublicKeyBase58    *string                               `json:"publicKeyBase58,omitempty"`
	PublicKeyMultibase *string                               `json:"publicKeyMultibase,omitempty"`
	PublicKeyJwk       *string                               `json:"publicKeyJwk,omitempty"`
	Metadata           *AssertionMethodJSONUnescapedMetadata `json:"metadata,omitempty"`
}

type AssertionMethodJSONUnescapedMetadata struct {
	ParticipantId *int    `json:"participantId"`
	ParamsRef     *string `json:"paramsRef"`
	CurveType     *string `json:"curveType"`
}
