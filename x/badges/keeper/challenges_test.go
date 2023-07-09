package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNoChallengesWorking() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesToApprovedIncomingTransfers = true
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesFromApprovedOutgoingTransfers = true
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].Challenges = []*types.Challenge{{}}

	CreateCollections(suite, wctx, collectionsToCreate)
	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	_, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.ChallengeSolution{}, false)
	suite.Require().Error(err, "Error getting user balance: %s")
}

func (suite *TestSuite) TestChallengesInvalidSolutions() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesToApprovedIncomingTransfers = true
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesFromApprovedOutgoingTransfers = true
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].Challenges = []*types.Challenge{{
		Root: "sample",
	}}

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	_, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.ChallengeSolution{}, false)
	suite.Require().Error(err, "Error getting user balance: %s")

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.ChallengeSolution{
		{
			Proof: &types.ClaimProof{
				Aunts: []*types.ClaimProofItem{},
				Leaf:  "sample",
			},
		},
	}, false)
	suite.Require().Error(err, "Error getting user balance: %s")
}
func (suite *TestSuite) TestSendAllToClaimsAccountTypeInvalid() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	aliceLeaf := "-" + alice + "-1-0-0"
	bobLeaf := "-" + bob + "-1-0-0"
	charlieLeaf := "-" + charlie + "-1-0-0"

	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
	leafHashes := make([][]byte, len(leafs))
	for i, leaf := range leafs {
		initialHash := sha256.Sum256(leaf)
		leafHashes[i] = initialHash[:]
	}

	levelTwoHashes := make([][]byte, 2)
	for i := 0; i < len(leafHashes); i += 2 {
		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
		levelTwoHashes[i/2] = iHash[:]
	}

	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
	rootHash := rootHashI[:]

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers, &types.CollectionApprovedTransfer{
		OverallApprovals: &types.ApprovalsTracker{
			Amounts: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10),
					BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			NumTransfers: sdk.NewUint(10),
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		Uri: "",
		Challenges: []*types.Challenge{
			{
				Root:                hex.EncodeToString(rootHash),
				ExpectedProofLength: sdk.NewUint(2),
				MaxOneUsePerLeaf:    true,
			},
		},
		ApprovalId:            "testing232",
		TransferTimes:        GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		OwnershipTimes: GetFullUintRanges(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

		OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
		IncrementBadgeIdsBy:                    sdk.NewUint(0),
		IncrementOwnershipTimesBy:              sdk.NewUint(0),
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: "",
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[1]),
									OnRight: true,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: aliceLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[1]),
									OnRight: true,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: aliceLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[1]),
									OnRight: true,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestFailsOnUseCreatorAddressAsLeaf() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	aliceLeaf := "-" + alice + "-1-0-0"
	bobLeaf := "-" + bob + "-1-0-0"
	charlieLeaf := "-" + charlie + "-1-0-0"

	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
	leafHashes := make([][]byte, len(leafs))
	for i, leaf := range leafs {
		initialHash := sha256.Sum256(leaf)
		leafHashes[i] = initialHash[:]
	}

	levelTwoHashes := make([][]byte, 2)
	for i := 0; i < len(leafHashes); i += 2 {
		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
		levelTwoHashes[i/2] = iHash[:]
	}

	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
	rootHash := rootHashI[:]

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers, &types.CollectionApprovedTransfer{
		OverallApprovals: &types.ApprovalsTracker{
			Amounts: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10),
					BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			NumTransfers: sdk.NewUint(10),
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		Uri: "",
		Challenges: []*types.Challenge{
			{
				Root:                    hex.EncodeToString(rootHash),
				ExpectedProofLength:     sdk.NewUint(2),
				MaxOneUsePerLeaf:        true,
				UseCreatorAddressAsLeaf: true,
			},
		},
		ApprovalId:            "testing232",
		TransferTimes:        GetFullUintRanges(),
		OwnershipTimes: GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

		OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
		IncrementBadgeIdsBy:                    sdk.NewUint(0),
		IncrementOwnershipTimesBy:              sdk.NewUint(0),
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: aliceLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[1]),
									OnRight: true,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestWrongExpectedProofLength() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	aliceLeaf := "-" + alice + "-1-0-0"
	bobLeaf := "-" + bob + "-1-0-0"
	charlieLeaf := "-" + charlie + "-1-0-0"

	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
	leafHashes := make([][]byte, len(leafs))
	for i, leaf := range leafs {
		initialHash := sha256.Sum256(leaf)
		leafHashes[i] = initialHash[:]
	}

	levelTwoHashes := make([][]byte, 2)
	for i := 0; i < len(leafHashes); i += 2 {
		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
		levelTwoHashes[i/2] = iHash[:]
	}

	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
	rootHash := rootHashI[:]

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers, &types.CollectionApprovedTransfer{
		OverallApprovals: &types.ApprovalsTracker{
			Amounts: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10),
					BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			NumTransfers: sdk.NewUint(10),
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		Uri: "",
		Challenges: []*types.Challenge{
			{
				Root:                hex.EncodeToString(rootHash),
				ExpectedProofLength: sdk.NewUint(5),
				MaxOneUsePerLeaf:    true,
			},
		},
		ApprovalId:            "testing232",
		TransferTimes:        GetFullUintRanges(),
		OwnershipTimes: GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

		OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
		IncrementBadgeIdsBy:                    sdk.NewUint(0),
		IncrementOwnershipTimesBy:              sdk.NewUint(0),
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf:  aliceLeaf,
							Aunts: []*types.ClaimProofItem{},
						},
					},
					{
						Proof: &types.ClaimProof{
							Leaf:  aliceLeaf,
							Aunts: []*types.ClaimProofItem{},
						},
					},
					{
						Proof: &types.ClaimProof{
							Leaf: aliceLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[1]),
									OnRight: true,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestIncrements() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	aliceLeaf := "-" + alice + "-1-0-0"
	bobLeaf := "-" + bob + "-1-0-0"
	charlieLeaf := "-" + charlie + "-1-0-0"

	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
	leafHashes := make([][]byte, len(leafs))
	for i, leaf := range leafs {
		initialHash := sha256.Sum256(leaf)
		leafHashes[i] = initialHash[:]
	}

	levelTwoHashes := make([][]byte, 2)
	for i := 0; i < len(leafHashes); i += 2 {
		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
		levelTwoHashes[i/2] = iHash[:]
	}

	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
	rootHash := rootHashI[:]

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers, &types.CollectionApprovedTransfer{
		OverallApprovals: &types.ApprovalsTracker{
			Amounts: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10),
					BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			NumTransfers: sdk.NewUint(10),
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		Uri: "",
		Challenges: []*types.Challenge{
			{
				Root:                             hex.EncodeToString(rootHash),
				ExpectedProofLength:              sdk.NewUint(2),
				MaxOneUsePerLeaf:                 true,
				UseLeafIndexForDistributionOrder: true,
			},
		},
		ApprovalId:            "testing232",
		TransferTimes:        GetFullUintRanges(),
		OwnershipTimes: 			GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

		OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
		IncrementBadgeIdsBy:                    sdk.NewUint(1),
		IncrementOwnershipTimesBy:              sdk.NewUint(0),
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: bobLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[0]),
									OnRight: false,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: aliceLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[1]),
									OnRight: true,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestIncrementsTransferAsMuchAsPossible() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	aliceLeaf := "-" + alice + "-1-0-0"
	bobLeaf := "-" + bob + "-1-0-0"
	charlieLeaf := "-" + charlie + "-1-0-0"

	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
	leafHashes := make([][]byte, len(leafs))
	for i, leaf := range leafs {
		initialHash := sha256.Sum256(leaf)
		leafHashes[i] = initialHash[:]
	}

	levelTwoHashes := make([][]byte, 2)
	for i := 0; i < len(leafHashes); i += 2 {
		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
		levelTwoHashes[i/2] = iHash[:]
	}

	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
	rootHash := rootHashI[:]

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers, &types.CollectionApprovedTransfer{
		OverallApprovals: &types.ApprovalsTracker{
			Amounts: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10),
					BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			NumTransfers: sdk.NewUint(10),
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		Uri: "",
		Challenges: []*types.Challenge{
			{
				Root:                             hex.EncodeToString(rootHash),
				ExpectedProofLength:              sdk.NewUint(2),
				MaxOneUsePerLeaf:                 true,
				UseLeafIndexForDistributionOrder: true,
			},
		},
		ApprovalId:            "testing232",
		TransferTimes:        GetFullUintRanges(),
		OwnershipTimes: GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

		OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
		IncrementBadgeIdsBy:                    sdk.NewUint(1),
		IncrementOwnershipTimesBy:              sdk.NewUint(0),
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				TransferAsMuchAsPossible: true,
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: aliceLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[1]),
									OnRight: true,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetOneUintRange(),
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)
}


func (suite *TestSuite) TestIncrementsTransferAsMuchAsPossibleGreaterAmount() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	aliceLeaf := "-" + alice + "-1-0-0"
	bobLeaf := "-" + bob + "-1-0-0"
	charlieLeaf := "-" + charlie + "-1-0-0"

	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
	leafHashes := make([][]byte, len(leafs))
	for i, leaf := range leafs {
		initialHash := sha256.Sum256(leaf)
		leafHashes[i] = initialHash[:]
	}

	levelTwoHashes := make([][]byte, 2)
	for i := 0; i < len(leafHashes); i += 2 {
		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
		levelTwoHashes[i/2] = iHash[:]
	}

	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
	rootHash := rootHashI[:]

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers, &types.CollectionApprovedTransfer{
		OverallApprovals: &types.ApprovalsTracker{
			Amounts: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			NumTransfers: sdk.NewUint(10),
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		Uri: "",
		Challenges: []*types.Challenge{
			{
				Root:                             hex.EncodeToString(rootHash),
				ExpectedProofLength:              sdk.NewUint(2),
				MaxOneUsePerLeaf:                 true,
				UseLeafIndexForDistributionOrder: true,
			},
		},
		ApprovalId:            "testing232",
		TransferTimes:        GetFullUintRanges(),
		OwnershipTimes:       GetFullUintRanges(),
		BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

		OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
		IncrementBadgeIdsBy:                    sdk.NewUint(1),
		IncrementOwnershipTimesBy:              sdk.NewUint(0),
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				TransferAsMuchAsPossible: true,
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: aliceLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[1]),
									OnRight: true,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetOneUintRange(),
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)


	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				TransferAsMuchAsPossible: true,
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: bobLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[0]),
									OnRight: false,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetTwoUintRanges(),
		OwnershipTimes: GetFullUintRanges(),
	}}, aliceBalance.Balances)
}



func (suite *TestSuite) TestIncrementsTransferAsMuchAsPossibleGreaterAmountSolo() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	aliceLeaf := "-" + alice + "-1-0-0"
	bobLeaf := "-" + bob + "-1-0-0"
	charlieLeaf := "-" + charlie + "-1-0-0"

	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
	leafHashes := make([][]byte, len(leafs))
	for i, leaf := range leafs {
		initialHash := sha256.Sum256(leaf)
		leafHashes[i] = initialHash[:]
	}

	levelTwoHashes := make([][]byte, 2)
	for i := 0; i < len(leafHashes); i += 2 {
		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
		levelTwoHashes[i/2] = iHash[:]
	}

	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
	rootHash := rootHashI[:]

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers, &types.CollectionApprovedTransfer{
		OverallApprovals: &types.ApprovalsTracker{
			Amounts: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			NumTransfers: sdk.NewUint(10),
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		Uri: "",
		Challenges: []*types.Challenge{
			{
				Root:                             hex.EncodeToString(rootHash),
				ExpectedProofLength:              sdk.NewUint(2),
				MaxOneUsePerLeaf:                 true,
				UseLeafIndexForDistributionOrder: true,
			},
		},
		ApprovalId:            "testing232",
		TransferTimes:        GetFullUintRanges(),
		OwnershipTimes:       GetFullUintRanges(),
		BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

		OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
		IncrementBadgeIdsBy:                    sdk.NewUint(1),
		IncrementOwnershipTimesBy:              sdk.NewUint(0),
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				TransferAsMuchAsPossible: true,
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: bobLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[0]),
									OnRight: false,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetTwoUintRanges(),
		OwnershipTimes: GetFullUintRanges(),
	}}, aliceBalance.Balances)
}


func (suite *TestSuite) TestIncrementsTransferGreaterThanMaxNumTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	aliceLeaf := "-" + alice + "-1-0-0"
	bobLeaf := "-" + bob + "-1-0-0"
	charlieLeaf := "-" + charlie + "-1-0-0"

	leafs := [][]byte{[]byte(aliceLeaf), []byte(bobLeaf), []byte(charlieLeaf), []byte(charlieLeaf)}
	leafHashes := make([][]byte, len(leafs))
	for i, leaf := range leafs {
		initialHash := sha256.Sum256(leaf)
		leafHashes[i] = initialHash[:]
	}

	levelTwoHashes := make([][]byte, 2)
	for i := 0; i < len(leafHashes); i += 2 {
		iHash := sha256.Sum256(append(leafHashes[i], leafHashes[i+1]...))
		levelTwoHashes[i/2] = iHash[:]
	}

	rootHashI := sha256.Sum256(append(levelTwoHashes[0], levelTwoHashes[1]...))
	rootHash := rootHashI[:]

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].ApprovedTransfers, &types.CollectionApprovedTransfer{
		OverallApprovals: &types.ApprovalsTracker{
			Amounts: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			NumTransfers: sdk.NewUint(1),
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		Uri: "",
		Challenges: []*types.Challenge{
			{
				Root:                             hex.EncodeToString(rootHash),
				ExpectedProofLength:              sdk.NewUint(2),
				MaxOneUsePerLeaf:                 true,
				UseLeafIndexForDistributionOrder: true,
			},
		},
		ApprovalId:            "testing232",
		TransferTimes:        GetFullUintRanges(),
		OwnershipTimes:       GetFullUintRanges(),
		BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

		OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
		IncrementBadgeIdsBy:                    sdk.NewUint(1),
		IncrementOwnershipTimesBy:              sdk.NewUint(0),
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				TransferAsMuchAsPossible: true,
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: aliceLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[1]),
									OnRight: true,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetOneUintRange(),
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)


	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				TransferAsMuchAsPossible: true,
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				Solutions: []*types.ChallengeSolution{
					{
						Proof: &types.ClaimProof{
							Leaf: bobLeaf,
							Aunts: []*types.ClaimProofItem{
								{
									Aunt:    hex.EncodeToString(leafHashes[0]),
									OnRight: false,
								},
								{
									Aunt:    hex.EncodeToString(levelTwoHashes[1]),
									OnRight: true,
								},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{}, aliceBalance.Balances)
}