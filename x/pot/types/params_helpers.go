package types

import "fmt"

const (
	// ModeStakedMultiplier means voting power = staked amount if credential balance > 0, else 0.
	// Full PoS economics — credential is a binary compliance gate.
	ModeStakedMultiplier = "staked_multiplier"

	// ModeEqual means voting power = 1 if credential balance >= MinCredentialBalance, else 0.
	// Democratic consensus — all credentialed validators have equal power.
	// Validators still self-delegate (security deposit for slashing).
	// NOTE: Not yet implemented — reserved for future use.
	ModeEqual = "equal"

	// ModeCredentialWeighted means voting power = credential token balance if >= MinCredentialBalance, else 0.
	// Power is proportional to credential token balance, decoupled from staked amount.
	// Validators still self-delegate (security deposit for slashing) but power comes from credential balance.
	// NOTE: Not yet implemented — reserved for future use.
	ModeCredentialWeighted = "credential_weighted"
)

// DefaultParams returns a default set of parameters with the module effectively disabled
// (collection ID 0 means no credential check is active).
// DefaultParams returns a default set of parameters with the module effectively disabled
// (collection ID 0 means no credential check is active).
func DefaultParams() Params {
	return Params{
		CredentialCollectionId: 0,
		CredentialTokenId:      0,
		MinCredentialBalance:   1,
		Mode:                   ModeStakedMultiplier,
	}
}

// Validate performs basic validation of params.
func (p Params) Validate() error {
	if p.Mode != ModeStakedMultiplier {
		return fmt.Errorf("invalid mode %q: only %q is currently supported", p.Mode, ModeStakedMultiplier)
	}
	if p.CredentialCollectionId > 0 && p.MinCredentialBalance == 0 {
		return fmt.Errorf("min_credential_balance must be > 0 when module is enabled")
	}
	return nil
}

// IsEnabled returns true if the module has a non-zero credential collection configured.
func (p Params) IsEnabled() bool {
	return p.CredentialCollectionId > 0
}
