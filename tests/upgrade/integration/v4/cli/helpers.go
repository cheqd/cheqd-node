package cli

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	didv2 "github.com/cheqd/cheqd-node/x/did/types"
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	solomachine "github.com/cosmos/ibc-go/v8/modules/light-clients/06-solomachine"
)

// The following structs are overridden from the tendermint codebase.
// They are used to parse the output of the `status` command.
// We need to override them because the tendermint codebase types are overridden
// by the cosmos-sdk codebase types.
// Also, ValidatorInfo.PubKey is replaced with cosmos-sdk crypto.PubKey, hence it needs
// to be parsed accordingly.
type NodeStatus struct {
	NodeInfo      DefaultNodeInfo `json:"node_info"`
	SyncInfo      SyncInfo        `json:"sync_info"`
	ValidatorInfo ValidatorInfo   `json:"validator_info"`
}

type DefaultNodeInfo struct {
	ProtocolVersion ProtocolVersion      `json:"protocol_version"`
	ID              string               `json:"id"`
	ListenAddr      string               `json:"listen_addr"`
	Network         string               `json:"network"`
	Version         string               `json:"version"`
	Channels        tmbytes.HexBytes     `json:"channels"`
	Moniker         string               `json:"moniker"`
	Other           DefaultNodeInfoOther `json:"other"`
}

type ProtocolVersion struct {
	P2P   uint64 `json:"p2p,string"`
	Block uint64 `json:"block,string"`
	App   uint64 `json:"app,string"`
}

type DefaultNodeInfoOther struct {
	TxIndex    string `json:"tx_index"`
	RPCAddress string `json:"rpc_address"`
}

type SyncInfo struct {
	LatestBlockHash   tmbytes.HexBytes `json:"latest_block_hash"`
	LatestAppHash     tmbytes.HexBytes `json:"latest_app_hash"`
	LatestBlockHeight int64            `json:"latest_block_height,string"`
	LatestBlockTime   time.Time        `json:"latest_block_time"`

	EarliestBlockHash   tmbytes.HexBytes `json:"earliest_block_hash"`
	EarliestAppHash     tmbytes.HexBytes `json:"earliest_app_hash"`
	EarliestBlockHeight int64            `json:"earliest_block_height,string"`
	EarliestBlockTime   time.Time        `json:"earliest_block_time"`

	CatchingUp bool `json:"catching_up"`
}

type ValidatorInfo struct {
	Address     tmbytes.HexBytes `json:"Address"`
	PubKey      interface{}      `json:"PubKey"`
	VotingPower int64            `json:"VotingPower,string"`
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
	return currentHeight + VotingPeriod/ExpectedBlockSeconds + ExtraBlocks*2, votingEndHeight + 100, nil
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
	)

	solomachine.RegisterInterfaces(interfaceRegistry)

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

func RunExportCommand(container, binary string) (out string, err error) {
	// get id of container
	containerId, err := GetContainerIDByName(container)
	if err != nil {
		return
	}

	// fetch latest height which will be used for export
	height, err := GetCurrentBlockHeight(container, binary)
	if err != nil {
		return
	}

	// stop container to make export command work
	_, err = LocalnetStopContainerWithId(containerId)
	if err != nil {
		return
	}

	// run export command
	out, err = LocalnetExecRunWithVolume(containerId, []string{binary, "export", fmt.Sprintf("--height=%d", height)})
	if err != nil {
		return
	}

	return
}
