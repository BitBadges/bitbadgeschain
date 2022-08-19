package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
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
			Creator: simAccount.Address.String(),
			BadgeId: r.Uint64(),
			Uri: &types.UriObject{
				Uri:          []byte(simtypes.RandStringOfLength(r, r.Intn(100))),
				Scheme:       uint64(r.Intn(10)),
				DecodeScheme: uint64(r.Intn(10)),
				IdxRangeToRemove: &types.IdRange{
					Start: uint64(r.Intn(10)),
					End:   uint64(r.Intn(10)),
				},
				InsertSubassetBytesIdx: uint64(r.Intn(10)),
				BytesToInsert:          []byte(simtypes.RandStringOfLength(r, r.Intn(100))),
				InsertIdIdx:            uint64(r.Intn(10)),
			},
		}

		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
