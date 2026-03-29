package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/pot/types"
)

// EndBlocker is called at the end of every block.
//
// Architecture (v2 — jail-based, no ValidatorUpdates):
//
//  1. Iterate ALL active validators via the ValidatorSetKeeper abstraction.
//  2. For each validator, check their credential token balance at the current block time.
//     This catches both transfer-based and time-dependent credential expiry.
//  3. Validators without credentials are disabled via ValidatorSetKeeper.DisableValidator().
//     For x/staking this calls Jail(); for PoA it sets power to 0.
//  4. Validators who regained credentials and were compliance-jailed by x/pot
//     are auto-enabled via ValidatorSetKeeper.EnableValidator(), but only if
//     CanSafelyEnable() returns true (respects tombstone/slashing jail for staking).
//  5. x/pot tracks its own "compliance-jailed" set in persistent state to
//     distinguish compliance-jailing from slashing-jailing.
//
// SAFETY:
//   - Never disables ALL validators. If disabling would leave zero active validators,
//     the remaining non-compliant validators are left enabled with a warning log.
//   - All state changes are done via CacheContext for atomicity — if endBlockerInner
//     returns an error, all changes are rolled back.
//   - Individual Disable/Enable calls are wrapped in per-validator recover blocks
//     as defense-in-depth.
//
// This EndBlocker returns only an error (via appmodule.HasEndBlocker), NOT
// ValidatorUpdates. This avoids the "validator EndBlock updates already set"
// panic that occurs when two modules return ValidatorUpdates.
func (k Keeper) EndBlocker(ctx sdk.Context) error {
	cacheCtx, writeCache := ctx.CacheContext()
	err := k.endBlockerInner(cacheCtx)
	if err != nil {
		k.Logger().Error("pot: EndBlocker failed, changes rolled back", "error", err)
		return nil // don't propagate error, just skip this block's changes
	}
	writeCache() // commit all changes atomically
	return nil
}

// endBlockerInner contains the core EndBlocker logic. It runs inside a CacheContext
// so that any error causes all state changes to be rolled back atomically.
func (k Keeper) endBlockerInner(ctx sdk.Context) error {
	params := k.GetParams(ctx)

	// If the module is not enabled (no credential collection configured), skip.
	if !params.IsEnabled() {
		return nil
	}

	// Collect all active validators and check credentials.
	type valInfo struct {
		operatorAddr  string
		consAddr      sdk.ConsAddress
		hasCredential bool
	}

	var validators []valInfo

	err := k.validatorSet.IterateActiveValidators(ctx, func(vi types.ValidatorInfo) bool {
		// Check credential balance using the validator's account address (bb1...).
		balance, err := k.tokenizationKeeper.GetCredentialBalance(
			ctx,
			params.CredentialCollectionId,
			params.CredentialTokenId,
			vi.OperatorAddr,
		)
		if err != nil {
			k.Logger().Error("pot: failed to get credential balance",
				"operator", vi.OperatorAddr, "error", err)
			// Treat errors as "has credential" to avoid disabling validators
			// due to a keeper bug.
			validators = append(validators, valInfo{
				operatorAddr:  vi.OperatorAddr,
				consAddr:      vi.ConsAddr,
				hasCredential: true,
			})
			return false
		}

		validators = append(validators, valInfo{
			operatorAddr:  vi.OperatorAddr,
			consAddr:      vi.ConsAddr,
			hasCredential: balance >= params.MinCredentialBalance,
		})
		return false // continue
	})
	if err != nil {
		k.Logger().Error("pot: failed to iterate active validators", "error", err)
		return nil
	}

	// Also check compliance-jailed validators that may have regained credentials.
	// These are no longer active, so IterateActiveValidators won't see them.
	complianceJailedAddrs := k.GetAllComplianceJailed(ctx)
	for _, addrBytes := range complianceJailedAddrs {
		consAddr := sdk.ConsAddress(addrBytes)

		// Check if this validator is already in our active list (shouldn't be, but guard).
		alreadySeen := false
		for _, vi := range validators {
			if vi.consAddr.Equals(consAddr) {
				alreadySeen = true
				break
			}
		}
		if alreadySeen {
			continue
		}

		// Look up the validator by consensus address (O(1) via the adapter).
		foundVal, err := k.validatorSet.GetValidatorByConsAddr(ctx, consAddr)
		if err != nil {
			// Validator no longer exists, clean up stale entry.
			k.RemoveComplianceJailed(ctx, consAddr)
			continue
		}

		balance, err := k.tokenizationKeeper.GetCredentialBalance(
			ctx,
			params.CredentialCollectionId,
			params.CredentialTokenId,
			foundVal.OperatorAddr,
		)
		if err != nil {
			k.Logger().Error("pot: failed to get credential balance for jailed validator",
				"operator", foundVal.OperatorAddr, "error", err)
			continue
		}

		if balance >= params.MinCredentialBalance {
			validators = append(validators, valInfo{
				operatorAddr:  foundVal.OperatorAddr,
				consAddr:      consAddr,
				hasCredential: true,
			})
		}
	}

	// SAFETY: Count how many active validators will remain after disabling.
	// Never disable ALL validators — that would halt the chain.
	bondedCount := 0
	toDisableCount := 0
	for _, vi := range validators {
		if !k.IsComplianceJailed(ctx, vi.consAddr) && !vi.hasCredential {
			toDisableCount++
		}
		// Count validators that are currently active (not already compliance-jailed).
		if !k.IsComplianceJailed(ctx, vi.consAddr) {
			bondedCount++
		}
	}

	// If disabling all non-compliant validators would leave zero active validators,
	// skip disabling entirely and log a critical warning.
	safeToDisable := true
	if bondedCount > 0 && toDisableCount >= bondedCount {
		k.Logger().Error("pot: SAFETY — refusing to disable all validators, would halt chain",
			"bonded", bondedCount,
			"to_disable", toDisableCount,
		)
		safeToDisable = false
	}

	// Process each validator.
	for _, vi := range validators {
		isComplianceJailed := k.IsComplianceJailed(ctx, vi.consAddr)

		switch {
		case vi.hasCredential && isComplianceJailed:
			// Regained credential — check if safe to enable, then try.
			if !k.validatorSet.CanSafelyEnable(ctx, vi.consAddr) {
				// For staking: tombstoned or slashing jail still active.
				// Check if tombstoned specifically to clean up.
				jailed, _ := k.validatorSet.IsValidatorJailed(ctx, vi.consAddr)
				if !jailed {
					// Validator is no longer jailed but CanSafelyEnable is false
					// — likely tombstoned. Clean up the stale entry.
					k.RemoveComplianceJailed(ctx, vi.consAddr)
					k.Logger().Info("pot: removed permanently disabled validator from compliance-jailed set",
						"operator", vi.operatorAddr,
					)
				} else {
					k.Logger().Info("pot: validator regained credential but not safe to enable yet, deferring",
						"operator", vi.operatorAddr,
					)
				}
			} else {
				k.safeEnable(ctx, vi.consAddr, vi.operatorAddr)
				k.RemoveComplianceJailed(ctx, vi.consAddr)
				k.Logger().Info("pot: auto-enabled validator (credential regained)",
					"operator", vi.operatorAddr,
				)
			}

		case !vi.hasCredential && !isComplianceJailed && safeToDisable:
			// Lost credential — disable for compliance.
			// Only mark compliance-jailed if disable actually succeeded.
			if k.safeDisable(ctx, vi.consAddr, vi.operatorAddr) {
				k.SetComplianceJailed(ctx, vi.consAddr)
				k.Logger().Info("pot: disabled validator for missing credential",
					"operator", vi.operatorAddr,
				)
			}

		// case vi.hasCredential && !isComplianceJailed: no action needed
		// case !vi.hasCredential && isComplianceJailed: already disabled, no action
		}
	}

	return nil
}

// safeDisable calls validatorSet.DisableValidator wrapped in a recover to prevent panics.
// Uses an inner CacheContext so partial writes from a panic are discarded.
// Returns true if the disable call succeeded, false otherwise.
func (k Keeper) safeDisable(ctx sdk.Context, consAddr sdk.ConsAddress, operatorAddr string) (success bool) {
	innerCtx, writeInner := ctx.CacheContext()
	defer func() {
		if r := recover(); r != nil {
			k.Logger().Error("pot: recovered from panic during DisableValidator",
				"operator", operatorAddr,
				"panic", fmt.Sprintf("%v", r),
			)
			success = false
			// writeInner is never called, so partial state is discarded
		}
	}()

	if err := k.validatorSet.DisableValidator(innerCtx, consAddr); err != nil {
		k.Logger().Warn("pot: failed to disable validator",
			"operator", operatorAddr,
			"error", err,
		)
		return false
	}
	writeInner() // commit only on success
	return true
}

// safeEnable calls validatorSet.EnableValidator wrapped in a recover to prevent panics.
// Uses an inner CacheContext so partial writes from a panic are discarded.
func (k Keeper) safeEnable(ctx sdk.Context, consAddr sdk.ConsAddress, operatorAddr string) {
	innerCtx, writeInner := ctx.CacheContext()
	defer func() {
		if r := recover(); r != nil {
			k.Logger().Error("pot: recovered from panic during EnableValidator",
				"operator", operatorAddr,
				"panic", fmt.Sprintf("%v", r),
			)
		}
	}()

	if err := k.validatorSet.EnableValidator(innerCtx, consAddr); err != nil {
		k.Logger().Error("pot: failed to enable validator",
			"operator", operatorAddr,
			"error", err,
		)
		return
	}
	writeInner() // commit only on success
}
