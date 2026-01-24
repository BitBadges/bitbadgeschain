package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	storetypes "cosmossdk.io/store/types"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibctypes "github.com/cosmos/ibc-go/v10/modules/core/types"
)

// noopParamSubspace is a no-op implementation of ParamSubspace for testing
type noopParamSubspace struct{}

var _ ibctypes.ParamSubspace = (*noopParamSubspace)(nil)

// GetParamSet implements ibctypes.ParamSubspace
func (n *noopParamSubspace) GetParamSet(ctx sdk.Context, ps paramtypes.ParamSet) {
	// No-op
}

// NewNoopParamSubspace creates a new no-op ParamSubspace
func NewNoopParamSubspace() ibctypes.ParamSubspace {
	return &noopParamSubspace{}
}

// noopUpgradeKeeper is a no-op implementation of UpgradeKeeper for testing
type noopUpgradeKeeper struct{}

var _ ibcclienttypes.UpgradeKeeper = (*noopUpgradeKeeper)(nil)

// GetUpgradePlan implements ibcclienttypes.UpgradeKeeper
func (n *noopUpgradeKeeper) GetUpgradePlan(ctx context.Context) (upgradetypes.Plan, error) {
	return upgradetypes.Plan{}, nil
}

// GetUpgradedClient implements ibcclienttypes.UpgradeKeeper
func (n *noopUpgradeKeeper) GetUpgradedClient(ctx context.Context, height int64) ([]byte, error) {
	return nil, nil
}

// SetUpgradedClient implements ibcclienttypes.UpgradeKeeper
func (n *noopUpgradeKeeper) SetUpgradedClient(ctx context.Context, planHeight int64, bz []byte) error {
	return nil
}

// GetUpgradedConsensusState implements ibcclienttypes.UpgradeKeeper
func (n *noopUpgradeKeeper) GetUpgradedConsensusState(ctx context.Context, lastHeight int64) ([]byte, error) {
	return nil, nil
}

// SetUpgradedConsensusState implements ibcclienttypes.UpgradeKeeper
func (n *noopUpgradeKeeper) SetUpgradedConsensusState(ctx context.Context, planHeight int64, bz []byte) error {
	return nil
}

// ScheduleUpgrade implements ibcclienttypes.UpgradeKeeper
func (n *noopUpgradeKeeper) ScheduleUpgrade(ctx context.Context, plan upgradetypes.Plan) error {
	return nil
}

// NewNoopUpgradeKeeper creates a minimal real UpgradeKeeper for testing
// ibckeeper.NewKeeper checks if UpgradeKeeper is empty, so we need a real instance
func NewNoopUpgradeKeeper() ibcclienttypes.UpgradeKeeper {
	// Create a minimal real UpgradeKeeper to pass the emptiness check
	// Use a memory store and minimal setup
	storeKey := storetypes.NewKVStoreKey(upgradetypes.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	
	registry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(registry)
	
	// Create a minimal upgrade keeper - this will pass the isEmpty check
	// The store doesn't need to be mounted since UpgradeKeeper won't be used in these tests
	upgradeKeeper := upgradekeeper.NewKeeper(
		map[int64]bool{}, // skipUpgradeHeights
		storeService,
		appCodec,
		"", // homePath - not needed for testing
		nil, // vsStore - version store, not needed for basic testing
		"", // authority - empty for testing
	)
	
	return upgradeKeeper
}

