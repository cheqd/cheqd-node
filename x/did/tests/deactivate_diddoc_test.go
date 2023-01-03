package tests

import (
	"fmt"

	"github.com/google/uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	testsetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	"github.com/cheqd/cheqd-node/x/did/types"
)

var _ = Describe("Deactivate DID tests", func() {
	var setup testsetup.TestSetup

	BeforeEach(func() {
		setup = testsetup.Setup()
	})

	It("Valid: Deactivate DID", func() {
		// Alice
		alice := setup.CreateSimpleDid()
		msg := &types.MsgDeactivateDidDocPayload{
			Id:        alice.Did,
			VersionId: uuid.NewString(),
		}

		signatures := []testsetup.SignInput{alice.DidDocInfo.SignInput}

		res, err := setup.DeactivateDidDoc(msg, signatures)
		Expect(err).To(BeNil())
		Expect(res.Value.Metadata.Deactivated).To(BeTrue())

		// Check that all versions are deactivated
		versions, err := setup.QueryAllDidDocVersionsMetadata(alice.Did)
		Expect(err).To(BeNil())

		for _, version := range versions.Versions {
			Expect(version.Deactivated).To(BeTrue())
		}
	})

	When("DID is not found", func() {
		It("Should return error", func() {
			NotFoundDID := testsetup.GenerateDID(testsetup.Base58_16bytes)

			msg := &types.MsgDeactivateDidDocPayload{
				Id:        NotFoundDID,
				VersionId: uuid.NewString(),
			}

			signatures := []testsetup.SignInput{}

			_, err := setup.DeactivateDidDoc(msg, signatures)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(NotFoundDID + ": DID Doc not found"))
		})
	})

	When("DID is already deactivated", func() {
		It("Should return error", func() {
			// Alice
			alice := setup.CreateSimpleDid()
			msg := &types.MsgDeactivateDidDocPayload{
				Id:        alice.Did,
				VersionId: uuid.NewString(),
			}

			signatures := []testsetup.SignInput{alice.DidDocInfo.SignInput}

			res, err := setup.DeactivateDidDoc(msg, signatures)
			Expect(err).To(BeNil())
			Expect(res.Value.Metadata.Deactivated).To(BeTrue())

			// Deactivate again
			_, err = setup.DeactivateDidDoc(msg, signatures)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(alice.DidDocInfo.Did + ": DID Doc already deactivated"))
		})
	})

	When("Signatures are invalid", func() {
		It("Should return an error", func() {
			// Alice
			alice := setup.CreateSimpleDid()
			// Bob
			bob := setup.CreateSimpleDid()

			msg := &types.MsgDeactivateDidDocPayload{
				Id:        alice.Did,
				VersionId: uuid.NewString(),
			}

			signatures := []testsetup.SignInput{bob.DidDocInfo.SignInput}

			_, err := setup.DeactivateDidDoc(msg, signatures)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(fmt.Sprintf("signer: %s: signature is required but not found", alice.DidDocInfo.Did)))
		})
	})
})
