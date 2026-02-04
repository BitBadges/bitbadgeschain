package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgCreateCollection(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Ensure we have valid accounts
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateCollection, "no accounts available"), nil, nil
		}

		simAccount := EnsureAccountExists(r, accs)

		// Generate valid token IDs (at least one required) using bounded ranges
		validTokenIds := GetBoundedTimelineTimes(r, r.Intn(3)+1, MinTimelineRange, MaxTimelineRange)

		// Generate collection permissions
		collectionPermissions := GetRandomCollectionPermissions(r, accs)

		// Generate collection metadata
		collectionMetadata := &types.CollectionMetadata{
			Uri:        "https://example.com/metadata/" + simtypes.RandStringOfLength(r, 10),
			CustomData: simtypes.RandStringOfLength(r, 20),
		}

		// Generate random collection approvals (optional, sometimes include mint approval)
		collectionApprovals := []*types.CollectionApproval{}
		if r.Intn(3) == 0 {
			// Add a mint approval
			mintApproval := GetRandomCollectionApproval(r, accs)
			mintApproval.FromListId = types.MintAddress
			mintApproval.ToListId = "All"
			mintApproval.InitiatedByListId = "All"
			if mintApproval.ApprovalCriteria == nil {
				mintApproval.ApprovalCriteria = &types.ApprovalCriteria{}
			}
			mintApproval.ApprovalCriteria.OverridesFromOutgoingApprovals = true
			mintApproval.ApprovalCriteria.OverridesToIncomingApprovals = true
			collectionApprovals = append(collectionApprovals, mintApproval)
		}
		// Sometimes add additional approvals
		if r.Intn(2) == 0 {
			collectionApprovals = append(collectionApprovals, GetRandomCollectionApproval(r, accs))
		}

		// Default balances should be empty (zero amounts are invalid)
		defaultBalances := &types.UserBalanceStore{
			Balances: []*types.Balance{},
		}

		msg := &types.MsgCreateCollection{
			Creator:               simAccount.Address.String(),
			DefaultBalances:       defaultBalances,
			ValidTokenIds:         validTokenIds,
			CollectionPermissions: collectionPermissions,
			Manager:               simAccount.Address.String(), // Creator is manager
			CollectionMetadata:    collectionMetadata,
			TokenMetadata:         []*types.TokenMetadata{},
			CustomData:            simtypes.RandStringOfLength(r, 20),
			CollectionApprovals:   collectionApprovals,
			Standards:             []string{},
			IsArchived:            false,
		}

		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
