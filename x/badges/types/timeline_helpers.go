package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetCurrentManager(ctx sdk.Context, collection *TokenCollection) string {
	return collection.Manager
}

func GetIsArchived(ctx sdk.Context, collection *TokenCollection) bool {
	return collection.IsArchived
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

func GetCollectionMetadataTimesAndValues(timeline []*CollectionMetadataTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range timeline {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.CollectionMetadata)
	}
	return times, values
}

func GetTokenMetadataTimesAndValues(timeline []*TokenMetadataTimeline) ([][]*UintRange, []interface{}) {
	times := [][]*UintRange{}
	values := []interface{}{}
	for _, timelineVal := range timeline {
		times = append(times, timelineVal.TimelineTimes)
		values = append(values, timelineVal.TokenMetadata)
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
