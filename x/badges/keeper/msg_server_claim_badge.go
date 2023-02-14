package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	if claim.Type == uint64(types.ClaimType_MerkleTree) {
		if msg.Proof.Leaf == "" {
			return nil, ErrLeafIsEmpty
		}

		hashedMsgLeaf := sha256.Sum256([]byte(msg.Proof.Leaf))
		leafHash := hashedMsgLeaf[:]

		usedKey = hex.EncodeToString(leafHash)

		currHash := leafHash
		for _, aunt := range msg.Proof.Aunts {
			decodedAunt, err := hex.DecodeString(aunt.Aunt)
			if err != nil {
				return nil, ErrDecodingHexString
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
			return nil, ErrRootHashInvalid
		}
	}

	//Assert not already claimed
	if k.StoreHasUsedClaimData(ctx, msg.CollectionId, msg.ClaimId, string(usedKey)) {
		return nil, ErrClaimAlreadyUsed
	}


	amountToClaim := []uint64{claim.AmountPerClaim}
	badgeIds := [][]*types.IdRange{claim.BadgeIds}
	toAddressNum := CreatorAccountNum
	if claim.Type == uint64(types.ClaimType_MerkleTree) {
		badgeIds = [][]*types.IdRange{}
		amountToClaim = []uint64{}

		//Split msg.Proof.Leaf string into 5 parts delimited by a comma
		//The first part is the code word
		//The second part is the address
		//The third part is the amount
		//The fourth part is the starting badge id
		//The fifth part is the ending badge id

		res := strings.Split(msg.Proof.Leaf, "-")
		if len(res) < 5 || (len(res) - 3) % 2 != 0 {
			return nil, ErrClaimDataInvalid
		}

		//convert res[1] to sdk.AccAddress
		if res[1] != "" {
			toAddress, err := sdk.AccAddressFromBech32(res[1])
			if err != nil {
				return nil, ErrInvalidAddress
			}
			toAddressNum = k.Keeper.GetOrCreateAccountNumberForAccAddressBech32(ctx, toAddress)

			if toAddressNum != CreatorAccountNum {
				return nil, ErrMustBeClaimee
			}
		}
		
		
		for i := 2; i < len(res); i += 3 {
			//convert res[i] to uint64
			amount, err := strconv.ParseUint(res[i], 10, 64)
			if err != nil {
				return nil, ErrClaimDataInvalid
			}
			amountToClaim = append(amountToClaim, amount)

			//convert res[i+1] to uint64
			startingBadgeId, err := strconv.ParseUint(res[i+1], 10, 64)
			if err != nil {
				return nil, ErrClaimDataInvalid
			}

			//convert res[i+2] to uint64
			endingBadgeId, err := strconv.ParseUint(res[i+2], 10, 64)
			if err != nil {
				return nil, ErrClaimDataInvalid
			}

			badgeIds = append(badgeIds, []*types.IdRange{&types.IdRange{
				Start: startingBadgeId,
				End: endingBadgeId,
			}})
		}
	}

	

	userBalance, found := k.GetUserBalanceFromStore(ctx, ConstructBalanceKey(toAddressNum, msg.CollectionId))
	if !found {
		userBalance = types.UserBalance{}
	}

	claimUserBalance := types.UserBalance{
		Balances: claim.Balances,
		Approvals: []*types.Approval{},
	}
	
	for i := 0; i < len(amountToClaim); i++ {
		userBalance, err = AddBalancesForIdRanges(userBalance, badgeIds[i], amountToClaim[i])
		if err != nil {
			return nil, err
		}

		claimUserBalance, err = SubtractBalancesForIdRanges(claimUserBalance, badgeIds[i], amountToClaim[i])
		if err != nil {
			return nil, err
		}
	}

	claim.Balances = claimUserBalance.Balances


	if claim.IncrementIdsBy > 0 && claim.Type == uint64(types.ClaimType_FirstCome) {
		for i := 0; i < len(badgeIds[0]); i++ {
			badgeIds[0][i].Start, err = SafeAdd(badgeIds[0][i].Start, claim.IncrementIdsBy)
			if err != nil {
				return nil, err
			}

			badgeIds[0][i].End, err = SafeAdd(badgeIds[0][i].End, claim.IncrementIdsBy)
			if err != nil {
				return nil, err
			}
		}
		claim.BadgeIds = badgeIds[0]
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

	collectionJson, err := json.Marshal(collection)
	if err != nil {
		return nil, err
	}

	userBalanceJson, err := json.Marshal(userBalance)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("collection", string(collectionJson)),
			sdk.NewAttribute("user_balance", string(userBalanceJson)),
			sdk.NewAttribute("to", fmt.Sprint(toAddressNum)),
			sdk.NewAttribute("claim_data", string(usedKey)),
		),
	)

	
	return &types.MsgClaimBadgeResponse{}, nil
}
