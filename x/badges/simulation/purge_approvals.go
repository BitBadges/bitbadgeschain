package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgPurgeApprovals(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Ensure we have valid accounts
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPurgeApprovals, "no accounts available"), nil, nil
		}

		// Try to get a known-good collection ID first
		collectionId, found := GetKnownGoodCollectionId(ctx, k)
		if !found {
			// Fallback: try to get a random existing collection
			collectionId, found = GetRandomCollectionId(r, ctx, k)
			if !found {
				// Try to create one first
				simAccount := EnsureAccountExists(r, accs)
				createdId, err := GetOrCreateCollection(ctx, k, simAccount.Address.String(), r, accs)
				if err != nil {
					return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPurgeApprovals, "no collections exist and failed to create one"), nil, nil
				}
				collectionId = createdId
			}
		}

		// Check if collection exists
		_, found = k.GetCollectionFromStore(ctx, collectionId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgPurgeApprovals, "collection not found"), nil, nil
		}

		simAccount := EnsureAccountExists(r, accs)

		// Generate specific approvals to purge (required - cannot be empty)
		// Randomly decide approval level
		approvalLevel := "collection"
		approverAddress := ""
		if r.Intn(3) == 0 {
			approvalLevel = "incoming"
			approverAddress = simAccount.Address.String()
		} else if r.Intn(2) == 0 {
			approvalLevel = "outgoing"
			approverAddress = simAccount.Address.String()
		}

		// Generate 1-3 approvals to purge
		count := r.Intn(3) + 1
		approvalsToPurge := []*types.ApprovalIdentifierDetails{}
		for i := 0; i < count; i++ {
			approvalsToPurge = append(approvalsToPurge, &types.ApprovalIdentifierDetails{
				ApprovalId:      simtypes.RandStringOfLength(r, 10),
				ApprovalLevel:   approvalLevel,
				ApproverAddress: approverAddress,
			})
		}

		// Determine purge options based on whether we're purging own approvals
		// If approverAddress is empty or matches creator, we're purging own approvals
		purgeExpired := true                // Required when purging own approvals
		purgeCounterpartyApprovals := false // Must be false when purging own approvals
		if approverAddress != "" && approverAddress != simAccount.Address.String() {
			// Purging someone else's approvals
			purgeExpired = r.Intn(2) == 0
			purgeCounterpartyApprovals = r.Intn(2) == 0
		}

		msg := &types.MsgPurgeApprovals{
			Creator:                    simAccount.Address.String(),
			CollectionId:               collectionId,
			PurgeExpired:               purgeExpired,
			ApproverAddress:            approverAddress, // Empty defaults to creator
			PurgeCounterpartyApprovals: purgeCounterpartyApprovals,
			ApprovalsToPurge:           approvalsToPurge,
		}

		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
