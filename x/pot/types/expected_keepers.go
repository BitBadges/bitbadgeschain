package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ---------------------------------------------------------------------------
// Abstract validator set interface (supports both x/staking and PoA)
// ---------------------------------------------------------------------------

// ValidatorInfo contains the information x/pot needs about a validator.
type ValidatorInfo struct {
	ConsAddr     sdk.ConsAddress
	OperatorAddr string // account address (bb1...) for credential lookup
	Power        int64  // current voting power
}

// ValidatorSetKeeper abstracts the validator set management layer.
// Implementations exist for both x/staking and PoA.
type ValidatorSetKeeper interface {
	// IterateActiveValidators walks all validators with positive power.
	IterateActiveValidators(ctx context.Context, fn func(val ValidatorInfo) bool) error

	// GetValidatorByConsAddr looks up a validator by consensus address.
	// Returns the validator info or an error if not found.
	GetValidatorByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (ValidatorInfo, error)

	// DisableValidator removes a validator from the active set.
	// For x/staking: calls Jail(). For PoA: sets power to 0.
	DisableValidator(ctx context.Context, consAddr sdk.ConsAddress) error

	// EnableValidator restores a validator to the active set.
	// For x/staking: calls Unjail(). For PoA: restores saved power.
	EnableValidator(ctx context.Context, consAddr sdk.ConsAddress) error

	// IsValidatorJailed returns whether the validator is currently disabled/jailed.
	IsValidatorJailed(ctx context.Context, consAddr sdk.ConsAddress) (bool, error)

	// CanSafelyEnable checks whether EnableValidator is safe to call.
	// For x/staking: checks not tombstoned and JailedUntil has passed.
	// For PoA: always returns true (no slashing concept).
	CanSafelyEnable(ctx context.Context, consAddr sdk.ConsAddress) bool
}

// PowerStore abstracts the persistence of saved validator power for the PoA adapter.
// The Keeper implements this interface so the PoA adapter can store/retrieve
// saved power without a circular dependency.
type PowerStore interface {
	GetSavedPower(ctx context.Context, consAddr sdk.ConsAddress) (int64, bool)
	SetSavedPower(ctx context.Context, consAddr sdk.ConsAddress, power int64)
	RemoveSavedPower(ctx context.Context, consAddr sdk.ConsAddress)
}

// ---------------------------------------------------------------------------
// Concrete keeper interfaces (used internally by the StakingAdapter)
// ---------------------------------------------------------------------------

// StakingKeeper defines the interface for the x/staking keeper that x/pot needs.
// Used by StakingValidatorSetAdapter — not referenced directly by the Keeper struct.
type StakingKeeper interface {
	// GetValidator returns the validator for the given operator address.
	GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error)

	// IterateBondedValidatorsByPower iterates over bonded validators sorted by power.
	IterateBondedValidatorsByPower(ctx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) error

	// IterateValidators iterates over all validators.
	IterateValidators(ctx context.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool)) error

	// PowerReduction returns the power reduction factor (tokens per unit of consensus power).
	PowerReduction(ctx context.Context) sdkmath.Int

	// Jail sets Jailed=true on a validator, causing x/staking's EndBlocker to
	// remove them from the active set (power → 0) via ValidatorUpdates.
	// Jailing does NOT slash — slashing is a separate mechanism.
	Jail(ctx context.Context, consAddr sdk.ConsAddress) error

	// Unjail sets Jailed=false on a validator, allowing x/staking's EndBlocker
	// to re-add them to the active set if they have enough stake.
	Unjail(ctx context.Context, consAddr sdk.ConsAddress) error

	// GetValidatorByConsAddr returns a validator by consensus address.
	GetValidatorByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (stakingtypes.Validator, error)
}

// SlashingKeeper defines the interface for the x/slashing keeper that x/pot needs.
// Used by StakingValidatorSetAdapter — not referenced directly by the Keeper struct.
type SlashingKeeper interface {
	// GetValidatorSigningInfo returns the signing info for a validator by consensus address.
	GetValidatorSigningInfo(ctx context.Context, address sdk.ConsAddress) (slashingtypes.ValidatorSigningInfo, error)

	// IsTombstoned returns true if a validator has been permanently removed from the validator set.
	IsTombstoned(ctx context.Context, consAddr sdk.ConsAddress) bool
}

// TokenizationKeeper defines the interface for the x/tokenization keeper that x/pot needs.
// We define a minimal interface so x/pot does not import x/tokenization's concrete types.
//
// The wiring code in app.go should create an adapter struct that implements this interface
// using the concrete x/tokenization keeper methods (GetCollectionFromStore, GetBalanceOrApplyDefault, etc.).
type TokenizationKeeper interface {
	// GetCredentialBalance returns the balance of a specific token for a given address.
	// collectionId and tokenId identify the credential token to check.
	// Returns the balance amount (0 if not found) and any error.
	GetCredentialBalance(ctx sdk.Context, collectionId uint64, tokenId uint64, address string) (uint64, error)
}
