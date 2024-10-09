//go:build integration

package integration

import (
	"fmt"
	"path/filepath"

	cli "github.com/cheqd/cheqd-node/tests/integration/cli"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration - Mint coins to given address", func() {
	var Proposal_id string
	It("should wait for node catching up", func() {
		By("pinging the node status until catching up is flagged as false ")
		err := cli.WaitForCaughtUp(cli.Validator0, cli.CliBinaryName, cli.VotingPeriod*6)
		Expect(err).To(BeNil())
	})
	It("should submit a mint proposal ", func() {
		By("passing the proposal file to the container")
		_, err := cli.LocalnetExecCopyAbsoluteWithPermissions(filepath.Join("proposal.json"), cli.DockerHome, cli.Validator0)
		Expect(err).To(BeNil())

		By("sending a SubmitProposal transaction from `validator0` container")
		res, err := cli.SubmitProposalTx(cli.Operator0, "proposal.json", cli.CliGasParams)
		Expect(err).To(BeNil())
		res, err = cli.QueryTxn(res.TxHash)
		Expect(err).To(BeNil())

		proposal_id, err := cli.GetProposalID(res.RawLog)

		fmt.Println("The proposal id>>>>>>>>>>>>>>>>>>>>>>>>>>>.", proposal_id)
		Proposal_id = proposal_id
		Expect(err).To(BeNil())

		Expect(res.Code).To(BeEquivalentTo(0))
	})
	It("keys list", func() {
		By("keys list in validator 0")
		keys, err := cli.QueryKeysList()
		Expect(err).To(BeNil())
		fmt.Println("keys arE>>>>>>>>>>>>>>>>>", keys)
	})

	// It("should wait for the proposal submission to be included in a block", func() {
	// 	By("getting the current block height")
	// 	currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
	// 	Expect(err).To(BeNil())

	// 	By("waiting for the proposal to be included in a block")
	// 	err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+1, 2)
	// 	Expect(err).To(BeNil())
	// })

	It("should vote for the mint proposal from `validator1` container", func() {
		By("sending a VoteProposal transaction from `validator1` container")
		res, err := cli.VoteProposalTx(cli.Operator1, Proposal_id, "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		fmt.Println("proposal_id>>>>>>>>>>Here>>>>>", Proposal_id)
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the mint proposal from `validator2` container", func() {
		By("sending a VoteProposal transaction from `validator2` container")
		res, err := cli.VoteProposalTx(cli.Operator2, Proposal_id, "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should vote for the mint proposal from `validator0` container", func() {
		By("sending a VoteProposal transaction from `validator0` container")
		res, err := cli.VoteProposalTx(cli.Operator0, Proposal_id, "yes", cli.CliGasParams)
		Expect(err).To(BeNil())
		Expect(res.Code).To(BeEquivalentTo(0))
	})

	It("should wait for the proposal to pass", func() {
		By("getting the current block height")
		currentHeight, err := cli.GetCurrentBlockHeight(cli.Validator0, cli.CliBinaryName)
		Expect(err).To(BeNil())

		By("waiting for the proposal to pass")
		err = cli.WaitForChainHeight(cli.Validator0, cli.CliBinaryName, currentHeight+20, 25)
		Expect(err).To(BeNil())
	})

	It("should check the proposal status to ensure it has passed", func() {
		By("sending a QueryProposal query from `validator0` container")
		proposal, err := cli.QueryProposal(cli.Validator0, Proposal_id)
		Expect(err).To(BeNil())
		Expect(proposal.Status).To(BeEquivalentTo(govtypesv1.StatusPassed))
	})

	It("should have the correct balance after minting", func() {
		By("querying the balance of the given address")
		bal, err := cli.QueryBalance("cheqd1lhl9g4rgldadgtz7v6rt50u45uhhj8hhv8d8uf", didtypes.BaseMinimalDenom)
		Expect(err).To(BeNil())
		Expect(bal.Amount.Int64()).To(Equal(int64(9000)))
	})
})
