package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) TransferManager(goCtx context.Context, msg *types.MsgTransferManager) (*types.MsgTransferManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)
	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")

	if err := k.AssertAccountNumbersAreRegistered(ctx, []uint64{msg.Address}); err != nil {
		return nil, ErrAccountsAreNotRegistered
	}

	badge, found := k.GetBadgeFromStore(ctx, msg.BadgeId)
	if !found {
		return nil, ErrBadgeNotExists
	}

	//Transfer to new manager but first need to check privileges
	if badge.Manager != CreatorAccountNum {
		return nil, ErrSenderIsNotManager
	}

	permissions := types.GetPermissions(badge.PermissionFlags)
	if !permissions.CanManagerTransfer() {
		return nil, ErrInvalidPermissions
	}

	ctx.GasMeter().ConsumeGas(TransferManagerCost, "transfer manager cost")
	requested := k.HasAddressRequestedManagerTransfer(ctx, msg.BadgeId, msg.Address)
	if !requested {
		return nil, ErrAddressNeedsToOptInAndRequestManagerTransfer
	}

	//TODO: other permissions such as remove force mint, etc

	badge.Manager = msg.Address

	if err := k.RemoveTransferManagerRequest(ctx, msg.BadgeId, msg.Address); err != nil {
		return nil, err
	}

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "TransferManager"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("NewManager", fmt.Sprint(msg.Address)),
		),
	)

	return &types.MsgTransferManagerResponse{}, nil
}
