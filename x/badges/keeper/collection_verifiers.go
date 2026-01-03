package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CollectionVerifier defines the interface for verifying collections before they are stored.
// All implementations should return an error if the collection fails verification.
// This allows developers to add custom validation logic that runs before collections are stored.
type CollectionVerifier interface {
	Name() string
	VerifyCollection(ctx sdk.Context, collection *types.TokenCollection) error
}

// NoOpCollectionVerifier is a no-op implementation of CollectionVerifier for use as a placeholder
// This can be used in app.go as a starting point that developers can replace with their own implementation
type NoOpCollectionVerifier struct{}

// Name returns the name of this verifier
func (v *NoOpCollectionVerifier) Name() string {
	return "NoOpCollectionVerifier"
}

// VerifyCollection performs no verification and always returns nil
func (v *NoOpCollectionVerifier) VerifyCollection(ctx sdk.Context, collection *types.TokenCollection) error {
	// No-op: always passes verification
	return nil
}
