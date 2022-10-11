package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/multiformats/go-multibase"
)

var _ = Describe("Update DID tests", func() {
	// params for cases
	var valid bool
	var signers []string
	var signerKeys []SignInput
	var msg *types.MsgUpdateDidPayload
	errMsg := ""

	var setup TestSetup
	var err error
	mainKeys := map[string]KeyPair{
		AliceKey1:    GenerateKeyPair(),
		AliceKey2:    GenerateKeyPair(),
		BobKey1:      GenerateKeyPair(),
		BobKey2:      GenerateKeyPair(),
		BobKey3:      GenerateKeyPair(),
		BobKey4:      GenerateKeyPair(),
		CharlieKey1:  GenerateKeyPair(),
		CharlieKey2:  GenerateKeyPair(),
		CharlieKey3:  GenerateKeyPair(),
		CharlieKey4:  GenerateKeyPair(),
		ImposterKey1: GenerateKeyPair(),
	}

	BeforeEach(func() {
		// setup
		valid = false
		signers = []string{}
		signerKeys = []SignInput{}
		msg = &types.MsgUpdateDidPayload{}
		errMsg = ""
	})

	AfterEach(func() {
		setup = InitEnv(mainKeys)

		for _, vm := range msg.VerificationMethod {
			if vm.PublicKeyMultibase == "" {
				vm.PublicKeyMultibase, err = multibase.Encode(multibase.Base58BTC, mainKeys[vm.Id].Public)
			}
			Expect(err).To(BeNil())
		}

		compiledKeys := []SignInput{}
		if len(signerKeys) > 0 {
			compiledKeys = signerKeys
		} else {
			for _, signer := range signers {
				compiledKeys = append(compiledKeys, SignInput{
					VerificationMethodId: signer,
					Key:                  mainKeys[signer].Private,
				})
			}
		}

		did, err := setup.SendUpdateDid(msg, compiledKeys)

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

	It("Valid: Key rotation works", func() {
		valid = true
		signerKeys = []SignInput{
			{
				VerificationMethodId: AliceKey1,
				Key:                  mainKeys[AliceKey1].Private,
			},
			{
				VerificationMethodId: AliceKey1,
				Key:                  mainKeys[AliceKey2].Private,
			},
		}
		msg = &types.MsgUpdateDidPayload{
			Id: AliceDID,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
					// TODO: Use multibase encoding
					PublicKeyMultibase: "z" + base58.Encode(mainKeys[AliceKey2].Public),
				},
			},
		}
	})

	It("Not Valid: replacing controller and Verification method ID does not work without new sign", func() {
		valid = false
		signers = []string{AliceKey2, BobKey1, AliceKey1}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{CharlieDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey2,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", CharlieDID)
	})

	It("Valid: replacing controller and Verification method ID works with all signatures", func() {
		valid = true
		signers = []string{AliceKey1, CharlieKey1, AliceKey2}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{CharlieDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey2,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", CharlieDID)
	})

	// Verification method's tests
	// cases:
	// - replacing VM controller works
	// - replacing VM controller does not work without new signature
	// - replacing VM controller does not work without old signature     ??????
	// - replacing VM doesn't work without new signature
	// - replacing VM doesn't work without old signature
	// - replacing VM works with all signatures
	// --- adding new VM works
	// --- adding new VM without new signature
	// --- adding new VM without old signature

	It("Valid: Replacing VM controller works with one signature", func() {
		valid = true
		signers = []string{AliceKey1, BobKey1}
		msg = &types.MsgUpdateDidPayload{
			Id: AliceDID,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: BobDID,
				},
			},
		}
	})

	It("Not Valid: Replacing VM controller does not work without new signature", func() {
		valid = false
		signers = []string{AliceKey1}
		msg = &types.MsgUpdateDidPayload{
			Id: AliceDID,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: BobDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID)
	})

	It("Not Valid: Replacing VM does not work without new signature", func() {
		valid = false
		signers = []string{AliceKey1}
		msg = &types.MsgUpdateDidPayload{
			Id: AliceDID,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey2,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one valid signature by %s (new version): invalid signature detected", AliceDID)
	})
	It("Not Valid: Replacing VM does not work without old signature", func() {
		valid = false
		signers = []string{AliceKey2}
		msg = &types.MsgUpdateDidPayload{
			Id: AliceDID,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey2,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", AliceDID)
	})

	It("Not Valid: Replacing VM works with all signatures", func() {
		valid = true
		signers = []string{AliceKey1, AliceKey2}
		msg = &types.MsgUpdateDidPayload{
			Id: AliceDID,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey2,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
	})

	// Adding VM

	It("Valid: Adding another verification method", func() {
		valid = true
		signers = []string{AliceKey1, BobKey1}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{AliceDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
				{
					Id:         AliceKey2,
					Type:       Ed25519VerificationKey2020,
					Controller: BobDID,
				},
			},
		}
	})

	It("Not Valid: Adding another verification method without new sign", func() {
		valid = false
		signers = []string{AliceKey1}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{AliceDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
				{
					Id:         AliceKey2,
					Type:       Ed25519VerificationKey2020,
					Controller: BobDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID)
	})

	It("Not Valid: Adding another verification method without old sign", func() {
		valid = false
		signers = []string{AliceKey2}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{AliceDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
				{
					Id:         AliceKey2,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", AliceDID)
	})

	// Controller's tests
	// cases:
	// - replacing Controller works with all signatures
	// - replacing Controller doesn't work without old signature
	// - replacing Controller doesn't work without new signature
	// --- adding Controller works with all signatures
	// --- adding Controller doesn't work without old signature
	// --- adding Controller doesn't work without new signature

	It("Valid: Replace controller works with all signatures", func() {
		valid = true
		signers = []string{BobKey1, AliceKey1}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{BobDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
	})

	It("Not Valid: Replace controller doesn't work without old signatures", func() {
		valid = false
		signers = []string{BobKey1}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{BobDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)
	})

	It("Not Valid: Replace controller doesn't work without new signatures", func() {
		valid = false
		signers = []string{AliceKey1}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{BobDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID)
	})

	// add controller

	It("Valid: Adding second controller works", func() {
		valid = true
		signers = []string{AliceKey1, CharlieKey3}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{AliceDID, CharlieDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
	})

	It("Not Valid: Adding controller without old signature", func() {
		valid = false
		signers = []string{BobKey1}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{AliceDID, BobDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)
	})

	It("Not Valid: Add controller without new signature doesn't work", func() {
		valid = false
		signers = []string{AliceKey1}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{AliceDID, BobDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID)
	})

	It("Valid: Adding verification method with the same controller works", func() {
		valid = true
		signers = []string{AliceKey1, AliceKey2}
		msg = &types.MsgUpdateDidPayload{
			Id:         AliceDID,
			Controller: []string{AliceDID},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         AliceKey2,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
				{
					Id:         AliceKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: AliceDID,
				},
			},
		}
	})

	It("Valid: Keeping VM with controller different then subject untouched during update should not require Bob signature", func() {
		valid = true
		signers = []string{CharlieKey1}
		msg = &types.MsgUpdateDidPayload{
			Id: CharlieDID,
			Authentication: []string{
				CharlieKey1,
				CharlieKey2,
				CharlieKey3,
			},

			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         CharlieKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: BobDID,
				},
				{
					Id:         CharlieKey2,
					Type:       Ed25519VerificationKey2020,
					Controller: BobDID,
				},
				{
					Id:         CharlieKey3,
					Type:       Ed25519VerificationKey2020,
					Controller: BobDID,
				},
				{
					Id:         CharlieKey4,
					Type:       Ed25519VerificationKey2020,
					Controller: CharlieDID,
				},
			},
		}
	})

	It("Valid: Removing verification method is possible with any kind of valid Bob's key", func() {
		valid = true
		signers = []string{BobKey1}
		msg = &types.MsgUpdateDidPayload{
			Id: BobDID,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:         BobKey1,
					Type:       Ed25519VerificationKey2020,
					Controller: BobDID,
				},
			},
		}
		errMsg = fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", BobDID)
	})
})
