package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/council/types"
)

// MsgRouter is the interface for dispatching sdk.Msg as the council's account.
// In production this is wired to baseapp.MsgServiceRouter; in tests use a mock.
type MsgRouter interface {
	// Handler returns a handler for the given sdk.Msg, or nil if not found.
	Handler(msg sdk.Msg) MsgHandler
}

// MsgHandler executes a single sdk.Msg and returns a result.
type MsgHandler func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error)

// Keeper holds references to external keepers and its own KV store.
type Keeper struct {
	cdc      codec.BinaryCodec
	logger   log.Logger
	storeKey storetypes.StoreKey

	tokenizationKeeper types.TokenizationKeeper
	msgRouter          MsgRouter
}

// NewKeeper creates a new x/council keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	logger log.Logger,
	tokenizationKeeper types.TokenizationKeeper,
	msgRouter MsgRouter,
) Keeper {
	return Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		logger:             logger,
		tokenizationKeeper: tokenizationKeeper,
		msgRouter:          msgRouter,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
