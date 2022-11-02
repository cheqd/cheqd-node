package types

import validation "github.com/go-ozzo/ozzo-validation/v4"

var _ IdentityMsg = &MsgDeactivateDidDocPayload{}

func (msg *MsgDeactivateDidDocPayload) GetSignBytes() []byte {
	bytes, err := msg.Marshal()
	if err != nil {
		panic(err)
	}

	return bytes
}

// Validation

func (msg MsgDeactivateDidDocPayload) Validate(allowedNamespaces []string) error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.Id, validation.Required, IsDID(allowedNamespaces)),
	)
}

func ValidMsgDeactivateDidPayloadRule(allowedNamespaces []string) *CustomErrorRule {
	return NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(*MsgDeactivateDidDocPayload)
		if !ok {
			panic("ValidMsgDeactivateDidPayloadRule must be only applied on MsgDeactivateDidPayload properties")
		}

		return casted.Validate(allowedNamespaces)
	})
}
