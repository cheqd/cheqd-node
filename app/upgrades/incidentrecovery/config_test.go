package incidentrecovery

import (
	"testing"

	"github.com/cosmos/btcutil/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// These tests exist specifically to prove the Juno-Prop-21-style failure
// mode (shipping with a wrong/placeholder/unreachable destination address)
// is actually caught by the code, not just asserted in a comment.

const (
	wellFormedDest = "cheqd1ufqwl05g8pcyrq92rdzx05ra0yvkqlgwpd5yzh"
	wellFormedAff1 = "cheqd1pav86kd046tv7em6nc2dgcycxdppdnyzu2xycv"
	wellFormedAff2 = "cheqd18wm3f6s7yw0f0la09rnczvu9xp8wkkxfue8j5m"
)

func TestValidate_RejectsPlaceholderDestination(t *testing.T) {
	cfg := &Config{
		UpgradeName:               "v4.3.0-incident-recovery",
		DestinationAddress:        placeholderDestination,
		DestinationAddressConfirm: placeholderDestination,
		AffectedAddresses:         []string{wellFormedAff1},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for placeholder destination_address, got nil")
	}
}

func TestValidate_RejectsNonCheqdDestination(t *testing.T) {
	// a validly-formatted bech32 address, but wrong HRP (an osmo address)
	addr := "osmo18ty4lfjecjdrmn8j3yhf3jyxrmtc0adz8lyv6y"
	cfg := &Config{
		UpgradeName:               "v4.3.0-incident-recovery",
		DestinationAddress:        addr,
		DestinationAddressConfirm: addr,
		AffectedAddresses:         []string{wellFormedAff1},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for non-cheqd-prefixed destination_address, got nil")
	}
}

func TestValidate_RejectsEmptyAffectedAddresses(t *testing.T) {
	cfg := &Config{
		UpgradeName:               "v4.3.0-incident-recovery",
		DestinationAddress:        wellFormedDest,
		DestinationAddressConfirm: wellFormedDest,
		AffectedAddresses:         []string{},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty affected_addresses, got nil")
	}
}

func TestValidate_RejectsDestinationAlsoBeingAffected(t *testing.T) {
	cfg := &Config{
		UpgradeName:               "v4.3.0-incident-recovery",
		DestinationAddress:        wellFormedDest,
		DestinationAddressConfirm: wellFormedDest,
		AffectedAddresses:         []string{wellFormedDest},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when destination_address duplicates an affected address, got nil")
	}
}

func TestValidate_AcceptsWellFormedConfig(t *testing.T) {
	cfg := &Config{
		UpgradeName:               "v4.3.0-incident-recovery",
		DestinationAddress:        wellFormedDest,
		DestinationAddressConfirm: wellFormedDest,
		AffectedAddresses:         []string{wellFormedAff1, wellFormedAff2},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected well-formed config to validate, got error: %v", err)
	}
}

// --- new tests below, for the two gaps fixed in this review ---

func TestValidate_RejectsMismatchedDestinationConfirm(t *testing.T) {
	cfg := &Config{
		UpgradeName:               "v4.3.0-incident-recovery",
		DestinationAddress:        wellFormedDest,
		DestinationAddressConfirm: wellFormedAff1, // deliberately different, well-formed address
		AffectedAddresses:         []string{wellFormedAff2},
	}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error when destination_address and destination_address_confirm differ, got nil")
	}
}

func TestValidate_RejectsUnsetDestinationConfirm(t *testing.T) {
	cfg := &Config{
		UpgradeName:        "v4.3.0-incident-recovery",
		DestinationAddress: wellFormedDest,
		// DestinationAddressConfirm deliberately left zero-value
		AffectedAddresses: []string{wellFormedAff1},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when destination_address_confirm is unset, got nil")
	}
}

// TestValidate_RejectsModuleAccountDestination proves the module-account
// guard actually fires against a real derived module address, re-encoded
// with the cheqd HRP — not just that the code compiles. Catches exactly
// the class of mistake where someone pastes an address that is valid but
// unreachable (no one holds a private key for it), the same root failure
// as Juno's Prop 21 incident, just a different specific cause.
func TestValidate_RejectsModuleAccountDestination(t *testing.T) {
	moduleAddr := authtypes.NewModuleAddress(stakingtypes.BondedPoolName)
	data, err := bech32.ConvertBits(moduleAddr.Bytes(), 8, 5, true)
	if err != nil {
		t.Fatalf("test setup: converting bits: %v", err)
	}
	cheqdEncodedModuleAddr, err := bech32.Encode(bech32HRP, data)
	if err != nil {
		t.Fatalf("test setup: encoding bech32: %v", err)
	}

	cfg := &Config{
		UpgradeName:               "v4.3.0-incident-recovery",
		DestinationAddress:        cheqdEncodedModuleAddr,
		DestinationAddressConfirm: cheqdEncodedModuleAddr,
		AffectedAddresses:         []string{wellFormedAff1},
	}
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("expected error for a destination address that is actually the %q module account (%s), got nil", stakingtypes.BondedPoolName, cheqdEncodedModuleAddr)
	}
}

func TestValidate_AcceptsWellFormedConfig_NonModuleDestination(t *testing.T) {
	// Sanity check that the module-account guard doesn't over-fire on an
	// ordinary address — wellFormedDest must NOT look like any known
	// module account, or this (and the earlier "accepts well-formed
	// config" test) would be a false negative.
	if err := rejectModuleAccountAddress(wellFormedDest); err != nil {
		t.Fatalf("expected wellFormedDest to NOT be flagged as a module account, got: %v", err)
	}
}
