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
		err := suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
			MappingId: check,
		})
		suite.Require().Error(err, "Error creating address mapping: %s", check)
	}

	autoFetched := []string{"Mint", "AllWithoutMint", "!(Mint)" , "None", alice, bob, charlie} //"122:323",
	for _, check := range autoFetched {
		err := suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
			MappingId: check,
		})
		suite.Require().Error(err, "Error creating address mapping: %s", check)

		mapping, err := suite.app.BadgesKeeper.GetAddressMappingById(suite.ctx, check)
		suite.Require().Nil(err, "Error getting address mapping: %s", check)
		suite.Require().NotNil(mapping, "Error getting address mapping: %s", check)
	}

	found, err := suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "Mint", "Mint")
	suite.Require().True(found, "Error checking mapping addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "Mint", alice)
	suite.Require().False(found, "Error checking mapping addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "AllWithoutMint:"+alice, alice)
	suite.Require().False(found, "Error checking mapping addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "AllWithoutMint:"+alice, bob)
	suite.Require().True(found, "Error checking mapping addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "AllWithoutMint:"+alice, "Mint")
	suite.Require().False(found, "Error checking mapping addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, bob+":"+alice, alice)
	suite.Require().True(found, "Error checking mapping addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, bob+":"+alice, bob)
	suite.Require().True(found, "Error checking mapping addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, bob+":"+alice, "Mint")
	suite.Require().False(found, "Error checking mapping addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Mint")

	// found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "Manager", alice,)
	// suite.Require().True(found, "Error checking mapping addresses: %s", "Manager")
	// suite.Require().Nil(err, "Error checking mapping addresses: %s", "Manager")

	// found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "Manager", "Mint",)
	// suite.Require().False(found, "Error checking mapping addresses: %s", "Manager")
	// suite.Require().Nil(err, "Error checking mapping addresses: %s", "Manager")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "AllWithoutMint", "Mint")
	suite.Require().False(found, "Error checking mapping addresses: %s", "AllWithoutMint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "AllWithoutMint")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "None", alice)
	suite.Require().False(found, "Error checking mapping addresses: %s", "None")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "None")

	// mapping, err := suite.app.BadgesKeeper.GetAddressMappingById(suite.ctx, "1:1",)
	// suite.Require().Nil(err, "Error getting address mapping: %s", "1:1")
	// suite.Require().Equal(mapping.MappingId, "1:1", "Error getting address mapping: %s", "1:1")
	// AssertUintsEqual(suite, mapping.Filters[0].Conditions[0].MustOwnBadges[0].BadgeIds[0].Start, sdkmath.NewUint(1))
	// AssertUintsEqual(suite, mapping.Filters[0].Conditions[0].MustOwnBadges[0].BadgeIds[0].End, sdkmath.NewUint(1))
	// AssertUintsEqual(suite, mapping.Filters[0].Conditions[0].MustOwnBadges[0].CollectionId, sdkmath.NewUint(1))
}

func (suite *TestSuite) TestStoreAddressMappings() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].AddressMappings = []*types.AddressMapping{
		{
			MappingId: "test1asdasfda",
			Addresses: []string{alice},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	mapping, err := GetAddressMapping(suite, suite.ctx, "test1asdasfda")
	suite.Require().Nil(err, "Error getting address mapping: %s", "test1asdasfda")
	suite.Require().NotNil(mapping, "Error getting address mapping: %s", "test1asdasfda")
	suite.Require().Equal(mapping.MappingId, "test1asdasfda", "Error getting address mapping: %s", "test1asdasfda")
}

func (suite *TestSuite) TestDuplicateStoreAddressMappings() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].AddressMappings = []*types.AddressMapping{
		{
			MappingId: "test1asdasfda",
			Addresses: []string{alice},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = CreateAddressMappings(suite, wctx, &types.MsgCreateAddressMappings{
		Creator:         alice,
		AddressMappings: collectionsToCreate[0].AddressMappings,
	})
	suite.Require().Error(err, "Error creating badge: %s")
}

// func (suite *TestSuite) TestAddressMappingsManagerOf() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustBeManager: []sdkmath.Uint{sdkmath.NewUint(1)},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	found, err := suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", bob, bob)
// 	suite.Require().True(found, "Error checking address mapping manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address mapping manager of: %s", "test")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", alice, bob)
// 	suite.Require().False(found, "Error checking address mapping manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address mapping manager of: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustNotBeManager: []sdkmath.Uint{sdkmath.NewUint(1)},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", bob, bob)
// 	suite.Require().False(found, "Error checking address mapping manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address mapping manager of: %s", "test")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", alice, bob)
// 	suite.Require().True(found, "Error checking address mapping manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address mapping manager of: %s", "test")
// }

// func (suite *TestSuite) TestAddressMappingCircularLookups() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustBeInMapping: []string{"test"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test2",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustBeInMapping: []string{"test3"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test3",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustBeInMapping: []string{"test2"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	_, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", bob, bob)
// 	suite.Require().Error(err, "Error checking address mapping circular lookups: %s", "test")

// 	_, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test2", bob, bob)
// 	suite.Require().Error(err, "Error checking address mapping circular lookups: %s", "test2")
// }

// func (suite *TestSuite) TestAddressMappingCircularLookupsInverted() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustNotBeInMapping: []string{"test"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test2",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustNotBeInMapping: []string{"test3"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test3",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustNotBeInMapping: []string{"test2"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	_, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", bob, bob)
// 	suite.Require().Error(err, "Error checking address mapping circular lookups: %s", "test")

// 	_, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test2", bob, bob)
// 	suite.Require().Error(err, "Error checking address mapping circular lookups: %s", "test2")
// }

// func (suite *TestSuite) TestAddressMappingMustBeInAnotherMapping() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test",
// 		Addresses: []string{alice},
// 		IncludeAddresses: true,
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test3",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustBeInMapping: []string{"test"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "mustnot",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustNotBeInMapping: []string{"test"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	found, err := suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", alice, bob)
// 	suite.Require().Nil(err, "Error checking alice which should be in address mapping")
// 	suite.Require().True(found, "Error checking address mapping")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", bob, bob)
// 	suite.Require().Nil(err, "Error checking bob which should not be in address mapping")
// 	suite.Require().False(found, "Error checking address mapping")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "mustnot", alice, bob)
// 	suite.Require().Nil(err, "Error checking alice which should be in address mapping")
// 	suite.Require().False(found, "Error checking address mapping")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "mustnot", bob, bob)
// 	suite.Require().Nil(err, "Error checking bob which should not be in address mapping")
// 	suite.Require().True(found, "Error checking address mapping")
// }

// func (suite *TestSuite) TestAddressMappingMoreThanOneCondition() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test",
// 		Addresses: []string{alice},
// 		IncludeAddresses: true,
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test3",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustBeInMapping: []string{"test"},
// 						MustBeManager: []sdkmath.Uint{sdkmath.NewUint(1)},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	found, err := suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", alice, bob)
// 	suite.Require().Nil(err, "Error checking alice which should be in address mapping")
// 	suite.Require().False(found, "Error checking address mapping")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", bob, bob)
// 	suite.Require().Nil(err, "Error checking bob which should not be in address mapping")
// 	suite.Require().False(found, "Error checking address mapping")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test3",
// 		Addresses: []string{},
// 		IncludeAddresses: false,
// 		Filters: []*types.AddressMappingFilters{
// 			{
// 				MustSatisfyMinX: sdkmath.NewUint(1),
// 				Conditions: []*types.AddressMappingConditions{
// 					{
// 						MustBeInMapping: []string{"test"},
// 						MustBeManager: []sdkmath.Uint{sdkmath.NewUint(1)},
// 					},
// 					{
// 						MustBeInMapping: []string{"test"},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", alice, bob)
// 	suite.Require().Nil(err, "Error checking alice which should be in address mapping")
// 	suite.Require().True(found, "Error checking address mapping")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", bob, bob)
// 	suite.Require().Nil(err, "Error checking bob which should not be in address mapping")
// 	suite.Require().False(found, "Error checking address mapping")
// }
