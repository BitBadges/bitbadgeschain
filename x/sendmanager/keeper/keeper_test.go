package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	module "github.com/bitbadges/bitbadgeschain/x/sendmanager/module"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	authority := authtypes.NewModuleAddress(types.GovModuleName)
	authorityBytes, err := addressCodec.StringToBytes(authority.String())
	if err != nil {
		t.Fatalf("failed to convert authority to bytes: %v", err)
	}

	// Create a mock BankKeeper
	mockBankKeeper := &mockBankKeeper{}

	// Create a mock DistributionKeeper
	mockDistributionKeeper := &mockDistributionKeeper{}

	k := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		addressCodec,
		authorityBytes,
		mockBankKeeper,
		mockDistributionKeeper,
	)

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          ctx,
		keeper:       k,
		addressCodec: addressCodec,
	}
}

// mockBankKeeper is a minimal mock implementation of BankKeeper for testing
type mockBankKeeper struct{}

func (m *mockBankKeeper) GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	allBalances := m.GetAllBalances(ctx, addr)
	amount := allBalances.AmountOf(denom)
	return sdk.Coin{Denom: denom, Amount: amount}
}

func (m *mockBankKeeper) SpendableCoins(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	return sdk.Coins{}
}

func (m *mockBankKeeper) SendCoins(ctx context.Context, from, to sdk.AccAddress, coins sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) MintCoins(ctx context.Context, moduleName string, coins sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) BurnCoins(ctx context.Context, moduleName string, coins sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, moduleName string, addr sdk.AccAddress, coins sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
	return nil
}

func (m *mockBankKeeper) GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	return sdk.Coins{}
}

// mockDistributionKeeper is a minimal mock implementation of DistributionKeeper for testing
type mockDistributionKeeper struct{}

func (m *mockDistributionKeeper) FundCommunityPool(ctx context.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	return nil
}
