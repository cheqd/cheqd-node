package tests

import (
	"fmt"

	. "github.com/cheqd/cheqd-node/x/cheqd/tests/setup"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("DID Doc update", func() {
	var setup TestSetup

	BeforeEach(func() {
		setup = Setup()
	})

	Describe("DIDDoc: update verification relationship", func() {
		var alice CreatedDidInfo
		var bob CreatedDidInfo
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			bob = setup.CreateDidWithExternalConterllers([]string{alice.Did}, []SignInput{alice.SignInput})

			msg = &types.MsgUpdateDidDocPayload{
				Id:         bob.Did,
				Controller: []string{alice.Did},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   bob.KeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           bob.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(bob.KeyPair.Public),
					},
				},
				Authentication:  []string{bob.KeyId},
				AssertionMethod: []string{bob.KeyId},
				VersionId:       bob.VersionId,
			}
		})

		It("Works with DID doc controllers signature", func() {
			signatures := []SignInput{alice.SignInput}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryDid(bob.Did)
			Expect(err).To(BeNil())
			Expect(msg.ToDidDoc()).To(Equal(*created.Value.DidDoc))
		})

		It("Doesn't work without controllers signatures", func() {
			signatures := []SignInput{}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", alice.Did)))
		})
	})

	Describe("DIDDoc: replacing controller", func() {
		var alice CreatedDidInfo
		var bob CreatedDidInfo
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			bob = setup.CreateSimpleDid()

			msg = &types.MsgUpdateDidDocPayload{
				Id:         alice.Did,
				Controller: []string{bob.Did},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   alice.KeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           alice.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				VersionId: alice.VersionId,
			}
		})

		It("Works with old and new signatures", func() {
			signatures := []SignInput{
				alice.SignInput,
				bob.SignInput,
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err).To(BeNil())

			// check
			updated, err := setup.QueryDid(alice.Did)
			Expect(err).To(BeNil())
			Expect(*updated.Value.DidDoc).To(Equal(msg.ToDidDoc()))
		})

		It("Doesn't work with only new controller signature", func() {
			signatures := []SignInput{
				bob.SignInput,
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", alice.Did)))
		})

		It("Doesn't work with only old controller signature", func() {
			signatures := []SignInput{
				alice.SignInput,
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", bob.Did)))
		})
	})

	Describe("DIDDoc: adding controller", func() {
		var alice CreatedDidInfo
		var bob CreatedDidInfo
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			bob = setup.CreateSimpleDid()

			msg = &types.MsgUpdateDidDocPayload{
				Id:         alice.Did,
				Controller: []string{alice.Did, bob.Did},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   alice.KeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           alice.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				VersionId: alice.VersionId,
			}
		})

		It("Works with old and new signatures", func() {
			signatures := []SignInput{
				alice.SignInput,
				bob.SignInput,
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err).To(BeNil())

			// check
			updated, err := setup.QueryDid(alice.Did)
			Expect(err).To(BeNil())
			Expect(*updated.Value.DidDoc).To(Equal(msg.ToDidDoc()))
		})

		It("Doesn't work with only new controller signatures", func() {
			signatures := []SignInput{
				bob.SignInput,
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", alice.Did)))
		})

		It("Doesn't work with only old controller signatures", func() {
			signatures := []SignInput{
				alice.SignInput,
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", bob.Did)))
		})
	})

	Describe("DIDDoc: Keeping VM with controller different then subject untouched during update", func() {
		var alice CreatedDidInfo
		var bob CreatedDidInfo
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			bob = setup.CreateSimpleDid()
			alice = setup.CreateDidWithExternalConterllers([]string{bob.Did}, []SignInput{bob.SignInput})

			msg = &types.MsgUpdateDidDocPayload{
				Id:         alice.Did,
				Controller: []string{bob.Did},
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   alice.KeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           alice.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				Authentication:  []string{alice.KeyId},
				AssertionMethod: []string{alice.KeyId}, // Adding new VM
				VersionId:       alice.VersionId,
			}
		})

		It("Doesn't require VM's controler signature", func() {
			signatures := []SignInput{
				bob.SignInput,
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryDid(alice.Did)
			Expect(err).To(BeNil())
			Expect(*created).ToNot(Equal(msg.ToDidDoc()))
		})
	})

	Describe("Verification method: key udpate", func() {
		var did CreatedDidInfo
		var newKeyPair KeyPair
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			did = setup.CreateSimpleDid()
			newKeyPair = GenerateKeyPair()
			msg = &types.MsgUpdateDidDocPayload{
				Id: did.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   did.KeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           did.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(newKeyPair.Public),
					},
				},
				VersionId: did.VersionId,
			}
		})

		It("Works with two signatures", func() {
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
			Expect(msg.ToDidDoc()).To(Equal(*created.Value.DidDoc))
		})

		It("Doesn't work without new signature", func() {
			signatures := []SignInput{did.SignInput}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (new version): invalid signature detected", did.Did)))
		})

		It("Doesn't work without old signature", func() {
			signatures := []SignInput{{
				VerificationMethodId: did.KeyId,
				Key:                  newKeyPair.Private,
			}}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", did.Did)))
		})
	})

	Describe("Verification method: controller update", func() {
		var alice CreatedDidInfo
		var bob CreatedDidInfo
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			bob = setup.CreateSimpleDid()

			msg = &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   alice.KeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           bob.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				Authentication: []string{alice.KeyId},
				VersionId:      alice.VersionId,
			}
		})

		It("Works with old and new controller signature", func() {
			signatures := []SignInput{alice.SignInput, bob.SignInput}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err).To(BeNil())

			// check
			updated, err := setup.QueryDid(alice.Did)
			Expect(err).To(BeNil())
			Expect(*updated.Value.DidDoc).To(Equal(msg.ToDidDoc()))
		})

		It("Doesn't work without old controller signature", func() {
			signatures := []SignInput{bob.SignInput}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", alice.Did)))
		})

		It("Doesn't work without new controller signature", func() {
			signatures := []SignInput{alice.SignInput}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one signature by %s: signature is required but not found", bob.Did)))
		})
	})

	Describe("Verification method: id update", func() {
		var alice CreatedDidInfo
		var newKeyId string
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			newKeyId = alice.Did + "#key-2"

			msg = &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   newKeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           alice.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				Authentication: []string{alice.KeyId},
				VersionId:      alice.VersionId,
			}
		})

		It("Doesn't work without new VM signature", func() {
			signatures := []SignInput{alice.SignInput}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (new version): invalid signature detected", alice.Did)))
		})

		It("Doesn't work without old VM signature", func() {
			signatures := []SignInput{
				{
					VerificationMethodId: newKeyId,
					Key:                  alice.KeyPair.Private,
				},
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", alice.Did)))
		})

		It("Works with new and old VM signature", func() {
			signatures := []SignInput{
				{
					VerificationMethodId: newKeyId,
					Key:                  alice.KeyPair.Private,
				},
				alice.SignInput,
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err).To(BeNil())

			// check
			updated, err := setup.QueryDid(alice.Did)
			Expect(err).To(BeNil())
			Expect(*updated.Value.DidDoc).To(Equal(msg.ToDidDoc()))
		})
	})

	Describe("Verification method: adding a new one", func() {
		var alice CreatedDidInfo
		var newKeyId string
		var newKey KeyPair
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()

			newKeyId = alice.Did + "#key-2"
			newKey = GenerateKeyPair()

			msg = &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   alice.KeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           alice.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
					{
						Id:                   newKeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           alice.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(newKey.Public),
					},
				},
				Authentication: []string{alice.KeyId},
				VersionId:      alice.VersionId,
			}
		})

		It("Works with only old VM signature", func() {
			signatures := []SignInput{
				alice.SignInput,
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryDid(alice.Did)
			Expect(err).To(BeNil())
			Expect(*created).ToNot(Equal(msg.ToDidDoc()))
		})

		It("Doesn't work with only new VM signature", func() {
			signatures := []SignInput{
				{
					VerificationMethodId: newKeyId,
					Key:                  newKey.Private,
				},
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (old version): invalid signature detected", alice.Did)))
		})
	})

	Describe("Verification method: removing existing one", func() {
		var alice CreatedDidInfo
		var secondKeyId string
		var secondKey KeyPair
		var secondSignInput SignInput
		var msg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()

			secondKeyId = alice.Did + "#key-2"
			secondKey = GenerateKeyPair()
			secondSignInput = SignInput{
				VerificationMethodId: secondKeyId,
				Key:                  secondKey.Private,
			}

			addSecondKeyMsg := &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   alice.KeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           alice.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
					{
						Id:                   secondKeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           alice.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(secondKey.Public),
					},
				},
				Authentication: []string{alice.KeyId},
				VersionId:      alice.VersionId,
			}

			_, err := setup.UpdateDid(addSecondKeyMsg, []SignInput{alice.SignInput})
			Expect(err).To(BeNil())

			msg = &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   alice.KeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           alice.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(alice.KeyPair.Public),
					},
				},
				Authentication: []string{alice.KeyId},
				VersionId:      alice.VersionId,
			}
		})

		It("Works with only first VM signature", func() {
			signatures := []SignInput{
				alice.SignInput,
			}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err).To(BeNil())

			// check
			created, err := setup.QueryDid(alice.Did)
			Expect(err).To(BeNil())
			Expect(*created).ToNot(Equal(msg.ToDidDoc()))
		})

		It("Doesn't work with only second VM signature (which is get deleted)", func() {
			signatures := []SignInput{secondSignInput}

			_, err := setup.UpdateDid(msg, signatures)
			Expect(err.Error()).To(ContainSubstring(fmt.Sprintf("there should be at least one valid signature by %s (new version): invalid signature detected", alice.Did)))
		})
	})

	Describe("Deactivating", func() {
		var alice CreatedDidInfo
		var bob CreatedDidInfo
		var updateMsg *types.MsgUpdateDidDocPayload

		BeforeEach(func() {
			alice = setup.CreateSimpleDid()
			bob = setup.CreateSimpleDid()

			updateMsg = &types.MsgUpdateDidDocPayload{
				Id: alice.Did,
				VerificationMethod: []*types.VerificationMethod{
					{
						Id:                   alice.DidInfo.KeyId,
						Type:                 types.Ed25519VerificationKey2020{}.Type(),
						Controller:           alice.DidInfo.Did,
						VerificationMaterial: BuildEd25519VerificationKey2020VerificationMaterial(alice.DidInfo.KeyPair.Public),
					},
				},
				Authentication: []string{alice.KeyId},
				VersionId:      alice.VersionId,
			}
		})

		When("Updating already deactivated DID", func() {
			It("Should fail with error", func() {
				// Deactivate DID
				deactivateMsg := &types.MsgDeactivateDidDocPayload{
					Id: alice.Did,
				}

				signatures := []SignInput{alice.DidInfo.SignInput}

				res, err := setup.DeactivateDid(deactivateMsg, signatures)
				Expect(err).To(BeNil())
				Expect(res.Value.Metadata.Deactivated).To(BeTrue())

				// Update deactivated DID
				signatures = []SignInput{
					alice.SignInput,
					bob.SignInput,
				}

				_, err = setup.UpdateDid(updateMsg, signatures)
				Expect(err.Error()).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring(alice.DidInfo.Did + ": DID Doc already deactivated"))
			})
		})
	})
})
