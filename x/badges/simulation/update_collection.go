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
			Creator:                          simAccount.Address.String(),
			UpdateCollectionPermissions:      r.Int63n(2) == 0,
			UpdateIsArchivedTimeline:         r.Int63n(2) == 0,
			UpdateManagerTimeline:            r.Int63n(2) == 0,
			UpdateCollectionMetadataTimeline: r.Int63n(2) == 0,
			UpdateBadgeMetadataTimeline:      r.Int63n(2) == 0,
			UpdateCustomDataTimeline:         r.Int63n(2) == 0,
			UpdateCollectionApprovals:        r.Int63n(2) == 0,
			UpdateStandardsTimeline:          r.Int63n(2) == 0,
			UpdateValidBadgeIds:              r.Int63n(2) == 0,

			CollectionId: sdkmath.NewUint(uint64(r.Int63n(5))),
			IsArchivedTimeline: []*types.IsArchivedTimeline{
				{
					IsArchived:    r.Int63n(2) == 0,
					TimelineTimes: GetTimelineTimes(r, 3),
				},
			},
			ValidBadgeIds: GetTimelineTimes(r, 3),
			CollectionApprovals: []*types.CollectionApproval{
				{
					FromListId:        GetRandomAddresses(r, 1, accs)[0],
					ToListId:          GetRandomAddresses(r, 1, accs)[0],
					InitiatedByListId: GetRandomAddresses(r, 1, accs)[0],
					TransferTimes:     GetTimelineTimes(r, 100),
					OwnershipTimes:    GetTimelineTimes(r, 100),
					BadgeIds:          GetTimelineTimes(r, 3),
				},
			},
			ManagerTimeline: []*types.ManagerTimeline{
				{
					Manager:       simAccount.Address.String(),
					TimelineTimes: GetTimelineTimes(r, 3),
				},
			},
			CollectionPermissions: GetRandomCollectionPermissions(r, accs),
		}

		// TODO: Handling the UpdateCollection simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "UpdateCollection simulation not implemented"), nil, nil
	}
}
