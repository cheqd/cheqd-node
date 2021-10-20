package tests

import (
	"crypto/ed25519"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateDID(t *testing.T) {
	prefilledDids := []struct {
		keys    map[string]KeyPair
		signers []string
		msg     *types.MsgCreateDid
	}{
		{
			keys: map[string]KeyPair{
				AliceKey1: GenerateKeyPair(),
			},
			signers: []string{AliceKey1},
			msg: &types.MsgCreateDid{
				Id:             AliceDID,
				Authentication: []string{AliceKey1},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         AliceKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: AliceDID,
					},
				},
			},
		},
		{
			keys: map[string]KeyPair{
				BobKey1: GenerateKeyPair(),
				BobKey2: GenerateKeyPair(),
				BobKey3: GenerateKeyPair(),
			},
			signers: []string{BobKey2},
			msg: &types.MsgCreateDid{
				Id: BobDID,
				Authentication: []string{
					BobKey1,
					BobKey2,
					BobKey3,
				},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         BobKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         BobKey3,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
		{
			keys: map[string]KeyPair{
				CharlieKey1: GenerateKeyPair(),
				CharlieKey2: GenerateKeyPair(),
				CharlieKey3: GenerateKeyPair(),
			},
			signers: []string{CharlieKey2},
			msg: &types.MsgCreateDid{
				Id: CharlieDID,
				Authentication: []string{
					CharlieKey1,
					CharlieKey2,
					CharlieKey3,
				},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         CharlieKey1,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         CharlieKey2,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
					{
						Id:         CharlieKey3,
						Type:       "Ed25519VerificationKey2020",
						Controller: BobDID,
					},
				},
			},
		},
	}

	setup := Setup()
	keys := map[string]KeyPair{}

	for _, prefilled := range prefilledDids {
		msg := prefilled.msg

		for i, key := range msg.Authentication {
			msg.VerificationMethod[i].PublicKeyMultibase = "z" + base58.Encode(prefilled.keys[key].PublicKey)
		}

		signerKeys := map[string]ed25519.PrivateKey{}
		for _, signer := range prefilled.signers {
			signerKeys[signer] = prefilled.keys[signer].PrivateKey
			keys[signer] = prefilled.keys[signer]
		}

		_, _ = setup.SendCreateDid(msg, signerKeys)
	}

	cases := []struct {
		valid   bool
		name    string
		keys    map[string]KeyPair
		signers []string
		msg     *types.MsgCreateDid
		errMsg  string
	}{
		{
			valid: true,
			name:  "Create DID works",
			keys: map[string]KeyPair{
				"did:cheqd:test:123456qwertyui#key-1": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:123456qwertyui#key-1"},
			msg: &types.MsgCreateDid{
				Id:             "did:cheqd:test:123456qwertyui",
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
				},
			},
		},
		{
			valid: true,
			name:  "Create DID with controller works",
			msg: &types.MsgCreateDid{
				Id:         "did:cheqd:test:controller1",
				Controller: []string{AliceDID, BobDID},
			},
			signers: []string{AliceKey1, BobKey2},
		},
		{
			valid: false,
			name:  "DID Doc already exists",
			keys: map[string]KeyPair{
				"did:cheqd:test:123456qwertyui#key-1": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:123456qwertyui#key-1"},
			msg: &types.MsgCreateDid{
				Id:             "did:cheqd:test:123456qwertyui",
				Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:123456qwertyui#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:123456qwertyui",
					},
				},
			},
			errMsg: "DID DOC already exists for DID did:cheqd:test:123456qwertyui: DID Doc exists",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.msg

			for i, key := range msg.Authentication {
				msg.VerificationMethod[i].PublicKeyMultibase = "z" + base58.Encode(tc.keys[key].PublicKey)
			}

			signerKeys := map[string]ed25519.PrivateKey{}
			for _, signer := range tc.signers {
				signerKeys[signer] = tc.keys[signer].PrivateKey
			}

			did, err := setup.SendCreateDid(msg, signerKeys)

			if tc.valid {
				require.Nil(t, err)
				require.Equal(t, did.Id, did.Id)
				require.Equal(t, did.Controller, did.Controller)
				require.Equal(t, did.VerificationMethod, did.VerificationMethod)
				require.Equal(t, did.Authentication, did.Authentication)
				require.Equal(t, did.AssertionMethod, did.AssertionMethod)
				require.Equal(t, did.CapabilityInvocation, did.CapabilityInvocation)
				require.Equal(t, did.CapabilityDelegation, did.CapabilityDelegation)
				require.Equal(t, did.KeyAgreement, did.KeyAgreement)
				require.Equal(t, did.AlsoKnownAs, did.AlsoKnownAs)
				require.Equal(t, did.Service, did.Service)
				require.Equal(t, did.Context, did.Context)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.errMsg, err.Error())
			}
		})
	}
}
