package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
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

	"github.com/bitbadges/bitbadgeschain/app/params"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
)

func TokenizationKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	// Add transient store for custom-hooks module (needed for deterministic error handling)
	transientStoreKey := customhookstypes.TransientStoreKey

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(registry)

	// Ensure SDK config is initialized with "bb" prefix before it gets sealed
	// This must be called before any address validation happens
	params.InitSDKConfigWithoutSeal()

	// Use bech32 codec with "bb" prefix
	bech32Codec := address.NewBech32Codec("bb")

	authorityAddr := authtypes.NewModuleAddress(govtypes.ModuleName)

	// Convert authority to "bb" prefix to match account keeper setup
	authorityStr, err := bech32Codec.BytesToString(authorityAddr)
	if err != nil {
		panic(err)
	}

	ak := accountkeeper.NewAccountKeeper(appCodec, runtime.NewKVStoreService(storeKey), func() sdk.AccountI { return &authtypes.BaseAccount{} }, map[string][]string{}, bech32Codec, "bb", authorityStr)

	bankKeeper := bankkeeper.NewBaseKeeper(appCodec, runtime.NewKVStoreService(storeKey), ak, map[string]bool{}, authorityStr, log.NewNopLogger())

	dk := distributionkeeper.Keeper{}

	// Create a mock SendManagerKeeper to avoid nil pointer panics
	mockSendManagerKeeper := &mockSendManagerKeeper{
		bankKeeper: bankKeeper,
	}

	k := keeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		authorityStr,
		bankKeeper,
		ak,
		dk,                    // DistributionKeeper
		mockSendManagerKeeper, // SendManagerKeeper
		nil,                   // IBCKeeperFn (IBC v10: capabilities removed)
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	// Initialize next collection ID to 1 (first collection will get ID 1)
	// This matches the behavior in InitGenesis
	k.SetNextCollectionId(ctx, math.NewUint(1))

	return &k, ctx
}

// mockSendManagerKeeper is a minimal mock implementation of SendManagerKeeper for testing
type mockSendManagerKeeper struct {
	bankKeeper bankkeeper.Keeper
}

func (m *mockSendManagerKeeper) SendCoinWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, toAddressAcc sdk.AccAddress, coin *sdk.Coin) error {
	// For tests, just use bank keeper directly
	return m.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(*coin))
}

func (m *mockSendManagerKeeper) SendCoinsWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, toAddressAcc sdk.AccAddress, coins sdk.Coins) error {
	// For tests, just use bank keeper directly
	// BankKeeper.SendCoins uses context.Context, so we need to wrap the sdk.Context
	return m.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, coins)
}

func (m *mockSendManagerKeeper) FundCommunityPoolWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, coins sdk.Coins) error {
	// For tests, we can't easily mock distribution keeper, so just return nil
	// This is acceptable for unit tests that don't test community pool funding
	return nil
}

func (m *mockSendManagerKeeper) GetBalanceWithAliasRouting(ctx sdk.Context, address sdk.AccAddress, denom string) (sdk.Coin, error) {
	// For tests, just use bank keeper directly
	// BankKeeper.GetBalance uses context.Context, so we need to wrap the sdk.Context
	return m.bankKeeper.GetBalance(ctx, address, denom), nil
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
	return m.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleName, toAddressAcc, coins)
}

func (m *mockSendManagerKeeper) SendCoinsFromAccountToModuleWithAliasRouting(ctx sdk.Context, fromAddressAcc sdk.AccAddress, moduleName string, coins sdk.Coins) error {
	// For tests, just use bank keeper directly
	// BankKeeper.SendCoinsFromAccountToModule uses context.Context, so we need to wrap the sdk.Context
	return m.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAddressAcc, moduleName, coins)
}

func (m *mockSendManagerKeeper) SpendFromCommunityPoolWithAliasRouting(ctx sdk.Context, toAddressAcc sdk.AccAddress, coins sdk.Coins) error {
	// For tests, just use bank keeper directly with distribution module
	// BankKeeper.SendCoinsFromModuleToAccount uses context.Context, so we need to wrap the sdk.Context
	return m.bankKeeper.SendCoinsFromModuleToAccount(ctx, "distribution", toAddressAcc, coins)
}
