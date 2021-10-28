package types

type (
	IdentityMsg interface {
		Validate(namespace string) error
		GetSigners() []Signer
	}

	Signer struct {
		Signer             string
		Authentication     []string
		VerificationMethod []*VerificationMethod
	}
)
