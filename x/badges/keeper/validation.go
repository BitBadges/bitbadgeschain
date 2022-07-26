package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Determines what to validate for each Msg
type UniversalValidationParams struct {
	Creator                      string
	BadgeId                      uint64
	AccountsThatCantEqualCreator []uint64
	SubbadgeRangesToValidate     []*types.IdRange
	AccountsToCheckRegistration  []uint64
	MustBeManager                bool
	CanFreeze                    bool
	CanCreateSubbadges           bool
	CanRevoke                    bool
	CanManagerTransfer           bool
	CanUpdateUris                bool
	CanUpdateBytes               bool
}

// Validates everything about the Msg is valid and returns (creatorNum, badge, permissions, error).
func (k Keeper) UniversalValidate(ctx sdk.Context, params UniversalValidationParams) (uint64, types.BitBadge, error) {
	CreatorAccountNum := k.MustGetAccountNumberForBech32AddressString(ctx, params.Creator)

	if len(params.AccountsThatCantEqualCreator) > 0 {
		for _, account := range params.AccountsThatCantEqualCreator {
			if account == CreatorAccountNum {
				return CreatorAccountNum, types.BitBadge{}, ErrAccountCanNotEqualCreator
			}
		}
	}

	if len(params.AccountsToCheckRegistration) > 0 {
		nextAccountNumber := k.accountKeeper.GetNextAccountNumber(ctx) //TODO: this increments each time we call this; we just want to get w/o incrementing 
		for _, accountNumber := range params.AccountsToCheckRegistration {
			//Probably a better way to do this such as only read once at beginning of block, but we check that addresses are valid and not > Next Account Number because then we would be sending to an unregistered address
			if accountNumber >= nextAccountNumber {
				return CreatorAccountNum, types.BitBadge{}, ErrAccountNotRegistered
			}
		}
	}

	// Assert badge and subbadge ranges exist and are well-formed
	badge, err := k.GetBadgeAndAssertSubbadgeRangesAreValid(ctx, params.BadgeId, params.SubbadgeRangesToValidate)
	if err != nil {
		return CreatorAccountNum, types.BitBadge{}, err
	}

	// Assert all permissions
	if params.MustBeManager && badge.Manager != CreatorAccountNum {
		return CreatorAccountNum, types.BitBadge{}, ErrSenderIsNotManager
	}

	permissions := types.GetPermissions(badge.Permissions)
	if params.CanFreeze && !permissions.CanFreeze {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	if params.CanCreateSubbadges && !permissions.CanCreate {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	if params.CanRevoke && !permissions.CanRevoke {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	if params.CanManagerTransfer && !permissions.CanManagerTransfer {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	if params.CanUpdateUris && !permissions.CanUpdateUris {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	if params.CanUpdateBytes && !permissions.CanUpdateBytes {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	return CreatorAccountNum, badge, nil
}
