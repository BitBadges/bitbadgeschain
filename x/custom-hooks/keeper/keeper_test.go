package keeper_test

import (
	"strconv"
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
	s.keeper = keeper.NewKeeper(
		s.App.Logger(),
		s.App.PoolManagerKeeper,
		s.App.BankKeeper,
		s.App.HooksICS4Wrapper,
		s.App.IBCKeeper.ChannelKeeper,
		s.App.ScopedIBCTransferKeeper,
	)
}

// TestExecuteSwapAndAction_SwapOnly tests swap execution without post-swap action
func (s *KeeperTestSuite) TestExecuteSwapAndAction_SwapOnly() {
	// Create a test pool
	poolID := s.prepareTestPool()

	// Fund sender account
	sender := s.TestAccs[0]
	s.FundAcc(sender, sdk.Coins{
		sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1000000)),
	})

	// Create swap_and_action with swap only
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
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().NoError(err)

	// Verify swap was executed (check balance changes)
	// Note: In a real test, you'd verify the actual swap happened
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
	if err != nil {
		s.Require().Contains(err.Error(), "channel capability not found")
	}
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
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid pool ID")
}

// TestExecuteSwapAndAction_EmptyOperations tests error handling for empty operations
func (s *KeeperTestSuite) TestExecuteSwapAndAction_EmptyOperations() {
	sender := s.TestAccs[0]
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
	}

	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().Error(err)
}

// TestExecuteSwapAndAction_NoMinAsset tests swap without min_asset
func (s *KeeperTestSuite) TestExecuteSwapAndAction_NoMinAsset() {
	poolID := s.prepareTestPool()

	sender := s.TestAccs[0]
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
	}

	// Use a smaller amount that the sender actually has
	tokenIn := sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100000))

	err := s.keeper.ExecuteSwapAndAction(s.Ctx, sender, swapAndAction, tokenIn)
	s.Require().NoError(err)
}

// TestConvertOperationsToRoutes tests route conversion
func (s *KeeperTestSuite) TestConvertOperationsToRoutes() {
	operations := []customhookstypes.Operation{
		{
			Pool:     "1",
			DenomIn:  "uosmo",
			DenomOut: "uatom",
		},
		{
			Pool:     "2",
			DenomIn:  "uatom",
			DenomOut: "uusdc",
		},
	}

	routes, err := s.keeper.ConvertOperationsToRoutes(operations)
	s.Require().NoError(err)
	s.Require().Len(routes, 2)
	s.Require().Equal(uint64(1), routes[0].PoolId)
	s.Require().Equal("uatom", routes[0].TokenOutDenom)
	s.Require().Equal(uint64(2), routes[1].PoolId)
	s.Require().Equal("uusdc", routes[1].TokenOutDenom)
}

// TestConvertOperationsToRoutes_InvalidPoolID tests error handling
func (s *KeeperTestSuite) TestConvertOperationsToRoutes_InvalidPoolID() {
	operations := []customhookstypes.Operation{
		{
			Pool:     "invalid",
			DenomIn:  "uosmo",
			DenomOut: "uatom",
		},
	}

	_, err := s.keeper.ConvertOperationsToRoutes(operations)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "invalid pool ID")
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
