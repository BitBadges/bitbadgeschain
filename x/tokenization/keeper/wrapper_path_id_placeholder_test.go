package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func placeholderTestCollection(denom string, allowOverride bool) ([]*types.MsgNewCollection, string) {
	colA := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	colA[0].TokensToCreate = []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}}, OwnershipTimes: GetFullUintRanges()}}
	colA[0].Transfers = nil
	colA[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{{
		Denom: denom,
		Conversion: &types.ConversionWithoutDenom{
			SideA: &types.ConversionSideA{Amount: sdkmath.NewUint(1)},
			SideB: []*types.Balance{{Amount: sdkmath.NewUint(1), OwnershipTimes: GetFullUintRanges(), TokenIds: GetOneUintRange()}},
		},
		Symbol:                         "FUNG",
		DenomUnits:                     []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "fung", IsDefaultDisplay: true}},
		AllowOverrideWithAnyValidToken: allowOverride,
	}}
	wrapperAddr := keeper.MustGenerateWrapperPathAddress(denom)
	wrap := &types.CollectionApproval{
		ApprovalId: "wrap", TransferTimes: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges(), TokenIds: GetFullUintRanges(),
		FromListId: "AllWithoutMint", ToListId: wrapperAddr, InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{AllowSpecialWrapping: true, MustPrioritize: true},
	}
	unwrap := &types.CollectionApproval{
		ApprovalId: "unwrap", TransferTimes: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges(), TokenIds: GetFullUintRanges(),
		FromListId: wrapperAddr, ToListId: "AllWithoutMint", InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{AllowSpecialWrapping: true, MustPrioritize: true},
	}
	colA[0].CollectionApprovals = []*types.CollectionApproval{rollbackTestMintApproval(), wrap, unwrap}
	return colA, wrapperAddr
}

func placeholderPrio(id string) []*types.ApprovalIdentifierDetails {
	return []*types.ApprovalIdentifierDetails{{ApprovalId: id, ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}}
}

func placeholderBal(amount uint64, id uint64) []*types.Balance {
	return []*types.Balance{{Amount: sdkmath.NewUint(amount), TokenIds: []*types.UintRange{{Start: sdkmath.NewUint(id), End: sdkmath.NewUint(id)}}, OwnershipTimes: GetFullUintRanges()}}
}

// allowOverrideWithAnyValidToken and the {id} placeholder must be configured together.
func (suite *TestSuite) TestWrapperPathOverrideRequiresIdPlaceholder() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	colA, _ := placeholderTestCollection("fung", true)
	suite.Require().Error(CreateCollections(suite, wctx, colA), "override flag without {id} must be rejected")

	colB, _ := placeholderTestCollection("fung{id}", false)
	suite.Require().Error(CreateCollections(suite, wctx, colB), "{id} without the override flag must be rejected")

	colC, _ := placeholderTestCollection("fung{id}", true)
	suite.Require().Nil(CreateCollections(suite, wctx, colC))

	colD, _ := placeholderTestCollection("fung{id}{id}", true)
	suite.Require().Error(CreateCollections(suite, wctx, colD), "only one {id} placeholder is allowed")

	colE, _ := placeholderTestCollection("ab{c}", false)
	suite.Require().Error(CreateCollections(suite, wctx, colE), "stray braces must be rejected")
}

// Each token id wraps into its own bank denom, so unwrapping returns exactly the deposited id.
func (suite *TestSuite) TestWrapperPathUnwrapReturnsDepositedTokenId() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	colA, wrapperAddr := placeholderTestCollection("fung{id}", true)
	suite.Require().Nil(CreateCollections(suite, wctx, colA))

	transfer := func(creator, from, to string, bal []*types.Balance, approval string) error {
		return TransferTokens(suite, wctx, &types.MsgTransferTokens{Creator: creator, CollectionId: sdkmath.NewUint(1), Transfers: []*types.Transfer{{From: from, ToAddresses: []string{to}, Balances: bal, PrioritizedApprovals: placeholderPrio(approval)}}})
	}

	// alice mints and wraps id 1; bob mints and wraps id 2
	suite.Require().Nil(transfer(alice, "Mint", alice, placeholderBal(1, 1), "mint-test"))
	suite.Require().Nil(transfer(alice, alice, wrapperAddr, placeholderBal(1, 1), "wrap"))
	suite.Require().Nil(transfer(bob, "Mint", bob, placeholderBal(1, 2), "mint-test"))
	suite.Require().Nil(transfer(bob, bob, wrapperAddr, placeholderBal(1, 2), "wrap"))

	bobAcc := sdk.MustAccAddressFromBech32(bob)
	aliceAcc := sdk.MustAccAddressFromBech32(alice)
	suite.Require().Equal(sdkmath.NewInt(1), suite.app.BankKeeper.GetBalance(suite.ctx, aliceAcc, keeper.WrappedDenomPrefix+"1:fung1").Amount)
	suite.Require().Equal(sdkmath.NewInt(1), suite.app.BankKeeper.GetBalance(suite.ctx, bobAcc, keeper.WrappedDenomPrefix+"1:fung2").Amount)
	suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, bobAcc, keeper.WrappedDenomPrefix+"1:fung1").Amount.IsZero())

	// bob cannot redeem alice's id 1 with his id-2 coin
	suite.Require().Error(transfer(bob, wrapperAddr, bob, placeholderBal(1, 1), "unwrap"), "bob deposited id 2 and must not receive id 1")
	bobBal, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Len(bobBal.Balances, 0)

	// round trip for both depositors
	suite.Require().Nil(transfer(bob, wrapperAddr, bob, placeholderBal(1, 2), "unwrap"))
	suite.Require().Nil(transfer(alice, wrapperAddr, alice, placeholderBal(1, 1), "unwrap"))
	bobBal, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	AssertBalancesEqual(suite, placeholderBal(1, 2), bobBal.Balances)
	aliceBal, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	AssertBalancesEqual(suite, placeholderBal(1, 1), aliceBal.Balances)
	suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, bobAcc, keeper.WrappedDenomPrefix+"1:fung2").Amount.IsZero())
	suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, aliceAcc, keeper.WrappedDenomPrefix+"1:fung1").Amount.IsZero())
}

// A stored path that carries the override flag but no {id} (pre-v35 state) converts only
// the SideB token id; the flag no longer substitutes the transferred id.
func (suite *TestSuite) TestLegacyOverridePathWithoutPlaceholderUsesSideBTokenId() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	colA, wrapperAddr := placeholderTestCollection("fung", false)
	suite.Require().Nil(CreateCollections(suite, wctx, colA))

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	collection.CosmosCoinWrapperPaths[0].AllowOverrideWithAnyValidToken = true
	suite.Require().Nil(suite.app.TokenizationKeeper.SetCollectionInStore(suite.ctx, collection, true))

	transfer := func(creator, from, to string, bal []*types.Balance, approval string) error {
		return TransferTokens(suite, wctx, &types.MsgTransferTokens{Creator: creator, CollectionId: sdkmath.NewUint(1), Transfers: []*types.Transfer{{From: from, ToAddresses: []string{to}, Balances: bal, PrioritizedApprovals: placeholderPrio(approval)}}})
	}
	suite.Require().Nil(transfer(bob, "Mint", bob, placeholderBal(1, 1), "mint-test"))
	suite.Require().Nil(transfer(bob, "Mint", bob, placeholderBal(1, 2), "mint-test"))

	suite.Require().Error(transfer(bob, bob, wrapperAddr, placeholderBal(1, 2), "wrap"), "id 2 is not the SideB id and must not wrap")
	suite.Require().Nil(transfer(bob, bob, wrapperAddr, placeholderBal(1, 1), "wrap"))
	bobAcc := sdk.MustAccAddressFromBech32(bob)
	suite.Require().Equal(sdkmath.NewInt(1), suite.app.BankKeeper.GetBalance(suite.ctx, bobAcc, keeper.WrappedDenomPrefix+"1:fung").Amount)
}
