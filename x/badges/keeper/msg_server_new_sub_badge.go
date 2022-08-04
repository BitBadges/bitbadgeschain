package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) NewSubBadge(goCtx context.Context, msg *types.MsgNewSubBadge) (*types.MsgNewSubBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	badge, found := k.GetBadgeFromStore(ctx, msg.Id)
	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")

	if !found {
		return nil, ErrBadgeNotExists
	}

	if badge.Manager != CreatorAccountNum {
		return nil, ErrSenderIsNotManager
	}

	permissions := types.GetPermissions(badge.PermissionFlags)
	if !permissions.CanCreateSubbadges() {
		return nil, ErrInvalidPermissions
	}
	originalSubassetId := badge.NextSubassetId
	for i, supply := range msg.Supplys {
		for j := uint64(0); j < msg.AmountsToCreate[i]; j++ {
			//Once here, we should be safe to mint
			//By default, we assume non fungible subbadge (i.e. supply == 1) so we don't store if supply == 1
			subasset_id := badge.NextSubassetId
			if supply != 1 {
				ctx.GasMeter().ConsumeGas(SubbadgeWithSupplyNotEqualToOne, "create new subbadge cost")
				hasLastEntry := false
				if len(badge.SubassetsTotalSupply) > 0 {
					hasLastEntry = true
				}

				lastEntry := &types.Subasset{}
				if hasLastEntry {
					lastEntry = badge.SubassetsTotalSupply[len(badge.SubassetsTotalSupply)-1]
				}

				if hasLastEntry && lastEntry.Supply == supply && lastEntry.EndId == subasset_id-1 {
					badge.SubassetsTotalSupply[len(badge.SubassetsTotalSupply)-1] = &types.Subasset{
						Supply:  lastEntry.Supply,
						StartId: lastEntry.StartId,
						EndId:   subasset_id,
					}
				} else {
					badge.SubassetsTotalSupply = append(badge.SubassetsTotalSupply, &types.Subasset{
						Supply:  supply,
						StartId: subasset_id,
						EndId:   subasset_id,
					})
				}
			}
			badge.NextSubassetId += 1

			//Mint the total supply of subbadge to the manager
			ManagerBalanceKey := GetBalanceKey(CreatorAccountNum, msg.Id, subasset_id)
			if err := k.AddToBadgeBalance(ctx, ManagerBalanceKey, supply); err != nil {
				return nil, err
			}
		}
	}

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "CreatedSubBadges"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("FirstCreatedID", fmt.Sprint(originalSubassetId)),
			sdk.NewAttribute("LastCreatedID", fmt.Sprint(badge.NextSubassetId - 1)),
		),
	)

	return &types.MsgNewSubBadgeResponse{
		SubassetId: badge.NextSubassetId,
	}, nil
}
