package gamm_test

import (
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	"github.com/bitbadges/bitbadgeschain/x/gamm/precompile/test/helpers"
)

// ConcurrencyTestSuite provides tests for parallel execution safety
type ConcurrencyTestSuite struct {
	apptesting.KeeperTestHelper

	Precompile *gamm.Precompile
	PoolId    uint64
}

func TestConcurrencyTestSuite(t *testing.T) {
	suite.Run(t, new(ConcurrencyTestSuite))
}

func (suite *ConcurrencyTestSuite) SetupTest() {
	suite.Reset()
	suite.Precompile = gamm.NewPrecompile(suite.App.GammKeeper)

	// Create a test pool
	alice := suite.TestAccs[0]
	largeAmount, _ := new(big.Int).SetString("10000000000000000000", 10)
	poolCreationCoins := sdk.NewCoins(
		sdk.NewCoin("uatom", osmomath.NewIntFromBigInt(largeAmount)), // Large amount for concurrent operations
		sdk.NewCoin("uosmo", osmomath.NewIntFromBigInt(largeAmount)),
	)
	suite.FundAcc(alice, poolCreationCoins)

	oneTrillion := osmomath.NewInt(1e12)
	poolAssets := []balancer.PoolAsset{
		{Token: sdk.NewCoin("uatom", oneTrillion), Weight: osmomath.NewInt(100)},
		{Token: sdk.NewCoin("uosmo", oneTrillion), Weight: osmomath.NewInt(100)},
	}

	poolParams := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"),
		ExitFee: osmomath.ZeroDec(),
	}

	msg := balancer.NewMsgCreateBalancerPool(alice, poolParams, poolAssets)
	poolId, err := suite.App.PoolManagerKeeper.CreatePool(suite.Ctx, msg)
	suite.Require().NoError(err)
	suite.PoolId = poolId
}

// TestConcurrency_ParallelQueries tests parallel query execution
// Queries are read-only and should be safe for concurrent access
func (suite *ConcurrencyTestSuite) TestConcurrency_ParallelQueries() {
	const numGoroutines = 10
	const numQueries = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*numQueries)

	// Run queries in parallel
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numQueries; j++ {
				// Test GetPool query
				method, found := suite.Precompile.ABI.Methods["getPool"]
				if !found {
					errors <- fmt.Errorf("getPool method not found")
					return
				}

				// Build JSON query
				queryJson, err := helpers.BuildGetPoolQueryJSON(suite.PoolId)
				if err != nil {
					errors <- err
					continue
				}

				// Pack method with JSON string
				input, err := helpers.PackMethodWithJSON(&method, queryJson)
				if err != nil {
					errors <- err
					continue
				}

				// Verify method exists (this is safe for concurrent access)
				gas := suite.Precompile.RequiredGas(input)
				if gas == 0 {
					errors <- fmt.Errorf("Gas should be greater than 0")
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
			suite.T().Logf("Concurrent query error: %v", err)
		}
	}

	suite.Equal(0, errorCount, "No errors should occur during concurrent queries")
}

// TestConcurrency_StateConsistency verifies state remains consistent
// This test verifies that concurrent operations don't corrupt state
func (suite *ConcurrencyTestSuite) TestConcurrency_StateConsistency() {
	// Get initial pool state
	poolBefore, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId)
	suite.Require().NoError(err)
	totalSharesBefore := poolBefore.GetTotalShares()

	// Perform multiple concurrent queries (read-only operations)
	const numGoroutines = 20
	var wg sync.WaitGroup
	var mu sync.Mutex
	shareCounts := make([]sdkmath.Int, 0, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId)
			if err != nil {
				return
			}
			mu.Lock()
			shareCounts = append(shareCounts, pool.GetTotalShares())
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Verify all queries returned the same value (state consistency)
	for _, shares := range shareCounts {
		suite.Equal(totalSharesBefore, shares, "All concurrent queries should return the same state")
	}

	// Verify final state matches initial state (no mutations from queries)
	poolAfter, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId)
	suite.Require().NoError(err)
	suite.Equal(totalSharesBefore, poolAfter.GetTotalShares(), "State should remain unchanged after read-only operations")
}

// TestConcurrency_PrecompileStructure tests that precompile structure is thread-safe
func (suite *ConcurrencyTestSuite) TestConcurrency_PrecompileStructure() {
	const numGoroutines = 50

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Access precompile structure concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Read precompile properties (should be safe)
			addr := suite.Precompile.ContractAddress
			if addr == (common.Address{}) {
				errors <- fmt.Errorf("Precompile address should not be zero")
				return
			}

			// Access ABI methods (read-only)
			method, found := suite.Precompile.ABI.Methods["getPool"]
			if !found {
				errors <- fmt.Errorf("getPool method should exist")
				return
			}

			// Verify method properties
			if method.Name != "getPool" {
				errors <- fmt.Errorf("Method name should be getPool")
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
		}
	}

	suite.Equal(0, errorCount, "No errors should occur during concurrent precompile access")
}

// TestConcurrency_ValidationFunctions tests concurrent validation
// Validation functions should be thread-safe (no shared mutable state)
func (suite *ConcurrencyTestSuite) TestConcurrency_ValidationFunctions() {
	const numGoroutines = 20

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Run validation functions concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Test ValidatePoolId
			err := gamm.ValidatePoolId(suite.PoolId)
			if err != nil {
				errors <- err
				return
			}

			// Test ValidateShareAmount (requires > 0)
			err = gamm.ValidateShareAmount(big.NewInt(int64(1000+id)), "shareAmount")
			if err != nil {
				errors <- err
				return
			}

			// Test CheckOverflow
			err = gamm.CheckOverflow(big.NewInt(int64(1000+id)), "amount")
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			errorCount++
			suite.T().Logf("Validation error: %v", err)
		}
	}

	suite.Equal(0, errorCount, "No errors should occur during concurrent validation")
}

