package custom_bindings

import (
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	anchorKeeper "github.com/bitbadges/bitbadgeschain/x/anchor/keeper"
	anchortypes "github.com/bitbadges/bitbadgeschain/x/anchor/types"
	badgeKeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	badgeTypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	mapsKeeper "github.com/bitbadges/bitbadgeschain/x/maps/keeper"
	mapsTypes "github.com/bitbadges/bitbadgeschain/x/maps/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PerformCustomBitBadgesModuleQuery(bk badgeKeeper.Keeper, ak anchorKeeper.Keeper, mk mapsKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		isBadgeModuleQuery := false
		var custom badgeCustomQuery
		err := json.Unmarshal(request, &custom)
		if err == nil {
			isBadgeModuleQuery = true
		}

		isAnchorModuleQuery := false
		var anchorCustom anchorCustomQuery
		err = json.Unmarshal(request, &anchorCustom)
		if err == nil {
			isAnchorModuleQuery = true
		}

		isMapsModuleQuery := false
		var mapsCustom mapsCustomQuery
		err = json.Unmarshal(request, &mapsCustom)
		if err == nil {
			isMapsModuleQuery = true
		}

		if isBadgeModuleQuery {
			return PerformCustomBadgeQuery(bk)(ctx, request)
		} else if isAnchorModuleQuery {
			return PerformCustomAnchorQuery(ak)(ctx, request)
		} else if isMapsModuleQuery {
			return PerformCustomMapsQuery(mk)(ctx, request)
		}

		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

func PerformCustomMapsQuery(mk mapsKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom mapsCustomQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		switch {
		case custom.QueryMap != nil:
			res, err := mk.Map(ctx, custom.QueryMap)
			if err != nil {
				return nil, err
			}
			return json.Marshal(mapsTypes.QueryGetMapResponse{Map: res.Map})
		case custom.QueryMapValue != nil:
			res, err := mk.MapValue(ctx, custom.QueryMapValue)
			if err != nil {
				return nil, err
			}
			return json.Marshal(mapsTypes.QueryGetMapValueResponse{Value: res.Value})
		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

func PerformCustomAnchorQuery(ak anchorKeeper.Keeper) wasmKeeper.CustomQuerier {
	var custom anchorCustomQuery
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		switch {
		case custom.QueryValueAtLocation != nil:
			res, err := ak.GetValueAtLocation(ctx, custom.QueryValueAtLocation)
			if err != nil {
				return nil, err
			}
			return json.Marshal(anchortypes.QueryGetValueAtLocationResponse{AnchorData: res.AnchorData})
		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

// WASM handler for contracts calling into the badges module
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

		case custom.QueryAddressList != nil:
			res, err := keeper.GetAddressList(ctx, custom.QueryAddressList)
			if err != nil {
				return nil, err
			}
			return json.Marshal(badgeTypes.QueryGetAddressListResponse{List: res.List})
		case custom.QueryApprovalTracker != nil:
			res, err := keeper.GetApprovalTracker(ctx, custom.QueryApprovalTracker)
			if err != nil {
				return nil, err
			}
			return json.Marshal(badgeTypes.QueryGetApprovalTrackerResponse{Tracker: res.Tracker})
		case custom.QueryGetChallengeTracker != nil:
			res, err := keeper.GetChallengeTracker(ctx, custom.QueryGetChallengeTracker)
			if err != nil {
				return nil, err
			}
			return json.Marshal(badgeTypes.QueryGetChallengeTrackerResponse{NumUsed: res.NumUsed})
		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

type badgeCustomQuery struct {
	QueryCollection          *badgeTypes.QueryGetCollectionRequest       `json:"queryCollection,omitempty"`
	QueryBalance             *badgeTypes.QueryGetBalanceRequest          `json:"queryBalance,omitempty"`
	QueryAddressList         *badgeTypes.QueryGetAddressListRequest      `json:"queryAddressList,omitempty"`
	QueryApprovalTracker     *badgeTypes.QueryGetApprovalTrackerRequest  `json:"queryApprovalTracker,omitempty"`
	QueryGetChallengeTracker *badgeTypes.QueryGetChallengeTrackerRequest `json:"queryGetChallengeTracker,omitempty"`
}

type anchorCustomQuery struct {
	QueryValueAtLocation *anchortypes.QueryGetValueAtLocationRequest `json:"queryValueAtLocation,omitempty"`
}

type mapsCustomQuery struct {
	QueryMap      *mapsTypes.QueryGetMapRequest      `json:"queryMap,omitempty"`
	QueryMapValue *mapsTypes.QueryGetMapValueRequest `json:"queryMapList,omitempty"`
}
