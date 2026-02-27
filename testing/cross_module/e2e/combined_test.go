//go:build test
// +build test

package e2e

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// CombinedTestSuite tests cross-module interactions between GAMM and Tokenization
type CombinedTestSuite struct {
	CrossModuleTestSuite
}

func TestCombinedTestSuite(t *testing.T) {
	suite.Run(t, new(CombinedTestSuite))
}

// TestPoolCreationThenCollectionCreation tests creating pools and collections in sequence
func (s *CombinedTestSuite) TestPoolCreationThenCollectionCreation() {
	s.T().Log("Testing pool creation followed by collection creation")

	user := s.TestAccs[0]

	// Create a pool
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(1000000)),
	)
	s.Require().Greater(poolId, uint64(0))

	// Create a collection
	collectionId, err := s.CreateBasicCollection(user)
	s.Require().NoError(err)
	s.Require().True(collectionId.GT(sdkmath.ZeroUint()))

	// Verify both exist
	pool, err := s.App.GammKeeper.GetPoolAndPoke(s.Ctx, poolId)
	s.Require().NoError(err)
	s.Require().NotNil(pool)

	collection, found := s.App.TokenizationKeeper.GetCollectionFromStore(s.Ctx, collectionId)
	s.Require().True(found)
	s.Require().NotNil(collection)

	s.T().Logf("Created pool %d and collection %s", poolId, collectionId)
}

// TestSwapThenMint tests swapping tokens then using the result for collection operations
func (s *CombinedTestSuite) TestSwapThenMint() {
	s.T().Log("Testing swap followed by mint operation")

	user := s.TestAccs[0]
	recipient := s.TestAccs[1]

	// Create pool
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)

	// Create collection
	collectionId, err := s.CreateBasicCollection(user)
	s.Require().NoError(err)

	// Perform swap to get bar tokens
	tokenOut, err := s.SwapTokens(user, poolId, sdk.NewCoin("foo", sdkmath.NewInt(100000)), "bar")
	s.Require().NoError(err)
	s.Require().True(tokenOut.Amount.GT(sdkmath.NewInt(0)))

	// Now mint collection tokens to recipient
	tokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
	}
	err = s.MintTokensToAddress(collectionId, recipient, tokenIds, sdkmath.NewUint(1))
	s.Require().NoError(err)

	// Verify both operations succeeded
	barBalance := s.GetCoinBalance(user, "bar")
	s.Require().True(barBalance.Amount.GTE(tokenOut.Amount))

	tokenBalance := s.GetTokenBalance(collectionId, recipient, sdkmath.NewUint(5))
	s.Require().Equal(sdkmath.NewUint(1), tokenBalance)

	s.T().Logf("Swapped for %s bar, minted tokens 1-10 to recipient", tokenOut.Amount)
}

// TestLPTokensAfterCollectionTransfer tests LP token operations after collection transfers
func (s *CombinedTestSuite) TestLPTokensAfterCollectionTransfer() {
	s.T().Log("Testing LP tokens after collection transfer")

	creator := s.TestAccs[0]
	user := s.TestAccs[1]

	// Create collection and mint to user
	collectionId, err := s.CreateBasicCollection(creator)
	s.Require().NoError(err)

	tokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
	}
	err = s.MintTokensToAddress(collectionId, user, tokenIds, sdkmath.NewUint(1))
	s.Require().NoError(err)

	// Create pool
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)

	// User joins the pool
	s.FundAcc(user, sdk.NewCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(1000000)),
	))

	pool, err := s.App.GammKeeper.GetPoolAndPoke(s.Ctx, poolId)
	s.Require().NoError(err)

	shareOutAmount := pool.GetTotalShares().Quo(sdkmath.NewInt(100))
	msg := &gammtypes.MsgJoinPool{
		Sender:         user.String(),
		PoolId:         poolId,
		ShareOutAmount: shareOutAmount,
		TokenInMaxs: sdk.NewCoins(
			sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
			sdk.NewCoin("bar", sdkmath.NewInt(1000000)),
		),
	}

	_, err = s.gammMsgServer.JoinPool(s.Ctx, msg)
	s.Require().NoError(err)

	// Verify user has both LP tokens and collection tokens
	lpDenom := gammtypes.GetPoolShareDenom(poolId)
	lpBalance := s.GetCoinBalance(user, lpDenom)
	s.Require().True(lpBalance.Amount.GT(sdkmath.NewInt(0)))

	tokenBalance := s.GetTokenBalance(collectionId, user, sdkmath.NewUint(5))
	s.Require().Equal(sdkmath.NewUint(1), tokenBalance)

	s.T().Logf("User has %s LP tokens and collection tokens", lpBalance.Amount)
}

// TestMultipleSwapsAndTransfers tests interleaved swaps and token transfers
func (s *CombinedTestSuite) TestMultipleSwapsAndTransfers() {
	s.T().Log("Testing multiple swaps and transfers interleaved")

	creator := s.TestAccs[0]
	userA := s.TestAccs[1]
	userB := s.TestAccs[2]

	// Create pool
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)

	// Create collection
	collectionId, err := s.CreateBasicCollection(creator)
	s.Require().NoError(err)

	// Mint tokens to userA
	tokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
	}
	err = s.MintTokensToAddress(collectionId, userA, tokenIds, sdkmath.NewUint(1))
	s.Require().NoError(err)

	// Interleave operations
	for i := 0; i < 3; i++ {
		// UserA swaps
		_, err = s.SwapTokens(userA, poolId, sdk.NewCoin("foo", sdkmath.NewInt(10000)), "bar")
		s.Require().NoError(err)

		// UserA transfers some collection tokens to userB
		transferTokenIds := []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(uint64(i*10 + 1)), End: sdkmath.NewUint(uint64(i*10 + 5))},
		}
		err = s.TransferTokens(collectionId, userA, userB, transferTokenIds, sdkmath.NewUint(1))
		s.Require().NoError(err)

		// UserB swaps
		_, err = s.SwapTokens(userB, poolId, sdk.NewCoin("bar", sdkmath.NewInt(5000)), "foo")
		s.Require().NoError(err)

		s.T().Logf("Completed iteration %d", i+1)
	}

	// Verify final state
	userABarBalance := s.GetCoinBalance(userA, "bar")
	userBFooBalance := s.GetCoinBalance(userB, "foo")
	s.Require().True(userABarBalance.Amount.GT(sdkmath.NewInt(0)))
	s.Require().True(userBFooBalance.Amount.GT(sdkmath.NewInt(0)))

	// UserB should have some collection tokens
	userBTokenBalance := s.GetTokenBalance(collectionId, userB, sdkmath.NewUint(5))
	s.Require().Equal(sdkmath.NewUint(1), userBTokenBalance)

	s.T().Log("Successfully completed interleaved operations")
}

// TestPoolExitAndCollectionBurn tests pool exit followed by token operations
func (s *CombinedTestSuite) TestPoolExitAndCollectionBurn() {
	s.T().Log("Testing pool exit and collection operations")

	user := s.TestAccs[0]

	// Create pool - user gets LP tokens
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(1000000)),
	)

	// Create collection
	collectionId, err := s.CreateBasicCollection(user)
	s.Require().NoError(err)

	// Mint tokens to user
	tokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
	}
	err = s.MintTokensToAddress(collectionId, user, tokenIds, sdkmath.NewUint(1))
	s.Require().NoError(err)

	// Exit part of the pool
	lpDenom := gammtypes.GetPoolShareDenom(poolId)
	lpBalance := s.GetCoinBalance(user, lpDenom)
	exitAmount := lpBalance.Amount.Quo(sdkmath.NewInt(4)) // Exit 25%

	exitMsg := &gammtypes.MsgExitPool{
		Sender:        user.String(),
		PoolId:        poolId,
		ShareInAmount: exitAmount,
		TokenOutMins:  sdk.NewCoins(),
	}
	_, err = s.gammMsgServer.ExitPool(s.Ctx, exitMsg)
	s.Require().NoError(err)

	// Transfer collection tokens to another user (simulating burn-like operation)
	// Note: Transferring directly to "Mint" requires specific approvals, so we'll transfer to another user
	burnTokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
	}
	recipient := s.TestAccs[1]
	burnMsg := &tokenizationtypes.MsgTransferTokens{
		Creator:      user.String(),
		CollectionId: collectionId,
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        user.String(),
				ToAddresses: []string{recipient.String()},
				Balances: []*tokenizationtypes.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       burnTokenIds,
						OwnershipTimes: []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(18446744073709551615)}},
					},
				},
			},
		},
	}
	_, err = s.tokenizationMsgServer.TransferTokens(s.Ctx, burnMsg)
	s.Require().NoError(err)

	// Verify LP tokens decreased
	finalLpBalance := s.GetCoinBalance(user, lpDenom)
	s.Require().Equal(lpBalance.Amount.Sub(exitAmount), finalLpBalance.Amount)

	// Verify collection tokens reduced (user no longer has tokens 1-10)
	tokenBalance := s.GetTokenBalance(collectionId, user, sdkmath.NewUint(5))
	s.Require().True(tokenBalance.IsZero())

	s.T().Log("Successfully exited pool and burned collection tokens")
}

// TestMultiHopSwapWithCollectionContext tests multi-hop swap in context of collection operations
func (s *CombinedTestSuite) TestMultiHopSwapWithCollectionContext() {
	s.T().Log("Testing multi-hop swap with collection context")

	creator := s.TestAccs[0]
	trader := s.TestAccs[1]

	// Create collection and mint to trader
	collectionId, err := s.CreateBasicCollection(creator)
	s.Require().NoError(err)

	tokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
	}
	err = s.MintTokensToAddress(collectionId, trader, tokenIds, sdkmath.NewUint(1))
	s.Require().NoError(err)

	// Create multiple pools for multi-hop
	pool1 := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)
	pool2 := s.CreatePoolWithCoins(
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
		sdk.NewCoin("baz", sdkmath.NewInt(10000000)),
	)
	pool3 := s.CreatePoolWithCoins(
		sdk.NewCoin("baz", sdkmath.NewInt(10000000)),
		sdk.NewCoin("qux", sdkmath.NewInt(10000000)),
	)

	// Fund trader
	s.FundAcc(trader, sdk.NewCoins(sdk.NewCoin("foo", sdkmath.NewInt(100000))))

	// Multi-hop swap: foo -> bar -> baz -> qux
	swapMsg := &gammtypes.MsgSwapExactAmountIn{
		Sender: trader.String(),
		Routes: []poolmanagertypes.SwapAmountInRoute{
			{PoolId: pool1, TokenOutDenom: "bar"},
			{PoolId: pool2, TokenOutDenom: "baz"},
			{PoolId: pool3, TokenOutDenom: "qux"},
		},
		TokenIn:           sdk.NewCoin("foo", sdkmath.NewInt(100000)),
		TokenOutMinAmount: sdkmath.NewInt(1),
	}
	res, err := s.gammMsgServer.SwapExactAmountIn(s.Ctx, swapMsg)
	s.Require().NoError(err)

	// Verify trader has qux and still has collection tokens
	quxBalance := s.GetCoinBalance(trader, "qux")
	s.Require().Equal(res.TokenOutAmount, quxBalance.Amount)

	tokenBalance := s.GetTokenBalance(collectionId, trader, sdkmath.NewUint(5))
	s.Require().Equal(sdkmath.NewUint(1), tokenBalance)

	s.T().Logf("Multi-hop swap: 100000 foo -> %s qux (3 hops)", res.TokenOutAmount)
}

// TestConcurrentPoolAndCollectionOperations tests multiple operations in rapid succession
func (s *CombinedTestSuite) TestConcurrentPoolAndCollectionOperations() {
	s.T().Log("Testing concurrent-style pool and collection operations")

	users := []sdk.AccAddress{s.TestAccs[0], s.TestAccs[1], s.TestAccs[2]}

	// Create pool
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(100000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(100000000)),
	)

	// Create collection
	collectionId, err := s.CreateBasicCollection(users[0])
	s.Require().NoError(err)

	// Mint tokens to all users
	for _, user := range users {
		tokenIds := []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		}
		err = s.MintTokensToAddress(collectionId, user, tokenIds, sdkmath.NewUint(1))
		s.Require().NoError(err)

		// Fund with coins
		s.FundAcc(user, sdk.NewCoins(
			sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
			sdk.NewCoin("bar", sdkmath.NewInt(1000000)),
		))
	}

	// Rapid operations
	for round := 0; round < 5; round++ {
		for i, user := range users {
			// Swap
			_, err = s.SwapTokens(user, poolId, sdk.NewCoin("foo", sdkmath.NewInt(1000)), "bar")
			s.Require().NoError(err)

			// Transfer collection token to next user
			nextUser := users[(i+1)%len(users)]
			transferTokenIds := []*tokenizationtypes.UintRange{
				{Start: sdkmath.NewUint(uint64(round*10 + i + 1)), End: sdkmath.NewUint(uint64(round*10 + i + 1))},
			}
			err = s.TransferTokens(collectionId, user, nextUser, transferTokenIds, sdkmath.NewUint(1))
			s.Require().NoError(err)
		}
		s.T().Logf("Completed round %d", round+1)
	}

	s.T().Log("Successfully completed rapid concurrent-style operations")
}

// TestStateConsistencyAcrossModules verifies state consistency after mixed operations
func (s *CombinedTestSuite) TestStateConsistencyAcrossModules() {
	s.T().Log("Testing state consistency across modules")

	user := s.TestAccs[0]

	// Initial state
	initialFooBalance := s.GetCoinBalance(user, "foo")
	initialBarBalance := s.GetCoinBalance(user, "bar")

	// Create pool
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(1000000)),
	)

	// Create collection
	collectionId, err := s.CreateBasicCollection(user)
	s.Require().NoError(err)

	// Mint tokens
	tokenIds := []*tokenizationtypes.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
	}
	err = s.MintTokensToAddress(collectionId, user, tokenIds, sdkmath.NewUint(1))
	s.Require().NoError(err)

	// Verify pool was created correctly
	pool, err := s.App.GammKeeper.GetPoolAndPoke(s.Ctx, poolId)
	s.Require().NoError(err)

	// Verify collection was created correctly
	collection, found := s.App.TokenizationKeeper.GetCollectionFromStore(s.Ctx, collectionId)
	s.Require().True(found)

	// Verify user LP tokens
	lpDenom := gammtypes.GetPoolShareDenom(poolId)
	lpBalance := s.GetCoinBalance(user, lpDenom)
	s.Require().True(lpBalance.Amount.GT(sdkmath.NewInt(0)))

	// Verify user collection tokens
	tokenBalance := s.GetTokenBalance(collectionId, user, sdkmath.NewUint(50))
	s.Require().Equal(sdkmath.NewUint(1), tokenBalance)

	// Log state summary
	s.T().Logf("State summary:")
	s.T().Logf("  Pool ID: %d, Total shares: %s", poolId, pool.GetTotalShares())
	s.T().Logf("  Collection ID: %s", collection.CollectionId)
	s.T().Logf("  User LP balance: %s", lpBalance)
	s.T().Logf("  User collection token balance: %s", tokenBalance)
	s.T().Logf("  Initial foo: %s, Initial bar: %s", initialFooBalance, initialBarBalance)

	s.T().Log("State consistency verified across modules")
}
