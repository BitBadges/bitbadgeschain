package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetFullIdRanges() []*types.IdRange {
	return []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:  sdk.NewUint(math.MaxUint64),
		},
	}
}

func GetBottomHalfIdRanges() []*types.IdRange {
	return []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:  sdk.NewUint(math.MaxUint32),
		},
	}
}

func GetTopHalfIdRanges() []*types.IdRange {
	return []*types.IdRange{
		{
			Start: sdk.NewUint(math.MaxUint32 + 1),
			End:  sdk.NewUint(math.MaxUint64),
		},
	}
}

func GetOneIdRange() []*types.IdRange {
	return []*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:  sdk.NewUint(1),
		},
	}
}

func GetCollectionsToCreate() []CollectionsToCreate {
	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				Creator: bob,
				ApprovedTransfersTimeline: []*types.CollectionApprovedTransferTimeline{
					{
						Times: GetFullIdRanges(),
						ApprovedTransfers: []*types.CollectionApprovedTransfer{
						{
							ToMappingId: "All",
							FromMappingId: "All",
							InitiatedByMappingId: "All",
							TransferTimes: GetFullIdRanges(),
							BadgeIds: GetFullIdRanges(),
							AllowedCombinations: []*types.IsCollectionTransferAllowed{
								{
									IsAllowed: true,
								},
							},
							Challenges: []*types.Challenge{},
							TrackerId: "test",
							IncrementIdsBy: sdk.NewUint(0),
							IncrementTimesBy: sdk.NewUint(0),
							PerAddressApprovals: &types.PerAddressApprovals{
								ApprovalsPerFromAddress: &types.ApprovalsTracker{
									Amounts: []*types.Balance{
										{
											Amount: sdk.NewUint(1),
											Times: GetFullIdRanges(),
											BadgeIds: GetFullIdRanges(),
										},
									},
									NumTransfers: sdk.NewUint(1000),
								},
							},
						}},
					},
				},
				DefaultApprovedIncomingTransfersTimeline: []*types.UserApprovedIncomingTransferTimeline{
					{
						ApprovedIncomingTransfers: []*types.UserApprovedIncomingTransfer{
							{
								FromMappingId: "All",
								InitiatedByMappingId: "All",
								TransferTimes: GetFullIdRanges(),
								BadgeIds: GetFullIdRanges(),
								AllowedCombinations: []*types.IsUserIncomingTransferAllowed{
									{
										IsAllowed: true,
									},
								},
								Challenges: []*types.Challenge{},
								TrackerId: "test",
								IncrementIdsBy: sdk.NewUint(0),
								IncrementTimesBy: sdk.NewUint(0),
								PerAddressApprovals: &types.PerAddressApprovals{
									ApprovalsPerFromAddress: &types.ApprovalsTracker{
										Amounts: []*types.Balance{
											{
												Amount: sdk.NewUint(1),
												Times: GetFullIdRanges(),
												BadgeIds: GetFullIdRanges(),
											},
										},
										NumTransfers: sdk.NewUint(1000),
									},
								},
							},
						},
						Times: GetFullIdRanges(),
					},
				},
				DefaultApprovedOutgoingTransfersTimeline: []*types.UserApprovedOutgoingTransferTimeline{
					{
						ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
							{
								ToMappingId: "All",
								InitiatedByMappingId: "All",
								TransferTimes: GetFullIdRanges(),
								BadgeIds: GetFullIdRanges(),
								AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
									{
										IsAllowed: true,
									},
								},
								Challenges: []*types.Challenge{},
								TrackerId: "test",
								IncrementIdsBy: sdk.NewUint(0),
								IncrementTimesBy: sdk.NewUint(0),
								PerAddressApprovals: &types.PerAddressApprovals{
									ApprovalsPerFromAddress: &types.ApprovalsTracker{
										Amounts: []*types.Balance{
											{
												Amount: sdk.NewUint(1),
												Times: GetFullIdRanges(),
												BadgeIds: GetFullIdRanges(),
											},
										},
										NumTransfers: sdk.NewUint(1000),
									},
								},
							},
						},
						Times: GetFullIdRanges(),
					},
				},
				BadgesToCreate: []*types.Balance{
					{
						Amount: sdk.NewUint(1),
						BadgeIds: GetFullIdRanges(),
						Times: GetFullIdRanges(),
					},
				},
				Permissions: &types.CollectionPermissions{
					CanCreateMoreBadges: []*types.ActionWithBadgeIdsPermission{
						{
							DefaultValues: &types.ActionWithBadgeIdsDefaultValues{
								BadgeIds: GetFullIdRanges(),
								PermittedTimes: GetFullIdRanges(),
								ForbiddenTimes: []*types.IdRange{},
							},
							Combinations: []*types.ActionWithBadgeIdsCombination{{}},
						},
					},
				},
	
			},
			Amount:  sdk.NewUint(1),
		
		},
	}

	return collectionsToCreate
}

func (suite *TestSuite) TestDeductFromOutgoing() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	
	CreateCollections(suite, wctx, GetCollectionsToCreate())
	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)


	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection, &bobBalance,	GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdk.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection, &bobBalance,	GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdk.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection, &bobBalance,	GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdk.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection, &bobBalance,	GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdk.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection, GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdk.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection, GetFullIdRanges(), GetFullIdRanges(), bob, alice, alice, sdk.NewUint(1), []*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}


func (suite *TestSuite) TestDeductFromOutgoingTwoSeparateTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	CreateCollections(suite, wctx, GetCollectionsToCreate())
	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)
	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)

	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx,	collection,	&bobBalance,	GetTopHalfIdRanges(),	GetFullIdRanges(), bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx,	collection,	&bobBalance,	GetTopHalfIdRanges(),	GetFullIdRanges(), bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&aliceBalance, GetTopHalfIdRanges(), GetFullIdRanges(), alice, bob, alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")


	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx,	collection,	&bobBalance,	GetTopHalfIdRanges(),	GetFullIdRanges(), bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx,	collection,	&bobBalance,	GetTopHalfIdRanges(),	GetFullIdRanges(), bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&aliceBalance, GetTopHalfIdRanges(), GetFullIdRanges(), alice, bob, alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestMaxOneTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].PerAddressApprovals.ApprovalsPerFromAddress.NumTransfers = sdk.NewUint(1)
	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)


	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestClaimIncrementsExceedsBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].IncrementIdsBy = sdk.NewUint(math.MaxUint64)
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].IncrementTimesBy = sdk.NewUint(math.MaxUint64)
	collectionsToCreate[0].Collection.DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].IncrementIdsBy = sdk.NewUint(math.MaxUint64)
	collectionsToCreate[0].Collection.DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].IncrementTimesBy = sdk.NewUint(math.MaxUint64)
	
	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestRequiresEquals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].RequireToDoesNotEqualInitiatedBy = true
	collectionsToCreate[0].Collection.DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].RequireFromDoesNotEqualInitiatedBy = true	

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	bob,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	charlie,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	
	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	charlie,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestSpecificApproved() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].InitiatedByMappingId = alice
	collectionsToCreate[0].Collection.DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].InitiatedByMappingId = alice	

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	charlie,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	charlie,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

}

func (suite *TestSuite) TestDefaults() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].InitiatedByMappingId = alice
	collectionsToCreate[0].Collection.DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].InitiatedByMappingId = alice	

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)


	err := suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	bob,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
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
							Start: sdk.NewUint(1),
							End: sdk.NewUint(1),
						},
					},
					AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
						{
							IsAllowed: false,
						},
					},
					Challenges: []*types.Challenge{},
					TrackerId: "test",
					IncrementIdsBy: sdk.NewUint(0),
					IncrementTimesBy: sdk.NewUint(0),
					PerAddressApprovals: &types.PerAddressApprovals{
						ApprovalsPerFromAddress: &types.ApprovalsTracker{
							Amounts: []*types.Balance{
								{
									Amount: sdk.NewUint(1),
									Times: GetFullIdRanges(),
									BadgeIds: GetFullIdRanges(),
								},
							},
							NumTransfers: sdk.NewUint(1000),
						},
					},
				},
				collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0],
			},
			Times: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetBottomHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
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
					TransferTimes:  []*types.IdRange{ { Start: sdk.NewUint(1), End: sdk.NewUint(1) } },
					BadgeIds: GetFullIdRanges(),
					AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
						{
							IsAllowed: false,
						},
					},
					Challenges: []*types.Challenge{},
					TrackerId: "test",
					IncrementIdsBy: sdk.NewUint(0),
					IncrementTimesBy: sdk.NewUint(0),
					PerAddressApprovals: &types.PerAddressApprovals{
						ApprovalsPerFromAddress: &types.ApprovalsTracker{
							Amounts: []*types.Balance{
								{
									Amount: sdk.NewUint(1),
									Times: GetFullIdRanges(),
									BadgeIds: GetFullIdRanges(),
								},
							},
							NumTransfers: sdk.NewUint(1000),
						},
					},
				},
				collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0],
			},
			Times: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
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
					BadgeIds:[]*types.IdRange{ { Start: sdk.NewUint(1), End: sdk.NewUint(1) } },
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
					IncrementIdsBy: sdk.NewUint(0),
					IncrementTimesBy: sdk.NewUint(0),
					PerAddressApprovals: &types.PerAddressApprovals{
						ApprovalsPerFromAddress: &types.ApprovalsTracker{
							Amounts: []*types.Balance{
								{
									Amount: sdk.NewUint(1),
									Times: GetFullIdRanges(),
									BadgeIds: GetFullIdRanges(),
								},
							},
							NumTransfers: sdk.NewUint(1000),
						},
					},
				},
			},
			Times: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetOneIdRange(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
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
					BadgeIds:[]*types.IdRange{ { Start: sdk.NewUint(1), End: sdk.NewUint(1) } },
					AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
						{
							IsAllowed: false,
						},
						{
							IsAllowed: true,
						},
					},
					Challenges: []*types.Challenge{},
					TrackerId: "test",
					IncrementIdsBy: sdk.NewUint(0),
					IncrementTimesBy: sdk.NewUint(0),
					PerAddressApprovals: &types.PerAddressApprovals{
						ApprovalsPerFromAddress: &types.ApprovalsTracker{
							Amounts: []*types.Balance{
								{
									Amount: sdk.NewUint(1),
									Times: GetFullIdRanges(),
									BadgeIds: GetFullIdRanges(),
								},
							},
							NumTransfers: sdk.NewUint(1000),
						},
					},
				},
			},
			Times: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetOneIdRange(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
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
					BadgeIds:[]*types.IdRange{ { Start: sdk.NewUint(1), End: sdk.NewUint(1) } },
					AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
						{
							IsAllowed: true,
						},
					},
					Challenges: []*types.Challenge{},
					TrackerId: "test",
					IncrementIdsBy: sdk.NewUint(0),
					IncrementTimesBy: sdk.NewUint(0),
					PerAddressApprovals: &types.PerAddressApprovals{
						ApprovalsPerFromAddress: &types.ApprovalsTracker{
							Amounts: []*types.Balance{
								{
									Amount: sdk.NewUint(1),
									Times: GetFullIdRanges(),
									BadgeIds: GetFullIdRanges(),
								},
							},
							NumTransfers: sdk.NewUint(1000),
						},
					},
				},
			},
			Times: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err := suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, collection,	&bobBalance, GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturned() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()
	// collectionsToCreate[0].Collection.ApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesFromApprovedOutgoingTransfers = true

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection,	GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(2, len(x), "Error deducting outgoing approvals")
	suite.Require().True(x[0].Outgoing != x[1].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesOutgoing() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.ApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesFromApprovedOutgoingTransfers = true

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection,	GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(1, len(x), "Error deducting outgoing approvals")
	suite.Require().False(x[0].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesIncoming() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.ApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesToApprovedIncomingTransfers = true

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection,	GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(1, len(x), "Error deducting outgoing approvals")
	suite.Require().True(x[0].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesBoth() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.ApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesToApprovedIncomingTransfers = true
	collectionsToCreate[0].Collection.ApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesFromApprovedOutgoingTransfers = true

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection,	GetTopHalfIdRanges(), GetFullIdRanges(),	bob,	alice,	alice,	sdk.NewUint(1),	[]*types.ChallengeSolution{},)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(0, len(x), "Error deducting outgoing approvals")
}

//TODO: Test transfer tracker ID after update approved transfers