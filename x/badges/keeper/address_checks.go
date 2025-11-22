package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckAddressChecks validates address checks for a given address
// Returns (deterministicErrorMsg, error) where deterministicErrorMsg is a deterministic error string
func (k Keeper) CheckAddressChecks(
	ctx sdk.Context,
	addressChecks *types.AddressChecks,
	address string,
) (string, error) {
	if addressChecks == nil {
		return "", nil
	}

	// Check WASM contract requirements
	if addressChecks.MustBeWasmContract {
		isWasm, err := k.IsWasmContract(ctx, address)
		if err != nil {
			return "", err
		}
		if !isWasm {
			detErrMsg := fmt.Sprintf("address %s must be a WASM contract", address)
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	if addressChecks.MustNotBeWasmContract {
		isWasm, err := k.IsWasmContract(ctx, address)
		if err != nil {
			return "", err
		}
		if isWasm {
			detErrMsg := fmt.Sprintf("address %s must not be a WASM contract", address)
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	// Check liquidity pool requirements
	if addressChecks.MustBeLiquidityPool {
		isPool, err := k.IsLiquidityPool(ctx, address)
		if err != nil {
			return "", err
		}
		if !isPool {
			detErrMsg := fmt.Sprintf("address %s must be a liquidity pool", address)
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	if addressChecks.MustNotBeLiquidityPool {
		isPool, err := k.IsLiquidityPool(ctx, address)
		if err != nil {
			return "", err
		}
		if isPool {
			detErrMsg := fmt.Sprintf("address %s must not be a liquidity pool", address)
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	return "", nil
}

// IsWasmContract checks if an address is a WASM contract
func (k Keeper) IsWasmContract(ctx sdk.Context, address string) (bool, error) {
	if k.wasmViewKeeper == nil {
		// If WasmViewKeeper is not set, we can't check - return false
		// This allows the feature to work even if WASM module is not available
		return false, nil
	}

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return false, err
	}

	// Check if contract info exists for this address
	return k.wasmViewKeeper.HasContractInfo(ctx, addr), nil
}

// IsLiquidityPool checks if an address is a liquidity pool
func (k Keeper) IsLiquidityPool(ctx sdk.Context, address string) (bool, error) {
	if k.gammKeeper == nil {
		// If GammKeeper is not set, we can't check - return false
		// This allows the feature to work even if gamm module is not available
		return false, nil
	}

	// Validate address format
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return false, err
	}

	// Check cache first - pool addresses are static, so if cached, it's valid
	// The cache is populated when pools are created (see x/gamm/keeper/pool_service.go)
	_, found := k.GetPoolIdFromAddressCache(ctx, address)
	if found {
		return true, nil
	}

	// If not in cache, it's not a pool (cache is populated when pools are created)
	return false, nil
}
