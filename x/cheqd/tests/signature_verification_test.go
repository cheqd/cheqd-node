package tests

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

var _ = Describe("Signature Verification while updating DID", func() {
	var setup TestSetup
	var aliceKeys, bobKeys map[string]ed25519.PrivateKey
	var aliceDid *types.MsgCreateDidPayload
	BeforeEach(func() {
		setup = Setup()
		aliceKeys, aliceDid, _ = setup.InitDid(AliceDID)
		bobKeys, _, _ = setup.InitDid(BobDID)
	})

	It("should have changed DIDDoc controller", func() {
		updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
		updatedDidDoc.Controller = append(updatedDidDoc.Controller, BobDID)
		receivedDid, _ := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(ConcatKeys(aliceKeys, bobKeys)))

		Expect(aliceDid.Controller).To(Not(Equal(receivedDid.Controller)))
		Expect([]string{AliceDID, BobDID}).To(Not(Equal(receivedDid.Controller)))
		Expect([]string{BobDID}, receivedDid.Controller)
	})

	When("Old signature in verification method is absent", func() {
		It("should fail", func() {
			updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
			updatedDidDoc.VerificationMethod[0].Type = types.Ed25519VerificationKey2020
			_, err := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(bobKeys))

			// check
			Expect(err).To(Not(BeNil()))
			Expect(err.Error()).To(Equal(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)))
		})
	})

	It("should fails cause we need old signature for changing verification method controller", func() {
		updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
		updatedDidDoc.VerificationMethod[0].Controller = BobDID
		_, err := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(bobKeys))

		// check
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(Equal(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)))
	})

	It("should fails cause we need old signature for changing DIDDoc controller", func() {
		updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
		updatedDidDoc.Controller = append(updatedDidDoc.Controller, BobDID)
		_, err := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(bobKeys))

		// check
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(Equal(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)))
	})
})

var _ = Describe("Signature Verification. Remove signature/VM", func() {
	var setup TestSetup
	var ApubKey, BpubKey ed25519.PublicKey
	var AprivKey, BprivKey ed25519.PrivateKey
	var aliceDid, bobDid *types.MsgCreateDidPayload
	var aliceKeys, bobKeys map[string]ed25519.PrivateKey

	BeforeEach(func() {
		setup = Setup()
		// Generate keys
		ApubKey, AprivKey, _ = ed25519.GenerateKey(rand.Reader)
		BpubKey, BprivKey, _ = ed25519.GenerateKey(rand.Reader)

		// Create dids
		aliceDid = setup.BuildMsgCreateDidPayload(AliceDID, AliceKey1, ApubKey)
		bobDid = setup.BuildMsgCreateDidPayload(BobDID, BobKey1, BpubKey)

		// Collect private keys
		aliceKeys = map[string]ed25519.PrivateKey{AliceKey1: AprivKey, BobKey1: BprivKey}
		bobKeys = map[string]ed25519.PrivateKey{BobKey1: BprivKey}

		// Add verification method
		aliceDid.VerificationMethod = append(aliceDid.VerificationMethod, &types.VerificationMethod{
			Id:                 AliceKey2,
			Controller:         BobDID,
			Type:               types.Ed25519VerificationKey2020,
			PublicKeyMultibase: MustEncodeBase58(BpubKey),
		})
	})

	It("should fails cause old signature is required for removing this signature", func() {
		// Send dids
		_, _ = setup.SendCreateDid(bobDid, bobKeys)
		_, _ = setup.SendCreateDid(aliceDid, aliceKeys)

		updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
		updatedDidDoc.VerificationMethod = []*types.VerificationMethod{aliceDid.VerificationMethod[0]}
		updatedDidDoc.Authentication = []string{aliceDid.Authentication[0]}
		_, err := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(bobKeys))

		// check
		Expect(err).To(Not(BeNil()))
		Expect(err.Error()).To(Equal(fmt.Sprintf("there should be at least one signature by %s (old version): signature is required but not found", AliceDID)))
	})

	It("should not fails while removing the whole verification method", func() {
		aliceDid.Authentication = append(aliceDid.Authentication, AliceKey2)

		// Send dids
		_, _ = setup.SendCreateDid(bobDid, bobKeys)
		_, _ = setup.SendCreateDid(aliceDid, aliceKeys)

		updatedDidDoc := setup.CreateToUpdateDid(aliceDid)
		updatedDidDoc.Authentication = []string{aliceDid.Authentication[0]}
		updatedDidDoc.VerificationMethod = []*types.VerificationMethod{aliceDid.VerificationMethod[0]}
		receivedDid, _ := setup.SendUpdateDid(updatedDidDoc, MapToListOfSignerKeys(ConcatKeys(aliceKeys, bobKeys)))

		// check
		Expect(len(aliceDid.VerificationMethod)).To(Not(Equal(len(receivedDid.VerificationMethod))))
		Expect(reflect.DeepEqual(aliceDid.VerificationMethod[0], receivedDid.VerificationMethod[0])).To(BeTrue())
	})
})
