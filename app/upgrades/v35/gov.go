package v35

import (
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
	)

	return res, nil
}
