package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestCollectionApprovalValidationDefersBalanceCleaning() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.Require().NoError(CreateCollections(suite, wctx, GetTransferableCollectionToCreateAllMintedToCreator(bob)))
	approval := GetBobApproval()
	approval.ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{UseOverallNumTransfers: true},
		ManualBalances: []*types.ManualBalances{{Balances: []*types.Balance{
			{Amount: sdkmath.OneUint(), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()},
			{Amount: sdkmath.OneUint(), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()},
		}}},
	}
	msg := &types.MsgSetCollectionApprovals{Creator: bob, CollectionId: sdkmath.OneUint(), CollectionApprovals: []*types.CollectionApproval{approval}}
	suite.Require().NotPanics(func() { suite.Require().NoError(msg.ValidateBasic()) })
	suite.Require().Len(approval.ApprovalCriteria.PredeterminedBalances.ManualBalances[0].Balances, 2)
	_, err := suite.msgServer.SetCollectionApprovals(wctx, msg)
	suite.Require().NoError(err)
	collection, err := GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	balances := collection.CollectionApprovals[0].ApprovalCriteria.PredeterminedBalances.ManualBalances[0].Balances
	suite.Require().Len(balances, 1)
	suite.Require().Equal(sdkmath.NewUint(2), balances[0].Amount)
}
