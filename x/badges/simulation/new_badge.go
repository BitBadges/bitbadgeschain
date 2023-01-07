package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgNewBadge(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		randomSubassets := []*types.SubassetSupplyAndAmount{}
		for i := 0; i < r.Intn(10); i++ {
			randomSubassets = append(randomSubassets, &types.SubassetSupplyAndAmount{
				Supply: r.Uint64(),
				Amount: r.Uint64(),
			});		
		}


		msg := &types.MsgNewBadge{
			Creator: simAccount.Address.String(),
			FreezeAddressRanges: []*types.IdRange{
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
			SubassetSupplysAndAmounts: randomSubassets,
			Permissions:             r.Uint64(),
			DefaultSubassetSupply:   r.Uint64(),
			Standard:                r.Uint64(),
			ArbitraryBytes:          simtypes.RandStringOfLength(r, r.Intn(256)),
			Uri: &types.UriObject{
				Uri:          simtypes.RandStringOfLength(r, r.Intn(100)),
				Scheme:       uint64(r.Intn(10)),
				DecodeScheme: uint64(r.Intn(10)),
				IdxRangeToRemove: &types.IdRange{
					Start: uint64(r.Intn(10)),
					End:   uint64(r.Intn(10)),
				},
				InsertSubassetBytesIdx: uint64(r.Intn(10)),
				BytesToInsert:          simtypes.RandStringOfLength(r, r.Intn(100)),
				InsertIdIdx:            uint64(r.Intn(10)),
			},
		}

		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
