package types

type (
	IdentityMsg interface {
		ValidateBasic(namespace string) error
		GetSigners() []Signer
	}

	Signer struct {
		Signer             string
		Authentication     []string
		VerificationMethod []*VerificationMethod
	}
)
