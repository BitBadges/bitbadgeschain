package simulation

import (
	"math/rand"

	"bitbadgeschain/x/offers/keeper"
	"bitbadgeschain/x/offers/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgRejectAndDeleteProposal(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgRejectAndDeleteProposal{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the RejectAndDeleteProposal simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "RejectAndDeleteProposal simulation not implemented"), nil, nil
	}
}
