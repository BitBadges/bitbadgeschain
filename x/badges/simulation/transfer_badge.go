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
		// Ensure we have valid accounts
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferTokens, "no accounts available"), nil, nil
		}
		
		// Try to get a known-good collection ID first
		collectionId, found := GetKnownGoodCollectionId(ctx, k)
		if !found {
			// Fallback: try to get a random existing collection
			collectionId, found = GetRandomCollectionId(r, ctx, k)
			if !found {
				// Try to create one first
				simAccount := EnsureAccountExists(r, accs)
				createdId, err := GetOrCreateCollection(ctx, k, simAccount.Address.String(), r, accs)
				if err != nil {
					return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgTransferTokens, "no collections exist and failed to create one"), nil, nil
				}
				collectionId = createdId
			}
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
			creator = EnsureAccountExists(r, accs)
			toAddresses := GetRandomAddresses(r, r.Intn(3)+1, accs)
			// Use valid balances from collection
			balances := GetValidBalancesFromCollection(r, collection, r.Intn(3)+1)
			transfers = []*types.Transfer{
				{
					From:        types.MintAddress,
					ToAddresses: toAddresses,
					Balances:    balances,
				},
			}
		} else {
			// Regular transfer: need sender with balance
			creator = EnsureAccountExists(r, accs)
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
				// Try minting instead if sender has no balance
				toAddresses := GetRandomAddresses(r, r.Intn(3)+1, accs)
				balances := GetValidBalancesFromCollection(r, collection, r.Intn(3)+1)
				transfers = []*types.Transfer{
					{
						From:        types.MintAddress,
						ToAddresses: toAddresses,
						Balances:    balances,
					},
				}
			} else {
				// Generate transfers using valid token IDs and amounts from balance
				toAddresses := GetRandomAddresses(r, r.Intn(3)+1, accs)
				// Use valid balances from collection
				balances := GetValidBalancesFromCollection(r, collection, r.Intn(2)+1)
				transfers = []*types.Transfer{
					{
						From:        fromAddress,
						ToAddresses: toAddresses,
						Balances:    balances,
					},
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
