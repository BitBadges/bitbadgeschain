package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func SimulateMsgNewBadge(
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
			SubassetSupplys: randomAccounts,
			SubassetAmountsToCreate: randomAmounts,
			Permissions: r.Uint64(),
			DefaultSubassetSupply: r.Uint64(),
			Standard: r.Uint64(),
			ArbitraryBytes:   []byte(simtypes.RandStringOfLength(r, r.Intn(256))),
			Uri:    &types.UriObject{
				Uri: []byte(simtypes.RandStringOfLength(r, r.Intn(100))),
				Scheme: uint64(r.Intn(10)),
				DecodeScheme: uint64(r.Intn(10)),
				IdxRangeToRemove: &types.IdRange{
					Start: uint64(r.Intn(10)),
					End: uint64(r.Intn(10)),
				},
				InsertSubassetBytesIdx: uint64(r.Intn(10)),
				BytesToInsert: []byte(simtypes.RandStringOfLength(r, r.Intn(100))),
				InsertIdIdx: uint64(r.Intn(10)),
			},
		}

		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
	}
}
