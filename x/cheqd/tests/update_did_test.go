package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("Update DID tests", func() {
	var setup TestSetup

	BeforeEach(func() {
		setup = Setup()
	})

	It("Valid: Key rotation works", func() {
		// Create did
		did := setup.CreateSimpleDid()

		// Update did
		newKeyPair := GenerateKeyPair()

		msg := &types.MsgUpdateDidPayload{
			Id: did.Did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 did.KeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         did.Did,
					PublicKeyMultibase: MustEncodeBase58(newKeyPair.Public),
				},
			},
			VersionId: did.VersionId,
		}

		signatures := []SignInput{
			did.SignInput, // Old signature
			{
				VerificationMethodId: did.KeyId, // New signature
				Key:                  newKeyPair.Private,
			},
		}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDid(did.Did)
		Expect(err).To(BeNil())
		Expect(msg.ToDid()).To(Equal(*created.Did))
	})

	It("Not Valid: Did doc update does not work without DID doc controllers signature", func() {
		// Setup
		alice := setup.CreateSimpleDid()
		bob := setup.CreateDidWithExternalConterllers([]string{alice.Did}, []SignInput{alice.SignInput})

		// Update did
		msg := &types.MsgUpdateDidPayload{
			Id:         bob.Did,
			Controller: []string{alice.Did},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 bob.KeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         bob.Did,
					PublicKeyMultibase: MustEncodeBase58(bob.KeyPair.Public),
				},
			},
			Authentication:  []string{bob.KeyId},
			AssertionMethod: []string{bob.KeyId},
			VersionId:       bob.VersionId,
		}

		signatures := []SignInput{}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", alice.Did)))
	})

	It("Valid: Did doc update does not works with DID doc controllers signature", func() {
		// Setup
		alice := setup.CreateSimpleDid()
		bob := setup.CreateDidWithExternalConterllers([]string{alice.Did}, []SignInput{alice.SignInput})

		// Update did
		msg := &types.MsgUpdateDidPayload{
			Id:         bob.Did,
			Controller: []string{alice.Did},
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 bob.KeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         bob.Did,
					PublicKeyMultibase: MustEncodeBase58(bob.KeyPair.Public),
				},
			},
			Authentication:  []string{bob.KeyId},
			AssertionMethod: []string{bob.KeyId},
			VersionId:       bob.VersionId,
		}

		signatures := []SignInput{alice.SignInput}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDid(bob.Did)
		Expect(err).To(BeNil())
		Expect(msg.ToDid()).To(Equal(*created.Did))
	})

	// Verification method's tests
	// cases:
	// - replacing VM controller works
	// - replacing VM controller does not work without new signature
	// - replacing VM controller does not work without old signature
	// - replacing VM doesn't work without new signature
	// - replacing VM doesn't work without old signature
	// - replacing VM works with all signatures
	// --- adding new VM works
	// --- adding new VM without new signature
	// --- adding new VM without old signature

	It("Valid: Replacing VM controller works with two signatures", func() {
		// Setup
		alice := setup.CreateSimpleDid()
		bob := setup.CreateSimpleDid()

		// Update did
		msg := &types.MsgUpdateDidPayload{
			Id: alice.Did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 alice.KeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         bob.Did,
					PublicKeyMultibase: MustEncodeBase58(alice.KeyPair.Public),
				},
			},
			Authentication: []string{alice.KeyId},
			VersionId:      alice.VersionId,
		}

		signatures := []SignInput{alice.SignInput, bob.SignInput}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		updated, err := setup.QueryDid(alice.Did)
		Expect(err).To(BeNil())
		Expect(*updated.Did).To(Equal(msg.ToDid()))
	})

	It("Not Valid: Replacing VM controller does not work without new signature", func() {
		// Setup
		alice := setup.CreateSimpleDid()
		bob := setup.CreateSimpleDid()

		// Update did
		msg := &types.MsgUpdateDidPayload{
			Id: alice.Did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 alice.KeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         bob.Did, // Previously alice
					PublicKeyMultibase: MustEncodeBase58(alice.KeyPair.Public),
				},
			},
			Authentication: []string{alice.KeyId},
			VersionId:      alice.VersionId,
		}

		signatures := []SignInput{alice.SignInput}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", bob.Did)))
	})

	It("Not Valid: Replacing VM controller does not work without old signature", func() {
		// Setup
		alice := setup.CreateSimpleDid()
		bob := setup.CreateSimpleDid()

		// Update did
		msg := &types.MsgUpdateDidPayload{
			Id: alice.Did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 alice.KeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         bob.Did,
					PublicKeyMultibase: MustEncodeBase58(alice.KeyPair.Public),
				},
			},
			Authentication: []string{alice.KeyId},
			VersionId:      alice.VersionId,
		}

		signatures := []SignInput{bob.SignInput}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", alice.Did)))
	})

	It("Not Valid: Replacing VM does not work without new signature", func() {
		// Setup
		alice := setup.CreateSimpleDid()

		// Update did
		newKeyId := alice.Did + "#key-2"

		msg := &types.MsgUpdateDidPayload{
			Id: alice.Did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 newKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         alice.Did,
					PublicKeyMultibase: MustEncodeBase58(alice.KeyPair.Public),
				},
			},
			Authentication: []string{alice.KeyId},
			VersionId:      alice.VersionId,
		}

		signatures := []SignInput{alice.SignInput}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (new version): invalid signature detected", alice.Did)))
	})

	It("Not Valid: Replacing VM does not work without old signature", func() {
		// Setup
		alice := setup.CreateSimpleDid()

		// Update did
		newKeyId := alice.Did + "#key-2"

		msg := &types.MsgUpdateDidPayload{
			Id: alice.Did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 newKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         alice.Did,
					PublicKeyMultibase: MustEncodeBase58(alice.KeyPair.Public),
				},
			},
			Authentication: []string{alice.KeyId},
			VersionId:      alice.VersionId,
		}

		signatures := []SignInput{
			{
				VerificationMethodId: newKeyId,
				Key:                  alice.KeyPair.Private,
			},
		}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", alice.Did)))
	})

	It("Not Valid: Replacing VM works with all signatures", func() {
		// Setup
		alice := setup.CreateSimpleDid()

		// Update did
		newKeyId := alice.Did + "#key-2"

		msg := &types.MsgUpdateDidPayload{
			Id: alice.Did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 newKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         alice.Did,
					PublicKeyMultibase: MustEncodeBase58(alice.KeyPair.Public),
				},
			},
			Authentication: []string{alice.KeyId},
			VersionId:      alice.VersionId,
		}

		signatures := []SignInput{
			alice.SignInput,
			{
				VerificationMethodId: newKeyId,
				Key:                  alice.KeyPair.Private,
			},
		}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDid(alice.Did)
		Expect(err).To(BeNil())
		Expect(*created).ToNot(Equal(msg.ToDid()))
	})

	// // Adding VM

	It("Valid: Adding another verification method", func() {
		// Setup
		alice := setup.CreateSimpleDid()

		// Update did
		newKeyId := alice.Did + "#key-2"
		newKey := GenerateKeyPair()

		msg := &types.MsgUpdateDidPayload{
			Id: alice.Did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 alice.KeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         alice.Did,
					PublicKeyMultibase: MustEncodeBase58(alice.KeyPair.Public),
				},
				{
					Id:                 newKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         alice.Did,
					PublicKeyMultibase: MustEncodeBase58(newKey.Public),
				},
			},
			Authentication: []string{alice.KeyId},
			VersionId:      alice.VersionId,
		}

		signatures := []SignInput{
			alice.SignInput,
			{
				VerificationMethodId: newKeyId,
				Key:                  newKey.Private,
			},
		}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err).To(BeNil())

		// check
		created, err := setup.QueryDid(alice.Did)
		Expect(err).To(BeNil())
		Expect(*created).ToNot(Equal(msg.ToDid()))
	})

	It("Not Valid: Adding another verification method without new sign", func() {
		// Setup
		alice := setup.CreateSimpleDid()

		// Update did
		newKeyId := alice.Did + "#key-2"
		newKey := GenerateKeyPair()

		msg := &types.MsgUpdateDidPayload{
			Id: alice.Did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 alice.KeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         alice.Did,
					PublicKeyMultibase: MustEncodeBase58(alice.KeyPair.Public),
				},
				{
					Id:                 newKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         alice.Did,
					PublicKeyMultibase: MustEncodeBase58(newKey.Public),
				},
			},
			Authentication: []string{alice.KeyId},
			VersionId:      alice.VersionId,
		}

		signatures := []SignInput{
			alice.SignInput,
		}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID)))
	})

	It("Not Valid: Adding another verification method without old sign", func() {
		// Setup
		alice := setup.CreateSimpleDid()

		// Update did
		newKeyId := alice.Did + "#key-2"
		newKey := GenerateKeyPair()

		msg := &types.MsgUpdateDidPayload{
			Id: alice.Did,
			VerificationMethod: []*types.VerificationMethod{
				{
					Id:                 alice.KeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         alice.Did,
					PublicKeyMultibase: MustEncodeBase58(alice.KeyPair.Public),
				},
				{
					Id:                 newKeyId,
					Type:               types.Ed25519VerificationKey2020,
					Controller:         alice.Did,
					PublicKeyMultibase: MustEncodeBase58(newKey.Public),
				},
			},
			Authentication: []string{alice.KeyId},
			VersionId:      alice.VersionId,
		}

		signatures := []SignInput{
			{
				VerificationMethodId: newKeyId,
				Key:                  newKey.Private,
			},
		}

		_, err := setup.UpdateDid(msg, signatures)
		Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", AliceDID)))
	})

	// // Controller's tests
	// // cases:
	// // - replacing Controller works with all signatures
	// // - replacing Controller doesn't work without old signature
	// // - replacing Controller doesn't work without new signature
	// // --- adding Controller works with all signatures
	// // --- adding Controller doesn't work without old signature
	// // --- adding Controller doesn't work without new signature

	// It("Valid: Replace controller works with all signatures", func() {
	// 	valid = true
	// 	signers = []string{BobKey1, AliceKey1}
	// 	msg = &types.MsgUpdateDidPayload{
	// 		Id:         AliceDID,
	// 		Controller: []string{BobDID},
	// 		VerificationMethod: []*types.VerificationMethod{
	// 			{
	// 				Id:         AliceKey1,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: AliceDID,
	// 			},
	// 		},
	// 	}
	// })

	// It("Not Valid: Replace controller doesn't work without old signatures", func() {
	// 	valid = false
	// 	signers = []string{BobKey1}
	// 	msg = &types.MsgUpdateDidPayload{
	// 		Id:         AliceDID,
	// 		Controller: []string{BobDID},
	// 		VerificationMethod: []*types.VerificationMethod{
	// 			{
	// 				Id:         AliceKey1,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: AliceDID,
	// 			},
	// 		},
	// 	}
	// 	errMsg = fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)
	// })

	// It("Not Valid: Replace controller doesn't work without new signatures", func() {
	// 	valid = false
	// 	signers = []string{AliceKey1}
	// 	msg = &types.MsgUpdateDidPayload{
	// 		Id:         AliceDID,
	// 		Controller: []string{BobDID},
	// 		VerificationMethod: []*types.VerificationMethod{
	// 			{
	// 				Id:         AliceKey1,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: AliceDID,
	// 			},
	// 		},
	// 	}
	// 	errMsg = fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID)
	// })

	// // add controller

	// It("Valid: Adding second controller works", func() {
	// 	valid = true
	// 	signers = []string{AliceKey1, CharlieKey3}
	// 	msg = &types.MsgUpdateDidPayload{
	// 		Id:         AliceDID,
	// 		Controller: []string{AliceDID, CharlieDID},
	// 		VerificationMethod: []*types.VerificationMethod{
	// 			{
	// 				Id:         AliceKey1,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: AliceDID,
	// 			},
	// 		},
	// 	}
	// })

	// It("Not Valid: Adding controller without old signature", func() {
	// 	valid = false
	// 	signers = []string{BobKey1}
	// 	msg = &types.MsgUpdateDidPayload{
	// 		Id:         AliceDID,
	// 		Controller: []string{AliceDID, BobDID},
	// 		VerificationMethod: []*types.VerificationMethod{
	// 			{
	// 				Id:         AliceKey1,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: AliceDID,
	// 			},
	// 		},
	// 	}
	// 	errMsg = fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)
	// })

	// It("Not Valid: Add controller without new signature doesn't work", func() {
	// 	valid = false
	// 	signers = []string{AliceKey1}
	// 	msg = &types.MsgUpdateDidPayload{
	// 		Id:         AliceDID,
	// 		Controller: []string{AliceDID, BobDID},
	// 		VerificationMethod: []*types.VerificationMethod{
	// 			{
	// 				Id:         AliceKey1,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: AliceDID,
	// 			},
	// 		},
	// 	}
	// 	errMsg = fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", BobDID)
	// })

	// It("Valid: Adding verification method with the same controller works", func() {
	// 	valid = true
	// 	signers = []string{AliceKey1, AliceKey2}
	// 	msg = &types.MsgUpdateDidPayload{
	// 		Id:         AliceDID,
	// 		Controller: []string{AliceDID},
	// 		VerificationMethod: []*types.VerificationMethod{
	// 			{
	// 				Id:         AliceKey2,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: AliceDID,
	// 			},
	// 			{
	// 				Id:         AliceKey1,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: AliceDID,
	// 			},
	// 		},
	// 	}
	// })

	// It("Valid: Keeping VM with controller different then subject untouched during update should not require Bob signature", func() {
	// 	valid = true
	// 	signers = []string{CharlieKey1}
	// 	msg = &types.MsgUpdateDidPayload{
	// 		Id: CharlieDID,
	// 		Authentication: []string{
	// 			CharlieKey1,
	// 			CharlieKey2,
	// 			CharlieKey3,
	// 		},

	// 		VerificationMethod: []*types.VerificationMethod{
	// 			{
	// 				Id:         CharlieKey1,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: BobDID,
	// 			},
	// 			{
	// 				Id:         CharlieKey2,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: BobDID,
	// 			},
	// 			{
	// 				Id:         CharlieKey3,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: BobDID,
	// 			},
	// 			{
	// 				Id:         CharlieKey4,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: CharlieDID,
	// 			},
	// 		},
	// 	}
	// })

	// It("Valid: Removing verification method is possible with any kind of valid Bob's key", func() {
	// 	valid = true
	// 	signers = []string{BobKey1}
	// 	msg = &types.MsgUpdateDidPayload{
	// 		Id: BobDID,
	// 		VerificationMethod: []*types.VerificationMethod{
	// 			{
	// 				Id:         BobKey1,
	// 				Type:       types.Ed25519VerificationKey2020,
	// 				Controller: BobDID,
	// 			},
	// 		},
	// 	}
	// 	errMsg = fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", BobDID)
	// })
})

// var _ = Describe("Signature Verification while updating DID", func() {
// 	var setup TestSetup
// 	var aliceKeys, bobKeys map[string]ed25519.PrivateKey
// 	var aliceDid *types.MsgCreateDidPayload
// 	BeforeEach(func() {
// 		setup = Setup()
// 		aliceKeys, aliceDid, _ = setup.InitDid(AliceDID)
// 		bobKeys, _, _ = setup.InitDid(BobDID)
// 	})

// 	It("should have changed DIDDoc controller", func() {
// 		updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
// 		updatedDidDoc.Controller = append(updatedDidDoc.Controller, BobDID)
// 		receivedDid, _ := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(ConcatKeys(aliceKeys, bobKeys)))

// 		Expect(aliceDid.Controller).To(Not(Equal(receivedDid.Controller)))
// 		Expect([]string{AliceDID, BobDID}).To(Not(Equal(receivedDid.Controller)))
// 		Expect([]string{BobDID}, receivedDid.Controller)
// 	})

// 	When("Old signature in verification method is absent", func() {
// 		It("should fail", func() {
// 			updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
// 			updatedDidDoc.VerificationMethod[0].Type = types.Ed25519VerificationKey2020
// 			_, err := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(bobKeys))

// 			// check
// 			Expect(err).To(Not(BeNil()))
// 			Expect(err.Error()).To(Equal(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)))
// 		})
// 	})

// 	It("should fails cause we need old signature for changing verification method controller", func() {
// 		updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
// 		updatedDidDoc.VerificationMethod[0].Controller = BobDID
// 		_, err := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(bobKeys))

// 		// check
// 		Expect(err).To(Not(BeNil()))
// 		Expect(err.Error()).To(Equal(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)))
// 	})

// 	It("should fails cause we need old signature for changing DIDDoc controller", func() {
// 		updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
// 		updatedDidDoc.Controller = append(updatedDidDoc.Controller, BobDID)
// 		_, err := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(bobKeys))

// 		// check
// 		Expect(err).To(Not(BeNil()))
// 		Expect(err.Error()).To(Equal(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)))
// 	})
// })

// var _ = Describe("Signature Verification. Remove signature/VM", func() {
// 	var setup TestSetup
// 	var ApubKey, BpubKey ed25519.PublicKey
// 	var AprivKey, BprivKey ed25519.PrivateKey
// 	var aliceDid, bobDid *types.MsgCreateDidPayload
// 	var aliceKeys, bobKeys map[string]ed25519.PrivateKey

// 	BeforeEach(func() {
// 		setup = Setup()
// 		// Generate keys
// 		ApubKey, AprivKey, _ = ed25519.GenerateKey(rand.Reader)
// 		BpubKey, BprivKey, _ = ed25519.GenerateKey(rand.Reader)

// 		// Create dids
// 		aliceDid = setup.BuildSimpleCreateDidPayload(AliceDID, AliceKey1, ApubKey)
// 		bobDid = setup.BuildSimpleCreateDidPayload(BobDID, BobKey1, BpubKey)

// 		// Collect private keys
// 		aliceKeys = map[string]ed25519.PrivateKey{AliceKey1: AprivKey, BobKey1: BprivKey}
// 		bobKeys = map[string]ed25519.PrivateKey{BobKey1: BprivKey}

// 		// Add verification method
// 		aliceDid.VerificationMethod = append(aliceDid.VerificationMethod, &types.VerificationMethod{
// 			Id:                 AliceKey2,
// 			Controller:         BobDID,
// 			Type:               types.Ed25519VerificationKey2020,
// 			PublicKeyMultibase: MustEncodeBase58(BpubKey),
// 		})
// 	})

// 	It("should fails cause old signature is required for removing this signature", func() {
// 		// Send dids
// 		_, _ = setup.SendCreateDid(bobDid, bobKeys)
// 		_, _ = setup.SendCreateDid(aliceDid, aliceKeys)

// 		updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
// 		updatedDidDoc.VerificationMethod = []*types.VerificationMethod{aliceDid.VerificationMethod[0]}
// 		updatedDidDoc.Authentication = []string{aliceDid.Authentication[0]}
// 		_, err := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(bobKeys))

// 		// check
// 		Expect(err).To(Not(BeNil()))
// 		Expect(err.Error()).To(Equal(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)))
// 	})

// 	It("should not fails while removing the whole verification method", func() {
// 		aliceDid.Authentication = append(aliceDid.Authentication, AliceKey2)

// 		// Send dids
// 		_, _ = setup.SendCreateDid(bobDid, bobKeys)
// 		_, _ = setup.SendCreateDid(aliceDid, aliceKeys)

// 		updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
// 		updatedDidDoc.Authentication = []string{aliceDid.Authentication[0]}
// 		updatedDidDoc.VerificationMethod = []*types.VerificationMethod{aliceDid.VerificationMethod[0]}
// 		receivedDid, _ := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(ConcatKeys(aliceKeys, bobKeys)))

// 		// check
// 		Expect(len(aliceDid.VerificationMethod)).To(Not(Equal(len(receivedDid.VerificationMethod))))
// 		Expect(reflect.DeepEqual(aliceDid.VerificationMethod[0], receivedDid.VerificationMethod[0])).To(BeTrue())
// 	})
// })
