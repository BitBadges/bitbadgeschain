package keeper_test

import (
	"math/big"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Ids that do not fit in a uint64 are rejected before any store key is built.
func (suite *TestSuite) TestDeleteCollectionRejectsIdAboveUint64() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	tooLarge := sdkmath.NewUintFromBigInt(new(big.Int).Lsh(big.NewInt(1), 64))
	suite.Require().NotPanics(func() {
		_, err := suite.msgServer.DeleteCollection(wctx, &types.MsgDeleteCollection{Creator: bob, CollectionId: tooLarge})
		suite.Require().Error(err)
	})
	suite.Require().NotPanics(func() {
		_, err := suite.msgServer.TransferTokens(wctx, &types.MsgTransferTokens{Creator: bob, CollectionId: sdkmath.Uint{}})
		suite.Require().Error(err)
	})
}

// A stored coin transfer with a negative amount fails the approval instead of panicking.
func (suite *TestSuite) TestNegativeCoinTransferAmountIsAnError() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	colA := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	negative := rollbackTestTransferApproval("negative-priced", []*types.CoinTransfer{
		{To: charlie, Coins: []*sdk.Coin{{Amount: sdkmath.NewInt(-5), Denom: "ubadge"}}},
	})
	colA[0].CollectionApprovals = []*types.CollectionApproval{rollbackTestMintApproval(), negative}
	colA[0].Transfers[0].PrioritizedApprovals = []*types.ApprovalIdentifierDetails{
		{ApprovalId: "mint-test", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
	}
	err := CreateCollections(suite, wctx, colA)
	if err != nil {
		// Rejected at validation time is the preferred outcome
		return
	}

	suite.Require().NotPanics(func() {
		err := TransferTokens(suite, wctx, &types.MsgTransferTokens{
			Creator:      bob,
			CollectionId: sdkmath.NewUint(1),
			Transfers: []*types.Transfer{{
				From:        bob,
				ToAddresses: []string{alice},
				Balances:    []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{ApprovalId: "negative-priced", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
				},
			}},
		})
		suite.Require().Error(err)
	})
}
