package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		tokenizationKeeper *tokenizationkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	tokenizationKeeper *tokenizationkeeper.Keeper,
) Keeper {
	return Keeper{
		cdc:                  cdc,
		storeService:         storeService,
		authority:            authority,
		logger:               logger,
		tokenizationKeeper:   tokenizationKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

