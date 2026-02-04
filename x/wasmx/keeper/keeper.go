package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	"github.com/bitbadges/bitbadgeschain/x/wasmx/types"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		ibcKeeperFn func() *ibckeeper.Keeper

		wasmKeeper wasmkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	ibcKeeperFn func() *ibckeeper.Keeper,
	wasmKeeper wasmkeeper.Keeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
		ibcKeeperFn:  ibcKeeperFn,
		wasmKeeper:   wasmKeeper,
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

// ----------------------------------------------------------------------------
// IBC Keeper Logic
// ----------------------------------------------------------------------------

// ChanCloseInit defines a wrapper function for the channel Keeper's function.
// IBC v10: capabilities removed, no capability parameter needed
func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return k.ibcKeeperFn().ChannelKeeper.ChanCloseInit(ctx, portID, channelID)
}

// ShouldBound checks if the IBC app module can be bound to the desired port
// IBC v10: ports are managed automatically, this always returns true
func (k Keeper) ShouldBound(ctx sdk.Context, portID string) bool {
	return true
}

// BindPort defines a wrapper function for the port Keeper's function in
// order to expose it to module's InitGenesis function
// IBC v10: ports are managed automatically, no binding needed
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	// In IBC v10, ports are managed automatically - no action needed
	return nil
}

// GetPort returns the portID for the IBC app module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the IBC app module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, []byte{})
	store.Set(types.PortKey, []byte(portID))
}

// IBC v10: Capability-related methods removed as capabilities are no longer used
