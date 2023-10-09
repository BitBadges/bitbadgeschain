package types

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

func ValidateApprovalTimeline(timeline []*CollectionApprovalTimeline, canChangeValues bool) error {
	err := *new(error)
	for _, timelineVal := range timeline {
		err = ValidateCollectionApprovals(timelineVal.CollectionApprovals, canChangeValues)
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

// func ValidateInheritedBalancesTimeline(timeline []*InheritedBalancesTimeline) error {
// 	for _, timelineVal := range timeline {
// 		err := ValidateInheritedBalances(timelineVal.InheritedBalances)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	times, _ := GetInheritedBalancesTimesAndValues(timeline)
// 	err := ValidateTimelineTimesDoNotOverlap(times)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

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

func ValidateContractAddressTimeline(timeline []*ContractAddressTimeline) error {
	times, _ := GetContractAddressTimesAndValues(timeline)
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

func ValidateUserOutgoingApprovalTimeline(timeline []*UserOutgoingApprovalTimeline, address string, canChangeValues bool) error {
	for _, timelineVal := range timeline {
		err := ValidateUserOutgoingApprovals(timelineVal.OutgoingApprovals, address, canChangeValues)
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

func ValidateUserIncomingApprovalTimeline(timeline []*UserIncomingApprovalTimeline, address string, canChangeValues bool) error {
	for _, timelineVal := range timeline {
		err := ValidateUserIncomingApprovals(timelineVal.IncomingApprovals, address, canChangeValues)
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
