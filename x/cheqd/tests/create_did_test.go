package tests

import (
	"crypto/ed25519"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
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
				"did:cheqd:test:alice#key-1": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:alice#key-1"},
			msg: &types.MsgCreateDid{
				Id:             "did:cheqd:test:alice",
				Authentication: []string{"did:cheqd:test:alice#key-1"},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:alice#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:alice",
					},
				},
			},
		},
		{
			keys: map[string]KeyPair{
				"did:cheqd:test:bob#key-1": GenerateKeyPair(),
				"did:cheqd:test:bob#key-2": GenerateKeyPair(),
				"did:cheqd:test:bob#key-3": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:bob#key-2"},
			msg: &types.MsgCreateDid{
				Id: "did:cheqd:test:bob",
				Authentication: []string{
					"did:cheqd:test:bob#key-1",
					"did:cheqd:test:bob#key-2",
					"did:cheqd:test:bob#key-3",
				},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:bob#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:bob",
					},
					{
						Id:         "did:cheqd:test:bob#key-2",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:bob",
					},
					{
						Id:         "did:cheqd:test:bob#key-3",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:bob",
					},
				},
			},
		},
		{
			keys: map[string]KeyPair{
				"did:cheqd:test:charlie#key-1": GenerateKeyPair(),
				"did:cheqd:test:charlie#key-2": GenerateKeyPair(),
				"did:cheqd:test:charlie#key-3": GenerateKeyPair(),
			},
			signers: []string{"did:cheqd:test:charlie#key-2"},
			msg: &types.MsgCreateDid{
				Id: "did:cheqd:test:charlie",
				Authentication: []string{
					"did:cheqd:test:charlie#key-1",
					"did:cheqd:test:charlie#key-2",
					"did:cheqd:test:charlie#key-3",
				},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:         "did:cheqd:test:charlie#key-1",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:bob",
					},
					{
						Id:         "did:cheqd:test:charlie#key-2",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:bob",
					},
					{
						Id:         "did:cheqd:test:charlie#key-3",
						Type:       "Ed25519VerificationKey2020",
						Controller: "did:cheqd:test:bob",
					},
				},
			},
		},
	}

	setup := Setup()

	var prefilledPrivateKeys = map[string]map[string]ed25519.PrivateKey{}

	for _, prefilled := range prefilledDids {
		msg := prefilled.msg

		for i, key := range msg.Authentication {
			msg.VerificationMethod[i].PublicKeyMultibase = "z" + base58.Encode(prefilled.keys[key].PublicKey)
		}

		signerKeys := map[string]ed25519.PrivateKey{}
		for _, signer := range prefilled.signers {
			signerKeys[signer] = prefilled.keys[signer].PrivateKey
		}

		prefilledPrivateKeys[msg.Id] = signerKeys
		_, _ = setup.SendCreateDid(msg, signerKeys)
	}

	/*
		cases := []struct {
			valid  bool
			name string
			keys map[string]KeyPair
			signers []string
			msg *types.MsgCreateDid
			errMsg string
		}{
			{
				valid: true,
				name: "Create DID works",
				keys: map[string]KeyPair{
					"did:cheqd:test:123456qwertyui#key-1" : GenerateKeyPair(),
				},
				signers: []string{"did:cheqd:test:123456qwertyui#key-1"},
				msg: &types.MsgCreateDid{
					Id: "did:cheqd:test:123456qwertyui",
					Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
					VerificationMethod: []*types.VerificationMethod {
						{
							Id: "did:cheqd:test:123456qwertyui#key-1",
							Type: "Ed25519VerificationKey2020",
							Controller: "did:cheqd:test:123456qwertyui",
						},
					},
				},
			},
			{
				valid: true,
				name: "Create DID with controller works",
				msg: &types.MsgCreateDid{
					Id: "did:cheqd:test:controller1",
					Controller: []string {"did:cheqd:test:alice", "did:cheqd:test:bob"},
				},
				signers: []string {"did:cheqd:test:alice#key-1", "did:cheqd:test:bob"},
				keys: map {prefilledPrivateKeys[],
			},
			{
				valid: false,
				name: "DID Doc already exists",
				keys: map[string]KeyPair{
					"did:cheqd:test:123456qwertyui#key-1" : GenerateKeyPair(),
				},
				signers: []string{"did:cheqd:test:123456qwertyui#key-1"},
				msg: &types.MsgCreateDid{
					Id: "did:cheqd:test:123456qwertyui",
					Authentication: []string{"did:cheqd:test:123456qwertyui#key-1"},
					VerificationMethod: []*types.VerificationMethod {
						{
							Id: "did:cheqd:test:123456qwertyui#key-1",
							Type: "Ed25519VerificationKey2020",
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
		}*/
}
