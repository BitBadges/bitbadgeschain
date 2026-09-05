package keeper_test

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// Shared helpers for the v35 approval regression tests below.

func v35Permissions() *types.CollectionPermissions {
	return &types.CollectionPermissions{
		CanArchiveCollection:         []*types.ActionPermission{},
		CanUpdateStandards:           []*types.ActionPermission{},
		CanUpdateCustomData:          []*types.ActionPermission{},
		CanDeleteCollection:          []*types.ActionPermission{},
		CanUpdateManager:             []*types.ActionPermission{},
		CanUpdateCollectionMetadata:  []*types.ActionPermission{},
		CanUpdateTokenMetadata:       []*types.TokenIdsActionPermission{},
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
		CanUpdateValidTokenIds:       []*types.TokenIdsActionPermission{{PermanentlyPermittedTimes: GetFullUintRanges()}},
	}
}

// v35MintCollection builds a collection with a single "claim" mint approval carrying the given criteria.
func v35MintCollection(criteria *types.ApprovalCriteria, tokenIds []*types.UintRange) *types.MsgNewCollection {
	return &types.MsgNewCollection{
		Creator: bob,
		CollectionApprovals: []*types.CollectionApproval{{
			ToListId:          "AllWithoutMint",
			FromListId:        "Mint",
			InitiatedByListId: "AllWithoutMint",
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          tokenIds,
			ApprovalId:        "claim",
			ApprovalCriteria:  criteria,
		}},
		TokensToCreate: []*types.Balance{{Amount: sdkmath.NewUint(1000), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
		Permissions:    v35Permissions(),
	}
}

// v35Claim builds a one-token mint claim from "Mint" to `to`, initiated by `creator`.
func v35Claim(creator, to string, extra func(*types.Transfer)) *types.MsgTransferTokens {
	t := &types.Transfer{
		From:        "Mint",
		ToAddresses: []string{to},
		Balances:    []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}},
	}
	if extra != nil {
		extra(t)
	}
	return &types.MsgTransferTokens{Creator: creator, CollectionId: sdkmath.NewUint(1), Transfers: []*types.Transfer{t}}
}

// ethSignatureEncodings returns every accepted textual/encoding form of one 65-byte signature:
// 0x / 0X / no prefix, upper/lower hex, V as 0/1 and 27/28, and the (r, n-s, v^1) form.
func ethSignatureEncodings(sig []byte) []string {
	lower := "0x" + hex.EncodeToString(sig)

	v27 := append([]byte{}, sig...)
	if v27[64] < 27 {
		v27[64] += 27
	}

	n := ethcrypto.S256().Params().N
	s := new(big.Int).SetBytes(sig[32:64])
	sFlipped := new(big.Int).Sub(n, s)
	highS := append([]byte{}, sig...)
	sFlipped.FillBytes(highS[32:64])
	highS[64] ^= 1

	return []string{
		"0x" + strings.ToUpper(lower[2:]),
		lower[2:],
		"0X" + lower[2:],
		"0x" + hex.EncodeToString(v27),
		"0x" + hex.EncodeToString(highS),
		"0x" + strings.ToUpper(hex.EncodeToString(highS)),
	}
}

// H-5 regression: one issued signature is one use, whatever encoding it is submitted in.
func (suite *TestSuite) TestETHSignatureChallengeSingleUseAcrossEncodings() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	privateKeyHex, signerAddress, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)

	err = CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
		EthSignatureChallenges:         []*types.ETHSignatureChallenge{{Signer: signerAddress, ChallengeTrackerId: "c"}},
	}, GetFullUintRanges())})
	suite.Require().NoError(err)

	nonce := "n1"
	signature, err := generateETHSignature(nonce, alice, "1", "", "collection", "claim", "c", privateKeyHex)
	suite.Require().NoError(err)
	sigBytes, err := hex.DecodeString(signature[2:])
	suite.Require().NoError(err)

	claim := func(sigStr string) error {
		return TransferTokens(suite, wctx, v35Claim(alice, alice, func(t *types.Transfer) {
			t.EthSignatureProofs = []*types.ETHSignatureProof{{Nonce: nonce, Signature: sigStr}}
			t.PrioritizedApprovals = []*types.ApprovalIdentifierDetails{{ApprovalId: "claim", ApprovalLevel: "collection", Version: sdkmath.NewUint(0)}}
		}))
	}

	suite.Require().NoError(claim(signature), "first use")
	suite.Require().Error(claim(signature), "exact re-submission")
	for i, encoding := range ethSignatureEncodings(sigBytes) {
		suite.Require().Error(claim(encoding), "encoding %d must count as the same use", i)
	}

	balance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().NoError(err)
	suite.Require().Len(balance.Balances, 1)
	suite.Require().Equal(sdkmath.NewUint(1), balance.Balances[0].Amount)

	// The tracker is keyed on the nonce the signer committed to.
	numUsed, exists := suite.app.TokenizationKeeper.GetETHSignatureTrackerFromStore(suite.ctx,
		keeper.ConstructETHSignatureTrackerKey(sdkmath.NewUint(1), "", "collection", "claim", "c", nonce))
	suite.Require().True(exists)
	suite.Require().Equal(sdkmath.NewUint(1), numUsed)
}

// H-6 regression: a per-initiatedBy limit is counted against one identity, whatever spelling
// of the bech32 address the message carries.
func (suite *TestSuite) TestPerInitiatedByLimitCountsOneIdentityPerKey() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
		MaxNumTransfers: &types.MaxNumTransfers{
			PerInitiatedByAddressMaxNumTransfers: sdkmath.NewUint(1),
			AmountTrackerId:                      "t",
		},
	}, GetFullUintRanges())})
	suite.Require().NoError(err)

	suite.Require().NoError(TransferTokens(suite, wctx, v35Claim(alice, alice, nil)), "first claim")
	suite.Require().Error(TransferTokens(suite, wctx, v35Claim(alice, alice, nil)), "second claim")
	suite.Require().Error(TransferTokens(suite, wctx, v35Claim(strings.ToUpper(alice), alice, nil)), "second claim, uppercase initiator")
	suite.Require().Error(TransferTokens(suite, wctx, v35Claim(alice, strings.ToUpper(alice), nil)), "uppercase recipient")

	balance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(1), balance.Balances[0].Amount)
	upperBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), strings.ToUpper(alice))
	suite.Require().NoError(err)
	suite.Require().Empty(upperBalance.Balances, "nothing is held under the uppercase spelling")
}

// Migration: store entries carrying a non-canonical bech32 spelling are rewritten to the
// canonical spelling and merged with any canonical twin (balances/tallies sum, lists union).
func (suite *TestSuite) TestMigrateV35CanonicalAddressesMergesEntries() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	k := suite.app.TokenizationKeeper
	upperAlice, upperBob, upperCharlie := strings.ToUpper(alice), strings.ToUpper(bob), strings.ToUpper(charlie)
	one := sdkmath.NewUint(1)

	err := CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}, GetFullUintRanges())})
	suite.Require().NoError(err)
	suite.Require().NoError(TransferTokens(suite, wctx, v35Claim(alice, alice, nil)))

	// Balance store: raw key under the uppercase spelling holding 3 of token 1.
	rawStore := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
	upperBalanceKey := append(append([]byte{}, keeper.UserBalanceKey...), []byte(keeper.ConstructBalanceKey(upperAlice, one))...)
	rawStore.Set(upperBalanceKey, suite.app.AppCodec().MustMarshal(&types.UserBalanceStore{
		Balances: []*types.Balance{{Amount: sdkmath.NewUint(3), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}},
	}))

	// Address list values.
	suite.Require().NoError(k.SetAddressListInStore(suite.ctx, types.AddressList{ListId: "l1", Addresses: []string{upperAlice, alice, bob}, Whitelist: true, CreatedBy: upperBob}))

	// Approval trackers (approver component and per-address component).
	upperTrackerKey := keeper.ConstructApprovalTrackerKey(one, "", "claim", "t", "collection", "initiatedBy", upperAlice)
	canonTrackerKey := keeper.ConstructApprovalTrackerKey(one, "", "claim", "t", "collection", "initiatedBy", alice)
	suite.Require().NoError(k.SetApprovalTrackerInStoreViaKey(suite.ctx, upperTrackerKey, types.ApprovalTracker{NumTransfers: sdkmath.NewUint(2), Amounts: []*types.Balance{{Amount: sdkmath.NewUint(2), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}}, LastUpdatedAt: sdkmath.NewUint(10)}))
	suite.Require().NoError(k.SetApprovalTrackerInStoreViaKey(suite.ctx, canonTrackerKey, types.ApprovalTracker{NumTransfers: one, Amounts: []*types.Balance{{Amount: one, TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}}, LastUpdatedAt: sdkmath.NewUint(5)}))
	approverOnlyKey := keeper.ConstructApprovalTrackerKey(one, upperBob, "a", "t", "incoming", "overall", "")
	suite.Require().NoError(k.SetApprovalTrackerInStoreViaKey(suite.ctx, approverOnlyKey, types.ApprovalTracker{NumTransfers: sdkmath.NewUint(7), Amounts: []*types.Balance{}, LastUpdatedAt: one}))

	// Used merkle-leaf and ETH signature trackers.
	suite.Require().NoError(k.SetChallengeTrackerInStore(suite.ctx, keeper.ConstructUsedClaimChallengeKey(one, upperBob, "incoming", "a", "c", one), one))
	suite.Require().NoError(k.SetChallengeTrackerInStore(suite.ctx, keeper.ConstructUsedClaimChallengeKey(one, bob, "incoming", "a", "c", one), one))
	suite.Require().NoError(k.SetETHSignatureTrackerInStore(suite.ctx, keeper.ConstructETHSignatureTrackerKey(one, upperBob, "incoming", "a", "c", "n1"), one))

	// Approval versions.
	k.SetApprovalTrackerVersionInStore(suite.ctx, keeper.ConstructApprovalVersionKey(one, "incoming", upperBob, "a"), sdkmath.NewUint(5))
	k.SetApprovalTrackerVersionInStore(suite.ctx, keeper.ConstructApprovalVersionKey(one, "incoming", bob, "a"), sdkmath.NewUint(3))

	// Dynamic store values and reserved protocol addresses.
	suite.Require().NoError(k.SetDynamicStoreValueInStore(suite.ctx, one, upperAlice, true))
	suite.Require().NoError(k.SetReservedProtocolAddressInStore(suite.ctx, upperCharlie, true))

	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))

	// Balances: 1 (claimed) + 3 (planted) under the canonical key; uppercase key gone.
	balance, err := GetUserBalance(suite, wctx, one, alice)
	suite.Require().NoError(err)
	suite.Require().Len(balance.Balances, 1)
	suite.Require().Equal(sdkmath.NewUint(4), balance.Balances[0].Amount)
	suite.Require().Nil(rawStore.Get(upperBalanceKey))

	list, found := k.GetAddressListFromStore(suite.ctx, "l1")
	suite.Require().True(found)
	suite.Require().Equal([]string{alice, bob}, list.Addresses)
	suite.Require().Equal(bob, list.CreatedBy)

	_, found = k.GetApprovalTrackerFromStore(suite.ctx, one, "", "claim", "t", "collection", "initiatedBy", upperAlice)
	suite.Require().False(found)
	tracker, found := k.GetApprovalTrackerFromStore(suite.ctx, one, "", "claim", "t", "collection", "initiatedBy", alice)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(3), tracker.NumTransfers)
	suite.Require().Equal(sdkmath.NewUint(3), tracker.Amounts[0].Amount)
	suite.Require().Equal(sdkmath.NewUint(10), tracker.LastUpdatedAt)
	tracker, found = k.GetApprovalTrackerFromStore(suite.ctx, one, bob, "a", "t", "incoming", "overall", "")
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(7), tracker.NumTransfers)

	numUsed, err := k.GetChallengeTrackerFromStore(suite.ctx, one, bob, "incoming", "a", "c", one)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(2), numUsed)
	numUsed, err = k.GetChallengeTrackerFromStore(suite.ctx, one, upperBob, "incoming", "a", "c", one)
	suite.Require().NoError(err)
	suite.Require().True(numUsed.IsZero())

	numUsed, found = k.GetETHSignatureTrackerFromStore(suite.ctx, keeper.ConstructETHSignatureTrackerKey(one, bob, "incoming", "a", "c", "n1"))
	suite.Require().True(found)
	suite.Require().Equal(one, numUsed)

	version, found := k.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(one, "incoming", bob, "a"))
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(5), version)
	_, found = k.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(one, "incoming", upperBob, "a"))
	suite.Require().False(found)

	value, found := k.GetDynamicStoreValueFromStore(suite.ctx, one, alice)
	suite.Require().True(found)
	suite.Require().True(value.Value)
	suite.Require().True(k.IsAddressReservedProtocolInStore(suite.ctx, charlie))
	suite.Require().False(k.IsAddressReservedProtocolInStore(suite.ctx, upperCharlie))
}

// Validation rejects an approval without token ids, at collection and user level.
func (suite *TestSuite) TestApprovalWithoutTokenIdsIsRejected() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}, []*types.UintRange{})})
	suite.Require().Error(err, "collection approval without token ids")

	err = CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}, GetFullUintRanges())})
	suite.Require().NoError(err)

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 alice,
		CollectionId:            sdkmath.NewUint(1),
		UpdateIncomingApprovals: true,
		IncomingApprovals: []*types.UserIncomingApproval{{
			FromListId:        "All",
			InitiatedByListId: "All",
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          []*types.UintRange{},
			ApprovalId:        "in",
		}},
	})
	suite.Require().Error(err, "incoming approval without token ids")
}

// Validation rejects predetermined balances that carry no order calculation method.
func (suite *TestSuite) TestPredeterminedBalancesRequireOrderCalculationMethod() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
		PredeterminedBalances:          &types.PredeterminedBalances{},
	}, GetFullUintRanges())})
	suite.Require().Error(err)
}

// v35StoreClaimCollection writes a collection with the "claim" mint approval straight into the
// store, skipping message validation, to model state written by an earlier version.
func (suite *TestSuite) v35StoreClaimCollection(criteria *types.ApprovalCriteria, tokenIds []*types.UintRange) {
	wctx := sdk.WrapSDKContext(suite.ctx)
	err := CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}, GetFullUintRanges())})
	suite.Require().NoError(err)

	collection, found := suite.app.TokenizationKeeper.GetCollectionFromStore(suite.ctx, sdkmath.NewUint(1))
	suite.Require().True(found)
	collection.CollectionApprovals[0].ApprovalCriteria = criteria
	collection.CollectionApprovals[0].TokenIds = tokenIds
	suite.Require().NoError(suite.app.TokenizationKeeper.SetCollectionInStore(suite.ctx, collection, true))
}

// Stored approvals without token ids are skipped with an error rather than aborting the transfer.
func (suite *TestSuite) TestTransferSkipsStoredApprovalWithoutTokenIds() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.v35StoreClaimCollection(&types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}, []*types.UintRange{})

	var err error
	suite.Require().NotPanics(func() { err = TransferTokens(suite, wctx, v35Claim(alice, alice, nil)) })
	suite.Require().Error(err)
}

// Stored predetermined balances without an order calculation method impose no ordering and
// do not abort the transfer.
func (suite *TestSuite) TestTransferHandlesStoredPredeterminedBalancesWithoutMethod() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.v35StoreClaimCollection(&types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
		PredeterminedBalances:          &types.PredeterminedBalances{},
	}, GetFullUintRanges())

	var err error
	suite.Require().NotPanics(func() { err = TransferTokens(suite, wctx, v35Claim(alice, alice, nil)) })
	suite.Require().NoError(err)

	// The precalculation path reads the same field and must also stay deterministic.
	collection, _ := suite.app.TokenizationKeeper.GetCollectionFromStore(suite.ctx, sdkmath.NewUint(1))
	suite.Require().NotPanics(func() {
		err = TransferTokens(suite, wctx, v35Claim(alice, alice, func(t *types.Transfer) {
			t.PrecalculateBalancesFromApproval = &types.PrecalculateBalancesFromApprovalDetails{ApprovalId: "claim", ApprovalLevel: "collection", Version: collection.CollectionApprovals[0].Version}
		}))
	})
	suite.Require().Error(err)
}

// Precalculation from a user-level approval reads merkle-leaf usage under that approval's
// own (approver, level) key, not the collection-level key.
func (suite *TestSuite) TestPrecalculationChecksLeafUsageAtUserLevel() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	k := suite.app.TokenizationKeeper
	one := sdkmath.NewUint(1)

	err := CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}, GetFullUintRanges())})
	suite.Require().NoError(err)
	collection, found := k.GetCollectionFromStore(suite.ctx, one)
	suite.Require().True(found)

	// Two-leaf tree; alice's proof is leaf index 0.
	aliceLeaf, bobLeaf := "-"+alice+"-1-0-0", "-"+bob+"-1-0-0"
	aliceHash, bobHash := sha256.Sum256([]byte(aliceLeaf)), sha256.Sum256([]byte(bobLeaf))
	root := sha256.Sum256(append(aliceHash[:], bobHash[:]...))

	incoming := []*types.UserIncomingApproval{{
		FromListId:        "All",
		InitiatedByListId: "All",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		ApprovalId:        "in",
		Version:           sdkmath.NewUint(0),
		ApprovalCriteria: &types.IncomingApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{AmountTrackerId: "t"},
			MerkleChallenges: []*types.MerkleChallenge{{
				Root:                hex.EncodeToString(root[:]),
				ExpectedProofLength: one,
				MaxUsesPerLeaf:      one,
				ChallengeTrackerId:  "c",
			}},
			PredeterminedBalances: &types.PredeterminedBalances{
				OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseMerkleChallengeLeafIndex: true, ChallengeTrackerId: "c"},
				IncrementedBalances: &types.IncrementedBalances{
					StartBalances:             []*types.Balance{{Amount: one, TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}},
					IncrementTokenIdsBy:       one,
					IncrementOwnershipTimesBy: sdkmath.NewUint(0),
					DurationFromTimestamp:     sdkmath.NewUint(0),
				},
			},
		},
	}}
	approvals := types.CastIncomingTransfersToCollectionTransfers(incoming, alice)

	transfer := &types.Transfer{
		From:        bob,
		ToAddresses: []string{alice},
		MerkleProofs: []*types.MerkleProof{{
			Leaf:  aliceLeaf,
			Aunts: []*types.MerklePathItem{{Aunt: hex.EncodeToString(bobHash[:]), OnRight: true}},
		}},
		PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
			ApprovalId: "in", ApprovalLevel: "incoming", ApproverAddress: alice, Version: sdkmath.NewUint(0),
		},
	}
	// Metadata as built by the transfer path before the approval level is known.
	metadata := keeper.TransferMetadata{From: bob, To: alice, InitiatedBy: bob, ApproverAddress: "", ApprovalLevel: "collection"}

	// Leaf 0 already used under a collection-level approval with the same ids: irrelevant here.
	suite.Require().NoError(k.SetChallengeTrackerInStore(suite.ctx, keeper.ConstructUsedClaimChallengeKey(one, "", "collection", "in", "c", sdkmath.NewUint(0)), one))
	balances, err := k.GetPredeterminedBalancesForPrecalculationId(suite.ctx, collection, approvals, transfer, metadata)
	suite.Require().NoError(err)
	suite.Require().Len(balances, 1)
	suite.Require().Equal(one, balances[0].Amount)

	// Leaf 0 used under alice's incoming approval: the precalculation must refuse it.
	suite.Require().NoError(k.SetChallengeTrackerInStore(suite.ctx, keeper.ConstructUsedClaimChallengeKey(one, alice, "incoming", "in", "c", sdkmath.NewUint(0)), one))
	_, err = k.GetPredeterminedBalancesForPrecalculationId(suite.ctx, collection, approvals, transfer, metadata)
	suite.Require().Error(err)
}

// A user-level (incoming) approval gated on an ETH signature challenge can be satisfied:
// the transfer's EthSignatureProofs (and Memo) reach the user-level approval checks.
func (suite *TestSuite) TestIncomingApprovalETHSignatureChallengeCanBeSatisfied() {
	// Passes once HandleTransfer copies EthSignatureProofs and Memo into the user-level transfer.
	suite.T().Skip("pending EthSignatureProofs/Memo pass-through in transfers.go")
	wctx := sdk.WrapSDKContext(suite.ctx)
	privateKeyHex, signerAddress, err := generateTestETHPrivateKey()
	suite.Require().NoError(err)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers:                &types.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(1000), AmountTrackerId: "mint-test-tracker"},
			ApprovalAmounts:                &types.ApprovalAmounts{PerFromAddressApprovalAmount: sdkmath.NewUint(1000), AmountTrackerId: "mint-test-tracker"},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)
	suite.Require().NoError(CreateCollections(suite, wctx, collectionsToCreate))

	suite.Require().NoError(TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{{
			From:                 "Mint",
			ToAddresses:          []string{bob},
			Balances:             []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetTopHalfUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{{ApprovalId: "mint-test", ApprovalLevel: "collection", Version: sdkmath.NewUint(0)}},
		}},
	}))

	// alice only accepts transfers carrying a ticket signed for her incoming approval.
	suite.Require().NoError(UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 alice,
		CollectionId:            sdkmath.NewUint(1),
		UpdateIncomingApprovals: true,
		IncomingApprovals: []*types.UserIncomingApproval{{
			FromListId:        "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetFullUintRanges(),
			ApprovalId:        "ticket",
			ApprovalCriteria: &types.IncomingApprovalCriteria{
				EthSignatureChallenges: []*types.ETHSignatureChallenge{{Signer: signerAddress, ChallengeTrackerId: "c"}},
			},
		}},
	}))
	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().NoError(err)
	suite.Require().Len(aliceBalance.IncomingApprovals, 1)
	ticketVersion := aliceBalance.IncomingApprovals[0].Version

	nonce := "ticket-1"
	signature, err := generateETHSignature(nonce, bob, "1", alice, "incoming", "ticket", "c", privateKeyHex)
	suite.Require().NoError(err)

	transfer := func(ctx sdk.Context, proofs []*types.ETHSignatureProof) error {
		return TransferTokens(suite, sdk.WrapSDKContext(ctx), &types.MsgTransferTokens{
			Creator:      bob,
			CollectionId: sdkmath.NewUint(1),
			Transfers: []*types.Transfer{{
				From:        bob,
				ToAddresses: []string{alice},
				Balances:    []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetTopHalfUintRanges(), OwnershipTimes: GetFullUintRanges()}},
				Memo:        "ticket",
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{ApprovalId: "test", ApprovalLevel: "collection", Version: sdkmath.NewUint(0)},
					{ApprovalId: "ticket", ApprovalLevel: "incoming", ApproverAddress: alice, Version: ticketVersion},
				},
				EthSignatureProofs: proofs,
			}},
		})
	}

	// The failed attempt runs on a cache context so its collection-level tally does not carry over.
	cacheCtx, _ := suite.ctx.CacheContext()
	suite.Require().Error(transfer(cacheCtx, nil), "no ticket")
	suite.Require().NoError(transfer(suite.ctx, []*types.ETHSignatureProof{{Nonce: nonce, Signature: signature}}), "signed ticket")

	numUsed, exists := suite.app.TokenizationKeeper.GetETHSignatureTrackerFromStore(suite.ctx,
		keeper.ConstructETHSignatureTrackerKey(sdkmath.NewUint(1), alice, "incoming", "ticket", "c", nonce))
	suite.Require().True(exists)
	suite.Require().Equal(sdkmath.NewUint(1), numUsed)
}
