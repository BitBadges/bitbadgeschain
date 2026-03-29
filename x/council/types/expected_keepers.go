package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TokenizationKeeper defines the interface for the x/tokenization keeper that x/council needs.
// We define a minimal interface so x/council does not import x/tokenization's concrete types.
//
// The wiring code in app.go should create an adapter struct that implements this interface
// using the concrete x/tokenization keeper methods.
type TokenizationKeeper interface {
	// GetCredentialBalance returns the balance of a specific credential token for a given address.
	// collectionId and tokenId identify the credential token to check.
	// Returns the balance amount (0 if not found) and any error.
	GetCredentialBalance(ctx sdk.Context, collectionId uint64, tokenId uint64, address string) (uint64, error)

	// GetTotalSupply returns the total supply of a specific credential token.
	// collectionId and tokenId identify the credential token.
	// Returns the total supply and any error.
	GetTotalSupply(ctx sdk.Context, collectionId uint64, tokenId uint64) (uint64, error)
}
