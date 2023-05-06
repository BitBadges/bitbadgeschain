package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ClaimBadge(goCtx context.Context, msg *types.MsgClaimBadge) (*types.MsgClaimBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := *new(error)

	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !found {
		return nil, ErrBadgeNotExists
	}
	
	claimId := msg.ClaimId
	claim, found := k.GetClaimFromStore(ctx, msg.CollectionId, msg.ClaimId)
	if !found {
		return nil, ErrClaimNotFound
	}

	//Assert claim is not expired
	if claim.TimeRange.Start > uint64(ctx.BlockTime().UnixMilli()) || claim.TimeRange.End < uint64(ctx.BlockTime().UnixMilli()) {
		return nil, ErrClaimTimeInvalid
	}

	//Check if address can claim
	numUsed := uint64(0)
	if claim.LimitOnePerAddress {
		numUsed, err = k.IncrementNumUsedForAddressInStore(ctx, msg.CollectionId, claimId, msg.Creator)
		if err != nil {
			return nil, err
		}

		if numUsed > 1 {
			return nil, ErrAddressAlreadyUsed
		}
	}

	//Check if claim is valid
	badgeIds := claim.BadgeIds
	codeRoot := claim.CodeRoot
	whitelistRoot := claim.WhitelistRoot
	amountToClaim := claim.Amount
	incrementIdsBy := claim.IncrementIdsBy

	if codeRoot != "" {
		if len(msg.CodeProof.Aunts) != int(claim.ExpectedCodeProofLength) {
			return nil, ErrCodeProofLengthInvalid
		}

		codeLeafIndex := GetLeafIndex(msg.CodeProof.Aunts)
		numUsed, err := k.IncrementNumUsedForCodeInStore(ctx, msg.CollectionId, claimId, codeLeafIndex)
		if err != nil {
			return nil, err
		}

		maxUses := uint64(1)
		if numUsed > maxUses {
			return nil, ErrCodeMaxUsesExceeded
		}

		//Check if claim is valid
		if msg.CodeProof.Leaf == "" {
			return nil, ErrLeafIsEmpty
		}

		err = CheckMerklePath(msg.CodeProof.Leaf, codeRoot, msg.CodeProof.Aunts)
		if err != nil {
			return nil, err
		}
	}

	if whitelistRoot != "" {
		if msg.WhitelistProof.Leaf != msg.Creator {
			return nil, ErrMustBeClaimee
		}

		whitelistLeafIndex := GetLeafIndex(msg.WhitelistProof.Aunts)
		numUsed, err = k.IncrementNumUsedForWhitelistIndexInStore(ctx, msg.CollectionId, claimId, whitelistLeafIndex)
		if err != nil {
			return nil, err
		}

		maxUses := uint64(1)
		if numUsed > maxUses {
			return nil, ErrCodeMaxUsesExceeded
		}

		//Check if claim is valid
		if msg.WhitelistProof.Leaf == "" {
			return nil, ErrLeafIsEmpty
		}

		err = CheckMerklePath(msg.WhitelistProof.Leaf, whitelistRoot, msg.WhitelistProof.Aunts)
		if err != nil {
			return nil, err
		}
	}

	/*
		Here, we do the following
		1. Get msg.Creator (to) balance from store
		2. Get claim balance. For compatibility, we set the claim balance as a UserBalanceStore
		3. Add the balances to toBalance and subtract the balances from claimBalanceStore
		4. Increment IDs if necessary
		5. Set everything in store
	*/

	toBalance, found := k.GetUserBalanceFromStore(ctx, ConstructBalanceKey(msg.Creator, msg.CollectionId))
	if !found {
		toBalance = types.UserBalanceStore{}
	}

	claimBalanceStore := types.UserBalanceStore{
		Balances:  claim.Balances,
		Approvals: []*types.Approval{},
	}

	toBalance.Balances, err = AddBalancesForIdRanges(toBalance.Balances, badgeIds, amountToClaim)
	if err != nil {
		return nil, err
	}

	claimBalanceStore.Balances, err = SubtractBalancesForIdRanges(claimBalanceStore.Balances, badgeIds, amountToClaim)
	if err != nil {
		return nil, err
	}

	claim.Balances = claimBalanceStore.Balances
	if incrementIdsBy > 0 {
		for i := 0; i < len(badgeIds); i++ {
			badgeIds[i].Start += incrementIdsBy
			badgeIds[i].End += incrementIdsBy
		}
	}

	err = k.SetClaimInStore(ctx, msg.CollectionId, claimId, claim)
	if err != nil {
		return nil, err
	}


	err = k.SetCollectionInStore(ctx, collection)
	if err != nil {
		return nil, err
	}

	err = k.SetUserBalanceInStore(ctx, ConstructBalanceKey(msg.Creator, msg.CollectionId), toBalance)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgClaimBadgeResponse{}, nil
}
