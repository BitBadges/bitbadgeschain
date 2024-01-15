package custom_bindings

import (
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	protocolKeeper "github.com/bitbadges/bitbadgeschain/x/protocols/keeper"
	protocolTypes "github.com/bitbadges/bitbadgeschain/x/protocols/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WASM handler for contracts calling into the protocols module
func PerformCustomProtocolQuery(keeper protocolKeeper.Keeper) wasmKeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var custom protocolCustomQuery
		err := json.Unmarshal(request, &custom)
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		switch {
		case custom.QueryGetProtocol != nil:
		res, err := keeper.GetProtocol(ctx, custom.QueryGetProtocol)
			if err != nil {
				return nil, err
			}
			return json.Marshal(protocolTypes.QueryGetProtocolResponse{Protocol: res.Protocol})
		case custom.QueryGetCollectionIdForProtocol != nil:
			res, err := keeper.GetCollectionIdForProtocol(ctx, custom.QueryGetCollectionIdForProtocol)
			if err != nil {
				return nil, err
			}
			return json.Marshal(protocolTypes.QueryGetCollectionIdForProtocolResponse{CollectionId: res.CollectionId})
		}
		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown Custom query variant")
	}
}

type protocolCustomQuery struct {
	QueryGetProtocol 								*protocolTypes.QueryGetProtocolRequest                   `json:"queryGetProtocol,omitempty"`
	QueryGetCollectionIdForProtocol 	*protocolTypes.QueryGetCollectionIdForProtocolRequest       `json:"queryGetCollectionIdForProtocol,omitempty"`
}
