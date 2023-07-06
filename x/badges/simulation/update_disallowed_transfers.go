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

// func SimulateMsgUpdateCollectionApprovedTransfers(
// 	ak types.AccountKeeper,
// 	bk types.BankKeeper,
// 	k keeper.Keeper,
// ) simtypes.Operation {
// 	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
// 	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
// 		simAccount, _ := simtypes.RandomAcc(r, accs)

// 		msg := &types.MsgUpdateCollectionApprovedTransfers{
// 			Creator:      simAccount.Address.String(),
// 			CollectionId: sdkmath.NewUint(r.Uint64()),
// 			ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 				{
// 					From: &types.AddressMapping{
// 						Addresses: []string{
// 							simAccount.Address.String(),
// 							simAccount.Address.String(),
// 							simAccount.Address.String(),
// 						},
// 						OnlySpecifiedAddresses: sdkmath.NewUint(r.Uint64()).Mod(sdkmath.NewUint(2)).IsZero(),
// 						ManagerOptions:       sdkmath.NewUint(r.Uint64()),
// 					},
// 				},
// 			},
// 		}

// 		return simtypes.NewOperationMsg(msg, true, "", types.ModuleCdc), nil, nil
// 	}
// }
