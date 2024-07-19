package keeper_test

import (
	"math"

	"bitbadgeschain/x/maps/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	badgeskeeper "bitbadgeschain/x/badges/keeper"
	badgestypes "bitbadgeschain/x/badges/types"
)

func GetDefaultCreateMsg() *types.MsgCreateMap {
	return &types.MsgCreateMap{
		Creator: bob,
		MapId:   "test",
		UpdateCriteria: &types.MapUpdateCriteria{
			ManagerOnly:         true,
			CollectionId:        sdkmath.NewUintFromString("0"),
			CreatorOnly:         false,
			FirstComeFirstServe: false,
		},
		ValueOptions: &types.ValueOptions{
			PermanentOnceSet: false,
			NoDuplicates:     false,
		},
		DefaultValue: "",
		ManagerTimeline: []*types.ManagerTimeline{
			{
				TimelineTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
				Manager: bob,
			},
		},
		MetadataTimeline: []*types.MapMetadataTimeline{
			{
				Metadata: &types.Metadata{
					Uri:        "test",
					CustomData: "",
				},
				TimelineTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
			},
		},
		Permissions: &types.MapPermissions{
			CanUpdateMetadata: []*types.TimedUpdatePermission{
				{
					PermanentlyPermittedTimes: []*types.UintRange{},
					PermanentlyForbiddenTimes: []*types.UintRange{},
					TimelineTimes: []*types.UintRange{
						{
							Start: sdkmath.NewUint(1),
							End:   sdkmath.NewUint(math.MaxUint64),
						},
					},
				},
			},
			CanUpdateManager: []*types.TimedUpdatePermission{
				{
					PermanentlyPermittedTimes: []*types.UintRange{},
					PermanentlyForbiddenTimes: []*types.UintRange{},
					TimelineTimes: []*types.UintRange{
						{
							Start: sdkmath.NewUint(1),
							End:   sdkmath.NewUint(math.MaxUint64),
						},
					},
				},
			},
			CanDeleteMap: []*types.ActionPermission{
				{
					PermanentlyPermittedTimes: []*types.UintRange{},
					PermanentlyForbiddenTimes: []*types.UintRange{},
				},
			},
		},
	}
}

func (suite *TestSuite) TestMaps() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateMap(suite, wctx, GetDefaultCreateMsg())
	suite.Require().Nil(err, "Error creating map: %s")

	currMap, err := GetMap(suite, wctx, "test")
	suite.Require().Nil(err, "Error getting map: %s")
	suite.Require().NotNil(currMap, "Error getting map: %s")
	suite.Require().Equal("test", currMap.MapId, "Error getting map: %s")

	suite.Require().Equal(bob, currMap.ManagerTimeline[0].Manager, "Error getting map: %s")

	updateMapMsg := &types.MsgUpdateMap{
		Creator:               bob,
		MapId:                 "test",
		UpdateManagerTimeline: true,
		ManagerTimeline: []*types.ManagerTimeline{
			{
				Manager: alice,
				TimelineTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
			},
		},
		UpdateMetadataTimeline: false,
		MetadataTimeline:       nil,
		UpdatePermissions:      false,
		Permissions:            nil,
	}

	updateMapMsg.Creator = alice
	err = UpdateMap(suite, wctx, updateMapMsg)
	suite.Require().Error(err, "Error updating map: %s")
	updateMapMsg.Creator = bob

	err = UpdateMap(suite, wctx, updateMapMsg)
	suite.Require().Nil(err, "Error updating map: %s")
	updateMapMsg.Creator = alice

	err = DeleteMap(suite, wctx, &types.MsgDeleteMap{
		Creator: bob,
		MapId:   "test",
	})
	suite.Require().Error(err, "Error deleting map: %s")

	err = DeleteMap(suite, wctx, &types.MsgDeleteMap{
		Creator: alice,
		MapId:   "test",
	})
	suite.Require().Nil(err, "Error deleting map: %s")
}

func (suite *TestSuite) TestPermissions() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.Permissions.CanUpdateManager = []*types.TimedUpdatePermission{
		{
			PermanentlyPermittedTimes: []*types.UintRange{},
			PermanentlyForbiddenTimes: []*types.UintRange{{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(math.MaxUint64),
			}},
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
			},
		},
	}

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")

	updateMapMsg := &types.MsgUpdateMap{
		Creator:               bob,
		MapId:                 "test",
		UpdateManagerTimeline: true,
		ManagerTimeline: []*types.ManagerTimeline{
			{
				Manager: alice,
				TimelineTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
			},
		},
	}

	err = UpdateMap(suite, wctx, updateMapMsg)
	suite.Require().Error(err, "Error updating map: %s")
}

func (suite *TestSuite) TestDeleteIsDisallowed() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.Permissions.CanDeleteMap = []*types.ActionPermission{
		{
			PermanentlyPermittedTimes: []*types.UintRange{},
			PermanentlyForbiddenTimes: []*types.UintRange{{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(math.MaxUint64),
			}},
		},
	}

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")

	err = DeleteMap(suite, wctx, &types.MsgDeleteMap{
		Creator: bob,
		MapId:   "test",
	})
	suite.Require().Error(err, "Error deleting map: %s")
}

func (suite *TestSuite) TestManagerOnly() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.ManagerOnly = true

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "test",
		Key:     alice,
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Error(err, "Error setting value: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: bob,
		MapId:   "test",
		Key:     alice,
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Nil(err, "Error setting value: %s")

	valStore, err := GetMapValue(suite, wctx, "test", alice)
	if err != nil {
		suite.Require().Error(err, "Error getting value: %s")
	}

	suite.Require().Equal("test", valStore.Value, "Error getting value: %s")
	suite.Require().Equal(valStore.LastSetBy, bob, "Error getting value: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "test",
		Key:     alice,
		Value:   "",
		Options: &types.SetOptions{},
	})
	suite.Require().Error(err, "Error setting value: %s")
}

func (suite *TestSuite) TestFirstComeFirstServe() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.FirstComeFirstServe = true

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "test",
		Key:     alice,
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Nil(err, "Error setting value: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: bob,
		MapId:   "test",
		Key:     alice,
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Error(err, "Error setting value: %s")

	valStore, err := GetMapValue(suite, wctx, "test", alice)
	if err != nil {
		suite.Require().Error(err, "Error getting value: %s")
	}

	suite.Require().Equal("test", valStore.Value, "Error getting value: %s")
	suite.Require().Equal(valStore.LastSetBy, alice, "Error getting value: %s")

	//unset value
	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "test",
		Key:     alice,
		Value:   "",
		Options: &types.SetOptions{},
	})
	suite.Require().Nil(err, "Error setting value: %s")

	valStore, err = GetMapValue(suite, wctx, "test", alice)
	if err != nil {
		suite.Require().Error(err, "Error getting value: %s")
	}

	suite.Require().Equal("", valStore.Value, "Error getting value: %s")
}

func (suite *TestSuite) TestNoDuplicates() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.CreatorOnly = true
	createMsg.UpdateCriteria.ManagerOnly = false
	createMsg.ValueOptions.NoDuplicates = true

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "test",
		Key:     alice,
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Nil(err, "Error setting value: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: bob,
		MapId:   "test",
		Key:     bob,
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Error(err, "Error setting value: %s")

	valStore, err := GetMapValue(suite, wctx, "test", alice)
	if err != nil {
		suite.Require().Error(err, "Error getting value: %s")
	}

	suite.Require().Equal("test", valStore.Value, "Error getting value: %s")
	suite.Require().Equal(valStore.LastSetBy, alice, "Error getting value: %s")

	//unset value
	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "test",
		Key:     alice,
		Value:   "",
		Options: &types.SetOptions{},
	})
	suite.Require().Nil(err, "Error setting value: %s")

	valStore, err = GetMapValue(suite, wctx, "test", alice)
	if err != nil {
		suite.Require().Error(err, "Error getting value: %s")
	}

	suite.Require().Equal("", valStore.Value, "Error getting value: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: bob,
		MapId:   "test",
		Key:     bob,
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Nil(err, "Error setting value: %s")
}

func (suite *TestSuite) TestPermanentOnceSet() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.ValueOptions.PermanentOnceSet = true
	createMsg.UpdateCriteria.CreatorOnly = true
	createMsg.UpdateCriteria.ManagerOnly = false

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "test",
		Key:     alice,
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Nil(err, "Error setting value: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "test",
		Key:     alice,
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Error(err, "Error setting value: %s")

	//Cant unset either
	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "test",
		Key:     alice,
		Value:   "",
		Options: &types.SetOptions{},
	})
	suite.Require().Error(err, "Error setting value: %s")
}

func (suite *TestSuite) TestUseMostRecentCollectionId() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.CollectionId = sdkmath.NewUintFromString("0")
	createMsg.UpdateCriteria.CreatorOnly = true
	createMsg.UpdateCriteria.ManagerOnly = false

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "test",
		Key:     alice,
		Value:   "test",
		Options: &types.SetOptions{
			UseMostRecentCollectionId: true,
		},
	})
	suite.Require().Nil(err, "Error setting value: %s")

	valStore, err := GetMapValue(suite, wctx, "test", alice)
	if err != nil {
		suite.Require().Error(err, "Error getting value: %s")
	}

	suite.Require().Equal("0", valStore.Value, "Error getting value: %s")
}

func (suite *TestSuite) TestReservedAddressMaps() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.CollectionId = sdkmath.NewUintFromString("0")
	createMsg.UpdateCriteria.CreatorOnly = true
	createMsg.UpdateCriteria.ManagerOnly = false
	createMsg.MapId = alice
	createMsg.Creator = alice

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")
}

func (suite *TestSuite) TestReservedAddressMapsNotOwnAddress() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.CollectionId = sdkmath.NewUintFromString("0")
	createMsg.UpdateCriteria.CreatorOnly = true
	createMsg.UpdateCriteria.ManagerOnly = false
	createMsg.MapId = alice
	createMsg.Creator = bob

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Error(err, "Error creating map: %s")
}

func (suite *TestSuite) TestCollectionIdReservedMaps() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.CreatorOnly = true
	createMsg.UpdateCriteria.ManagerOnly = false
	createMsg.MapId = "1234"
	createMsg.Creator = alice

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Error(err, "Error creating map: %s")
}

func (suite *TestSuite) TestCollectionIdReservedMapsWithManager() {
	suite.app.BadgesKeeper.SetCollectionInStore(suite.ctx, &badgestypes.BadgeCollection{
		CollectionId: sdkmath.NewUint(1),
		ManagerTimeline: []*badgestypes.ManagerTimeline{
			{
				Manager: alice,
				TimelineTimes: []*badgestypes.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
			},
		},
	})

	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.CreatorOnly = true
	createMsg.UpdateCriteria.ManagerOnly = false
	createMsg.MapId = "1"
	createMsg.Creator = alice

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")
}

func (suite *TestSuite) TestCollectionIdReservedMapsWithManagerNotManager() {
	suite.app.BadgesKeeper.SetCollectionInStore(suite.ctx, &badgestypes.BadgeCollection{
		CollectionId: sdkmath.NewUint(1),
		ManagerTimeline: []*badgestypes.ManagerTimeline{
			{
				Manager: alice,
				TimelineTimes: []*badgestypes.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
			},
		},
	})

	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.CreatorOnly = true
	createMsg.UpdateCriteria.ManagerOnly = false
	createMsg.MapId = "1"
	createMsg.Creator = bob

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Error(err, "Error creating map: %s")
}

func (suite *TestSuite) TestCollectionIdCriteria() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.CollectionId = sdkmath.NewUintFromString("1")
	createMsg.UpdateCriteria.CreatorOnly = true
	createMsg.UpdateCriteria.ManagerOnly = false
	createMsg.MapId = "1"
	createMsg.Creator = alice

	suite.app.BadgesKeeper.SetCollectionInStore(suite.ctx, &badgestypes.BadgeCollection{
		CollectionId: sdkmath.NewUint(1),
		ManagerTimeline: []*badgestypes.ManagerTimeline{
			{
				Manager: alice,
				TimelineTimes: []*badgestypes.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
			},
		},
		BalancesType: "Standard",
	})

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")

	suite.app.BadgesKeeper.SetUserBalanceInStore(suite.ctx, badgeskeeper.ConstructBalanceKey(alice, sdkmath.NewUint(1)), &badgestypes.UserBalanceStore{
		Balances: []*badgestypes.Balance{
			{
				BadgeIds: []*badgestypes.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(1),
					},
				},
				Amount: sdkmath.NewUint(1),
				OwnershipTimes: []*badgestypes.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
			},
		},
	})

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "1",
		Key:     "1",
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Nil(err, "Error setting value: %s")

	valStore, err := GetMapValue(suite, wctx, "1", "1")
	if err != nil {
		suite.Require().Error(err, "Error getting value: %s")
	}

	suite.Require().Equal("test", valStore.Value, "Error getting value: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "1",
		Key:     "2",
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Error(err, "Error setting value: %s")
}

func (suite *TestSuite) TestInheritManagerFromCollection() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	createMsg := GetDefaultCreateMsg()
	createMsg.UpdateCriteria.ManagerOnly = false
	createMsg.UpdateCriteria.CreatorOnly = false
	createMsg.UpdateCriteria.FirstComeFirstServe = true
	createMsg.MapId = "testing"
	createMsg.Creator = alice
	createMsg.ManagerTimeline = []*types.ManagerTimeline{}
	createMsg.InheritManagerTimelineFrom = sdkmath.NewUint(1)

	err := CreateMap(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating map: %s")

	suite.app.BadgesKeeper.SetCollectionInStore(suite.ctx, &badgestypes.BadgeCollection{
		CollectionId: sdkmath.NewUint(1),
		ManagerTimeline: []*badgestypes.ManagerTimeline{
			{
				Manager: alice,
				TimelineTimes: []*badgestypes.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
			},
		},
		BalancesType: "Standard",
	})

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: alice,
		MapId:   "testing",
		Key:     "test",
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Nil(err, "Error setting value: %s")

	err = SetValue(suite, wctx, &types.MsgSetValue{
		Creator: bob,
		MapId:   "testing",
		Key:     "test",
		Value:   "test",
		Options: &types.SetOptions{},
	})
	suite.Require().Error(err, "Error setting value: %s")
}
