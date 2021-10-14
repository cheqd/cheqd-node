package types

var _ IdentityMsg = &MsgCreateCredDef{}

func NewMsgCreateCredDef(id string, schemaId string, tag string, signatureType string, controller []string, value *MsgCreateCredDef_ClType) *MsgCreateCredDef {
	return &MsgCreateCredDef{
		Id:            id,
		SchemaId:      schemaId,
		Tag:           tag,
		SignatureType: signatureType,
		Value:         value,
		Controller:    controller,
	}
}

func (msg *MsgCreateCredDef) GetSigners() []Signer {
	result := make([]Signer, len(msg.Controller))

	for i, signer := range msg.Controller {
		result[i] = Signer{
			Signer: signer,
		}
	}

	return result
}

func (msg *MsgCreateCredDef) ValidateBasic() error {
	return nil
}
