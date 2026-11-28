package types

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
)

// ValidateApprovalCriteriaErgonomics enforces strict BB-375 §2.2 mutual-exclusivity
// and satisfiability invariants on ApprovalCriteria without modifying permissions.proto.
func ValidateApprovalCriteriaErgonomics(criteria *ApprovalCriteria) error {
	if criteria == nil {
		return nil
	}

	// 1. PredeterminedBalances XOR ApprovalAmounts
	hasPredetermined := criteria.PredeterminedBalances != nil &&
		!PredeterminedBalancesIsBasicallyNil(criteria.PredeterminedBalances)
	hasApprovalAmounts := criteria.ApprovalAmounts != nil &&
		!ApprovalAmountsIsBasicallyNil(criteria.ApprovalAmounts)
	if hasPredetermined && hasApprovalAmounts {
		return sdkerrors.Wrap(
			ErrInvalidRequest,
			"predeterminedBalances and approvalAmounts are mutually exclusive and cannot both be active",
		)
	}

	// 2. DurationFromTimestamp XOR RecurringOwnershipTimes XOR IncrementOwnershipTimesBy
	if criteria.PredeterminedBalances != nil && criteria.PredeterminedBalances.IncrementedBalances != nil {
		inc := criteria.PredeterminedBalances.IncrementedBalances
		hasDuration := !inc.DurationFromTimestamp.IsNil() && !inc.DurationFromTimestamp.IsZero()
		hasIncrement := !inc.IncrementOwnershipTimesBy.IsNil() && !inc.IncrementOwnershipTimesBy.IsZero()
		hasRecurring := inc.RecurringOwnershipTimes != nil &&
			((!inc.RecurringOwnershipTimes.StartTime.IsNil() && !inc.RecurringOwnershipTimes.StartTime.IsZero()) ||
				(!inc.RecurringOwnershipTimes.IntervalLength.IsNil() && !inc.RecurringOwnershipTimes.IntervalLength.IsZero()) ||
				(!inc.RecurringOwnershipTimes.ChargePeriodLength.IsNil() && !inc.RecurringOwnershipTimes.ChargePeriodLength.IsZero()))

		activeTimeModes := 0
		if hasDuration {
			activeTimeModes++
		}
		if hasIncrement {
			activeTimeModes++
		}
		if hasRecurring {
			activeTimeModes++
		}
		if activeTimeModes > 1 {
			return sdkerrors.Wrap(
				ErrInvalidRequest,
				"durationFromTimestamp, recurringOwnershipTimes, and incrementOwnershipTimesBy are mutually exclusive",
			)
		}
	}

	// 3. Contradictory initiator party-relation booleans
	if criteria.RequireToEqualsInitiatedBy && criteria.RequireToDoesNotEqualInitiatedBy {
		return sdkerrors.Wrap(
			ErrInvalidRequest,
			"requireToEqualsInitiatedBy and requireToDoesNotEqualInitiatedBy cannot both be true",
		)
	}
	if criteria.RequireFromEqualsInitiatedBy && criteria.RequireFromDoesNotEqualInitiatedBy {
		return sdkerrors.Wrap(
			ErrInvalidRequest,
			"requireFromEqualsInitiatedBy and requireFromDoesNotEqualInitiatedBy cannot both be true",
		)
	}

	// 4. Contradictory AddressChecks booleans
	if err := validateAddressChecksConsistency("senderChecks", criteria.SenderChecks); err != nil {
		return err
	}
	if err := validateAddressChecksConsistency("recipientChecks", criteria.RecipientChecks); err != nil {
		return err
	}
	if err := validateAddressChecksConsistency("initiatorChecks", criteria.InitiatorChecks); err != nil {
		return err
	}

	return nil
}

// ValidateOutgoingApprovalCriteriaErgonomics validates an OutgoingApprovalCriteria
// using the canonical collection-criteria ergonomics rules.
func ValidateOutgoingApprovalCriteriaErgonomics(criteria *OutgoingApprovalCriteria) error {
	if criteria == nil {
		return nil
	}
	return ValidateApprovalCriteriaErgonomics(CastOutgoingApprovalCriteriaToCollectionApprovalCriteria(criteria))
}

// ValidateIncomingApprovalCriteriaErgonomics validates an IncomingApprovalCriteria
// using the canonical collection-criteria ergonomics rules.
func ValidateIncomingApprovalCriteriaErgonomics(criteria *IncomingApprovalCriteria) error {
	if criteria == nil {
		return nil
	}
	return ValidateApprovalCriteriaErgonomics(CastIncomingApprovalCriteriaToCollectionApprovalCriteria(criteria))
}

func validateAddressChecksConsistency(field string, checks *AddressChecks) error {
	if checks == nil {
		return nil
	}
	if checks.MustBeEvmContract && checks.MustNotBeEvmContract {
		return sdkerrors.Wrapf(ErrInvalidRequest, "%s: mustBeEvmContract and mustNotBeEvmContract cannot both be true", field)
	}
	if checks.MustBeLiquidityPool && checks.MustNotBeLiquidityPool {
		return sdkerrors.Wrapf(ErrInvalidRequest, "%s: mustBeLiquidityPool and mustNotBeLiquidityPool cannot both be true", field)
	}
	return nil
}

// ValidateAddressListPredicate rejects malformed string-DSL address list IDs and
// delimiter-injection patterns (BB-375 §2.2 / H-14).
func ValidateAddressListPredicate(listId string) error {
	if listId == "" {
		return sdkerrors.Wrap(ErrInvalidAddress, "address list predicate cannot be empty")
	}
	for _, r := range listId {
		if r <= 0x20 || r == 0x7f {
			return sdkerrors.Wrapf(ErrInvalidAddress, "address list predicate %q contains whitespace or control characters", listId)
		}
	}

	body := listId
	if strings.HasPrefix(body, "!") {
		if len(body) == 1 {
			return sdkerrors.Wrap(ErrInvalidAddress, "negated address list predicate cannot be bare '!'")
		}
		if strings.HasPrefix(body, "!(") {
			if !strings.HasSuffix(body, ")") {
				return sdkerrors.Wrapf(ErrInvalidAddress, "address list predicate %q has unclosed '!(' group", listId)
			}
			body = body[2 : len(body)-1]
			if body == "" {
				return sdkerrors.Wrap(ErrInvalidAddress, "address list predicate '!()' cannot be empty")
			}
		} else {
			body = body[1:]
		}
		if strings.HasPrefix(body, "!") {
			return sdkerrors.Wrapf(ErrInvalidAddress, "address list predicate %q cannot use double negation", listId)
		}
	}

	if strings.ContainsAny(body, "()") {
		return sdkerrors.Wrapf(ErrInvalidAddress, "address list predicate %q contains unbalanced or nested parentheses", listId)
	}

	if strings.HasPrefix(body, "AllWithout") && body != "AllWithoutMint" {
		rest := strings.TrimPrefix(body, "AllWithout")
		if rest == "" {
			return sdkerrors.Wrap(ErrInvalidAddress, "AllWithout predicate requires at least one address")
		}
		body = rest
	}

	if strings.Contains(body, ":") {
		parts := strings.Split(body, ":")
		for _, part := range parts {
			if part == "" {
				return sdkerrors.Wrapf(ErrInvalidAddress, "address list predicate %q contains empty colon-delimited token", listId)
			}
			if strings.Contains(part, "!") {
				return sdkerrors.Wrapf(ErrInvalidAddress, "address list predicate %q contains embedded '!' delimiter injection", listId)
			}
		}
	}

	return nil
}

// NormalizeCollectionApprovalDefaults applies BB-375 §2.3 sane protocol defaults
// in-place to a CollectionApproval, eliminating common developer and AI agent traps:
//  1. FromListId == "Mint" automatically sets OverridesFromOutgoingApprovals = true.
//  2. AllowBackedMinting == true or AllowSpecialWrapping == true automatically sets MustPrioritize = true.
//  3. Non-auto-scannable approvals normalize MustPrioritize = true.
func NormalizeCollectionApprovalDefaults(approval *CollectionApproval) {
	if approval == nil {
		return
	}

	if approval.FromListId == MintAddress {
		if approval.ApprovalCriteria == nil {
			approval.ApprovalCriteria = &ApprovalCriteria{}
		}
		approval.ApprovalCriteria.OverridesFromOutgoingApprovals = true
	}

	if approval.ApprovalCriteria == nil {
		return
	}

	if approval.ApprovalCriteria.AllowBackedMinting || approval.ApprovalCriteria.AllowSpecialWrapping {
		approval.ApprovalCriteria.MustPrioritize = true
	}

	EnforceMustPrioritizeForNonAutoScannable(approval.ApprovalCriteria)
}

// NormalizeDefaultBalances applies BB-375 §2.3 safe defaults for a collection's
// UserBalanceStore defaultBalances: AutoApproveAllIncomingTransfers defaults to
// true unless explicitly opted out via optOutAutoApproveIncoming, preventing
// silent mint/airdrop failures for fresh recipient addresses.
func NormalizeDefaultBalances(defaults *UserBalanceStore, optOutAutoApproveIncoming bool) *UserBalanceStore {
	if defaults == nil {
		defaults = &UserBalanceStore{}
	}
	defaults.AutoApproveAllIncomingTransfers = !optOutAutoApproveIncoming
	defaults.AutoApproveSelfInitiatedOutgoingTransfers = true
	defaults.AutoApproveSelfInitiatedIncomingTransfers = true
	return defaults
}
