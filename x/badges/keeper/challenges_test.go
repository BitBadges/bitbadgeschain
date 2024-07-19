package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"math"

	"bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNoMerkleChallengeWorking() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.MerkleChallenges = []*types.MerkleChallenge{}

	CreateCollections(suite, wctx, collectionsToCreate)
	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	_, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite.ctx,
		collection,
		&types.Transfer{
			From:        bob,
			ToAddresses: []string{alice},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       GetTopHalfUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
		alice,
		alice,
	)
	suite.Require().Nil(err, "Error getting user balance: %s")
}

func (suite *TestSuite) TestMerkleChallengeInvalidSolutions() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.MerkleChallenges = []*types.MerkleChallenge{
		{
			Root:               "sample",
			ChallengeTrackerId: "testchallenge",
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	_, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite.ctx,
		collection,
		&types.Transfer{
			From:        bob,
			ToAddresses: []string{alice},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       GetTopHalfUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
		alice,
		alice,
	)

	suite.Require().Error(err, "Error getting user balance: %s")

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite.ctx,
		collection,
		&types.Transfer{
			From:        bob,
			ToAddresses: []string{alice},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       GetTopHalfUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			MerkleProofs: []*types.MerkleProof{

				{
					Aunts: []*types.MerklePathItem{},
					Leaf:  "sample",
				},
			},
		},
		alice,
		alice,
	)
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{

		ApprovalId: "asadsdas",
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{

					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
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
		CollectionId: sdkmath.NewUint(1),
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
		CollectionId: sdkmath.NewUint(1),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{

					Root:                    hex.EncodeToString(rootHash),
					ExpectedProofLength:     sdkmath.NewUint(2),
					MaxUsesPerLeaf:          sdkmath.NewUint(1),
					UseCreatorAddressAsLeaf: true,
				},
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{
					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(5),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
					UseMerkleChallengeLeafIndex: true,
					ChallengeTrackerId:          "testchallenge",
				},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				},
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{
					ChallengeTrackerId:  "testchallenge",
					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId = "testchallenge"
	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
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
		CollectionId: sdkmath.NewUint(1),
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

func (suite *TestSuite) TestIncrementsMismatchingTrackerId() {
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
					UseMerkleChallengeLeafIndex: true,
					ChallengeTrackerId:          "testchallenge",
				},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				},
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{
					ChallengeTrackerId:  "testchallenge",
					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId = "mismatched id"
	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
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
		CollectionId: sdkmath.NewUint(1),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				},
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{

					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				},
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{

					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				},
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseMerkleChallengeLeafIndex: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				},
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{

					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			}, OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals: true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				},
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseMerkleChallengeLeafIndex: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				},
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{

					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId:        "asadsdas",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				},
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				},
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{
					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
		Amount: sdkmath.NewUint(1),
		BadgeIds: []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		BadgeIds: []*types.UintRange{
			GetTwoUintRanges()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
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
		OwnershipTimes: GetFullUintRanges(),
	})
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UsePerToAddressNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				},
			},
			MerkleChallenges: []*types.MerkleChallenge{
				{
					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId:        "asadsdas",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
		Amount: sdkmath.NewUint(1),
		BadgeIds: []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		BadgeIds: []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				},
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{

					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
		Amount: sdkmath.NewUint(1),
		BadgeIds: []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		BadgeIds: []*types.UintRange{
			GetTwoUintRanges()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				ManualBalances: []*types.ManualBalances{
					{
						Balances: []*types.Balance{
							{
								BadgeIds:       GetOneUintRange(),
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
					},

					{
						Balances: []*types.Balance{
							{
								BadgeIds:       GetTopHalfUintRanges(),
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
					},
				},
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{

					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
		Amount: sdkmath.NewUint(1),
		BadgeIds: []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		BadgeIds: []*types.UintRange{
			GetTopHalfUintRanges()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							BadgeIds:       GetBottomHalfUintRanges(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementBadgeIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				},
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{
					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
				},
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
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
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						BadgeIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(2),
						OwnershipTimes: GetFullUintRanges(),
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
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						BadgeIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
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
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						BadgeIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestMustOwnBadgesMustSatisfyForAllAssets() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestMustOwnBadgesMustSatisfyForAllAssets2() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(2),
			},
			BadgeIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(2),
			},
			BadgeIds:                GetFullUintRanges(),
			OwnershipTimes:          GetFullUintRanges(),
			MustSatisfyForAllAssets: true,
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestMustOwnBadgesMustOwnOne() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestMustOwnBadgesMustOwnOne2() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(2),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(2),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestMustOwnBadgesDoesntOwnBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(0),
				End:   sdkmath.NewUint(0),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
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
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.MustOwnBadges = []*types.MustOwnBadges{
		{
			CollectionId: sdkmath.NewUint(1),
			AmountRange: &types.UintRange{
				Start: sdkmath.NewUint(2),
				End:   sdkmath.NewUint(100),
			},
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestMultipleApprovalCriteria() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		// ManualBalances: &types.ManualBalances{},
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					BadgeIds:       GetBottomHalfUintRanges(),
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementBadgeIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	deepCopy := *collectionsToCreate[0].CollectionApprovals[1]

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &deepCopy)

	collectionsToCreate[0].CollectionApprovals[2].ApprovalId = "fasdfasdf"
	collectionsToCreate[0].CollectionApprovals[2].ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(1000),
		},
		ApprovalAmounts: &types.ApprovalAmounts{
			PerFromAddressApprovalAmount: sdkmath.NewUint(1),
		},
		PredeterminedBalances: &types.PredeterminedBalances{
			// ManualBalances: &types.ManualBalances{},
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						BadgeIds:       GetTopHalfUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				IncrementBadgeIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	//Fails because we do not take overflows
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")
}

func (suite *TestSuite) TestMultipleApprovalCriteriaPrioritizedApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		collectionsToCreate[0].CollectionApprovals[0],
		{
			FromListId:        "AllWithoutMint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			ApprovalId:        "asadsdas",

			TransferTimes:  GetFullUintRanges(),
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			ApprovalCriteria: &types.ApprovalCriteria{
				ApprovalAmounts: &types.ApprovalAmounts{},
				MaxNumTransfers: &types.MaxNumTransfers{},
				PredeterminedBalances: &types.PredeterminedBalances{
					// ManualBalances: &types.ManualBalances{},
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy:       sdkmath.NewUint(0),
						IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					},
					OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
						UseOverallNumTransfers: true,
					},
				},
			},
		},
		{
			FromListId:        "AllWithoutMint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			ApprovalId:        "target approval",

			TransferTimes:  GetFullUintRanges(),
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			ApprovalCriteria: &types.ApprovalCriteria{
				ApprovalAmounts: &types.ApprovalAmounts{},
				MaxNumTransfers: &types.MaxNumTransfers{},
				PredeterminedBalances: &types.PredeterminedBalances{
					// ManualBalances: &types.ManualBalances{},
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(1000)}},
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy:       sdkmath.NewUint(0),
						IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					},
					OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
						UseOverallNumTransfers: true,
					},
				},
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	bobBalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].Amount)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].BadgeIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(math.MaxUint64), bobBalance.Balances[0].BadgeIds[0].End)

	//Fails because we do not take overflows
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances:    []*types.Balance{},
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "target approval",
					ApprovalLevel:   "collection",
					ApproverAddress: "",
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "target approval",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].Amount)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].BadgeIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(999), bobBalance.Balances[0].BadgeIds[0].End)

	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[1].Amount)
	suite.Require().Equal(sdkmath.NewUint(1001), bobBalance.Balances[1].BadgeIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(math.MaxUint64), bobBalance.Balances[1].BadgeIds[0].End)
}

func (suite *TestSuite) TestMultipleApprovalCriteriaPrioritizedApprovalsOnlyCheckPrioritized() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		collectionsToCreate[0].CollectionApprovals[0],
		{
			FromListId:        "AllWithoutMint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			ApprovalId:        "asadsdas",

			TransferTimes:  GetFullUintRanges(),
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			ApprovalCriteria: &types.ApprovalCriteria{
				ApprovalAmounts: &types.ApprovalAmounts{},
				MaxNumTransfers: &types.MaxNumTransfers{},
				PredeterminedBalances: &types.PredeterminedBalances{
					// ManualBalances: &types.ManualBalances{},
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy:       sdkmath.NewUint(0),
						IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					},
					OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
						UseOverallNumTransfers: true,
					},
				},
			},
		},
		{
			FromListId:        "AllWithoutMint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			ApprovalId:        "target approval",

			TransferTimes:  GetFullUintRanges(),
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			ApprovalCriteria: &types.ApprovalCriteria{
				ApprovalAmounts: &types.ApprovalAmounts{},
				MaxNumTransfers: &types.MaxNumTransfers{},
				PredeterminedBalances: &types.PredeterminedBalances{
					// ManualBalances: &types.ManualBalances{},
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(1000)}},
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
						IncrementBadgeIdsBy:       sdkmath.NewUint(0),
						IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					},
					OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
						UseOverallNumTransfers: true,
					},
				},
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	bobBalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].Amount)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].BadgeIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(math.MaxUint64), bobBalance.Balances[0].BadgeIds[0].End)

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances:    []*types.Balance{},
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "target approval",
					ApprovalLevel:   "collection",
					ApproverAddress: "",
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "random approval",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances:    []*types.Balance{},
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "target approval",
					ApprovalLevel:   "collection",
					ApproverAddress: "",
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "random approval",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	bobBalance, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].Amount)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].BadgeIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(999), bobBalance.Balances[0].BadgeIds[0].End)

	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[1].Amount)
	suite.Require().Equal(sdkmath.NewUint(1001), bobBalance.Balances[1].BadgeIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(math.MaxUint64), bobBalance.Balances[1].BadgeIds[0].End)
}

func (suite *TestSuite) TestMultipleApprovalCriteriaSameAmountTrackerId() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria = nil
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria = nil

	collectionsToCreate[0].CollectionApprovals[1].ApprovalId = "asadsdasfghdsfasdfasdf"
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		// ManualBalances: &types.ManualBalances{},
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					BadgeIds:       GetBottomHalfUintRanges(),
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementBadgeIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	deepCopy := *collectionsToCreate[0].CollectionApprovals[1]
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &deepCopy)
	collectionsToCreate[0].CollectionApprovals[2].ApprovalId = "asadsdasfghaaadsd"

	collectionsToCreate[0].CollectionApprovals[2].ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(1000),
		},
		ApprovalAmounts: &types.ApprovalAmounts{
			PerFromAddressApprovalAmount: sdkmath.NewUint(1),
		},
		PredeterminedBalances: &types.PredeterminedBalances{
			// ManualBalances: &types.ManualBalances{},
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						BadgeIds:       GetBottomHalfUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				IncrementBadgeIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdasfghd",

		TransferTimes:     GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	//Fails because we do not take overflows
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        alice,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge: %s")
}
