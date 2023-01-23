package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgUpdateDisallowedTransfers(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)


		msg := &types.MsgUpdateDisallowedTransfers{
			Creator: simAccount.Address.String(),
			CollectionId: r.Uint64(),
			DisallowedTransfers: []*types.TransferMapping{
				{
					From: &types.Addresses{
						AccountNums: []*types.IdRange{
							{
								Start: r.Uint64(),
								End:   r.Uint64(),
							},
							{
								Start: r.Uint64(),
								End:   r.Uint64(),
							},
							{
								Start: r.Uint64(),
								End:   r.Uint64(),
							},
						},
						ManagerOptions: types.ManagerOptions_Neutral,
					},
				},
			},
		}

		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
