package custom_bindings

import (
	"encoding/json"

	anchorKeeper "github.com/bitbadges/bitbadgeschain/x/anchor/keeper"
	anchortypes "github.com/bitbadges/bitbadgeschain/x/anchor/types"
	tokenizationKeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenTypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	gammKeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	gammTypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	managersplitterKeeper "github.com/bitbadges/bitbadgeschain/x/managersplitter/keeper"
	managersplitterTypes "github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	mapsKeeper "github.com/bitbadges/bitbadgeschain/x/maps/keeper"
	mapsTypes "github.com/bitbadges/bitbadgeschain/x/maps/types"

	sdkerrors "cosmossdk.io/errors"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func PerformCustomBitBadgesModuleQuery(bk tokenizationKeeper.Keeper, ak anchorKeeper.Keeper, mk mapsKeeper.Keeper, gk gammKeeper.Keeper, msk managersplitterKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		isTokenizationModuleQuery := false
		var custom tokenizationCustomQuery
		err := json.Unmarshal(request, &custom)
		if err == nil {
			isTokenizationModuleQuery = true
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

		isGammModuleQuery := false
		var gammCustom gammCustomQuery
		err = json.Unmarshal(request, &gammCustom)
		if err == nil {
			isGammModuleQuery = true
		}

		isManagersplitterModuleQuery := false
		var managersplitterCustom managersplitterCustomQuery
		err = json.Unmarshal(request, &managersplitterCustom)
		if err == nil {
			isManagersplitterModuleQuery = true
		}

		if isTokenizationModuleQuery {
			return PerformCustomTokenizationQuery(bk)(ctx, request)
		} else if isAnchorModuleQuery {
			return PerformCustomAnchorQuery(ak)(ctx, request)
		} else if isMapsModuleQuery {
			return PerformCustomMapsQuery(mk)(ctx, request)
		} else if isGammModuleQuery {
			return PerformCustomGammQuery(gk)(ctx, request)
		} else if isManagersplitterModuleQuery {
			return PerformCustomManagersplitterQuery(msk)(ctx, request)
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

// WASM handler for contracts calling into the tokens module
func PerformCustomTokenizationQuery(keeper tokenizationKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom tokenizationCustomQuery
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
			return json.Marshal(tokenTypes.QueryGetCollectionResponse{Collection: res.Collection})
		case custom.QueryBalance != nil:
			res, err := keeper.GetBalance(ctx, custom.QueryBalance)
			if err != nil {
				return nil, err
			}
			return json.Marshal(tokenTypes.QueryGetBalanceResponse{Balance: res.Balance})

		case custom.QueryAddressList != nil:
			res, err := keeper.GetAddressList(ctx, custom.QueryAddressList)
			if err != nil {
				return nil, err
			}
			return json.Marshal(tokenTypes.QueryGetAddressListResponse{List: res.List})
		case custom.QueryApprovalTracker != nil:
			res, err := keeper.GetApprovalTracker(ctx, custom.QueryApprovalTracker)
			if err != nil {
				return nil, err
			}
			return json.Marshal(tokenTypes.QueryGetApprovalTrackerResponse{Tracker: res.Tracker})
		case custom.QueryGetChallengeTracker != nil:
			res, err := keeper.GetChallengeTracker(ctx, custom.QueryGetChallengeTracker)
			if err != nil {
				return nil, err
			}
			return json.Marshal(tokenTypes.QueryGetChallengeTrackerResponse{NumUsed: res.NumUsed})
		case custom.QueryGetETHSignatureTracker != nil:
			res, err := keeper.GetETHSignatureTracker(ctx, custom.QueryGetETHSignatureTracker)
			if err != nil {
				return nil, err
			}
			return json.Marshal(tokenTypes.QueryGetETHSignatureTrackerResponse{NumUsed: res.NumUsed})
		case custom.QueryGetWrappableBalances != nil:
			res, err := keeper.GetWrappableBalances(ctx, custom.QueryGetWrappableBalances)
			if err != nil {
				return nil, err
			}
			return json.Marshal(tokenTypes.QueryGetWrappableBalancesResponse{Amount: res.Amount})
		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

func PerformCustomGammQuery(gk gammKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom gammCustomQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		// Create a querier to handle the gRPC-style queries
		querier := gammKeeper.NewQuerier(gk)

		// Convert sdk.Context to context.Context for gRPC methods
		grpcCtx := sdk.WrapSDKContext(ctx)

		switch {
		case custom.QueryPool != nil:
			res, err := querier.Pool(grpcCtx, custom.QueryPool)
			if err != nil {
				return nil, err
			}
			return json.Marshal(gammTypes.QueryPoolResponse{Pool: res.Pool})
		case custom.QueryPools != nil:
			res, err := querier.Pools(grpcCtx, custom.QueryPools)
			if err != nil {
				return nil, err
			}
			return json.Marshal(gammTypes.QueryPoolsResponse{Pools: res.Pools, Pagination: res.Pagination})
		case custom.QueryPoolType != nil:
			res, err := querier.PoolType(grpcCtx, custom.QueryPoolType)
			if err != nil {
				return nil, err
			}
			return json.Marshal(gammTypes.QueryPoolTypeResponse{PoolType: res.PoolType})
		case custom.QueryPoolsWithFilter != nil:
			res, err := querier.PoolsWithFilter(grpcCtx, custom.QueryPoolsWithFilter)
			if err != nil {
				return nil, err
			}
			return json.Marshal(gammTypes.QueryPoolsWithFilterResponse{Pools: res.Pools, Pagination: res.Pagination})
		case custom.QueryNumPools != nil:
			res, err := querier.NumPools(grpcCtx, custom.QueryNumPools)
			if err != nil {
				return nil, err
			}
			return json.Marshal(gammTypes.QueryNumPoolsResponse{NumPools: res.NumPools})
		case custom.QueryTotalLiquidity != nil:
			res, err := querier.TotalLiquidity(grpcCtx, custom.QueryTotalLiquidity)
			if err != nil {
				return nil, err
			}
			return json.Marshal(gammTypes.QueryTotalLiquidityResponse{Liquidity: res.Liquidity})
		case custom.QueryTotalPoolLiquidity != nil:
			res, err := querier.TotalPoolLiquidity(grpcCtx, custom.QueryTotalPoolLiquidity)
			if err != nil {
				return nil, err
			}
			return json.Marshal(gammTypes.QueryTotalPoolLiquidityResponse{Liquidity: res.Liquidity})
		case custom.QuerySpotPrice != nil:
			res, err := querier.SpotPrice(grpcCtx, custom.QuerySpotPrice)
			if err != nil {
				return nil, err
			}
			return json.Marshal(gammTypes.QuerySpotPriceResponse{SpotPrice: res.SpotPrice})
		case custom.QueryPoolParams != nil:
			res, err := querier.PoolParams(grpcCtx, custom.QueryPoolParams)
			if err != nil {
				return nil, err
			}
			return json.Marshal(gammTypes.QueryPoolParamsResponse{Params: res.Params})
		case custom.QueryTotalShares != nil:
			res, err := querier.TotalShares(grpcCtx, custom.QueryTotalShares)
			if err != nil {
				return nil, err
			}
			return json.Marshal(gammTypes.QueryTotalSharesResponse{TotalShares: res.TotalShares})
		}

		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

func PerformCustomManagersplitterQuery(msk managersplitterKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom managersplitterCustomQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		// Create a querier to handle the gRPC-style queries
		grpcCtx := sdk.WrapSDKContext(ctx)

		switch {
		case custom.QueryManagerSplitter != nil:
			res, err := msk.ManagerSplitter(grpcCtx, custom.QueryManagerSplitter)
			if err != nil {
				return nil, err
			}
			return json.Marshal(managersplitterTypes.QueryGetManagerSplitterResponse{ManagerSplitter: res.ManagerSplitter})
		case custom.QueryAllManagerSplitters != nil:
			res, err := msk.AllManagerSplitters(grpcCtx, custom.QueryAllManagerSplitters)
			if err != nil {
				return nil, err
			}
			return json.Marshal(managersplitterTypes.QueryAllManagerSplittersResponse{ManagerSplitters: res.ManagerSplitters, Pagination: res.Pagination})
		case custom.QueryParams != nil:
			res, err := msk.Params(grpcCtx, custom.QueryParams)
			if err != nil {
				return nil, err
			}
			return json.Marshal(managersplitterTypes.QueryParamsResponse{Params: res.Params})
		}

		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

type tokenizationCustomQuery struct {
	QueryCollection             *tokenTypes.QueryGetCollectionRequest          `json:"queryCollection,omitempty"`
	QueryBalance                *tokenTypes.QueryGetBalanceRequest             `json:"queryBalance,omitempty"`
	QueryAddressList            *tokenTypes.QueryGetAddressListRequest         `json:"queryAddressList,omitempty"`
	QueryApprovalTracker        *tokenTypes.QueryGetApprovalTrackerRequest     `json:"queryApprovalTracker,omitempty"`
	QueryGetChallengeTracker    *tokenTypes.QueryGetChallengeTrackerRequest    `json:"queryGetChallengeTracker,omitempty"`
	QueryGetETHSignatureTracker *tokenTypes.QueryGetETHSignatureTrackerRequest `json:"queryGetETHSignatureTracker,omitempty"`
	QueryGetWrappableBalances   *tokenTypes.QueryGetWrappableBalancesRequest   `json:"queryGetWrappableBalances,omitempty"`
}

type anchorCustomQuery struct {
	QueryValueAtLocation *anchortypes.QueryGetValueAtLocationRequest `json:"queryValueAtLocation,omitempty"`
}

type mapsCustomQuery struct {
	QueryMap      *mapsTypes.QueryGetMapRequest      `json:"queryMap,omitempty"`
	QueryMapValue *mapsTypes.QueryGetMapValueRequest `json:"queryMapList,omitempty"`
}

type gammCustomQuery struct {
	QueryPool               *gammTypes.QueryPoolRequest               `json:"queryPool,omitempty"`
	QueryPools              *gammTypes.QueryPoolsRequest              `json:"queryPools,omitempty"`
	QueryPoolType           *gammTypes.QueryPoolTypeRequest           `json:"queryPoolType,omitempty"`
	QueryPoolsWithFilter    *gammTypes.QueryPoolsWithFilterRequest    `json:"queryPoolsWithFilter,omitempty"`
	QueryNumPools           *gammTypes.QueryNumPoolsRequest           `json:"queryNumPools,omitempty"`
	QueryTotalLiquidity     *gammTypes.QueryTotalLiquidityRequest     `json:"queryTotalLiquidity,omitempty"`
	QueryTotalPoolLiquidity *gammTypes.QueryTotalPoolLiquidityRequest `json:"queryTotalPoolLiquidity,omitempty"`
	QuerySpotPrice          *gammTypes.QuerySpotPriceRequest          `json:"querySpotPrice,omitempty"`
	QueryPoolParams         *gammTypes.QueryPoolParamsRequest         `json:"queryPoolParams,omitempty"`
	QueryTotalShares        *gammTypes.QueryTotalSharesRequest        `json:"queryTotalShares,omitempty"`
}

type managersplitterCustomQuery struct {
	QueryManagerSplitter     *managersplitterTypes.QueryGetManagerSplitterRequest  `json:"queryManagerSplitter,omitempty"`
	QueryAllManagerSplitters *managersplitterTypes.QueryAllManagerSplittersRequest `json:"queryAllManagerSplitters,omitempty"`
	QueryParams              *managersplitterTypes.QueryParamsRequest              `json:"queryParams,omitempty"`
}
