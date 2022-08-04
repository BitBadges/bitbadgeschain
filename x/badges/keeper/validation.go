package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k Keeper) UniversalValidateMsgAndReturnMsgInfo(ctx sdk.Context, MsgCreator string, AddressesToValidate []uint64, BadgeId uint64, SubbadgeId uint64, MustBeManager bool) (uint64, types.BitBadge, types.PermissionFlags, error) {
	CreatorAccountNum := k.MustGetAccountNumberForBech32AddressString(ctx, MsgCreator)

	if err := k.AssertAccountNumbersAreRegistered(ctx, AddressesToValidate); err != nil {
		return CreatorAccountNum, types.BitBadge{}, types.PermissionFlags{}, err
	}

	badge, err := k.AssertBadgeAndSubBadgeExistsAndReturnBadge(ctx, BadgeId, SubbadgeId)
	if err != nil {
		return CreatorAccountNum, types.BitBadge{}, types.PermissionFlags{}, err
	}

	if MustBeManager && badge.Manager != CreatorAccountNum {
		return CreatorAccountNum, types.BitBadge{}, types.PermissionFlags{}, ErrSenderIsNotManager
	}

	permissions := types.GetPermissions(badge.PermissionFlags)

	return CreatorAccountNum, badge, permissions, nil
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
		postApprovalBadgeBalanceInfo, err := k.RemoveBalanceFromApproval(ctx, newBadgeBalanceInfo, amount, requester, subbadgeId) //if pending and cancelled, this approval will be added back
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
		unfrozen_addresses := badge.FreezeAddresses
		for _, unfrozen_address := range unfrozen_addresses {
			if unfrozen_address == address {
				can_transfer = true
			}
		}
	} else {
		frozen_addresses := badge.FreezeAddresses
		can_transfer = true
		for _, frozen_address := range frozen_addresses {
			if frozen_address == address {
				can_transfer = false
			}
		}
	}

	return can_transfer
}
