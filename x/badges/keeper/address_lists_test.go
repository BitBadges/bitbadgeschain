package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestReservedIds() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	toCheck := []string{"Mint", "!dadsf", "1asdsdfa:1234", "All", "AllWithoutMint", "!(Mint)", "AllWithoutMint:" + alice, "None", alice, bob, charlie} //"122:323",
	for _, check := range toCheck {
		err := suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
			ListId: check,
		})
		suite.Require().Error(err, "Error creating address list: %s", check)
	}

	autoFetched := []string{"Mint", "AllWithoutMint", "!(Mint)" , "None", alice, bob, charlie} //"122:323",
	for _, check := range autoFetched {
		err := suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
			ListId: check,
		})
		suite.Require().Error(err, "Error creating address list: %s", check)

		list, err := suite.app.BadgesKeeper.GetAddressListById(suite.ctx, check)
		suite.Require().Nil(err, "Error getting address list: %s", check)
		suite.Require().NotNil(list, "Error getting address list: %s", check)
	}

	found, err := suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "Mint", "Mint")
	suite.Require().True(found, "Error checking list addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking list addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "Mint", alice)
	suite.Require().False(found, "Error checking list addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking list addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "AllWithoutMint:"+alice, alice)
	suite.Require().False(found, "Error checking list addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking list addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "AllWithoutMint:"+alice, bob)
	suite.Require().True(found, "Error checking list addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking list addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "AllWithoutMint:"+alice, "Mint")
	suite.Require().False(found, "Error checking list addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking list addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, bob+":"+alice, alice)
	suite.Require().True(found, "Error checking list addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking list addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, bob+":"+alice, bob)
	suite.Require().True(found, "Error checking list addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking list addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, bob+":"+alice, "Mint")
	suite.Require().False(found, "Error checking list addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking list addresses: %s", "Mint")

	// found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "Manager", alice,)
	// suite.Require().True(found, "Error checking list addresses: %s", "Manager")
	// suite.Require().Nil(err, "Error checking list addresses: %s", "Manager")

	// found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "Manager", "Mint",)
	// suite.Require().False(found, "Error checking list addresses: %s", "Manager")
	// suite.Require().Nil(err, "Error checking list addresses: %s", "Manager")

	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "AllWithoutMint", "Mint")
	suite.Require().False(found, "Error checking list addresses: %s", "AllWithoutMint")
	suite.Require().Nil(err, "Error checking list addresses: %s", "AllWithoutMint")

	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "None", alice)
	suite.Require().False(found, "Error checking list addresses: %s", "None")
	suite.Require().Nil(err, "Error checking list addresses: %s", "None")

	// list, err := suite.app.BadgesKeeper.GetAddressListById(suite.ctx, "1:1",)
	// suite.Require().Nil(err, "Error getting address list: %s", "1:1")
	// suite.Require().Equal(list.ListId, "1:1", "Error getting address list: %s", "1:1")
	// AssertUintsEqual(suite, list.Filters[0].Conditions[0].MustOwnBadges[0].BadgeIds[0].Start, sdkmath.NewUint(1))
	// AssertUintsEqual(suite, list.Filters[0].Conditions[0].MustOwnBadges[0].BadgeIds[0].End, sdkmath.NewUint(1))
	// AssertUintsEqual(suite, list.Filters[0].Conditions[0].MustOwnBadges[0].CollectionId, sdkmath.NewUint(1))
}

func (suite *TestSuite) TestStoreAddressLists() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].AddressLists = []*types.AddressList{
		{
			ListId: "test1asdasfda",
			Addresses: []string{alice},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	list, err := GetAddressList(suite, suite.ctx, "test1asdasfda")
	suite.Require().Nil(err, "Error getting address list: %s", "test1asdasfda")
	suite.Require().NotNil(list, "Error getting address list: %s", "test1asdasfda")
	suite.Require().Equal(list.ListId, "test1asdasfda", "Error getting address list: %s", "test1asdasfda")
}

func (suite *TestSuite) TestDuplicateStoreAddressLists() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].AddressLists = []*types.AddressList{
		{
			ListId: "test1asdasfda",
			Addresses: []string{alice},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = CreateAddressLists(suite, wctx, &types.MsgCreateAddressLists{
		Creator:         alice,
		AddressLists: collectionsToCreate[0].AddressLists,
	})
	suite.Require().Error(err, "Error creating badge: %s")
}

// func (suite *TestSuite) TestAddressListsManagerOf() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustBeManager: []sdkmath.Uint{sdkmath.NewUint(1)},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	found, err := suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test", bob, bob)
// 	suite.Require().True(found, "Error checking address list manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address list manager of: %s", "test")

// 	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test", alice, bob)
// 	suite.Require().False(found, "Error checking address list manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address list manager of: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustNotBeManager: []sdkmath.Uint{sdkmath.NewUint(1)},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test", bob, bob)
// 	suite.Require().False(found, "Error checking address list manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address list manager of: %s", "test")

// 	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test", alice, bob)
// 	suite.Require().True(found, "Error checking address list manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address list manager of: %s", "test")
// }

// func (suite *TestSuite) TestAddressListCircularLookups() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustBeInList: []string{"test"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test2",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustBeInList: []string{"test3"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test3",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustBeInList: []string{"test2"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	_, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test", bob, bob)
// 	suite.Require().Error(err, "Error checking address list circular lookups: %s", "test")

// 	_, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test2", bob, bob)
// 	suite.Require().Error(err, "Error checking address list circular lookups: %s", "test2")
// }

// func (suite *TestSuite) TestAddressListCircularLookupsInverted() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustNotBeInList: []string{"test"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test2",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustNotBeInList: []string{"test3"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test3",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustNotBeInList: []string{"test2"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	_, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test", bob, bob)
// 	suite.Require().Error(err, "Error checking address list circular lookups: %s", "test")

// 	_, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test2", bob, bob)
// 	suite.Require().Error(err, "Error checking address list circular lookups: %s", "test2")
// }

// func (suite *TestSuite) TestAddressListMustBeInAnotherList() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test",
// 		Addresses: []string{alice},
// 		Allowlist: true,
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test3",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustBeInList: []string{"test"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "mustnot",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustNotBeInList: []string{"test"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	found, err := suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test3", alice, bob)
// 	suite.Require().Nil(err, "Error checking alice which should be in address list")
// 	suite.Require().True(found, "Error checking address list")

// 	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test3", bob, bob)
// 	suite.Require().Nil(err, "Error checking bob which should not be in address list")
// 	suite.Require().False(found, "Error checking address list")

// 	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "mustnot", alice, bob)
// 	suite.Require().Nil(err, "Error checking alice which should be in address list")
// 	suite.Require().False(found, "Error checking address list")

// 	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "mustnot", bob, bob)
// 	suite.Require().Nil(err, "Error checking bob which should not be in address list")
// 	suite.Require().True(found, "Error checking address list")
// }

// func (suite *TestSuite) TestAddressListMoreThanOneCondition() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test",
// 		Addresses: []string{alice},
// 		Allowlist: true,
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test3",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustBeInList: []string{"test"},
// 						MustBeManager: []sdkmath.Uint{sdkmath.NewUint(1)},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	found, err := suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test3", alice, bob)
// 	suite.Require().Nil(err, "Error checking alice which should be in address list")
// 	suite.Require().False(found, "Error checking address list")

// 	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test3", bob, bob)
// 	suite.Require().Nil(err, "Error checking bob which should not be in address list")
// 	suite.Require().False(found, "Error checking address list")

// 	err = suite.app.BadgesKeeper.CreateAddressList(suite.ctx, &types.AddressList{
// 		ListId: "test3",
// 		Addresses: []string{},
// 		Allowlist: false,
// 		Filters: []*types.AddressListFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressListConditions{
// 					{
// 						MustBeInList: []string{"test"},
// 						MustBeManager: []sdkmath.Uint{sdkmath.NewUint(1)},
// 					},
// 					{
// 						MustBeInList: []string{"test"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address list: %s", "test")

// 	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test3", alice, bob)
// 	suite.Require().Nil(err, "Error checking alice which should be in address list")
// 	suite.Require().True(found, "Error checking address list")

// 	found, err = suite.app.BadgesKeeper.CheckAddresses(suite.ctx, "test3", bob, bob)
// 	suite.Require().Nil(err, "Error checking bob which should not be in address list")
// 	suite.Require().False(found, "Error checking address list")
// }
