package types

var _ IdentityMsg = &MsgUpdateDidPayload{}

func (msg *MsgUpdateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}
