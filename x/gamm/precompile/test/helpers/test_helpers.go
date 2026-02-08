package helpers

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"

	sdkmath "cosmossdk.io/math"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/app"
	"github.com/bitbadges/bitbadgeschain/app/params"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
)

// init ensures SDK config is initialized with "bb" prefix before any tests run
// This must be called before any address operations to ensure correct Bech32 prefix
func init() {
	// Initialize SDK config with "bb" prefix if not already set
	// This is safe to call multiple times - it will only set if not already "bb"
	params.InitSDKConfigWithoutSeal()
}

// TestSuite provides common test utilities and fixtures
// Uses apptesting.KeeperTestHelper for full app integration
type TestSuite struct {
	apptesting.KeeperTestHelper

	Precompile  *gamm.Precompile
	MsgServer   gammtypes.MsgServer
	QueryClient gammtypes.QueryClient

	// Test addresses (EVM format)
	AliceEVM   common.Address
	BobEVM     common.Address
	CharlieEVM common.Address

	// Test addresses (Cosmos format)
	Alice   sdk.AccAddress
	Bob     sdk.AccAddress
	Charlie sdk.AccAddress

	// Test pool data
	PoolId uint64
}

var defaultTestStartTime = time.Now().UTC()

// NewTestSuite creates a new test suite with initialized keeper and context
// Uses apptesting.KeeperTestHelper which provides a full app instance
// t is required - it will be used for logging in Setup()
func NewTestSuite(t *testing.T) *TestSuite {
	// Ensure SDK config is initialized with "bb" prefix before any address operations
	// This must be called before creating addresses or calling keeper functions
	params.InitSDKConfigWithoutSeal()

	ts := &TestSuite{}
	
	// Set T() directly on the embedded suite
	// We need T() to be set for Setup() to work
	ts.SetT(t)
	
	// Now we can call Setup() which requires T()
	ts.App = app.Setup(false)
	
	// Manually set up the context and query helper (setupGeneral logic)
	ts.Ctx = ts.App.BaseApp.NewContextLegacy(false, cmtproto.Header{Height: 1, ChainID: "bitbadges-1", Time: defaultTestStartTime})
	ts.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: ts.App.GRPCQueryRouter(),
		Ctx:             ts.Ctx,
	}
	ts.TestAccs = []sdk.AccAddress{}
	
	precompile := gamm.NewPrecompile(ts.App.GammKeeper)
	msgServer := gammkeeper.NewMsgServerImpl(&ts.App.GammKeeper)

	// Create test EVM addresses
	aliceEVM := common.HexToAddress("0x1111111111111111111111111111111111111111")
	bobEVM := common.HexToAddress("0x2222222222222222222222222222222222222222")
	charlieEVM := common.HexToAddress("0x3333333333333333333333333333333333333333")

	// Convert to Cosmos addresses
	alice := sdk.AccAddress(aliceEVM.Bytes())
	bob := sdk.AccAddress(bobEVM.Bytes())
	charlie := sdk.AccAddress(charlieEVM.Bytes())

	ts.Precompile = precompile
	ts.MsgServer = msgServer
	ts.QueryClient = gammtypes.NewQueryClient(ts.QueryHelper)
	ts.AliceEVM = aliceEVM
	ts.BobEVM = bobEVM
	ts.CharlieEVM = charlieEVM
	ts.Alice = alice
	ts.Bob = bob
	ts.Charlie = charlie
	ts.PoolId = 0

	// Fund test accounts
	defaultCoins := sdk.NewCoins(
		sdk.NewCoin("uosmo", osmomath.NewInt(1_000_000_000_000_000_000)),
		sdk.NewCoin("uion", osmomath.NewInt(1_000_000_000_000_000_000)),
	)
	ts.FundAcc(alice, defaultCoins)
	ts.FundAcc(bob, defaultCoins)
	ts.FundAcc(charlie, defaultCoins)

	return ts
}

// CreateMockContract creates a mock EVM contract for testing
func (ts *TestSuite) CreateMockContract(caller common.Address, input []byte) *vm.Contract {
	precompileAddr := common.HexToAddress(gamm.GammPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
	if len(input) > 0 {
		contract.SetCallCode(common.Hash{}, input)
	}
	return contract
}

// CreateTestBalancerPool creates a basic balancer pool for testing
func (ts *TestSuite) CreateTestBalancerPool(creator sdk.AccAddress, poolAssets []balancer.PoolAsset) (uint64, error) {
	poolParams := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"), // 2.5% swap fee
		ExitFee: osmomath.ZeroDec(),
	}

	msg := balancer.NewMsgCreateBalancerPool(creator, poolParams, poolAssets)
	
	// Use pool manager to create pool
	poolId, err := ts.App.PoolManagerKeeper.CreatePool(ts.Ctx, msg)
	if err != nil {
		return 0, err
	}

	ts.PoolId = poolId
	return poolId, nil
}

// CreateDefaultTestPool creates a default two-asset balancer pool for testing
func (ts *TestSuite) CreateDefaultTestPool(creator sdk.AccAddress) (uint64, error) {
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

	return ts.CreateTestBalancerPool(creator, poolAssets)
}

// FundAccount funds an account with coins
// Uses the KeeperTestHelper's FundAcc method
func (ts *TestSuite) FundAccount(addr sdk.AccAddress, coins sdk.Coins) {
	ts.FundAcc(addr, coins)
}

// CreateMockMethod creates a mock abi.Method for testing
func CreateMockMethod(name string, inputs, outputs abi.Arguments) abi.Method {
	// Create a mock method ID (first 4 bytes of keccak256 hash)
	methodSig := name
	hash := common.BytesToHash([]byte(methodSig))
	methodID := hash[:4]

	return abi.Method{
		Name:    name,
		RawName: name,
		Type:    abi.Function,
		Inputs:  inputs,
		Outputs: outputs,
		ID:      methodID,
	}
}

// BigIntToInt converts *big.Int to sdkmath.Int
func BigIntToInt(bi *big.Int) sdkmath.Int {
	return sdkmath.NewIntFromBigInt(bi)
}

// IntToBigInt converts sdkmath.Int to *big.Int
func IntToBigInt(i sdkmath.Int) *big.Int {
	return i.BigInt()
}

