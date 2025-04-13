package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/anchor/keeper"
	"github.com/bitbadges/bitbadgeschain/x/anchor/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgAddCustomData(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgAddCustomData{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the AddCustomData simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "AddCustomData simulation not implemented"), nil, nil
	}
}
