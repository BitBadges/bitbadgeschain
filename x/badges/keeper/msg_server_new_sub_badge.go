package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) NewSubBadge(goCtx context.Context, msg *types.MsgNewSubBadge) (*types.MsgNewSubBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.Id,
		MustBeManager: true,
		CanCreateSubbadges: true,
	})
	if err != nil {
		return nil, err
	}
	
	managerBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.Id)
	managerBalanceInfo, found := k.GetUserBalanceFromStore(ctx, managerBalanceKey)
	if !found {
		managerBalanceInfo = types.UserBalanceInfo{}
	}

	originalSubassetId := badge.NextSubassetId
	newSubassetSupplys := badge.SubassetSupplys
	defaultSupply := badge.DefaultSubassetSupply
	if badge.DefaultSubassetSupply == 0 {
		defaultSupply = 1
	}

	// Update supplys and mint total supply for each to manager. Don't store if supply == default
	for i, supply := range msg.Supplys {
		for j := uint64(0); j < msg.AmountsToCreate[i]; j++ {
			nextSubassetId := badge.NextSubassetId

			// We conventionalize supply == 0 as default, so we don't store if it is the default
			if supply != 0 && supply != defaultSupply {
				ctx.GasMeter().ConsumeGas(SubbadgeWithSupplyNotEqualToOne, "create new subbadge cost")

				newSubassetSupplys = UpdateBalanceForId(nextSubassetId, supply, newSubassetSupplys)
			}
			badge.NextSubassetId += 1

			managerBalanceInfo, err = AddBalanceForId(ctx, managerBalanceInfo, nextSubassetId, supply)
			if err != nil {
				return nil, err
			}
		}
	}
	badge.SubassetSupplys = newSubassetSupplys

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	if err := k.SetUserBalanceInStore(ctx, managerBalanceKey, managerBalanceInfo); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "CreatedSubBadges"),
			sdk.NewAttribute("BadgeId", fmt.Sprint(badge.Id)),
			sdk.NewAttribute("FirstId", fmt.Sprint(originalSubassetId)),
			sdk.NewAttribute("LastId", fmt.Sprint(badge.NextSubassetId-1)),
		),
	)

	return &types.MsgNewSubBadgeResponse{
		SubassetId: badge.NextSubassetId,
	}, nil
}
