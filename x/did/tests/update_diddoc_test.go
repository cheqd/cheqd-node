package tests

import (
	"fmt"

	. "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/cheqd-node/x/did/types"
)

var _ = Describe("DIDDoc update", func() {
	var setup TestSetup

	BeforeEach(func() {
		setup = Setup()
	})

	Describe("DIDDoc: Update verification relationship", func() {
		var alice CreatedDidDocInfo
		var bob CreatedDidDocInfo
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			bob = setup.CreateDidDocWithExternalControllers([]string{alice.Did}, []SignInput{alice.SignInput})

			msg = &types.MsgUpdateDidDocPayload{
				Id:         bob.Did,
				Controller: []string{alice.Did},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     bob.KeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             bob.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(bob.KeyPair.Public),
					},
				},
				Authentication:  []string{bob.KeyID},
				AssertionMethod: []string{bob.KeyID},
				VersionId:       uuid.NewString(),
			}
		})

		It("Works with DIDDoc controller signatures", func() {
			signatures := []SignInput{alice.SignInput}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryDidDoc(bob.Did)
			Expect(err).To(BeNil())
			Expect(msg.ToDidDoc()).To(Equal(*created.DidDocWithMetadata.DidDoc))
		})

		It("Creates a new DIDDoc version in case of success", func() {
			signatures := []SignInput{alice.SignInput}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err).To(BeNil())

			// check latest version
			created, err := setup.QueryDidDoc(bob.Did)
			Expect(err).To(BeNil())
			Expect(msg.ToDidDoc()).To(Equal(*created.DidDocWithMetadata.DidDoc))

			// query the first version
			v1, err := setup.QueryDidDocVersion(bob.Did, bob.VersionID)
			Expect(err).To(BeNil())
			Expect(*v1.DidDocWithMetadata.DidDoc).To(Equal(bob.Msg.ToDidDoc()))
			Expect(v1.DidDocWithMetadata.Metadata.VersionId).To(Equal(bob.VersionID))
			Expect(v1.DidDocWithMetadata.Metadata.NextVersionId).To(Equal(msg.VersionId))

			// query the second version
			v2, err := setup.QueryDidDocVersion(bob.Did, msg.VersionId)
			Expect(err).To(BeNil())
			Expect(*v2.DidDocWithMetadata.DidDoc).To(Equal(msg.ToDidDoc()))
			Expect(v2.DidDocWithMetadata.Metadata.VersionId).To(Equal(msg.VersionId))
			Expect(v2.DidDocWithMetadata.Metadata.PreviousVersionId).To(Equal(bob.VersionID))

			// query all versions
			versions, err := setup.QueryAllDidDocVersionsMetadata(bob.Did)
			Expect(err).To(BeNil())
			Expect(versions.Versions).To(HaveLen(2))
			Expect(versions.Versions).To(ContainElement(v1.DidDocWithMetadata.Metadata))
			Expect(versions.Versions).To(ContainElement(v2.DidDocWithMetadata.Metadata))
		})

		It("Doesn't work without controller signatures", func() {
			signatures := []SignInput{}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", alice.Did)))
		})
	})

	Describe("DIDDoc: Replacing controller", func() {
		var alice CreatedDidDocInfo
		var bob CreatedDidDocInfo
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			bob = setup.CreateSimpleDid()

			msg = &types.MsgUpdateDidDocPayload{
				Id:         alice.Did,
				Controller: []string{bob.Did},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     alice.KeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             alice.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				VersionId: uuid.NewString(),
			}
		})

		It("Works with old and new signatures", func() {
			signatures := []SignInput{
				alice.SignInput,
				bob.SignInput,
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err).To(BeNil())

			// check
			updated, err := setup.QueryDidDoc(alice.Did)
			Expect(err).To(BeNil())
			Expect(*updated.DidDocWithMetadata.DidDoc).To(Equal(msg.ToDidDoc()))
		})

		It("Doesn't work with only new controller signature", func() {
			signatures := []SignInput{
				bob.SignInput,
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", alice.Did)))
		})

		It("Doesn't work with only old controller signature", func() {
			signatures := []SignInput{
				alice.SignInput,
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", bob.Did)))
		})
	})

	Describe("DIDDoc: Adding controller", func() {
		var alice CreatedDidDocInfo
		var bob CreatedDidDocInfo
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			bob = setup.CreateSimpleDid()

			msg = &types.MsgUpdateDidDocPayload{
				Id:         alice.Did,
				Controller: []string{alice.Did, bob.Did},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     alice.KeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             alice.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				VersionId: uuid.NewString(),
			}
		})

		It("Works with old and new signatures", func() {
			signatures := []SignInput{
				alice.SignInput,
				bob.SignInput,
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err).To(BeNil())

			// check
			updated, err := setup.QueryDidDoc(alice.Did)
			Expect(err).To(BeNil())
			Expect(*updated.DidDocWithMetadata.DidDoc).To(Equal(msg.ToDidDoc()))
		})

		It("Doesn't work with only new controller signatures", func() {
			signatures := []SignInput{
				bob.SignInput,
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", alice.Did)))
		})

		It("Doesn't work with only old controller signatures", func() {
			signatures := []SignInput{
				alice.SignInput,
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", bob.Did)))
		})
	})

	Describe("DIDDoc: Keep verification method with controller different than subject untouched during update", func() {
		var alice CreatedDidDocInfo
		var bob CreatedDidDocInfo
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			bob = setup.CreateSimpleDid()
			alice = setup.CreateDidDocWithExternalControllers([]string{bob.Did}, []SignInput{bob.SignInput})

			msg = &types.MsgUpdateDidDocPayload{
				Id:         alice.Did,
				Controller: []string{bob.Did},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     alice.KeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             alice.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				Authentication:  []string{alice.KeyID},
				AssertionMethod: []string{alice.KeyID}, // Adding new verification method
				VersionId:       uuid.NewString(),
			}
		})

		It("Doesn't require verification method controller signature", func() {
			signatures := []SignInput{
				bob.SignInput,
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryDidDoc(alice.Did)
			Expect(err).To(BeNil())
			Expect(*created).ToNot(Equal(msg.ToDidDoc()))
		})
	})

	Describe("Verification method: key udpate", func() {
		var did CreatedDidDocInfo
		var newKeyPair KeyPair
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			did = setup.CreateSimpleDid()
			newKeyPair = GenerateKeyPair()
			msg = &types.MsgUpdateDidDocPayload{
				Id: did.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     did.KeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             did.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(newKeyPair.Public),
					},
				},
				VersionId: uuid.NewString(),
			}
		})

		It("Works with old and new signatures", func() {
			signatures := []SignInput{
				did.SignInput, // Old signature
				{
					VerificationMethodID: did.KeyID, // New signature
					Key:                  newKeyPair.Private,
				},
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryDidDoc(did.Did)
			Expect(err).To(BeNil())
			Expect(msg.ToDidDoc()).To(Equal(*created.DidDocWithMetadata.DidDoc))
		})

		It("Doesn't work without new signature", func() {
			signatures := []SignInput{did.SignInput}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (new version): invalid signature detected", did.Did)))
		})

		It("Doesn't work without old signature", func() {
			signatures := []SignInput{{
				VerificationMethodID: did.KeyID,
				Key:                  newKeyPair.Private,
			}}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", did.Did)))
		})
	})

	Describe("Verification method: controller update", func() {
		var alice CreatedDidDocInfo
		var bob CreatedDidDocInfo
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			bob = setup.CreateSimpleDid()

			msg = &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     alice.KeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             bob.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				Authentication: []string{alice.KeyID},
				VersionId:      uuid.NewString(),
			}
		})

		It("Works with old and new controller signatures", func() {
			signatures := []SignInput{alice.SignInput, bob.SignInput}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err).To(BeNil())

			// check
			updated, err := setup.QueryDidDoc(alice.Did)
			Expect(err).To(BeNil())
			Expect(*updated.DidDocWithMetadata.DidDoc).To(Equal(msg.ToDidDoc()))
		})

		It("Doesn't work without old controller signature", func() {
			signatures := []SignInput{bob.SignInput}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", alice.Did)))
		})

		It("Doesn't work without new controller signature", func() {
			signatures := []SignInput{alice.SignInput}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", bob.Did)))
		})
	})

	Describe("Verification method: ID update", func() {
		var alice CreatedDidDocInfo
		var newKeyID string
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			newKeyID = alice.Did + "#key-2"

			msg = &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     newKeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             alice.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				Authentication: []string{alice.KeyID},
				VersionId:      uuid.NewString(),
			}
		})

		It("Doesn't work without new verification method signature", func() {
			signatures := []SignInput{alice.SignInput}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (new version): invalid signature detected", alice.Did)))
		})

		It("Doesn't work without old verification method signature", func() {
			signatures := []SignInput{
				{
					VerificationMethodID: newKeyID,
					Key:                  alice.KeyPair.Private,
				},
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", alice.Did)))
		})

		It("Works with new and old verification method signature", func() {
			signatures := []SignInput{
				{
					VerificationMethodID: newKeyID,
					Key:                  alice.KeyPair.Private,
				},
				alice.SignInput,
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err).To(BeNil())

			// check
			updated, err := setup.QueryDidDoc(alice.Did)
			Expect(err).To(BeNil())
			Expect(*updated.DidDocWithMetadata.DidDoc).To(Equal(msg.ToDidDoc()))
		})
	})

	Describe("Verification method: adding a new one", func() {
		var alice CreatedDidDocInfo
		var newKeyID string
		var newKey KeyPair
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()

			newKeyID = alice.Did + "#key-2"
			newKey = GenerateKeyPair()

			msg = &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     alice.KeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             alice.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
					{
						Id:                     newKeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             alice.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(newKey.Public),
					},
				},
				Authentication: []string{alice.KeyID},
				VersionId:      uuid.NewString(),
			}
		})

		It("Works with only old verification method signature", func() {
			signatures := []SignInput{
				alice.SignInput,
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryDidDoc(alice.Did)
			Expect(err).To(BeNil())
			Expect(*created).ToNot(Equal(msg.ToDidDoc()))
		})

		It("Doesn't work with only new verification method signature", func() {
			signatures := []SignInput{
				{
					VerificationMethodID: newKeyID,
					Key:                  newKey.Private,
				},
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", alice.Did)))
		})
	})

	Describe("Verification method: removing existing", func() {
		var alice CreatedDidDocInfo
		var secondKeyID string
		var secondKey KeyPair
		var secondSignInput SignInput
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()

			secondKeyID = alice.Did + "#key-2"
			secondKey = GenerateKeyPair()
			secondSignInput = SignInput{
				VerificationMethodID: secondKeyID,
				Key:                  secondKey.Private,
			}

			addSecondKeyMsg := &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     alice.KeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             alice.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
					{
						Id:                     secondKeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             alice.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(secondKey.Public),
					},
				},
				Authentication: []string{alice.KeyID},
				VersionId:      uuid.NewString(),
			}

			_, err := setup.UpdateDidDoc(addSecondKeyMsg, []SignInput{alice.SignInput})
			Expect(err).To(BeNil())

			msg = &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     alice.KeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             alice.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				Authentication: []string{alice.KeyID},
				VersionId:      uuid.NewString(),
			}
		})

		It("Works with only first verification method signature", func() {
			signatures := []SignInput{
				alice.SignInput,
			}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryDidDoc(alice.Did)
			Expect(err).To(BeNil())
			Expect(*created).ToNot(Equal(msg.ToDidDoc()))
		})

		It("Doesn't work with only second verification method signature (which will be deleted)", func() {
			signatures := []SignInput{secondSignInput}

			_, err := setup.UpdateDidDoc(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (new version): invalid signature detected", alice.Did)))
		})
	})

	Describe("Deactivating", func() {
		var alice CreatedDidDocInfo
		var bob CreatedDidDocInfo
		var updateMsg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			bob = setup.CreateSimpleDid()

			updateMsg = &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                     alice.DidDocInfo.KeyID,
						VerificationMethodType: types.Ed25519VerificationKey2020Type,
						Controller:             alice.DidDocInfo.Did,
						VerificationMaterial:   GenerateEd25519VerificationKey2020VerificationMaterial(alice.DidDocInfo.KeyPair.Public),
					},
				},
				Authentication: []string{alice.KeyID},
				VersionId:      uuid.NewString(),
			}
		})

		When("Updating already deactivated DID", func() {
			It("Should fail with error", func() {
				// Deactivate DID
				deactivateMsg := &types.MsgDeactivateDidDocPayload{
					Id:        alice.Did,
					VersionId: uuid.NewString(),
				}

				signatures := []SignInput{alice.DidDocInfo.SignInput}

				res, err := setup.DeactivateDidDoc(deactivateMsg, signatures)
				Expect(err).To(BeNil())
				Expect(res.Value.Metadata.Deactivated).To(BeTrue())

				// Update deactivated DID
				signatures = []SignInput{
					alice.SignInput,
					bob.SignInput,
				}

				_, err = setup.UpdateDidDoc(updateMsg, signatures)
				Expect(err.Error()).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring(alice.DidDocInfo.Did + ": DID Doc already deactivated"))
			})
		})
	})
})
