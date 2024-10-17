package types

type AssertionMethodJSONUnescaped struct {
	Id                 string  `json:"id"`
	Type               string  `json:"type"`
	Controller         string  `json:"controller"`
	PublicKeyBase58    *string `json:"publicKeyBase58,omitempty"`
	PublicKeyMultibase *string `json:"publicKeyMultibase,omitempty"`
	PublicKeyJwk       *string `json:"publicKeyJwk,omitempty"`
}
