package cli

import (
	"fmt"
	"sync"
	"time"

	integrationhelpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
	tmcoretypes "github.com/tendermint/tendermint/rpc/core/types"
)

func GetNodeStatus(container string, binary string) (tmcoretypes.ResultStatus, error) {
	out, err := LocalnetExecExec(container, binary, "status")
	if err != nil {
		return tmcoretypes.ResultStatus{}, err
	}
	var result tmcoretypes.ResultStatus
	err = integrationhelpers.Codec.UnmarshalInterfaceJSON([]byte(out), &result)
	if err != nil {
		return tmcoretypes.ResultStatus{}, err
	}
	return result, nil
}

func GetCurrentBlockHeight(container string, binary string) error {
	status, err := GetNodeStatus(container, binary)
	if err != nil {
		return err
	}
	CURRENT_HEIGHT = status.SyncInfo.LatestBlockHeight
	return nil
}

func GetVotingEndHeight() error {
	VOTING_END_HEIGHT = CURRENT_HEIGHT + VOTING_PERIOD/EXPECTED_BLOCK_SECONDS + EXTRA_BLOCKS
	return nil
}

func CalculateUpgradeHeight(container string, binary string) error {
	err := GetCurrentBlockHeight(container, binary)
	if err != nil {
		return err
	}
	err = GetVotingEndHeight()
	if err != nil {
		return err
	}
	UPGRADE_HEIGHT = CURRENT_HEIGHT + VOTING_PERIOD/EXPECTED_BLOCK_SECONDS + EXTRA_BLOCKS*2
	return nil
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
		fmt.Printf("Container %s reached height %d after %d seconds of waiting", container, height, *waited)
		*waited = period + 1
		return
	}

	if *waited == period {
		fmt.Printf("Container %s did not reach height %d after %d seconds of waiting", container, height, *waited)
		return
	}

	fmt.Printf("Container %s is at height %d after %d seconds of waiting, with a max waiting period of %d", container, status.SyncInfo.LatestBlockHeight, *waited, period)
}
