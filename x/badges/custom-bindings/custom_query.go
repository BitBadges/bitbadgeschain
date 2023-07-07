package custom_bindings

import (
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	badgeKeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	badgeTypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PerformCustomBadgeQuery(keeper badgeKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom badgeCustomQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		switch {
		case custom.QueryCollection != nil:
			res, err := keeper.GetCollection(ctx, custom.QueryCollection)
			if err != nil {
				return nil, err
			}
			return json.Marshal(badgeTypes.QueryGetCollectionResponse{Collection: res.Collection})
		case custom.QueryBalance != nil:
			res, err := keeper.GetBalance(ctx, custom.QueryBalance)
			if err != nil {
				return nil, err
			}
			return json.Marshal(badgeTypes.QueryGetBalanceResponse{Balance: res.Balance})

		case custom.QueryAddressMapping != nil:
			res, err := keeper.GetAddressMapping(ctx, custom.QueryAddressMapping)
			if err != nil {
				return nil, err
			}
			return json.Marshal(badgeTypes.QueryGetAddressMappingResponse{Mapping: res.Mapping})
		case custom.QueryApprovalsTracker != nil:
			res, err := keeper.GetApprovalsTracker(ctx, custom.QueryApprovalsTracker)
			if err != nil {
				return nil, err
			}
			return json.Marshal(badgeTypes.QueryGetApprovalsTrackerResponse{Tracker: res.Tracker})
		case custom.QueryGetNumUsedForChallenge != nil:
			res, err := keeper.GetNumUsedForChallenge(ctx, custom.QueryGetNumUsedForChallenge)
			if err != nil {
				return nil, err
			}
			return json.Marshal(badgeTypes.QueryGetNumUsedForChallengeResponse{NumUsed: res.NumUsed})
		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

type badgeCustomQuery struct {
	QueryCollection             *badgeTypes.QueryGetCollectionRequest          `json:"queryCollection,omitempty"`
	QueryBalance                *badgeTypes.QueryGetBalanceRequest             `json:"queryBalance,omitempty"`
	QueryAddressMapping         *badgeTypes.QueryGetAddressMappingRequest      `json:"queryAddressMapping,omitempty"`
	QueryApprovalsTracker       *badgeTypes.QueryGetApprovalsTrackerRequest    `json:"queryApprovalsTracker,omitempty"`
	QueryGetNumUsedForChallenge *badgeTypes.QueryGetNumUsedForChallengeRequest `json:"queryGetNumUsedForChallenge,omitempty"`
}
