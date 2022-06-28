package types

var _ IdentityMsg = &MsgDeactivateDidPayload{}

func (msg *MsgDeactivateDidPayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

func (msg *MsgDeactivateDidPayload) ToDid() Did {
	return Did{
		Id: msg.Id,
	}
}

// Validation

func (msg MsgDeactivateDidPayload) Validate(allowedNamespaces []string) error {
	return msg.ToDid().Validate(allowedNamespaces)
}

func ValidMsgDeactivateDidPayloadRule(allowedNamespaces []string) *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(*MsgDeactivateDidPayload)
		if !ok {
			panic("ValidMsgDeactivateDidPayloadRule must be only applied on MsgDeactivateDidPayload properties")
		}

		return casted.Validate(allowedNamespaces)
	})
}
