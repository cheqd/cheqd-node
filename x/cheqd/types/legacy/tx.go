package legacy

type IdentityMsg interface {
	GetSignBytes() []byte
}
