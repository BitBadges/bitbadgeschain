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

	toCheck := []string{"Mint", "Manager", "!dadsf", "1asdsdfa:1234", "All", "None", alice, bob, charlie} //"122:323",
	for _, check := range toCheck {
		err := suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
			MappingId: check,
		})	
		suite.Require().Error(err, "Error creating address mapping: %s", check)
	}

	autoFetched := []string{"Mint", "Manager", "All", "None", alice, bob, charlie} //"122:323",
	for _, check := range autoFetched {
		err := suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
			MappingId: check,
		})	
		suite.Require().Error(err, "Error creating address mapping: %s", check)

		mapping, err := suite.app.BadgesKeeper.GetAddressMapping(suite.ctx, check, alice)
		suite.Require().Nil(err, "Error getting address mapping: %s", check)
		suite.Require().NotNil(mapping, "Error getting address mapping: %s", check)
	}

	found, err := suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "Mint", "Mint", alice, []string{})
	suite.Require().True(found, "Error checking mapping addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Mint")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "Mint", alice, alice, []string{})
	suite.Require().False(found, "Error checking mapping addresses: %s", "Mint")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Mint")


	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "Manager", alice, alice, []string{})
	suite.Require().True(found, "Error checking mapping addresses: %s", "Manager")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Manager")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "Manager", "Mint", alice, []string{})
	suite.Require().False(found, "Error checking mapping addresses: %s", "Manager")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Manager")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "Manager", "Manager", alice, []string{})
	suite.Require().True(found, "Error checking mapping addresses: %s", "Manager")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "Manager")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "All", "Manager", alice, []string{})
	suite.Require().True(found, "Error checking mapping addresses: %s", "All")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "All")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "All", "Mint", alice, []string{})
	suite.Require().False(found, "Error checking mapping addresses: %s", "All")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "All")

	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "None", alice, alice, []string{})
	suite.Require().False(found, "Error checking mapping addresses: %s", "None")
	suite.Require().Nil(err, "Error checking mapping addresses: %s", "None")

	// mapping, err := suite.app.BadgesKeeper.GetAddressMapping(suite.ctx, "1:1", alice)
	// suite.Require().Nil(err, "Error getting address mapping: %s", "1:1")
	// suite.Require().Equal(mapping.MappingId, "1:1", "Error getting address mapping: %s", "1:1")
	// AssertUintsEqual(suite, mapping.Filters[0].Conditions[0].MustOwnBadges[0].BadgeIds[0].Start, sdkmath.NewUint(1))
	// AssertUintsEqual(suite, mapping.Filters[0].Conditions[0].MustOwnBadges[0].BadgeIds[0].End, sdkmath.NewUint(1))
	// AssertUintsEqual(suite, mapping.Filters[0].Conditions[0].MustOwnBadges[0].CollectionId, sdkmath.NewUint(1))
}

// func (suite *TestSuite) TestAddressMappingsManagerOf() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test",
// 		Addresses: []string{},
// 		IncludeOnlySpecified: false,
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

// 	found, err := suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", bob, bob, []string{})
// 	suite.Require().True(found, "Error checking address mapping manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address mapping manager of: %s", "test")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", alice, bob, []string{})
// 	suite.Require().False(found, "Error checking address mapping manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address mapping manager of: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test",
// 		Addresses: []string{},
// 		IncludeOnlySpecified: false,
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

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", bob, bob, []string{})
// 	suite.Require().False(found, "Error checking address mapping manager of: %s", "test")
// 	suite.Require().Nil(err, "Error checking address mapping manager of: %s", "test")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", alice, bob, []string{})
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
// 		IncludeOnlySpecified: false,
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
// 		IncludeOnlySpecified: false,
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
// 		IncludeOnlySpecified: false,
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


// 	_, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", bob, bob, []string{})
// 	suite.Require().Error(err, "Error checking address mapping circular lookups: %s", "test")

// 	_, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test2", bob, bob, []string{})
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
// 		IncludeOnlySpecified: false,
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
// 		IncludeOnlySpecified: false,
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
// 		IncludeOnlySpecified: false,
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


// 	_, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test", bob, bob, []string{})
// 	suite.Require().Error(err, "Error checking address mapping circular lookups: %s", "test")

// 	_, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test2", bob, bob, []string{})
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
// 		IncludeOnlySpecified: true,
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test3",
// 		Addresses: []string{},
// 		IncludeOnlySpecified: false,
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
// 		IncludeOnlySpecified: false,
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

// 	found, err := suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", alice, bob, []string{})
// 	suite.Require().Nil(err, "Error checking alice which should be in address mapping")
// 	suite.Require().True(found, "Error checking address mapping")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", bob, bob, []string{})
// 	suite.Require().Nil(err, "Error checking bob which should not be in address mapping")
// 	suite.Require().False(found, "Error checking address mapping")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "mustnot", alice, bob, []string{})
// 	suite.Require().Nil(err, "Error checking alice which should be in address mapping")
// 	suite.Require().False(found, "Error checking address mapping")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "mustnot", bob, bob, []string{})
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
// 		IncludeOnlySpecified: true,
// 	})
// 	suite.Require().Nil(err, "Error creating address mapping: %s", "test")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test3",
// 		Addresses: []string{},
// 		IncludeOnlySpecified: false,
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

// 	found, err := suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", alice, bob, []string{})
// 	suite.Require().Nil(err, "Error checking alice which should be in address mapping")
// 	suite.Require().False(found, "Error checking address mapping")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", bob, bob, []string{})
// 	suite.Require().Nil(err, "Error checking bob which should not be in address mapping")
// 	suite.Require().False(found, "Error checking address mapping")

// 	err = suite.app.BadgesKeeper.CreateAddressMapping(suite.ctx, &types.AddressMapping{
// 		MappingId: "test3",
// 		Addresses: []string{},
// 		IncludeOnlySpecified: false,
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

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", alice, bob, []string{})
// 	suite.Require().Nil(err, "Error checking alice which should be in address mapping")
// 	suite.Require().True(found, "Error checking address mapping")

// 	found, err = suite.app.BadgesKeeper.CheckMappingAddresses(suite.ctx, "test3", bob, bob, []string{})
// 	suite.Require().Nil(err, "Error checking bob which should not be in address mapping")
// 	suite.Require().False(found, "Error checking address mapping")
// }