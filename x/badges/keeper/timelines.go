package keeper

import (
	proto "github.com/gogo/protobuf/proto"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//TODO: DRY and clean this file up; a lot of repeated code; also work on naming conventions

func GetCombosForStringTimelines(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
	x := []*types.UniversalPermissionDetails{}
	if oldValue.(string) != newValue.(string) {
		x = append(x, &types.UniversalPermissionDetails{})
	}
	return x
}

func ValidateCollectionApprovedTransfersUpdate(ctx sdk.Context, oldApprovedTransfers []*types.CollectionApprovedTransferTimeline, newApprovedTransfers []*types.CollectionApprovedTransferTimeline, canUpdateApprovedTransfers []*types.CollectionApprovedTransferPermission) error {
	oldTimes, oldValues := GetCollectionApprovedTransferTimesAndValues(oldApprovedTransfers)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetCollectionApprovedTransferTimesAndValues(newApprovedTransfers)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, []*types.CollectionApprovedTransfer{}, func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldApprovedTransfers := oldValue.([]*types.CollectionApprovedTransfer)
		firstMatchesForOld := types.GetFirstMatchOnly(CastCollectionApprovedTransferToUniversalPermission(oldApprovedTransfers))

		newApprovedTransfers := newValue.([]*types.CollectionApprovedTransfer)
		firstMatchesForNew := types.GetFirstMatchOnly(CastCollectionApprovedTransferToUniversalPermission(newApprovedTransfers))

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			oldVal := oldDetails.ArbitraryValue.([]*types.CollectionApprovedTransfer)
			newVal := newDetails.ArbitraryValue.([]*types.CollectionApprovedTransfer)
		
			different := false
			if len(newVal) != len(oldVal) {
				different = true
			} else {
				for i := 0; i < len(newVal); i++ {
					if newVal[i].RequireToEqualsInitiatedBy != oldVal[i].RequireToEqualsInitiatedBy ||	
						newVal[i].RequireFromEqualsInitiatedBy != oldVal[i].RequireFromEqualsInitiatedBy ||
						newVal[i].RequireToDoesNotEqualInitiatedBy != oldVal[i].RequireToDoesNotEqualInitiatedBy ||
						newVal[i].RequireFromDoesNotEqualInitiatedBy != oldVal[i].RequireFromDoesNotEqualInitiatedBy ||
						newVal[i].OverridesFromApprovedOutgoingTransfers != oldVal[i].OverridesFromApprovedOutgoingTransfers ||
						newVal[i].OverridesToApprovedIncomingTransfers != oldVal[i].OverridesToApprovedIncomingTransfers ||
						newVal[i].TrackerId != oldVal[i].TrackerId ||
						newVal[i].Uri != oldVal[i].Uri ||
						newVal[i].CustomData != oldVal[i].CustomData ||
						newVal[i].MaxNumTransfers != oldVal[i].MaxNumTransfers {
							different = true
							break
					}

					if len(newVal[i].AllowedCombinations) != len(oldVal[i].AllowedCombinations) {
						different = true
						break
					} else {
						for j := 0; j < len(newVal[i].AllowedCombinations); j++ {
							x, err := proto.Marshal(newVal[i].AllowedCombinations[j])
							if err != nil {
								panic(err)
							}

							y, err := proto.Marshal(oldVal[i].AllowedCombinations[j])
							if err != nil {
								panic(err)
							}

							if string(x) != string(y) {
								different = true
								break
							}
						}
					}

					if len(newVal[i].Challenges) != len(oldVal[i].Challenges) {
						different = true
						break
					} else {
						for j := 0; j < len(newVal[i].Challenges); j++ {
							x, err := proto.Marshal(newVal[i].Challenges[j])
							if err != nil {
								panic(err)
							}

							y, err := proto.Marshal(oldVal[i].Challenges[j])
							if err != nil {
								panic(err)
							}

							if string(x) != string(y) {
								different = true
								break
							}
						}
					}

					x, err := proto.Marshal(newVal[i].Approvals)
					if err != nil {
						panic(err)
					}

					y, err := proto.Marshal(oldVal[i].Approvals)
					if err != nil {
						panic(err)
					}

					if string(x) != string(y) {
						different = true
						break
					}

					x, err = proto.Marshal(newVal[i].PerAddressApprovals)
					if err != nil {
						panic(err)
					}

					y, err = proto.Marshal(oldVal[i].PerAddressApprovals)
					if err != nil {
						panic(err)
					}

					if string(x) != string(y) {
						different = true
						break
					}

					x, err = proto.Marshal(newVal[i].PerAddressMaxNumTransfers)
					if err != nil {
						panic(err)
					}

					y, err = proto.Marshal(oldVal[i].PerAddressMaxNumTransfers)
					if err != nil {
						panic(err)
					}

					if string(x) != string(y) {
						different = true
						break
					}


				}
			}
				
			if different {
				detailsToReturn = append(detailsToReturn, overlap)
			}
		}

		

		//If there are combinations in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn
	})

	err := CheckCollectionApprovedTransferPermission(ctx, detailsToCheck, canUpdateApprovedTransfers)
	if err != nil {
		return err
	}

	return nil
}


func ValidateUserApprovedOutgoingTransfersUpdate(ctx sdk.Context, oldApprovedTransfers []*types.UserApprovedOutgoingTransferTimeline, newApprovedTransfers []*types.UserApprovedOutgoingTransferTimeline, canUpdateApprovedTransfers []*types.UserApprovedTransferPermission) error {
	oldTimes, oldValues := GetUserApprovedOutgoingTransferTimesAndValues(oldApprovedTransfers)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetUserApprovedOutgoingTransferTimesAndValues(newApprovedTransfers)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, []*types.UserApprovedOutgoingTransfer{}, func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldApprovedTransfers := oldValue.([]*types.UserApprovedOutgoingTransfer)
		firstMatchesForOld := types.GetFirstMatchOnly(CastUserApprovedOutgoingTransferToUniversalPermission(oldApprovedTransfers))

		newApprovedTransfers := newValue.([]*types.UserApprovedOutgoingTransfer)
		firstMatchesForNew := types.GetFirstMatchOnly(CastUserApprovedOutgoingTransferToUniversalPermission(newApprovedTransfers))

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			oldVal := oldDetails.ArbitraryValue.([]*types.UserApprovedOutgoingTransfer)
			newVal := newDetails.ArbitraryValue.([]*types.UserApprovedOutgoingTransfer)
			
			different := false
			if len(newVal) != len(oldVal) {
				different = true
			} else {
				for i := 0; i < len(newVal); i++ {
					if newVal[i].RequireToEqualsInitiatedBy != oldVal[i].RequireToEqualsInitiatedBy ||	
						newVal[i].RequireToDoesNotEqualInitiatedBy != oldVal[i].RequireToDoesNotEqualInitiatedBy ||
						newVal[i].TrackerId != oldVal[i].TrackerId ||
						newVal[i].Uri != oldVal[i].Uri ||
						newVal[i].CustomData != oldVal[i].CustomData ||
						newVal[i].MaxNumTransfers != oldVal[i].MaxNumTransfers {
							different = true
							break
					}

					if len(newVal[i].AllowedCombinations) != len(oldVal[i].AllowedCombinations) {
						different = true
						break
					} else {
						for j := 0; j < len(newVal[i].AllowedCombinations); j++ {
							x, err := proto.Marshal(newVal[i].AllowedCombinations[j])
							if err != nil {
								panic(err)
							}

							y, err := proto.Marshal(oldVal[i].AllowedCombinations[j])
							if err != nil {
								panic(err)
							}

							if string(x) != string(y) {
								different = true
								break
							}
						}
					}

					if len(newVal[i].Challenges) != len(oldVal[i].Challenges) {
						different = true
						break
					} else {
						for j := 0; j < len(newVal[i].Challenges); j++ {
							x, err := proto.Marshal(newVal[i].Challenges[j])
							if err != nil {
								panic(err)
							}

							y, err := proto.Marshal(oldVal[i].Challenges[j])
							if err != nil {
								panic(err)
							}

							if string(x) != string(y) {
								different = true
								break
							}
						}
					}

					x, err := proto.Marshal(newVal[i].Approvals)
					if err != nil {
						panic(err)
					}

					y, err := proto.Marshal(oldVal[i].Approvals)
					if err != nil {
						panic(err)
					}

					if string(x) != string(y) {
						different = true
						break
					}

					x, err = proto.Marshal(newVal[i].PerAddressApprovals)
					if err != nil {
						panic(err)
					}

					y, err = proto.Marshal(oldVal[i].PerAddressApprovals)
					if err != nil {
						panic(err)
					}

					if string(x) != string(y) {
						different = true
						break
					}

					x, err = proto.Marshal(newVal[i].PerAddressMaxNumTransfers)
					if err != nil {
						panic(err)
					}

					y, err = proto.Marshal(oldVal[i].PerAddressMaxNumTransfers)
					if err != nil {
						panic(err)
					}

					if string(x) != string(y) {
						different = true
						break
					}
				}
			}
				
			if different {
				detailsToReturn = append(detailsToReturn, overlap)
			}
		}

		//If there are combinations in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn
	})

	err := CheckUserApprovedTransferPermission(ctx, detailsToCheck, canUpdateApprovedTransfers)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserApprovedIncomingTransfersUpdate(ctx sdk.Context, oldApprovedTransfers []*types.UserApprovedIncomingTransferTimeline, newApprovedTransfers []*types.UserApprovedIncomingTransferTimeline, canUpdateApprovedTransfers []*types.UserApprovedTransferPermission) error {
	oldTimes, oldValues := GetUserApprovedIncomingTransferTimesAndValues(oldApprovedTransfers)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetUserApprovedIncomingTransferTimesAndValues(newApprovedTransfers)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, []*types.UserApprovedOutgoingTransfer{}, func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldApprovedTransfers := oldValue.([]*types.UserApprovedIncomingTransfer)
		firstMatchesForOld := types.GetFirstMatchOnly(CastUserApprovedIncomingTransferToUniversalPermission(oldApprovedTransfers))

		newApprovedTransfers := newValue.([]*types.UserApprovedIncomingTransfer)
		firstMatchesForNew := types.GetFirstMatchOnly(CastUserApprovedIncomingTransferToUniversalPermission(newApprovedTransfers))

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			oldVal := oldDetails.ArbitraryValue.([]*types.UserApprovedIncomingTransfer)
			newVal := newDetails.ArbitraryValue.([]*types.UserApprovedIncomingTransfer)
			
			different := false
			if len(newVal) != len(oldVal) {
				different = true
			} else {
				for i := 0; i < len(newVal); i++ {
					if newVal[i].RequireFromDoesNotEqualInitiatedBy != oldVal[i].RequireFromDoesNotEqualInitiatedBy ||	
						newVal[i].RequireFromEqualsInitiatedBy != oldVal[i].RequireFromEqualsInitiatedBy ||
						newVal[i].TrackerId != oldVal[i].TrackerId ||
						newVal[i].Uri != oldVal[i].Uri ||
						newVal[i].CustomData != oldVal[i].CustomData ||
						newVal[i].MaxNumTransfers != oldVal[i].MaxNumTransfers {
							different = true
							break
					}

					if len(newVal[i].AllowedCombinations) != len(oldVal[i].AllowedCombinations) {
						different = true
						break
					} else {
						for j := 0; j < len(newVal[i].AllowedCombinations); j++ {
							x, err := proto.Marshal(newVal[i].AllowedCombinations[j])
							if err != nil {
								panic(err)
							}

							y, err := proto.Marshal(oldVal[i].AllowedCombinations[j])
							if err != nil {
								panic(err)
							}

							if string(x) != string(y) {
								different = true
								break
							}
						}
					}

					if len(newVal[i].Challenges) != len(oldVal[i].Challenges) {
						different = true
						break
					} else {
						for j := 0; j < len(newVal[i].Challenges); j++ {
							x, err := proto.Marshal(newVal[i].Challenges[j])
							if err != nil {
								panic(err)
							}

							y, err := proto.Marshal(oldVal[i].Challenges[j])
							if err != nil {
								panic(err)
							}

							if string(x) != string(y) {
								different = true
								break
							}
						}
					}

					x, err := proto.Marshal(newVal[i].Approvals)
					if err != nil {
						panic(err)
					}

					y, err := proto.Marshal(oldVal[i].Approvals)
					if err != nil {
						panic(err)
					}

					if string(x) != string(y) {
						different = true
						break
					}

					x, err = proto.Marshal(newVal[i].PerAddressApprovals)
					if err != nil {
						panic(err)
					}

					y, err = proto.Marshal(oldVal[i].PerAddressApprovals)
					if err != nil {
						panic(err)
					}

					if string(x) != string(y) {
						different = true
						break
					}

					x, err = proto.Marshal(newVal[i].PerAddressMaxNumTransfers)
					if err != nil {
						panic(err)
					}

					y, err = proto.Marshal(oldVal[i].PerAddressMaxNumTransfers)
					if err != nil {
						panic(err)
					}

					if string(x) != string(y) {
						different = true
						break
					}
				}
			}
				
			if different {
				detailsToReturn = append(detailsToReturn, overlap)
			}
		}

		//If there are combinations in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn
	})

	err := CheckUserApprovedTransferPermission(ctx, detailsToCheck, canUpdateApprovedTransfers)
	if err != nil {
		return err
	}

	return nil
}

func ValidateBadgeMetadataUpdate(ctx sdk.Context, oldBadgeMetadata []*types.BadgeMetadataTimeline, newBadgeMetadata []*types.BadgeMetadataTimeline, canUpdateBadgeMetadata []*types.TimedUpdateWithBadgeIdsPermission) error {
	oldTimes, oldValues := GetBadgeMetadataTimesAndValues(oldBadgeMetadata)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetBadgeMetadataTimesAndValues(newBadgeMetadata)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, []*types.BadgeMetadata{}, func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldBadgeMetadata := oldValue.([]*types.BadgeMetadata)
		firstMatchesForOld := types.GetFirstMatchOnly(CastBadgeMetadataToUniversalPermission(oldBadgeMetadata))

		newBadgeMetadata := newValue.([]*types.BadgeMetadata)
		firstMatchesForNew := types.GetFirstMatchOnly(CastBadgeMetadataToUniversalPermission(newBadgeMetadata))

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			oldVal := oldDetails.ArbitraryValue.(string)
			newVal := newDetails.ArbitraryValue.(string)
			
			if newVal != oldVal {
				detailsToReturn = append(detailsToReturn, overlap)
			}
		}

		//If metadata is in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn
	})

	err := CheckTimedUpdateWithBadgeIdsPermission(ctx, detailsToCheck, canUpdateBadgeMetadata)
	if err != nil {
		return err
	}

	return nil
}

func ValidateCollectionMetadataUpdate(ctx sdk.Context, oldCollectionMetadata []*types.CollectionMetadataTimeline, newCollectionMetadata []*types.CollectionMetadataTimeline, canUpdateCollectionMetadata []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := GetCollectionMetadataTimesAndValues(oldCollectionMetadata)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetCollectionMetadataTimesAndValues(newCollectionMetadata)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, &types.CollectionMetadata{}, func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
		detailsToCheck := []*types.UniversalPermissionDetails{}
		oldVal := oldValue.(*types.CollectionMetadata)
		newVal := newValue.(*types.CollectionMetadata)

		if oldVal.Uri != newVal.Uri || oldVal.CustomData != newVal.CustomData {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{})
		}
		return detailsToCheck
	})

	
	err := CheckTimedUpdatePermission(ctx, detailsToCheck, canUpdateCollectionMetadata)
	if err != nil {
		return err
	}

	return nil
}

func ValidateOffChainBalancesMetadataUpdate(ctx sdk.Context, oldOffChainBalancesMetadata []*types.OffChainBalancesMetadataTimeline, newOffChainBalancesMetadata []*types.OffChainBalancesMetadataTimeline, canUpdateOffChainBalancesMetadata []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := GetOffChainBalancesMetadataTimesAndValues(oldOffChainBalancesMetadata)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetOffChainBalancesMetadataTimesAndValues(newOffChainBalancesMetadata)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, &types.OffChainBalancesMetadata{}, func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
		detailsToCheck := []*types.UniversalPermissionDetails{}
		oldVal := oldValue.(*types.OffChainBalancesMetadata)
		newVal := newValue.(*types.OffChainBalancesMetadata)

		if oldVal.Uri != newVal.Uri || oldVal.CustomData != newVal.CustomData {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{})
		}
		return detailsToCheck
	})

	err := CheckTimedUpdatePermission(ctx, detailsToCheck, canUpdateOffChainBalancesMetadata)
	if err != nil {
		return err
	}

	return nil
}

func ValidateInheritedBalancesUpdate(ctx sdk.Context, oldInheritedBalances []*types.InheritedBalancesTimeline, newInheritedBalances []*types.InheritedBalancesTimeline, canUpdateInheritedBalances []*types.TimedUpdateWithBadgeIdsPermission) error {
	oldTimes, oldValues := GetInheritedBalancesTimesAndValues(oldInheritedBalances)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetInheritedBalancesTimesAndValues(newInheritedBalances)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	detailsToCheck := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, []*types.InheritedBalance{}, func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
		//Cast to UniversalPermissionDetails for comaptibility with these overlap functions and get first matches only (i.e. first match for each badge ID)
		oldInheritedBalances := oldValue.([]*types.InheritedBalance)
		firstMatchesForOld := types.GetFirstMatchOnly(CastInheritedBalancesToUniversalPermission(oldInheritedBalances))

		newInheritedBalances := newValue.([]*types.InheritedBalance)
		firstMatchesForNew := types.GetFirstMatchOnly(CastInheritedBalancesToUniversalPermission(newInheritedBalances))

		//For every badge, we need to check if the new provided value is different in any way from the old value for each badge ID
		//The overlapObjects from GetOverlapsAndNonOverlaps will return which badge IDs overlap
		detailsToReturn := []*types.UniversalPermissionDetails{}
		overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
		for _, overlapObject := range overlapObjects {
			overlap := overlapObject.Overlap
			oldDetails := overlapObject.FirstDetails
			newDetails := overlapObject.SecondDetails
			oldVal := oldDetails.ArbitraryValue.([]*types.InheritedBalance)
			newVal := newDetails.ArbitraryValue.([]*types.InheritedBalance)

			different := false
			if len(newVal) != len(oldVal) {
				different = true
			} else {
				for i := 0; i < len(newVal); i++ {
					x, err := proto.Marshal(newVal[i])
					if err != nil {
						panic(err)
					}

					y, err := proto.Marshal(oldVal[i])
					if err != nil {
						panic(err)
					}

					if string(x) != string(y) {
						different = true
						break
					}
				}
			}

			if different {
				detailsToReturn = append(detailsToReturn, overlap)
			}
		}

		//If metadata is in old but not new, then it is considered updated. If it is in new but not old, then it is considered updated.
		detailsToReturn = append(detailsToReturn, inOldButNotNew...)
		detailsToReturn = append(detailsToReturn, inNewButNotOld...)

		return detailsToReturn
	})

	err := CheckTimedUpdateWithBadgeIdsPermission(ctx, detailsToCheck, canUpdateInheritedBalances)
	if err != nil {
		return err
	}

	return nil
}

func ValidateManagerUpdate(ctx sdk.Context, oldManager []*types.ManagerTimeline, newManager []*types.ManagerTimeline, canUpdateManager []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := GetManagerTimesAndValues(oldManager)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetManagerTimesAndValues(newManager)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	updatedTimelineTimes := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, "", GetCombosForStringTimelines)

	if err := CheckTimedUpdatePermission(ctx, updatedTimelineTimes, canUpdateManager); err != nil {
		return err
	}

	return nil
}

func ValidateCustomDataUpdate(ctx sdk.Context, oldCustomData []*types.CustomDataTimeline, newCustomData []*types.CustomDataTimeline, canUpdateCustomData []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := GetCustomDataTimesAndValues(oldCustomData)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetCustomDataTimesAndValues(newCustomData)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	updatedTimelineTimes := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, "", GetCombosForStringTimelines)

	if err := CheckTimedUpdatePermission(ctx, updatedTimelineTimes, canUpdateCustomData); err != nil {
		return err
	}

	return nil
}

func ValidateStandardsUpdate(ctx sdk.Context, oldStandards []*types.StandardTimeline, newStandards []*types.StandardTimeline, canUpdateStandards []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := GetStandardsTimesAndValues(oldStandards)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetStandardsTimesAndValues(newStandards)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	updatedTimelineTimes := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, "", func (oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
		oldVal := oldValue.([]string)
		newVal := newValue.([]string)
		
		if len(oldVal) != len(newVal) {
			return []*types.UniversalPermissionDetails{
				{},
			}
		} else {
			for i := 0; i < len(oldVal); i++ {
				if oldVal[i] != newVal[i] {
					return []*types.UniversalPermissionDetails{
						{},
					}
				}
			}
		}

		return []*types.UniversalPermissionDetails{}
	})

	if err := CheckTimedUpdatePermission(ctx, updatedTimelineTimes, canUpdateStandards); err != nil {
		return err
	}

	return nil
}

func ValidateContractAddressUpdate(ctx sdk.Context, oldContractAddress []*types.ContractAddressTimeline, newContractAddress []*types.ContractAddressTimeline, canUpdateContractAddress []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := GetContractAddressTimesAndValues(oldContractAddress)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetContractAddressTimesAndValues(newContractAddress)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	updatedTimelineTimes := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, "", GetCombosForStringTimelines)

	if err := CheckTimedUpdatePermission(ctx, updatedTimelineTimes, canUpdateContractAddress); err != nil {
		return err
	}

	return nil
}

func ValidateIsArchivedUpdate(ctx sdk.Context, oldIsArchived []*types.IsArchivedTimeline, newIsArchived []*types.IsArchivedTimeline, canUpdateIsArchived []*types.TimedUpdatePermission) error {
	oldTimes, oldValues := GetIsArchivedTimesAndValues(oldIsArchived)
	oldTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(oldTimes, oldValues)

	newTimes, newValues := GetIsArchivedTimesAndValues(newIsArchived)
	newTimelineFirstMatches := GetPotentialUpdatesForTimelineValues(newTimes, newValues)

	updatedTimelineTimes := GetUpdateCombinationsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, false, func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
		oldVal := oldValue.(bool)
		newVal := newValue.(bool)
		if oldVal != newVal {
			return []*types.UniversalPermissionDetails{
				{},
			}
		}
		return []*types.UniversalPermissionDetails{}
	})

	if err := CheckTimedUpdatePermission(ctx, updatedTimelineTimes, canUpdateIsArchived); err != nil {
		return err
	}

	return nil
}