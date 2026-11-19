package incidentrecovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos/btcutil/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Config is the configurable input to the incident-recovery upgrade handler.
// It is carried inside the governance proposal itself, in the upgrade
// plan's Info field (see ParsePlanInfo) — nothing here is hardcoded in Go
// source, and nothing is read from any node's local disk.
type Config struct {
	UpgradeName        string `json:"upgrade_name"`
	DestinationAddress string `json:"destination_address"`
	// DestinationAddressConfirm must be typed/pasted independently and must
	// match DestinationAddress exactly, byte for byte. This exists
	// specifically because a single unconfirmed destination field is what
	// let Juno's Prop 21 copy-paste error through undetected. Two people
	// (or the same person, twice, from two independent sources) agreeing on
	// the exact same string is a much stronger check than one field typed
	// once.
	DestinationAddressConfirm string   `json:"destination_address_confirm"`
	AffectedAddresses         []string `json:"affected_addresses"`
}

const placeholderDestination = "PLACEHOLDER_UNVERIFIED_DO_NOT_USE"

const bech32HRP = "cheqd"

// PlanInfoKey is the top-level key of the upgrade plan's Info JSON that
// holds the recovery Config.
const PlanInfoKey = "incident_recovery"

// planInfoBinariesKey is the key Cosmovisor reads from the same Info JSON to
// auto-download the upgrade binary. It is the only other top-level key
// allowed, so a typo in PlanInfoKey fails loudly instead of being ignored.
const planInfoBinariesKey = "binaries"

// ParsePlanInfo extracts and validates the recovery Config from a
// software-upgrade plan's Info string. The expected shape is:
//
//	{
//	  "binaries": { "linux/amd64": "https://...?checksum=sha256:..." },
//	  "incident_recovery": {
//	    "upgrade_name": "v4.3.0-incident-recovery",
//	    "destination_address": "cheqd1...",
//	    "destination_address_confirm": "cheqd1...",
//	    "affected_addresses": ["cheqd1...", "..."]
//	  }
//	}
//
// "binaries" is optional and ignored here; Cosmovisor uses it (and ignores
// "incident_recovery") to download the binary.
//
// Carrying the config in Plan.Info rather than in a file on each node's
// disk means every node evaluates byte-identical input, because Plan.Info
// came through consensus. Voters approve the actual addresses and
// destination, not just an upgrade name or a hash of a file they may never
// see. And no operator has to place a file correctly for the chain to
// upgrade.
//
// Decoding is strict: unknown top-level keys, unknown fields inside
// "incident_recovery" and trailing data are all rejected. A typo such as
// "affected_adresses" would otherwise decode as an empty list or be silently
// dropped. Validate then applies the same checks as before, so the
// placeholder destination, a destination/confirm mismatch, module-account
// addresses, duplicates and non-cheqd addresses are all refused.
func ParsePlanInfo(info string) (*Config, error) {
	info = strings.TrimSpace(info)
	if info == "" {
		return nil, fmt.Errorf("upgrade plan info is empty: it must carry the %q config JSON that voters approved", PlanInfoKey)
	}

	var top map[string]json.RawMessage
	if err := decodeStrict([]byte(info), &top); err != nil {
		return nil, fmt.Errorf("upgrade plan info is not a JSON object: %w", err)
	}
	for k := range top {
		if k != PlanInfoKey && k != planInfoBinariesKey {
			return nil, fmt.Errorf("upgrade plan info has unexpected top-level key %q (allowed: %q, %q)", k, PlanInfoKey, planInfoBinariesKey)
		}
	}
	raw, ok := top[PlanInfoKey]
	if !ok {
		return nil, fmt.Errorf("upgrade plan info has no %q key", PlanInfoKey)
	}

	var cfg Config
	if err := decodeStrict(raw, &cfg); err != nil {
		return nil, fmt.Errorf("upgrade plan info %q: %w", PlanInfoKey, err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("upgrade plan info %q: invalid config: %w", PlanInfoKey, err)
	}
	return &cfg, nil
}

// decodeStrict decodes exactly one JSON value into v, rejecting unknown
// struct fields and any trailing data.
func decodeStrict(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if rest := bytes.TrimSpace(data[dec.InputOffset():]); len(rest) > 0 {
		return fmt.Errorf("unexpected data after the JSON value")
	}
	return nil
}

// Validate checks the config is safe to use. Called by ParsePlanInfo, and
// again by the handler immediately before it executes, as a second
// independent check at the point of use.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.UpgradeName) == "" {
		return fmt.Errorf("upgrade_name must be set")
	}
	if c.UpgradeName != UpgradeName {
		// Catches a config drifted from the registered handler name — e.g.
		// an operator edits upgrade_name in the YAML without realizing it
		// must match the compiled-in name app.go actually registered via
		// SetUpgradeHandler. Mismatched, the upgrade would fire (the chain
		// halts at the right height) but this handler would never run at
		// all — a silent no-op, not a loud failure, which is worse.
		return fmt.Errorf("upgrade_name %q does not match the compiled-in UpgradeName %q — these must match exactly or this handler will never actually run when the upgrade fires", c.UpgradeName, UpgradeName)
	}

	if c.DestinationAddress == "" || c.DestinationAddress == placeholderDestination {
		return fmt.Errorf("destination_address is still the placeholder (%q) — this must be replaced with a real, independently-verified address before this config can be used for anything", placeholderDestination)
	}

	if c.DestinationAddressConfirm == "" || c.DestinationAddressConfirm == placeholderDestination {
		return fmt.Errorf("destination_address_confirm is still unset/placeholder — set it to the same value as destination_address, entered independently, as a transcription check")
	}

	if c.DestinationAddress != c.DestinationAddressConfirm {
		return fmt.Errorf("destination_address (%q) and destination_address_confirm (%q) do not match — refusing, this is exactly the kind of mismatch that must be caught before it reaches an upgrade handler", c.DestinationAddress, c.DestinationAddressConfirm)
	}

	if err := validateCheqdAddress(c.DestinationAddress); err != nil {
		return fmt.Errorf("destination_address: %w", err)
	}

	if err := rejectModuleAccountAddress(c.DestinationAddress); err != nil {
		return fmt.Errorf("destination_address: %w", err)
	}

	if len(c.AffectedAddresses) == 0 {
		return fmt.Errorf("affected_addresses must not be empty")
	}

	seen := make(map[string]bool, len(c.AffectedAddresses))
	for i, addr := range c.AffectedAddresses {
		if err := validateCheqdAddress(addr); err != nil {
			return fmt.Errorf("affected_addresses[%d] (%q): %w", i, addr, err)
		}
		// A module account address in affected_addresses would be a typo or
		// a malicious draft, not a real compromised wallet — module
		// accounts aren't held by anyone and have no attacker to recover
		// from. Sweeping one would move protocol funds (e.g. the fee
		// collector or bonded/not-bonded pools) to the destination instead.
		if err := rejectModuleAccountAddress(addr); err != nil {
			return fmt.Errorf("affected_addresses[%d] (%q): %w", i, addr, err)
		}
		if addr == c.DestinationAddress {
			return fmt.Errorf("affected_addresses[%d] (%q) is the same as destination_address — refusing, this is almost certainly a config error", i, addr)
		}
		if seen[addr] {
			return fmt.Errorf("affected_addresses[%d] (%q) is a duplicate", i, addr)
		}
		seen[addr] = true
	}

	return nil
}

// validateCheqdAddress checks that addr bech32-decodes with the "cheqd" HRP.
// Deliberately uses the raw bech32 decoder rather than the cosmos-sdk
// AccAddress helper, so this check doesn't depend on (or silently pass due
// to) the sdk global Bech32 prefix config being set a particular way in
// whatever binary links this package.
func validateCheqdAddress(addr string) error {
	// bech32.Decode accepts an all-uppercase address and normalizes it to
	// lowercase internally before decoding — but every other check in this
	// file (destination == affected address, duplicate detection,
	// destination == destination_address_confirm) compares the raw strings
	// as given, not the decoded bytes. An address and its all-uppercase
	// equivalent decode to the identical account but would NOT be caught as
	// equal/duplicate by those checks. Rejecting anything but lowercase
	// here closes that gap cheaply, rather than switching every comparison
	// in this file to compare decoded bytes.
	if addr != strings.ToLower(addr) {
		return fmt.Errorf("address must be lowercase (got %q) — this handler's duplicate/equality checks compare raw strings, not decoded bytes, so an uppercase variant of an already-listed address would not be caught as the same address", addr)
	}
	hrp, _, err := bech32.Decode(addr, 1023)
	if err != nil {
		return fmt.Errorf("not a valid bech32 address: %w", err)
	}
	if hrp != bech32HRP {
		return fmt.Errorf("expected %q-prefixed address, got prefix %q", bech32HRP, hrp)
	}
	return nil
}

// knownModuleAccountNames lists the standard cosmos-sdk module accounts
// whose addresses are derived deterministically from their name, not from
// any keypair. Sending recovered funds to one of these by mistake would be
// a different flavor of Juno's "sent to an address nobody controls" bug —
// technically valid, practically unreachable (no one holds a private key
// for a module account; funds would be stuck in protocol internals rather
// than usable custody). This list should be extended if cheqd-node
// registers additional modules with their own module accounts.
var knownModuleAccountNames = []string{
	authtypes.FeeCollectorName,
	stakingtypes.BondedPoolName,
	stakingtypes.NotBondedPoolName,
	distrtypes.ModuleName,
	minttypes.ModuleName,
	govtypes.ModuleName,
}

// rejectModuleAccountAddress errors if addr matches any known module
// account's deterministically-derived address.
//
// Deliberately compares raw address bytes, not bech32 strings:
// authtypes.NewModuleAddress(name).String() formats using the SDK's
// *global* Bech32 prefix config, which defaults to "cosmos" and is never
// set by this standalone package — comparing bech32 strings directly would
// silently never match anything and give false confidence. Decoding addr
// ourselves and comparing raw bytes sidesteps that global-config dependency
// entirely, the same reasoning as validateCheqdAddress's use of the raw
// bech32 decoder.
func rejectModuleAccountAddress(addr string) error {
	_, data, err := bech32.Decode(addr, 1023)
	if err != nil {
		return fmt.Errorf("not a valid bech32 address: %w", err)
	}
	rawBytes, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return fmt.Errorf("decoding bech32 data part: %w", err)
	}

	for _, name := range knownModuleAccountNames {
		moduleAddr := authtypes.NewModuleAddress(name)
		if string(moduleAddr.Bytes()) == string(rawBytes) {
			return fmt.Errorf("address %q is the %q module account — not a real custody destination, refusing", addr, name)
		}
	}
	return nil
}
