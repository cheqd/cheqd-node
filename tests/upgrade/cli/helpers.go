package cli

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmcoretypes "github.com/tendermint/tendermint/rpc/core/types"
)

// The following structs are overridden from the tendermint codebase.
// They are used to parse the output of the `status` command.
// We need to override them because the tendermint codebase types are overridden
// by the cosmos-sdk codebase types.
type NodeStatus struct {
	NodeInfo      DefaultNodeInfo           `json:"NodeInfo"`
	SyncInfo      tmcoretypes.SyncInfo      `json:"SyncInfo"`
	ValidatorInfo tmcoretypes.ValidatorInfo `json:"ValidatorInfo"`
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

func GetNodeStatus(container string, binary string) (NodeStatus, error) {
	out, err := LocalnetExecExec(container, binary, "status")
	if err != nil {
		return NodeStatus{}, err
	}
	fmt.Println("out", out)
	var result NodeStatus
	err = json.Unmarshal([]byte(out), &result)
	if err != nil {
		return NodeStatus{}, err
	}
	return result, nil
}

func GetCurrentBlockHeight(container string, binary string) (int64, error) {
	status, err := GetNodeStatus(container, binary)
	fmt.Println("status", status)
	fmt.Println("current height", status.SyncInfo.LatestBlockHeight)
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func GetVotingEndHeight(currentHeight int64) (int64, error) {
	return currentHeight + VOTING_PERIOD/EXPECTED_BLOCK_SECONDS + EXTRA_BLOCKS, nil
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
	return currentHeight + VOTING_PERIOD/EXPECTED_BLOCK_SECONDS + EXTRA_BLOCKS*2, votingEndHeight, nil
}

// Added to wait for the upgrade to be applied.
// NOTE: This can be extended to run concurrent waits for multiple containers.
func WaitForChainHeight(container string, binary string, height int64, period int64) error {
	var waited int64
	var waitInterval int64 = 1
	var wg sync.WaitGroup

	for waited < period {
		wg.Add(1)
		go waitCallback(container, binary, height, period, &waited, &waitInterval, &wg)
		wg.Wait()
	}

	if waited == period {
		return fmt.Errorf("timeout waiting for chain height")
	}

	return nil
}

func waitCallback(container string, binary string, height int64, period int64, waited *int64, waitInterval *int64, wg *sync.WaitGroup) {
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
