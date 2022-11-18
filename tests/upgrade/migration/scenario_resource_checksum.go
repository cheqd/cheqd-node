package migration

import (
	// "fmt"
	// "path/filepath"

	// migrationsetup "github.com/cheqd/cheqd-node/tests/upgrade/migration/setup"

	// didtestssetup "github.com/cheqd/cheqd-node/x/did/tests/setup"
	didtypes "github.com/cheqd/cheqd-node/x/did/types"
	didtypesv1 "github.com/cheqd/cheqd-node/x/did/types/v1"

	resourcetypes "github.com/cheqd/cheqd-node/x/resource/types"
	resourcetypesv1 "github.com/cheqd/cheqd-node/x/resource/types/v1"
)

// var (
// 	didDoc                         didtypesv1.MsgCreateDidPayload
// 	didInfo                        migrationsetup.MinimalDidDocInfoV1
// 	existingChecksumResource       resourcetypesv1.MsgCreateResourcePayload
// 	expectedChecksumResourceHeader resourcetypes.Metadata
// )

func InitDataChunkForResourceChecksum() (DataChunk, error) {
	// err := Loader(filepath.Join("payload", "existing", "v1", "diddoc_uuid.json"), &didDoc)
	// if err != nil {
	// 	fmt.Println("Error loading didDoc")
	// 	return DataChunk{}, err
	// }
	// var signInput SignInput
	// err = Loader(filepath.Join("keys", "signinput_uuid.json"), &signInput)
	// if err != nil {
	// 	fmt.Println("Error loading signInput")
	// 	return DataChunk{}, err
	// }
	// err = Loader(filepath.Join("payload", "existing", "v1", "resource_checksum.json"), &existingChecksumResource)
	// if err != nil {
	// 	fmt.Println("Error loading existingChecksumResource")
	// 	return DataChunk{}, err
	// }
	// err = Loader(filepath.Join("payload", "expected", "v2", "resource_checksum.json"), &expectedChecksumResourceHeader)
	// if err != nil {
	// 	fmt.Println("Error loading expectedChecksumResourceHeader")
	// 	return DataChunk{}, err
	// }

	dataChunk := DataChunk{
		[]resourcetypesv1.Resource{},
		[]didtypesv1.StateValue{},
		[]resourcetypes.Resource{},
		[]didtypes.DidDocWithMetadata{},
	}
	return dataChunk, nil
}
