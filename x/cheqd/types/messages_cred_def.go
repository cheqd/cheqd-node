package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	"strings"
)

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
	if !strings.HasSuffix(msg.Id, "/credDef") {
		return ErrBadRequest.Wrap("Id must end with resource type '/credDef'")
	}

	if utils.IsNotDid(msg.Id[:len(msg.Id)-8]) {
		return ErrBadRequestIsNotDid.Wrap("Id")
	}

	if msg.Value == nil || msg.Value.Size() == 0 {
		return ErrBadRequestIsRequired.Wrap("Value")
	}

	if len(msg.SchemaId) == 0 {
		return ErrBadRequestIsRequired.Wrap("SchemaId")
	}

	if len(msg.SignatureType) == 0 {
		return ErrBadRequestIsRequired.Wrap("SignatureType")
	}

	if utils.IsNotCredDefSignatureType(msg.SignatureType) {
		return ErrBadRequest.Wrapf("%s is not allowed signature type", msg.SignatureType)
	}

	if len(msg.Controller) == 0 {
		return ErrBadRequestIsRequired.Wrap("Controller")
	}

	if notValid, i := utils.ArrayContainsNotDid(msg.Controller); notValid {
		return ErrBadRequestIsNotDid.Wrapf("Controller item %s", msg.Controller[i])
	}

	return nil
}
