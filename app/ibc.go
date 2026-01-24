package app

import (
	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	icamodule "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/types"
	ibccallbacks "github.com/cosmos/ibc-go/v10/modules/apps/callbacks"
	ibctransfer "github.com/cosmos/ibc-go/v10/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v10/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ibc "github.com/cosmos/ibc-go/v10/modules/core"
	ibcclienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v10/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	solomachine "github.com/cosmos/ibc-go/v10/modules/light-clients/06-solomachine"
	ibctm "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"

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

	packetforward "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward/types"

	customhooks "github.com/bitbadges/bitbadgeschain/x/custom-hooks"
	customhookskeeper "github.com/bitbadges/bitbadgeschain/x/custom-hooks/keeper"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
	ibchookstypes "github.com/bitbadges/bitbadgeschain/x/ibc-hooks/types"
	ibcratelimithooks "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/hooks"
	ibcratelimitkeeper "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	ibcratelimitmodule "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/module"
	ibcratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
)

// CombinedIBCHooks combines rate limit and custom hooks
type CombinedIBCHooks struct {
	RateLimitOverrideHooks *ibcratelimithooks.RateLimitOverrideHooks
	CustomHooks            *customhooks.CustomHooks
}

// Implement hook interfaces by delegating to the appropriate hook
var (
	_ ibchooks.OnRecvPacketOverrideHooks = &CombinedIBCHooks{}
	_ ibchooks.SendPacketOverrideHooks   = &CombinedIBCHooks{}
)

func (h *CombinedIBCHooks) OnRecvPacketOverride(im ibchooks.IBCMiddleware, ctx sdk.Context, channelID string, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
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
		return h.RateLimitOverrideHooks.OnRecvPacketOverride(wrappedIM, ctx, channelID, packet, relayer)
	}

	// If no rate limit hooks, use custom hooks if available
	if h.CustomHooks != nil {
		return h.CustomHooks.OnRecvPacketOverride(im, ctx, channelID, packet, relayer)
	}

	// Fallback to default behavior
	return im.App.OnRecvPacket(ctx, channelID, packet, relayer)
}

func (h *CombinedIBCHooks) SendPacketOverride(i ibchooks.ICS4Middleware, ctx sdk.Context, sourcePort string, sourceChannel string, timeoutHeight ibcclienttypes.Height, timeoutTimestamp uint64, data []byte) (uint64, error) {
	// Rate limit hooks take precedence - check first
	if h.RateLimitOverrideHooks != nil {
		return h.RateLimitOverrideHooks.SendPacketOverride(i, ctx, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
	}

	// Fallback to default behavior
	return i.SendPacket(ctx, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
}

// registerIBCModules register IBC keepers and non dependency inject modules.
func (app *App) registerIBCModules(appOpts servertypes.AppOptions) error {
	// set up non depinject support modules store keys
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(ibcexported.StoreKey),
		storetypes.NewKVStoreKey(ibctransfertypes.StoreKey),
		storetypes.NewKVStoreKey(icahosttypes.StoreKey),
		storetypes.NewKVStoreKey(icacontrollertypes.StoreKey),
		storetypes.NewKVStoreKey(packetforwardtypes.StoreKey),
		storetypes.NewKVStoreKey(ibchookstypes.StoreKey),
		storetypes.NewKVStoreKey(ibcratelimittypes.StoreKey),
		customhookstypes.TransientStoreKey,
	); err != nil {
		return err
	}

	// register the key tables for legacy param subspaces
	app.ParamsKeeper.Subspace(ibchookstypes.ModuleName).WithKeyTable(ibchookstypes.ParamKeyTable())

	// Create IBC keeper (IBC v10 - no capability keeper needed)
	app.IBCKeeper = ibckeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(ibcexported.StoreKey)),
		app.GetSubspace(ibcexported.ModuleName),
		app.UpgradeKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Register the proposal types
	// Deprecated: Avoid adding new handlers, instead use the new proposal flow
	// by granting the governance module the right to execute the message.
	// See: https://docs.cosmos.network/main/modules/gov#proposal-messages
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler)

	// Create IBC transfer keeper (IBC v10 - updated API)
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(ibctransfertypes.StoreKey)),
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, // ICS4Wrapper
		app.IBCKeeper.ChannelKeeper, // ChannelKeeper
		app.MsgServiceRouter(),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Build custom non-depinject modules (IBC v10 - updated API)
	// Note: PacketForwardKeeper needs to be created after transfer stack is built
	// We'll create it after the transfer stack is set up
	app.PacketForwardKeeper = packetforwardkeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(packetforwardtypes.StoreKey)),
		app.TransferKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.BankKeeper,
		app.IBCKeeper.ChannelKeeper, // ICS4Wrapper
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create interchain account keepers (IBC v10 - updated API)
	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(icahosttypes.StoreKey)),
		app.GetSubspace(icatypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, // ICS4Wrapper
		app.IBCKeeper.ChannelKeeper, // ChannelKeeper
		app.AccountKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	app.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		app.appCodec,
		runtime.NewKVStoreService(app.GetKey(icacontrollertypes.StoreKey)),
		app.GetSubspace(icatypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, // ICS4Wrapper
		app.IBCKeeper.ChannelKeeper, // ChannelKeeper
		app.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	app.GovKeeper.SetLegacyRouter(govRouter)

	// integration point for custom authentication modules (IBC v10 - no noAuthzModule)
	// Create basic ICA controller module for router (will be wrapped with callbacks later)
	icaControllerIBCModule := icacontroller.NewIBCMiddleware(app.ICAControllerKeeper)
	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	// Create static IBC router (transfer route will be added after building full stack with hooks)
	ibcRouter := porttypes.NewRouter().
		AddRoute(icacontrollertypes.SubModuleName, icaControllerIBCModule).
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)

	anchorIBCModule := anchormodule.NewIBCModule(app.AnchorKeeper)
	ibcRouter.AddRoute(anchormoduletypes.ModuleName, anchorIBCModule)
	mapsIBCModule := mapsmodule.NewIBCModule(app.MapsKeeper)
	ibcRouter.AddRoute(mapsmoduletypes.ModuleName, mapsIBCModule)
	wasmxIBCModule := wasmxmodule.NewIBCModule(app.WasmxKeeper)
	ibcRouter.AddRoute(wasmxmoduletypes.ModuleName, wasmxIBCModule)
	badgesIBCModule := badgesmodule.NewIBCModule(app.BadgesKeeper)
	ibcRouter.AddRoute(badgesmoduletypes.ModuleName, badgesIBCModule)

	// this line is used by starport scaffolding # ibc/app/module

	// wasmd 0.61.6 uses IBC v10 natively - NewIBCHandler needs IBCContractKeeper, ChannelKeeper, ICS20TransferPortSource, and appVersionGetter
	// Create appVersionGetter adapter
	appVersionGetter := &appVersionGetterAdapter{app: app}
	wasmStack := wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper, app.TransferKeeper, appVersionGetter)
	ibcRouter.AddRoute(wasmtypes.ModuleName, wasmStack)

	// Create Interchain Accounts Stack (IBC v10 - no fee middleware, no noAuthzModule)
	// SendPacket, since it is originating from the application to core IBC:
	// icaAuthModuleKeeper.SendTx -> icaController.SendPacket -> channel.SendPacket
	var icaControllerStack porttypes.IBCModule
	icaControllerStack = icacontroller.NewIBCMiddleware(app.ICAControllerKeeper)
	// IBC callbacks: wasmd 0.61.6 doesn't fully implement ContractKeeper - use no-op implementation
	// If callbacks are needed in the future, we'll need to create a ContractKeeper adapter
	icaControllerStack = ibccallbacks.NewIBCMiddleware(icaControllerStack, app.IBCKeeper.ChannelKeeper, NewNoopContractKeeper(), wasm.DefaultMaxIBCCallbackGas)
	icaICS4Wrapper := icaControllerStack.(porttypes.ICS4Wrapper)
	app.ICAControllerKeeper.WithICS4Wrapper(icaICS4Wrapper)

	// Create Transfer Stack (IBC v10 - callbacks wraps transfer, packetforward wraps callbacks)
	var transferStack porttypes.IBCModule
	transferStack = ibctransfer.NewIBCModule(app.TransferKeeper)
	// callbacks wraps the transfer stack as its base app, and uses PacketForwardKeeper as the ICS4Wrapper
	// i.e. packet-forward-middleware is higher on the stack and sits between callbacks and the ibc channel keeper
	// Since this is the lowest level middleware of the transfer stack, it should be the first entrypoint for transfer keeper's
	// WriteAcknowledgement.
	// IBC callbacks: wasmd 0.61.6 doesn't fully implement ContractKeeper - use no-op implementation
	// If callbacks are needed in the future, we'll need to create a ContractKeeper adapter
	cbStack := ibccallbacks.NewIBCMiddleware(transferStack, app.PacketForwardKeeper, NewNoopContractKeeper(), wasm.DefaultMaxIBCCallbackGas)
	transferStack = packetforward.NewIBCMiddleware(
		cbStack,
		app.PacketForwardKeeper,
		0, // retries on timeout
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp, // forward timeout
	)
	app.TransferICS4Wrapper = app.PacketForwardKeeper

	// Setup Custom Hooks Keeper with the proper ICS4Wrapper
	bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	// Pass pointer to GammKeeper to avoid copying the keeper (which contains storeKey)
	customHooksKeeper := customhookskeeper.NewKeeper(
		app.Logger(),
		&app.GammKeeper,
		app.BankKeeper,
		&app.BadgesKeeper,
		&app.SendmanagerKeeper,
		app.TransferKeeper,
		app.TransferICS4Wrapper,
		app.IBCKeeper.ChannelKeeper,
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
	// Use IBCKeeper.ChannelKeeper as the ICS4Wrapper since it implements the interface
	// The channel field is only used for SendPacket operations, not OnRecvPacket hooks
	app.HooksICS4Wrapper = ibchooks.NewICS4Middleware(
		app.IBCKeeper.ChannelKeeper,
		combinedHooks,
	)

	// Add IBC Hooks middleware last (outermost) to have full control over acknowledgements
	// This ensures error acknowledgements from hooks are properly returned
	hooksTransferModule := ibchooks.NewIBCMiddleware(transferStack, &app.HooksICS4Wrapper)
	transferStack = hooksTransferModule
	app.TransferKeeper.WithICS4Wrapper(cbStack)

	// Add the transfer stack (with hooks) to the router
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)

	app.IBCKeeper.SetRouter(ibcRouter)

	// register IBC modules (IBC v10 - no capability or fee modules)
	if err := app.RegisterModules(
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.TransferKeeper),
		icamodule.NewAppModule(&app.ICAControllerKeeper, &app.ICAHostKeeper),
		ibctm.NewAppModule(ibctm.NewLightClientModule(app.appCodec, ibcclienttypes.NewStoreProvider(runtime.NewKVStoreService(app.GetKey(ibcexported.StoreKey))))),
		solomachine.NewAppModule(solomachine.NewLightClientModule(app.appCodec, ibcclienttypes.NewStoreProvider(runtime.NewKVStoreService(app.GetKey(ibcexported.StoreKey))))),
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
		icatypes.ModuleName:           icamodule.AppModule{},
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
