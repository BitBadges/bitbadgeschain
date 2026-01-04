package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgCreateDynamicStore(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Ensure we have valid accounts
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateDynamicStore, "no accounts available"), nil, nil
		}
		
		simAccount := EnsureAccountExists(r, accs)

		// Random boolean for defaultValue
		defaultValue := r.Intn(2) == 0

		// Random string for uri (sometimes empty, sometimes with value)
		uri := ""
		if r.Intn(2) == 0 {
			uri = simtypes.RandStringOfLength(r, 20)
		}

		// Random string for customData (sometimes empty, sometimes with value)
		customData := ""
		if r.Intn(2) == 0 {
			customData = simtypes.RandStringOfLength(r, 30)
		}

		msg := &types.MsgCreateDynamicStore{
			Creator:      simAccount.Address.String(),
			DefaultValue: defaultValue,
			Uri:          uri,
			CustomData:   customData,
		}

		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

