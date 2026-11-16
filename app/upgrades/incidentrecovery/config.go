package incidentrecovery

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/btcutil/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"gopkg.in/yaml.v3"
)

// Config is the configurable input to the incident-recovery upgrade handler.
// Loaded from a YAML file (see config/addresses.example.yaml) — nothing here
// is hardcoded in Go source, so the affected-address list and destination
// can be changed without recompiling.
type Config struct {
	UpgradeName        string `yaml:"upgrade_name"`
	DestinationAddress string `yaml:"destination_address"`
	// DestinationAddressConfirm must be typed/pasted independently and must
	// match DestinationAddress exactly, byte for byte. This exists
	// specifically because a single unconfirmed destination field is what
	// let Juno's Prop 21 copy-paste error through undetected. Two people
	// (or the same person, twice, from two independent sources) agreeing on
	// the exact same string is a much stronger check than one field typed
	// once.
	DestinationAddressConfirm string   `yaml:"destination_address_confirm"`
	AffectedAddresses         []string `yaml:"affected_addresses"`
}

const placeholderDestination = "PLACEHOLDER_UNVERIFIED_DO_NOT_USE"

const bech32HRP = "cheqd"

// LoadConfig reads and validates the config file at path. It deliberately
// fails loudly and refuses to return a usable Config if the destination
// address is still the placeholder, is malformed, or if the affected-address
// list is empty — the Juno Prop 21 incident ($36M sent to an inaccessible
// address via a copy-paste error) is exactly the failure mode these checks
// exist to catch before it ever reaches a running handler.
func LoadConfig(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}
	return parseAndValidate(raw, path)
}

// LoadConfigVerified is LoadConfig plus a binding check: the raw config
// file's sha256 must match expectedSHA256Hex exactly (case-insensitively).
//
// This exists because the config file is read from each node's own local
// disk, not from on-chain state. Without this check, a software-upgrade
// governance proposal only ever approves an upgrade *name* — the actual
// addresses and destination are whatever each validator happens to have
// in a local YAML file, which nothing guarantees is identical across
// nodes. Two nodes with even slightly different (but both individually
// "valid") files would silently compute different post-upgrade state and
// diverge or halt unpredictably.
//
// The expected hash should be committed to the governance proposal's own
// Plan.Info field (see upgrade.go), which IS part of consensus — every
// node evaluates the exact same expected hash. A node whose local file
// doesn't match it fails loudly and halts, rather than silently applying
// different data; a node with no file at all fails the same way. This
// doesn't guarantee every node's operator used the right file, but it
// guarantees any node that proceeds did so starting from byte-identical
// configuration, and anyone verifying the proposal can confirm the
// expected hash matches the config their multisig/leadership approved
// before voting yes.
func LoadConfigVerified(path string, expectedSHA256Hex string) (*Config, error) {
	expectedSHA256Hex = strings.ToLower(strings.TrimSpace(expectedSHA256Hex))
	if expectedSHA256Hex == "" {
		return nil, fmt.Errorf("refusing to load config: no expected sha256 was supplied (governance plan.Info was empty) — this config would be used completely unverified against what was actually voted on")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	actual := ConfigSHA256(raw)
	if actual != expectedSHA256Hex {
		return nil, fmt.Errorf("config file %q does not match the governance-approved config: expected sha256 %s, got %s — refusing to run with a config that was not the one voted on", path, expectedSHA256Hex, actual)
	}

	return parseAndValidate(raw, path)
}

// ConfigSHA256 returns the lowercase hex-encoded sha256 digest of raw config
// file bytes, in the same form expected by LoadConfigVerified and recorded
// in the governance proposal's Plan.Info.
func ConfigSHA256(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func parseAndValidate(raw []byte, path string) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config %q: %w", path, err)
	}

	return &cfg, nil
}

// Validate checks the config is safe to use. Called by LoadConfig, and
// should also be called again immediately before the upgrade handler
// actually executes, as a second independent check at the point of use.
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
