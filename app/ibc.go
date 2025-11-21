package app

import (
	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	icamodule "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	solomachine "github.com/cosmos/ibc-go/v8/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	// this line is used by starport scaffolding # ibc/app/import
	anchormodule "github.com/bitbadges/bitbadgeschain/x/anchor/module"
	anchormoduletypes "github.com/bitbadges/bitbadgeschain/x/anchor/types"
	badgesmoduletypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	mapsmodule "github.com/bitbadges/bitbadgeschain/x/maps/module"
	mapsmoduletypes "github.com/bitbadges/bitbadgeschain/x/maps/types"
	wasmxmodule "github.com/bitbadges/bitbadgeschain/x/wasmx/module"
	wasmxmoduletypes "github.com/bitbadges/bitbadgeschain/x/wasmx/types"

	badgesmodule "github.com/bitbadges/bitbadgeschain/x/badges/module"

	wasm "github.com/CosmWasm/wasmd/x/wasm"

	packetforward "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v8/packetforward/types"
	ibccallbacks "github.com/cosmos/ibc-go/modules/apps/callbacks"

	customhooks "github.com/bitbadges/bitbadgeschain/x/custom-hooks"
	customhookskeeper "github.com/bitbadges/bitbadgeschain/x/custom-hooks/keeper"
	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
	ibchookstypes "github.com/bitbadges/bitbadgeschain/x/ibc-hooks/types"
	ibcratelimithooks "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/hooks"
	ibcratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	ibcratelimitmodule "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/module"
	ibcratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

// CombinedIBCHooks combines rate limit and custom hooks
type CombinedIBCHooks struct {
	RateLimitOverrideHooks *ibcratelimithooks.RateLimitOverrideHooks
	CustomHooks            *customhooks.CustomHooks
}

// Implement hook interfaces by delegating to the appropriate hook
var _ ibchooks.OnRecvPacketOverrideHooks = &CombinedIBCHooks{}
var _ ibchooks.SendPacketOverrideHooks = &CombinedIBCHooks{}

func (h *CombinedIBCHooks) OnRecvPacketOverride(im ibchooks.IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	// Rate limit hooks take precedence - check first
	if h.RateLimitOverrideHooks != nil {
		// Create a wrapper IBCMiddleware that chains custom hooks after rate limit
		// The rate limit hooks will call im.App.OnRecvPacket, so we wrap that to include custom hooks
		wrappedApp := &customHooksWrapper{
			app:         im.App,
			customHooks: h.CustomHooks,
		}
		// Create a new IBCMiddleware with the wrapped app
		wrappedIM := ibchooks.NewIBCMiddleware(wrappedApp, im.ICS4Middleware)
		return h.RateLimitOverrideHooks.OnRecvPacketOverride(wrappedIM, ctx, packet, relayer)
	}

	// If no rate limit hooks, use custom hooks if available
	if h.CustomHooks != nil {
		return h.CustomHooks.OnRecvPacketOverride(im, ctx, packet, relayer)
	}

	// Fallback to default behavior
	return im.App.OnRecvPacket(ctx, packet, relayer)
}

func (h *CombinedIBCHooks) SendPacketOverride(i ibchooks.ICS4Middleware, ctx sdk.Context, chanCap *capabilitytypes.Capability, sourcePort string, sourceChannel string, timeoutHeight ibcclienttypes.Height, timeoutTimestamp uint64, data []byte) (uint64, error) {
	// Rate limit hooks take precedence - check first
	if h.RateLimitOverrideHooks != nil {
		return h.RateLimitOverrideHooks.SendPacketOverride(i, ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// Fallback to default behavior
	return i.SendPacket(ctx, chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// registerIBCModules register IBC keepers and non dependency inject modules.
func (app *App) registerIBCModules(appOpts servertypes.AppOptions) error {
	// set up non depinject support modules store keys
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(capabilitytypes.StoreKey),
		storetypes.NewKVStoreKey(ibcexported.StoreKey),
		storetypes.NewKVStoreKey(ibctransfertypes.StoreKey),
		storetypes.NewKVStoreKey(ibcfeetypes.StoreKey),
		storetypes.NewKVStoreKey(icahosttypes.StoreKey),
		storetypes.NewKVStoreKey(icacontrollertypes.StoreKey),
		storetypes.NewKVStoreKey(packetforwardtypes.StoreKey),
		storetypes.NewKVStoreKey(ibchookstypes.StoreKey),
		storetypes.NewKVStoreKey(ibcratelimittypes.StoreKey),
		storetypes.NewMemoryStoreKey(capabilitytypes.MemStoreKey),
		storetypes.NewTransientStoreKey(paramstypes.TStoreKey),
	); err != nil {
		return err
	}

	// register the key tables for legacy param subspaces
	keyTable := ibcclienttypes.ParamKeyTable()
	keyTable.RegisterParamSet(&ibcconnectiontypes.Params{})
	app.ParamsKeeper.Subspace(ibcexported.ModuleName).WithKeyTable(keyTable)
	app.ParamsKeeper.Subspace(ibctransfertypes.ModuleName).WithKeyTable(ibctransfertypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(icacontrollertypes.SubModuleName).WithKeyTable(icacontrollertypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(icahosttypes.SubModuleName).WithKeyTable(icahosttypes.ParamKeyTable())
	app.ParamsKeeper.Subspace(ibchookstypes.ModuleName).WithKeyTable(ibchookstypes.ParamKeyTable())

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		app.AppCodec(),
		app.GetKey(capabilitytypes.StoreKey),
		app.GetMemKey(capabilitytypes.MemStoreKey),
	)

	// add capability keeper and ScopeToModule for ibc module
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedIBCTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)

	// Create IBC keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		app.appCodec,
		app.GetKey(ibcexported.StoreKey),
		app.GetSubspace(ibcexported.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Register the proposal types
	// Deprecated: Avoid adding new handlers, instead use the new proposal flow
	// by granting the governance module the right to execute the message.
	// See: https://docs.cosmos.network/main/modules/gov#proposal-messages
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler)

	//Build custom non-depinject modules
	app.PacketForwardKeeper = packetforwardkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(packetforwardtypes.StoreKey),
		app.TransferKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.BankKeeper,
		app.IBCFeeKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		app.appCodec, app.GetKey(ibcfeetypes.StoreKey),
		app.IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper, app.AccountKeeper, app.BankKeeper,
	)

	// Create IBC transfer keeper
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(ibctransfertypes.StoreKey),
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCFeeKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedIBCTransferKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create interchain account keepers
	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(icahosttypes.StoreKey),
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.ICAHostKeeper.WithQueryRouter(app.GRPCQueryRouter())

	app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(icacontrollertypes.StoreKey),
		app.GetSubspace(icacontrollertypes.SubModuleName),
		app.IBCFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper,
		app.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.GovKeeper.SetLegacyRouter(govRouter)

	// integration point for custom authentication modules
	var noAuthzModule porttypes.IBCModule
	icaControllerIBCModule := ibcfee.NewIBCMiddleware(
		icacontroller.NewIBCMiddleware(noAuthzModule, app.ICAControllerKeeper),
		app.IBCFeeKeeper,
	)

	icaHostIBCModule := ibcfee.NewIBCMiddleware(icahost.NewIBCModule(app.ICAHostKeeper), app.IBCFeeKeeper)

	// Create static IBC router (transfer route will be added after building full stack with hooks)
	ibcRouter := porttypes.NewRouter().
		AddRoute(icacontrollertypes.SubModuleName, icaControllerIBCModule).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)

	anchorIBCModule := ibcfee.NewIBCMiddleware(anchormodule.NewIBCModule(app.AnchorKeeper), app.IBCFeeKeeper)
	ibcRouter.AddRoute(anchormoduletypes.ModuleName, anchorIBCModule)
	mapsIBCModule := ibcfee.NewIBCMiddleware(mapsmodule.NewIBCModule(app.MapsKeeper), app.IBCFeeKeeper)
	ibcRouter.AddRoute(mapsmoduletypes.ModuleName, mapsIBCModule)
	wasmxIBCModule := ibcfee.NewIBCMiddleware(wasmxmodule.NewIBCModule(app.WasmxKeeper), app.IBCFeeKeeper)
	ibcRouter.AddRoute(wasmxmoduletypes.ModuleName, wasmxIBCModule)
	badgesIBCModule := ibcfee.NewIBCMiddleware(badgesmodule.NewIBCModule(app.BadgesKeeper), app.IBCFeeKeeper)
	ibcRouter.AddRoute(badgesmoduletypes.ModuleName, badgesIBCModule)

	// this line is used by starport scaffolding # ibc/app/module

	var wasmStack porttypes.IBCModule
	wasmStackIBCHandler := wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper, app.IBCFeeKeeper)
	wasmStack = wasmStackIBCHandler
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, app.IBCFeeKeeper)
	ibcRouter.AddRoute(wasmtypes.ModuleName, wasmStack)

	// Create Interchain Accounts Stack
	// SendPacket, since it is originating from the application to core IBC:
	// icaAuthModuleKeeper.SendTx -> icaController.SendPacket -> fee.SendPacket -> channel.SendPacket
	var icaControllerStack porttypes.IBCModule
	icaControllerStack = icacontroller.NewIBCMiddleware(noAuthzModule, app.ICAControllerKeeper)
	icaControllerStack = icacontroller.NewIBCMiddleware(icaControllerStack, app.ICAControllerKeeper)
	icaControllerStack = ibccallbacks.NewIBCMiddleware(icaControllerStack, app.IBCFeeKeeper, wasmStackIBCHandler, wasm.DefaultMaxIBCCallbackGas)
	icaICS4Wrapper := icaControllerStack.(porttypes.ICS4Wrapper)
	icaControllerStack = ibcfee.NewIBCMiddleware(icaControllerStack, app.IBCFeeKeeper)
	app.ICAControllerKeeper.WithICS4Wrapper(icaICS4Wrapper)

	// Create Transfer Stack
	var transferStack porttypes.IBCModule
	transferStack = transfer.NewIBCModule(app.TransferKeeper)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, app.IBCFeeKeeper)
	transferStack = ibccallbacks.NewIBCMiddleware(transferStack, app.IBCFeeKeeper, wasmStackIBCHandler, wasm.DefaultMaxIBCCallbackGas)
	app.TransferICS4Wrapper = transferStack.(porttypes.ICS4Wrapper)

	// Setup Custom Hooks Keeper with the proper ICS4Wrapper
	bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	// Pass pointer to GammKeeper to avoid copying the keeper (which contains storeKey)
	customHooksKeeper := customhookskeeper.NewKeeper(
		app.Logger(),
		&app.GammKeeper,
		app.BankKeeper,
		&app.BadgesKeeper,
		app.TransferICS4Wrapper,
		app.IBCKeeper.ChannelKeeper,
		scopedIBCTransferKeeper,
	)

	// Setup Custom Hooks (standalone)
	customHooks := customhooks.NewCustomHooks(customHooksKeeper, bech32Prefix)

	// Setup IBC Rate Limit Keeper
	// Authority defaults to gov module account
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	app.IBCRateLimitKeeper = ibcratelimitkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(ibcratelimittypes.StoreKey),
		app.BankKeeper,
		authority,
	)

	// IBC rate limit module will be registered via RegisterModules below

	// Setup IBC Rate Limit Hooks
	rateLimitOverrideHooks := ibcratelimithooks.NewRateLimitOverrideHooks(app.IBCRateLimitKeeper)

	// Combine hooks: rate limit override hooks take precedence, then custom hooks
	// The rate limit hooks need to be checked first to reject packets before processing
	combinedHooks := &CombinedIBCHooks{
		RateLimitOverrideHooks: rateLimitOverrideHooks,
		CustomHooks:            customHooks,
	}

	// Setup ICS4 Wrapper for hooks (with rate limit + custom hooks)
	// Use IBCFeeKeeper as the ICS4Wrapper since it implements the interface
	// The channel field is only used for SendPacket operations, not OnRecvPacket hooks
	app.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		app.IBCFeeKeeper,
		combinedHooks,
	)

	// Add packetforward middleware first (before hooks)
	transferStack = packetforward.NewIBCMiddleware(
		transferStack,
		app.PacketForwardKeeper,
		0, // retries on timeout
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp, // forward timeout
	)

	// Add IBC Hooks middleware last (outermost) to have full control over acknowledgements
	// This ensures error acknowledgements from hooks are properly returned
	hooksTransferModule := ibchooks.NewIBCMiddleware(transferStack, &app.HooksICS4Wrapper)
	transferStack = hooksTransferModule
	app.TransferKeeper.WithICS4Wrapper(app.TransferICS4Wrapper)

	// Add the transfer stack (with hooks) to the router
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)

	// RecvPacket, message that originates from core IBC and goes down to app, the flow is:
	// channel.RecvPacket -> fee.OnRecvPacket -> icaHost.OnRecvPacket
	var icaHostStack porttypes.IBCModule
	icaHostStack = icahost.NewIBCModule(app.ICAHostKeeper)
	icaHostStack = ibcfee.NewIBCMiddleware(icaHostStack, app.IBCFeeKeeper)

	app.IBCKeeper.SetRouter(ibcRouter)

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedIBCTransferKeeper = scopedIBCTransferKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper
	app.ScopedICAControllerKeeper = scopedICAControllerKeeper

	// register IBC modules
	if err := app.RegisterModules(
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.TransferKeeper),
		ibcfee.NewAppModule(app.IBCFeeKeeper),
		icamodule.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		capability.NewAppModule(app.appCodec, *app.CapabilityKeeper, false),
		ibctm.NewAppModule(),
		solomachine.NewAppModule(),
		packetforward.NewAppModule(app.PacketForwardKeeper, app.GetSubspace(packetforwardtypes.ModuleName)),
		ibcratelimitmodule.NewAppModule(app.appCodec, app.IBCRateLimitKeeper),
	); err != nil {
		return err
	}

	return nil
}

// RegisterIBC Since the IBC modules don't support dependency injection,
// we need to manually register the modules on the client side.
// This needs to be removed after IBC supports App Wiring.
func RegisterIBC(registry cdctypes.InterfaceRegistry) map[string]appmodule.AppModule {
	modules := map[string]appmodule.AppModule{
		ibcexported.ModuleName:        ibc.AppModule{},
		ibctransfertypes.ModuleName:   ibctransfer.AppModule{},
		ibcfeetypes.ModuleName:        ibcfee.AppModule{},
		icatypes.ModuleName:           icamodule.AppModule{},
		capabilitytypes.ModuleName:    capability.AppModule{},
		ibctm.ModuleName:              ibctm.AppModule{},
		solomachine.ModuleName:        solomachine.AppModule{},
		packetforwardtypes.ModuleName: packetforward.AppModule{},
		ibcratelimittypes.ModuleName:  ibcratelimitmodule.AppModule{},
	}

	for name, m := range modules {
		module.CoreAppModuleBasicAdaptor(name, m).RegisterInterfaces(registry)
	}

	return modules
}
