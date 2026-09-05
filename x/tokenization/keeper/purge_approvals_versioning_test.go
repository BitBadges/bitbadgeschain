package keeper_test

import (
	"time"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func purgeTestOutgoingApproval(id string, start, end sdkmath.Uint, purgeable bool) *types.UserOutgoingApproval {
	return &types.UserOutgoingApproval{
		ToListId:          alice,
		InitiatedByListId: bob,
		TransferTimes:     []*types.UintRange{{Start: start, End: end}},
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        id,
		ApprovalCriteria: &types.OutgoingApprovalCriteria{
			AutoDeletionOptions: &types.AutoDeletionOptions{AllowPurgeIfExpired: purgeable},
		},
	}
}

// A third-party purge removes only the purged approvals; the target's other approvals keep
// their stored version even when their stored form predates mustPrioritize normalisation.
func (suite *PurgeApprovalsTestSuite) TestPurgeDoesNotReversionUntouchedApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.Require().Nil(CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        bob,
		TokensToCreate: []*types.Balance{{TokenIds: GetFullUintRanges()}},
	}}))

	now := uint64(time.Now().UnixMilli())
	expired := purgeTestOutgoingApproval("expired", sdkmath.NewUint(now-2000000), sdkmath.NewUint(now-1000000), true)
	live := purgeTestOutgoingApproval("live", sdkmath.NewUint(1), sdkmath.NewUint(now+10000000000), false)
	// Non-auto-scannable criteria, stored without the mustPrioritize normalisation
	live.ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{{To: alice, Coins: []*sdk.Coin{{Denom: "ubadge", Amount: sdkmath.NewInt(1)}}}}
	live.ApprovalCriteria.MaxNumTransfers = &types.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(1), AmountTrackerId: "live-tracker"}

	suite.Require().Nil(UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator: bob, CollectionId: sdkmath.NewUint(1),
		OutgoingApprovals: []*types.UserOutgoingApproval{expired, live}, UpdateOutgoingApprovals: true,
	}))

	balance, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	suite.Require().Len(balance.OutgoingApprovals, 2)
	for _, approval := range balance.OutgoingApprovals {
		if approval.ApprovalId == "live" {
			approval.ApprovalCriteria.MustPrioritize = false
		}
	}
	suite.Require().Nil(suite.app.TokenizationKeeper.SetUserBalanceInStore(suite.ctx, keeper.ConstructBalanceKey(bob, sdkmath.NewUint(1)), balance, true))
	liveVersionBefore := balance.OutgoingApprovals[1].Version

	_, err = suite.msgServer.PurgeApprovals(wctx, &types.MsgPurgeApprovals{
		Creator: alice, CollectionId: sdkmath.NewUint(1), ApproverAddress: bob, PurgeCounterpartyApprovals: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{{ApprovalId: "expired", ApprovalLevel: "outgoing", ApproverAddress: bob, Version: balance.OutgoingApprovals[0].Version}},
	})
	suite.Require().Nil(err)

	after, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	suite.Require().Len(after.OutgoingApprovals, 1)
	suite.Require().Equal("live", after.OutgoingApprovals[0].ApprovalId)
	suite.Require().True(after.OutgoingApprovals[0].Version.Equal(liveVersionBefore), "untouched approval must keep its version (before %s, after %s)", liveVersionBefore, after.OutgoingApprovals[0].Version)
}
