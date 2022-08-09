package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Determines what to check universally for each Msg
type UniversalValidationParams struct {
	Creator 					string
	BadgeId 					uint64
	AccountsToCheckIfRegistered []uint64
	AccountsThatCantEqualCreator[]uint64
	SubbadgeRangesToValidate	[]*types.NumberRange
	MustBeManager				bool
	CanFreeze					bool
	CanCreateSubbadges			bool
	CanRevoke					bool
	CanManagerTransfer			bool
	CanUpdateUris				bool
}

// Validates everything about the Msg is valid and returns (creatorNum, badge, permissions, error). 
func (k Keeper) UniversalValidate(ctx sdk.Context, params UniversalValidationParams) (uint64, types.BitBadge, error) {
	CreatorAccountNum := k.MustGetAccountNumberForBech32AddressString(ctx, params.Creator)

	// Check if accounts are registered
	if len(params.AccountsToCheckIfRegistered) > 0 {
		if err := k.AssertAccountNumbersAreRegistered(ctx, params.AccountsToCheckIfRegistered); err != nil {
			return CreatorAccountNum, types.BitBadge{}, err
		}
	}	

	if len(params.AccountsThatCantEqualCreator) > 0 {
		for _, account := range params.AccountsThatCantEqualCreator {
			if account == CreatorAccountNum {
				return CreatorAccountNum, types.BitBadge{}, ErrAccountCanNotEqualCreator
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

	permissions := types.GetPermissions(badge.PermissionFlags)
	if params.CanFreeze && !permissions.CanFreeze() {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	if params.CanCreateSubbadges && !permissions.CanCreateSubbadges() {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	if params.CanRevoke && !permissions.CanRevoke() {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	if params.CanManagerTransfer && !permissions.CanManagerTransfer() {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	if params.CanUpdateUris && !permissions.CanUpdateUris() {
		return CreatorAccountNum, types.BitBadge{}, ErrInvalidPermissions
	}

	return CreatorAccountNum, badge, nil
}


func (k Keeper) HandlePreTransfer(ctx sdk.Context, badgeBalanceInfo types.BadgeBalanceInfo, badge types.BitBadge, badgeId uint64, subbadgeId uint64, from uint64, to uint64, requester uint64, amount uint64) (types.BadgeBalanceInfo, error) {
	newBadgeBalanceInfo := badgeBalanceInfo
	permissions := types.GetPermissions(badge.PermissionFlags)

	can_transfer := AccountNotFrozen(badge, permissions, from)
	if !can_transfer {
		return badgeBalanceInfo, ErrAddressFrozen
	}

	// Check and handle approvals if requester != from
	if from != requester {
		postApprovalBadgeBalanceInfo, err := k.RemoveBalanceFromApproval(ctx, newBadgeBalanceInfo, amount, requester, types.NumberRange{Start: subbadgeId, End: subbadgeId}) //if pending and cancelled, this approval will be added back
		newBadgeBalanceInfo = postApprovalBadgeBalanceInfo
		if err != nil {
			return badgeBalanceInfo, err
		}
	}

	return newBadgeBalanceInfo, nil
}

func AccountNotFrozen(badge types.BitBadge, permissions types.PermissionFlags, address uint64) bool {
	frozen_by_default := permissions.FrozenByDefault()

	can_transfer := false
	if frozen_by_default {
		unfrozen_address_ranges := badge.FreezeAddressRanges
		for _, unfrozen_address_range := range unfrozen_address_ranges {
			if unfrozen_address_range.Start <= address && unfrozen_address_range.End >= address {
				can_transfer = true
			} else if unfrozen_address_range.Start == address && unfrozen_address_range.End == 0 {
				can_transfer = true
			}
		}
	} else {
		frozen_address_ranges := badge.FreezeAddressRanges
		can_transfer = true
		for _, frozen_address_ranges := range frozen_address_ranges {
			if frozen_address_ranges.Start <= address && frozen_address_ranges.End >= address {
				can_transfer = false
			} else if frozen_address_ranges.Start == address && frozen_address_ranges.End == 0 {
				can_transfer = false
			}
		}
	}

	return can_transfer
}
