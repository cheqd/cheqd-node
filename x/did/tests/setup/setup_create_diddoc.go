package setup

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/google/uuid"
)

func (s *TestSetup) CreateDid(payload *types.MsgCreateDidDocPayload, signInputs []SignInput) (*types.MsgCreateDidDocResponse, error) {
	signBytes := payload.GetSignBytes()
	signatures := make([]*types.SignInfo, 0, len(signInputs))

	for _, input := range signInputs {
		signature := ed25519.Sign(input.Key, signBytes)

		signatures = append(signatures, &types.SignInfo{
			VerificationMethodId: input.VerificationMethodID,
			Signature:            signature,
		})
	}

	msg := &types.MsgCreateDidDoc{
		Payload:    payload,
		Signatures: signatures,
	}

	return s.MsgServer.CreateDidDoc(s.StdCtx, msg)
}

func (s *TestSetup) BuildDidDocWithCustomDID(did string) DidDocInfo {
	_, _, collectionID := utils.MustSplitDID(did)

	keyPair := GenerateKeyPair()
	keyID := did + "#key-1"

	msg := &types.MsgCreateDidDocPayload{
		Id: did,
		VerificationMethod: []*types.VerificationMethod{
			{
				Id:                     keyID,
				VerificationMethodType: types.Ed25519VerificationKey2020Type,
				Controller:             did,
				VerificationMaterial:   BuildEd25519VerificationKey2020VerificationMaterial(keyPair.Public),
			},
		},
		Authentication: []string{keyID},
		VersionId:      uuid.NewString(),
	}

	signInput := SignInput{
		VerificationMethodID: keyID,
		Key:                  keyPair.Private,
	}

	return DidDocInfo{
		Did:          did,
		CollectionID: collectionID,
		KeyPair:      keyPair,
		KeyID:        keyID,
		Msg:          msg,
		SignInput:    signInput,
	}
}

func (s *TestSetup) BuildDidDocWithCustomID(uuid string) DidDocInfo {
	did := "did:cheqd:" + DidNamespace + ":" + uuid
	return s.BuildDidDocWithCustomDID(did)
}

func (s *TestSetup) BuildSimpleDidDoc() DidDocInfo {
	did := GenerateDID(Base58_16bytes)
	return s.BuildDidDocWithCustomDID(did)
}

func (s *TestSetup) CreateCustomDidDoc(info DidDocInfo) CreatedDidDocInfo {
	created, err := s.CreateDid(info.Msg, []SignInput{info.SignInput})
	if err != nil {
		panic(err)
	}

	return CreatedDidDocInfo{
		DidDocInfo: info,
		VersionID:  created.Value.Metadata.VersionId,
	}
}

func (s *TestSetup) CreateSimpleDid() CreatedDidDocInfo {
	did := s.BuildSimpleDidDoc()
	return s.CreateCustomDidDoc(did)
}

func (s *TestSetup) CreateDidDocWithExternalControllers(controllers []string, signInputs []SignInput) CreatedDidDocInfo {
	did := s.BuildSimpleDidDoc()
	did.Msg.Controller = append(did.Msg.Controller, controllers...)

	created, err := s.CreateDid(did.Msg, append(signInputs, did.SignInput))
	if err != nil {
		panic(err)
	}

	return CreatedDidDocInfo{
		DidDocInfo: did,
		VersionID:  created.Value.Metadata.VersionId,
	}
}
