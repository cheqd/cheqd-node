package types

import (
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ didtypes.IdentityMsg = &MsgCreateResourcePayload{}

func (msg *MsgCreateResourcePayload) GetSignBytes() []byte {
	bytes, err := msg.Marshal()
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgCreateResourcePayload) ToResource() ResourceWithMetadata {
	return ResourceWithMetadata{
		Metadata: &Metadata{
			CollectionId: msg.CollectionId,
			Id:           msg.Id,
			Name:         msg.Name,
			Version:      msg.Version,
			ResourceType: msg.ResourceType,
			AlsoKnownAs:  msg.AlsoKnownAs,
		},
		Resource: &Resource{
			Data: msg.Data,
		},
	}
}

// Validation

func (msg MsgCreateResourcePayload) Validate() error {
	return validation.ValidateStruct(&msg,
		validation.Field(&msg.CollectionId, validation.Required, didtypes.IsID()),
		validation.Field(&msg.Id, validation.Required, didtypes.IsUUID()),
		validation.Field(&msg.Name, validation.Required, validation.Length(1, 64)),
		validation.Field(&msg.Version, validation.Length(1, 64)),
		validation.Field(&msg.ResourceType, validation.Required, validation.Length(1, 64)),
		validation.Field(&msg.AlsoKnownAs, validation.Each(ValidAlternativeURI())),
		validation.Field(&msg.Data, validation.Required, validation.Length(1, 200*1024)), // 200KB
	)
}

func ValidMsgCreateResourcePayload() *didtypes.CustomErrorRule {
	return didtypes.NewCustomErrorRule(func(value interface{}) error {
		casted, ok := value.(*MsgCreateResourcePayload)
		if !ok {
			panic("ValidMsgCreateResourcePayload must be only applied on MsgCreateDidPayload properties")
		}

		return casted.Validate()
	})
}

// Normalize

func (msg *MsgCreateResourcePayload) Normalize() {
	msg.CollectionId = didutils.NormalizeID(msg.CollectionId)
	msg.Id = didutils.NormalizeUUID(msg.Id)
}
