package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgUniversalUpdateCollection(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgUniversalUpdateCollection{
			Creator:                     simAccount.Address.String(),
			UpdateCollectionPermissions: r.Int63n(2) == 0,
			UpdateIsArchived:            r.Int63n(2) == 0,
			UpdateManager:               r.Int63n(2) == 0,
			UpdateCollectionMetadata:    r.Int63n(2) == 0,
			UpdateTokenMetadata:         r.Int63n(2) == 0,
			UpdateCustomData:            r.Int63n(2) == 0,
			UpdateCollectionApprovals:   r.Int63n(2) == 0,
			UpdateStandards:             r.Int63n(2) == 0,
			UpdateValidTokenIds:         r.Int63n(2) == 0,

			CollectionId:  sdkmath.NewUint(uint64(r.Int63n(5))),
			IsArchived:    r.Int63n(2) == 0,
			ValidTokenIds: GetTimelineTimes(r, 3),
			CollectionApprovals: []*types.CollectionApproval{
				{
					FromListId:        GetRandomAddresses(r, 1, accs)[0],
					ToListId:          GetRandomAddresses(r, 1, accs)[0],
					InitiatedByListId: GetRandomAddresses(r, 1, accs)[0],
					TransferTimes:     GetTimelineTimes(r, 100),
					OwnershipTimes:    GetTimelineTimes(r, 100),
					TokenIds:          GetTimelineTimes(r, 3),
				},
			},
			Manager:               simAccount.Address.String(),
			CollectionPermissions: GetRandomCollectionPermissions(r, accs),
		}

		// TODO: Handling the UpdateCollection simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "UpdateCollection simulation not implemented"), nil, nil
	}
}
