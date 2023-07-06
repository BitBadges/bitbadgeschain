package keeper_test

import (
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestDeductFromOutgoing() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, GetCollectionsToCreate())
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)


	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection, bobBalance,	GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection, bobBalance,	GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection, bobBalance,	GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection, bobBalance,	GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection, GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection, GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}


func (suite *TestSuite) TestDeductFromOutgoingTwoSeparateTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, GetCollectionsToCreate())
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)
	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx,	collection,	bobBalance,	GetTopHalfIdRanges(),	GetFullIdRanges(), bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx,	collection,	bobBalance,	GetTopHalfIdRanges(),	GetFullIdRanges(), bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	aliceBalance, GetTopHalfIdRanges(), GetFullIdRanges(), alice, bob, alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")


	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx,	collection,	bobBalance,	GetTopHalfIdRanges(),	GetFullIdRanges(), bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx,	collection,	bobBalance,	GetTopHalfIdRanges(),	GetFullIdRanges(), bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	aliceBalance, GetTopHalfIdRanges(), GetFullIdRanges(), alice, bob, alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestMaxOneTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].PerAddressApprovals.ApprovalsPerFromAddress.NumTransfers = sdkmath.NewUint(1)
	err := CreateCollections(suite, wctx, collectionsToCreate)
suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)


	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestClaimIncrementsExceedsBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].IncrementBadgeIdsBy = sdkmath.NewUint(math.MaxUint64)
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].IncrementOwnershipTimesBy = sdkmath.NewUint(math.MaxUint64)
	collectionsToCreate[0].Collection.DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].IncrementBadgeIdsBy = sdkmath.NewUint(math.MaxUint64)
	collectionsToCreate[0].Collection.DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].IncrementOwnershipTimesBy = sdkmath.NewUint(math.MaxUint64)
	
	err := CreateCollections(suite, wctx, collectionsToCreate)
suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestRequiresEquals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].RequireToDoesNotEqualInitiatedBy = true
	collectionsToCreate[0].Collection.DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].RequireFromDoesNotEqualInitiatedBy = true	

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	bob,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	charlie,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	
	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	charlie,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestSpecificApproved() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].InitiatedByMappingId = alice
	collectionsToCreate[0].Collection.DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].InitiatedByMappingId = alice	

	err := CreateCollections(suite, wctx, collectionsToCreate)
suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	charlie,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	charlie,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

}

func (suite *TestSuite) TestDefaults() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].InitiatedByMappingId = alice
	collectionsToCreate[0].Collection.DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].InitiatedByMappingId = alice	

	err := CreateCollections(suite, wctx, collectionsToCreate)
suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)


	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	bob,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestFirstMatchOnly() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].InitiatedByMappingId = alice
	newOutgoingTimeline := []*types.UserApprovedOutgoingTransferTimeline{	{
			ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
				{
					ToMappingId: "All",
					InitiatedByMappingId: alice,
					TransferTimes: GetFullIdRanges(),
					BadgeIds: []*types.IdRange{
						{
							Start: sdkmath.NewUint(1),
							End: sdkmath.NewUint(1),
						},
					},
					AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
						{
							IsAllowed: false,
						},
					},
					Challenges: []*types.Challenge{},
					TrackerId: "test",
					IncrementBadgeIdsBy: sdkmath.NewUint(0),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					PerAddressApprovals: &types.PerAddressApprovals{
						ApprovalsPerFromAddress: &types.ApprovalsTracker{
							Amounts: []*types.Balance{
								{
									Amount: sdkmath.NewUint(1),
									OwnershipTimes: GetFullIdRanges(),
									BadgeIds: GetFullIdRanges(),
								},
							},
							NumTransfers: sdkmath.NewUint(1000),
						},
					},
				},
				collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0],
			},
			TimelineTimes: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}


func (suite *TestSuite) TestFirstMatchOnlyWrongTime() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()

	newOutgoingTimeline := []*types.UserApprovedOutgoingTransferTimeline{	{
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
				{
					ToMappingId: "All",
					InitiatedByMappingId: alice,
					TransferTimes:  []*types.IdRange{ { Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1) } },
					BadgeIds: GetFullIdRanges(),
					AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
						{
							IsAllowed: false,
						},
					},
					Challenges: []*types.Challenge{},
					TrackerId: "test",
					IncrementBadgeIdsBy: sdkmath.NewUint(0),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					PerAddressApprovals: &types.PerAddressApprovals{
						ApprovalsPerFromAddress: &types.ApprovalsTracker{
							Amounts: []*types.Balance{
								{
									Amount: sdkmath.NewUint(1),
									OwnershipTimes: GetFullIdRanges(),
									BadgeIds: GetFullIdRanges(),
								},
							},
							NumTransfers: sdkmath.NewUint(1000),
						},
					},
				},
				collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0],
			},
			TimelineTimes: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestCombinations() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()

	newOutgoingTimeline := []*types.UserApprovedOutgoingTransferTimeline{	{
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
				{
					ToMappingId: "All",
					InitiatedByMappingId: alice,
					TransferTimes:  GetFullIdRanges(),
					BadgeIds:[]*types.IdRange{ { Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1) } },
					AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
						{
							IsAllowed: true,
						},
						{
							InvertBadgeIds: true,
							IsAllowed: false,
						},
					},
					Challenges: []*types.Challenge{},
					TrackerId: "test",
					IncrementBadgeIdsBy: sdkmath.NewUint(0),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					PerAddressApprovals: &types.PerAddressApprovals{
						ApprovalsPerFromAddress: &types.ApprovalsTracker{
							Amounts: []*types.Balance{
								{
									Amount: sdkmath.NewUint(1),
									OwnershipTimes: GetFullIdRanges(),
									BadgeIds: GetFullIdRanges(),
								},
							},
							NumTransfers: sdkmath.NewUint(1000),
						},
					},
				},
			},
			TimelineTimes: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetOneIdRange(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}


func (suite *TestSuite) TestCombinationsOrder() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()

	newOutgoingTimeline := []*types.UserApprovedOutgoingTransferTimeline{	{
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
				{
					ToMappingId: "All",
					InitiatedByMappingId: alice,
					TransferTimes:  GetFullIdRanges(),
					BadgeIds:[]*types.IdRange{ { Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1) } },
					AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
						{
							IsAllowed: false,
						},
						{
							InvertInitiatedBy: true,
							IsAllowed: true,
						},
					},
					Challenges: []*types.Challenge{},
					TrackerId: "test",
					IncrementBadgeIdsBy: sdkmath.NewUint(0),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					PerAddressApprovals: &types.PerAddressApprovals{
						ApprovalsPerFromAddress: &types.ApprovalsTracker{
							Amounts: []*types.Balance{
								{
									Amount: sdkmath.NewUint(1),
									OwnershipTimes: GetFullIdRanges(),
									BadgeIds: GetFullIdRanges(),
								},
							},
							NumTransfers: sdkmath.NewUint(1000),
						},
					},
				},
			},
			TimelineTimes: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetOneIdRange(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestNotExplicitlyDefined() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()

	newOutgoingTimeline := []*types.UserApprovedOutgoingTransferTimeline{	{
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
				{
					ToMappingId: "All",
					InitiatedByMappingId: alice,
					TransferTimes:  GetFullIdRanges(),
					BadgeIds:[]*types.IdRange{ { Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1) } },
					AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
						{
							IsAllowed: true,
						},
					},
					Challenges: []*types.Challenge{},
					TrackerId: "test",
					IncrementBadgeIdsBy: sdkmath.NewUint(0),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					PerAddressApprovals: &types.PerAddressApprovals{
						ApprovalsPerFromAddress: &types.ApprovalsTracker{
							Amounts: []*types.Balance{
								{
									Amount: sdkmath.NewUint(1),
									OwnershipTimes: GetFullIdRanges(),
									BadgeIds: GetFullIdRanges(),
								},
							},
							NumTransfers: sdkmath.NewUint(1000),
						},
					},
				},
			},
			TimelineTimes: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	bobBalance, GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturned() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()
	// collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesFromApprovedOutgoingTransfers = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection,	GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(2, len(x), "Error deducting outgoing approvals")
	suite.Require().True(x[0].Outgoing != x[1].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesOutgoing() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesFromApprovedOutgoingTransfers = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection,	GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(1, len(x), "Error deducting outgoing approvals")
	suite.Require().False(x[0].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesIncoming() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesToApprovedIncomingTransfers = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection,	GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(1, len(x), "Error deducting outgoing approvals")
	suite.Require().True(x[0].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesBoth() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesToApprovedIncomingTransfers = true
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesFromApprovedOutgoingTransfers = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
suite.Require().Nil(err, "error creating badges")


	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection,	GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdkmath.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(0, len(x), "Error deducting outgoing approvals")
}

//TODO: Test transfer tracker ID after update approved transfers