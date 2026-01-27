package keeper

import (
	badgesmodulekeeper "github.com/bitbadges/bitbadgeschain/x/badges/keeper"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
)

// Keeper manages the twofa module state and logic
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
	logger       log.Logger

	badgesKeeper *badgesmodulekeeper.Keeper
}

// NewKeeper creates a new twofa Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	badgesKeeper *badgesmodulekeeper.Keeper,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		logger:       logger,
		badgesKeeper: badgesKeeper,
	}
}

// GetBadgesKeeper returns the badges keeper (for testing)
func (k Keeper) GetBadgesKeeper() *badgesmodulekeeper.Keeper {
	return k.badgesKeeper
}

// Logger returns the logger
func (k Keeper) Logger() log.Logger {
	return k.logger
}

