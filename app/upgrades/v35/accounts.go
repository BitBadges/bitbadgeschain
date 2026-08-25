package v35

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	vestingexported "github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// AccountMigrationResult reports what RescaleVestingAccounts touched.
type AccountMigrationResult struct {
	Scanned  int
	Rescaled int
}

// RescaleVestingAccounts converts the amounts recorded on vesting accounts.
//
// A vesting account's spendable balance is derived, not stored: x/auth computes
// it as (bank balance - locked), where locked comes from OriginalVesting and the
// vesting schedule held on the account itself. RedenominateBank scales the bank
// balance; if the account's own figures stay at the old scale, the account
// suddenly looks 10^9 times richer in spendable terms and the vesting schedule
// stops constraining anything.
//
// That is the worst possible direction for this bug: it silently unlocks funds
// that are supposed to be locked, with no error and nothing in the logs.
//
// Every vesting account type embeds BaseVestingAccount, so OriginalVesting,
// DelegatedFree and DelegatedVesting are common to all of them. Periodic and
// permanent-locked accounts add their own schedules, handled below.
func RescaleVestingAccounts(ctx sdk.Context, ak authkeeper.AccountKeeper) (AccountMigrationResult, error) {
	var res AccountMigrationResult

	// Collect before writing: IterateAccounts walks an open store iterator.
	//
	// Selected by interface rather than by a list of concrete types on purpose.
	// An allowlist of the four SDK vesting types silently skips anything else
	// that locks coins — a chain-specific or newly added vesting account — and
	// skipping is the dangerous direction here: the account's own figures stay
	// at the old scale while its bank balance grows 10^9x, which unlocks funds
	// that are supposed to be locked. Selecting by interface makes the default
	// branch of the type switch below reachable, so an unhandled type fails the
	// upgrade instead of quietly passing.
	var vesting []sdk.AccountI
	ak.IterateAccounts(ctx, func(acc sdk.AccountI) bool {
		res.Scanned++
		if _, ok := acc.(vestingexported.VestingAccount); ok {
			vesting = append(vesting, acc)
		}
		return false
	})

	for _, acc := range vesting {
		switch v := acc.(type) {
		case *vestingtypes.ContinuousVestingAccount:
			scaleBaseVesting(v.BaseVestingAccount)
		case *vestingtypes.DelayedVestingAccount:
			scaleBaseVesting(v.BaseVestingAccount)
		case *vestingtypes.PermanentLockedAccount:
			scaleBaseVesting(v.BaseVestingAccount)
		case *vestingtypes.PeriodicVestingAccount:
			scaleBaseVesting(v.BaseVestingAccount)
			// Each period releases a fixed amount; the periods must sum to
			// OriginalVesting, so they scale with it or the schedule stops
			// adding up and the tail of the vesting never unlocks.
			for i := range v.VestingPeriods {
				v.VestingPeriods[i].Amount = convertCoins(v.VestingPeriods[i].Amount, legacyDenom(), newDenom())
			}
		default:
			// Reachable: the collection above selects every VestingAccount, not
			// just the four handled here. See gamm.go for the same stance on
			// unrecognised pool types.
			return res, fmt.Errorf(
				"unhandled vesting account type %T for %s: refusing to leave its locked amounts unconverted",
				acc, acc.GetAddress())
		}

		ak.SetAccount(ctx, acc)
		res.Rescaled++
	}

	ctx.Logger().Info(
		"v35: rescaled vesting accounts",
		"factor", ConversionFactor.String(),
		"scanned", res.Scanned,
		"rescaled", res.Rescaled,
	)
	return res, nil
}

func scaleBaseVesting(b *vestingtypes.BaseVestingAccount) {
	if b == nil {
		return
	}
	b.OriginalVesting = convertCoins(b.OriginalVesting, legacyDenom(), newDenom())
	b.DelegatedFree = convertCoins(b.DelegatedFree, legacyDenom(), newDenom())
	b.DelegatedVesting = convertCoins(b.DelegatedVesting, legacyDenom(), newDenom())
}
