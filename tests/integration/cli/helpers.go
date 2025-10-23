package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"maps"
	"sort"
	"strings"
	"sync"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	resourcev2 "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	didv2 "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

// The following structs are overridden from the tendermint codebase.
// They are used to parse the output of the `status` command.
// We need to override them because the tendermint codebase types are overridden
// by the cosmos-sdk codebase types.
// Also, ValidatorInfo.PubKey is replaced with cosmos-sdk crypto.PubKey, hence it needs
// to be parsed accordingly.
type NodeStatus struct {
	SyncInfo SyncInfo `json:"sync_info"`
}

type SyncInfo struct {
	LatestBlockHeight int64 `json:"latest_block_height,string"`
	CatchingUp        bool  `json:"catching_up"`
}

func GetNodeStatus(container string, binary string) (NodeStatus, error) {
	out, err := LocalnetExecExec(container, binary, "status")
	if err != nil {
		return NodeStatus{}, err
	}
	var result NodeStatus
	err = json.Unmarshal([]byte(out), &result)
	if err != nil {
		return NodeStatus{}, err
	}
	return result, nil
}

func GetCurrentBlockHeight(container string, binary string) (int64, error) {
	status, err := GetNodeStatus(container, binary)
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func GetVotingEndHeight(currentHeight int64) (int64, error) {
	return currentHeight + VotingPeriod/ExpectedBlockSeconds + ExtraBlocks, nil
}

func CalculateUpgradeHeight(container string, binary string) (int64, int64, error) {
	currentHeight, err := GetCurrentBlockHeight(container, binary)
	if err != nil {
		return 0, 0, err
	}
	votingEndHeight, err := GetVotingEndHeight(currentHeight)
	if err != nil {
		return 0, 0, err
	}
	return currentHeight + VotingPeriod/ExpectedBlockSeconds + ExtraBlocks*2, votingEndHeight, nil
}

// Added to wait for the upgrade to be applied.
// NOTE: This can be extended to run concurrent waits for multiple containers.
func WaitForChainHeight(container string, binary string, height int64, period int64) error {
	var waited int64
	var waitInterval int64 = 1
	var wg sync.WaitGroup

	for waited < period {
		wg.Add(1)
		go waitHeightCallback(container, binary, height, period, &waited, &waitInterval, &wg)
		wg.Wait()
	}

	if waited == period {
		return fmt.Errorf("timeout waiting for chain height")
	}

	return nil
}

func WaitForCaughtUp(container string, binary string, period int64) error {
	var waited int64
	var waitInterval int64 = 1
	var wg sync.WaitGroup

	for waited < period {
		wg.Add(1)
		go waitCaughtUpCallback(container, binary, period, &waited, &waitInterval, &wg)
		wg.Wait()
	}

	if waited == period {
		return fmt.Errorf("timeout waiting for chain height")
	}

	return nil
}

func waitHeightCallback(container string, binary string, height int64, period int64, waited *int64, waitInterval *int64, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Duration(*waitInterval) * time.Second)
	*waited += *waitInterval

	status, err := GetNodeStatus(container, binary)
	if err != nil {
		panic(err)
	}

	if status.SyncInfo.LatestBlockHeight >= height {
		fmt.Printf("Container %s reached height %d after %d seconds of waiting.\n", container, height, *waited)
		*waited = period + 1
		return
	}

	if *waited == period {
		fmt.Printf("Container %s did not reach height %d after %d seconds of waiting.\n", container, height, *waited)
		return
	}

	fmt.Printf("Container %s is at height %d after %d seconds of waiting, with a max waiting period of %d.\n", container, status.SyncInfo.LatestBlockHeight, *waited, period)
}

func waitCaughtUpCallback(container string, binary string, period int64, waited *int64, waitInterval *int64, wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Duration(*waitInterval) * time.Second)
	*waited += *waitInterval

	status, err := GetNodeStatus(container, binary)
	if err != nil {
		panic(err)
	}

	if !status.SyncInfo.CatchingUp {
		fmt.Printf("Container %s is caught up after %d seconds of waiting.\n", container, *waited)
		*waited = period + 1
		return
	}

	if *waited == period {
		fmt.Printf("Container %s is not caught up after %d seconds of waiting.\n", container, *waited)
		return
	}

	fmt.Printf("Container %s is still catching up after %d seconds of waiting, with a max waiting period of %d.\n", container, *waited, period)
}

func TrimExtraLineOffset(input string, offset int) string {
	return strings.Join(strings.Split(input, "\n")[offset:], "")
}

func MakeCodecWithExtendedRegistry() codec.Codec {
	interfaceRegistry := types.NewInterfaceRegistry()

	// TODO: Remove nolint after cheqd-node release v1.x is successful
	// Register the interfaces from the cosmos-sdk codebase.
	interfaceRegistry.RegisterImplementations(
		(*govtypesv1beta1.Content)(nil),
		//nolint: staticcheck
		&upgradetypes.SoftwareUpgradeProposal{},
		&paramproposal.ParameterChangeProposal{},
	)

	interfaceRegistry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&upgradetypes.MsgSoftwareUpgrade{},
		&govtypesv1.MsgExecLegacyContent{},
		&didv2.MsgBurn{},
		&didv2.MsgMint{},
		&didv2.MsgUpdateParams{},
	)
	interfaceRegistry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&resourcev2.MsgUpdateParams{},
	)

	return codec.NewProtoCodec(interfaceRegistry)
}

// issue with type in proposal messages struct, fix it
func convertProposalJSON(input string) ([]byte, error) {
	var data map[string]any
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal input: %w", err)
	}

	proposal, ok := data["proposal"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("proposal not found or invalid")
	}

	rawMessages, ok := proposal["messages"].([]any)
	if !ok {
		return nil, fmt.Errorf("messages not found or invalid")
	}

	var newMessages []any
	for _, msg := range rawMessages {
		msgMap, ok := msg.(map[string]any)
		if !ok {
			continue
		}

		newMsg := make(map[string]any)
		newMsg["@type"] = msgMap["type"]

		if valueMap, ok := msgMap["value"].(map[string]any); ok {
			maps.Copy(newMsg, valueMap)
		}

		newMessages = append(newMessages, newMsg)
	}
	proposal["messages"] = newMessages

	// Marshal to []byte
	output, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal output: %w", err)
	}

	return output, nil
}

// CalculateVoteHash calculates the hash for an exchange rate vote using the same algorithm
// that the Oracle module uses. This simulates the hash calculation that happens in the price feeder.
func CalculateVoteHash(salt string, exchangeRates string, voter string) string {
	// Parse the exchange rates string (e.g. "CHEQ:1.2,BTC:30000.0")
	rateMap := parseExchangeRates(exchangeRates)

	// Create a sorted slice of denoms for deterministic ordering
	denoms := make([]string, 0, len(rateMap))
	for denom := range rateMap {
		denoms = append(denoms, denom)
	}
	sort.Strings(denoms)

	// Build the vote message as salt:denom1:rate1,denom2:rate2,...
	// This is the format required by the Oracle module
	voteMsg := salt
	for _, denom := range denoms {
		voteMsg += ":" + denom + ":" + rateMap[denom]
	}

	// Calculate SHA256 hash
	hash := sha256.Sum256([]byte(voteMsg))

	// Return hex-encoded hash string
	return hex.EncodeToString(hash[:])[:40] // Oracle module uses first 20 bytes
}

// ParseSalt generates a random salt string for use in voting
// In a real scenario, the price feeder would generate this randomly
func GenerateSalt() string {
	// For testing, we use a deterministic "random" string
	// In a real price feeder, this would be a random hex string
	return "2d1a938d791590770571807dc584e4f4ec5641dd3ec66f0c68d3b2c8a7522a63"
}

// parseExchangeRates converts a comma-separated exchange rate string like "CHEQ:1.2,BTC:30000.0"
// into a map of denom -> rate strings
func parseExchangeRates(exchangeRates string) map[string]string {
	rates := make(map[string]string)
	pairs := strings.Split(exchangeRates, ",")

	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			denom := parts[0]
			rate := parts[1]
			rates[denom] = rate
		}
	}

	return rates
}

// ValidatePrevoteHash verifies that a hash matches what would be generated from the salt and rates
func ValidatePrevoteHash(hash string, salt string, exchangeRates string, voter string) bool {
	calculatedHash := CalculateVoteHash(salt, exchangeRates, voter)
	return hash == calculatedHash
}

// ConstructAggregateVoteMsg creates a properly formatted exchange rate string for voting
// This ensures rates are consistently ordered by denom
func ConstructAggregateVoteMsg(rates map[string]string) string {
	// Create a sorted slice of denoms for deterministic ordering
	denoms := make([]string, 0, len(rates))
	for denom := range rates {
		denoms = append(denoms, denom)
	}
	sort.Strings(denoms)

	// Build the formatted string
	pairs := make([]string, 0, len(rates))
	for _, denom := range denoms {
		pairs = append(pairs, denom+":"+rates[denom])
	}

	return strings.Join(pairs, ",")
}
