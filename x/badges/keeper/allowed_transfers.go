package keeper

import "github.com/bitbadges/bitbadgeschain/x/badges/types"

func AssertIsFrozenLogicForApprovedTransfers(prevApprovedTransfers []*types.CollectionApprovedTransfer, newApprovedTransfers []*types.CollectionApprovedTransfer) error {
	//Check all previous allowed transfers to remain frozen
	for idx, prevAllowedTransfer := range prevApprovedTransfers {
		if prevAllowedTransfer.IsFrozen {
			if len(newApprovedTransfers) <= idx {
				return ErrApprovedTransfersMustBeFrozen
			}

			if !newApprovedTransfers[idx].IsFrozen {
				return ErrApprovedTransfersMustBeFrozen
			}
		}
	}

	//Check to make sure that if any allowed transfers are frozen, all previous allowed transfers are frozen
	foundNonFrozen := false
	for _, allowedTransfer := range newApprovedTransfers {
		if allowedTransfer.IsFrozen {
			if foundNonFrozen {
				return ErrApprovedTransfersMustBeFrozen
			}
		} else {
			foundNonFrozen = true
		}
	}

	return nil
}

func GetApprovedTransfersToStore(collection types.BadgeCollection, msgApprovedTransfers []*types.CollectionApprovedTransfer) ([]*types.CollectionApprovedTransfer, bool) {
	needToValidateUpdateCollectionApprovedTransfers := false
	newApprovedTransfers := collection.ApprovedTransfers
	if msgApprovedTransfers != nil && len(msgApprovedTransfers) > 0 {
		newApprovedTransfers = msgApprovedTransfers

		for idx, allowedTransfer := range collection.ApprovedTransfers {
			if idx >= len(msgApprovedTransfers) {
				needToValidateUpdateCollectionApprovedTransfers = true
				break
			}

			//TODO: This may work or may not. We need to deep equals compare these
			if allowedTransfer != nil && allowedTransfer != msgApprovedTransfers[idx] {
				needToValidateUpdateCollectionApprovedTransfers = true
				break
			}
		}
	}

	return newApprovedTransfers, needToValidateUpdateCollectionApprovedTransfers, nil
}
