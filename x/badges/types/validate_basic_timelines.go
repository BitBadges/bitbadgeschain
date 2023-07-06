package types

func ValidateTimelineTimesDoNotOverlap(times [][]*IdRange) error {
	handledBadgeIds := []*IdRange{}
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

func ValidateApprovedTransferTimeline(timeline []*CollectionApprovedTransferTimeline) error {
	err := *new(error)
	for _, timelineVal := range timeline {
		for _, approvedTransfer := range timelineVal.ApprovedTransfers {
			err = ValidateCollectionApprovedTransfer(approvedTransfer)
			if err != nil {
				return err
			}
		}
	}

	times, _ := GetCollectionApprovedTransferTimesAndValues(timeline)
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

func ValidateBadgeMetadataTimeline(timeline []*BadgeMetadataTimeline) error {
	for _, timelineVal := range timeline {
		err := ValidateBadgeMetadata(timelineVal.BadgeMetadata)
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

func ValidateInheritedBalancesTimeline(timeline []*InheritedBalancesTimeline) error {
	for _, timelineVal := range timeline {
		err := ValidateInheritedBalances(timelineVal.InheritedBalances)
		if err != nil {
			return err
		}
	}

	times, _ := GetInheritedBalancesTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateStandardsTimeline(timeline []*StandardTimeline) error {
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


func ValidateUserApprovedOutgoingTransferTimeline(timeline []*UserApprovedOutgoingTransferTimeline, address string) error {
	for _, timelineVal := range timeline {
		for _, approvedTransfer := range timelineVal.ApprovedOutgoingTransfers {
			err := ValidateUserApprovedOutgoingTransfer(approvedTransfer, address)
			if err != nil {
				return err
			}
		}
	}

	times, _ := GetUserApprovedOutgoingTransferTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserApprovedIncomingTransferTimeline(timeline []*UserApprovedIncomingTransferTimeline, address string) error {
	for _, timelineVal := range timeline {
		for _, approvedTransfer := range timelineVal.ApprovedIncomingTransfers {
			err := ValidateUserApprovedIncomingTransfer(approvedTransfer, address)
			if err != nil {
				return err
			}
		}
	}

	times, _ := GetUserApprovedIncomingTransferTimesAndValues(timeline)
	err := ValidateTimelineTimesDoNotOverlap(times)
	if err != nil {
		return err
	}

	return nil
}