package keeper

import (
	"fmt"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"

	"cosmossdk.io/core/store"
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		authority    []byte
		addressCodec interface{} // address.Codec - stored for msg server

		bankKeeper         types.BankKeeper
		distributionKeeper types.DistributionKeeper

		// prefixToRouter maps denom prefixes to their corresponding routers
		// Prefixes are registered when routers are registered
		prefixToRouter map[string]types.AliasDenomRouter

		// registeredPrefixes maintains the list of all registered prefixes in registration order
		registeredPrefixes []string
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
		prefixToRouter:     make(map[string]types.AliasDenomRouter),
		registeredPrefixes: []string{},
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// RegisterRouter registers an alias denom router for a specific prefix
// The prefix is stored globally and used for all routing operations
// This ensures that specific prefixes always route to the correct router
func (k *Keeper) RegisterRouter(prefix string, router types.AliasDenomRouter) error {
	// Validate prefix is not empty
	if prefix == "" {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "prefix cannot be empty")
	}

	// Check if prefix is already registered
	if _, exists := k.prefixToRouter[prefix]; exists {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "prefix %s is already registered", prefix)
	}

	// Register the prefix and router
	k.prefixToRouter[prefix] = router
	k.registeredPrefixes = append(k.registeredPrefixes, prefix)
	return nil
}

// GetRegisteredPrefixes returns all registered prefixes in registration order
func (k Keeper) GetRegisteredPrefixes() []string {
	// Return a copy to prevent external modification
	result := make([]string, len(k.registeredPrefixes))
	copy(result, k.registeredPrefixes)
	return result
}

// getRouterForDenom returns the router for a given denom based on prefix matching
// Returns the router and true if a matching prefix is found, nil and false otherwise
// Uses the globally stored registered prefixes
func (k Keeper) getRouterForDenom(denom string) (types.AliasDenomRouter, bool) {
	// Check all registered prefixes in registration order
	for _, prefix := range k.registeredPrefixes {
		if strings.HasPrefix(denom, prefix) {
			router, exists := k.prefixToRouter[prefix]
			if !exists {
				// This should never happen if registration is correct, but handle it gracefully
				return nil, false
			}
			return router, true
		}
	}
	return nil, false
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SendCoinWithAliasRouting sends a coin using the appropriate routing (alias denom or bank)
// For alias denoms, it finds the appropriate router based on prefix matching and uses SendNativeTokensViaAliasDenom
// For regular denoms, it uses bank keeper SendCoins
// Fails if a prefix matches but no router is registered for it
func (k Keeper) SendCoinWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	toAddressAcc sdk.AccAddress,
	coin *sdk.Coin,
) error {
	// Check if this denom matches any known prefix
	router, found := k.getRouterForDenom(coin.Denom)
	if found {
		// Prefix matched and router is registered - use it
		amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
		return router.SendNativeTokensViaAliasDenom(ctx, fromAddressAcc.String(), toAddressAcc.String(), coin.Denom, amountUint)
	}

	// Check if prefix matched but no router was registered (error condition)
	// This should not happen if getRouterForDenom is working correctly, but check for safety
	for _, prefix := range k.registeredPrefixes {
		if strings.HasPrefix(coin.Denom, prefix) {
			// Prefix matches but no router registered - this is an error
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "denom %s matches prefix %s but no router is registered for this prefix", coin.Denom, prefix)
		}
	}

	// No prefix matched, use standard bank routing
	return k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(*coin))
}

// SendCoinsWithAliasRouting sends coins using the appropriate routing (alias denom or bank)
// For alias denoms, it finds the appropriate router based on prefix matching and uses SendNativeTokensViaAliasDenom
// For regular denoms, it uses bank keeper SendCoins
// Fails if a prefix matches but no router is registered for it
func (k Keeper) SendCoinsWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	toAddressAcc sdk.AccAddress,
	coins sdk.Coins,
) error {
	// Process each coin individually
	for _, coin := range coins {
		// Check if this denom matches any known prefix
		router, found := k.getRouterForDenom(coin.Denom)
		if found {
			// Prefix matched and router is registered - use it
			amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			err := router.SendNativeTokensViaAliasDenom(ctx, fromAddressAcc.String(), toAddressAcc.String(), coin.Denom, amountUint)
			if err != nil {
				return err
			}
			continue
		}

		// Check if prefix matched but no router was registered (error condition)
		// This should not happen if getRouterForDenom is working correctly, but check for safety
		matchedPrefix := ""
		for _, prefix := range k.registeredPrefixes {
			if strings.HasPrefix(coin.Denom, prefix) {
				matchedPrefix = prefix
				break
			}
		}
		if matchedPrefix != "" {
			// Prefix matches but no router registered - this is an error
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "denom %s matches prefix %s but no router is registered for this prefix", coin.Denom, matchedPrefix)
		}

		// No prefix matched, use standard bank routing
		err := k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}

	return nil
}

// FundCommunityPoolWithAliasRouting funds the community pool, using alias denom routing if needed
// For alias denoms, it finds the appropriate router and uses FundCommunityPoolViaAliasDenom
// For regular denoms, it uses distribution keeper FundCommunityPool
func (k Keeper) FundCommunityPoolWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	coins sdk.Coins,
) error {
	// Get community pool module address (distribution module)
	moduleAddress := authtypes.NewModuleAddress(distrtypes.ModuleName).String()

	for _, coin := range coins {
		// Check if this denom matches any known prefix
		router, found := k.getRouterForDenom(coin.Denom)
		if found {
			// Prefix matched and router is registered - use it
			amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			err := router.FundCommunityPoolViaAliasDenom(ctx, fromAddressAcc.String(), moduleAddress, coin.Denom, amountUint)
			if err != nil {
				return err
			}
			continue
		}

		// Check if prefix matched but no router was registered (error condition)
		matchedPrefix := ""
		for _, prefix := range k.registeredPrefixes {
			if strings.HasPrefix(coin.Denom, prefix) {
				matchedPrefix = prefix
				break
			}
		}
		if matchedPrefix != "" {
			// Prefix matches but no router registered - this is an error
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "denom %s matches prefix %s but no router is registered for this prefix", coin.Denom, matchedPrefix)
		}

		// No prefix matched, use standard distribution routing
		err := k.distributionKeeper.FundCommunityPool(ctx, sdk.NewCoins(coin), fromAddressAcc)
		if err != nil {
			return err
		}
	}

	return nil
}

// SpendFromCommunityPoolWithAliasRouting spends from the community pool, using alias denom routing if needed
// For alias denoms, it finds the appropriate router and uses SpendFromCommunityPoolViaAliasDenom (standard send)
// For regular denoms, it uses bank keeper SendCoinsFromModuleToAccount with the distribution module
func (k Keeper) SpendFromCommunityPoolWithAliasRouting(
	ctx sdk.Context,
	toAddressAcc sdk.AccAddress,
	coins sdk.Coins,
) error {
	// Get community pool module address (distribution module)
	moduleName := distrtypes.ModuleName
	moduleAddress := authtypes.NewModuleAddress(moduleName)

	for _, coin := range coins {
		// Check if this denom matches any known prefix
		router, found := k.getRouterForDenom(coin.Denom)
		if found {
			// Prefix matched and router is registered - use standard send for badges keeper
			amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			err := router.SpendFromCommunityPoolViaAliasDenom(ctx, moduleAddress.String(), toAddressAcc.String(), coin.Denom, amountUint)
			if err != nil {
				return err
			}
			continue
		}

		// Check if prefix matched but no router was registered (error condition)
		matchedPrefix := ""
		for _, prefix := range k.registeredPrefixes {
			if strings.HasPrefix(coin.Denom, prefix) {
				matchedPrefix = prefix
				break
			}
		}
		if matchedPrefix != "" {
			// Prefix matches but no router registered - this is an error
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "denom %s matches prefix %s but no router is registered for this prefix", coin.Denom, matchedPrefix)
		}

		// No prefix matched, use standard bank routing
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleName, toAddressAcc, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}

	return nil
}

// GetBalanceWithAliasRouting gets the balance for a specific denom, handling alias denom routing
// Mirrors bankKeeper.GetBalance but routes alias denoms through their respective routers
// For alias denoms (e.g., badgeslp:), uses the router's GetBalanceWithAliasRouting
// For regular denoms, uses bankKeeper.GetBalance
func (k Keeper) GetBalanceWithAliasRouting(ctx sdk.Context, address sdk.AccAddress, denom string) (sdk.Coin, error) {
	// Check if this denom matches any known prefix
	router, found := k.getRouterForDenom(denom)
	if found {
		// Prefix matched and router is registered - use it
		return router.GetBalanceWithAliasRouting(ctx, address, denom)
	}

	// No prefix matched, use standard bank routing
	return k.bankKeeper.GetBalance(ctx, address, denom), nil
}

// IsICS20Compatible checks if a denom is ICS20 compatible (i.e., doesn't match any alias denom prefix)
// Returns true if and only if the denom doesn't match any registered prefix (is a bank keeper denom)
// Returns false if the denom matches any registered prefix (is an alias denom)
func (k Keeper) IsICS20Compatible(ctx sdk.Context, denom string) bool {
	// Check if this denom matches any known prefix
	_, found := k.getRouterForDenom(denom)
	// Return true if no prefix matched (ICS20 compatible), false if prefix matched (not ICS20 compatible)
	return !found
}

// StandardName returns the standard name for a denom type
// For alias denoms (e.g., badgeslp:, badges:), returns "BitBadges"
// For regular ICS20 denoms, returns "ICS20"
func (k Keeper) StandardName(ctx sdk.Context, denom string) string {
	// Check if this denom matches any known prefix
	_, found := k.getRouterForDenom(denom)
	if found {
		// Alias denom - return "BitBadges"
		return "x/badges"
	}
	// Regular ICS20 denom
	return "x/bank"
}

// SendCoinsFromModuleToAccountWithAliasRouting sends coins from a module account to an account, using alias denom routing if needed
// For alias denoms, it finds the appropriate router and uses SendNativeTokensViaAliasDenom (standard send)
// For regular denoms, it uses bank keeper SendCoinsFromModuleToAccount
func (k Keeper) SendCoinsFromModuleToAccountWithAliasRouting(
	ctx sdk.Context,
	moduleName string,
	toAddressAcc sdk.AccAddress,
	coins sdk.Coins,
) error {
	moduleAddress := authtypes.NewModuleAddress(moduleName)

	for _, coin := range coins {
		// Check if this denom matches any known prefix
		router, found := k.getRouterForDenom(coin.Denom)
		if found {
			// Prefix matched and router is registered - use adapter method
			amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			err := router.SendFromModuleToAccountViaAliasDenom(ctx, moduleAddress.String(), toAddressAcc.String(), coin.Denom, amountUint)
			if err != nil {
				return err
			}
			continue
		}

		// Check if prefix matched but no router was registered (error condition)
		matchedPrefix := ""
		for _, prefix := range k.registeredPrefixes {
			if strings.HasPrefix(coin.Denom, prefix) {
				matchedPrefix = prefix
				break
			}
		}
		if matchedPrefix != "" {
			// Prefix matches but no router registered - this is an error
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "denom %s matches prefix %s but no router is registered for this prefix", coin.Denom, matchedPrefix)
		}

		// No prefix matched, use standard bank routing
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleName, toAddressAcc, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}

	return nil
}

// SendCoinsFromAccountToModuleWithAliasRouting sends coins from an account to a module account, using alias denom routing if needed
// For alias denoms, it finds the appropriate router and uses SendNativeTokensViaAliasDenom (standard send)
// For regular denoms, it uses bank keeper SendCoinsFromAccountToModule
func (k Keeper) SendCoinsFromAccountToModuleWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	moduleName string,
	coins sdk.Coins,
) error {
	moduleAddress := authtypes.NewModuleAddress(moduleName)

	for _, coin := range coins {
		// Check if this denom matches any known prefix
		router, found := k.getRouterForDenom(coin.Denom)
		if found {
			// Prefix matched and router is registered - use adapter method
			amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			err := router.SendFromAccountToModuleViaAliasDenom(ctx, fromAddressAcc.String(), moduleAddress.String(), coin.Denom, amountUint)
			if err != nil {
				return err
			}
			continue
		}

		// Check if prefix matched but no router was registered (error condition)
		matchedPrefix := ""
		for _, prefix := range k.registeredPrefixes {
			if strings.HasPrefix(coin.Denom, prefix) {
				matchedPrefix = prefix
				break
			}
		}
		if matchedPrefix != "" {
			// Prefix matches but no router registered - this is an error
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "denom %s matches prefix %s but no router is registered for this prefix", coin.Denom, matchedPrefix)
		}

		// No prefix matched, use standard bank routing
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAddressAcc, moduleName, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}

	return nil
}
