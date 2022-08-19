package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func SimulateMsgRequestTransferManager(
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
		msg := &types.MsgRequestTransferManager{
			Creator: simAccount.Address.String(),
			BadgeId: r.Uint64(),
			Add:	 randBool,
		}
		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
