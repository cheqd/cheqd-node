package types

var _ IdentityMsg = &MsgCreateSchema{}

func NewMsgCreateSchema(id string, typeSchema string, name string, version string, attrNames []string, controller []string) *MsgCreateSchema {
	return &MsgCreateSchema{
		Id:         id,
		Type:       typeSchema,
		Name:       name,
		Version:    version,
		AttrNames:  attrNames,
		Controller: controller,
	}
}

func (msg *MsgCreateSchema) GetSigners() []Signer {
	result := make([]Signer, len(msg.Controller))

	for i, signer := range msg.Controller {
		result[i] = Signer{
			Signer: signer,
		}
	}

	return result
}

func (msg *MsgCreateSchema) ValidateBasic() error {
	return nil
}
