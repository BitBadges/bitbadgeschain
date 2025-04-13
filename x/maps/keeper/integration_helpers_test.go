package keeper_test

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/maps/types"
)

/* Query helpers */

func GetMap(suite *TestSuite, ctx context.Context, id string) (*types.Map, error) {
	res, err := suite.app.MapsKeeper.Map(ctx, &types.QueryGetMapRequest{MapId: id})
	if err != nil {
		return &types.Map{}, err
	}

	return res.Map, nil
}

func GetMapValue(suite *TestSuite, ctx context.Context, id string, key string) (*types.ValueStore, error) {
	res, err := suite.app.MapsKeeper.MapValue(ctx, &types.QueryGetMapValueRequest{MapId: id, Key: key})
	if err != nil {
		return &types.ValueStore{}, err
	}

	return res.Value, nil
}

// /* Msg helpers */

func CreateMap(suite *TestSuite, ctx context.Context, msg *types.MsgCreateMap) error {
	_, err := suite.msgServer.CreateMap(ctx, msg)
	return err
}

func SetValue(suite *TestSuite, ctx context.Context, msg *types.MsgSetValue) error {
	_, err := suite.msgServer.SetValue(ctx, msg)
	return err
}

func DeleteMap(suite *TestSuite, ctx context.Context, msg *types.MsgDeleteMap) error {
	_, err := suite.msgServer.DeleteMap(ctx, msg)
	return err
}

func UpdateMap(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateMap) error {
	_, err := suite.msgServer.UpdateMap(ctx, msg)
	return err
}
