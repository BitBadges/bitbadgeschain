package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) PruneBalances(goCtx context.Context, msg *types.MsgPruneBalances) (*types.MsgPruneBalancesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	nextBadgeId := k.GetNextAssetId(ctx)
	
	//Anyone should be able to call this to clear the network

	for _, badgeId := range msg.BadgeIds {
		if badgeId < nextBadgeId && !k.StoreHasBadgeID(ctx, badgeId) {
			ctx.BlockGasMeter().RefundGas(PruneBalanceRefundAmountPerBadge, "prune balances")
			for _, address := range msg.Addresses {
				k.DeleteBadgeBalanceFromStore(ctx, GetBalanceKey(address, badgeId))
				ctx.BlockGasMeter().RefundGas(PruneBalanceRefundAmountPerAddress, "prune balances")
			}
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "PrunedBalances"),
			sdk.NewAttribute("Addresses", fmt.Sprint(msg.Addresses)),
			sdk.NewAttribute("BadgeIds", fmt.Sprint(msg.BadgeIds)),
		),
	)

	return &types.MsgPruneBalancesResponse{}, nil
}
