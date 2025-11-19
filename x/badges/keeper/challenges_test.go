package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"time"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

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
					TokenIds:       GetTopHalfUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
		},
		keeper.TransferMetadata{
			To:              alice,
			From:            bob,
			InitiatedBy:     alice,
			ApproverAddress: "",
			ApprovalLevel:   "collection",
		},
		&keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}},
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
					TokenIds:       GetTopHalfUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
		keeper.TransferMetadata{
			To:              alice,
			From:            bob,
			InitiatedBy:     alice,
			ApproverAddress: "",
			ApprovalLevel:   "collection",
		},
		&keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}},
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
					TokenIds:       GetTopHalfUintRanges(),
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
		keeper.TransferMetadata{
			To:              alice,
			From:            bob,
			InitiatedBy:     alice,
			ApproverAddress: "",
			ApprovalLevel:   "collection",
		},
		&keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}},
	)
	suite.Require().Error(err, "Error getting user balance: %s")
}

func (suite *TestSuite) TestLeafSignature() {
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
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{}, &types.CollectionApproval{

		ApprovalId: "asadsdas",
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			MerkleChallenges: []*types.MerkleChallenge{
				{

					Root:                hex.EncodeToString(rootHash),
					ExpectedProofLength: sdkmath.NewUint(2),
					MaxUsesPerLeaf:      sdkmath.NewUint(1),
					LeafSigner:          "0xa612B14Ff99DAe9FBC9613bF4553781086c5F887",
				},
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	transfers := &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	}

	err = TransferTokens(suite, wctx, transfers)
	suite.Require().Error(err, "Error transferring token: %s")

	transfers.Transfers[0].MerkleProofs[0].LeafSignature = "0x53277c915e10b01e32878284293809e171976a8e987f211c3d106a2afccdd85072cfab0f12188ff653f8638d73c0e9b18c8e9892b6aa484799f4076002a66cb71b"
	err = TransferTokens(suite, wctx, transfers)
	suite.Require().Nil(err, "Error transferring token: %s")
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
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{}, &types.CollectionApproval{

		ApprovalId: "asadsdas",
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
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
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
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
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Error(err, "Error transferring token: %s")
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
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
		TokenIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
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
	suite.Require().Error(err, "Error transferring token: %s")
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
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
		TokenIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
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
	suite.Require().Error(err, "Error transferring token: %s")
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
					UseMerkleChallengeLeafIndex: true,
					ChallengeTrackerId:          "testchallenge",
				},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
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
		TokenIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId = "testchallenge"
	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
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
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
					UseMerkleChallengeLeafIndex: true,
					ChallengeTrackerId:          "testchallenge",
				},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
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
		TokenIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId = "mismatched id"
	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
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
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
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
	suite.Require().Error(err, "Error transferring token: %s")
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
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
		TokenIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
					Version:         sdkmath.NewUint(0),
				},
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       GetOneUintRange(),
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
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
		TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
					Version:         sdkmath.NewUint(0),
				},
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       GetOneUintRange(),
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
					Version:         sdkmath.NewUint(0),
				},
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       GetTwoUintRanges(),
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseMerkleChallengeLeafIndex: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
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
		TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
					Version:         sdkmath.NewUint(0),
				},
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       GetTwoUintRanges(),
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseMerkleChallengeLeafIndex: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
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
		TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
					Version:         sdkmath.NewUint(0),
				},
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       GetOneUintRange(),
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
					Version:         sdkmath.NewUint(0),
				},
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						TokenIds:       GetFullUintRanges(),
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
	suite.Require().Error(err, "Error transferring token: %s")
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
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
		TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
					Version:         sdkmath.NewUint(0),
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
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
	collectionsToCreate[0].TokensToCreate = append(collectionsToCreate[0].TokensToCreate, &types.Balance{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       GetFullUintRanges(),
		OwnershipTimes: GetFullUintRanges(),
	})
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UsePerToAddressNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
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
		TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
					Version:         sdkmath.NewUint(0),
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
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
		TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
					Version:         sdkmath.NewUint(0),
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				ManualBalances: []*types.ManualBalances{
					{
						Balances: []*types.Balance{
							{
								TokenIds:       GetOneUintRange(),
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
					},

					{
						Balances: []*types.Balance{
							{
								TokenIds:       GetTopHalfUintRanges(),
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
		TokenIds:          GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
					ApprovalId:      "asadsdas",
					ApproverAddress: "",
					ApprovalLevel:   "collection",
					Version:         sdkmath.NewUint(0),
				},
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
	suite.Require().Nil(err, "Error transferring token: %s")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting user balance: %s")

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
			GetOneUintRange()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
	}}, bobBalance.Balances)

	AssertBalancesEqual(suite, []*types.Balance{{
		Amount: sdkmath.NewUint(1),
		TokenIds: []*types.UintRange{
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(1),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetBottomHalfUintRanges(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:       sdkmath.NewUint(1),
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
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
		TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)
	// collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10000),
						TokenIds:       GetFullUintRanges(),
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
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetOneUintRange(),
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
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
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
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob, alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetOneUintRange(),
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
	suite.Require().Error(err, "Error transferring token: %s")
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
					TokenIds:       GetBottomHalfUintRanges(),
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(0),
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
			AmountTrackerId:        "test-tracker",
		},
		ApprovalAmounts: &types.ApprovalAmounts{
			PerFromAddressApprovalAmount: sdkmath.NewUint(1),
			AmountTrackerId:              "test-tracker",
		},
		PredeterminedBalances: &types.PredeterminedBalances{
			// ManualBalances: &types.ManualBalances{},
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						TokenIds:       GetTopHalfUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(0),
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	//Fails because we do not take overflows
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")
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
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			ApprovalCriteria: &types.ApprovalCriteria{
				ApprovalAmounts: &types.ApprovalAmounts{},
				MaxNumTransfers: &types.MaxNumTransfers{},
				PredeterminedBalances: &types.PredeterminedBalances{
					// ManualBalances: &types.ManualBalances{},
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
						IncrementTokenIdsBy:       sdkmath.NewUint(0),
						IncrementOwnershipTimesBy: sdkmath.NewUint(0),
						DurationFromTimestamp:     sdkmath.NewUint(0),
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
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			ApprovalCriteria: &types.ApprovalCriteria{
				ApprovalAmounts: &types.ApprovalAmounts{},
				MaxNumTransfers: &types.MaxNumTransfers{},
				PredeterminedBalances: &types.PredeterminedBalances{
					// ManualBalances: &types.ManualBalances{},
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(1000)}},
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
						IncrementTokenIdsBy:       sdkmath.NewUint(0),
						IncrementOwnershipTimesBy: sdkmath.NewUint(0),
						DurationFromTimestamp:     sdkmath.NewUint(0),
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
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].TokenIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(math.MaxUint64), bobBalance.Balances[0].TokenIds[0].End)

	//Fails because we do not take overflows
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
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
					Version:         sdkmath.NewUint(0),
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")

	bobBalance, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].Amount)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].TokenIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(999), bobBalance.Balances[0].TokenIds[0].End)

	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[1].Amount)
	suite.Require().Equal(sdkmath.NewUint(1001), bobBalance.Balances[1].TokenIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(math.MaxUint64), bobBalance.Balances[1].TokenIds[0].End)
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
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			ApprovalCriteria: &types.ApprovalCriteria{
				ApprovalAmounts: &types.ApprovalAmounts{},
				MaxNumTransfers: &types.MaxNumTransfers{},
				PredeterminedBalances: &types.PredeterminedBalances{
					// ManualBalances: &types.ManualBalances{},
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
						IncrementTokenIdsBy:       sdkmath.NewUint(0),
						IncrementOwnershipTimesBy: sdkmath.NewUint(0),
						DurationFromTimestamp:     sdkmath.NewUint(0),
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
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			ApprovalCriteria: &types.ApprovalCriteria{
				ApprovalAmounts: &types.ApprovalAmounts{},
				MaxNumTransfers: &types.MaxNumTransfers{},
				PredeterminedBalances: &types.PredeterminedBalances{
					// ManualBalances: &types.ManualBalances{},
					IncrementedBalances: &types.IncrementedBalances{
						StartBalances: []*types.Balance{
							{
								TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(1000)}},
								Amount:         sdkmath.NewUint(1),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
						IncrementTokenIdsBy:       sdkmath.NewUint(0),
						IncrementOwnershipTimesBy: sdkmath.NewUint(0),
						DurationFromTimestamp:     sdkmath.NewUint(0),
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
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].TokenIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(math.MaxUint64), bobBalance.Balances[0].TokenIds[0].End)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
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
					Version:         sdkmath.NewUint(0),
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "random approval",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
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
					Version:         sdkmath.NewUint(0),
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")

	bobBalance, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].Amount)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].TokenIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(999), bobBalance.Balances[0].TokenIds[0].End)

	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[1].Amount)
	suite.Require().Equal(sdkmath.NewUint(1001), bobBalance.Balances[1].TokenIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(math.MaxUint64), bobBalance.Balances[1].TokenIds[0].End)
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
					TokenIds:       GetBottomHalfUintRanges(),
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(0),
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
			AmountTrackerId:        "test-tracker",
		},
		ApprovalAmounts: &types.ApprovalAmounts{
			PerFromAddressApprovalAmount: sdkmath.NewUint(1),
			AmountTrackerId:              "test-tracker",
		},
		PredeterminedBalances: &types.PredeterminedBalances{
			// ManualBalances: &types.ManualBalances{},
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						TokenIds:       GetBottomHalfUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(0),
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
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},

			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdasfghd",

		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		OwnershipTimes:    GetFullUintRanges(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	//Fails because we do not take overflows
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        alice,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetBottomHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token: %s")
}

func (suite *TestSuite) TestSequentialTransferApprovalDurationFromNow() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria = nil
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria = nil

	collectionsToCreate[0].CollectionApprovals[1].ApprovalId = "asadsdasfghdsfasdfasdf"

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{},
		MaxNumTransfers: &types.MaxNumTransfers{},
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(1000),
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		}},
	})
	suite.Require().Nil(err)
	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)

	startTime := aliceBalance.Balances[0].OwnershipTimes[0].Start
	endTime := aliceBalance.Balances[0].OwnershipTimes[0].End
	suite.Require().Equal(startTime.Add(sdkmath.NewUint(999)), endTime)

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)
	suite.Require().Equal(bobBalance.Balances[0].Amount, sdkmath.NewUint(1))
	bob0EndTime := bobBalance.Balances[0].OwnershipTimes[0].End
	bob1StartTime := bobBalance.Balances[1].OwnershipTimes[0].Start

	// 1000 difference between the two transfers + 1 on start to account for the inclusive start time
	suite.Require().Equal(bob0EndTime.Add(sdkmath.NewUint(1001)), bob1StartTime)
}

func (suite *TestSuite) TestCantSetBothIncrementOwnershipTimesByAndApprovalDurationFromNow() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria = &types.OutgoingApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				DurationFromTimestamp:     sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(1),
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "Error creating tokens")
}

func (suite *TestSuite) TestCantSetBothIncrementOwnershipTimesByAndRecurringOwnershipTimes() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria = &types.OutgoingApprovalCriteria{
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				DurationFromTimestamp:     sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(1),
			},
		},
	}
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.PredeterminedBalances.IncrementedBalances.RecurringOwnershipTimes = &types.RecurringOwnershipTimes{
		StartTime:          sdkmath.NewUint(10),
		IntervalLength:     sdkmath.NewUint(10),
		ChargePeriodLength: sdkmath.NewUint(10),
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "Error creating tokens")
}

func (suite *TestSuite) TestSequentialTransferApprovalDurationFromNowWithTimestampOverride() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria = nil
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria = nil

	collectionsToCreate[0].CollectionApprovals[1].ApprovalId = "asadsdasfghdsfasdfasdf"

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{},
		MaxNumTransfers: &types.MaxNumTransfers{},
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(1000),
				AllowOverrideTimestamp:    true,
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
			PrecalculationOptions: &types.PrecalculationOptions{
				OverrideTimestamp: sdkmath.NewUint(10),
			},
		}},
	})
	suite.Require().Nil(err)
	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)

	startTime := aliceBalance.Balances[0].OwnershipTimes[0].Start
	endTime := aliceBalance.Balances[0].OwnershipTimes[0].End
	suite.Require().Equal(startTime.Add(sdkmath.NewUint(999)).Uint64(), endTime.Uint64())

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)
	suite.Require().Equal(bobBalance.Balances[0].Amount.Uint64(), sdkmath.NewUint(1).Uint64())
	bob0EndTime := bobBalance.Balances[0].OwnershipTimes[0].End
	bob1StartTime := bobBalance.Balances[1].OwnershipTimes[0].Start

	// 1000 difference between the two transfers + 1 on either end to account for the inclusive end time
	suite.Require().Equal(bob0EndTime.Add(sdkmath.NewUint(1001)).Uint64(), bob1StartTime.Uint64())

	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].Start.Uint64(), sdkmath.NewUint(10).Uint64())
}

func (suite *TestSuite) TestSequentialTransferApprovalDurationFromNowWithTimestampOverrideNotAllowed() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria = nil
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria = nil

	collectionsToCreate[0].CollectionApprovals[1].ApprovalId = "asadsdasfghdsfasdfasdf"

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{},
		MaxNumTransfers: &types.MaxNumTransfers{},
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(1000),
				AllowOverrideTimestamp:    false,
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
			PrecalculationOptions: &types.PrecalculationOptions{
				OverrideTimestamp: sdkmath.NewUint(10),
			},
		}},
	})
	suite.Require().Nil(err)

	// We did not use the overridden timestamp
	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().NotEqual(aliceBalance.Balances[0].OwnershipTimes[0].Start, sdkmath.NewUint(10))
}

func (suite *TestSuite) TestRecurringOwnershipTimes() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria = nil
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria = nil

	collectionsToCreate[0].CollectionApprovals[1].ApprovalId = "asadsdasfghdsfasdfasdf"

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{},
		MaxNumTransfers: &types.MaxNumTransfers{},
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(), // doesn't matter
					},
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(0),
				AllowOverrideTimestamp:    false,
				RecurringOwnershipTimes: &types.RecurringOwnershipTimes{
					StartTime:          sdkmath.NewUint(100),
					IntervalLength:     sdkmath.NewUint(100),
					ChargePeriodLength: sdkmath.NewUint(10),
				},
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// outside grace period
	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(9949))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		}},
	})
	suite.Require().Error(err)

	// before start time
	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(1))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		}},
	})
	suite.Require().Error(err)

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(9999))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		}},
	})
	suite.Require().Nil(err)

	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].Start, sdkmath.NewUint(10000))
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].End, sdkmath.NewUint(10000).Add(sdkmath.NewUint(99)))

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(10099))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		}},
	})
	suite.Require().Nil(err)

	aliceBalance, _ = GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].Start, sdkmath.NewUint(10000))
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].End, sdkmath.NewUint(10000).Add(sdkmath.NewUint(199)))

	// skip a month
	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(10299))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		}},
	})
	suite.Require().Nil(err)

	aliceBalance, _ = GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].Start, sdkmath.NewUint(10000))
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].End, sdkmath.NewUint(10000).Add(sdkmath.NewUint(199)))
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[1].Start, sdkmath.NewUint(10300))
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[1].End, sdkmath.NewUint(10300).Add(sdkmath.NewUint(99)))
}

func (suite *TestSuite) TestSubscriptionApproach() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria = nil
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria = nil

	collectionsToCreate[0].CollectionApprovals[1].ApprovalId = "asadsdasfghdsfasdfasdf"

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{},
		MaxNumTransfers: &types.MaxNumTransfers{},
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(), // doesn't matter
					},
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(1000),
				AllowOverrideTimestamp:    true,
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
	}
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "transfer-approval",
		FromListId:        collectionsToCreate[0].DefaultIncomingApprovals[0].FromListId,
		InitiatedByListId: collectionsToCreate[0].DefaultIncomingApprovals[0].InitiatedByListId,
		ToListId:          collectionsToCreate[0].DefaultOutgoingApprovals[0].ToListId,
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TransferTimes:     GetFullUintRanges(),
	})

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(1500))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
			PrecalculationOptions: &types.PrecalculationOptions{
				OverrideTimestamp: sdkmath.NewUint(2000),
			},
		}},
	})
	suite.Require().Nil(err)

	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].Start, sdkmath.NewUint(2000))
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].End, sdkmath.NewUint(2000).Add(sdkmath.NewUint(999)))

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 charlie,
		CollectionId:            collection.CollectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals: []*types.UserIncomingApproval{
			{
				ApprovalId:        "asadsdasfghdsfasdfasdf",
				FromListId:        "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				TokenIds:          GetFullUintRanges(),

				ApprovalCriteria: &types.IncomingApprovalCriteria{
					MaxNumTransfers: &types.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(0),
					},
					ApprovalAmounts: &types.ApprovalAmounts{
						OverallApprovalAmount: sdkmath.NewUint(0),
					},
					PredeterminedBalances: &types.PredeterminedBalances{
						IncrementedBalances: &types.IncrementedBalances{
							StartBalances: []*types.Balance{
								{Amount: sdkmath.NewUint(1), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()},
							},
							IncrementTokenIdsBy:       sdkmath.NewUint(0),
							IncrementOwnershipTimesBy: sdkmath.NewUint(0),
							DurationFromTimestamp:     sdkmath.NewUint(0),
							AllowOverrideTimestamp:    true,
							RecurringOwnershipTimes: &types.RecurringOwnershipTimes{
								StartTime:          sdkmath.NewUint(1000),
								IntervalLength:     sdkmath.NewUint(1000),
								ChargePeriodLength: sdkmath.NewUint(1000),
							},
						},
						OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
							UseOverallNumTransfers: true,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        alice,
			ToAddresses: []string{charlie},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "incoming", ApproverAddress: charlie, Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "incoming",
				ApproverAddress: charlie,
				Version:         sdkmath.NewUint(0),
			},
		}},
	})
	suite.Require().Nil(err)

	aliceBalance, _ = GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().Equal(len(aliceBalance.Balances), 0)
}

func (suite *TestSuite) TestRecurringOwnershipTimesChargeFirstInterval() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria = nil
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria = nil

	collectionsToCreate[0].CollectionApprovals[1].ApprovalId = "asadsdasfghdsfasdfasdf"

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{},
		MaxNumTransfers: &types.MaxNumTransfers{},
		PredeterminedBalances: &types.PredeterminedBalances{
			IncrementedBalances: &types.IncrementedBalances{
				StartBalances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(), // doesn't matter
					},
				},
				IncrementTokenIdsBy:       sdkmath.NewUint(0),
				IncrementOwnershipTimesBy: sdkmath.NewUint(0),
				DurationFromTimestamp:     sdkmath.NewUint(0),
				AllowOverrideTimestamp:    false,
				RecurringOwnershipTimes: &types.RecurringOwnershipTimes{
					StartTime:          sdkmath.NewUint(100),
					IntervalLength:     sdkmath.NewUint(100),
					ChargePeriodLength: sdkmath.NewUint(10),
				},
			},
			OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
				UseOverallNumTransfers: true,
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(99))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances:    []*types.Balance{},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "asadsdasfghdsfasdfasdf", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
			PrecalculateBalancesFromApproval: &types.ApprovalIdentifierDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		}},
	})
	suite.Require().Nil(err)

	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].Start, sdkmath.NewUint(100))
	suite.Require().Equal(aliceBalance.Balances[0].OwnershipTimes[0].End, sdkmath.NewUint(100).Add(sdkmath.NewUint(99)))
}

func (suite *TestSuite) TestTokenIdsOverride() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
					UseOverallNumTransfers: true,
				},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:            sdkmath.NewUint(1),
					IncrementOwnershipTimesBy:      sdkmath.NewUint(0),
					DurationFromTimestamp:          sdkmath.NewUint(0),
					AllowOverrideWithAnyValidToken: true,
				},
			},

			MerkleChallenges:               []*types.MerkleChallenge{},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId = "testchallenge"
	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err)

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances.IncrementedBalances.IncrementTokenIdsBy = sdkmath.NewUint(0)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	msg := &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
				MerkleProofs:         []*types.MerkleProof{},
			},
		},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Error(err, "Error transferring token: %s")

	msg.Transfers[0].PrecalculationOptions = &types.PrecalculationOptions{
		TokenIdsOverride: []*types.UintRange{{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(2)}},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Error(err)

	msg.Transfers[0].PrecalculationOptions.TokenIdsOverride = []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Nil(err)
}

func (suite *TestSuite) TestTokenIdsOverrideWithMoreThanOneBadge() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
					UseOverallNumTransfers: true,
				},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:            sdkmath.NewUint(1),
					IncrementOwnershipTimesBy:      sdkmath.NewUint(0),
					DurationFromTimestamp:          sdkmath.NewUint(0),
					AllowOverrideWithAnyValidToken: true,
				},
			},

			MerkleChallenges:               []*types.MerkleChallenge{},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId = "testchallenge"
	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err)

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances.IncrementedBalances.IncrementTokenIdsBy = sdkmath.NewUint(0)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	msg := &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
				MerkleProofs:         []*types.MerkleProof{},
			},
		},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Error(err, "Error transferring token: %s")

	msg.Transfers[0].PrecalculationOptions = &types.PrecalculationOptions{
		TokenIdsOverride: []*types.UintRange{{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(2)}},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Error(err)

	msg.Transfers[0].PrecalculationOptions = &types.PrecalculationOptions{
		TokenIdsOverride: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Error(err)

	msg.Transfers[0].PrecalculationOptions.TokenIdsOverride = []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Nil(err)
}

func (suite *TestSuite) TestTokenIdsOverrideNoPrecalcSpecified() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := *new(error)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "test-tracker",
			},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
					UseOverallNumTransfers: true,
				},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances: []*types.Balance{
						{
							TokenIds:       GetOneUintRange(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					IncrementTokenIdsBy:            sdkmath.NewUint(1),
					IncrementOwnershipTimesBy:      sdkmath.NewUint(0),
					DurationFromTimestamp:          sdkmath.NewUint(0),
					AllowOverrideWithAnyValidToken: true,
				},
			},

			MerkleChallenges:               []*types.MerkleChallenge{},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},

		ApprovalId: "asadsdas",

		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "Mint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
	})

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId = "testchallenge"
	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err)

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.PredeterminedBalances.IncrementedBalances.IncrementTokenIdsBy = sdkmath.NewUint(0)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	msg := &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
				MerkleProofs:         []*types.MerkleProof{},
			},
		},
	}

	// msg.Transfers[0].PrecalculationOptions.TokenIdsOverride = []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Error(err)
}
