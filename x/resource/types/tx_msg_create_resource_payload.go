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
	return Resource{
		CollectionId: msg.CollectionId,
		Id:           msg.Id,
		Name:         msg.Name,
		ResourceType: msg.ResourceType,
		MimeType:     msg.MimeType,
		Data:         msg.Data,
		Created:      "",
		Checksum:     "",
	}
}

// Validation

func (msg MsgCreateResourcePayload) Validate() error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.CollectionId, validation.Required, cheqdTypes.IsID()),
		validation.Field(&msg.Id, validation.Required, IsUUID()),
		validation.Field(&msg.Name, validation.Required, validation.Length(1, 64)),
		// TODO: add validation for resource type
		// TODO: add validation for mime type
		validation.Field(&msg.Data, validation.Required, validation.Length(1, 1024*1024)), // 1MB
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
