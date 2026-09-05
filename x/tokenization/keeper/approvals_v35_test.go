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
