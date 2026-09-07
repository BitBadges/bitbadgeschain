package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// A collection cannot be deleted while a wrapper path still escrows tokens for
// outstanding bank coins.
func (suite *TestSuite) TestDeleteCollectionRefusedWhileWrapperSupplyOutstanding() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	colA, wrapperAddr := placeholderTestCollection("fung", false)
	suite.Require().Nil(CreateCollections(suite, wctx, colA))

	transfer := func(from, to string, bal []*types.Balance, approval string) error {
		return TransferTokens(suite, wctx, &types.MsgTransferTokens{Creator: bob, CollectionId: sdkmath.NewUint(1), Transfers: []*types.Transfer{{From: from, ToAddresses: []string{to}, Balances: bal, PrioritizedApprovals: placeholderPrio(approval)}}})
	}
	suite.Require().Nil(transfer("Mint", bob, placeholderBal(1, 1), "mint-test"))
	suite.Require().Nil(transfer(bob, wrapperAddr, placeholderBal(1, 1), "wrap"))

	deleteMsg := &types.MsgDeleteCollection{Creator: bob, CollectionId: sdkmath.NewUint(1)}
	suite.Require().Error(DeleteCollection(suite, wctx, deleteMsg), "wrapped supply is outstanding")
	_, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "collection must still exist")

	suite.Require().Nil(transfer(wrapperAddr, bob, placeholderBal(1, 1), "unwrap"))
	suite.Require().Nil(DeleteCollection(suite, wctx, deleteMsg))
	_, err = GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Error(err)

	// Collection stats are purged with the collection
	_, found := suite.app.TokenizationKeeper.GetCollectionStatsFromStore(suite.ctx, sdkmath.NewUint(1))
	suite.Require().False(found, "collection stats must be purged on delete")
}

// A backed collection cannot be deleted while unbacked tokens are in circulation, since the
// escrowed bank coins are only redeemable through the collection.
func (suite *TestSuite) TestDeleteCollectionRefusedWhileBackedSupplyCirculating() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{Amount: sdkmath.NewUint(1), Denom: "ibc/deletetest"},
				SideB: []*types.Balance{{Amount: sdkmath.NewUint(1), OwnershipTimes: GetFullUintRanges(), TokenIds: GetOneUintRange()}},
			},
		},
	}
	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "delete-test")
	suite.Require().Nil(err)
	backedPath := collection.Invariants.CosmosCoinBackedPath

	bobAccAddr := sdk.MustAccAddressFromBech32(bob)
	coin := sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))}
	suite.Require().Nil(suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coin))
	suite.Require().Nil(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, coin))

	transfer := func(from, to string) error {
		return TransferTokens(suite, wctx, &types.MsgTransferTokens{
			Creator: bob, CollectionId: sdkmath.NewUint(1),
			Transfers: []*types.Transfer{{From: from, ToAddresses: []string{to}, Balances: placeholderBal(1, 1), PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection)}},
		})
	}
	suite.Require().Nil(transfer(backedPath.Address, bob), "unback")

	deleteMsg := &types.MsgDeleteCollection{Creator: bob, CollectionId: sdkmath.NewUint(1)}
	suite.Require().Error(DeleteCollection(suite, wctx, deleteMsg), "backed supply is circulating")

	suite.Require().Nil(transfer(bob, backedPath.Address), "back")
	suite.Require().Nil(DeleteCollection(suite, wctx, deleteMsg))
}
