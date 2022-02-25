package types

var _ IdentityMsg = &MsgCreateDidPayload{}

func (msg *MsgCreateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}
