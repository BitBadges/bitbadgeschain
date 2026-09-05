package keeper_test

import (
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
