package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// New supply can only enter the collection for token ids inside ValidTokenIds.
func (suite *TestSuite) TestMintOutsideValidTokenIdsIsRejected() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	colA := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	colA[0].TokensToCreate = []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}}
	colA[0].CollectionApprovals = []*types.CollectionApproval{rollbackTestMintApproval()}
	colA[0].Transfers = nil
	suite.Require().Nil(CreateCollections(suite, wctx, colA))
	col, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	AssertUintRangesEqual(suite, col.ValidTokenIds, GetOneUintRange())

	mint := func(id uint64) error {
		return TransferTokens(suite, wctx, &types.MsgTransferTokens{
			Creator:      bob,
			CollectionId: sdkmath.NewUint(1),
			Transfers: []*types.Transfer{{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances:    []*types.Balance{{Amount: sdkmath.NewUint(5), TokenIds: []*types.UintRange{{Start: sdkmath.NewUint(id), End: sdkmath.NewUint(id)}}, OwnershipTimes: GetFullUintRanges()}},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{ApprovalId: "mint-test", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
				},
			}},
		})
	}

	suite.Require().Error(mint(500), "minting a token id outside ValidTokenIds must fail")
	suite.Require().Nil(mint(1), "minting a valid token id must still work")

	total, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Total")
	suite.Require().Nil(err)
	suite.Require().Nil(types.ValidateBalancesWithinValidTokenIds(total.Balances, col.ValidTokenIds))
}

// Shrinking ValidTokenIds is gated by CanUpdateValidTokenIds for the removed ids, exactly like growing it.
func (suite *TestSuite) TestShrinkValidTokenIdsRequiresPermission() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].TokensToCreate = []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		OwnershipTimes: GetFullUintRanges(),
	}}
	suite.Require().Nil(CreateCollections(suite, wctx, collectionsToCreate))

	// Lock the id set for every id
	suite.Require().Nil(UpdateCollection(suite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:                     bob,
		CollectionId:                sdkmath.NewUint(1),
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{{
				TokenIds:                  GetFullUintRanges(),
				PermanentlyForbiddenTimes: GetFullUintRanges(),
			}},
		},
	}))

	err := UpdateCollection(suite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		UpdateValidTokenIds: true,
		ValidTokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
	})
	suite.Require().Error(err, "shrinking a locked id set must fail")

	col, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	AssertUintRangesEqual(suite, col.ValidTokenIds, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}})

	// Re-submitting the identical set is a no-op and stays allowed
	suite.Require().Nil(UpdateCollection(suite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		UpdateValidTokenIds: true,
		ValidTokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
	}))
}

func (suite *TestSuite) TestShrinkValidTokenIdsWithPermission() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].TokensToCreate = []*types.Balance{{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		OwnershipTimes: GetFullUintRanges(),
	}}
	suite.Require().Nil(CreateCollections(suite, wctx, collectionsToCreate))

	suite.Require().Nil(UpdateCollection(suite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:                     bob,
		CollectionId:                sdkmath.NewUint(1),
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{{
				TokenIds:                  GetFullUintRanges(),
				PermanentlyPermittedTimes: GetFullUintRanges(),
			}},
		},
	}))

	suite.Require().Nil(UpdateCollection(suite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		UpdateValidTokenIds: true,
		ValidTokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}},
	}))
	col, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	AssertUintRangesEqual(suite, col.ValidTokenIds, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)}})
}
