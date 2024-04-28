package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func ValidateTimelineTimesDoNotOverlap(times [][]*UintRange) error {
	handledBadgeIds := []*UintRange{}
	for _, time := range times {
		if len(time) == 0 {
			return ErrNoTimelineTimeSpecified
		}

		err := AssertRangesDoNotOverlapAtAll(time, handledBadgeIds)
		if err != nil {
			return err
		}

		handledBadgeIds = append(handledBadgeIds, time...)
	}
	return nil
}

func ValidateApprovalTimeline(ctx sdk.Context, timeline []*CollectionApprovalTimeline, canChangeValues bool) error {
	err := *new(error)
	for _, timelineVal := range timeline {
		err = ValidateCollectionApprovals(ctx, timelineVal.CollectionApprovals, canChangeValues)
		if err != nil {
			return err
		}
	}

	times, _ := GetCollectionApprovalTimesAndValues(timeline)
	err = ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateOffChainBalancesMetadataTimeline(timeline []*OffChainBalancesMetadataTimeline) error {
	for _, timelineVal := range timeline {
		err := ValidateURI(timelineVal.OffChainBalancesMetadata.Uri)
		if err != nil {
			return err
		}
	}

	times, _ := GetOffChainBalancesMetadataTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateBadgeMetadataTimeline(timeline []*BadgeMetadataTimeline, canChangeValues bool) error {
	for _, timelineVal := range timeline {
		err := ValidateBadgeMetadata(timelineVal.BadgeMetadata, canChangeValues)
		if err != nil {
			return err
		}
	}

	times, _ := GetBadgeMetadataTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateCollectionMetadataTimeline(timeline []*CollectionMetadataTimeline) error {
	for _, timelineVal := range timeline {
		err := ValidateURI(timelineVal.CollectionMetadata.Uri)
		if err != nil {
			return err
		}
	}

	times, _ := GetCollectionMetadataTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateStandardsTimeline(timeline []*StandardsTimeline) error {
	times, _ := GetStandardsTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateCustomDataTimeline(timeline []*CustomDataTimeline) error {
	times, _ := GetCustomDataTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateManagerTimeline(timeline []*ManagerTimeline) error {
	times, _ := GetManagerTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateIsArchivedTimeline(timeline []*IsArchivedTimeline) error {
	times, _ := GetIsArchivedTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserOutgoingApprovalTimeline(ctx sdk.Context, timeline []*UserOutgoingApprovalTimeline, address string, canChangeValues bool) error {
	for _, timelineVal := range timeline {
		err := ValidateUserOutgoingApprovals(ctx, timelineVal.OutgoingApprovals, address, canChangeValues)
		if err != nil {
			return err
		}
	}

	times, _ := GetUserOutgoingApprovalTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserIncomingApprovalTimeline(ctx sdk.Context, timeline []*UserIncomingApprovalTimeline, address string, canChangeValues bool) error {
	for _, timelineVal := range timeline {
		err := ValidateUserIncomingApprovals(ctx, timelineVal.IncomingApprovals, address, canChangeValues)
		if err != nil {
			return err
		}
	}

	times, _ := GetUserIncomingApprovalTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}
