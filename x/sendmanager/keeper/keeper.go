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
// Prevents overlapping prefixes to avoid routing ambiguity
func (k *Keeper) RegisterRouter(prefix string, router types.AliasDenomRouter) error {
	// Validate prefix is not empty
	if prefix == "" {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "prefix cannot be empty")
	}

	// Check if prefix is already registered
	if _, exists := k.prefixToRouter[prefix]; exists {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "prefix %s is already registered", prefix)
	}

	// Check for overlapping prefixes (prevent sub-prefix or super-prefix conflicts)
	// A prefix overlaps if one is a prefix of the other
	for existingPrefix := range k.prefixToRouter {
		if strings.HasPrefix(prefix, existingPrefix) {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "prefix %s overlaps with existing prefix %s (new prefix starts with existing prefix)", prefix, existingPrefix)
		}
		if strings.HasPrefix(existingPrefix, prefix) {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "prefix %s overlaps with existing prefix %s (existing prefix starts with new prefix)", prefix, existingPrefix)
		}
	}

	// Register the prefix and router
	k.prefixToRouter[prefix] = router
	k.registeredPrefixes = append(k.registeredPrefixes, prefix)
	
	// Sort registeredPrefixes by length (longest first) for longest-prefix matching
	// This ensures that longer prefixes are checked before shorter ones
	k.sortPrefixesByLength()
	
	return nil
}

// sortPrefixesByLength sorts registeredPrefixes by length in descending order (longest first)
// This enables longest-prefix matching in getRouterForDenom
func (k *Keeper) sortPrefixesByLength() {
	// Create a slice of prefix-length pairs
	type prefixLen struct {
		prefix string
		length int
	}
	pairs := make([]prefixLen, len(k.registeredPrefixes))
	for i, prefix := range k.registeredPrefixes {
		pairs[i] = prefixLen{prefix: prefix, length: len(prefix)}
	}
	
	// Sort by length descending (longest first), then by string order for stability
	for i := 0; i < len(pairs)-1; i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].length < pairs[j].length || (pairs[i].length == pairs[j].length && pairs[i].prefix > pairs[j].prefix) {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
	
	// Rebuild registeredPrefixes in sorted order
	k.registeredPrefixes = make([]string, len(pairs))
	for i, pair := range pairs {
		k.registeredPrefixes[i] = pair.prefix
	}
}

// GetRegisteredPrefixes returns all registered prefixes in registration order
func (k Keeper) GetRegisteredPrefixes() []string {
	// Return a copy to prevent external modification
	result := make([]string, len(k.registeredPrefixes))
	copy(result, k.registeredPrefixes)
	return result
}

// getRouterForDenom returns the router for a given denom based on longest-prefix matching
// Returns the router and true if a matching prefix is found, nil and false otherwise
// Uses longest-prefix matching: if multiple prefixes match, the longest one is used
// registeredPrefixes is sorted by length (longest first) to enable this
func (k Keeper) getRouterForDenom(denom string) (types.AliasDenomRouter, bool) {
	// Validate denom is not empty
	if denom == "" {
		return nil, false
	}
	
	// Check all registered prefixes (sorted by length, longest first)
	// This implements longest-prefix matching
	for _, prefix := range k.registeredPrefixes {
		if strings.HasPrefix(denom, prefix) {
			router, exists := k.prefixToRouter[prefix]
			// Security: MED-003 - Router inconsistency detection
			// If prefix exists in registeredPrefixes but not in prefixToRouter, this indicates
			// data corruption or an inconsistent state. This should never happen in normal operation.
			// Panic to detect corruption early - this is a critical invariant violation.
			if !exists {
				panic(fmt.Sprintf("sendmanager: router inconsistency detected - prefix '%s' exists in registeredPrefixes but not in prefixToRouter. This indicates data corruption.", prefix))
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
func (k Keeper) SendCoinWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	toAddressAcc sdk.AccAddress,
	coin *sdk.Coin,
) error {
	// Validate coin denom is not empty
	if coin.Denom == "" {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
	}

	// Check if this denom matches any known prefix using longest-prefix matching
	router, found := k.getRouterForDenom(coin.Denom)
	if found {
		// Prefix matched and router is registered - use it
		amountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
		return router.SendNativeTokensViaAliasDenom(ctx, fromAddressAcc.String(), toAddressAcc.String(), coin.Denom, amountUint)
	}

	// No prefix matched, use standard bank routing
	return k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(*coin))
}

// SendCoinsWithAliasRouting sends coins using the appropriate routing (alias denom or bank)
// For alias denoms, it finds the appropriate router based on prefix matching and uses SendNativeTokensViaAliasDenom
// For regular denoms, it uses bank keeper SendCoins
func (k Keeper) SendCoinsWithAliasRouting(
	ctx sdk.Context,
	fromAddressAcc sdk.AccAddress,
	toAddressAcc sdk.AccAddress,
	coins sdk.Coins,
) error {
	// Process each coin individually
	for _, coin := range coins {
		// Validate coin denom is not empty
		if coin.Denom == "" {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
		}

		// Check if this denom matches any known prefix using longest-prefix matching
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
		// Validate coin denom is not empty
		if coin.Denom == "" {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
		}

		// Check if this denom matches any known prefix using longest-prefix matching
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
		// Validate coin denom is not empty
		if coin.Denom == "" {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
		}

		// Check if this denom matches any known prefix using longest-prefix matching
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
// For alias denoms (e.g., tokenizationlp:), uses the router's GetBalanceWithAliasRouting
// For regular denoms, uses bankKeeper.GetBalance
func (k Keeper) GetBalanceWithAliasRouting(ctx sdk.Context, address sdk.AccAddress, denom string) (sdk.Coin, error) {
	// Validate denom is not empty
	if denom == "" {
		return sdk.Coin{}, sdkerrors.Wrapf(types.ErrInvalidRequest, "denom cannot be empty")
	}

	// Check if this denom matches any known prefix using longest-prefix matching
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
	// Validate denom is not empty
	if denom == "" {
		return true // Empty denom is considered ICS20 compatible (will be handled by bank keeper)
	}

	// Check if this denom matches any known prefix using longest-prefix matching
	_, found := k.getRouterForDenom(denom)
	// Return true if no prefix matched (ICS20 compatible), false if prefix matched (not ICS20 compatible)
	return !found
}

// StandardName returns the standard name for a denom type
// For alias denoms (e.g., tokenizationlp:, tokenization:), returns "BitBadges"
// For regular ICS20 denoms, returns "ICS20"
func (k Keeper) StandardName(ctx sdk.Context, denom string) string {
	// Validate denom is not empty
	if denom == "" {
		return "x/bank" // Empty denom defaults to bank
	}

	// Check if this denom matches any known prefix using longest-prefix matching
	_, found := k.getRouterForDenom(denom)
	if found {
		// Alias denom - return "BitBadges"
		return "x/tokenization"
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
		// Validate coin denom is not empty
		if coin.Denom == "" {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
		}

		// Check if this denom matches any known prefix using longest-prefix matching
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
		// Validate coin denom is not empty
		if coin.Denom == "" {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin denom cannot be empty")
		}

		// Check if this denom matches any known prefix using longest-prefix matching
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

		// No prefix matched, use standard bank routing
		err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAddressAcc, moduleName, sdk.NewCoins(coin))
		if err != nil {
			return err
		}
	}

	return nil
}
