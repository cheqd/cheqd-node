package scenarios

import (
	migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/unit/setup"
)

type IBuilder interface {
	buildExistingDids() error
	buildExistingResources() error
	buildExpectedDids() error
	buildExpectedResources() error
	BuildDataSet(setup migrationsetup.TestSetup) (IDataSet, error)
}
