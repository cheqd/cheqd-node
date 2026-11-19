package incidentrecovery

import (
	"encoding/json"
	"strings"
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

// TestValidate_RejectsModuleAccountInAffectedAddresses is the same guard as
// TestValidate_RejectsModuleAccountDestination, but for affected_addresses —
// third-pass review finding: the original code only ever applied the
// module-account check to destination_address, so a typo'd or malicious
// affected_addresses entry could sweep a module account's (e.g. the fee
// collector's) balance to the destination instead of recovering a real
// compromised wallet.
func TestValidate_RejectsModuleAccountInAffectedAddresses(t *testing.T) {
	moduleAddr := authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName)
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
		DestinationAddress:        wellFormedDest,
		DestinationAddressConfirm: wellFormedDest,
		AffectedAddresses:         []string{wellFormedAff1, cheqdEncodedModuleAddr},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected error for an affected address that is actually the %q module account (%s), got nil", stakingtypes.NotBondedPoolName, cheqdEncodedModuleAddr)
	}
}

// TestValidate_RejectsMixedCaseAddress covers a third-pass review finding:
// bech32.Decode accepts an all-uppercase address and silently lowercases it
// before decoding, but every equality/duplicate check in Validate compares
// the raw strings as given. An address and its uppercase equivalent decode
// to the identical account but would NOT be caught as the same address by
// those checks (e.g. duplicate affected addresses, or destination ==
// affected). Requiring lowercase input closes that gap.
func TestValidate_RejectsMixedCaseAddress(t *testing.T) {
	upper := strings.ToUpper(wellFormedAff1)
	cfg := &Config{
		UpgradeName:               "v4.3.0-incident-recovery",
		DestinationAddress:        wellFormedDest,
		DestinationAddressConfirm: wellFormedDest,
		AffectedAddresses:         []string{upper},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected an uppercase address (%s) to be rejected even though it bech32-decodes successfully", upper)
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

// --- ParsePlanInfo tests ---
//
// The config travels in the governance-approved plan.Info, so every node
// evaluates byte-identical input. These tests pin the accepted shape and
// make sure anything malformed fails loudly rather than decoding as an
// empty or partial config.

func planInfoJSON(t *testing.T, cfg *Config, extra map[string]any) string {
	t.Helper()
	top := map[string]any{PlanInfoKey: cfg}
	for k, v := range extra {
		top[k] = v
	}
	b, err := json.Marshal(top)
	if err != nil {
		t.Fatalf("test setup: marshalling plan info: %v", err)
	}
	return string(b)
}

func TestParsePlanInfo_AcceptsConfigWithCosmovisorBinaries(t *testing.T) {
	info := planInfoJSON(t, testConfig(wellFormedAff1, wellFormedAff2), map[string]any{
		"binaries": map[string]string{
			"linux/amd64": "https://example.com/cheqd-noded-linux-amd64.tar.gz?checksum=sha256:00",
		},
	})
	cfg, err := ParsePlanInfo(info)
	if err != nil {
		t.Fatalf("expected plan info with binaries to be accepted, got: %v", err)
	}
	if cfg.DestinationAddress != wellFormedDest || len(cfg.AffectedAddresses) != 2 {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestParsePlanInfo_AcceptsConfigWithoutBinaries(t *testing.T) {
	if _, err := ParsePlanInfo(planInfoJSON(t, testConfig(wellFormedAff1), nil)); err != nil {
		t.Fatalf("expected plan info without binaries to be accepted, got: %v", err)
	}
}

func TestParsePlanInfo_Rejects(t *testing.T) {
	valid := planInfoJSON(t, testConfig(wellFormedAff1), nil)
	cases := map[string]string{
		"empty":                     "",
		"whitespace":                "   ",
		"old bare sha256 format":    "7d07abb1203403ee43eecd9a8e624d677358a3f972e1d1c0527234fbf0a30582",
		"url":                       "https://example.com/upgrade-info.json",
		"json array":                `[1,2]`,
		"missing incident_recovery": `{"binaries":{}}`,
		"unknown top-level key":     strings.Replace(valid, `"incident_recovery"`, `"incident_recovery_typo"`, 1),
		"extra top-level key":       strings.TrimSuffix(valid, "}") + `,"notes":"x"}`,
		"unknown config field":      strings.Replace(valid, `"affected_addresses"`, `"affected_adresses"`, 1),
		"trailing data":             valid + `{}`,
		"trailing bracket":          valid + `]`,
		"placeholder destination":   planInfoJSON(t, &Config{UpgradeName: UpgradeName, DestinationAddress: placeholderDestination, DestinationAddressConfirm: placeholderDestination, AffectedAddresses: []string{wellFormedAff1}}, nil),
		"destination mismatch":      planInfoJSON(t, &Config{UpgradeName: UpgradeName, DestinationAddress: wellFormedDest, DestinationAddressConfirm: wellFormedAff2, AffectedAddresses: []string{wellFormedAff1}}, nil),
		"wrong upgrade name":        planInfoJSON(t, &Config{UpgradeName: "v9", DestinationAddress: wellFormedDest, DestinationAddressConfirm: wellFormedDest, AffectedAddresses: []string{wellFormedAff1}}, nil),
		"empty affected list":       planInfoJSON(t, &Config{UpgradeName: UpgradeName, DestinationAddress: wellFormedDest, DestinationAddressConfirm: wellFormedDest}, nil),
	}
	for name, info := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := ParsePlanInfo(info); err == nil {
				t.Fatalf("expected plan info to be rejected, got nil error for: %s", info)
			}
		})
	}
}
