package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgCreateAddressLists(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Ensure we have valid accounts
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgCreateAddressLists, "no accounts available"), nil, nil
		}

		simAccount := EnsureAccountExists(r, accs)

		// Generate 1-3 address lists
		count := r.Intn(3) + 1
		addressLists := []*types.AddressListInput{}
		usedListIds := make(map[string]bool)

		for i := 0; i < count; i++ {
			// Generate unique list ID
			listId := simtypes.RandStringOfLength(r, 10)
			for usedListIds[listId] {
				listId = simtypes.RandStringOfLength(r, 10)
			}
			usedListIds[listId] = true

			// Generate addresses (1-5 addresses per list)
			addressCount := r.Intn(5) + 1
			addresses := GetRandomAddresses(r, addressCount, accs)

			// Ensure no duplicate addresses in the list
			uniqueAddresses := make(map[string]bool)
			uniqueList := []string{}
			for _, addr := range addresses {
				if !uniqueAddresses[addr] {
					uniqueAddresses[addr] = true
					uniqueList = append(uniqueList, addr)
				}
			}

			addressLists = append(addressLists, &types.AddressListInput{
				Addresses: uniqueList,
				ListId:    listId,
			})
		}

		msg := &types.MsgCreateAddressLists{
			Creator:      simAccount.Address.String(),
			AddressLists: addressLists,
		}

		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
