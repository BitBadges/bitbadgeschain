package keeper_test

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
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
		"collection",
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
		"collection",
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
		"collection",
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
				PrecalculationOptions: &types.PrecalculationOptions{
					OverrideTimestamp: sdkmath.NewUint(10),
				},
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
				PrecalculationOptions: &types.PrecalculationOptions{
					OverrideTimestamp: sdkmath.NewUint(10),
				},
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
				ApprovalId:      "asadsdasfghdsfasdfasdf",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
				PrecalculationOptions: &types.PrecalculationOptions{
					OverrideTimestamp: sdkmath.NewUint(2000),
				},
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
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

	if msg.Transfers[0].PrecalculateBalancesFromApproval == nil {
		msg.Transfers[0].PrecalculateBalancesFromApproval = &types.PrecalculateBalancesFromApprovalDetails{}
	}
	msg.Transfers[0].PrecalculateBalancesFromApproval.PrecalculationOptions = &types.PrecalculationOptions{
		TokenIdsOverride: []*types.UintRange{{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(2)}},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Error(err)

	msg.Transfers[0].PrecalculateBalancesFromApproval.PrecalculationOptions.TokenIdsOverride = []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}

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

	if msg.Transfers[0].PrecalculateBalancesFromApproval == nil {
		msg.Transfers[0].PrecalculateBalancesFromApproval = &types.PrecalculateBalancesFromApprovalDetails{}
	}
	msg.Transfers[0].PrecalculateBalancesFromApproval.PrecalculationOptions = &types.PrecalculationOptions{
		TokenIdsOverride: []*types.UintRange{{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(2)}},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Error(err)

	msg.Transfers[0].PrecalculateBalancesFromApproval.PrecalculationOptions = &types.PrecalculationOptions{
		TokenIdsOverride: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Error(err)

	msg.Transfers[0].PrecalculateBalancesFromApproval.PrecalculationOptions.TokenIdsOverride = []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}

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

// Helper functions for ETH signature testing

// generateTestETHPrivateKey generates a deterministic test Ethereum private key
// Returns (privateKeyHex, addressHex, error)
func generateTestETHPrivateKey() (string, string, error) {
	// Use a fixed private key for deterministic testing
	// This is a well-known test private key: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
	privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	privateKey, err := ethcrypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", err
	}

	address := ethcrypto.PubkeyToAddress(*publicKeyECDSA)

	return privateKeyHex, address.Hex(), nil
}

// generateETHSignature generates a valid ETH signature for the given message components
// Signature scheme: ETHSign(nonce + "-" + initiatorAddress + "-" + collectionId + "-" + approverAddress + "-" + approvalLevel + "-" + approvalId + "-" + challengeId)
// Uses ethers.signMessage format (EIP-191) which prefixes message with "\x19Ethereum Signed Message:\n<length>"
func generateETHSignature(nonce, initiatorAddress, collectionId, approverAddress, approvalLevel, approvalId, challengeId string, privateKeyHex string) (string, error) {
	// Construct the message to sign
	message := nonce + "-" + initiatorAddress + "-" + collectionId + "-" + approverAddress + "-" + approvalLevel + "-" + approvalId + "-" + challengeId

	// Get private key
	privateKey, err := ethcrypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}

	// Use ethers.signMessage format (EIP-191): prefix with "\x19Ethereum Signed Message:\n<length>"
	// This matches what ethers.Wallet.signMessage() does
	messageBytes := []byte(message)
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(messageBytes))
	prefixedMessage := append([]byte(prefix), messageBytes...)

	// Hash the prefixed message and sign
	hash := ethcrypto.Keccak256Hash(prefixedMessage)
	signature, err := ethcrypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	// Convert to hex string (with 0x prefix for compatibility with sigverify)
	return "0x" + hex.EncodeToString(signature), nil
}

// TestETHSignatureChallenge_ValidSignature tests that a valid ETH signature is accepted
func (suite *TestSuite) TestETHSignatureChallenge_ValidSignature() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Generate test private key and address
	_, signerAddress, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)

	// Create collection with ETH signature challenge
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
		{
			Signer:             signerAddress,
			ChallengeTrackerId: "test-challenge-1",
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval so we can mint badges to bob
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob so he can transfer them
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
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "mint-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().NoError(err, "Error minting badges to bob")

	// Generate valid signature
	nonce := "test-nonce-123"
	privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	signature, err := generateETHSignature(
		nonce,
		alice,              // initiatorAddress
		"1",                // collectionId
		"",                 // approverAddress (empty for collection level)
		"collection",       // approvalLevel
		"test",             // approvalId
		"test-challenge-1", // challengeId
		privateKeyHex,
	)
	suite.Require().NoError(err)

	// Create transfer with valid signature
	msg := &types.MsgTransferTokens{
		Creator:      alice,
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
				EthSignatureProofs: []*types.ETHSignatureProof{
					{
						Nonce:     nonce,
						Signature: signature,
					},
				},
			},
		},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().NoError(err, "Valid signature should be accepted")

	// Verify signature tracker was incremented
	signatureKey := keeper.ConstructETHSignatureTrackerKey(
		sdkmath.NewUint(1),
		"",                 // approverAddress
		"collection",       // approvalLevel
		"test",             // approvalId
		"test-challenge-1", // challengeId
		signature,
	)
	numUsed, exists := suite.app.BadgesKeeper.GetETHSignatureTrackerFromStore(suite.ctx, signatureKey)
	suite.Require().True(exists, "Signature tracker should exist")
	suite.Require().Equal(sdkmath.NewUint(1), numUsed, "Signature should be used once")
}

// TestETHSignatureChallenge_InvalidSignature tests various invalid signature scenarios
func (suite *TestSuite) TestETHSignatureChallenge_InvalidSignature() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Generate test private key and address
	_, signerAddress, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)

	// Create collection with ETH signature challenge
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
		{
			Signer:             signerAddress,
			ChallengeTrackerId: "test-challenge-1",
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Test 1: Wrong signer (signature from different address)
	suite.Run("WrongSigner", func() {
		// Generate a different private key
		wrongPrivateKey, err := ethcrypto.GenerateKey()
		suite.Require().NoError(err)
		wrongPrivateKeyHex := hex.EncodeToString(ethcrypto.FromECDSA(wrongPrivateKey))

		nonce := "test-nonce-123"
		signature, err := generateETHSignature(
			nonce,
			alice,
			"1",
			"",
			"collection",
			"test",
			"test-challenge-1",
			wrongPrivateKeyHex,
		)
		suite.Require().NoError(err)

		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     nonce,
							Signature: signature,
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Signature from wrong signer should be rejected")
	})

	// Test 2: Wrong message (signature for different context)
	suite.Run("WrongMessage", func() {
		privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		nonce := "test-nonce-123"
		// Sign for different collectionId
		signature, err := generateETHSignature(
			nonce,
			alice,
			"2", // Wrong collectionId
			"",
			"collection",
			"test",
			"test-challenge-1",
			privateKeyHex,
		)
		suite.Require().NoError(err)

		msg := &types.MsgTransferTokens{
			Creator:      alice,
			CollectionId: sdkmath.NewUint(1), // But trying to use for collection 1
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
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     nonce,
							Signature: signature,
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Signature for wrong context should be rejected")
	})

	// Test 3: Empty signature
	suite.Run("EmptySignature", func() {
		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     "test-nonce-123",
							Signature: "", // Empty signature
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Empty signature should be rejected")
	})

	// Test 4: Empty nonce
	suite.Run("EmptyNonce", func() {
		privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		signature, err := generateETHSignature(
			"test-nonce-123",
			alice,
			"1",
			"",
			"collection",
			"test",
			"test-challenge-1",
			privateKeyHex,
		)
		suite.Require().NoError(err)

		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     "", // Empty nonce
							Signature: signature,
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Empty nonce should be rejected")
	})

	// Test 5: Malformed signature
	suite.Run("MalformedSignature", func() {
		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     "test-nonce-123",
							Signature: "0xinvalid", // Malformed signature
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Malformed signature should be rejected")
	})
}

// TestETHSignatureChallenge_SignatureReuse tests that signatures cannot be reused
func (suite *TestSuite) TestETHSignatureChallenge_SignatureReuse() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Generate test private key and address
	_, signerAddress, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)

	// Create collection with ETH signature challenge
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
		{
			Signer:             signerAddress,
			ChallengeTrackerId: "test-challenge-1",
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval so we can mint badges to bob
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob so he can transfer them
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(2), // Mint 2 for reuse test
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "mint-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().NoError(err, "Error minting badges to bob")

	// Generate valid signature
	nonce := "test-nonce-123"
	privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	signature, err := generateETHSignature(
		nonce,
		alice,
		"1",
		"",
		"collection",
		"test",
		"test-challenge-1",
		privateKeyHex,
	)
	suite.Require().NoError(err)

	// First transfer - should succeed
	msg1 := &types.MsgTransferTokens{
		Creator:      alice,
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
				EthSignatureProofs: []*types.ETHSignatureProof{
					{
						Nonce:     nonce,
						Signature: signature,
					},
				},
			},
		},
	}

	err = TransferTokens(suite, wctx, msg1)
	suite.Require().NoError(err, "First use of signature should succeed")

	// Verify signature tracker shows it was used
	signatureKey := keeper.ConstructETHSignatureTrackerKey(
		sdkmath.NewUint(1),
		"",
		"collection",
		"test",
		"test-challenge-1",
		signature,
	)
	numUsed, exists := suite.app.BadgesKeeper.GetETHSignatureTrackerFromStore(suite.ctx, signatureKey)
	suite.Require().True(exists)
	suite.Require().Equal(sdkmath.NewUint(1), numUsed, "Signature should be marked as used once")

	// Second transfer with same signature - should fail
	msg2 := &types.MsgTransferTokens{
		Creator:      alice,
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
				EthSignatureProofs: []*types.ETHSignatureProof{
					{
						Nonce:     nonce,
						Signature: signature, // Same signature
					},
				},
			},
		},
	}

	err = TransferTokens(suite, wctx, msg2)
	suite.Require().Error(err, "Reusing signature should be rejected")
}

// TestETHSignatureChallenge_CrossContextRejection tests that signatures are rejected when used in different contexts
func (suite *TestSuite) TestETHSignatureChallenge_CrossContextRejection() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Generate test private key and address
	_, signerAddress, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)

	privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	// Test 1: Different collectionId
	suite.Run("DifferentCollectionId", func() {
		// Create two collections
		collectionsToCreate := GetCollectionsToCreate()
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
			{
				Signer:             signerAddress,
				ChallengeTrackerId: "test-challenge-1",
			},
		}
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

		err = CreateCollections(suite, wctx, collectionsToCreate)
		suite.Require().NoError(err)

		// Create second collection
		collectionsToCreate2 := GetCollectionsToCreate()
		collectionsToCreate2[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
			{
				Signer:             signerAddress,
				ChallengeTrackerId: "test-challenge-1",
			},
		}
		collectionsToCreate2[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
		collectionsToCreate2[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

		err = CreateCollections(suite, wctx, collectionsToCreate2)
		suite.Require().NoError(err)

		// Generate signature for collection 1
		nonce := "test-nonce-123"
		signature, err := generateETHSignature(
			nonce,
			alice,
			"1", // collectionId 1
			"",
			"collection",
			"test",
			"test-challenge-1",
			privateKeyHex,
		)
		suite.Require().NoError(err)

		// Try to use it for collection 2
		msg := &types.MsgTransferTokens{
			Creator:      alice,
			CollectionId: sdkmath.NewUint(2), // Different collection
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
					PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(2)),
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     nonce,
							Signature: signature,
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Signature for different collectionId should be rejected")
	})

	// Test 2: Different approverAddress - tested via DifferentApprovalLevel

	// Test 3: Different approvalLevel
	suite.Run("DifferentApprovalLevel", func() {
		collectionsToCreate := GetCollectionsToCreate()
		// Add outgoing approval
		collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
			{
				Signer:             signerAddress,
				ChallengeTrackerId: "test-challenge-1",
			},
		}

		err = CreateCollections(suite, wctx, collectionsToCreate)
		suite.Require().NoError(err)

		// Generate signature for collection level
		nonce := "test-nonce-789"
		signature, err := generateETHSignature(
			nonce,
			alice,
			"1",
			"",           // approverAddress
			"collection", // approvalLevel
			"test",
			"test-challenge-1",
			privateKeyHex,
		)
		suite.Require().NoError(err)

		// Try to use it for outgoing level (which has approverAddress = bob)
		// The signature was signed with approverAddress = "", but outgoing level uses approverAddress = bob
		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
						{
							ApprovalId:      "test",
							ApprovalLevel:   "outgoing",
							ApproverAddress: bob,
							Version:         sdkmath.NewUint(0),
						},
					},
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     nonce,
							Signature: signature,
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Signature for different approvalLevel should be rejected")
	})

	// Test 4: Different approvalId
	suite.Run("DifferentApprovalId", func() {
		collectionsToCreate := GetCollectionsToCreate()
		// Add second approval with different ID
		collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
			ApprovalId:        "test-2",
			ToListId:          "AllWithoutMint",
			FromListId:        "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetFullUintRanges(),
			ApprovalCriteria: &types.ApprovalCriteria{
				EthSignatureChallenges: []*types.ETHSignatureChallenge{
					{
						Signer:             signerAddress,
						ChallengeTrackerId: "test-challenge-1",
					},
				},
				OverridesToIncomingApprovals:   true,
				OverridesFromOutgoingApprovals: true,
			},
		})

		err = CreateCollections(suite, wctx, collectionsToCreate)
		suite.Require().NoError(err)

		// Generate signature for approvalId = "test"
		nonce := "test-nonce-abc"
		signature, err := generateETHSignature(
			nonce,
			alice,
			"1",
			"",
			"collection",
			"test", // approvalId
			"test-challenge-1",
			privateKeyHex,
		)
		suite.Require().NoError(err)

		// Try to use it for approvalId = "test-2"
		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
						{
							ApprovalId:      "test-2", // Different approvalId
							ApprovalLevel:   "collection",
							ApproverAddress: "",
							Version:         sdkmath.NewUint(0),
						},
					},
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     nonce,
							Signature: signature,
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Signature for different approvalId should be rejected")
	})

	// Test 5: Different challengeId
	suite.Run("DifferentChallengeId", func() {
		collectionsToCreate := GetCollectionsToCreate()
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
			{
				Signer:             signerAddress,
				ChallengeTrackerId: "test-challenge-1",
			},
			{
				Signer:             signerAddress,
				ChallengeTrackerId: "test-challenge-2", // Different challenge
			},
		}
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

		err = CreateCollections(suite, wctx, collectionsToCreate)
		suite.Require().NoError(err)

		// Generate signature for challengeId = "test-challenge-1"
		nonce := "test-nonce-xyz"
		signature, err := generateETHSignature(
			nonce,
			alice,
			"1",
			"",
			"collection",
			"test",
			"test-challenge-1", // challengeId
			privateKeyHex,
		)
		suite.Require().NoError(err)

		// Try to use it for challengeId = "test-challenge-2"
		// We need to provide a signature for challenge-2, but we're providing one for challenge-1
		// The code checks all challenges, so if we only provide signature for challenge-1, challenge-2 will fail
		// Actually, the code requires ALL challenges to be satisfied, so we need signatures for both
		// But if we provide signature for challenge-1 when challenge-2 is being checked, it should fail
		// Let's test by only providing the signature for challenge-1 when both are required
		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     nonce,
							Signature: signature, // Signature for challenge-1, but challenge-2 also needs a signature
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Missing signature for challenge-2 should cause failure")
	})
}

// TestETHSignatureChallenge_MultipleChallenges tests multiple ETH signature challenges
func (suite *TestSuite) TestETHSignatureChallenge_MultipleChallenges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Generate test private key and address
	_, signerAddress, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)

	// Create collection with multiple ETH signature challenges
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
		{
			Signer:             signerAddress,
			ChallengeTrackerId: "test-challenge-1",
		},
		{
			Signer:             signerAddress,
			ChallengeTrackerId: "test-challenge-2",
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval so we can mint badges to bob
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob so he can transfer them
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
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "mint-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().NoError(err, "Error minting badges to bob")

	privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	// Test: All signatures valid - should succeed
	suite.Run("AllSignaturesValid", func() {
		nonce1 := "test-nonce-1"
		signature1, err := generateETHSignature(
			nonce1,
			alice,
			"1",
			"",
			"collection",
			"test",
			"test-challenge-1",
			privateKeyHex,
		)
		suite.Require().NoError(err)

		nonce2 := "test-nonce-2"
		signature2, err := generateETHSignature(
			nonce2,
			alice,
			"1",
			"",
			"collection",
			"test",
			"test-challenge-2",
			privateKeyHex,
		)
		suite.Require().NoError(err)

		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     nonce1,
							Signature: signature1,
						},
						{
							Nonce:     nonce2,
							Signature: signature2,
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().NoError(err, "All valid signatures should be accepted")
	})

	// Test: Missing one signature - should fail
	suite.Run("MissingSignature", func() {
		nonce1 := "test-nonce-3"
		signature1, err := generateETHSignature(
			nonce1,
			alice,
			"1",
			"",
			"collection",
			"test",
			"test-challenge-1",
			privateKeyHex,
		)
		suite.Require().NoError(err)

		// Only provide signature for challenge-1, missing challenge-2
		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     nonce1,
							Signature: signature1,
						},
						// Missing signature for challenge-2
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Missing signature should cause failure")
	})

	// Test: One invalid signature - should fail
	suite.Run("InvalidSignature", func() {
		nonce1 := "test-nonce-4"
		signature1, err := generateETHSignature(
			nonce1,
			alice,
			"1",
			"",
			"collection",
			"test",
			"test-challenge-1",
			privateKeyHex,
		)
		suite.Require().NoError(err)

		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     nonce1,
							Signature: signature1,
						},
						{
							Nonce:     "test-nonce-5",
							Signature: "0xinvalid", // Invalid signature for challenge-2
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Invalid signature should cause failure")
	})
}

// TestETHSignatureChallenge_MissingSignature tests that missing signatures are rejected
func (suite *TestSuite) TestETHSignatureChallenge_MissingSignature() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Generate test private key and address
	_, signerAddress, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)

	// Create collection with ETH signature challenge
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
		{
			Signer:             signerAddress,
			ChallengeTrackerId: "test-challenge-1",
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Create transfer without ETH signature proofs
	msg := &types.MsgTransferTokens{
		Creator:      alice,
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
				EthSignatureProofs:   []*types.ETHSignatureProof{}, // Empty - no signatures
			},
		},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().Error(err, "Missing signature should be rejected")
}

// TestETHSignatureChallenge_EmptyChallenge tests edge cases with empty/nil challenges
func (suite *TestSuite) TestETHSignatureChallenge_EmptyChallenge() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Test 1: Empty signer address
	suite.Run("EmptySigner", func() {
		collectionsToCreate := GetCollectionsToCreate()
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
			{
				Signer:             "", // Empty signer
				ChallengeTrackerId: "test-challenge-1",
			},
		}
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

		err := CreateCollections(suite, wctx, collectionsToCreate)
		suite.Require().NoError(err)

		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     "test-nonce",
							Signature: "0x1234",
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().Error(err, "Empty signer should cause error")
	})
}

// TestETHSignatureChallenge_WithOtherCriteria tests integration with other approval criteria
func (suite *TestSuite) TestETHSignatureChallenge_WithOtherCriteria() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Generate test private key and address
	_, signerAddress, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)

	// Test: ETH signature challenge combined with Merkle challenge
	suite.Run("WithMerkleChallenge", func() {
		// Create Merkle tree
		aliceLeaf := "-" + alice + "-1-0-0"
		leafs := [][]byte{[]byte(aliceLeaf)}
		leafHashes := make([][]byte, len(leafs))
		for i, leaf := range leafs {
			initialHash := sha256.Sum256(leaf)
			leafHashes[i] = initialHash[:]
		}
		rootHash := hex.EncodeToString(leafHashes[0])

		collectionsToCreate := GetCollectionsToCreate()
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.MerkleChallenges = []*types.MerkleChallenge{
			{
				Root:                rootHash,
				ExpectedProofLength: sdkmath.NewUint(0),
				MaxUsesPerLeaf:      sdkmath.NewUint(1),
				ChallengeTrackerId:  "merkle-challenge-1",
			},
		}
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
			{
				Signer:             signerAddress,
				ChallengeTrackerId: "test-challenge-1",
			},
		}
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

		// Add mint approval so we can mint badges to bob
		collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
			ToListId:          "AllWithoutMint",
			FromListId:        "Mint",
			InitiatedByListId: "AllWithoutMint",
			TransferTimes:     GetFullUintRanges(),
			TokenIds:          GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "mint-test",
			ApprovalCriteria: &types.ApprovalCriteria{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(1000),
					AmountTrackerId:        "mint-test-tracker",
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
					AmountTrackerId:              "mint-test-tracker",
				},
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
			},
		}}, collectionsToCreate[0].CollectionApprovals...)

		err = CreateCollections(suite, wctx, collectionsToCreate)
		suite.Require().NoError(err)

		// Mint badges to bob so he can transfer them
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
							TokenIds:       GetTopHalfUintRanges(),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
						{
							ApprovalId:      "mint-test",
							ApprovalLevel:   "collection",
							ApproverAddress: "",
							Version:         sdkmath.NewUint(0),
						},
					},
				},
			},
		})
		suite.Require().NoError(err, "Error minting badges to bob")

		privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
		nonce := "test-nonce-combined"
		signature, err := generateETHSignature(
			nonce,
			alice,
			"1",
			"",
			"collection",
			"test",
			"test-challenge-1",
			privateKeyHex,
		)
		suite.Require().NoError(err)

		// Transfer with both Merkle proof and ETH signature
		msg := &types.MsgTransferTokens{
			Creator:      alice,
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
					MerkleProofs: []*types.MerkleProof{
						{
							Leaf:  aliceLeaf,
							Aunts: []*types.MerklePathItem{},
						},
					},
					EthSignatureProofs: []*types.ETHSignatureProof{
						{
							Nonce:     nonce,
							Signature: signature,
						},
					},
				},
			},
		}

		err = TransferTokens(suite, wctx, msg)
		suite.Require().NoError(err, "Both challenges satisfied should succeed")
	})
}

// TestETHSignatureChallenge_TrackerQuery tests the signature tracker query functionality
func (suite *TestSuite) TestETHSignatureChallenge_TrackerQuery() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Generate test private key and address
	_, signerAddress, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)

	// Create collection with ETH signature challenge
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.EthSignatureChallenges = []*types.ETHSignatureChallenge{
		{
			Signer:             signerAddress,
			ChallengeTrackerId: "test-challenge-1",
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval so we can mint badges to bob
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob so he can transfer them
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
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "mint-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().NoError(err, "Error minting badges to bob")

	// Generate and use signature
	nonce := "test-nonce-query"
	privateKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	signature, err := generateETHSignature(
		nonce,
		alice,
		"1",
		"",
		"collection",
		"test",
		"test-challenge-1",
		privateKeyHex,
	)
	suite.Require().NoError(err)

	msg := &types.MsgTransferTokens{
		Creator:      alice,
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
				EthSignatureProofs: []*types.ETHSignatureProof{
					{
						Nonce:     nonce,
						Signature: signature,
					},
				},
			},
		},
	}

	err = TransferTokens(suite, wctx, msg)
	suite.Require().NoError(err)

	// Query the signature tracker
	response, err := suite.app.BadgesKeeper.GetETHSignatureTracker(
		wctx,
		&types.QueryGetETHSignatureTrackerRequest{
			CollectionId:       "1",
			ApproverAddress:    "",
			ApprovalLevel:      "collection",
			ApprovalId:         "test",
			ChallengeTrackerId: "test-challenge-1",
			Signature:          signature,
		},
	)
	suite.Require().NoError(err)
	suite.Require().Equal("1", response.NumUsed, "Signature should be used once")
}
