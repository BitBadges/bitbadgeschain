package simulation

// import (
// sdkmath "cosmossdk.io/math"
// 	"math/rand"

// 	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
// 	"github.com/bitbadges/bitbadgeschain/x/badges/types"
// 	"github.com/cosmos/cosmos-sdk/baseapp"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
// )

// func SimulateMsgTransferBadge(
// 	ak types.AccountKeeper,
// 	bk types.BankKeeper,
// 	k keeper.Keeper,
// ) simtypes.Operation {
// 	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
// 	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
// 		simAccount, _ := simtypes.RandomAcc(r, accs)

// 		randomAccounts := []string{}
// 		for i := 0; i < r.Intn(10); i++ {
// 			randomAccounts = append(randomAccounts, simAccount.Address.String())
// 		}

// 		randomAmounts := []sdkmath.Uint{}
// 		for i := 0; i < r.Intn(10); i++ {
// 			randomAmounts = append(randomAmounts, sdk.NewUint(r.Uint64()))
// 		}

// 		msg := &types.MsgTransferBadge{
// 			Creator:      simAccount.Address.String(),
// 			From:         simAccount.Address.String(),
// 			CollectionId: sdk.NewUint(r.Uint64()),
// 			Transfers: []*types.Transfer{
// 				{
// 					ToAddresses: randomAccounts,
// 					Balances: []*types.Balance{
// 						{
// 							Amount: sdk.NewUint(r.Uint64()),
// 							BadgeIds: []*types.IdRange{
// 								{
// 									Start: sdk.NewUint(r.Uint64()),
// 									End:   sdk.NewUint(r.Uint64()),
// 								},
// 								{
// 									Start: sdk.NewUint(r.Uint64()),
// 									End:   sdk.NewUint(r.Uint64()),
// 								},
// 								{
// 									Start: sdk.NewUint(r.Uint64()),
// 									End:   sdk.NewUint(r.Uint64()),
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}

// 		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
// 	}
// }
