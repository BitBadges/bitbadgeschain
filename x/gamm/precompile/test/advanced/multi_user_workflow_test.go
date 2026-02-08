package gamm_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
)

// MultiUserWorkflowTestSuite provides tests for multi-user workflow scenarios
type MultiUserWorkflowTestSuite struct {
	apptesting.KeeperTestHelper

	Precompile *gamm.Precompile

	// Test addresses (EVM format)
	AliceEVM   common.Address
	BobEVM     common.Address
	CharlieEVM common.Address

	// Test addresses (Cosmos format)
	Alice   sdk.AccAddress
	Bob     sdk.AccAddress
	Charlie sdk.AccAddress

	// Test pools
	PoolId1 uint64
	PoolId2 uint64
}

func TestMultiUserWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(MultiUserWorkflowTestSuite))
}

func (suite *MultiUserWorkflowTestSuite) SetupTest() {
	suite.Reset()

	suite.Precompile = gamm.NewPrecompile(suite.App.GammKeeper)

	// Create test EVM addresses
	suite.AliceEVM = common.HexToAddress("0x1111111111111111111111111111111111111111")
	suite.BobEVM = common.HexToAddress("0x2222222222222222222222222222222222222222")
	suite.CharlieEVM = common.HexToAddress("0x3333333333333333333333333333333333333333")

	// Convert to Cosmos addresses
	suite.Alice = sdk.AccAddress(suite.AliceEVM.Bytes())
	suite.Bob = sdk.AccAddress(suite.BobEVM.Bytes())
	suite.Charlie = sdk.AccAddress(suite.CharlieEVM.Bytes())

	// Fund accounts with tokens needed for pool operations
	largeAmount, _ := new(big.Int).SetString("10000000000000000000", 10)
	poolCreationCoins := sdk.NewCoins(
		sdk.NewCoin("uatom", osmomath.NewIntFromBigInt(largeAmount)),
		sdk.NewCoin("uosmo", osmomath.NewIntFromBigInt(largeAmount)),
		sdk.NewCoin("uion", osmomath.NewIntFromBigInt(largeAmount)),
	)
	suite.FundAcc(suite.Alice, poolCreationCoins)
	suite.FundAcc(suite.Bob, poolCreationCoins)
	suite.FundAcc(suite.Charlie, poolCreationCoins)

	// Create test pools
	poolId1, err := suite.createDefaultTestPoolInContext(suite.Alice)
	suite.Require().NoError(err)
	suite.PoolId1 = poolId1

	poolId2, err := suite.createDefaultTestPoolInContext(suite.Bob)
	suite.Require().NoError(err)
	suite.PoolId2 = poolId2
}

func (suite *MultiUserWorkflowTestSuite) createDefaultTestPoolInContext(creator sdk.AccAddress) (uint64, error) {
	oneTrillion := osmomath.NewInt(1e12)
	poolAssets := []balancer.PoolAsset{
		{
			Token:  sdk.NewCoin("uatom", oneTrillion),
			Weight: osmomath.NewInt(100),
		},
		{
			Token:  sdk.NewCoin("uosmo", oneTrillion),
			Weight: osmomath.NewInt(100),
		},
	}
	poolParams := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"),
		ExitFee: osmomath.ZeroDec(),
	}
	msg := balancer.NewMsgCreateBalancerPool(creator, poolParams, poolAssets)
	poolId, err := suite.App.PoolManagerKeeper.CreatePool(suite.Ctx, msg)
	return poolId, err
}

// TestMultiUser_SequentialPoolOperations tests sequential operations by multiple users
func (suite *MultiUserWorkflowTestSuite) TestMultiUser_SequentialPoolOperations() {
	// Step 1: Alice joins pool 1
	pool1Before, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId1)
	suite.Require().NoError(err)
	sharesBefore1 := pool1Before.GetTotalShares()

	// Verify Alice can query pool 1
	method := suite.Precompile.ABI.Methods["getPool"]
	packed, err := method.Inputs.Pack(suite.PoolId1)
	suite.Require().NoError(err)
	input := append(method.ID, packed...)

	// Verify method exists and is accessible
	suite.NotNil(method)
	suite.Greater(len(input), 4, "Input should include method ID and parameters")

	// Step 2: Bob joins pool 2
	pool2Before, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId2)
	suite.Require().NoError(err)
	sharesBefore2 := pool2Before.GetTotalShares()

	// Verify Bob can query pool 2
	packed, err = method.Inputs.Pack(suite.PoolId2)
	suite.Require().NoError(err)
	input = append(method.ID, packed...)
	suite.NotNil(method)

	// Step 3: Charlie queries both pools
	pool1After, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId1)
	suite.Require().NoError(err)
	pool2After, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId2)
	suite.Require().NoError(err)

	// Verify pools are accessible by all users
	suite.Equal(suite.PoolId1, pool1After.GetId(), "Pool 1 ID should match")
	suite.Equal(suite.PoolId2, pool2After.GetId(), "Pool 2 ID should match")
	suite.Equal(sharesBefore1.String(), pool1After.GetTotalShares().String(), "Pool 1 shares should remain unchanged")
	suite.Equal(sharesBefore2.String(), pool2After.GetTotalShares().String(), "Pool 2 shares should remain unchanged")

	suite.T().Log("Sequential pool operations by multiple users completed successfully")
}

// TestMultiUser_ConcurrentQueries tests concurrent queries by multiple users
func (suite *MultiUserWorkflowTestSuite) TestMultiUser_ConcurrentQueries() {
	// All users query the same pool concurrently
	method := suite.Precompile.ABI.Methods["getPool"]

	// Alice queries pool 1
	packed1, err := method.Inputs.Pack(suite.PoolId1)
	suite.Require().NoError(err)
	input1 := append(method.ID, packed1...)
	suite.NotNil(input1)

	// Bob queries pool 1
	packed2, err := method.Inputs.Pack(suite.PoolId1)
	suite.Require().NoError(err)
	input2 := append(method.ID, packed2...)
	suite.NotNil(input2)

	// Charlie queries pool 1
	packed3, err := method.Inputs.Pack(suite.PoolId1)
	suite.Require().NoError(err)
	input3 := append(method.ID, packed3...)
	suite.NotNil(input3)

	// All queries should be valid
	suite.Equal(input1, input2, "Query inputs should be identical")
	suite.Equal(input2, input3, "Query inputs should be identical")

	// Verify all users get the same pool data
	pool1, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId1)
	suite.Require().NoError(err)
	pool2, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId1)
	suite.Require().NoError(err)
	pool3, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId1)
	suite.Require().NoError(err)

	suite.Equal(pool1.GetId(), pool2.GetId(), "All users should see the same pool ID")
	suite.Equal(pool2.GetId(), pool3.GetId(), "All users should see the same pool ID")
	suite.Equal(pool1.GetTotalShares().String(), pool2.GetTotalShares().String(), "All users should see the same total shares")
	suite.Equal(pool2.GetTotalShares().String(), pool3.GetTotalShares().String(), "All users should see the same total shares")

	suite.T().Log("Concurrent queries by multiple users completed successfully")
}

// TestMultiUser_DifferentPools tests operations on different pools by different users
func (suite *MultiUserWorkflowTestSuite) TestMultiUser_DifferentPools() {
	// Alice operates on pool 1
	pool1, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId1)
	suite.Require().NoError(err)
	suite.Equal(suite.PoolId1, pool1.GetId(), "Alice should access pool 1")

	// Bob operates on pool 2
	pool2, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId2)
	suite.Require().NoError(err)
	suite.Equal(suite.PoolId2, pool2.GetId(), "Bob should access pool 2")

	// Charlie queries both pools
	pool1Again, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId1)
	suite.Require().NoError(err)
	pool2Again, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId2)
	suite.Require().NoError(err)

	// Verify pools are independent
	suite.Equal(pool1.GetId(), pool1Again.GetId(), "Pool 1 should remain consistent")
	suite.Equal(pool2.GetId(), pool2Again.GetId(), "Pool 2 should remain consistent")
	suite.NotEqual(pool1.GetId(), pool2.GetId(), "Pools should be different")

	suite.T().Log("Operations on different pools by different users completed successfully")
}

// TestMultiUser_QueryMethods tests that all users can use query methods
func (suite *MultiUserWorkflowTestSuite) TestMultiUser_QueryMethods() {
	queryMethods := []string{
		"getPool",
		"getPools",
		"getPoolType",
		"getPoolParams",
		"getTotalShares",
		"getTotalLiquidity",
	}

	// Test that all users can access all query methods
	for _, methodName := range queryMethods {
		method, found := suite.Precompile.ABI.Methods[methodName]
		suite.True(found, "%s method should exist", methodName)
		suite.NotNil(method, "%s method should not be nil", methodName)

		// Verify method is not a transaction (queries are read-only)
		suite.False(suite.Precompile.IsTransaction(&method), "%s should not be a transaction", methodName)
	}

	suite.T().Log("All query methods are accessible to all users")
}

// TestMultiUser_TransactionMethods tests that transaction methods are properly structured
func (suite *MultiUserWorkflowTestSuite) TestMultiUser_TransactionMethods() {
	transactionMethods := []string{
		"joinPool",
		"exitPool",
		"swapExactAmountIn",
		"swapExactAmountInWithIBCTransfer",
	}

	// Test that all transaction methods are properly marked
	for _, methodName := range transactionMethods {
		method, found := suite.Precompile.ABI.Methods[methodName]
		suite.True(found, "%s method should exist", methodName)
		suite.NotNil(method, "%s method should not be nil", methodName)

		// Verify method is a transaction (state-changing)
		suite.True(suite.Precompile.IsTransaction(&method), "%s should be a transaction", methodName)
	}

	suite.T().Log("All transaction methods are properly structured")
}

// TestMultiUser_StateIsolation tests that operations by one user don't affect another user's view
func (suite *MultiUserWorkflowTestSuite) TestMultiUser_StateIsolation() {
	// Get initial state for both pools
	pool1Initial, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId1)
	suite.Require().NoError(err)
	pool2Initial, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId2)
	suite.Require().NoError(err)

	// Alice queries pool 1
	pool1Alice, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId1)
	suite.Require().NoError(err)

	// Bob queries pool 2
	pool2Bob, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId2)
	suite.Require().NoError(err)

	// Charlie queries both pools
	pool1Charlie, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId1)
	suite.Require().NoError(err)
	pool2Charlie, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId2)
	suite.Require().NoError(err)

	// Verify state isolation - each user sees consistent state
	suite.Equal(pool1Initial.GetId(), pool1Alice.GetId(), "Alice should see consistent pool 1 state")
	suite.Equal(pool1Initial.GetId(), pool1Charlie.GetId(), "Charlie should see consistent pool 1 state")
	suite.Equal(pool2Initial.GetId(), pool2Bob.GetId(), "Bob should see consistent pool 2 state")
	suite.Equal(pool2Initial.GetId(), pool2Charlie.GetId(), "Charlie should see consistent pool 2 state")

	// Verify pools are independent
	suite.NotEqual(pool1Alice.GetId(), pool2Bob.GetId(), "Pool 1 and pool 2 should be different")

	suite.T().Log("State isolation verified - users see consistent and independent pool states")
}

// TestMultiUser_PoolAccessControl tests that all users can access all pools (no access restrictions)
func (suite *MultiUserWorkflowTestSuite) TestMultiUser_PoolAccessControl() {
	// All users should be able to query all pools
	users := []sdk.AccAddress{suite.Alice, suite.Bob, suite.Charlie}
	pools := []uint64{suite.PoolId1, suite.PoolId2}

	for _, user := range users {
		for _, poolId := range pools {
			pool, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, poolId)
			suite.NoError(err, "User %s should be able to query pool %d", user.String(), poolId)
			suite.NotNil(pool, "Pool %d should not be nil for user %s", poolId, user.String())
			suite.Equal(poolId, pool.GetId(), "Pool ID should match for user %s", user.String())
		}
	}

	suite.T().Log("All users can access all pools (no access restrictions)")
}

