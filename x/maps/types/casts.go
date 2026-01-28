package types

import (
	tokentypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastUintRanges(ranges []*UintRange) []*tokentypes.UintRange {
	castedRanges := make([]*tokentypes.UintRange, len(ranges))
	for i, rangeVal := range ranges {
		castedRanges[i] = &tokentypes.UintRange{
			Start: rangeVal.Start,
			End:   rangeVal.End,
		}
	}
	return castedRanges
}

func GetCurrentManagerForMap(ctx sdk.Context, currMap *Map, collection *tokentypes.TokenCollection) string {
	if !currMap.InheritManagerFrom.IsNil() && !currMap.InheritManagerFrom.IsZero() {
		if collection == nil {
			panic("Token collection must be provided if map is inheriting manager from a collection")
		}

		return tokentypes.GetCurrentManager(ctx, collection)
	} else {
		return currMap.Manager
	}
}

func CastActionPermission(perm *ActionPermission) *tokentypes.ActionPermission {
	return &tokentypes.ActionPermission{
		PermanentlyPermittedTimes: CastUintRanges(perm.PermanentlyPermittedTimes),
		PermanentlyForbiddenTimes: CastUintRanges(perm.PermanentlyForbiddenTimes),
	}
}

func CastActionPermissions(perms []*ActionPermission) []*tokentypes.ActionPermission {
	casted := make([]*tokentypes.ActionPermission, len(perms))
	for i, perm := range perms {
		casted[i] = CastActionPermission(perm)
	}
	return casted
}

// Helper functions to convert maps metadata to badges CollectionMetadata
func CastMapMetadataToCollectionMetadata(metadata *Metadata) *tokentypes.CollectionMetadata {
	if metadata == nil {
		return nil
	}
	return &tokentypes.CollectionMetadata{
		Uri:        metadata.Uri,
		CustomData: metadata.CustomData,
	}
}
