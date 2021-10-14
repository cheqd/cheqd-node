package types

import "github.com/cheqd/cheqd-node/x/cheqd/utils"

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
	if utils.IsNotDid(msg.Id) {
		return ErrBadRequestIsNotDid.Wrap("Id")
	}

	if len(msg.Type) == 0 {
		return ErrBadRequestIsRequired.Wrap("Type")
	}

	if len(msg.Type) == 0 {
		return ErrBadRequestIsRequired.Wrap("Type")
	}

	if len(msg.Type) == 0 {
		return ErrBadRequestIsRequired.Wrap("Type")
	}

	if len(msg.Type) == 0 {
		return ErrBadRequestIsRequired.Wrap("Type")
	}

	if valid, i := utils.ArrayContainsNotDid(msg.Controller); !valid {
		return ErrBadRequestIsNotDid.Wrapf("Controller item %s", msg.Controller[i])
	}

	return nil
}
