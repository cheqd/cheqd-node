package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	didkeeper "github.com/cheqd/cheqd-node/x/did/keeper"
	resourcekeeper "github.com/cheqd/cheqd-node/x/resource/keeper"
)

func NormalizeDids(didKeeper didkeeper.Keeper, resourceKeeper resourcekeeper.Keeper) error {
	earliestCreatedDidDocs := make(map[string]uint64)



	iterator := didKeeper.DidDocIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		didDoc, err := didKeeper.GetDidDoc(ctx, string(iterator.Key()))
		if err != nil {
			return err
		}

		didDoc.Normalize()

		didKeeper.SetDidDoc(ctx, didDoc)
	}

	return nil
}

}
