package types

import (
	cheqdTypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cheqdTypes.IdentityMsg = &MsgCreateResourcePayload{}

func (msg *MsgCreateResourcePayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

func (msg *MsgCreateResourcePayload) ToResource() Resource {
	created := ""
	checksum := ""
	return Resource{
		CollectionId: msg.CollectionId,
		Id:           msg.Id,
		Name:         msg.Name,
		ResourceType: msg.ResourceType,
		MimeType:     msg.MimeType,
		Data:         msg.Data,
		Created:      created,
		Checksum:     checksum,
	}
}

// Validation

func (msg MsgCreateResourcePayload) Validate() error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.CollectionId, validation.Required, IsUUID()),	// TODO: Wrong
	)
}

func ValidMsgCreateResourcePayload() *cheqdTypes.CustomErrorRule {
	return cheqdTypes.NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(*MsgCreateResourcePayload)
		if !ok {
			panic("ValidMsgCreateResourcePayload must be only applied on MsgCreateDidPayload properties")
		}

		return casted.Validate()
	})
}
