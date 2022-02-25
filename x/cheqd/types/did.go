package types

import "github.com/multiformats/go-multibase"

var _ StateValueData = &Did{}

func NewVerificationMethod(id string, type_ string, controller string, publicKeyJwk []*KeyValuePair, publicKeyMultibase string) *VerificationMethod {
	return &VerificationMethod{
		Id:                 id,
		Type:               type_,
		Controller:         controller,
		PublicKeyJwk:       publicKeyJwk,
		PublicKeyMultibase: publicKeyMultibase,
	}
}

func (v VerificationMethod) GetPublicKey() ([]byte, error) {
	if len(v.PublicKeyMultibase) > 0 {
		_, key, err := multibase.Decode(v.PublicKeyMultibase)
		if err != nil {
			return nil, ErrInvalidPublicKey.Wrapf("Cannot decode verification method '%s' public key", v.Id)
		}
		return key, nil
	}

	if len(v.PublicKeyJwk) > 0 {
		return nil, ErrInvalidPublicKey.Wrap("JWK format not supported")
	}

	return nil, ErrInvalidPublicKey.Wrapf("verification method '%s' public key not found", v.Id)
}
