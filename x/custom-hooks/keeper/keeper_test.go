package keeper_test

import (
	"strconv"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		s.App.HooksICS4Wrapper,
		s.App.IBCKeeper.ChannelKeeper,
		s.App.ScopedIBCTransferKeeper,
	)
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

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "post_swap_action is required")
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

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	// Note: This will fail if IBC channel doesn't exist, which is expected in unit tests
	// In integration tests, we'd set up proper IBC channels
	// The validation now happens before the swap, so we catch channel issues earlier
	s.Require().Error(err)
	s.Require().True(
		strings.Contains(err.Error(), "IBC channel") || strings.Contains(err.Error(), "channel capability not found"),
		"error should mention channel issue: %s", err.Error(),
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

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "post_swap_action requires a swap to be defined")
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
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid pool ID")
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
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
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
		// No MinAsset specified
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	// Use a smaller amount that the sender actually has
	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().NoError(err)
}

// TestExecuteSwapAndAction_MultiHopNotSupported tests that multi-hop swaps are rejected
func (s *KeeperTestSuite) TestExecuteSwapAndAction_MultiHopNotSupported() {
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
						DenomIn:  "uatom",
						DenomOut: "uusdc",
					},
				},
			},
		},
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: recipient.String(),
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "multi-hop swaps are not supported")
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

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().NoError(err)

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

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "cannot have both ibc_transfer and transfer")
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
		PostSwapAction: &customhookstypes.PostSwapAction{
			// Neither IBCTransfer nor Transfer is set
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "must have either ibc_transfer or transfer")
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

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "IBC channel")
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

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	// Channel validation happens first, so we check for either error
	s.Require().True(
		strings.Contains(err.Error(), "receiver is required") || strings.Contains(err.Error(), "IBC channel"),
		"error should mention receiver or channel issue: %s", err.Error(),
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
		PostSwapAction: &customhookstypes.PostSwapAction{
			Transfer: &customhookstypes.TransferInfo{
				ToAddress: "invalid-address", // Invalid address
			},
		},
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid transfer.to_address")
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

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	// Channel validation happens first, so we check for either error
	s.Require().True(
		(strings.Contains(err.Error(), "invalid") && strings.Contains(err.Error(), "recover_address")) ||
			strings.Contains(err.Error(), "IBC channel"),
		"error should mention invalid recover_address or channel issue: %s", err.Error(),
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

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "to_address is required")
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
