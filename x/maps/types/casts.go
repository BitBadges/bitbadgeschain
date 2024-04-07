package types

import (
	"math"

	sdkmath "cosmossdk.io/math"
	badgetypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
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

func GetCurrentManagerForMap(ctx sdk.Context, currMap *Map, collection *badgetypes.BadgeCollection) string {
	if !currMap.InheritManagerTimelineFrom.IsNil() && !currMap.InheritManagerTimelineFrom.IsZero() {
		if collection == nil {
			panic("Badge collection must be provided if map is inheriting manager timeline from a collection")
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

func CastIsEditablePermission(perm *IsEditablePermission) *badgetypes.CollectionApprovalPermission {
	return &badgetypes.CollectionApprovalPermission{
		ApprovalId:        perm.KeyListId,
		FromListId:        "All",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes: []*badgetypes.UintRange{
			{Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64)},
		},
		BadgeIds: []*badgetypes.UintRange{
			{Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64)},
		},
		OwnershipTimes: []*badgetypes.UintRange{
			{Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64)},
		},
		PermanentlyPermittedTimes: CastUintRanges(perm.PermanentlyPermittedTimes),
		PermanentlyForbiddenTimes: CastUintRanges(perm.PermanentlyForbiddenTimes),
	}
}

func CastIsEditablePermissions(perms []*IsEditablePermission) []*badgetypes.CollectionApprovalPermission {
	casted := make([]*badgetypes.CollectionApprovalPermission, len(perms))
	for i, perm := range perms {
		casted[i] = CastIsEditablePermission(perm)
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
