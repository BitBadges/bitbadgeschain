package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		bankKeeper types.BankKeeper
		accountKeeper types.AccountKeeper

		ApprovedContractAddresses []string
		PayoutAddress             string
		EnableCoinTransfers       bool
		AllowedDenoms             []string
		FixedCostPerTransfer      string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	ApprovedContractAddresses []string,
	PayoutAddress string,
	EnableCoinTransfers bool,
	AllowedDenoms []string,
	FixedCostPerTransfer string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:                       cdc,
		storeService:              storeService,
		authority:                 authority,
		logger:                    logger,
		bankKeeper:                bankKeeper,
		accountKeeper:             accountKeeper,
		ApprovedContractAddresses: ApprovedContractAddresses,
		PayoutAddress:             PayoutAddress,
		EnableCoinTransfers:       EnableCoinTransfers,
		AllowedDenoms:             AllowedDenoms,
		FixedCostPerTransfer:      FixedCostPerTransfer,
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

// // ChanCloseInit defines a wrapper function for the channel Keeper's function.
// func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
// 	capName := host.ChannelCapabilityPath(portID, channelID)
// 	chanCap, ok := k.ScopedKeeper().GetCapability(ctx, capName)
// 	if !ok {
// 		return errorsmod.Wrapf(channeltypes.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
// 	}
// 	return k.ibcKeeperFn().ChannelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
// }

// // ShouldBound checks if the IBC app module can be bound to the desired port
// func (k Keeper) ShouldBound(ctx sdk.Context, portID string) bool {
// 	scopedKeeper := k.ScopedKeeper()
// 	if scopedKeeper == nil {
// 		return false
// 	}
// 	_, ok := scopedKeeper.GetCapability(ctx, host.PortPath(portID))
// 	return !ok
// }

// // BindPort defines a wrapper function for the port Keeper's function in
// // order to expose it to module's InitGenesis function
// func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
// 	cap := k.ibcKeeperFn().PortKeeper.BindPort(ctx, portID)
// 	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
// }

// // GetPort returns the portID for the IBC app module. Used in ExportGenesis
// func (k Keeper) GetPort(ctx sdk.Context) string {
// 	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
// 	store := prefix.NewStore(storeAdapter, []byte{})
// 	return string(store.Get(types.PortKey))
// }

// // SetPort sets the portID for the IBC app module. Used in InitGenesis
// func (k Keeper) SetPort(ctx sdk.Context, portID string) {
// 	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
// 	store := prefix.NewStore(storeAdapter, []byte{})
// 	store.Set(types.PortKey, []byte(portID))
// }

// // AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
// func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
// 	return k.ScopedKeeper().AuthenticateCapability(ctx, cap, name)
// }

// // ClaimCapability allows the IBC app module to claim a capability that core IBC
// // passes to it
// func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
// 	return k.ScopedKeeper().ClaimCapability(ctx, cap, name)
// }

// // ScopedKeeper returns the ScopedKeeper
// func (k Keeper) ScopedKeeper() exported.ScopedKeeper {
// 	return k.capabilityScopedFn(types.ModuleName)
// }
