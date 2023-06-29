package simulation

// import (
// 	"math/rand"

// 	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
// 	"github.com/bitbadges/bitbadgeschain/x/badges/types"
// 	"github.com/cosmos/cosmos-sdk/baseapp"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
// )

// func SimulateMsgSetApproval(
// 	ak types.AccountKeeper,
// 	bk types.BankKeeper,
// 	k keeper.Keeper,
// ) simtypes.Operation {
// 	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
// 	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
// 		simAccount, _ := simtypes.RandomAcc(r, accs)
// 		msg := &types.MsgSetApproval{
// 			Creator: simAccount.Address.String(),
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdk.NewUint(r.Uint64()),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdk.NewUint(r.Uint64()),
// 							End:   sdk.NewUint(r.Uint64()),
// 						},
// 						{
// 							Start: sdk.NewUint(r.Uint64()),
// 							End:   sdk.NewUint(r.Uint64()),
// 						},
// 						{
// 							Start: sdk.NewUint(r.Uint64()),
// 							End:   sdk.NewUint(r.Uint64()),
// 						},
// 					},
// 				},
// 			},
// 			CollectionId: sdk.NewUint(r.Uint64()),
// 		}

// 		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
// 	}
// }
