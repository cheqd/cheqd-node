package types

type IdentityMsg interface {
	ValidateDynamic(namespace string) error
	GetSigners() []Signer
	GetSignBytes() []byte
}

//TODO: Get rid of this
type Signer struct {
	Signer             string
	Authentication     []string
	VerificationMethod []*VerificationMethod
}
