package v35

import (
	"bytes"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// GovMigrationResult reports what RescaleGovDeposits touched.
type GovMigrationResult struct {
	Proposals int
	Deposits  int
	// FlaggedMessages counts in-flight proposals whose Messages name the retired
	// denom. Reported, never rewritten — see the note on RescaleGovDeposits.
	FlaggedMessages int
}

// RescaleGovDeposits converts the deposit amounts recorded on in-flight
// proposals.
//
// MigrateDenomParams already moves the min-deposit *parameters*. This handles
// the other half: deposits that have actually been paid. The coins themselves
// sit on the gov module account and are moved by RedenominateBank, but the
// per-depositor records and each proposal's TotalDeposit are gov's own
// accounting.
//
// Getting this wrong is not cosmetic. Refunds pay out from the Deposit records,
// so leaving them at the old scale would refund 10^-9 of what was paid and strand
// the remainder on the module account. And because MinDeposit scales while
// TotalDeposit would not, every proposal still in its deposit period would
// instantly look underfunded and fail to enter voting.
func RescaleGovDeposits(ctx sdk.Context, gk *govkeeper.Keeper) (GovMigrationResult, error) {
	var res GovMigrationResult

	// Collect first: these walk open store iterators, and writing back inside
	// the callback mutates the store underneath them.
	type depositEntry struct {
		key     collections.Pair[uint64, sdk.AccAddress]
		deposit govv1.Deposit
	}
	var (
		proposals []govv1.Proposal
		deposits  []depositEntry
	)

	if err := gk.Proposals.Walk(ctx, nil, func(_ uint64, proposal govv1.Proposal) (bool, error) {
		proposals = append(proposals, proposal)
		return false, nil
	}); err != nil {
		return res, fmt.Errorf("reading proposals: %w", err)
	}

	if err := gk.Deposits.Walk(ctx, nil, func(key collections.Pair[uint64, sdk.AccAddress], deposit govv1.Deposit) (bool, error) {
		deposits = append(deposits, depositEntry{key: key, deposit: deposit})
		return false, nil
	}); err != nil {
		return res, fmt.Errorf("reading deposits: %w", err)
	}

	for _, proposal := range proposals {
		if flagLegacyDenomInMessages(ctx, proposal) {
			res.FlaggedMessages++
		}

		proposal.TotalDeposit = convertCoins(proposal.TotalDeposit, legacyDenom(), newDenom())
		if err := gk.SetProposal(ctx, proposal); err != nil {
			return res, fmt.Errorf("setting proposal %d: %w", proposal.Id, err)
		}
		res.Proposals++
	}

	for _, d := range deposits {
		d.deposit.Amount = convertCoins(d.deposit.Amount, legacyDenom(), newDenom())
		if err := gk.SetDeposit(ctx, d.deposit); err != nil {
			return res, fmt.Errorf("setting deposit on proposal %d: %w", d.deposit.ProposalId, err)
		}
		res.Deposits++
	}

	ctx.Logger().Info(
		"v35: rescaled gov deposits",
		"factor", ConversionFactor.String(),
		"proposals", res.Proposals,
		"deposits", res.Deposits,
		"flagged_messages", res.FlaggedMessages,
	)

	return res, nil
}

// flagLegacyDenomInMessages warns about an undecided proposal whose payload
// names the retired denom, and reports whether it found one.
//
// A proposal's Messages are executed *after* it passes, so a proposal still in
// its deposit or voting period at the upgrade height carries an Any-packed
// payload written against the old denom into a chain that no longer has one.
// The consequences split two ways:
//
//   - Most are harmless and fail closed. A community-pool spend of ubadge, for
//     instance, executes against a denom with zero supply and simply errors; the
//     proposal fails and nothing is lost.
//   - A MsgUpdateParams is not harmless. Mainnet's four ubadge-bearing proposals
//     are all of exactly this shape (x/tokenization allowed_denoms). One of those
//     passing after the upgrade would write "ubadge" straight back into the
//     params this migration just moved, silently reverting part of it.
//
// They are flagged rather than rewritten, and rewriting is the wrong call twice
// over. Messages are arbitrary Any-packed protos, so a generic rewrite would
// have to guess which string fields are native-denom references and which are
// not — and a governance payload is a thing voters approved verbatim, not
// something an upgrade handler should edit underneath them.
//
// They are flagged rather than rejected for a sharper reason: erroring here
// would be the same free DoS the precisebank reserve check used to be. Anyone
// can file a deposit-period proposal containing the string "ubadge" for the cost
// of a minimum deposit, and it would halt the chain at the upgrade height.
//
// At the time of writing mainnet has 43 proposals, all PASSED, and none in a
// deposit or voting period — so this logs nothing today. It exists for the
// window between now and the upgrade height.
func flagLegacyDenomInMessages(ctx sdk.Context, proposal govv1.Proposal) bool {
	switch proposal.Status {
	case govv1.StatusDepositPeriod, govv1.StatusVotingPeriod:
	default:
		// Anything already passed, rejected or failed will never be executed
		// again, so its payload cannot affect post-upgrade state.
		return false
	}

	found := false
	for _, msg := range proposal.Messages {
		if msg == nil {
			continue
		}
		if bytes.Contains(msg.Value, []byte(legacyDenom())) {
			found = true
			break
		}
	}
	if !found {
		return false
	}

	ctx.Logger().Error(
		"v35: an undecided governance proposal names the retired denom in its messages; "+
			"it is NOT rewritten, and executing it after the upgrade may fail or undo part of this migration",
		"proposal", proposal.Id,
		"status", proposal.Status.String(),
		"denom", legacyDenom(),
	)
	return true
}
