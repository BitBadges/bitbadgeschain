package types

import (
	tokentypes "github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
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
	if !currMap.InheritManagerTimelineFrom.IsNil() && !currMap.InheritManagerTimelineFrom.IsZero() {
		if collection == nil {
			panic("Token collection must be provided if map is inheriting manager timeline from a collection")
		}

		return tokentypes.GetCurrentManager(ctx, collection)
	} else {
		blockTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currManager := ""
		for _, managerTimelineVal := range currMap.ManagerTimeline {
			found, err := tokentypes.SearchUintRangesForUint(blockTime, CastUintRanges(managerTimelineVal.TimelineTimes))
			if found || err != nil {
				currManager = managerTimelineVal.Manager
				break
			}
		}
		return currManager
	}
}

func CastActionPermission(perm *ActionPermission) *tokentypes.ActionPermission {
	return &tokentypes.ActionPermission{
		PermanentlyPermittedTimes: CastUintRanges(perm.PermanentlyPermittedTimes),
		PermanentlyForbiddenTimes: CastUintRanges(perm.PermanentlyForbiddenTimes),
	}
}

func CastTimedUpdatePermission(perm *TimedUpdatePermission) *tokentypes.TimedUpdatePermission {
	return &tokentypes.TimedUpdatePermission{
		PermanentlyPermittedTimes: CastUintRanges(perm.PermanentlyPermittedTimes),
		PermanentlyForbiddenTimes: CastUintRanges(perm.PermanentlyForbiddenTimes),
		TimelineTimes:             CastUintRanges(perm.TimelineTimes),
	}
}

func CastActionPermissions(perms []*ActionPermission) []*tokentypes.ActionPermission {
	casted := make([]*tokentypes.ActionPermission, len(perms))
	for i, perm := range perms {
		casted[i] = CastActionPermission(perm)
	}
	return casted
}

func CastTimedUpdatePermissions(perms []*TimedUpdatePermission) []*tokentypes.TimedUpdatePermission {
	casted := make([]*tokentypes.TimedUpdatePermission, len(perms))
	for i, perm := range perms {
		casted[i] = CastTimedUpdatePermission(perm)
	}
	return casted
}

func CastManagerTimeline(timeline *ManagerTimeline) *tokentypes.ManagerTimeline {
	return &tokentypes.ManagerTimeline{
		Manager:       timeline.Manager,
		TimelineTimes: CastUintRanges(timeline.TimelineTimes),
	}
}

func CastManagerTimelineArray(timelines []*ManagerTimeline) []*tokentypes.ManagerTimeline {
	casted := make([]*tokentypes.ManagerTimeline, len(timelines))
	for i, timeline := range timelines {
		casted[i] = CastManagerTimeline(timeline)
	}
	return casted
}

func CastMetadataTimeline(timeline *MapMetadataTimeline) *tokentypes.CollectionMetadataTimeline {
	return &tokentypes.CollectionMetadataTimeline{
		CollectionMetadata: &tokentypes.CollectionMetadata{
			Uri:        timeline.Metadata.Uri,
			CustomData: timeline.Metadata.CustomData,
		},
		TimelineTimes: CastUintRanges(timeline.TimelineTimes),
	}
}

func CastMetadataTimelineArray(timelines []*MapMetadataTimeline) []*tokentypes.CollectionMetadataTimeline {
	casted := make([]*tokentypes.CollectionMetadataTimeline, len(timelines))
	for i, timeline := range timelines {
		casted[i] = CastMetadataTimeline(timeline)
	}
	return casted
}
