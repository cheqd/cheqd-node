package v1

type IdentityMsg interface {
	GetSignBytes() []byte
}
