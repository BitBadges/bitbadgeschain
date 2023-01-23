package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgMintBadge(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		randomSubassets := []*types.BadgeSupplyAndAmount{}
		for i := 0; i < r.Intn(10); i++ {
			randomSubassets = append(randomSubassets, &types.BadgeSupplyAndAmount{
				Supply: r.Uint64(),
				Amount: r.Uint64(),
			});		
		}

		msg := &types.MsgMintBadge{
			Creator:         simAccount.Address.String(),
			CollectionId:         r.Uint64(),
			BadgeSupplys: randomSubassets,
		}

		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
