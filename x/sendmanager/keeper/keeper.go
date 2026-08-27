package keeper

import (
	"fmt"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"

	"cosmossdk.io/core/store"
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/log/v2"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// AliasDenomPrefix is the hardcoded prefix for tokenization alias denoms.
// Any denom starting with this prefix is routed through the alias router
// instead of x/bank. This is a compile-time constant — no dynamic registration needed.
const AliasDenomPrefix = "badgeslp:"

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		authority    []byte
		addressCodec interface{} // address.Codec - stored for msg server

		bankKeeper         types.BankKeeper
		distributionKeeper types.DistributionKeeper

		// aliasRouter is the router for badgeslp: prefixed denoms.
		// Set once via SetAliasRouter after both keepers are created.
		// Uses a pointer so value-copies of the Keeper (e.g., via depinject interface
		// satisfaction) share the same router reference.
		aliasRouter *types.AliasDenomRouter
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	addressCodec interface{}, // address.Codec
	authority []byte,
	bankKeeper types.BankKeeper,
	distributionKeeper types.DistributionKeeper,
) Keeper {
	return Keeper{
		cdc:                cdc,
		storeService:       storeService,
		logger:             nil, // Will be set via SetLogger if needed
		authority:          authority,
		addressCodec:       addressCodec,
		bankKeeper:         bankKeeper,
		distributionKeeper: distributionKeeper,
		aliasRouter:        new(types.AliasDenomRouter), // allocate pointer, nil value inside
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// SetAliasRouter sets the alias denom router. Must be called after both keepers are created.
// Uses the shared pointer so all value-copies of the Keeper see the same router.
func (k *Keeper) SetAliasRouter(router types.AliasDenomRouter) {
	*k.aliasRouter = router
}

// RegisterRouter is kept for backward compatibility with existing app.go code.
// It ignores the prefix parameter (hardcoded to badgeslp:) and delegates to SetAliasRouter.
func (k *Keeper) RegisterRouter(prefix string, router types.AliasDenomRouter) error {
	if prefix != AliasDenomPrefix {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "only prefix %q is supported, got %q", AliasDenomPrefix, prefix)
	}
	k.SetAliasRouter(router)
	return nil
}

// GetRegisteredPrefixes returns the supported prefixes (for backward compatibility).
func (k Keeper) GetRegisteredPrefixes() []string {
	if k.aliasRouter != nil && *k.aliasRouter != nil {
		return []string{AliasDenomPrefix}
	}
	return []string{}
}

// getRouterForDenom checks if the denom starts with the alias prefix and returns the router.
// Simple string prefix check — no dynamic registry, no stale-copy issues.
func (k Keeper) getRouterForDenom(denom string) (types.AliasDenomRouter, bool) {
	if denom == "" {
		return nil, false
	}
	if strings.HasPrefix(denom, AliasDenomPrefix) && k.aliasRouter != nil && *k.aliasRouter != nil {
		return *k.aliasRouter, true
	}
	return nil, false
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SendCoinWithAliasRouting sends a coin using the appropriate routing (alias denom or bank)
func (k Keeper) SendCoinWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	toAddressAcc sdk.AccAddress,
	coin *sdk.Coin,
) error {
	if coin.Denom == "" {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
	}

	router, found := k.getRouterForDenom(coin.Denom)
	if found {
		if coin.Amount.IsNegative() {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin amount cannot be negative: %s", coin.Denom)
		}
		amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
		return router.SendNativeTokensViaAliasDenom(ctx, fromAddressAcc.String(), toAddressAcc.String(), coin.Denom, amountUint)
	}

	return k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(*coin))
}

// SendCoinsWithAliasRouting sends coins using the appropriate routing (alias denom or bank)
func (k Keeper) SendCoinsWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	toAddressAcc sdk.AccAddress,
	coins sdk.Coins,
) error {
	for _, coin := range coins {
		if coin.Denom == "" {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
		}

		router, found := k.getRouterForDenom(coin.Denom)
		if found {
			if coin.Amount.IsNegative() {
				return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin amount cannot be negative: %s", coin.Denom)
			}
			amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			err := router.SendNativeTokensViaAliasDenom(ctx, fromAddressAcc.String(), toAddressAcc.String(), coin.Denom, amountUint)
			if err != nil {
				return err
			}
			continue
		}

		err := k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}

	return nil
}

// FundCommunityPoolWithAliasRouting funds the community pool, using alias denom routing if needed
func (k Keeper) FundCommunityPoolWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	coins sdk.Coins,
) error {
	moduleAddress := authtypes.NewModuleAddress(distrtypes.ModuleName).String()

	for _, coin := range coins {
		if coin.Denom == "" {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
		}

		router, found := k.getRouterForDenom(coin.Denom)
		if found {
			if coin.Amount.IsNegative() {
				return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin amount cannot be negative: %s", coin.Denom)
			}
			amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			err := router.FundCommunityPoolViaAliasDenom(ctx, fromAddressAcc.String(), moduleAddress, coin.Denom, amountUint)
			if err != nil {
				return err
			}
			continue
		}

		err := k.distributionKeeper.FundCommunityPool(ctx, sdk.NewCoins(coin), fromAddressAcc)
		if err != nil {
			return err
		}
	}

	return nil
}

// SpendFromCommunityPoolWithAliasRouting spends from the community pool, using alias denom routing if needed
func (k Keeper) SpendFromCommunityPoolWithAliasRouting(
	ctx sdk.Context,
	toAddressAcc sdk.AccAddress,
	coins sdk.Coins,
) error {
	moduleName := distrtypes.ModuleName
	moduleAddress := authtypes.NewModuleAddress(moduleName)

	for _, coin := range coins {
		if coin.Denom == "" {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
		}

		router, found := k.getRouterForDenom(coin.Denom)
		if found {
			if coin.Amount.IsNegative() {
				return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin amount cannot be negative: %s", coin.Denom)
			}
			amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			err := router.SpendFromCommunityPoolViaAliasDenom(ctx, moduleAddress.String(), toAddressAcc.String(), coin.Denom, amountUint)
			if err != nil {
				return err
			}
			continue
		}

		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleName, toAddressAcc, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}

	return nil
}

// GetBalanceWithAliasRouting gets the balance for a specific denom, handling alias denom routing
func (k Keeper) GetBalanceWithAliasRouting(ctx sdk.Context, address sdk.AccAddress, denom string) (sdk.Coin, error) {
	if denom == "" {
		return sdk.Coin{}, sdkerrors.Wrapf(types.ErrInvalidRequest, "denom cannot be empty")
	}

	router, found := k.getRouterForDenom(denom)
	if found {
		return router.GetBalanceWithAliasRouting(ctx, address, denom)
	}

	return k.bankKeeper.GetBalance(ctx, address, denom), nil
}

// IsICS20Compatible checks if a denom is ICS20 compatible (not an alias denom)
func (k Keeper) IsICS20Compatible(ctx sdk.Context, denom string) bool {
	if denom == "" {
		return true
	}
	_, found := k.getRouterForDenom(denom)
	return !found
}

// StandardName returns the standard name for a denom type
func (k Keeper) StandardName(ctx sdk.Context, denom string) string {
	if denom == "" {
		return "x/bank"
	}
	_, found := k.getRouterForDenom(denom)
	if found {
		return "x/tokenization"
	}
	return "x/bank"
}

// SendCoinsFromModuleToAccountWithAliasRouting sends coins from a module account to an account
func (k Keeper) SendCoinsFromModuleToAccountWithAliasRouting(
	ctx sdk.Context,
	moduleName string,
	toAddressAcc sdk.AccAddress,
	coins sdk.Coins,
) error {
	moduleAddress := authtypes.NewModuleAddress(moduleName)

	for _, coin := range coins {
		if coin.Denom == "" {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
		}

		router, found := k.getRouterForDenom(coin.Denom)
		if found {
			if coin.Amount.IsNegative() {
				return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin amount cannot be negative: %s", coin.Denom)
			}
			amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			err := router.SendFromModuleToAccountViaAliasDenom(ctx, moduleAddress.String(), toAddressAcc.String(), coin.Denom, amountUint)
			if err != nil {
				return err
			}
			continue
		}

		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleName, toAddressAcc, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}

	return nil
}

// SendCoinsFromAccountToModuleWithAliasRouting sends coins from an account to a module account
func (k Keeper) SendCoinsFromAccountToModuleWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	moduleName string,
	coins sdk.Coins,
) error {
	moduleAddress := authtypes.NewModuleAddress(moduleName)

	for _, coin := range coins {
		if coin.Denom == "" {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
		}

		router, found := k.getRouterForDenom(coin.Denom)
		if found {
			if coin.Amount.IsNegative() {
				return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin amount cannot be negative: %s", coin.Denom)
			}
			amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			err := router.SendFromAccountToModuleViaAliasDenom(ctx, fromAddressAcc.String(), moduleAddress.String(), coin.Denom, amountUint)
			if err != nil {
				return err
			}
			continue
		}

		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAddressAcc, moduleName, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}

	return nil
}
