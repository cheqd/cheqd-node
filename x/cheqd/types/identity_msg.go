package types

type IdentityMsg interface {
	GetSignBytes() []byte
}
