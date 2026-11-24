package v36

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/poolmanager"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const (
	UpgradeName          = "v36"
	LegacyParamsStoreKey = "params"
)

// Keepers is the set of keepers and stores the v36 handler needs.
type Keepers struct {
	Account         authkeeper.AccountKeeper
	LegacyAmino     *codec.LegacyAmino
	LegacyParamsKey storetypes.StoreKey
	PoolManager     *poolmanager.Keeper
}

// CustomUpgradeHandlerLogic migrates poolmanager parameters from the legacy
// x/params store into the poolmanager module's own store and ensures the
// staking key rotation fee pool module account is initialized before running
// module migrations for Cosmos SDK v0.55.
func CustomUpgradeHandlerLogic(ctx context.Context, k Keepers) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := EnsureKeyRotationFeePoolAccount(sdkCtx, k.Account); err != nil {
		return fmt.Errorf("v36: key rotation fee pool account: %w", err)
	}
	if err := MigratePoolManagerParams(sdkCtx, k.LegacyAmino, k.LegacyParamsKey, k.PoolManager); err != nil {
		return fmt.Errorf("v36: poolmanager params migration: %w", err)
	}
	return nil
}

// EnsureKeyRotationFeePoolAccount ensures the key_rotation_fee_pool module
// account exists as a ModuleAccount with Burner permission, converting an
// existing BaseAccount at that address if one was created before v36.
func EnsureKeyRotationFeePoolAccount(ctx sdk.Context, ak authkeeper.AccountKeeper) error {
	addr, perms := ak.GetModuleAddressAndPermissions(stakingtypes.KeyRotationFeePoolName)
	if addr == nil {
		return fmt.Errorf("module account %s is not registered", stakingtypes.KeyRotationFeePoolName)
	}

	existing := ak.GetAccount(ctx, addr)
	if existing == nil {
		_ = ak.GetModuleAccount(ctx, stakingtypes.KeyRotationFeePoolName)
		return nil
	}
	if macc, ok := existing.(sdk.ModuleAccountI); ok && macc.HasPermission(authtypes.Burner) {
		return nil
	}

	base := authtypes.NewBaseAccount(existing.GetAddress(), nil, existing.GetAccountNumber(), existing.GetSequence())
	ak.SetModuleAccount(ctx, authtypes.NewModuleAccount(base, stakingtypes.KeyRotationFeePoolName, perms...))
	return nil
}

// MigratePoolManagerParams reads legacy poolmanager subspace entries from the
// raw "params" KV store (if present), overlays them onto DefaultParams,
// validates the result, and writes it to poolmanager's module store.
func MigratePoolManagerParams(
	ctx sdk.Context,
	legacyAmino *codec.LegacyAmino,
	legacyParamsKey storetypes.StoreKey,
	pmKeeper *poolmanager.Keeper,
) error {
	params := poolmanagertypes.DefaultParams()
	if legacyParamsKey != nil && legacyAmino != nil {
		store := ctx.KVStore(legacyParamsKey)
		prefix := append([]byte(poolmanagertypes.ModuleName), '/')
		readKey := func(paramKey []byte, dst interface{}) error {
			fullKey := append(append([]byte{}, prefix...), paramKey...)
			bz := store.Get(fullKey)
			if len(bz) == 0 {
				return nil
			}
			if err := legacyAmino.UnmarshalJSON(bz, dst); err != nil {
				return fmt.Errorf("unmarshal legacy param %q: %w", string(paramKey), err)
			}
			return nil
		}

		if err := readKey(poolmanagertypes.KeyDefaultTakerFee, &params.TakerFeeParams.DefaultTakerFee); err != nil {
			return err
		}
		if err := readKey(poolmanagertypes.KeyOsmoTakerFeeDistribution, &params.TakerFeeParams.OsmoTakerFeeDistribution); err != nil {
			return err
		}
		if err := readKey(poolmanagertypes.KeyNonOsmoTakerFeeDistribution, &params.TakerFeeParams.NonOsmoTakerFeeDistribution); err != nil {
			return err
		}
		if err := readKey(poolmanagertypes.KeyAdminAddresses, &params.TakerFeeParams.AdminAddresses); err != nil {
			return err
		}
		if err := readKey(poolmanagertypes.KeyCommunityPoolDenomToSwapNonWhitelistedAssetsTo, &params.TakerFeeParams.CommunityPoolDenomToSwapNonWhitelistedAssetsTo); err != nil {
			return err
		}
		if err := readKey(poolmanagertypes.KeyReducedTakerFeeByWhitelist, &params.TakerFeeParams.ReducedFeeWhitelist); err != nil {
			return err
		}
		if err := readKey(poolmanagertypes.KeyCommunityPoolDenomWhitelist, &params.TakerFeeParams.CommunityPoolDenomWhitelist); err != nil {
			return err
		}
	}

	if err := params.Validate(); err != nil {
		return fmt.Errorf("validate migrated poolmanager params: %w", err)
	}
	pmKeeper.SetParams(ctx, params)
	return nil
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k Keepers,
) func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		if err := CustomUpgradeHandlerLogic(ctx, k); err != nil {
			return nil, err
		}
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
