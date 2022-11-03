package setup

import (
	"crypto/ed25519"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
)

func (s *TestSetup) CreateDid(payload *types.MsgCreateDidDocPayload, signInputs []SignInput) (*types.MsgCreateDidDocResponse, error) {
	signBytes := payload.GetSignBytes()
	var signatures []*types.SignInfo

	for _, input := range signInputs {
		signature := ed25519.Sign(input.Key, signBytes)

		signatures = append(signatures, &types.SignInfo{
			VerificationMethodId: input.VerificationMethodId,
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
	_, _, collectionId := utils.MustSplitDID(did)

	keyPair := GenerateKeyPair()
	keyId := did + "#key-1"

	msg := &types.MsgCreateDidDocPayload{
		Id: did,
		VerificationMethod: []*types.VerificationMethod{
			{
				Id:                   keyId,
				Type:                 types.Ed25519VerificationKey2020{}.Type(),
				Controller:           did,
				VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(keyPair.Public),
			},
		},
		Authentication: []string{keyId},
	}

	signInput := SignInput{
		VerificationMethodId: keyId,
		Key:                  keyPair.Private,
	}

	return DidDocInfo{
		Did:          did,
		CollectionId: collectionId,
		KeyPair:      keyPair,
		KeyId:        keyId,
		Msg:          msg,
		SignInput:    signInput,
	}
}

func (s *TestSetup) BuildDidDocWithCustomId(uuid string) DidDocInfo {
	did := "did:cheqd:" + DID_NAMESPACE + ":" + uuid
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
		VersionId:  created.Value.Metadata.VersionId,
	}
}

func (s *TestSetup) CreateSimpleDid() CreatedDidDocInfo {
	did := s.BuildSimpleDidDoc()
	return s.CreateCustomDidDoc(did)
}

func (s *TestSetup) CreateDidDocWithExternalConterllers(controllers []string, signInputs []SignInput) CreatedDidDocInfo {
	did := s.BuildSimpleDidDoc()
	did.Msg.Controller = append(did.Msg.Controller, controllers...)

	created, err := s.CreateDid(did.Msg, append(signInputs, did.SignInput))
	if err != nil {
		panic(err)
	}

	return CreatedDidDocInfo{
		DidDocInfo: did,
		VersionId:  created.Value.Metadata.VersionId,
	}
}
