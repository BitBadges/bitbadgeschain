package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgTransferTokens(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Get a random existing collection
		collectionId, found := GetRandomCollectionId(r, ctx, k)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferTokens, "no collections exist"), nil, nil
		}
		
		// Check if collection exists
		collection, found := k.GetCollectionFromStore(ctx, collectionId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferTokens, "collection not found"), nil, nil
		}
		
		// Randomly decide between minting and regular transfer
		isMint := r.Intn(3) == 0 // 33% chance of minting
		
		var transfers []*types.Transfer
		var creator simtypes.Account
		
		if isMint {
			// Minting: from Mint address
			creator, _ = simtypes.RandomAcc(r, accs)
			toAddresses := GetRandomAddresses(r, r.Intn(3)+1, accs)
			balances := GetRandomValidBalances(r, 3)
			// Use valid token IDs from collection
			for _, balance := range balances {
				balance.TokenIds = GetRandomValidTokenIds(r, collection, 1)
			}
			transfers = []*types.Transfer{
				{
					From:        types.MintAddress,
					ToAddresses: toAddresses,
					Balances:    balances,
				},
			}
		} else {
			// Regular transfer: need sender with balance
			creator, _ = simtypes.RandomAcc(r, accs)
			fromAddress := creator.Address.String()
			
			// Get sender's balance
			balance, _ := k.GetBalanceOrApplyDefault(ctx, collection, fromAddress)
			
			// Check if sender has any balance
			hasBalance := false
			for _, bal := range balance.Balances {
				if !bal.Amount.IsZero() {
					hasBalance = true
					break
				}
			}
			
			if !hasBalance {
				return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferTokens, "sender has no balance"), nil, nil
			}
			
			// Generate transfers using valid token IDs and amounts from balance
			toAddresses := GetRandomAddresses(r, r.Intn(3)+1, accs)
			transfers = []*types.Transfer{
				{
					From:        fromAddress,
					ToAddresses: toAddresses,
					Balances:    GetRandomValidBalances(r, 2),
				},
			}
			
			// Use valid token IDs from collection
			for _, transfer := range transfers {
				for _, balance := range transfer.Balances {
					balance.TokenIds = GetRandomValidTokenIds(r, collection, 1)
				}
			}
		}
		
		msg := &types.MsgTransferTokens{
			Creator:      creator.Address.String(),
			CollectionId: collectionId,
			Transfers:    transfers,
		}
		
		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}
		
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
