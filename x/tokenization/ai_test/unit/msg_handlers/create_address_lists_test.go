package msg_handlers_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type CreateAddressListsTestSuite struct {
	testutil.AITestSuite
}

func TestCreateAddressListsSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(CreateAddressListsTestSuite))
}

// TestCreateAddressLists_SingleList tests successfully creating a single address list
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_SingleList() {
	msg := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "testList1",
				Addresses: []string{suite.Alice, suite.Bob},
				Whitelist: true,
				Uri:       "https://example.com/list",
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "creating address list should succeed")

	// Verify list was persisted
	list, found := suite.Keeper.GetAddressListFromStore(suite.Ctx, "testList1")
	suite.Require().True(found, "list should exist")
	suite.Require().Equal("testList1", list.ListId)
	suite.Require().Equal(2, len(list.Addresses))
	suite.Require().True(list.Whitelist)
	suite.Require().Equal(suite.Manager, list.CreatedBy)
}

// TestCreateAddressLists_MultipleLists tests successfully creating multiple lists at once
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_MultipleLists() {
	msg := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{ListId: "listA", Addresses: []string{suite.Alice}, Whitelist: true},
			{ListId: "listB", Addresses: []string{suite.Bob}, Whitelist: false},
			{ListId: "listC", Addresses: []string{suite.Charlie}, Whitelist: true},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Verify all lists exist
	listA, found := suite.Keeper.GetAddressListFromStore(suite.Ctx, "listA")
	suite.Require().True(found)
	suite.Require().True(listA.Whitelist)
	suite.Require().Equal(suite.Manager, listA.CreatedBy)

	listB, found := suite.Keeper.GetAddressListFromStore(suite.Ctx, "listB")
	suite.Require().True(found)
	suite.Require().False(listB.Whitelist)
	suite.Require().Equal(suite.Manager, listB.CreatedBy)

	listC, found := suite.Keeper.GetAddressListFromStore(suite.Ctx, "listC")
	suite.Require().True(found)
	suite.Require().True(listC.Whitelist)
	suite.Require().Equal(suite.Manager, listC.CreatedBy)
}

// TestCreateAddressLists_ListPersistence tests that list is correctly persisted to state
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_ListPersistence() {
	addresses := []string{suite.Alice, suite.Bob, suite.Charlie}
	msg := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:     "persistenceTest",
				Addresses:  addresses,
				Whitelist:  false,
				Uri:        "https://example.com/persistence",
				CustomData: "test custom data",
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Verify all fields are persisted correctly
	list, found := suite.Keeper.GetAddressListFromStore(suite.Ctx, "persistenceTest")
	suite.Require().True(found, "list should exist in store")
	suite.Require().Equal("persistenceTest", list.ListId)
	suite.Require().Equal(3, len(list.Addresses))
	suite.Require().Equal(suite.Alice, list.Addresses[0])
	suite.Require().Equal(suite.Bob, list.Addresses[1])
	suite.Require().Equal(suite.Charlie, list.Addresses[2])
	suite.Require().False(list.Whitelist)
	suite.Require().Equal("https://example.com/persistence", list.Uri)
	suite.Require().Equal("test custom data", list.CustomData)
	suite.Require().Equal(suite.Manager, list.CreatedBy)
}

// TestCreateAddressLists_CreatedByAutoSet tests that createdBy is automatically set to creator
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_CreatedByAutoSet() {
	// Test with Alice as creator
	msg := &types.MsgCreateAddressLists{
		Creator: suite.Alice,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "aliceList",
				Addresses: []string{suite.Bob},
				Whitelist: true,
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	list, found := suite.Keeper.GetAddressListFromStore(suite.Ctx, "aliceList")
	suite.Require().True(found)
	suite.Require().Equal(suite.Alice, list.CreatedBy, "createdBy should be auto-set to creator")

	// Test with Bob as creator
	msg2 := &types.MsgCreateAddressLists{
		Creator: suite.Bob,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "bobList",
				Addresses: []string{suite.Alice},
				Whitelist: true,
			},
		},
	}

	_, err = suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	list2, found := suite.Keeper.GetAddressListFromStore(suite.Ctx, "bobList")
	suite.Require().True(found)
	suite.Require().Equal(suite.Bob, list2.CreatedBy, "createdBy should be auto-set to creator")
}

// TestCreateAddressLists_RetrievalAfterCreation tests list retrieval after creation
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_RetrievalAfterCreation() {
	msg := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "retrievalTest",
				Addresses: []string{suite.Alice, suite.Bob},
				Whitelist: true,
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Test retrieval using GetAddressListById (which handles reserved IDs and inversion)
	list, err := suite.Keeper.GetAddressListById(suite.Ctx, "retrievalTest")
	suite.Require().NoError(err, "should be able to retrieve list by ID")
	suite.Require().Equal("retrievalTest", list.ListId)
	suite.Require().Equal(2, len(list.Addresses))
	suite.Require().True(list.Whitelist)
}

// TestCreateAddressLists_EmptyAddressList tests creating a list with no addresses
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_EmptyAddressList() {
	msg := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "emptyList",
				Addresses: []string{},
				Whitelist: true,
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "creating empty address list should succeed")

	list, found := suite.Keeper.GetAddressListFromStore(suite.Ctx, "emptyList")
	suite.Require().True(found)
	suite.Require().Equal(0, len(list.Addresses))
}

// TestCreateAddressLists_DuplicateListId tests that duplicate list IDs fail
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_DuplicateListId() {
	// Create first list
	msg1 := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "duplicateTest",
				Addresses: []string{suite.Alice},
				Whitelist: true,
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Try to create another list with the same ID
	msg2 := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "duplicateTest",
				Addresses: []string{suite.Bob},
				Whitelist: false,
			},
		},
	}

	_, err = suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "creating list with duplicate ID should fail")
}

// TestCreateAddressLists_InvalidListId tests that invalid list IDs fail
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_InvalidListId() {
	// Test with list ID starting with !
	msg := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "!invalidList",
				Addresses: []string{suite.Alice},
				Whitelist: true,
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "list ID starting with ! should fail")
}

// TestCreateAddressLists_InvalidAddress tests that invalid addresses fail
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_InvalidAddress() {
	msg := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "invalidAddressTest",
				Addresses: []string{"invalid_address"},
				Whitelist: true,
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "creating list with invalid address should fail")
}

// TestCreateAddressLists_WhitelistVsBlacklist tests whitelist and blacklist behavior
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_WhitelistVsBlacklist() {
	// Create whitelist
	msg1 := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "whitelistTest",
				Addresses: []string{suite.Alice},
				Whitelist: true,
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Create blacklist
	msg2 := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "blacklistTest",
				Addresses: []string{suite.Alice},
				Whitelist: false,
			},
		},
	}

	_, err = suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Verify whitelist
	whitelist, found := suite.Keeper.GetAddressListFromStore(suite.Ctx, "whitelistTest")
	suite.Require().True(found)
	suite.Require().True(whitelist.Whitelist, "whitelist flag should be true")

	// Verify blacklist
	blacklist, found := suite.Keeper.GetAddressListFromStore(suite.Ctx, "blacklistTest")
	suite.Require().True(found)
	suite.Require().False(blacklist.Whitelist, "whitelist flag should be false for blacklist")
}

// TestCreateAddressLists_CheckAddresses tests checking if an address is in a list
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_CheckAddresses() {
	// Create a whitelist containing Alice
	msg := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "checkAddressTest",
				Addresses: []string{suite.Alice},
				Whitelist: true,
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Check Alice is in the list (whitelist)
	isAliceInList, err := suite.Keeper.CheckAddresses(suite.Ctx, "checkAddressTest", suite.Alice)
	suite.Require().NoError(err)
	suite.Require().True(isAliceInList, "Alice should be in the whitelist")

	// Check Bob is NOT in the list (whitelist)
	isBobInList, err := suite.Keeper.CheckAddresses(suite.Ctx, "checkAddressTest", suite.Bob)
	suite.Require().NoError(err)
	suite.Require().False(isBobInList, "Bob should NOT be in the whitelist")
}

// TestCreateAddressLists_InvalidCreator tests that invalid creator fails ValidateBasic
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_InvalidCreator() {
	msg := &types.MsgCreateAddressLists{
		Creator: "invalid_creator",
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "testList",
				Addresses: []string{suite.Alice},
				Whitelist: true,
			},
		},
	}

	// ValidateBasic should fail for invalid creator
	err := msg.ValidateBasic()
	suite.Require().Error(err, "invalid creator should fail ValidateBasic")
}

// TestCreateAddressLists_DuplicateAddressesInList tests that duplicate addresses in a list fail
func (suite *CreateAddressListsTestSuite) TestCreateAddressLists_DuplicateAddressesInList() {
	msg := &types.MsgCreateAddressLists{
		Creator: suite.Manager,
		AddressLists: []*types.AddressListInput{
			{
				ListId:    "duplicateAddressTest",
				Addresses: []string{suite.Alice, suite.Alice}, // Duplicate Alice
				Whitelist: true,
			},
		},
	}

	_, err := suite.MsgServer.CreateAddressLists(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "duplicate addresses in list should fail")
}
