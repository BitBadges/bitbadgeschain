package keeper_test

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
)

/* Query helpers */

func GetManagerSplitter(suite *TestSuite, ctx context.Context, address string) (*types.ManagerSplitter, error) {
	res, err := suite.queryClient.ManagerSplitter(ctx, &types.QueryGetManagerSplitterRequest{Address: address})
	if err != nil {
		return &types.ManagerSplitter{}, err
	}

	return res.ManagerSplitter, nil
}

/* Msg helpers */

func CreateManagerSplitter(suite *TestSuite, ctx context.Context, msg *types.MsgCreateManagerSplitter) (*types.MsgCreateManagerSplitterResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	res, err := suite.msgServer.CreateManagerSplitter(ctx, msg)
	return res, err
}

func UpdateManagerSplitter(suite *TestSuite, ctx context.Context, msg *types.MsgUpdateManagerSplitter) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.UpdateManagerSplitter(ctx, msg)
	return err
}

func DeleteManagerSplitter(suite *TestSuite, ctx context.Context, msg *types.MsgDeleteManagerSplitter) error {
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	_, err = suite.msgServer.DeleteManagerSplitter(ctx, msg)
	return err
}

func ExecuteUniversalUpdateCollection(suite *TestSuite, ctx context.Context, msg *types.MsgExecuteUniversalUpdateCollection) (*types.MsgExecuteUniversalUpdateCollectionResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	res, err := suite.msgServer.ExecuteUniversalUpdateCollection(ctx, msg)
	return res, err
}
