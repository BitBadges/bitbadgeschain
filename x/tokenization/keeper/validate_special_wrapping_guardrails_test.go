package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// createWrapperCollection creates a collection with CosmosCoinWrapperPaths and returns it.
// The collection has a transferable approval but NO allowSpecialWrapping approvals,
// so tests can add their own with specific configurations.
func (suite *TestSuite) createWrapperCollection() *types.TokenCollection {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "testcoin",
			Conversion: &types.ConversionWithoutDenom{
				SideA: &types.ConversionSideA{
					Amount: sdkmath.NewUint(1),
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
			Symbol:     "TEST",
			DenomUnits: []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "testcoin", IsDefaultDisplay: true}},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating wrapper collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "error getting wrapper collection")
	suite.Require().Equal(1, len(collection.CosmosCoinWrapperPaths))

	return collection
}

// TestSpecialWrappingGuardrails_ValidWrapApproval tests that an approval with
// allowSpecialWrapping=true, wrapper address on toListId, and mustPrioritize=true passes.
func (suite *TestSuite) TestSpecialWrappingGuardrails_ValidWrapApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collection := suite.createWrapperCollection()
	wrapperAddr := collection.CosmosCoinWrapperPaths[0].Address

	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        "wrap-to-wrapper",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          wrapperAddr,
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowSpecialWrapping: true,
			MustPrioritize:       true,
		},
	})

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Nil(err, "Should pass: wrapper address on toListId with mustPrioritize=true")
}

// TestSpecialWrappingGuardrails_ValidUnwrapApproval tests that an approval with
// allowSpecialWrapping=true, wrapper address on fromListId, and mustPrioritize=true passes.
func (suite *TestSuite) TestSpecialWrappingGuardrails_ValidUnwrapApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collection := suite.createWrapperCollection()
	wrapperAddr := collection.CosmosCoinWrapperPaths[0].Address

	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        "unwrap-from-wrapper",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        wrapperAddr,
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowSpecialWrapping: true,
			MustPrioritize:       true,
		},
	})

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Nil(err, "Should pass: wrapper address on fromListId with mustPrioritize=true")
}

// TestSpecialWrappingGuardrails_NeitherSideMatches_Rejected tests that an approval with
// allowSpecialWrapping=true but no wrapper address on either side is rejected.
func (suite *TestSuite) TestSpecialWrappingGuardrails_NeitherSideMatches_Rejected() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	_ = suite.createWrapperCollection()

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        "wrapping-neither",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowSpecialWrapping: true,
			MustPrioritize:       true,
		},
	})

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Error(err, "Should reject: neither side is wrapper address")
	suite.Require().Contains(err.Error(), "neither fromListId nor toListId is exactly a wrapper path address")
}

// TestSpecialWrappingGuardrails_MustPrioritize_Rejected tests that an approval with
// allowSpecialWrapping=true but mustPrioritize=false is rejected.
func (suite *TestSuite) TestSpecialWrappingGuardrails_MustPrioritize_Rejected() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collection := suite.createWrapperCollection()
	wrapperAddr := collection.CosmosCoinWrapperPaths[0].Address

	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        "wrapping-no-prioritize",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          wrapperAddr,
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowSpecialWrapping: true,
			MustPrioritize:       false, // Should be rejected
		},
	})

	err := UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	suite.Require().Error(err, "Should reject: mustPrioritize is false")
	suite.Require().Contains(err.Error(), "mustPrioritize is not set to true")
}
