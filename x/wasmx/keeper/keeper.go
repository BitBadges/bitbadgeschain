package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/bitbadges/bitbadgeschain/x/wasmx/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

// Keeper of this module maintains collections of wasmx.
type Keeper struct {
	storeKey   storetypes.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramtypes.Subspace

	accountKeeper authkeeper.AccountKeeper
	bankKeeper    types.BankKeeper
	wasmKeeper    wasmkeeper.Keeper
}

// NewKeeper creates new instances of the wasmx Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	ak authkeeper.AccountKeeper,
	bk types.BankKeeper,
	wk wasmkeeper.Keeper,
) Keeper {

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		paramSpace: paramSpace,

		storeKey:      storeKey,
		cdc:           cdc,
		accountKeeper: ak,
		bankKeeper:    bk,
		wasmKeeper: 	 wk,
	}
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

func (k *Keeper) getStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k *Keeper) SetWasmKeeper(wk wasmkeeper.Keeper) {
	k.wasmKeeper = wk
}
