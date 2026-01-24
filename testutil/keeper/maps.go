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
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/maps/keeper"
	"github.com/bitbadges/bitbadgeschain/x/maps/types"

	badgeskeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
)

func MapsKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(registry)

	// IBC v10: capabilities removed, portKeeper doesn't need scopedKeeper
	// Create a minimal IBCKeeper for testing - PortKeeper is created internally
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	ibcK := ibckeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(storeKey),
		NewNoopParamSubspace(), // ParamSubspace - no-op for testing
		NewNoopUpgradeKeeper(), // UpgradeKeeper - no-op for testing
		authority.String(), // authority - required, use gov module address
	)

	k := keeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(storeKey),
		log.NewNopLogger(),
		"",
		func() *ibckeeper.Keeper {
			return ibcK
		},
		badgeskeeper.Keeper{},
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, ctx
}
