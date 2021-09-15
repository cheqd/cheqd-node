package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		CredDefList: []*CredDef{},
		SchemaList:  []*Schema{},
		AttribList:  []*Attrib{},
		DidList:     []*Did{},
		NymList:     []*Nym{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated ID in credDef
	credDefIdMap := make(map[uint64]bool)

	for _, elem := range gs.CredDefList {
		if _, ok := credDefIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for credDef")
		}
		credDefIdMap[elem.Id] = true
	}
	// Check for duplicated ID in schema
	schemaIdMap := make(map[uint64]bool)

	for _, elem := range gs.SchemaList {
		if _, ok := schemaIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for schema")
		}
		schemaIdMap[elem.Id] = true
	}
	// Check for duplicated ID in attrib
	attribIdMap := make(map[uint64]bool)

	for _, elem := range gs.AttribList {
		if _, ok := attribIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for attrib")
		}
		attribIdMap[elem.Id] = true
	}
	// Check for duplicated ID in did
	didIdMap := make(map[uint64]bool)

	for _, elem := range gs.DidList {
		if _, ok := didIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for did")
		}
		didIdMap[elem.Id] = true
	}
	// Check for duplicated ID in nym
	nymIdMap := make(map[uint64]bool)

	for _, elem := range gs.NymList {
		if _, ok := nymIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for nym")
		}
		nymIdMap[elem.Id] = true
	}

	return nil
}
