package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// func (k msgServer) ClaimBadge(goCtx context.Context, msg *types.MsgClaimBadge) (*types.MsgClaimBadgeResponse, error) {
// 	ctx := sdk.UnwrapSDKContext(goCtx)
// 	err := *new(error)

// 	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
// 	if !found {
// 		return nil, ErrBadgeNotExists
// 	}

// 	if collection.BalancesType.LTE(sdk.NewUint(0)) {
// 		return nil, ErrOffChainBalances
// 	}

// 	claimId := msg.ClaimId
// 	claim, found := k.GetClaimFromStore(ctx, msg.CollectionId, msg.ClaimId)
// 	if !found {
// 		return nil, ErrClaimNotFound
// 	}

// 	//Check if solutions matches challenges length
// 	if len(msg.Solutions) != len(claim.Challenges) {
// 		return nil, ErrSolutionsLengthInvalid
// 	}

// 	if !claim.IsAssignable && msg.Recipient != "" {
// 		return nil, ErrClaimNotAssignable
// 	}

// 	recipient := msg.Recipient
// 	if recipient == "" {
// 		recipient = msg.Creator
// 	}

// 	//Assert claim is not expired
// 	blockTime := sdk.NewUint(uint64(ctx.BlockTime().UnixMilli()))
// 	validTime := false
// 	for _, interval := range claim.TimeIntervals {
// 		if interval.Start.GT(blockTime) || interval.End.LT(blockTime) {
// 			continue
// 		}
// 		validTime = true
// 	}

// 	if !validTime {
// 		return nil, ErrInvalidTime
// 	}

// 	//Check if address can claim
// 	numUsed := sdk.NewUint(0)
// 	if !claim.NumClaimsPerAddress.IsZero() {
// 		numUsed, err = k.IncrementNumUsedForAddressInStore(ctx, msg.CollectionId, claimId, msg.Creator)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if numUsed.GT(claim.NumClaimsPerAddress) {
// 			return nil, ErrAddressMaxUsesExceeded
// 		}
// 	}

// 	increment := sdk.NewUint(0)
// 	if !claim.IncrementIdsBy.IsZero() {
// 		increment = claim.IncrementIdsBy.Mul(claim.TotalClaimsProcessed)
// 	}
// 	//Check if solutions are valid
// 	for idx, challenge := range claim.Challenges {
// 		root := challenge.Root
// 		useCreatorAddressAsLeaf := challenge.IsWhitelistTree
// 		expectedProofLength := challenge.ExpectedProofLength
// 		solution := msg.Solutions[idx]
// 		challengeId := sdk.NewUint(uint64(idx))

// 		if root != "" {
// 			if len(msg.Solutions[idx].Proof.Aunts) != int(expectedProofLength.Uint64()) {
// 				return nil, ErrProofLengthInvalid
// 			}

// 			if useCreatorAddressAsLeaf {
// 				solution.Proof.Leaf = msg.Creator //overwrites it
// 			}

// 			leafIndex := GetLeafIndex(solution.Proof.Aunts)
// 			if challenge.UseLeafIndexForBadgeIds {
// 				//Get leftmost leaf index for layer === expectedProofLength
// 				leftmostLeafIndex := sdk.NewUint(1)
// 				for i := sdk.NewUint(0); i.LT(expectedProofLength); i = i.Add(sdk.NewUint(1)) {
// 					leftmostLeafIndex = leftmostLeafIndex.Mul(sdk.NewUint(2))
// 				}

// 				increment = leafIndex.Sub(leftmostLeafIndex)
// 			}

// 			if challenge.MaxOneUsePerLeaf {
// 				numUsed, err := k.IncrementNumUsedForChallengeInStore(ctx, msg.CollectionId, claimId, challengeId, leafIndex)
// 				if err != nil {
// 					return nil, err
// 				}

// 				maxUses := sdk.NewUint(1)
// 				if numUsed.GT(maxUses) {
// 					return nil, ErrChallengeMaxUsesExceeded
// 				}
// 			}

// 			//Check if claim is valid
// 			if solution.Proof.Leaf == "" {
// 				return nil, ErrLeafIsEmpty
// 			}

// 			err = CheckMerklePath(solution.Proof.Leaf, root, solution.Proof.Aunts)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	}

// 	/*
// 		Here, we do the following
// 		1. Get msg.Creator (to) balance from store
// 		2. Get claim balance. For compatibility, we set the claim balance as a UserBalanceStore
// 		3. Add the balances to toBalance and subtract the balances from claimBalanceStore
// 		4. Increment IDs if necessary
// 		5. Set everything in store
// 	*/

// 	toBalance, found := k.GetUserBalanceFromStore(ctx, ConstructBalanceKey(recipient, msg.CollectionId))
// 	if !found {
// 		toBalance = types.UserBalanceStore{}
// 	}

// 	claimBalanceStore := types.UserBalanceStore{
// 		Balances:  claim.Balances,
// 		Approvals: []*types.Approval{},
// 	}

// 	//Get the current claim amounts
// 	//We already calculated the necessary increment above
// 	currentClaimAmounts := []*types.Balance{}
// 	for _, balance := range claim.StartingClaimAmounts {
// 		incrementedBadgeIds := []*types.IdRange{}
// 		for _, idRange := range balance.BadgeIds {
// 			incrementedBadgeIds = append(incrementedBadgeIds, &types.IdRange{
// 				Start: idRange.Start.Add(increment),
// 				End:   idRange.End.Add(increment),
// 			})
// 		}

// 		currentClaimAmounts = append(currentClaimAmounts, &types.Balance{
// 			Amount:   balance.Amount,
// 			BadgeIds: incrementedBadgeIds,
// 		})
// 	}

// 	badgeIds := []*types.IdRange{}
// 	for _, balance := range currentClaimAmounts {
// 		badgeIds = append(badgeIds, balance.BadgeIds...)
// 	}

// 	//Check if the claim is allowed (requiresApproval is ignored because it is a mint)
// 	allowed, _ := IsTransferAllowed(ctx, badgeIds, collection, "Mint", msg.Creator, "Mint")
// 	if !allowed {
// 		return nil, ErrMintNotAllowed
// 	}

// 	for _, balance := range currentClaimAmounts {
// 		toBalance.Balances, err = AddBalancesForIdRanges(toBalance.Balances, balance.BadgeIds, balance.Amount)
// 		if err != nil {
// 			return nil, err
// 		}

// 		claimBalanceStore.Balances, err = SubtractBalancesForIdRanges(claimBalanceStore.Balances, balance.BadgeIds, balance.Amount)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	claim.TotalClaimsProcessed = claim.TotalClaimsProcessed.Add(sdk.NewUint(uint64(1)))

// 	claim.Balances = claimBalanceStore.Balances

// 	err = k.SetClaimInStore(ctx, msg.CollectionId, claimId, claim)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = k.SetCollectionInStore(ctx, collection)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = k.SetUserBalanceInStore(ctx, ConstructBalanceKey(recipient, msg.CollectionId), toBalance)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ctx.EventManager().EmitEvent(
// 		sdk.NewEvent(sdk.EventTypeMessage,
// 			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
// 			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
// 		),
// 	)

// 	return &types.MsgClaimBadgeResponse{}, nil
// }
