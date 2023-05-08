package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgUpdateUris(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		msg := &types.MsgUpdateUris{
			Creator:       simAccount.Address.String(),
			CollectionId:  sdk.NewUint(r.Uint64()),
			CollectionUri: simtypes.RandStringOfLength(r, r.Intn(100)),
			BadgeUris:     []*types.BadgeUri{},
		}

		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
