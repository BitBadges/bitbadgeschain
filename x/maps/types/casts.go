package types

import (
	badgetypes "github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastUintRanges(ranges []*UintRange) []*badgetypes.UintRange {
	castedRanges := make([]*badgetypes.UintRange, len(ranges))
	for i, rangeVal := range ranges {
		castedRanges[i] = &badgetypes.UintRange{
			Start: rangeVal.Start,
			End:   rangeVal.End,
		}
	}
	return castedRanges
}

func GetCurrentManagerForMap(ctx sdk.Context, currMap *Map, collection *badgetypes.TokenCollection) string {
	if !currMap.InheritManagerTimelineFrom.IsNil() && !currMap.InheritManagerTimelineFrom.IsZero() {
		if collection == nil {
			panic("Token collection must be provided if map is inheriting manager timeline from a collection")
		}

		return badgetypes.GetCurrentManager(ctx, collection)
	} else {
		blockTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currManager := ""
		for _, managerTimelineVal := range currMap.ManagerTimeline {
			found, err := badgetypes.SearchUintRangesForUint(blockTime, CastUintRanges(managerTimelineVal.TimelineTimes))
			if found || err != nil {
				currManager = managerTimelineVal.Manager
				break
			}
		}
		return currManager
	}
}

func CastActionPermission(perm *ActionPermission) *badgetypes.ActionPermission {
	return &badgetypes.ActionPermission{
		PermanentlyPermittedTimes: CastUintRanges(perm.PermanentlyPermittedTimes),
		PermanentlyForbiddenTimes: CastUintRanges(perm.PermanentlyForbiddenTimes),
	}
}

func CastTimedUpdatePermission(perm *TimedUpdatePermission) *badgetypes.TimedUpdatePermission {
	return &badgetypes.TimedUpdatePermission{
		PermanentlyPermittedTimes: CastUintRanges(perm.PermanentlyPermittedTimes),
		PermanentlyForbiddenTimes: CastUintRanges(perm.PermanentlyForbiddenTimes),
		TimelineTimes:             CastUintRanges(perm.TimelineTimes),
	}
}

func CastActionPermissions(perms []*ActionPermission) []*badgetypes.ActionPermission {
	casted := make([]*badgetypes.ActionPermission, len(perms))
	for i, perm := range perms {
		casted[i] = CastActionPermission(perm)
	}
	return casted
}

func CastTimedUpdatePermissions(perms []*TimedUpdatePermission) []*badgetypes.TimedUpdatePermission {
	casted := make([]*badgetypes.TimedUpdatePermission, len(perms))
	for i, perm := range perms {
		casted[i] = CastTimedUpdatePermission(perm)
	}
	return casted
}

func CastManagerTimeline(timeline *ManagerTimeline) *badgetypes.ManagerTimeline {
	return &badgetypes.ManagerTimeline{
		Manager:       timeline.Manager,
		TimelineTimes: CastUintRanges(timeline.TimelineTimes),
	}
}

func CastManagerTimelineArray(timelines []*ManagerTimeline) []*badgetypes.ManagerTimeline {
	casted := make([]*badgetypes.ManagerTimeline, len(timelines))
	for i, timeline := range timelines {
		casted[i] = CastManagerTimeline(timeline)
	}
	return casted
}

func CastMetadataTimeline(timeline *MapMetadataTimeline) *badgetypes.CollectionMetadataTimeline {
	return &badgetypes.CollectionMetadataTimeline{
		CollectionMetadata: &badgetypes.CollectionMetadata{
			Uri:        timeline.Metadata.Uri,
			CustomData: timeline.Metadata.CustomData,
		},
		TimelineTimes: CastUintRanges(timeline.TimelineTimes),
	}
}

func CastMetadataTimelineArray(timelines []*MapMetadataTimeline) []*badgetypes.CollectionMetadataTimeline {
	casted := make([]*badgetypes.CollectionMetadataTimeline, len(timelines))
	for i, timeline := range timelines {
		casted[i] = CastMetadataTimeline(timeline)
	}
	return casted
}
