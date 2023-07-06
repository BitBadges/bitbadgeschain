package simulation

import (
	"math/rand"

	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgMintAndDistributeBadges(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		randomSubassets := []*types.Balance{}
		for i := 0; i < r.Intn(10); i++ {
			start := sdkmath.NewUint(r.Uint64())
			randomSubassets = append(randomSubassets, &types.Balance{
				Amount: sdkmath.NewUint(r.Uint64()),
				BadgeIds: []*types.IdRange{
					{
						Start: start,
						End:   start.Add(sdkmath.NewUint(r.Uint64())),
					},
				},
			})
		}

		msg := &types.MsgMintAndDistributeBadges{
			Creator:        simAccount.Address.String(),
			CollectionId:   sdkmath.NewUint(r.Uint64()),
			BadgesToCreate: randomSubassets,
		}

		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
