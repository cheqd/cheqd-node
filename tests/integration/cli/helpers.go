package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	didv2 "github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// The following structs are overridden from the tendermint codebase.
// They are used to parse the output of the `status` command.
// We need to override them because the tendermint codebase types are overridden
// by the cosmos-sdk codebase types.
// Also, ValidatorInfo.PubKey is replaced with cosmos-sdk crypto.PubKey, hence it needs
// to be parsed accordingly.
type NodeStatus struct {
	SyncInfo SyncInfo `json:"SyncInfo"`
}

type SyncInfo struct {
	LatestBlockHeight int64 `json:"latest_block_height,string"`
	CatchingUp        bool  `json:"catching_up"`
}

func GetNodeStatus(container string, binary string) (NodeStatus, error) {
	out, err := LocalnetExecExec(container, binary, "status", "--log_format", OutputFormat)

	fmt.Println("out, err>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>.", out, err)
	if err != nil {
		fmt.Printf("\"returning error while fetching\": %v\n", "returning error while fetching")
		fmt.Printf("err: %v\n", err)
		return NodeStatus{}, err
	}

	extractedJSON, err := extractOnlyJSON(out)
	if err != nil {
		return NodeStatus{}, err
	}
	out = extractedJSON

	var result NodeStatus
	err = json.Unmarshal([]byte(out), &result)
	if err != nil {
		fmt.Printf("\"returning error while unmarshalling\": %v\n", "returning error while unmarshalling")
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
		// panic(err)
		fmt.Printf("err>>>>>>>>>>>>>>>>>>>>>>>>>>>>>...: %v\n", err)
		return
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

	return codec.NewProtoCodec(interfaceRegistry)
}
