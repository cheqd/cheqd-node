package tests

import (
	"crypto/ed25519"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/btcsuite/btcutil/base58"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/multiformats/go-multibase"
)

var _ = Describe("Create DID tests", func() {
	// params for cases
	var valid bool
	var keys map[string]KeyPair
	var signers []string
	var msg *types.MsgCreateDidPayload
	var errMsg = ""

	var setup TestSetup
	var err error
	var mainKeys = GenerateTestKeys()

	BeforeEach(func() {
		// setup
		valid = false
		keys = map[string]KeyPair{}
		signers = []string{}
		msg = &types.MsgCreateDidPayload{}
		errMsg = ""
	})

	AfterEach(func() {
		
		setup = InitEnv(mainKeys)

		for _, vm := range msg.VerificationMethod {
			if vm.PublicKeyMultibase == "" {
				vm.PublicKeyMultibase, err = multibase.Encode(multibase.Base58BTC, keys[vm.Id].PublicKey)
			}
			Expect(err).To(BeNil())
		}

		signerKeys := map[string]ed25519.PrivateKey{}
		for _, signer := range signers {
			signerKeys[signer] = keys[signer].PrivateKey
		}

		did, err := setup.SendCreateDid(msg, signerKeys)

		if valid {
			Expect(err).To(BeNil())
			Expect(msg.Id).To(Equal(did.Id))
			Expect(msg.Controller).To(Equal(did.Controller))
			Expect(msg.VerificationMethod).To(Equal(did.VerificationMethod))
			Expect(msg.Authentication).To(Equal(did.Authentication))
			Expect(msg.AssertionMethod).To(Equal(did.AssertionMethod))
			Expect(msg.CapabilityInvocation).To(Equal(did.CapabilityInvocation))
			Expect(msg.CapabilityDelegation).To(Equal(did.CapabilityDelegation))
			Expect(msg.KeyAgreement).To(Equal(did.KeyAgreement))
			Expect(msg.AlsoKnownAs).To(Equal(did.AlsoKnownAs))
			Expect(msg.Service).To(Equal(did.Service))
			Expect(msg.Context).To(Equal(did.Context))
		} else {
			Expect(err).To(HaveOccurred())
			Expect(errMsg).To(Equal(err.Error()))
		}
	})

	It("Valid: Works", func() {
		valid = true
		keys = map[string]KeyPair{
			ImposterKey1: GenerateKeyPair(),
		}
		signers = []string{ImposterKey1}
		msg = &types.MsgCreateDidPayload{
			Id:             ImposterDID,
			Authentication: []string{ImposterKey1},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         ImposterKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: ImposterDID,
				},
			},
		}
	})

	It("Valid: Works with Key Agreement", func() {
		valid = true
		keys = map[string]KeyPair{
			ImposterKey1: GenerateKeyPair(),
			AliceKey1:    mainKeys[AliceKey1],
		}
		signers = []string{ImposterKey1, AliceKey1}
		msg = &types.MsgCreateDidPayload{
			Id:           ImposterDID,
			KeyAgreement: []string{ImposterKey1},
			Controller:   []string{AliceDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         ImposterKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: ImposterDID,
				},
			},
		}
	})

	It("Valid: Works with Assertion Method", func() {
		valid = true
		keys = map[string]KeyPair{
			ImposterKey1: GenerateKeyPair(),
			AliceKey1:    mainKeys[AliceKey1],
		}
		signers = []string{AliceKey1, ImposterKey1}
		msg = &types.MsgCreateDidPayload{
			Id:              ImposterDID,
			AssertionMethod: []string{ImposterKey1},
			Controller:      []string{AliceDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         ImposterKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: ImposterDID,
				},
			},
		}
	})

	It("Valid: Works with Capability Delegation", func() {
		valid = true
		keys = map[string]KeyPair{
			ImposterKey1: GenerateKeyPair(),
			AliceKey1:    mainKeys[AliceKey1],
		}
		signers = []string{AliceKey1, ImposterKey1}
		msg = &types.MsgCreateDidPayload{
			Id:                   ImposterDID,
			CapabilityDelegation: []string{ImposterKey1},
			Controller:           []string{AliceDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         ImposterKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: ImposterDID,
				},
			},
		}
	})

	It("Valid: Works with Capability Invocation", func() {
		valid = true
		keys = map[string]KeyPair{
			ImposterKey1: GenerateKeyPair(),
			AliceKey1:    mainKeys[AliceKey1],
		}
		signers = []string{AliceKey1, ImposterKey1}
		msg= &types.MsgCreateDidPayload{
			Id:                   ImposterDID,
			CapabilityInvocation: []string{ImposterKey1},
			Controller:           []string{AliceDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         ImposterKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: ImposterDID,
				},
			},
		}
	})

	It("Valid: With controller works", func() {
		valid = true
		msg = &types.MsgCreateDidPayload{
			Id:         ImposterDID,
			Controller: []string{AliceDID, BobDID},
		}
		signers = []string{AliceKey1, BobKey3}
		keys = map[string]KeyPair{
			AliceKey1: mainKeys[AliceKey1],
			BobKey3:   mainKeys[BobKey3],
		}
	})

	It("Valid: Full message works", func() {
		valid = true
		keys = map[string]KeyPair{
			"did:cheqd:test:yyyyyyyyyyyyyyyy#key-1": GenerateKeyPair(),
			"did:cheqd:test:yyyyyyyyyyyyyyyy#key-2": GenerateKeyPair(),
			"did:cheqd:test:yyyyyyyyyyyyyyyy#key-3": GenerateKeyPair(),
			"did:cheqd:test:yyyyyyyyyyyyyyyy#key-4": GenerateKeyPair(),
			"did:cheqd:test:yyyyyyyyyyyyyyyy#key-5": GenerateKeyPair(),
			AliceKey1:                               mainKeys[AliceKey1],
			BobKey1:                                 mainKeys[BobKey1],
			BobKey2:                                 mainKeys[BobKey2],
			BobKey3:                                 mainKeys[BobKey3],
			CharlieKey1:                             mainKeys[CharlieKey1],
			CharlieKey2:                             mainKeys[CharlieKey2],
			CharlieKey3:                             mainKeys[CharlieKey3],
		}
		signers = []string{
			"did:cheqd:test:yyyyyyyyyyyyyyyy#key-1",
			"did:cheqd:test:yyyyyyyyyyyyyyyy#key-5",
			AliceKey1,
			BobKey1,
			BobKey2,
			BobKey3,
			CharlieKey1,
			CharlieKey2,
			CharlieKey3,
		}
		msg = &types.MsgCreateDidPayload{
			Id: "did:cheqd:test:yyyyyyyyyyyyyyyy",
			Authentication: []string{
				"did:cheqd:test:yyyyyyyyyyyyyyyy#key-1",
				"did:cheqd:test:yyyyyyyyyyyyyyyy#key-5",
			},
			Context:              []string{"abc", "de"},
			CapabilityInvocation: []string{"did:cheqd:test:yyyyyyyyyyyyyyyy#key-2"},
			CapabilityDelegation: []string{"did:cheqd:test:yyyyyyyyyyyyyyyy#key-3"},
			KeyAgreement:         []string{"did:cheqd:test:yyyyyyyyyyyyyyyy#key-4"},
			AlsoKnownAs:          []string{"SomeUri"},
			Service: []*types.Service{
				{
					Id:              "did:cheqd:test:yyyyyyyyyyyyyyyy#service-1",
					Type:            "DIDCommMessaging",
					ServiceEndpoint: "ServiceEndpoint",
				},
			},
			Controller: []string{"did:cheqd:test:yyyyyyyyyyyyyyyy", AliceDID, BobDID, CharlieDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         "did:cheqd:test:yyyyyyyyyyyyyyyy#key-1",
					Type:       Ed25519VerificationKey2020,
					Controller: "did:cheqd:test:yyyyyyyyyyyyyyyy",
				},
				{
					Id:         "did:cheqd:test:yyyyyyyyyyyyyyyy#key-2",
					Type:       Ed25519VerificationKey2020,
					Controller: "did:cheqd:test:yyyyyyyyyyyyyyyy",
				},
				{
					Id:         "did:cheqd:test:yyyyyyyyyyyyyyyy#key-3",
					Type:       Ed25519VerificationKey2020,
					Controller: "did:cheqd:test:yyyyyyyyyyyyyyyy",
				},
				{
					Id:         "did:cheqd:test:yyyyyyyyyyyyyyyy#key-4",
					Type:       "Ed25519VerificationKey2020",
					Controller: "did:cheqd:test:yyyyyyyyyyyyyyyy",
				},
				{
					Id:         "did:cheqd:test:yyyyyyyyyyyyyyyy#key-5",
					Type:       "Ed25519VerificationKey2020",
					Controller: "did:cheqd:test:yyyyyyyyyyyyyyyy",
				},
			},
		}
	})

// ************************** 
// ***** Negative cases *****
// **************************

	It("Not Valid: Second controller did not sign request", func() {
		valid = false
		msg = &types.MsgCreateDidPayload{
			Id:         ImposterDID,
			Controller: []string{AliceDID, BobDID},
		}
		signers = []string{AliceKey1}
		keys = map[string]KeyPair{
			AliceKey1: mainKeys[AliceKey1],
		}
		errMsg = fmt.Sprintf("signer: %s: signature is required but not found", BobDID)
	})

	It("Not Valid: No signature", func(){
		valid = false
		msg = &types.MsgCreateDidPayload{
			Id:         ImposterDID,
			Controller: []string{AliceDID, BobDID},
		}
		errMsg = fmt.Sprintf("signer: %s: signature is required but not found", AliceDID)
	})
	
	It("Not Valid: Controller not found", func(){
		valid = false
		msg = &types.MsgCreateDidPayload{
			Id:         ImposterDID,
			Controller: []string{AliceDID, NotFounDID},
		}
		signers = []string{AliceKey1, ImposterKey1}
		keys = map[string]KeyPair{
			AliceKey1:    mainKeys[AliceKey1],
			ImposterKey1: GenerateKeyPair(),
		}
		errMsg = fmt.Sprintf("%s: DID Doc not found", NotFounDID)	
	})

	It("Not Valid: Wrong signature", func(){
		valid = false
		msg = &types.MsgCreateDidPayload{
			Id:         ImposterDID,
			Controller: []string{AliceDID},
		}
		signers = []string{AliceKey1}
		keys = map[string]KeyPair{
			AliceKey1: mainKeys[BobKey1],
		}
		errMsg = fmt.Sprintf("method id: %s: invalid signature detected", AliceKey1)
	})

	It("Not Valid: DID signed by wrong controller", func(){
		valid =false
		msg = &types.MsgCreateDidPayload{
			Id:             ImposterDID,
			Authentication: []string{ImposterKey1},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 ImposterKey1,
					Type:               Ed25519VerificationKey2020,
					Controller:         ImposterDID,
					PublicKeyMultibase: "z" + base58.Encode(mainKeys[ImposterKey1].PublicKey),
				},
			},
		}
		signers = []string{AliceKey1}
		keys = map[string]KeyPair{
			AliceKey1: mainKeys[AliceKey1],
		}
		errMsg = fmt.Sprintf("signer: %s: signature is required but not found", ImposterDID)
	})

	It("Not Valid: DID self-signed by not existing verification method", func(){
		valid = false
		msg = &types.MsgCreateDidPayload{
			Id:             ImposterDID,
			Authentication: []string{ImposterKey1},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 ImposterKey1,
					Type:               Ed25519VerificationKey2020,
					Controller:         ImposterDID,
					PublicKeyMultibase: "z" + base58.Encode(mainKeys[ImposterKey1].PublicKey),
				},
			},
		}
		signers = []string{ImposterKey2}
		keys = map[string]KeyPair{
			ImposterKey2: GenerateKeyPair(),
		}
		errMsg = fmt.Sprintf("%s: verification method not found", ImposterKey2)
	})

	It("Not Valid: Self-signature not found", func(){
		valid = false
		msg = &types.MsgCreateDidPayload{
			Id:             ImposterDID,
			Controller:     []string{AliceDID, ImposterDID},
			Authentication: []string{ImposterKey1},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 ImposterKey1,
					Type:               Ed25519VerificationKey2020,
					Controller:         ImposterDID,
					PublicKeyMultibase: "z" + base58.Encode(mainKeys[ImposterKey1].PublicKey),
				},
			},
		}
		signers = []string{AliceKey1, ImposterKey2}
		keys = map[string]KeyPair{
			AliceKey1:    mainKeys[AliceKey1],
			ImposterKey2: GenerateKeyPair(),
		}
		errMsg = fmt.Sprintf("%s: verification method not found", ImposterKey2)
	})

	It("Not Valid: DID Doc already exists", func(){
		valid = false
		keys = map[string]KeyPair{
			CharlieKey1: GenerateKeyPair(),
		}
		signers = []string{CharlieKey1}
		msg = &types.MsgCreateDidPayload{
			Id:             CharlieDID,
			Authentication: []string{CharlieKey1},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         CharlieKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: CharlieDID,
				},
			},
		}
		errMsg = fmt.Sprintf("%s: DID Doc exists", CharlieDID)
	})
})

var _ = Describe("TestHandler", func() {
	It("Fails cause DIDDoc is already exists", func() {
		setup := Setup()

		_, _, _ = setup.InitDid(AliceDID)
		_, _, err := setup.InitDid(AliceDID)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(fmt.Sprintf("%s: DID Doc exists", AliceDID)))
	})
})
