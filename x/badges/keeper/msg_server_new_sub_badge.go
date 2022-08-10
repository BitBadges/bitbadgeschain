package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) NewSubBadge(goCtx context.Context, msg *types.MsgNewSubBadge) (*types.MsgNewSubBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validationParams := UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.Id,
		MustBeManager: true,
		CanCreateSubbadges: true,
	}

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, validationParams)
	if err != nil {
		return nil, err
	}
	

	ManagerBalanceKey := ConstructBalanceKey(CreatorAccountNum, msg.Id)
	badgeBalanceInfo, found := k.GetBadgeBalanceFromStore(ctx, ManagerBalanceKey)
	if !found {
		badgeBalanceInfo = types.BadgeBalanceInfo{}
	}

	originalSubassetId := badge.NextSubassetId
	new_amounts := badge.SubassetsTotalSupply
	for i, supply := range msg.Supplys {
		for j := uint64(0); j < msg.AmountsToCreate[i]; j++ {
			//Once here, we should be safe to mint
			//We don't need to store if subbadge supply == default
			subasset_id := badge.NextSubassetId
			defaultSupply := badge.DefaultSubassetSupply
			if badge.DefaultSubassetSupply == 0 {
				defaultSupply = 1
			}

			//default to supply = default when supply is 0
			if supply == 0 {
				supply = defaultSupply
			}

			if supply != defaultSupply {
				ctx.GasMeter().ConsumeGas(SubbadgeWithSupplyNotEqualToOne, "create new subbadge cost")

				new_amounts = UpdateBalanceForSubbadgeId(subasset_id, supply, new_amounts)
			}
			badge.NextSubassetId += 1

			//Mint the total supply of subbadge to the manager
			newBadgeBalanceInfo, err := k.AddBalanceForSubbadgeId(ctx, badgeBalanceInfo, subasset_id, supply)
			if err != nil {
				return nil, err
			}
			badgeBalanceInfo = newBadgeBalanceInfo
		}
	}
	badge.SubassetsTotalSupply = new_amounts

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	if err := k.SetBadgeBalanceInStore(ctx, ManagerBalanceKey, badgeBalanceInfo); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "CreatedSubBadges"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("FirstCreatedID", fmt.Sprint(originalSubassetId)),
			sdk.NewAttribute("LastCreatedID", fmt.Sprint(badge.NextSubassetId-1)),
		),
	)

	return &types.MsgNewSubBadgeResponse{
		SubassetId: badge.NextSubassetId,
	}, nil
}
