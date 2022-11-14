package types

import (
	cheqdtypes "github.com/canow-co/cheqd-node/x/cheqd/types"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cheqdtypes.IdentityMsg = &MsgCreateResourcePayload{}

func (msg *MsgCreateResourcePayload) GetSignBytes() []byte {
	return ModuleCdc.MustMarshal(msg)
}

func (msg *MsgCreateResourcePayload) ToResource() Resource {
	return Resource{
		Header: &ResourceHeader{
			CollectionId: msg.CollectionId,
			Id:           msg.Id,
			Name:         msg.Name,
			ResourceType: msg.ResourceType,
		},
		Data: msg.Data,
	}
}

// Validation

func (msg MsgCreateResourcePayload) Validate() error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.CollectionId, validation.Required, cheqdtypes.IsID()),
		validation.Field(&msg.Id, validation.Required, cheqdtypes.IsUUID()),
		validation.Field(&msg.Name, validation.Required, validation.Length(1, 64)),
		validation.Field(&msg.ResourceType, validation.Required, validation.Length(1, 64)),
		validation.Field(&msg.Data, validation.Required, validation.Length(1, 200*1024)), // 200KB
	)
}

func ValidMsgCreateResourcePayload() *cheqdtypes.CustomErrorRule {
	return cheqdtypes.NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(*MsgCreateResourcePayload)
		if !ok {
			panic("ValidMsgCreateResourcePayload must be only applied on MsgCreateDidPayload properties")
		}

		return casted.Validate()
	})
}
