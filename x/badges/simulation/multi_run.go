package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

// MultiRunOperation wraps an operation to run it multiple times with different parameters
// Returns the first successful operation, or the last attempt if all fail
func MultiRunOperation(
	op simtypes.Operation,
	numRuns int,
) simtypes.Operation {
	if numRuns <= 0 {
		numRuns = DefaultMultiRunAttempts
	}
	
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var lastOpMsg simtypes.OperationMsg
		var lastFutureOps []simtypes.FutureOperation
		var lastErr error
		
		// Try the operation multiple times with different random seeds
		for attempt := 0; attempt < numRuns; attempt++ {
			// Create a new random source for this attempt to get different parameters
			// Use the original seed plus attempt number to get variation
			attemptRand := rand.New(rand.NewSource(r.Int63() + int64(attempt)))
			
			opMsg, futureOps, err := op(attemptRand, app, ctx, accs, chainID)
			
			// If we got a successful operation (not NoOpMsg), return it immediately
			if err == nil && opMsg.OK {
				return opMsg, futureOps, err
			}
			
			// Store the last attempt for fallback
			lastOpMsg = opMsg
			lastFutureOps = futureOps
			lastErr = err
		}
		
		// Return the last attempt if all failed
		return lastOpMsg, lastFutureOps, lastErr
	}
}

