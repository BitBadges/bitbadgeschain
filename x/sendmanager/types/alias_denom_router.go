package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AliasDenomRouter defines the interface that modules must implement
// to participate in alias denom routing (e.g., badges: or badgeslp: prefixes)
type AliasDenomRouter interface {
	// CheckIsAliasDenom checks if a given denom is an alias denom handled by this router
	// Returns true if this router can handle the denom, false otherwise
	CheckIsAliasDenom(ctx sdk.Context, denom string) bool

	// SendNativeTokensViaAliasDenom sends native tokens using the alias denom routing
	// This is called when a coin with an alias denom needs to be sent using the module's handled alias approach
	SendNativeTokensViaAliasDenom(ctx sdk.Context, recipientAddress string, toAddress string, denom string, amount sdkmath.Uint) error

	// FundCommunityPoolViaAliasDenom funds the community pool using alias denom routing
	// This handles the alias denom-specific logic for funding the community pool (e.g., setting auto-approvals)
	FundCommunityPoolViaAliasDenom(ctx sdk.Context, fromAddress string, toAddress string, denom string, amount sdkmath.Uint) error

	// SpendFromCommunityPoolViaAliasDenom spends from the community pool using alias denom routing
	// This handles the alias denom-specific logic for spending from the community pool (e.g., standard send)
	SpendFromCommunityPoolViaAliasDenom(ctx sdk.Context, fromAddress string, toAddress string, denom string, amount sdkmath.Uint) error

	// SendFromModuleToAccountViaAliasDenom sends from a module account to an account using alias denom routing
	// This handles the alias denom-specific logic for module-to-account transfers (e.g., standard send)
	SendFromModuleToAccountViaAliasDenom(ctx sdk.Context, moduleAddress string, toAddress string, denom string, amount sdkmath.Uint) error

	// SendFromAccountToModuleViaAliasDenom sends from an account to a module account using alias denom routing
	// This handles the alias denom-specific logic for account-to-module transfers (e.g., standard send)
	SendFromAccountToModuleViaAliasDenom(ctx sdk.Context, fromAddress string, moduleAddress string, denom string, amount sdkmath.Uint) error

	// GetBalanceWithAliasRouting gets the balance for a specific denom, handling alias denom routing
	// For alias denoms (e.g., badgeslp:), this may use custom logic (e.g., getMaxWrappableAmount flow)
	// Returns the coin balance for the given address and denom
	GetBalanceWithAliasRouting(ctx sdk.Context, address sdk.AccAddress, denom string) (sdk.Coin, error)
}
