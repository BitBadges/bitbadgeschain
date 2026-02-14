package helpers

import (
	"encoding/json"
	"math/big"
	"testing"

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
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	sendmanagertypes "github.com/bitbadges/bitbadgeschain/x/sendmanager/types"

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

	Precompile  *sendmanager.Precompile
	MsgServer   sendmanagertypes.MsgServer
	QueryClient sendmanagertypes.QueryClient

	// Test addresses (EVM format)
	AliceEVM   common.Address
	BobEVM     common.Address
	CharlieEVM common.Address

	// Test addresses (Cosmos format)
	Alice   sdk.AccAddress
	Bob     sdk.AccAddress
	Charlie sdk.AccAddress
}

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
	ts.Ctx = ts.App.BaseApp.NewContextLegacy(false, cmtproto.Header{Height: 1, ChainID: "bitbadges-1"})
	ts.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: ts.App.GRPCQueryRouter(),
		Ctx:             ts.Ctx,
	}
	ts.TestAccs = []sdk.AccAddress{}

	precompile := sendmanager.NewPrecompile(ts.App.SendmanagerKeeper)
	msgServer := sendmanagerkeeper.NewMsgServerImpl(ts.App.SendmanagerKeeper)

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
	ts.QueryClient = sendmanagertypes.NewQueryClient(ts.QueryHelper)
	ts.AliceEVM = aliceEVM
	ts.BobEVM = bobEVM
	ts.CharlieEVM = charlieEVM
	ts.Alice = alice
	ts.Bob = bob
	ts.Charlie = charlie

	// Fund test accounts
	defaultCoins := sdk.NewCoins(
		sdk.NewCoin("stake", sdkmath.NewInt(1_000_000_000_000)),
		sdk.NewCoin("uosmo", sdkmath.NewInt(1_000_000_000_000)),
	)
	ts.FundAcc(alice, defaultCoins)
	ts.FundAcc(bob, defaultCoins)
	ts.FundAcc(charlie, defaultCoins)

	return ts
}

// CreateMockContract creates a mock EVM contract for testing
func (ts *TestSuite) CreateMockContract(caller common.Address, input []byte) *vm.Contract {
	precompileAddr := common.HexToAddress(sendmanager.SendManagerPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
	if len(input) > 0 {
		// SetupABI reads from contract.Input, so we need to set it explicitly
		contract.Input = input
		contract.SetCallCode(common.Hash{}, input)
	}
	return contract
}

// PackMethodWithJSON packs a method call with a JSON string argument
func PackMethodWithJSON(method *abi.Method, jsonStr string) ([]byte, error) {
	return method.Inputs.Pack(jsonStr)
}

// PackMethodCall packs a method ID with its arguments for precompile calls
func PackMethodCall(method *abi.Method, jsonStr string) ([]byte, error) {
	args, err := method.Inputs.Pack(jsonStr)
	if err != nil {
		return nil, err
	}
	// Prepend method ID (first 4 bytes)
	methodID := method.ID
	return append(methodID[:], args...), nil
}

// BuildQueryJSON builds a JSON string from a map
func BuildQueryJSON(data map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// BuildSendJSON builds a JSON string for MsgSendWithAliasRouting
// amount should be a Coins array: [{"denom": "stake", "amount": "1000"}]
func BuildSendJSON(fromAddress string, toAddress string, amount string, denom string) (string, error) {
	msg := map[string]interface{}{
		"from_address": fromAddress,
		"to_address":   toAddress,
		"amount": []map[string]interface{}{
			{
				"denom":  denom,
				"amount": amount,
			},
		},
	}
	return BuildQueryJSON(msg)
}
