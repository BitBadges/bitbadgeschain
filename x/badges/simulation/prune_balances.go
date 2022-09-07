package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgPruneBalances(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		randomAccounts := []uint64{}
		for i := 0; i < r.Intn(10); i++ {
			randomAccounts = append(randomAccounts, r.Uint64())
		}

		randomAmounts := []uint64{}
		for i := 0; i < r.Intn(10); i++ {
			randomAmounts = append(randomAmounts, r.Uint64())
		}

		msg := &types.MsgPruneBalances{
			Creator:   simAccount.Address.String(),
			Addresses: randomAccounts,
			BadgeIds:  randomAmounts,
		}

		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
