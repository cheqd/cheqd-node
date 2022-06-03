package types

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
	//TODO: implementation
	return nil
}
