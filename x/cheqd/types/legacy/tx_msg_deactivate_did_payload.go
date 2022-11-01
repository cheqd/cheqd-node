package legacy

import validation "github.com/go-ozzo/ozzo-validation/v4"

var _ IdentityMsg = &MsgDeactivateDidPayload{}

func (msg *MsgDeactivateDidPayload) GetSignBytes() []byte {
	bytes, err := msg.Marshal()
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgDeactivateDidPayload) ToDid() Did {
	return Did{
		Id: msg.Id,
	}
}

// Validation

func (msg MsgDeactivateDidPayload) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.Id, validation.Required, IsDID(allowedNamespaces)),
	)
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
