package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ClaimBadge(goCtx context.Context, msg *types.MsgClaimBadge) (*types.MsgClaimBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := *new(error)
	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !found {
		return nil, ErrBadgeNotExists
	}

	if uint64(len(collection.Claims)) <= msg.ClaimId {
		return nil, ErrClaimNotFound
	}

	claim := collection.Claims[msg.ClaimId]

	//Assert claim is not expired
	if claim.TimeRange.Start > uint64(ctx.BlockTime().Unix()) || claim.TimeRange.End < uint64(ctx.BlockTime().Unix()) {
		return nil, ErrClaimTimeInvalid
	}

	usedKey := msg.Creator
	if claim.Type == uint64(types.ClaimType_AccountNum) {
		hashedMsgLeaf := sha256.Sum256([]byte(msg.Proof.Leaf))
		leafHash := hashedMsgLeaf[:]

		usedKey = string(leafHash)

		currHash := leafHash
		for _, aunt := range msg.Proof.Aunts {
			decodedAunt, err := hex.DecodeString(aunt.Aunt)
			if err != nil {
				return nil, ErrRootHashInvalid
				// return nil, sdkerrors.Wrapf(err, "Proof is invalid %s %d", hex.EncodeToString(currHash), len(aunt.Aunt))
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
		if hexCurrHash != claim.Data {
			// return nil, sdkerrors.Wrapf(ErrRootHashInvalid, "Proof is invalid %s %s", hexCurrHash, claim.Data)
			return nil, ErrRootHashInvalid
		}
	}

	//Assert not already claimed
	if k.StoreHasUsedClaimData(ctx, msg.CollectionId, msg.ClaimId, string(usedKey)) {
		return nil, ErrClaimAlreadyUsed
	}


	
	amountToClaim := claim.AmountPerClaim
	badgeIds := claim.BadgeIds
	toAddressNum := CreatorAccountNum
	if claim.Type == uint64(types.ClaimType_AccountNum) {
		//Split msg.Proof.Leaf string into 5 parts delimited by a comma
		//The first part is the code word
		//The second part is the address
		//The third part is the amount
		//The fourth part is the starting badge id
		//The fifth part is the ending badge id

		res := strings.Split(msg.Proof.Leaf, "-")
		if len(res) != 5 {
			return nil, ErrClaimDataInvalid
		}

		//convert res[1] to sdk.AccAddress
		if res[1] != "" {
			toAddress, err := sdk.AccAddressFromBech32(res[1])
			if err != nil {
				return nil, ErrInvalidAddress
			}
			toAddressNum = k.Keeper.GetOrCreateAccountNumberForAccAddressBech32(ctx, toAddress)
		}
		
		
		
		//convert res[2] to uint64
		amountToClaim, err = strconv.ParseUint(res[2], 10, 64)
		if err != nil {
			return nil, ErrClaimDataInvalid
		}

		//convert res[3] to uint64
		startingBadgeId, err := strconv.ParseUint(res[3], 10, 64)
		if err != nil {
			return nil, ErrClaimDataInvalid
		}

		//convert res[4] to uint64
		endingBadgeId, err := strconv.ParseUint(res[4], 10, 64)
		if err != nil {
			return nil, ErrClaimDataInvalid
		}

		badgeIds = &types.IdRange{
			Start: startingBadgeId,
			End:   endingBadgeId,
		}
		
	}




	

	userBalance, found := k.GetUserBalanceFromStore(ctx, ConstructBalanceKey(toAddressNum, msg.CollectionId))
	if !found {
		userBalance = types.UserBalance{}
	}
	
	userBalance, err = AddBalancesForIdRanges(userBalance, []*types.IdRange{badgeIds}, amountToClaim)
	if err != nil {
		return nil, err
	}

	claimUserBalance := types.UserBalance{
		Balances: []*types.Balance{
			claim.Balance,
		},
		Approvals: []*types.Approval{},
	}

	
	claimUserBalance, err = SubtractBalancesForIdRanges(claimUserBalance, []*types.IdRange{badgeIds}, amountToClaim)
	if err != nil {
		return nil, err
	}

	claim.Balance = claimUserBalance.Balances[0]


	if claim.IncrementIdsBy > 0 {
		badgeIds.Start, err = SafeAdd(badgeIds.Start, claim.IncrementIdsBy)
		if err != nil {
			return nil, err
		}

		badgeIds.End, err = SafeAdd(badgeIds.End, claim.IncrementIdsBy)
		if err != nil {
			return nil, err
		}
		claim.BadgeIds = badgeIds
	}

	collection.Claims[msg.ClaimId] = claim
	
	err = k.SetCollectionInStore(ctx, collection)
	if err != nil {
		return nil, err
	}

	err = k.SetUserBalanceInStore(ctx, ConstructBalanceKey(toAddressNum, msg.CollectionId), userBalance)
	if err != nil {
		return nil, err
	}

	//Claim and mark as claimed
	err = k.SetUsedClaimDataInStore(ctx, msg.CollectionId, msg.ClaimId, string(usedKey))
	if err != nil {
		return nil, err
	}
	
	return &types.MsgClaimBadgeResponse{}, nil
}
