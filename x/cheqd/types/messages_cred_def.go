package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

var _ IdentityMsg = &MsgCreateCredDef{}

func NewMsgCreateCredDef(id string, schemaId string, tag string, signatureType string, controller []string, value *MsgCreateCredDef_ClType) *MsgCreateCredDef {
	return &MsgCreateCredDef{
		Id:         id,
		SchemaId:   schemaId,
		Tag:        tag,
		Type:       signatureType,
		Value:      value,
		Controller: controller,
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

func (msg *MsgCreateCredDef) GetDid() string {
	return utils.GetDidFromCredDef(msg.Id)
}

func (msg *MsgCreateCredDef) ValidateBasic() error {
	if !utils.IsCredDef(msg.Id) {
		return ErrBadRequest.Wrap("Id must end with resource type '/credDef'")
	}

	if utils.IsNotDid(msg.GetDid()) {
		return ErrBadRequestIsNotDid.Wrap("Id")
	}

	if msg.Value == nil || msg.Value.Size() == 0 {
		return ErrBadRequestIsRequired.Wrap("Value")
	}

	if len(msg.SchemaId) == 0 {
		return ErrBadRequestIsRequired.Wrap("SchemaId")
	}

	if len(msg.Type) == 0 {
		return ErrBadRequestIsRequired.Wrap("SignatureType")
	}

	if utils.IsNotCredDefSignatureType(msg.Type) {
		return ErrBadRequest.Wrapf("%s is not allowed type", msg.Type)
	}

	if len(msg.Controller) == 0 {
		return ErrBadRequestIsRequired.Wrap("Controller")
	}

	if notValid, i := utils.ArrayContainsNotDid(msg.Controller); notValid {
		return ErrBadRequestIsNotDid.Wrapf("Controller item %s", msg.Controller[i])
	}

	return nil
}
