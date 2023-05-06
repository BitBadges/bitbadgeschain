package keeper

import "github.com/bitbadges/bitbadgeschain/x/badges/types"

func GetUrisToStoreAndPermissionsToCheck(collection types.BadgeCollection, msgCollectionUri string, msgBadgeUris []*types.BadgeUri, msgBalancesUri string) (newCollectionUri string, newBadgeUris []*types.BadgeUri, newBalanceUri string, needToValidateUpdateMetadataUris bool, needToValidateUpdateBalanceUri bool) {
	needToValidateUpdateMetadataUris = false
	needToValidateUpdateBalanceUri = false

	newCollectionUri = collection.CollectionUri
	newBadgeUris = collection.BadgeUris
	newBalanceUri = collection.BalancesUri

	if msgCollectionUri != "" && msgCollectionUri != collection.CollectionUri {
		needToValidateUpdateMetadataUris = true
		newCollectionUri = msgCollectionUri
	}

	if msgBalancesUri != "" && msgBalancesUri != collection.BalancesUri {
		needToValidateUpdateBalanceUri = true
		newBalanceUri = msgBalancesUri
	}

	if msgBadgeUris != nil && len(msgBadgeUris) > 0 {
		newBadgeUris = msgBadgeUris

		for idx, badgeUri := range collection.BadgeUris {
			if msgBadgeUris[idx].Uri != badgeUri.Uri {
				needToValidateUpdateMetadataUris = true
				break
			}

			if len(msgBadgeUris[idx].BadgeIds) != len(badgeUri.BadgeIds) {
				needToValidateUpdateMetadataUris = true
				break
			}

			for j, badgeIdRange := range badgeUri.BadgeIds {
				if badgeIdRange.Start != msgBadgeUris[idx].BadgeIds[j].Start || badgeIdRange.End != msgBadgeUris[idx].BadgeIds[j].End {
					needToValidateUpdateMetadataUris = true
					break
				}
			}
		}
	}

	return newCollectionUri, newBadgeUris, newBalanceUri, needToValidateUpdateMetadataUris, needToValidateUpdateBalanceUri
}