package keeper

import "github.com/bitbadges/bitbadgeschain/x/badges/types"

func AssertIsFrozenLogicIsMaintained(prevBadgeMetadata []*types.BadgeMetadata, newBadgeMetadata []*types.BadgeMetadata) error {
	//Check badge metadata for isFrozen logic. If previously frozen, must be frozen now.
	//TODO: And must be the same
	// for idx, badgeMetadata := range prevBadgeMetadata {
	// 	if badgeMetadata.IsFrozen {
	// 		if len(newBadgeMetadata) <= idx {
	// 			return ErrBadgeMetadataMustBeFrozen
	// 		}

	// 		if !newBadgeMetadata[idx].IsFrozen {
	// 			return ErrBadgeMetadataMustBeFrozen
	// 		}
	// 	}
	// }

	// //Check to make sure that only first X are frozen
	// stillFrozen := true
	// for _, badgeMetadata := range newBadgeMetadata {
	// 	if badgeMetadata.IsFrozen {
	// 		if !stillFrozen {
	// 			return ErrBadgeMetadataMustBeFrozen
	// 		}
	// 	} else {
	// 		stillFrozen = false
	// 	}
	// }

	return nil
}

// func GetUrisToStoreAndPermissionsToCheck(collection *types.BadgeCollection, msgCollectionMetadata *types.CollectionMetadata, msgBadgeMetadata []*types.BadgeMetadata, msgOffChainBalancesMetadata *types.OffChainBalancesMetadata) (newCollectionMetadata *types.CollectionMetadata, newBadgeMetadata []*types.BadgeMetadata, newBalanceUri *types.OffChainBalancesMetadata, needToValidateUpdateCollectionMetadata bool, needToValidateUpdateBadgeMetadata bool, needToValidateUpdateBalanceUri bool) {
// 	needToValidateUpdateCollectionMetadata = false
// 	needToValidateUpdateBalanceUri = false
// 	needToValidateUpdateBadgeMetadata = false

// 	newCollectionMetadata = collection.CollectionMetadata
// 	newBadgeMetadata = collection.BadgeMetadata
// 	newBalanceUri = collection.OffChainBalancesMetadata

// 	if msgCollectionMetadata != nil &&
// 		(msgCollectionMetadata.Uri != collection.CollectionMetadata.Uri || msgCollectionMetadata.CustomData != collection.CollectionMetadata.CustomData) {
// 		needToValidateUpdateCollectionMetadata = true
// 		newCollectionMetadata = msgCollectionMetadata
// 	}

// 	if msgOffChainBalancesMetadata != nil &&
// 		(msgOffChainBalancesMetadata.Uri != collection.OffChainBalancesMetadata.Uri || msgOffChainBalancesMetadata.CustomData != collection.OffChainBalancesMetadata.CustomData) {
// 		needToValidateUpdateBalanceUri = true
// 		newBalanceUri = msgOffChainBalancesMetadata
// 	}

// 	if msgBadgeMetadata != nil && len(msgBadgeMetadata) > 0 {
// 		newBadgeMetadata = msgBadgeMetadata

// 		for idx, badgeMetadata := range collection.BadgeMetadata {
// 			if idx >= len(msgBadgeMetadata) {
// 				needToValidateUpdateBadgeMetadata = true
// 				break
// 			}

// 			if msgBadgeMetadata[idx].Uri != badgeMetadata.Uri || msgBadgeMetadata[idx].CustomData != badgeMetadata.CustomData {
// 				needToValidateUpdateBadgeMetadata = true
// 				break
// 			}

// 			if len(msgBadgeMetadata[idx].BadgeIds) != len(badgeMetadata.BadgeIds) {
// 				needToValidateUpdateBadgeMetadata = true
// 				break
// 			}

// 			for j, badgeIdRange := range badgeMetadata.BadgeIds {
// 				if badgeIdRange.Start != msgBadgeMetadata[idx].BadgeIds[j].Start || badgeIdRange.End != msgBadgeMetadata[idx].BadgeIds[j].End {
// 					needToValidateUpdateBadgeMetadata = true
// 					break
// 				}
// 			}
// 		}
// 	}

// 	return newCollectionMetadata, newBadgeMetadata, newBalanceUri, needToValidateUpdateCollectionMetadata, needToValidateUpdateBadgeMetadata, needToValidateUpdateBalanceUri
// }
