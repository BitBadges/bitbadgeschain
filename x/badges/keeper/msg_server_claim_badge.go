package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	if claim.TimeRange.Start > uint64(ctx.BlockTime().Unix()) || claim.TimeRange.End < uint64(ctx.BlockTime().Unix()) {
		return nil, ErrClaimTimeInvalid
	}

	//Check if claim is valid
	badgeIds := claim.BadgeIds
	origBadgeIds := []*types.IdRange{}
	for _, badgeId := range badgeIds {
		origBadgeIds = append(origBadgeIds, &types.IdRange{
			Start: badgeId.Start,
			End:   badgeId.End,
		})
	}

	codeRoot := claim.CodeRoot
	whitelistRoot := claim.WhitelistRoot
	amountToClaim := claim.Amount
	incrementIdsBy := claim.IncrementIdsBy

	if codeRoot != "" {
		if len(msg.CodeProof.Aunts) != int(claim.ExpectedMerkleProofLength) {
			return nil, ErrCodeProofLengthInvalid
		}

		codeLeafIndex := uint64(1)
		//iterate through msg.CodeProof.Aunts backwards
		for i := len(msg.CodeProof.Aunts) - 1; i >= 0; i-- {
			aunt := msg.CodeProof.Aunts[i]
			onRight := aunt.OnRight

			if onRight {
				codeLeafIndex = codeLeafIndex * 2
			} else {
				codeLeafIndex = codeLeafIndex*2 + 1
			}
		}

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

		hashedMsgLeaf := sha256.Sum256([]byte(msg.CodeProof.Leaf))
		leafHash := hashedMsgLeaf[:]

		str := ""
		str = str + msg.CodeProof.Leaf + " - " + hex.EncodeToString(leafHash)

		currHash := leafHash
		for _, aunt := range msg.CodeProof.Aunts {
			decodedAunt, err := hex.DecodeString(aunt.Aunt)
			if err != nil {
				return nil, ErrDecodingHexString
			}

			if aunt.OnRight {
				parentHash := sha256.Sum256(append(currHash, decodedAunt...))
				currHash = parentHash[:]
			} else {
				parentHash := sha256.Sum256(append(decodedAunt, currHash...))
				currHash = parentHash[:]
			}

			str = str + " - " + hex.EncodeToString(currHash)
		}

		hexCurrHash := hex.EncodeToString(currHash)

		str = str + " - ROOT: " + codeRoot

		if hexCurrHash != codeRoot {
			return nil, sdkerrors.Wrapf(ErrRootHashInvalid, "Got: %s", str)
		}
	}

	numUsed := uint64(0)
	if claim.RestrictOptions == 2 { //by address
		numUsed, err = k.IncrementNumUsedForAddressInStore(ctx, msg.CollectionId, claimId, msg.Creator)
		if err != nil {
			return nil, err
		}

		if numUsed > 1 {
			return nil, ErrAddressAlreadyUsed
		}
	}

	if whitelistRoot != "" {
		if msg.WhitelistProof.Leaf != msg.Creator {
			return nil, ErrMustBeClaimee
		}

		if claim.RestrictOptions == 1 { //whitelist index
			whitelistLeafIndex := uint64(1)
			//iterate through msg.WhitelistProof.Aunts backwards
			for i := len(msg.WhitelistProof.Aunts) - 1; i >= 0; i-- {
				aunt := msg.WhitelistProof.Aunts[i]
				onRight := aunt.OnRight

				if onRight {
					whitelistLeafIndex = whitelistLeafIndex * 2
				} else {
					whitelistLeafIndex = whitelistLeafIndex*2 + 1
				}
			}
			numUsed, err = k.IncrementNumUsedForWhitelistIndexInStore(ctx, msg.CollectionId, claimId, whitelistLeafIndex)
			if err != nil {
				return nil, err
			}
		}

		maxUses := uint64(1)
		if numUsed > maxUses {
			return nil, ErrCodeMaxUsesExceeded
		}

		//Check if claim is valid
		if msg.WhitelistProof.Leaf == "" {
			return nil, ErrLeafIsEmpty
		}

		hashedMsgLeaf := sha256.Sum256([]byte(msg.WhitelistProof.Leaf))
		currHash := hashedMsgLeaf[:]

		for _, aunt := range msg.WhitelistProof.Aunts {
			decodedAunt, err := hex.DecodeString(aunt.Aunt)
			if err != nil {
				return nil, ErrDecodingHexString
			}

			if aunt.OnRight {
				parentHash := sha256.Sum256(append(currHash, decodedAunt...))
				currHash = parentHash[:]
			} else {
				parentHash := sha256.Sum256(append(decodedAunt, currHash...))
				currHash = parentHash[:]
			}
		}

		hexCurrHash := hex.EncodeToString(currHash)
		if hexCurrHash != whitelistRoot {
			return nil, ErrRootHashInvalid
		}
	}

	userBalance, found := k.GetUserBalanceFromStore(ctx, ConstructBalanceKey(msg.Creator, msg.CollectionId))
	if !found {
		userBalance = types.UserBalanceStore{}
	}

	claimUserBalance := types.UserBalanceStore{
		Balances:  claim.Balances,
		Approvals: []*types.Approval{},
	}

	userBalance, err = AddBalancesForIdRanges(userBalance, badgeIds, amountToClaim)
	if err != nil {
		return nil, err
	}

	claimUserBalance, err = SubtractBalancesForIdRanges(claimUserBalance, badgeIds, amountToClaim)
	if err != nil {
		return nil, err
	}

	claim.Balances = claimUserBalance.Balances

	if incrementIdsBy > 0 {
		_, err := k.IncrementNumUsedForClaimInStore(ctx, msg.CollectionId, claimId)
		if err != nil {
			return nil, err
		}

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

	err = k.SetUserBalanceInStore(ctx, ConstructBalanceKey(msg.Creator, msg.CollectionId), userBalance)
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
