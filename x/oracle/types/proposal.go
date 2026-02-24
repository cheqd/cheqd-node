package types

import (
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

var (
	proposalTypeMsgLegacyGovUpdateParams                = MsgLegacyGovUpdateParams{}.String()
	proposalTypeMsgGovUpdateParams                      = MsgGovUpdateParams{}.String()
	proposalTypeMsgGovCancelUpdateParams                = MsgGovCancelUpdateParamPlan{}.String()
	proposalTypeMsgGovAddDenoms                         = MsgGovAddDenoms{}.String()
	proposalTypeMsgGovRemoveCurrencyPairProviders       = MsgGovRemoveCurrencyPairProviders{}.String()
	proposalTypeMsgGovRemoveCurrencyDeviationThresholds = MsgGovRemoveCurrencyDeviationThresholds{}.String()
)

func init() {
	gov.RegisterProposalType(proposalTypeMsgLegacyGovUpdateParams)
	gov.RegisterProposalType(proposalTypeMsgGovUpdateParams)
	gov.RegisterProposalType(proposalTypeMsgGovCancelUpdateParams)
	gov.RegisterProposalType(proposalTypeMsgGovAddDenoms)
	gov.RegisterProposalType(proposalTypeMsgGovRemoveCurrencyPairProviders)
	gov.RegisterProposalType(proposalTypeMsgGovRemoveCurrencyDeviationThresholds)
}

// Implements Proposal Interface
var _ gov.Content = &MsgLegacyGovUpdateParams{}

// GetTitle returns the title of a community pool spend proposal.
func (msg *MsgLegacyGovUpdateParams) GetTitle() string { return msg.Title }

// GetDescription returns the description of a community pool spend proposal.
func (msg *MsgLegacyGovUpdateParams) GetDescription() string { return msg.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (msg *MsgLegacyGovUpdateParams) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (msg *MsgLegacyGovUpdateParams) ProposalType() string {
	return proposalTypeMsgLegacyGovUpdateParams
}

// Implements Proposal Interface
var _ gov.Content = &MsgGovUpdateParams{}

// GetTitle returns the title of a community pool spend proposal.
func (msg *MsgGovUpdateParams) GetTitle() string { return msg.Title }

// GetDescription returns the description of a community pool spend proposal.
func (msg *MsgGovUpdateParams) GetDescription() string { return msg.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (msg *MsgGovUpdateParams) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (msg *MsgGovUpdateParams) ProposalType() string { return proposalTypeMsgGovUpdateParams }

// Implements Proposal Interface
var _ gov.Content = &MsgGovCancelUpdateParamPlan{}

// GetTitle returns the title of a community pool spend proposal.
func (msg *MsgGovCancelUpdateParamPlan) GetTitle() string { return msg.Title }

// GetDescription returns the description of a community pool spend proposal.
func (msg *MsgGovCancelUpdateParamPlan) GetDescription() string { return msg.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (msg *MsgGovCancelUpdateParamPlan) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (msg *MsgGovCancelUpdateParamPlan) ProposalType() string {
	return proposalTypeMsgGovCancelUpdateParams
}

// Implements Proposal Interface
var _ gov.Content = &MsgGovAddDenoms{}

// GetTitle returns the title of a community pool spend proposal.
func (msg *MsgGovAddDenoms) GetTitle() string { return msg.Title }

// GetDescription returns the description of a community pool spend proposal.
func (msg *MsgGovAddDenoms) GetDescription() string { return msg.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (msg *MsgGovAddDenoms) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (msg *MsgGovAddDenoms) ProposalType() string {
	return proposalTypeMsgGovAddDenoms
}

// Implements Proposal Interface
var _ gov.Content = &MsgGovRemoveCurrencyPairProviders{}

// GetTitle returns the title of a community pool spend proposal.
func (msg *MsgGovRemoveCurrencyPairProviders) GetTitle() string { return msg.Title }

// GetDescription returns the description of a community pool spend proposal.
func (msg *MsgGovRemoveCurrencyPairProviders) GetDescription() string { return msg.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (msg *MsgGovRemoveCurrencyPairProviders) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (msg *MsgGovRemoveCurrencyPairProviders) ProposalType() string {
	return proposalTypeMsgGovRemoveCurrencyPairProviders
}

// Implements Proposal Interface
var _ gov.Content = &MsgGovRemoveCurrencyDeviationThresholds{}

// GetTitle returns the title of a community pool spend proposal.
func (msg *MsgGovRemoveCurrencyDeviationThresholds) GetTitle() string { return msg.Title }

// GetDescription returns the description of a community pool spend proposal.
func (msg *MsgGovRemoveCurrencyDeviationThresholds) GetDescription() string { return msg.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (msg *MsgGovRemoveCurrencyDeviationThresholds) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (msg *MsgGovRemoveCurrencyDeviationThresholds) ProposalType() string {
	return proposalTypeMsgGovRemoveCurrencyDeviationThresholds
}
