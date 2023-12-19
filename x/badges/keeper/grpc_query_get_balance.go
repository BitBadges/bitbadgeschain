package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Queries a balance for the given address and badgeId and returns its contents.
func (k Keeper) GetBalance(goCtx context.Context, req *types.QueryGetBalanceRequest) (*types.QueryGetBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	//Assert that initiatedBy owns the required badges
	balances := &types.UserBalanceStore{}
	isBlank := false

	currCollectionId := req.CollectionId
	currCollection := &types.BadgeCollection{}


	//Check if the collection has inherited balances
	collection, found := k.GetCollectionFromStore(ctx, currCollectionId)
	if !found {
		isBlank = true
	} else {
		currCollection = collection
		isStandardBalances := collection.BalancesType == "Standard"
		if isStandardBalances  || req.Address == "Mint"  || req.Address == "Total" {
			initiatedByBalanceKey := ConstructBalanceKey(req.Address, currCollectionId)
			initiatedByBalance, found := k.GetUserBalanceFromStore(ctx, initiatedByBalanceKey)
			if found {
				balances = initiatedByBalance
			} else {
				isBlank = true
			}
		} else {
			return nil, sdkerrors.Wrapf(ErrWrongBalancesType, "unsupported balances type %s %s", collection.BalancesType, collection.CollectionId)
		}
	}

	if !isBlank {
		return &types.QueryGetBalanceResponse{
			Balance: balances,
		}, nil
	} else {
		blankUserBalance := &types.UserBalanceStore{
			Balances:          []*types.Balance{},
			OutgoingApprovals: currCollection.DefaultUserOutgoingApprovals,
			IncomingApprovals: currCollection.DefaultUserIncomingApprovals,
			AutoApproveSelfInitiatedOutgoingTransfers: currCollection.DefaultAutoApproveSelfInitiatedOutgoingTransfers,
			AutoApproveSelfInitiatedIncomingTransfers: currCollection.DefaultAutoApproveSelfInitiatedIncomingTransfers,
			UserPermissions: currCollection.DefaultUserPermissions,
		}
		return &types.QueryGetBalanceResponse{
			Balance: blankUserBalance,
		}, nil
	}
}
