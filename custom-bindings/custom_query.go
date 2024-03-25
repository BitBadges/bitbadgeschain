package custom_bindings

import (
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	badgeKeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	badgeTypes "github.com/bitbadges/bitbadgeschain/x/badges/types"

	protocolKeeper "github.com/bitbadges/bitbadgeschain/x/protocols/keeper"
	protocolTypes "github.com/bitbadges/bitbadgeschain/x/protocols/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PerformCustomBitBadgesModuleQuery(bk badgeKeeper.Keeper, pk protocolKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		isBadgeModuleQuery := false
		var custom badgeCustomQuery
		err := json.Unmarshal(request, &custom)
		if err == nil {
			isBadgeModuleQuery = true
		}

		if isBadgeModuleQuery {
			switch {
			case custom.QueryCollection != nil:
				res, err := bk.GetCollection(ctx, custom.QueryCollection)
				if err != nil {
					return nil, err
				}
				return json.Marshal(badgeTypes.QueryGetCollectionResponse{Collection: res.Collection})
			case custom.QueryBalance != nil:
				res, err := bk.GetBalance(ctx, custom.QueryBalance)
				if err != nil {
					return nil, err
				}
				return json.Marshal(badgeTypes.QueryGetBalanceResponse{Balance: res.Balance})

			case custom.QueryAddressList != nil:
				res, err := bk.GetAddressList(ctx, custom.QueryAddressList)
				if err != nil {
					return nil, err
				}
				return json.Marshal(badgeTypes.QueryGetAddressListResponse{List: res.List})
			case custom.QueryApprovalTracker != nil:
				res, err := bk.GetApprovalTracker(ctx, custom.QueryApprovalTracker)
				if err != nil {
					return nil, err
				}
				return json.Marshal(badgeTypes.QueryGetApprovalTrackerResponse{Tracker: res.Tracker})
			case custom.QueryGetChallengeTracker != nil:
				res, err := bk.GetChallengeTracker(ctx, custom.QueryGetChallengeTracker)
				if err != nil {
					return nil, err
				}
				return json.Marshal(badgeTypes.QueryGetChallengeTrackerResponse{NumUsed: res.NumUsed})
			}
		} else {
			var custom protocolCustomQuery
			err := json.Unmarshal(request, &custom)
			if err != nil {
				return nil, sdkerrors.Wrap(err, err.Error())
			}

			switch {
			case custom.QueryGetProtocol != nil:
				res, err := pk.GetProtocol(ctx, custom.QueryGetProtocol)
				if err != nil {
					return nil, err
				}
				return json.Marshal(protocolTypes.QueryGetProtocolResponse{Protocol: res.Protocol})
			case custom.QueryGetCollectionIdForProtocol != nil:
				res, err := pk.GetCollectionIdForProtocol(ctx, custom.QueryGetCollectionIdForProtocol)
				if err != nil {
					return nil, err
				}
				return json.Marshal(protocolTypes.QueryGetCollectionIdForProtocolResponse{CollectionId: res.CollectionId})
			}
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

type protocolCustomQuery struct {
	QueryGetProtocol                *protocolTypes.QueryGetProtocolRequest                `json:"queryGetProtocol,omitempty"`
	QueryGetCollectionIdForProtocol *protocolTypes.QueryGetCollectionIdForProtocolRequest `json:"queryGetCollectionIdForProtocol,omitempty"`
}
