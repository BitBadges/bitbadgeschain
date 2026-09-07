package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// A backed collection always gets the forbidden Mint permission first, even when the creator
// pre-seeds a permission with the same list ids but permitted times.
func (suite *TestSuite) TestBackedCollectionForbidsMintApprovalsEvenWhenPreSeededPermitted() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filtered := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filtered = append(filtered, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filtered
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{Amount: sdkmath.NewUint(1), Denom: "ibc/permtest"},
				SideB: []*types.Balance{{Amount: sdkmath.NewUint(1), OwnershipTimes: GetFullUintRanges(), TokenIds: GetOneUintRange()}},
			},
		},
	}
	collectionsToCreate[0].Permissions.CanUpdateCollectionApprovals = []*types.CollectionApprovalPermission{{
		FromListId: "Mint", ToListId: "All", InitiatedByListId: "All", ApprovalId: "All",
		TransferTimes: GetFullUintRanges(), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges(),
		PermanentlyPermittedTimes: GetFullUintRanges(),
	}}
	suite.Require().Nil(CreateCollections(suite, wctx, collectionsToCreate))

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	first := collection.CollectionPermissions.CanUpdateCollectionApprovals[0]
	suite.Require().Equal("Mint", first.FromListId)
	suite.Require().Len(first.PermanentlyPermittedTimes, 0)
	AssertUintRangesEqual(suite, GetFullUintRanges(), first.PermanentlyForbiddenTimes)

	// And a Mint approval can indeed not be added afterwards
	approvals := append([]*types.CollectionApproval{}, collection.CollectionApprovals...)
	approvals = append(approvals, rollbackTestMintApproval())
	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{Creator: bob, CollectionId: sdkmath.NewUint(1), CollectionApprovals: approvals})
	suite.Require().Error(err)
}
