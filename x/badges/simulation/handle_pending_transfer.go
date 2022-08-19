package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func SimulateMsgHandlePendingTransfer(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		randInt := r.Uint64()
		randBool := false
		if randInt % 2 == 0 {
			randBool = true
		}

		randInt2 := r.Uint64()
		randBool2 := false
		if randInt2 % 2 == 0 {
			randBool2 = true
		}


		msg := &types.MsgHandlePendingTransfer{
			Creator: simAccount.Address.String(),
			Accept: randBool,
			ForcefulAccept: randBool2,
			BadgeId: r.Uint64(),
			NonceRanges:[]*types.IdRange{
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
		}

		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
