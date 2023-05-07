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

	//Check if claim is allowed
 	allowed, _ :=	IsTransferAllowed(collection, types.GetPermissions(collection.Permissions), "Mint", msg.Creator, "Mint")
	if !allowed {
		return nil, ErrMintNotAllowed
	}

	claimId := msg.ClaimId
	claim, found := k.GetClaimFromStore(ctx, msg.CollectionId, msg.ClaimId)
	if !found {
		return nil, ErrClaimNotFound
	}

	//Check if solutions matches challenges length
	if len(msg.Solutions) != len(claim.Challenges) {
		return nil, ErrSolutionsLengthInvalid
	}

	//Assert claim is not expired
	if claim.TimeRange.Start > uint64(ctx.BlockTime().UnixMilli()) || claim.TimeRange.End < uint64(ctx.BlockTime().UnixMilli()) {
		return nil, ErrClaimTimeInvalid
	}

	//Check if address can claim
	numUsed := uint64(0)
	if claim.NumClaimsPerAddress > 0 {
		numUsed, err = k.IncrementNumUsedForAddressInStore(ctx, msg.CollectionId, claimId, msg.Creator)
		if err != nil {
			return nil, err
		}

		if numUsed > claim.NumClaimsPerAddress {
			return nil, ErrAddressMaxUsesExceeded
		}
	}


	//Check if solutions are valid
	for idx, challenge := range claim.Challenges {
		root := challenge.Root
		useCreatorAddressAsLeaf := challenge.UseCreatorAddressAsLeaf
		expectedProofLength := challenge.ExpectedProofLength
		solution := msg.Solutions[idx]
		challengeId := uint64(idx)

		if root != "" {
			if len(msg.Solutions[idx].Proof.Aunts) != int(expectedProofLength) {
				return nil, ErrProofLengthInvalid
			}

			if useCreatorAddressAsLeaf {
				solution.Proof.Leaf = msg.Creator //overwrites it
			}

			leafIndex := GetLeafIndex(solution.Proof.Aunts)
			numUsed, err := k.IncrementNumUsedForChallengeInStore(ctx, msg.CollectionId, claimId, challengeId, leafIndex)
			if err != nil {
				return nil, err
			}

			maxUses := uint64(1)
			if numUsed > maxUses {
				return nil, ErrChallengeMaxUsesExceeded
			}

			//Check if claim is valid
			if solution.Proof.Leaf == "" {
				return nil, ErrLeafIsEmpty
			}

			err = CheckMerklePath(solution.Proof.Leaf, root, solution.Proof.Aunts)
			if err != nil {
				return nil, err
			}
		}
	}

	//Check if claim is valid
	incrementIdsBy := claim.IncrementIdsBy


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
		Balances:  claim.UndistributedBalances,
		Approvals: []*types.Approval{},
	}
 
	for _, balance := range claim.CurrentClaimAmounts {
		toBalance.Balances, err = AddBalancesForIdRanges(toBalance.Balances, balance.BadgeIds, balance.Amount)
		if err != nil {
			return nil, err
		}

		claimBalanceStore.Balances, err = SubtractBalancesForIdRanges(claimBalanceStore.Balances, balance.BadgeIds, balance.Amount)
		if err != nil {
			return nil, err
		}

		if incrementIdsBy > 0 {
			for i := 0; i < len(balance.BadgeIds); i++ {
				balance.BadgeIds[i].Start += incrementIdsBy
				balance.BadgeIds[i].End += incrementIdsBy
			}
		}
	}

	claim.UndistributedBalances = claimBalanceStore.Balances

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
