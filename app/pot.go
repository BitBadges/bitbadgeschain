package app

import (
	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	potkeeper "github.com/bitbadges/bitbadgeschain/x/pot/keeper"
	pot "github.com/bitbadges/bitbadgeschain/x/pot/module"
	pottypes "github.com/bitbadges/bitbadgeschain/x/pot/types"
)

// registerPotModule registers the x/pot module stores, keeper, and AppModule.
// This follows the same manual registration pattern as registerGammModules.
func (app *App) registerPotModule(_ servertypes.AppOptions) error {
	// Register store keys (no transient store needed — x/pot uses persistent KV only).
	storeKey := storetypes.NewKVStoreKey(pottypes.StoreKey)
	if err := app.RegisterStores(storeKey); err != nil {
		return err
	}

	// Create the tokenization adapter (wraps the concrete tokenization keeper).
	tokenizationAdapter := potkeeper.NewTokenizationKeeperAdapter(app.TokenizationKeeper)

	// Governance authority for MsgUpdateParams.
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	// Create the staking adapter (for chains using x/staking).
	// For PoA chains, replace this with potkeeper.NewPoAAdapter(poaKeeper, potKeeper).
	validatorSet := potkeeper.NewStakingAdapter(app.StakingKeeper, app.SlashingKeeper)

	// Create the pot keeper.
	app.PotKeeper = potkeeper.NewKeeper(
		app.appCodec,
		storeKey,
		app.Logger(),
		authority.String(),
		tokenizationAdapter,
		validatorSet,
	)

	// Create and register the AppModule.
	potModule := pot.NewAppModule(app.appCodec, app.PotKeeper)
	if err := app.RegisterModules(potModule); err != nil {
		return err
	}

	return nil
}

// RegisterPot registers the x/pot module on the client side so that
// its CLI commands (tx pot, query pot) and interface types are available
// in the binary. This mirrors the RegisterGamm / RegisterIBC pattern.
func RegisterPot(registry cdctypes.InterfaceRegistry) map[string]appmodule.AppModule {
	modules := map[string]appmodule.AppModule{
		pottypes.ModuleName: pot.AppModule{},
	}

	for name, m := range modules {
		module.CoreAppModuleBasicAdaptor(name, m).RegisterInterfaces(registry)
	}

	return modules
}
