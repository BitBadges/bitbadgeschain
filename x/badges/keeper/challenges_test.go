package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNoMerkleChallengesWorking() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].OverridesToApprovedIncomingTransfers = true
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].OverridesFromApprovedOutgoingTransfers = true
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].MerkleChallenges = []*types.MerkleChallenge{{}}

	CreateCollections(suite, wctx, collectionsToCreate)
	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	_, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error getting user balance: %s")
}

func (suite *TestSuite) TestMerkleChallengesInvalidSolutions() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].OverridesToApprovedIncomingTransfers = true
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].OverridesFromApprovedOutgoingTransfers = true
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].MerkleChallenges = []*types.MerkleChallenge{{
		ChallengeId: "testchallenge",
		Root: "sample",
	}}

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	_, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error getting user balance: %s")

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{
		
			{
				Aunts: []*types.MerklePathItem{},
				Leaf:  "sample",
			},
		
	})
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
				
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                hex.EncodeToString(rootHash),
						ExpectedProofLength: sdk.NewUint(2),
						MaxOneUsePerLeaf:    true,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		OwnedTimes:       GetFullUintRanges(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

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
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: "",
							Aunts: []*types.MerklePathItem{
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
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                    hex.EncodeToString(rootHash),
						ExpectedProofLength:     sdk.NewUint(2),
						MaxOneUsePerLeaf:        true,
						UseCreatorAddressAsLeaf: true,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes: GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

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
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                hex.EncodeToString(rootHash),
						ExpectedProofLength: sdk.NewUint(5),
						MaxOneUsePerLeaf:    true,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes: GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

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
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf:  aliceLeaf,
							Aunts: []*types.MerklePathItem{},
						
					},
					{
						
							Leaf:  aliceLeaf,
							Aunts: []*types.MerklePathItem{},
						
					},
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
				PredeterminedBalances: &types.PredeterminedBalances{
					OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{ UseMerkleChallengeLeafIndex: true },
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds: GetOneUintRange(),
								Amount:   sdkmath.NewUint(1),
								OwnedTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy: sdkmath.NewUint(1),
						IncrementOwnedTimesBy: sdkmath.NewUint(0),
					},
				},
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                             hex.EncodeToString(rootHash),
						ExpectedProofLength:              sdk.NewUint(2),
						MaxOneUsePerLeaf:                 true,
						UseLeafIndexForTransferOrder: true,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},

		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes: 			GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

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
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: bobLeaf,
							Aunts: []*types.MerklePathItem{
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
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
					
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
			
				PredeterminedBalances: &types.PredeterminedBalances{
					OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{ UseOverallNumTransfers: true },
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds: GetOneUintRange(),
								Amount:   sdkmath.NewUint(1),
								OwnedTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy: sdkmath.NewUint(1),
						IncrementOwnedTimesBy: sdkmath.NewUint(0),
					},
				},
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                             hex.EncodeToString(rootHash),
						ExpectedProofLength:              sdk.NewUint(2),
						MaxOneUsePerLeaf:                 true,
						UseLeafIndexForTransferOrder: true,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes: 			GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateFromApproval: &types.ApprovalIdDetails{
					ApprovalId:   "testing232",
					Address: "",
					ApprovalLevel: "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetOneUintRange(),
		OwnedTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
				PredeterminedBalances: &types.PredeterminedBalances{
					OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{ UseOverallNumTransfers: true },
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds: GetOneUintRange(),
								Amount:   sdkmath.NewUint(1),
								OwnedTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy: sdkmath.NewUint(1),
						IncrementOwnedTimesBy: sdkmath.NewUint(0),
					},
				},
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                             hex.EncodeToString(rootHash),
						ExpectedProofLength:              sdk.NewUint(2),
						MaxOneUsePerLeaf:                 true,
						UseLeafIndexForTransferOrder: 		true,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes:       GetFullUintRanges(),
		BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateFromApproval: &types.ApprovalIdDetails{
					ApprovalId:   "testing232",
					Address: "",
					ApprovalLevel: "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetOneUintRange(),
		OwnedTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)


	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateFromApproval: &types.ApprovalIdDetails{
					ApprovalId:   "testing232",
					Address: "",
					ApprovalLevel: "collection",
				},
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: bobLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetTwoUintRanges(),
		OwnedTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(1),
				},
				PredeterminedBalances: &types.PredeterminedBalances{
					OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{ UseMerkleChallengeLeafIndex: true },
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds: GetOneUintRange(),
								Amount:   sdkmath.NewUint(1),
								OwnedTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy: sdkmath.NewUint(1),
						IncrementOwnedTimesBy: sdkmath.NewUint(0),
					},
				},
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                             hex.EncodeToString(rootHash),
						ExpectedProofLength:              sdk.NewUint(2),
						MaxOneUsePerLeaf:                 true,
						UseLeafIndexForTransferOrder: true,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes:       GetFullUintRanges(),
		BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateFromApproval: &types.ApprovalIdDetails{
					ApprovalId:   "testing232",
					Address: "",
					ApprovalLevel: "collection",
				},
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: bobLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetTwoUintRanges(),
		OwnedTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(1),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(1),
				},
				PredeterminedBalances: &types.PredeterminedBalances{
					OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{ UseMerkleChallengeLeafIndex: true },
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds: GetOneUintRange(),
								Amount:   sdkmath.NewUint(1),
								OwnedTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy: sdkmath.NewUint(1),
						IncrementOwnedTimesBy: sdkmath.NewUint(0),
					},
				},
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                             hex.EncodeToString(rootHash),
						ExpectedProofLength:              sdk.NewUint(2),
						MaxOneUsePerLeaf:                 true,
						UseLeafIndexForTransferOrder: true,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes:       GetFullUintRanges(),
		BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateFromApproval: &types.ApprovalIdDetails{
					ApprovalId:   "testing232",
					Address: "",
					ApprovalLevel: "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetOneUintRange(),
		OwnedTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)


	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateFromApproval: &types.ApprovalIdDetails{
					ApprovalId:   "testing232",
					Address: "",
					ApprovalLevel: "collection",
				},
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: bobLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestIncrementsTransferAsMuchAsPossibleOneTx() {
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(1),
				},
				PredeterminedBalances: &types.PredeterminedBalances{
					OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{ UseOverallNumTransfers: true },
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds: GetOneUintRange(),
								Amount:   sdkmath.NewUint(1),
								OwnedTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy: sdkmath.NewUint(1),
						IncrementOwnedTimesBy: sdkmath.NewUint(0),
					},
				},
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                             hex.EncodeToString(rootHash),
						ExpectedProofLength:              sdk.NewUint(2),
						MaxOneUsePerLeaf:                 true,
						UseLeafIndexForTransferOrder: false,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes:       GetFullUintRanges(),
		BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateFromApproval: &types.ApprovalIdDetails{
					ApprovalId:   "testing232",
					Address: "",
					ApprovalLevel: "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
					{
						
							Leaf: bobLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnedTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       []*types.UintRange{
			GetTwoUintRanges()[0],
		},
		OwnedTimes: GetFullUintRanges(),
	}}, aliceBalance.Balances)
}


func (suite *TestSuite) TestIncrementsUsingPerToAddressNumTransfers() {
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
	collectionsToCreate[0].BadgesToCreate = append(collectionsToCreate[0].BadgesToCreate, &types.Balance{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetFullUintRanges(),
		OwnedTimes: GetFullUintRanges(),
	})
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
				PredeterminedBalances: &types.PredeterminedBalances{
					OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{ UsePerToAddressNumTransfers: true },
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds: GetOneUintRange(),
								Amount:   sdkmath.NewUint(1),
								OwnedTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy: sdkmath.NewUint(1),
						IncrementOwnedTimesBy: sdkmath.NewUint(0),
					},
				},
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                             hex.EncodeToString(rootHash),
						ExpectedProofLength:              sdk.NewUint(2),
						MaxOneUsePerLeaf:                 true,
						UseLeafIndexForTransferOrder: false,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes:       GetFullUintRanges(),
		BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateFromApproval: &types.ApprovalIdDetails{
					ApprovalId:   "testing232",
					Address: "",
					ApprovalLevel: "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
					{
						
							Leaf: bobLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnedTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnedTimes: GetFullUintRanges(),
	}}, aliceBalance.Balances)
}


func (suite *TestSuite) TestIncrementsTransferAsMuchAsPossibleOneTxWithLeafIndexOrder() {
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(1),
				},
				PredeterminedBalances: &types.PredeterminedBalances{
					OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{ UseOverallNumTransfers: true },
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds: GetOneUintRange(),
								Amount:   sdkmath.NewUint(1),
								OwnedTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy: sdkmath.NewUint(1),
						IncrementOwnedTimesBy: sdkmath.NewUint(0),
					},
				},
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                             hex.EncodeToString(rootHash),
						ExpectedProofLength:              sdk.NewUint(2),
						MaxOneUsePerLeaf:                 true,
						UseLeafIndexForTransferOrder: true,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes:       GetFullUintRanges(),
		BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateFromApproval: &types.ApprovalIdDetails{
					ApprovalId:   "testing232",
					Address: "",
					ApprovalLevel: "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
					{
						
							Leaf: bobLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnedTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       []*types.UintRange{
			GetTwoUintRanges()[0],
		},
		OwnedTimes: GetFullUintRanges(),
	}}, aliceBalance.Balances)
}



func (suite *TestSuite) TestManualTransferDefinitionWithIncrements() {
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(1),
				},
				PredeterminedBalances: &types.PredeterminedBalances{
					OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{ UseOverallNumTransfers: true },
					ManualBalances: []*types.ManualBalances{
						{
							Balances: []*types.Balance{
								{
									BadgeIds: GetOneUintRange(),
									Amount:   sdkmath.NewUint(1),
									OwnedTimes: GetFullUintRanges(),
								},
							},
						},
					
						{
							Balances: []*types.Balance{
								{
									BadgeIds: GetTopHalfUintRanges(),
									Amount:   sdkmath.NewUint(1),
									OwnedTimes: GetFullUintRanges(),
								},
							},
						},
					},
				},
			
				Uri: "",
				MerkleChallenges: []*types.MerkleChallenge{
					{
						ChallengeId: "testchallenge",
						Root:                             hex.EncodeToString(rootHash),
						ExpectedProofLength:              sdk.NewUint(2),
						MaxOneUsePerLeaf:                 true,
						UseLeafIndexForTransferOrder: true,
					},
				},
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		OwnedTimes:       GetFullUintRanges(),
		BadgeIds:             GetFullUintRanges(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateFromApproval: &types.ApprovalIdDetails{
					ApprovalId:   "testing232",
					Address: "",
					ApprovalLevel: "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
					{
						
							Leaf: bobLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnedTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       []*types.UintRange{
			GetTopHalfUintRanges()[0],
		},
		OwnedTimes: GetFullUintRanges(),
	}}, aliceBalance.Balances)
}




func (suite *TestSuite) TestRequestMalformedPredeterminedTransfer() {
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
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{ UseOverallNumTransfers: true },
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds: GetBottomHalfUintRanges(),
							Amount:   sdkmath.NewUint(1),
							OwnedTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy: sdkmath.NewUint(1),
					IncrementOwnedTimesBy: sdkmath.NewUint(0),
				},
			},
			Uri: "",
			MerkleChallenges: []*types.MerkleChallenge{
				{
					ChallengeId: "testchallenge",
					Root:                             hex.EncodeToString(rootHash),
					ExpectedProofLength:              sdk.NewUint(2),
					MaxOneUsePerLeaf:                 true,
					UseLeafIndexForTransferOrder: true,
				},
			},
			ApprovalId:            "testing232",
			OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
		},
	},
	AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},	
	TransferTimes:        GetFullUintRanges(),
		OwnedTimes:       GetFullUintRanges(),
		BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

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
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						BadgeIds: GetOneUintRange(),
						Amount:   sdkmath.NewUint(2),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						BadgeIds: GetFullUintRanges(),
						Amount:   sdkmath.NewUint(1),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						BadgeIds: GetOneUintRange(),
						Amount:   sdkmath.NewUint(1),
						OwnedTimes: GetFullUintRanges(),
					},
				},
				MerkleProofs: []*types.MerkleProof{
					{
						
							Leaf: aliceLeaf,
							Aunts: []*types.MerklePathItem{
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
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestMustOwnBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[1].ApprovalDetails[0].MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdk.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds: GetFullUintRanges(),
			OwnedTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
			
				Uri: "",
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
				OverridesToApprovedIncomingTransfers:   true,
			},

		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		OwnedTimes:       GetFullUintRanges(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

		
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnedTimes: 		GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestMustOwnBadgesDoesntOwnBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[1].ApprovalDetails[0].MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdk.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds: GetFullUintRanges(),
			OwnedTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
					MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
				
				Uri: "",
					
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		OwnedTimes:       		GetFullUintRanges(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnedTimes: 		GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}


func (suite *TestSuite) TestMustOwnBadgesMustOwnZero() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[1].ApprovalDetails[0].MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdk.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(0),
				End:   sdkmath.NewUint(0),
			},
			BadgeIds: GetFullUintRanges(),
			OwnedTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
				Uri: "",
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
				OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		OwnedTimes:       		GetFullUintRanges(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

		
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnedTimes: 		GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnedTimes: 		GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")
}


func (suite *TestSuite) TestMustOwnBadgesMustOwnGreaterThan() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[1].ApprovalDetails[0].MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdk.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(2),
				End:   sdkmath.NewUint(100),
			},
			BadgeIds: GetFullUintRanges(),
			OwnedTimes: GetFullUintRanges(),
		},
	}
	
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers = append(collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers, &types.CollectionApprovedTransfer{
		ApprovalDetails: []*types.ApprovalDetails{
			{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(10),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					OverallApprovalAmount: sdkmath.NewUint(10),
				},
				
				Uri: "",
				ApprovalId:            "testing232",
				OverridesFromApprovedOutgoingTransfers: true,
		OverridesToApprovedIncomingTransfers:   true,
			},
		},
		AllowedCombinations: []*types.IsCollectionTransferAllowed{{
			IsAllowed: true,
		}},
		TransferTimes:        GetFullUintRanges(),
		BadgeIds:             GetOneUintRange(),
		OwnedTimes:       GetFullUintRanges(),
		FromMappingId:        "Mint",
		ToMappingId:          "All",
		InitiatedByMappingId: "All",

	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnedTimes: 		GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdk.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnedTimes: 		GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}