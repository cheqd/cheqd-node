package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

var _ IdentityMsg = &MsgCreateSchema{}
var MaxAttrNamesCount = 125

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

func (msg *MsgCreateSchema) GetDid() string {
	return utils.GetDidFromSchema(msg.Id)
}

func (msg *MsgCreateSchema) ValidateBasic(namespace string) error {
	if !utils.IsSchema(msg.Id) {
		return ErrBadRequest.Wrap("Id must end with resource type '?service=CL-Schema'")
	}

	if utils.IsNotDid(namespace, msg.GetDid()) {
		return ErrBadRequestIsNotDid.Wrap("Id")
	}

	if len(msg.Type) == 0 {
		return ErrBadRequestIsRequired.Wrap("Type")
	}

	if utils.IsNotSchemaType(msg.Type) {
		return ErrBadRequest.Wrapf("%s is not allowed type", msg.Type)
	}

	if len(msg.AttrNames) == 0 {
		return ErrBadRequestIsRequired.Wrap("AttrNames")
	}

	if len(msg.AttrNames) > MaxAttrNamesCount {
		return ErrBadRequest.Wrapf("AttrNames: Expected max length 125, got: %d", len(msg.AttrNames))
	}

	if len(msg.Name) == 0 {
		return ErrBadRequestIsRequired.Wrap("Name")
	}

	if len(msg.Version) == 0 {
		return ErrBadRequestIsRequired.Wrap("Version")
	}

	if len(msg.Controller) == 0 {
		return ErrBadRequestIsRequired.Wrap("Controller")
	}

	if notValid, i := utils.ArrayContainsNotDid(namespace, msg.Controller); notValid {
		return ErrBadRequestIsNotDid.Wrapf("Controller item %s", msg.Controller[i])
	}

	return nil
}
