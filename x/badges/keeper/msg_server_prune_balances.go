package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) PruneBalances(goCtx context.Context, msg *types.MsgPruneBalances) (*types.MsgPruneBalancesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	nextBadgeId := k.GetNextBadgeId(ctx)

	// Don't have to call UniversalValidate(). Anyone should be able to call this to prune unnecessary balances.

	// For every (badgeId, address) pair, make sure the badge has been self destructed, and then delete the balance.
	for _, badgeId := range msg.BadgeIds {
		if badgeId < nextBadgeId && !k.StoreHasBadgeID(ctx, badgeId) {
			ctx.BlockGasMeter().RefundGas(PruneBalanceRefundAmountPerBadge, "prune balances refund per badge")
			for _, address := range msg.Addresses {
				k.DeleteUserBalanceFromStore(ctx, ConstructBalanceKey(address, badgeId))
				ctx.BlockGasMeter().RefundGas(PruneBalanceRefundAmountPerAddress, "prune balances refund per address")
			}
		} else {
			return nil, ErrBadgeNotExists
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "PrunedBalances"),
		),
	)

	return &types.MsgPruneBalancesResponse{}, nil
}
