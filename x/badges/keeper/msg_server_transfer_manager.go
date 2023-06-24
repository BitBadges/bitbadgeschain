package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ManagerTimelineDetails struct {
	TimelineTime *types.IdRange
	Manager      string
}

func GetManagerTimesAndValues(managerTimeline []*types.ManagerTimeline) ([][]*types.IdRange, []interface{}) {
	times := [][]*types.IdRange{}
	values := []interface{}{}
	for _, timelineVal := range managerTimeline {
		times = append(times, timelineVal.Times)
		values = append(values, timelineVal.Manager)
	}
	return times, values
}

func (k msgServer) UpdateManager(goCtx context.Context, msg *types.MsgUpdateManager) (*types.MsgUpdateManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                 msg.Creator,
		CollectionId:            msg.CollectionId,
		MustBeManager:           true,
	})
	if err != nil {
		return nil, err
	}

	oldTimes, oldValues := GetManagerTimesAndValues(collection.ManagerTimeline)
	oldTimelineFirstMatches := GetFirstMatchOnlyForTimeline(oldTimes, oldValues)

	newTimes, newValues := GetManagerTimesAndValues(msg.ManagerTimeline)
	newTimelineFirstMatches := GetFirstMatchOnlyForTimeline(newTimes, newValues)

	updatedTimelineTimes := GetDetailsToCheck(oldTimelineFirstMatches, newTimelineFirstMatches, func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails {
		detailsToCheck := []*types.UniversalPermissionDetails{}
		if oldValue.(string) != newValue.(string) {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{})
		}
		return detailsToCheck
	})

	if err := CheckTimedUpdatePermission(ctx, updatedTimelineTimes, collection.Permissions.CanUpdateManager); err != nil {
		return nil, err
	}

	collection.ManagerTimeline = msg.ManagerTimeline


	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgUpdateManagerResponse{}, nil
}
