package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func BadgesKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	ak := accountkeeper.NewAccountKeeper(appCodec, runtime.NewKVStoreService(storeKey), func() sdk.AccountI { return &authtypes.BaseAccount{} }, map[string][]string{}, address.NewBech32Codec("bb"), "bb", authority.String())

	bankKeeper := bankkeeper.NewBaseKeeper(appCodec, runtime.NewKVStoreService(storeKey), ak, map[string]bool{}, authority.String(), log.NewNopLogger())

	dk := distributionkeeper.Keeper{}

	// Create a mock SendManagerKeeper to avoid nil pointer panics
	mockSendManagerKeeper := &mockSendManagerKeeper{
		bankKeeper: bankKeeper,
	}

	k := keeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authority.String(),
		bankKeeper,
		ak,
		dk,                    // DistributionKeeper
		mockSendManagerKeeper, // SendManagerKeeper
		nil,                   // IBCKeeperFn
		nil,                   // CapabilityScopedFn
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, ctx
}

// mockSendManagerKeeper is a minimal mock implementation of SendManagerKeeper for testing
type mockSendManagerKeeper struct {
	bankKeeper bankkeeper.Keeper
}

func (m *mockSendManagerKeeper) SendCoinWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, toAddressAcc sdk.AccAddress, coin *sdk.Coin) error {
	// For tests, just use bank keeper directly
	// BankKeeper.SendCoins uses context.Context, so we need to wrap the sdk.Context
	return m.bankKeeper.SendCoins(sdk.WrapSDKContext(ctx), fromAddressAcc, toAddressAcc, sdk.NewCoins(*coin))
}

func (m *mockSendManagerKeeper) SendCoinsWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, toAddressAcc sdk.AccAddress, coins sdk.Coins) error {
	// For tests, just use bank keeper directly
	// BankKeeper.SendCoins uses context.Context, so we need to wrap the sdk.Context
	return m.bankKeeper.SendCoins(sdk.WrapSDKContext(ctx), fromAddressAcc, toAddressAcc, coins)
}

func (m *mockSendManagerKeeper) FundCommunityPoolWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, coins sdk.Coins) error {
	// For tests, we can't easily mock distribution keeper, so just return nil
	// This is acceptable for unit tests that don't test community pool funding
	return nil
}

func (m *mockSendManagerKeeper) GetBalanceWithAliasRouting(ctx sdk.Context, address sdk.AccAddress, denom string) (sdk.Coin, error) {
	// For tests, just use bank keeper directly
	// BankKeeper.GetBalance uses context.Context, so we need to wrap the sdk.Context
	return m.bankKeeper.GetBalance(sdk.WrapSDKContext(ctx), address, denom), nil
}

func (m *mockSendManagerKeeper) IsICS20Compatible(ctx sdk.Context, denom string) bool {
	// For tests, assume all denoms are ICS20 compatible
	return true
}

func (m *mockSendManagerKeeper) StandardName(ctx sdk.Context, denom string) string {
	// For tests, return "x/bank" for all denoms
	return "x/bank"
}

func (m *mockSendManagerKeeper) SendCoinsFromModuleToAccountWithAliasRouting(ctx sdk.Context, moduleName string, toAddressAcc sdk.AccAddress, coins sdk.Coins) error {
	// For tests, just use bank keeper directly
	// BankKeeper.SendCoinsFromModuleToAccount uses context.Context, so we need to wrap the sdk.Context
	return m.bankKeeper.SendCoinsFromModuleToAccount(sdk.WrapSDKContext(ctx), moduleName, toAddressAcc, coins)
}

func (m *mockSendManagerKeeper) SendCoinsFromAccountToModuleWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, moduleName string, coins sdk.Coins) error {
	// For tests, just use bank keeper directly
	// BankKeeper.SendCoinsFromAccountToModule uses context.Context, so we need to wrap the sdk.Context
	return m.bankKeeper.SendCoinsFromAccountToModule(sdk.WrapSDKContext(ctx), fromAddressAcc, moduleName, coins)
}

func (m *mockSendManagerKeeper) SpendFromCommunityPoolWithAliasRouting(ctx sdk.Context, toAddressAcc sdk.AccAddress, coins sdk.Coins) error {
	// For tests, just use bank keeper directly with distribution module
	// BankKeeper.SendCoinsFromModuleToAccount uses context.Context, so we need to wrap the sdk.Context
	return m.bankKeeper.SendCoinsFromModuleToAccount(sdk.WrapSDKContext(ctx), "distribution", toAddressAcc, coins)
}
