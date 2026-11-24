package poolmanager

import (
	"github.com/bitbadges/bitbadgeschain/third_party/osmoutils"
	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey

	gammKeeper          *gammkeeper.Keeper
	bankKeeper          types.BankI
	accountKeeper       types.AccountI
	communityPoolKeeper types.CommunityPoolI
	stakingKeeper       types.StakingKeeper

	// routes is a map to get the pool module by id.
	routes map[types.PoolType]types.PoolModuleI

	// poolModules is a list of all pool modules.
	// It is used when an operation has to be applied to all pool
	// modules. Since map iterations are non-deterministic, we
	// use this list to ensure deterministic iteration.
	poolModules []types.PoolModuleI
}

func NewKeeper(storeKey storetypes.StoreKey, gammKeeper *gammkeeper.Keeper, bankKeeper types.BankI, accountKeeper types.AccountI, communityPoolKeeper types.CommunityPoolI, stakingKeeper types.StakingKeeper) *Keeper {
	routesMap := map[types.PoolType]types.PoolModuleI{
		types.Balancer:   gammKeeper,
		types.Stableswap: gammKeeper,
	}

	routesList := []types.PoolModuleI{
		gammKeeper,
	}

	return &Keeper{
		storeKey:            storeKey,
		gammKeeper:          gammKeeper,
		bankKeeper:          bankKeeper,
		accountKeeper:       accountKeeper,
		communityPoolKeeper: communityPoolKeeper,
		routes:              routesMap,
		poolModules:         routesList,
		stakingKeeper:       stakingKeeper,
	}
}

// ResetCaches is retained for callers that previously reset the route cache.
// Routes are now read directly from the caller's context.
func (k *Keeper) ResetCaches() {}

// GetParams returns the total set of poolmanager parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	found, err := osmoutils.Get(store, types.KeyParams, &params)
	if err != nil {
		panic(err)
	}
	if !found {
		return types.DefaultParams()
	}
	if params.TakerFeeParams.AdminAddresses == nil {
		params.TakerFeeParams.AdminAddresses = []string{}
	}
	if params.TakerFeeParams.ReducedFeeWhitelist == nil {
		params.TakerFeeParams.ReducedFeeWhitelist = []string{}
	}
	if params.TakerFeeParams.CommunityPoolDenomWhitelist == nil {
		params.TakerFeeParams.CommunityPoolDenomWhitelist = []string{}
	}
	return params
}

// SetParams sets the total set of poolmanager parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	osmoutils.MustSet(store, types.KeyParams, &params)
}

// InitGenesis initializes the poolmanager module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	k.SetNextPoolId(ctx, genState.NextPoolId)
	if err := genState.Validate(); err != nil {
		panic(err)
	}

	k.SetParams(ctx, genState.Params)

	for _, poolRoute := range genState.PoolRoutes {
		k.SetPoolRoute(ctx, poolRoute.PoolId, poolRoute.PoolType)
	}

	// We track taker fees generated in the module's KVStore.
	for _, coin := range genState.TakerFeesTracker.TakerFeesToStakers {
		if err := k.UpdateTakerFeeTrackerForStakersByDenom(ctx, coin.Denom, coin.Amount); err != nil {
			panic(err)
		}
	}
	for _, coin := range genState.TakerFeesTracker.TakerFeesToCommunityPool {
		if err := k.UpdateTakerFeeTrackerForCommunityPoolByDenom(ctx, coin.Denom, coin.Amount); err != nil {
			panic(err)
		}
	}
	k.SetTakerFeeTrackerStartHeight(ctx, genState.TakerFeesTracker.HeightAccountingStartsFrom)

	// Set the pool volumes KVStore.
	for _, poolVolume := range genState.PoolVolumes {
		k.SetVolume(ctx, poolVolume.PoolId, poolVolume.PoolVolume)
	}

	// Set the denom pair taker fees KVStore.
	for _, denomPairTakerFee := range genState.DenomPairTakerFeeStore {
		k.SetDenomPairTakerFee(ctx, denomPairTakerFee.TokenInDenom, denomPairTakerFee.TokenOutDenom, denomPairTakerFee.TakerFee)
	}
}

// ExportGenesis returns the poolmanager module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	pools, err := k.AllPools(ctx)
	if err != nil {
		panic(err)
	}

	// Utilize poolVolumes struct to export pool volumes from KVStore.
	poolVolumes := make([]*types.PoolVolume, len(pools))
	for i, pool := range pools {
		poolVolume := k.GetTotalVolumeForPool(ctx, pool.GetId())
		poolVolumes[i] = &types.PoolVolume{
			PoolId:     pool.GetId(),
			PoolVolume: poolVolume,
		}
	}

	// Utilize denomPairTakerFee struct to export taker fees from KVStore.
	denomPairTakerFees, err := k.GetAllTradingPairTakerFees(ctx)
	if err != nil {
		panic(err)
	}

	// Export KVStore values to the genesis state so they can be imported in init genesis.
	takerFeesTracker := types.TakerFeesTracker{
		TakerFeesToStakers:         k.GetTakerFeeTrackerForStakers(ctx),
		TakerFeesToCommunityPool:   k.GetTakerFeeTrackerForCommunityPool(ctx),
		HeightAccountingStartsFrom: k.GetTakerFeeTrackerStartHeight(ctx),
	}
	return &types.GenesisState{
		Params:                 k.GetParams(ctx),
		NextPoolId:             k.GetNextPoolId(ctx),
		PoolRoutes:             k.getAllPoolRoutes(ctx),
		TakerFeesTracker:       &takerFeesTracker,
		PoolVolumes:            poolVolumes,
		DenomPairTakerFeeStore: denomPairTakerFees,
	}
}

// GetNextPoolId returns the next pool id.
func (k Keeper) GetNextPoolId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	nextPoolId := gogotypes.UInt64Value{}
	osmoutils.MustGet(store, types.KeyNextGlobalPoolId, &nextPoolId)
	return nextPoolId.Value
}

// SetNextPoolId sets next pool Id.
func (k Keeper) SetNextPoolId(ctx sdk.Context, poolId uint64) {
	store := ctx.KVStore(k.storeKey)
	osmoutils.MustSet(store, types.KeyNextGlobalPoolId, &gogotypes.UInt64Value{Value: poolId})
}

// SetStakingKeeper sets staking keeper
func (k *Keeper) SetStakingKeeper(stakingKeeper types.StakingKeeper) {
	k.stakingKeeper = stakingKeeper
}

func (k *Keeper) BeginBlock(ctx sdk.Context) {}

// AlloyedAssetCompositionUpdateRate is the rate in blocks at which the taker fee share alloy composition is updated in the end block.
var AlloyedAssetCompositionUpdateRate = int64(700)

// EndBlock updates the taker fee share alloy composition for all registered alloyed pools
// if the current block height is a multiple of the alloyedAssetCompositionUpdateRate.
func (k *Keeper) EndBlock(ctx sdk.Context) {
}
