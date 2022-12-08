package types

import (
	"fmt"
)

const (
	DefaultDidNamespace       = "testnet"
	DefaultCreateDidTxFee     = 50e9                   // 50 CHEQ or 50000000000 ncheq
	DefaultUpdateDidTxFee     = 25e9                   // 25 CHEQ or 25000000000 ncheq
	DefaultDeactivateDidTxFee = 10e9                   // 10 CHEQ or 10000000000 ncheq
	DefaultBurnFactor         = "0.500000000000000000" // 0.5 or 50%
)

// DefaultGenesis returns the default `did` genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		VersionSets:  []*DidDocVersionSet{},
		DidNamespace: DefaultDidNamespace,
		FeeParams:    DefaultFeeParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	err := gs.ValidateNoDuplicates()
	if err != nil {
		return err
	}

	err = gs.ValidateVersionSets()
	if err != nil {
		return err
	}

	err = gs.FeeParams.ValidateBasic()
	if err != nil {
		return err
	}

	return gs.ValidateBasic()
}

func (gs GenesisState) ValidateNoDuplicates() error {
	// Check for duplicates in version set
	didCache := make(map[string]bool)

	for _, versionSet := range gs.VersionSets {
		did := versionSet.DidDocs[0].DidDoc.Id
		if _, ok := didCache[did]; ok {
			return fmt.Errorf("duplicated didDoc found with id %s", did)
		}

		didCache[did] = true

		// Check for duplicates in didDoc versions
		versionCache := make(map[string]bool)

		for _, didDoc := range versionSet.DidDocs {
			version := didDoc.Metadata.VersionId
			if _, ok := versionCache[version]; ok {
				return fmt.Errorf("duplicated didDoc version found with id %s and version %s", did, version)
			}

			versionCache[version] = true
		}

		// Check that latest version is present
		if _, ok := versionCache[versionSet.LatestVersion]; !ok {
			return fmt.Errorf("latest version not found in didDoc with id %s", did)
		}
	}

	return nil
}

func (gs GenesisState) ValidateVersionSets() error {
	for _, versionSet := range gs.VersionSets {
		did := versionSet.DidDocs[0].DidDoc.Id

		for _, didDoc := range versionSet.DidDocs {
			if did != didDoc.DidDoc.Id {
				return fmt.Errorf("diddoc %s does not belong to version set %s", didDoc.DidDoc.Id, did)
			}
		}
	}

	return nil
}

func (gs GenesisState) ValidateBasic() error {
	for _, versionSet := range gs.VersionSets {
		for _, didDoc := range versionSet.DidDocs {
			err := didDoc.DidDoc.Validate(nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
