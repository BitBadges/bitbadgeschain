package testutil

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/app/params"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	module "github.com/bitbadges/bitbadgeschain/x/sendmanager/module"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

// AITestSuite provides a comprehensive test suite for AI-generated tests
type AITestSuite struct {
	suite.Suite

	Keeper      keeper.Keeper
	Ctx         sdk.Context
	MsgServer   types.MsgServer
	AddressCodec address.Codec
	MockBank    *MockBankKeeper

	// Test addresses
	Alice   string
	Bob     string
	Charlie string
	Authority string
}

// SetupTest initializes the test suite with a fresh keeper and context
func (suite *AITestSuite) SetupTest() {
	// Ensure SDK config is initialized with "bb" prefix before it gets sealed
	params.InitSDKConfigWithoutSeal()
	
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec("bb")
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(suite.T(), storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)
	authorityBytes, err := addressCodec.StringToBytes(authority.String())
	suite.Require().NoError(err, "failed to convert authority to bytes")

	// Create mock keepers
	mockBankKeeper := &MockBankKeeper{}
	mockDistributionKeeper := &MockDistributionKeeper{}

	k := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		addressCodec,
		authorityBytes,
		mockBankKeeper,
		mockDistributionKeeper,
	)

	// Initialize params
	err = k.SetParams(ctx, types.DefaultParams())
	suite.Require().NoError(err, "failed to set params")

	suite.Keeper = k
	suite.Ctx = ctx
	suite.MsgServer = keeper.NewMsgServerImpl(k)
	suite.MockBank = mockBankKeeper
	suite.AddressCodec = addressCodec

	// Initialize test addresses
	suite.Alice = "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
	suite.Bob = "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	suite.Charlie = "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme"
	suite.Authority = authority.String()
}

// MockBankKeeper is a mock implementation of BankKeeper for testing
type MockBankKeeper struct {
	balances map[string]map[string]sdk.Coin // address -> denom -> coin
}

func NewMockBankKeeper() *MockBankKeeper {
	return &MockBankKeeper{
		balances: make(map[string]map[string]sdk.Coin),
	}
}

func (m *MockBankKeeper) GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	allBalances := m.GetAllBalances(ctx, addr)
	amount := allBalances.AmountOf(denom)
	return sdk.Coin{Denom: denom, Amount: amount}
}

func (m *MockBankKeeper) SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	return m.GetAllBalances(ctx, addr)
}

func (m *MockBankKeeper) SendCoins(ctx context.Context, from, to sdk.AccAddress, coins sdk.Coins) error {
	// Simple mock implementation - just track balances
	if m.balances == nil {
		m.balances = make(map[string]map[string]sdk.Coin)
	}
	
	fromStr := from.String()
	toStr := to.String()
	
	// Deduct from sender
	if m.balances[fromStr] == nil {
		m.balances[fromStr] = make(map[string]sdk.Coin)
	}
		for _, coin := range coins {
			if existing, ok := m.balances[fromStr][coin.Denom]; ok {
				newAmount := existing.Amount.Sub(coin.Amount)
				if newAmount.IsNegative() {
					return fmt.Errorf("insufficient funds: %s", coin.String())
				}
				if newAmount.IsZero() {
					delete(m.balances[fromStr], coin.Denom)
				} else {
					m.balances[fromStr][coin.Denom] = sdk.NewCoin(coin.Denom, newAmount)
				}
			} else {
				// No existing balance, but trying to send - insufficient funds
				return fmt.Errorf("insufficient funds: %s", coin.String())
			}
		}
	
	// Add to recipient
	if m.balances[toStr] == nil {
		m.balances[toStr] = make(map[string]sdk.Coin)
	}
	for _, coin := range coins {
		if existing, ok := m.balances[toStr][coin.Denom]; ok {
			m.balances[toStr][coin.Denom] = sdk.NewCoin(coin.Denom, existing.Amount.Add(coin.Amount))
		} else {
			m.balances[toStr][coin.Denom] = coin
		}
	}
	
	return nil
}

func (m *MockBankKeeper) MintCoins(ctx context.Context, moduleName string, coins sdk.Coins) error {
	return nil
}

func (m *MockBankKeeper) BurnCoins(ctx context.Context, moduleName string, coins sdk.Coins) error {
	return nil
}

func (m *MockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, moduleName string, addr sdk.AccAddress, coins sdk.Coins) error {
	return m.SendCoins(ctx, sdk.AccAddress(authtypes.NewModuleAddress(moduleName)), addr, coins)
}

func (m *MockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
	return m.SendCoins(ctx, addr, sdk.AccAddress(authtypes.NewModuleAddress(moduleName)), coins)
}

func (m *MockBankKeeper) GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	if m.balances == nil {
		return sdk.Coins{}
	}
	
	addrStr := addr.String()
	if balances, ok := m.balances[addrStr]; ok {
		coins := sdk.Coins{}
		for _, coin := range balances {
			coins = coins.Add(coin)
		}
		return coins
	}
	return sdk.Coins{}
}

// SetBalance sets a balance for testing
func (m *MockBankKeeper) SetBalance(addr sdk.AccAddress, coin sdk.Coin) {
	if m.balances == nil {
		m.balances = make(map[string]map[string]sdk.Coin)
	}
	addrStr := addr.String()
	if m.balances[addrStr] == nil {
		m.balances[addrStr] = make(map[string]sdk.Coin)
	}
	m.balances[addrStr][coin.Denom] = coin
}

// MockDistributionKeeper is a mock implementation of DistributionKeeper for testing
type MockDistributionKeeper struct{}

func (m *MockDistributionKeeper) FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	return nil
}

// RunTestSuite runs a test suite
func RunTestSuite(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}

