package keeper

import (
	"bytes"
	"context"
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/merkle"
)

func (k msgServer) ClaimBadge(goCtx context.Context, msg *types.MsgClaimBadge) (*types.MsgClaimBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	claim, found := k.GetClaimFromStore(ctx, msg.ClaimId)
	if !found {
		return nil, ErrClaimNotFound
	}

	claimData := msg.Leaf
	if claim.Type == types.ClaimType_AccountNum {
		//Assert claimData is either the account number or cosmos address

		accountNumBytes := []byte(strconv.FormatUint(CreatorAccountNum, 10))
		addressBytes := []byte(msg.Creator)

		if !bytes.Equal(claimData, accountNumBytes) && !bytes.Equal(claimData, addressBytes) {
			return nil, ErrClaimDataInvalid
		}
	}

	//Assert not already claimed
	if k.StoreHasUsedClaimData(ctx, msg.CollectionId, msg.ClaimId, string(claimData)) {
		return nil, ErrClaimAlreadyUsed
	}

	//Convert claim.Data to bytes
	rootHash := claim.Data
	
	proof := merkle.Proof{
		Total:    msg.Proof.Total,
		Index:    msg.Proof.Index,
		LeafHash: msg.Proof.LeafHash,
		Aunts:    msg.Proof.Aunts,
	}

	err := proof.Verify(rootHash, claimData);
	if err != nil {
		return nil, ErrClaimDataInvalid
	}

	userBalance, found := k.GetUserBalanceFromStore(ctx, ConstructBalanceKey(CreatorAccountNum, msg.CollectionId))
	if !found {
		userBalance = types.UserBalance{}
	}
	
	userBalance, err = AddBalancesForIdRanges(userBalance, claim.Balance.BadgeIds, claim.AmountPerClaim)
	if err != nil {
		return nil, err
	}

	claimUserBalance := types.UserBalance{
		Balances: []*types.Balance{
			claim.Balance,
		},
		Approvals: []*types.Approval{},
	}

	
	claimUserBalance, err = SubtractBalancesForIdRanges(claimUserBalance, claim.Balance.BadgeIds, claim.AmountPerClaim)
	if err != nil {
		return nil, err
	}

	claim.Balance = claimUserBalance.Balances[0]

	err = k.SetClaimInStore(ctx, msg.ClaimId, claim)
	if err != nil {
		return nil, err
	}

	err = k.SetUserBalanceInStore(ctx, ConstructBalanceKey(CreatorAccountNum, msg.CollectionId), GetBalanceToInsertToStorage(userBalance))
	if err != nil {
		return nil, err
	}

	//Claim and mark as claimed
	err = k.SetUsedClaimDataInStore(ctx, msg.CollectionId, msg.ClaimId, string(claimData))
	if err != nil {
		return nil, err
	}
	
	return &types.MsgClaimBadgeResponse{}, nil
}