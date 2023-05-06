package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgUpdateAllowedTransfers(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		msg := &types.MsgUpdateAllowedTransfers{
			Creator:      simAccount.Address.String(),
			CollectionId: r.Uint64(),
			AllowedTransfers: []*types.TransferMapping{
				{
					From: &types.AddressesMapping{
						Addresses: []string{
							simAccount.Address.String(),
							simAccount.Address.String(),
							simAccount.Address.String(),
						},
						IncludeOnlySpecified: r.Uint64() % 2 == 0,
						ManagerOptions: r.Uint64(),
					},
				},
			},
		}

		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
