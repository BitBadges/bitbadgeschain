package types

func CastUserApprovedTransferPermissionToUniversalPermission(UserApprovedTransferPermission []*UserApprovedTransferPermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, UserApprovedTransferPermission := range UserApprovedTransferPermission {
		castedCombinations := []*UniversalCombination{}
		for _, UserApprovedOutgoingTransferCombination := range UserApprovedTransferPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: UserApprovedOutgoingTransferCombination.BadgeIdsOptions,
				PermittedTimesOptions: UserApprovedOutgoingTransferCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: UserApprovedOutgoingTransferCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: UserApprovedOutgoingTransferCombination.TimelineTimesOptions,
				TransferTimesOptions: UserApprovedOutgoingTransferCombination.TransferTimesOptions,
				ToMappingIdOptions: UserApprovedOutgoingTransferCombination.ToMappingIdOptions,
				InitiatedByMappingIdOptions: UserApprovedOutgoingTransferCombination.InitiatedByMappingIdOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				BadgeIds: UserApprovedTransferPermission.DefaultValues.BadgeIds,
				TimelineTimes: UserApprovedTransferPermission.DefaultValues.TimelineTimes,
				TransferTimes: UserApprovedTransferPermission.DefaultValues.TransferTimes,
				ToMappingId: UserApprovedTransferPermission.DefaultValues.ToMappingId,
				InitiatedByMappingId: UserApprovedTransferPermission.DefaultValues.InitiatedByMappingId,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesInitiatedByMappingId: true,
				PermittedTimes: UserApprovedTransferPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: UserApprovedTransferPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}

func CastActionPermissionToUniversalPermission(actionPermission []*ActionPermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, actionPermission := range actionPermission {
		castedCombinations := []*UniversalCombination{}
		for _, actionCombination := range actionPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				PermittedTimesOptions: actionCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: actionCombination.ForbiddenTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				PermittedTimes: actionPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: actionPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}

func CastCollectionApprovedTransferPermissionToUniversalPermission(collectionUpdatePermission []*CollectionApprovedTransferPermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, collectionUpdatePermission := range collectionUpdatePermission {
		castedCombinations := []*UniversalCombination{}
		for _, collectionUpdateCombination := range collectionUpdatePermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				PermittedTimesOptions: collectionUpdateCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: collectionUpdateCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: collectionUpdateCombination.TimelineTimesOptions,
				TransferTimesOptions: collectionUpdateCombination.TransferTimesOptions,
				ToMappingIdOptions: collectionUpdateCombination.ToMappingIdOptions,
				FromMappingIdOptions: collectionUpdateCombination.FromMappingIdOptions,
				InitiatedByMappingIdOptions: collectionUpdateCombination.InitiatedByMappingIdOptions,
				BadgeIdsOptions: collectionUpdateCombination.BadgeIdsOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				TimelineTimes: collectionUpdatePermission.DefaultValues.TimelineTimes,
				TransferTimes: collectionUpdatePermission.DefaultValues.TransferTimes,
				ToMappingId: collectionUpdatePermission.DefaultValues.ToMappingId,
				FromMappingId: collectionUpdatePermission.DefaultValues.FromMappingId,
				InitiatedByMappingId: collectionUpdatePermission.DefaultValues.InitiatedByMappingId,
				BadgeIds: collectionUpdatePermission.DefaultValues.BadgeIds,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesFromMappingId: true,
				UsesInitiatedByMappingId: true,
				PermittedTimes: collectionUpdatePermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: collectionUpdatePermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}


func CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(timedUpdateWithBadgeIdsPermission []*TimedUpdateWithBadgeIdsPermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, timedUpdateWithBadgeIdsPermission := range timedUpdateWithBadgeIdsPermission {
		castedCombinations := []*UniversalCombination{}
		for _, timedUpdateWithBadgeIdsCombination := range timedUpdateWithBadgeIdsPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: timedUpdateWithBadgeIdsCombination.BadgeIdsOptions,
				PermittedTimesOptions: timedUpdateWithBadgeIdsCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: timedUpdateWithBadgeIdsCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: timedUpdateWithBadgeIdsCombination.TimelineTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				TimelineTimes: timedUpdateWithBadgeIdsPermission.DefaultValues.TimelineTimes,
				BadgeIds: timedUpdateWithBadgeIdsPermission.DefaultValues.BadgeIds,
				UsesTimelineTimes: true,
				UsesBadgeIds: true,
				PermittedTimes: timedUpdateWithBadgeIdsPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: timedUpdateWithBadgeIdsPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}

func CastTimedUpdatePermissionToUniversalPermission(timedUpdatePermission []*TimedUpdatePermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, timedUpdatePermission := range timedUpdatePermission {
		castedCombinations := []*UniversalCombination{}
		for _, timedUpdateCombination := range timedUpdatePermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				PermittedTimesOptions: timedUpdateCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: timedUpdateCombination.ForbiddenTimesOptions,
				TimelineTimesOptions: timedUpdateCombination.TimelineTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				TimelineTimes: timedUpdatePermission.DefaultValues.TimelineTimes,
				UsesTimelineTimes: true,
				PermittedTimes: timedUpdatePermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: timedUpdatePermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}


func CastActionWithBadgeIdsPermissionToUniversalPermission(actionWithBadgeIdsPermission []*ActionWithBadgeIdsPermission) []*UniversalPermission {
	castedPermissions := []*UniversalPermission{}
	for _, actionWithBadgeIdsPermission := range actionWithBadgeIdsPermission {
		castedCombinations := []*UniversalCombination{}
		for _, actionWithBadgeIdsCombination := range actionWithBadgeIdsPermission.Combinations {
			castedCombinations = append(castedCombinations, &UniversalCombination{
				BadgeIdsOptions: actionWithBadgeIdsCombination.BadgeIdsOptions,
				PermittedTimesOptions: actionWithBadgeIdsCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: actionWithBadgeIdsCombination.ForbiddenTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &UniversalPermission{
			DefaultValues: &UniversalDefaultValues{
				BadgeIds: actionWithBadgeIdsPermission.DefaultValues.BadgeIds,
				UsesBadgeIds: true,
				PermittedTimes: actionWithBadgeIdsPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: actionWithBadgeIdsPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions
}