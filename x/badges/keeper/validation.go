package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Determines what to validate for each Msg
type UniversalValidationParams struct {
	Creator                      string
	CollectionId                 sdk.Uint
	AccountsThatCantEqualCreator []string
	BadgeIdRangesToValidate      []*types.IdRange
	MustBeManager                bool
	CanFreeze                    bool
	CanCreateMoreBadges          bool
	CanUpdateAllowed          	 bool
	CanManagerBeTransferred      bool
	CanUpdateMetadataUris        bool
	CanUpdateBytes               bool
	CanDelete                    bool
	CanUpdateBalancesUri				 bool
}

// Validates everything about the Msg is valid and returns (creatorNum, badge, permissions, error).
func (k Keeper) UniversalValidate(ctx sdk.Context, params UniversalValidationParams) (types.BadgeCollection, error) {
	
	if len(params.AccountsThatCantEqualCreator) > 0 {
		for _, account := range params.AccountsThatCantEqualCreator {
			if account == params.Creator {
				return types.BadgeCollection{}, ErrAccountCanNotEqualCreator
			}
		}
	}
	// if len(params.AccountsToCheckRegistration) > 0 {
		// We have three options here. I currently think doing nothing is what I will go with.

		// 1. Do nothing and put blame on users
		// 2. Check if account exists through GetNextAccountNumber() and if not, return error. Less Gas. Increments account number each time.
		//Since account number is incremented each time, there will be gaps in the account numbers, and that blame is placed on the users.
		// 3. Check if account exists through HasAccountAddressByID() and if not, return error. More Gas

		//Option 2
		// nextAccountNumber := k.accountKeeper.GetNextAccountNumber(ctx)
		// for _, accountNumber := range params.AccountsToCheckRegistration {
		// 	//Probably a better way to do this such as only read once at beginning of block, but we check that addresses are valid and not > Next Account Number because then we would be sending to an unregistered address
		// 	if accountNumber >= nextAccountNumber {
		// 		return types.BadgeCollection{}, ErrAccountNotRegistered
		// 	}
		// }

		//Option 3
		// for _, accountNumber := range params.AccountsToCheckRegistration {
		// 	if !k.accountKeeper.HasAccountAddressByID(ctx, accountNumber) {
		// 		return types.BadgeCollection{}, ErrAccountNotRegistered
		// 	}
		// }
	// }

	// Assert collection and badgeId ranges exist and are well-formed
	badge, err := k.GetCollectionAndAssertBadgeIdsAreValid(ctx, params.CollectionId, params.BadgeIdRangesToValidate)
	if err != nil {
		return types.BadgeCollection{}, err
	}

	// Assert all permissions
	if params.MustBeManager && badge.Manager != params.Creator {
		return types.BadgeCollection{}, ErrSenderIsNotManager
	}

	permissions := types.GetPermissions(badge.Permissions)
	if params.CanUpdateAllowed && !permissions.CanUpdateAllowed {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanCreateMoreBadges && !permissions.CanCreateMoreBadges {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanManagerBeTransferred && !permissions.CanManagerBeTransferred {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanUpdateMetadataUris && !permissions.CanUpdateMetadataUris {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanUpdateBytes && !permissions.CanUpdateBytes {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanDelete && !permissions.CanDelete {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanUpdateBalancesUri && !permissions.CanUpdateBalancesUri {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	return badge, nil
}
