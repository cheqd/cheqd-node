package types

type (
	IdentityMsg interface {
		ValidateBasic() error
		GetSigners() []Signer
	}

	Signer struct {
		Signer             string
		Authentication     []string
		VerificationMethod []*VerificationMethod
	}
)
