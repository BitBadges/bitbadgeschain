package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgCreateAddressLists(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgCreateAddressLists{
			Creator: simAccount.Address.String(),
			AddressLists: []*types.AddressList{
				{
					Addresses:  []string{simAccount.Address.String()},
					Uri:        "",
					CustomData: "",
					ListId:     simtypes.RandStringOfLength(r, 10),
					CreatedBy:  simAccount.Address.String(),
				},
			},
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
