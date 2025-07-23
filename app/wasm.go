package app

// this line is used by starport scaffolding # ibc/app/import

// // registerIBCModules register IBC keepers and non dependency inject modules.
// func (app *App) registerWasmModules(appOpts servertypes.AppOptions) (store.KVStoreService, error) {
// 	wasmKey := storetypes.NewKVStoreKey(wasmtypes.StoreKey)
// 	wasmxKey := storetypes.NewKVStoreKey(wasmxmoduletypes.StoreKey)

// 	// set up non depinject support modules store keys
// 	if err := app.RegisterStores(
// 		wasmKey,
// 		wasmxKey,
// 	); err != nil {
// 		return nil, err
// 	}

// 	app.ParamsKeeper.Subspace(wasmtypes.ModuleName)

// 	wasmDir := filepath.Join(DefaultNodeHome, "wasm")
// 	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
// 	if err != nil {
// 		panic(fmt.Sprintf("error while reading wasm config: %s", err))
// 	}

// 	customEncoderOptions := GetCustomMsgEncodersOptions()
// 	customQueryOptions := GetCustomMsgQueryOptions(app.BadgesKeeper, app.AnchorKeeper, app.MapsKeeper)
// 	wasmOpts := append(customEncoderOptions, customQueryOptions...)
// 	availableCapabilities := wasmkeeper.BuiltInCapabilities()
// 	availableCapabilities = append(availableCapabilities, "bitbadges")

// 	storeService := runtime.NewKVStoreService(wasmKey)

// 	app.WasmKeeper = wasmkeeper.NewKeeper(
// 		app.appCodec,
// 		storeService,
// 		app.AccountKeeper,
// 		app.BankKeeper,
// 		app.StakingKeeper,
// 		distrkeeper.NewQuerier(app.DistrKeeper),
// 		app.IBCFeeKeeper, // ISC4 Wrapper: fee IBC middleware
// 		app.IBCKeeper.ChannelKeeper,
// 		app.IBCKeeper.PortKeeper,
// 		app.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName),
// 		app.TransferKeeper,
// 		app.MsgServiceRouter(),
// 		app.GRPCQueryRouter(),
// 		wasmDir,
// 		wasmConfig,
// 		availableCapabilities,
// 		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
// 		wasmOpts...,
// 	)

// 	app.WasmxKeeper = wasmxkeeper.NewKeeper(
// 		app.appCodec,
// 		storeService,
// 		app.Logger(),
// 		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
// 		func() *ibckeeper.Keeper { return app.IBCKeeper },
// 		app.CapabilityKeeper.ScopeToModule(wasmxmoduletypes.ModuleName),
// 		app.WasmKeeper,
// 	)

// 	if err := app.RegisterModules(
// 		wasm.NewAppModule(
// 			app.appCodec,
// 			&app.WasmKeeper,
// 			app.StakingKeeper,
// 			app.AccountKeeper,
// 			app.BankKeeper,
// 			app.MsgServiceRouter(),
// 			app.GetSubspace(wasmtypes.ModuleName)),
// 		wasmx.NewAppModule(
// 			app.appCodec,
// 			app.WasmxKeeper,
// 			app.AccountKeeper,
// 			app.BankKeeper,
// 			app.WasmKeeper,
// 		),
// 	); err != nil {
// 		return nil, err
// 	}

// 	return storeService, nil
// }

// // RegisterIBC Since the IBC modules don't support dependency injection,
// // we need to manually register the modules on the client side.
// // This needs to be removed after IBC supports App Wiring.
// func RegisterWasm(registry cdctypes.InterfaceRegistry) map[string]appmodule.AppModule {
// 	modules := map[string]appmodule.AppModule{
// 		wasmtypes.ModuleName:        wasm.AppModule{},
// 		wasmxmoduletypes.ModuleName: wasmx.AppModule{},
// 	}

// 	for name, m := range modules {
// 		module.CoreAppModuleBasicAdaptor(name, m).RegisterInterfaces(registry)
// 	}

// 	return modules
// }
