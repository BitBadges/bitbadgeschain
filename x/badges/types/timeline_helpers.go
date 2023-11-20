package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetCurrentManager(ctx sdk.Context, collection *BadgeCollection) string {
	blockTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	managerTimeline := collection.ManagerTimeline
	for _, managerTimelineVal := range managerTimeline {
		found, err := SearchUintRangesForUint(blockTime, managerTimelineVal.TimelineTimes)
		if found || err != nil {
			return managerTimelineVal.Manager
		}
	}

	return ""
}

func GetIsArchived(ctx sdk.Context, collection *BadgeCollection) bool {
	blockTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	isArchivedTimeline := collection.IsArchivedTimeline
	for _, isArchivedTimelineVal := range isArchivedTimeline {
		found, err := SearchUintRangesForUint(blockTime, isArchivedTimelineVal.TimelineTimes)
		if found || err != nil {
			return isArchivedTimelineVal.IsArchived
		}
	}

	return false
}

func GetIsArchivedTimesAndValues(isArchivedTimeline []*IsArchivedTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range isArchivedTimeline {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.IsArchived)
	}
	return times, values
}

func GetCollectionApprovalTimesAndValues(approvals []*CollectionApprovalTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range approvals {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.CollectionApprovals)
	}
	return times, values
}

func GetUserOutgoingApprovalTimesAndValues(approvals []*UserOutgoingApprovalTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range approvals {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.OutgoingApprovals)
	}
	return times, values
}

func GetUserIncomingApprovalTimesAndValues(approvals []*UserIncomingApprovalTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range approvals {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.IncomingApprovals)
	}
	return times, values
}

func GetInheritedBalancesTimesAndValues(inheritedBalances []*InheritedBalancesTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range inheritedBalances {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.InheritedBalances)
	}

	return times, values
}

func GetOffChainBalancesMetadataTimesAndValues(inheritedBalancesMetadata []*OffChainBalancesMetadataTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range inheritedBalancesMetadata {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.OffChainBalancesMetadata)
	}
	return times, values
}

func GetCollectionMetadataTimesAndValues(timeline []*CollectionMetadataTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range timeline {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.CollectionMetadata)
	}
	return times, values
}

func GetBadgeMetadataTimesAndValues(timeline []*BadgeMetadataTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range timeline {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.BadgeMetadata)
	}
	return times, values
}

func GetManagerTimesAndValues(managerTimeline []*ManagerTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range managerTimeline {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.Manager)
	}
	return times, values
}

func GetCustomDataTimesAndValues(customDataTimeline []*CustomDataTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range customDataTimeline {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.CustomData)
	}
	return times, values
}

func GetStandardsTimesAndValues(standardsTimeline []*StandardsTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range standardsTimeline {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.Standards)
	}
	return times, values
}
