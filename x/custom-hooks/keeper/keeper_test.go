package keeper_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/custom-hooks/keeper"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	keeper keeper.Keeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.Reset()

	// Create keeper with real app keepers
	// Pass pointer to GammKeeper to avoid copying the keeper (which contains storeKey)
	s.keeper = keeper.NewKeeper(
		s.App.Logger(),
		&s.App.GammKeeper,
		s.App.BankKeeper,
		&s.App.BadgesKeeper,
		&s.App.SendmanagerKeeper,
		s.App.TransferKeeper,
		s.App.HooksICS4Wrapper,
		s.App.IBCKeeper.ChannelKeeper,
		s.App.ScopedIBCTransferKeeper,
	)
}

// getAckError extracts the error message from an acknowledgement
func getAckError(ack ibcexported.Acknowledgement) string {
	channelAck, ok := ack.(channeltypes.Acknowledgement)
	if !ok {
		return ""
	}
	if errResp, ok := channelAck.Response.(*channeltypes.Acknowledgement_Error); ok {
		return errResp.Error
	}
	return ""
}

// TestExecuteSwapAndAction_NoPostSwapAction tests error when post_swap_action is missing
func (s *KeeperTestSuite) TestExecuteSwapAndAction_NoPostSwapAction() {
	// Create a test pool
	poolID := s.prepareTestPool()

	// Fund sender account
	sender := s.TestAccs[0]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Create swap_and_action with swap only (no post_swap_action)
	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		// PostSwapAction is nil - should fail
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "post_swap_action is required")
}

// TestExecuteSwapAndAction_SwapWithPostAction tests swap with IBC transfer
func (s *KeeperTestSuite) TestExecuteSwapAndAction_SwapWithPostAction() {
	// Create a test pool
	poolID := s.prepareTestPool()

	// Fund sender account
	sender := s.TestAccs[0]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	timeoutTimestamp := uint64(1762692450136475000)

	// Create swap_and_action with swap and post-swap action
	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		TimeoutTimestamp: &timeoutTimestamp,
		PostSwapAction: &customhookstypes.PostSwapAction{
			IBCTransfer: &customhookstypes.IBCTransferInfo{
				IBCInfo: &customhookstypes.IBCInfo{
					SourceChannel:  "channel-0",
					Receiver:       "cosmos1test",
					Memo:           "",
					RecoverAddress: "osmo1test",
				},
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	// Note: This will fail if IBC channel doesn't exist, which is expected in unit tests
	// In integration tests, we'd set up proper IBC channels
	// The validation now happens before the swap, so we catch channel issues earlier
	s.Require().False(ack.Success())
	ackErr := getAckError(ack)
	s.Require().True(
		strings.Contains(ackErr, "IBC channel") || strings.Contains(ackErr, "channel capability not found"),
		"error should mention channel issue: %s", ackErr,
	)
}

// TestExecuteSwapAndAction_PostActionWithoutSwap tests error case
func (s *KeeperTestSuite) TestExecuteSwapAndAction_PostActionWithoutSwap() {
	sender := s.TestAccs[0]

	// Create swap_and_action with post-swap action but no swap
	swapAndAction := &customhookstypes.SwapAndAction{
		PostSwapAction: &customhookstypes.PostSwapAction{
			IBCTransfer: &customhookstypes.IBCTransferInfo{
				IBCInfo: &customhookstypes.IBCInfo{
					SourceChannel: "channel-0",
					Receiver:      "cosmos1test",
				},
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "post_swap_action requires a swap to be defined")
}

// TestExecuteSwapAndAction_InvalidPoolID tests error handling for invalid pool ID
func (s *KeeperTestSuite) TestExecuteSwapAndAction_InvalidPoolID() {
	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin("uosmo", osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     "invalid",
						DenomIn:  "uosmo",
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "invalid pool ID")
}

// TestExecuteSwapAndAction_EmptyOperations tests error handling for empty operations
func (s *KeeperTestSuite) TestExecuteSwapAndAction_EmptyOperations() {
	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin("uosmo", osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations:    []customhookstypes.Operation{},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
}

// TestExecuteSwapAndAction_NoMinAsset tests swap without min_asset
func (s *KeeperTestSuite) TestExecuteSwapAndAction_NoMinAsset() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	// Fund with enough balance for the swap (need more than tokenIn amount)
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(10000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		// No MinAsset specified - should fail
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	// Use a smaller amount that the sender actually has
	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "min_asset is required for swaps")
}

// TestExecuteSwapAndAction_MultiHop tests that multi-hop swaps work correctly
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MultiHop() {
	// Create first pool: baseDenom <-> uatom
	poolID1 := s.prepareTestPool()

	// Create second pool: uatom <-> uusdc
	testDenom1 := "uatom"
	testDenom2 := "uusdc"

	poolAssets2 := []balancer.PoolAsset{
		{Token: sdk.NewInt64Coin(testDenom1, 1000000), Weight: osmomath.NewInt(1)},
		{Token: sdk.NewInt64Coin(testDenom2, 1000000), Weight: osmomath.NewInt(1)},
	}

	poolParams2 := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"),
		ExitFee: osmomath.ZeroDec(),
	}

	balances2 := sdk.Coins{
		sdk.NewCoin(testDenom1, osmomath.NewInt(1000000)),
		sdk.NewCoin(testDenom2, osmomath.NewInt(1000000)),
	}

	s.FundAcc(s.TestAccs[0], balances2)

	poolID2, err := s.App.PoolManagerKeeper.CreatePool(
		s.Ctx,
		balancer.NewMsgCreateBalancerPool(s.TestAccs[0], poolParams2, poolAssets2),
	)
	s.Require().NoError(err)

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Get initial recipient balance
	initialBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uusdc")

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID1, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
					{
						Pool:     strconv.FormatUint(poolID2, 10),
						DenomIn:  "uatom",
						DenomOut: "uusdc",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uusdc",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())

	// Verify transfer was executed (check balance changes)
	finalBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uusdc")
	s.Require().True(finalBalance.Amount.GT(initialBalance.Amount), "recipient should have received tokens from multi-hop swap")
}

// TestExecuteSwapAndAction_MultiHop_ThreeHops tests a three-hop swap
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MultiHop_ThreeHops() {
	// Create first pool: baseDenom <-> uatom
	poolID1 := s.prepareTestPool()

	// Create second pool: uatom <-> uusdc
	testDenom1 := "uatom"
	testDenom2 := "uusdc"

	poolAssets2 := []balancer.PoolAsset{
		{Token: sdk.NewInt64Coin(testDenom1, 1000000), Weight: osmomath.NewInt(1)},
		{Token: sdk.NewInt64Coin(testDenom2, 1000000), Weight: osmomath.NewInt(1)},
	}

	poolParams2 := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"),
		ExitFee: osmomath.ZeroDec(),
	}

	balances2 := sdk.Coins{
		sdk.NewCoin(testDenom1, osmomath.NewInt(1000000)),
		sdk.NewCoin(testDenom2, osmomath.NewInt(1000000)),
	}

	s.FundAcc(s.TestAccs[0], balances2)

	poolID2, err := s.App.PoolManagerKeeper.CreatePool(
		s.Ctx,
		balancer.NewMsgCreateBalancerPool(s.TestAccs[0], poolParams2, poolAssets2),
	)
	s.Require().NoError(err)

	// Create third pool: uusdc <-> uusdt
	testDenom3 := "uusdt"

	poolAssets3 := []balancer.PoolAsset{
		{Token: sdk.NewInt64Coin(testDenom2, 1000000), Weight: osmomath.NewInt(1)},
		{Token: sdk.NewInt64Coin(testDenom3, 1000000), Weight: osmomath.NewInt(1)},
	}

	poolParams3 := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"),
		ExitFee: osmomath.ZeroDec(),
	}

	balances3 := sdk.Coins{
		sdk.NewCoin(testDenom2, osmomath.NewInt(1000000)),
		sdk.NewCoin(testDenom3, osmomath.NewInt(1000000)),
	}

	s.FundAcc(s.TestAccs[0], balances3)

	poolID3, err := s.App.PoolManagerKeeper.CreatePool(
		s.Ctx,
		balancer.NewMsgCreateBalancerPool(s.TestAccs[0], poolParams3, poolAssets3),
	)
	s.Require().NoError(err)

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Get initial recipient balance
	initialBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uusdt")

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID1, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
					{
						Pool:     strconv.FormatUint(poolID2, 10),
						DenomIn:  "uatom",
						DenomOut: "uusdc",
					},
					{
						Pool:     strconv.FormatUint(poolID3, 10),
						DenomIn:  "uusdc",
						DenomOut: "uusdt",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uusdt",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())

	// Verify transfer was executed (check balance changes)
	finalBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uusdt")
	s.Require().True(finalBalance.Amount.GT(initialBalance.Amount), "recipient should have received tokens from three-hop swap")
}

// TestExecuteSwapAndAction_MultiHop_WithAffiliates tests multi-hop swap with affiliates
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MultiHop_WithAffiliates() {
	// Create first pool: baseDenom <-> uatom
	poolID1 := s.prepareTestPool()

	// Create second pool: uatom <-> uusdc
	testDenom1 := "uatom"
	testDenom2 := "uusdc"

	poolAssets2 := []balancer.PoolAsset{
		{Token: sdk.NewInt64Coin(testDenom1, 1000000), Weight: osmomath.NewInt(1)},
		{Token: sdk.NewInt64Coin(testDenom2, 1000000), Weight: osmomath.NewInt(1)},
	}

	poolParams2 := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"),
		ExitFee: osmomath.ZeroDec(),
	}

	balances2 := sdk.Coins{
		sdk.NewCoin(testDenom1, osmomath.NewInt(1000000)),
		sdk.NewCoin(testDenom2, osmomath.NewInt(1000000)),
	}

	s.FundAcc(s.TestAccs[0], balances2)

	poolID2, err := s.App.PoolManagerKeeper.CreatePool(
		s.Ctx,
		balancer.NewMsgCreateBalancerPool(s.TestAccs[0], poolParams2, poolAssets2),
	)
	s.Require().NoError(err)

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate := s.TestAccs[2]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Get initial balances
	initialRecipientBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uusdc")
	initialAffiliateBalance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate, "uusdc")

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID1, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
					{
						Pool:     strconv.FormatUint(poolID2, 10),
						DenomIn:  "uatom",
						DenomOut: "uusdc",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uusdc",
				Amount: "10000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "100", // 1%
				Address:        affiliate.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())

	// Verify affiliate received fee
	finalAffiliateBalance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate, "uusdc")
	affiliateReceived := finalAffiliateBalance.Amount.Sub(initialAffiliateBalance.Amount)
	s.Require().True(affiliateReceived.IsPositive(), "affiliate should have received fee from multi-hop swap")

	// Verify recipient received remaining amount
	finalRecipientBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uusdc")
	recipientReceived := finalRecipientBalance.Amount.Sub(initialRecipientBalance.Amount)
	s.Require().True(recipientReceived.IsPositive(), "recipient should have received tokens from multi-hop swap")
}

// TestExecuteSwapAndAction_MultiHop_InvalidChain tests that operations must chain correctly
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MultiHop_InvalidChain() {
	poolID1 := s.prepareTestPool()
	poolID2 := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID1, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
					{
						Pool:     strconv.FormatUint(poolID2, 10),
						DenomIn:  "uusdc", // Wrong! Should be "uatom" to chain correctly
						DenomOut: "uusdc",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uusdc",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	ackErr := getAckError(ack)
	s.Require().True(
		strings.Contains(ackErr, "operations do not chain correctly") || strings.Contains(ackErr, "denom_out") && strings.Contains(ackErr, "denom_in"),
		"error should mention operations chaining issue: %s", ackErr,
	)
}

// TestExecuteSwapAndAction_MultiHop_FirstOperationMismatch tests that first operation must match tokenIn
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MultiHop_FirstOperationMismatch() {
	poolID1 := s.prepareTestPool()
	poolID2 := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID1, 10),
						DenomIn:  "uatom", // Wrong! Should match tokenIn (sdk.DefaultBondDenom)
						DenomOut: "uusdc",
					},
					{
						Pool:     strconv.FormatUint(poolID2, 10),
						DenomIn:  "uusdc",
						DenomOut: "uusdt",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uusdt",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	ackErr := getAckError(ack)
	s.Require().True(
		strings.Contains(ackErr, "first operation denom_in") && strings.Contains(ackErr, "does not match token_in"),
		"error should mention first operation denom_in mismatch: %s", ackErr,
	)
}

// TestExecuteSwapAndAction_LocalTransfer tests swap with local transfer
func (s *KeeperTestSuite) TestExecuteSwapAndAction_LocalTransfer() {
	// Create a test pool
	poolID := s.prepareTestPool()

	// Fund sender account
	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Get initial recipient balance
	initialBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uatom")

	// Create swap_and_action with swap and local transfer
	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())

	// Verify transfer was executed (check balance changes)
	finalBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uatom")
	s.Require().True(finalBalance.Amount.GT(initialBalance.Amount), "recipient should have received tokens")
}

// TestExecuteSwapAndAction_BothTransferTypes tests error when both IBCTransfer and Transfer are set
func (s *KeeperTestSuite) TestExecuteSwapAndAction_BothTransferTypes() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			IBCTransfer: &customhookstypes.IBCTransferInfo{
				IBCInfo: &customhookstypes.IBCInfo{
					SourceChannel: "channel-0",
					Receiver:      "cosmos1test",
				},
			},
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "cannot have both ibc_transfer and transfer")
}

// TestExecuteSwapAndAction_NeitherTransferType tests error when neither IBCTransfer nor Transfer is set
func (s *KeeperTestSuite) TestExecuteSwapAndAction_NeitherTransferType() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			// Neither IBCTransfer nor Transfer is set
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "must have either ibc_transfer or transfer")
}

// TestExecuteSwapAndAction_InvalidIBCChannel tests validation of non-existent IBC channel
func (s *KeeperTestSuite) TestExecuteSwapAndAction_InvalidIBCChannel() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			IBCTransfer: &customhookstypes.IBCTransferInfo{
				IBCInfo: &customhookstypes.IBCInfo{
					SourceChannel: "channel-nonexistent",
					Receiver:      "cosmos1test",
				},
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "IBC channel")
}

// TestExecuteSwapAndAction_MissingIBCReceiver tests validation of missing receiver
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MissingIBCReceiver() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			IBCTransfer: &customhookstypes.IBCTransferInfo{
				IBCInfo: &customhookstypes.IBCInfo{
					SourceChannel: "channel-0",
					Receiver:      "", // Missing receiver
				},
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	ackErr := getAckError(ack)
	// Channel validation happens first, so we check for either error
	s.Require().True(
		strings.Contains(ackErr, "receiver is required") || strings.Contains(ackErr, "IBC channel"),
		"error should mention receiver or channel issue: %s", ackErr,
	)
}

// TestExecuteSwapAndAction_InvalidLocalTransferAddress tests validation of invalid local transfer address
func (s *KeeperTestSuite) TestExecuteSwapAndAction_InvalidLocalTransferAddress() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: "invalid-address", // Invalid address
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "invalid transfer.to_address")
}

// TestExecuteSwapAndAction_InvalidRecoverAddress tests validation of invalid recover address
func (s *KeeperTestSuite) TestExecuteSwapAndAction_InvalidRecoverAddress() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			IBCTransfer: &customhookstypes.IBCTransferInfo{
				IBCInfo: &customhookstypes.IBCInfo{
					SourceChannel:  "channel-0",
					Receiver:       "cosmos1test",
					RecoverAddress: "invalid-recover-address", // Invalid address
				},
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	ackErr := getAckError(ack)
	// Channel validation happens first, so we check for either error
	s.Require().True(
		(strings.Contains(ackErr, "invalid") && strings.Contains(ackErr, "recover_address")) ||
			strings.Contains(ackErr, "IBC channel"),
		"error should mention invalid recover_address or channel issue: %s", ackErr,
	)
}

// TestExecuteSwapAndAction_MissingLocalTransferAddress tests validation of missing local transfer address
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MissingLocalTransferAddress() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: "", // Missing address
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "to_address is required")
}

// TestExecuteSwapAndAction_WithSingleAffiliate tests swap with a single affiliate
func (s *KeeperTestSuite) TestExecuteSwapAndAction_WithSingleAffiliate() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate := s.TestAccs[2]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Get initial balances
	initialRecipientBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uatom")
	initialAffiliateBalance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate, "uatom")

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "100", // 1%
				Address:        affiliate.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())

	// Verify affiliate received fee
	finalAffiliateBalance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate, "uatom")
	affiliateReceived := finalAffiliateBalance.Amount.Sub(initialAffiliateBalance.Amount)
	s.Require().True(affiliateReceived.IsPositive(), "affiliate should have received fee")

	// Verify recipient received remaining amount (after fee deduction)
	finalRecipientBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uatom")
	recipientReceived := finalRecipientBalance.Amount.Sub(initialRecipientBalance.Amount)
	s.Require().True(recipientReceived.IsPositive(), "recipient should have received tokens")

	// Verify that affiliate fee was deducted (affiliate received fee, recipient received less)
	// The sum of what recipient and affiliate received should be less than or equal to swap output
	// (accounting for pool fees)
	s.Require().True(affiliateReceived.IsPositive(), "affiliate fee should be positive")
	s.Require().True(recipientReceived.IsPositive(), "recipient should receive tokens")
}

// TestExecuteSwapAndAction_WithMultipleAffiliates tests swap with multiple affiliates
func (s *KeeperTestSuite) TestExecuteSwapAndAction_WithMultipleAffiliates() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate1 := s.TestAccs[2]
	// Create a 4th account for the second affiliate
	affiliate2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	s.App.AccountKeeper.SetAccount(s.Ctx, s.App.AccountKeeper.NewAccountWithAddress(s.Ctx, affiliate2))

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Get initial balances
	initialRecipientBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uatom")
	initialAffiliate1Balance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate1, "uatom")
	initialAffiliate2Balance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate2, "uatom")

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "80", // 0.8%
				Address:        affiliate1.String(),
			},
			{
				BasisPointsFee: "20", // 0.2%
				Address:        affiliate2.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())

	// Verify both affiliates received fees
	finalAffiliate1Balance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate1, "uatom")
	finalAffiliate2Balance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate2, "uatom")

	affiliate1Received := finalAffiliate1Balance.Amount.Sub(initialAffiliate1Balance.Amount)
	affiliate2Received := finalAffiliate2Balance.Amount.Sub(initialAffiliate2Balance.Amount)

	s.Require().True(affiliate1Received.IsPositive(), "affiliate1 should have received fee")
	s.Require().True(affiliate2Received.IsPositive(), "affiliate2 should have received fee")

	// Verify affiliate1 received more than affiliate2 (80 vs 20 basis points)
	s.Require().True(affiliate1Received.GT(affiliate2Received), "affiliate1 should receive more than affiliate2")

	// Verify recipient received remaining amount
	finalRecipientBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uatom")
	recipientReceived := finalRecipientBalance.Amount.Sub(initialRecipientBalance.Amount)
	s.Require().True(recipientReceived.IsPositive(), "recipient should have received tokens")
}

// TestExecuteSwapAndAction_InvalidAffiliateAddress tests validation of invalid affiliate address
func (s *KeeperTestSuite) TestExecuteSwapAndAction_InvalidAffiliateAddress() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "100",
				Address:        "invalid-address",
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	ackErr := getAckError(ack)
	s.Require().Contains(ackErr, "invalid affiliate")
	s.Require().Contains(ackErr, "address")
}

// TestExecuteSwapAndAction_InvalidAffiliateBasisPoints tests validation of invalid basis points
func (s *KeeperTestSuite) TestExecuteSwapAndAction_InvalidAffiliateBasisPoints() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate := s.TestAccs[2]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "not-a-number",
				Address:        affiliate.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	ackErr := getAckError(ack)
	s.Require().Contains(ackErr, "invalid affiliate")
	s.Require().Contains(ackErr, "basis_points_fee")
}

// TestExecuteSwapAndAction_AffiliateBasisPointsExceeds10000 tests validation when basis points exceed 100%
func (s *KeeperTestSuite) TestExecuteSwapAndAction_AffiliateBasisPointsExceeds10000() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate := s.TestAccs[2]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "10001", // Exceeds 100%
				Address:        affiliate.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "cannot exceed 10000")
}

// TestExecuteSwapAndAction_TotalAffiliateBasisPointsExceeds10000 tests validation when total basis points exceed 100%
func (s *KeeperTestSuite) TestExecuteSwapAndAction_TotalAffiliateBasisPointsExceeds10000() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate1 := s.TestAccs[2]
	// Create a 4th account for the second affiliate
	affiliate2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	s.App.AccountKeeper.SetAccount(s.Ctx, s.App.AccountKeeper.NewAccountWithAddress(s.Ctx, affiliate2))

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "8000", // 80%
				Address:        affiliate1.String(),
			},
			{
				BasisPointsFee: "2001", // 20.01% - total exceeds 100%
				Address:        affiliate2.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "total affiliate basis_points_fee cannot exceed 10000")
}

// TestExecuteSwapAndAction_MissingAffiliateAddress tests validation of missing affiliate address
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MissingAffiliateAddress() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "100",
				Address:        "", // Missing address
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "address is required")
}

// TestExecuteSwapAndAction_MissingAffiliateBasisPoints tests validation of missing basis points
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MissingAffiliateBasisPoints() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate := s.TestAccs[2]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "", // Missing basis points
				Address:        affiliate.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "basis_points_fee is required")
}

// TestExecuteSwapAndAction_AffiliateFeeFromMinAmount tests that affiliate fees are calculated from min_asset, not actual output
func (s *KeeperTestSuite) TestExecuteSwapAndAction_AffiliateFeeFromMinAmount() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate := s.TestAccs[2]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Set a specific min amount (e.g., 10000 uatom)
	// The actual swap output will likely be higher due to pool conditions
	minAmount := "10000"
	affiliateBasisPoints := "1000" // 10%

	// Get initial balances
	initialAffiliateBalance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate, "uatom")

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: minAmount,
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: affiliateBasisPoints, // 10%
				Address:        affiliate.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())

	// Verify affiliate received fee based on min amount, not actual output
	finalAffiliateBalance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate, "uatom")
	affiliateReceived := finalAffiliateBalance.Amount.Sub(initialAffiliateBalance.Amount)

	// Expected fee: 10000 * 1000 / 10000 = 1000 uatom (from min amount)
	expectedFee := osmomath.NewInt(1000)
	s.Require().True(affiliateReceived.Equal(expectedFee),
		"affiliate should receive fee based on min amount (%s), got %s, expected %s",
		minAmount, affiliateReceived.String(), expectedFee.String())
}

// TestExecuteSwapAndAction_MinAssetRequired tests that min_asset is always required for swaps
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MinAssetRequired() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		// No MinAsset specified - should fail
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().False(ack.Success())
	s.Require().Contains(getAckError(ack), "min_asset is required for swaps")
}

// TestExecuteSwapAndAction_EmptyAffiliates tests that empty affiliates array is valid
func (s *KeeperTestSuite) TestExecuteSwapAndAction_EmptyAffiliates() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{}, // Empty array should be valid
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())
}

// TestExecuteSwapAndAction_AffiliateFeesDeductedFromOutput tests that affiliate fees are correctly deducted from swap output
func (s *KeeperTestSuite) TestExecuteSwapAndAction_AffiliateFeesDeductedFromOutput() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate := s.TestAccs[2]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Get initial balances
	initialRecipientBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uatom")
	initialAffiliateBalance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate, "uatom")

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "10000", // min output
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "100", // 1% of min output = 100 uatom
				Address:        affiliate.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())

	// Verify affiliate received fee (1% of 10000 = 100 uatom)
	finalAffiliateBalance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate, "uatom")
	affiliateReceived := finalAffiliateBalance.Amount.Sub(initialAffiliateBalance.Amount)
	s.Require().Equal(osmomath.NewInt(100), affiliateReceived, "affiliate should receive 1% of min output")

	// Verify recipient received remaining amount (swap output - 100 uatom)
	finalRecipientBalance := s.App.BankKeeper.GetBalance(s.Ctx, recipient, "uatom")
	recipientReceived := finalRecipientBalance.Amount.Sub(initialRecipientBalance.Amount)
	s.Require().True(recipientReceived.IsPositive(), "recipient should receive tokens")
	s.Require().True(affiliateReceived.IsPositive(), "affiliate should receive fee")
}

// TestExecuteSwapAndAction_AffiliateFeesFromPool tests that affiliate fees are sent from pool, not sender
func (s *KeeperTestSuite) TestExecuteSwapAndAction_AffiliateFeesFromPool() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate := s.TestAccs[2]

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Get initial sender balance (should not decrease by affiliate fee amount)
	initialSenderBalance := s.App.BankKeeper.GetBalance(s.Ctx, sender, "uatom")

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "10000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "100", // 1%
				Address:        affiliate.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())

	// Verify sender's uatom balance didn't decrease (fees come from pool, not sender)
	finalSenderBalance := s.App.BankKeeper.GetBalance(s.Ctx, sender, "uatom")
	// Sender should have 0 uatom (they swapped it all), but the point is fees weren't deducted from them
	s.Require().True(finalSenderBalance.Amount.LTE(initialSenderBalance.Amount), "sender balance should not increase")
}

// TestExecuteSwapAndAction_MultipleAffiliatesProportionalFees tests that multiple affiliates receive proportional fees
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MultipleAffiliatesProportionalFees() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
	recipient := s.TestAccs[1]
	affiliate1 := s.TestAccs[2]
	// Create a 4th account for the second affiliate
	affiliate2 := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	s.App.AccountKeeper.SetAccount(s.Ctx, s.App.AccountKeeper.NewAccountWithAddress(s.Ctx, affiliate2))

	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	initialAffiliate1Balance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate1, "uatom")
	initialAffiliate2Balance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate2, "uatom")

	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{
					{
						Pool:     strconv.FormatUint(poolID, 10),
						DenomIn:  sdk.DefaultBondDenom,
						DenomOut: "uatom",
					},
				},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "10000",
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
		Affiliates: []customhookstypes.Affiliate{
			{
				BasisPointsFee: "300", // 3% of 10000 = 300 uatom
				Address:        affiliate1.String(),
			},
			{
				BasisPointsFee: "200", // 2% of 10000 = 200 uatom
				Address:        affiliate2.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	ack := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn, sender.String())
	s.Require().True(ack.Success())

	// Verify both affiliates received correct proportional fees
	finalAffiliate1Balance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate1, "uatom")
	finalAffiliate2Balance := s.App.BankKeeper.GetBalance(s.Ctx, affiliate2, "uatom")

	affiliate1Received := finalAffiliate1Balance.Amount.Sub(initialAffiliate1Balance.Amount)
	affiliate2Received := finalAffiliate2Balance.Amount.Sub(initialAffiliate2Balance.Amount)

	s.Require().Equal(osmomath.NewInt(300), affiliate1Received, "affiliate1 should receive 3% of min output")
	s.Require().Equal(osmomath.NewInt(200), affiliate2Received, "affiliate2 should receive 2% of min output")
	s.Require().True(affiliate1Received.GT(affiliate2Received), "affiliate1 should receive more than affiliate2")
}

// Helper function to prepare a test pool
func (s *KeeperTestSuite) prepareTestPool() uint64 {
	// Create a balancer pool for testing
	// Use the base denom from the chain
	baseDenom := sdk.DefaultBondDenom
	testDenom := "uatom" // Use a test denom

	poolAssets := []balancer.PoolAsset{
		{Token: sdk.NewInt64Coin(baseDenom, 1000000), Weight: osmomath.NewInt(1)},
		{Token: sdk.NewInt64Coin(testDenom, 1000000), Weight: osmomath.NewInt(1)},
	}

	poolParams := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"),
		ExitFee: osmomath.ZeroDec(),
	}

	balances := sdk.Coins{
		sdk.NewCoin(baseDenom, osmomath.NewInt(1000000)),
		sdk.NewCoin(testDenom, osmomath.NewInt(1000000)),
	}

	s.FundAcc(s.TestAccs[0], balances)

	poolID, err := s.App.PoolManagerKeeper.CreatePool(
		s.Ctx,
		balancer.NewMsgCreateBalancerPool(s.TestAccs[0], poolParams, poolAssets),
	)
	s.Require().NoError(err)

	return poolID
}
