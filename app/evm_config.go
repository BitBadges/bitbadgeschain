//go:build !test
// +build !test

package app

import (
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
)

// configureEVMKeeper applies EVM configuration for production builds.
//
// cosmos/evm v0.7 removed Keeper.WithDefaultEvmCoinInfo. Coin info is now
// process-global state seeded through the EVM configurator, so it is set here
// rather than by decorating the keeper. The keeper is returned unchanged so the
// call site in registerEVMModules is unaffected.
//
// The chain stays 9-decimal (ubadge), with abadge as the 18-decimal extended
// denom the EVM operates in. See ticket #0467 for the open question of
// migrating to 18 decimals, which upstream now treats as the only fully
// supported configuration.
func configureEVMKeeper(keeper *evmkeeper.Keeper) *evmkeeper.Keeper {
	// Deliberately does NOT call evmmodule.SetGlobalConfigVariables.
	//
	// That function registers the extended EIP activators into a process-global
	// registry and is not idempotent — a second call panics with
	// "duplicate activation: 2 is already present in [...]".
	//
	// Upstream x/vm already calls it exactly once, guarded by the module's own
	// sync.Once, from InitGenesis (fresh chain) / PreBlock / HydrateGlobals
	// (existing chain). Calling it here as well meant the chain panicked at the
	// first block on a new chain and at startup on an upgraded one. Seeding the
	// globals is upstream's job; ours is to make sure genesis carries the right
	// denominations for it to read (see app/evm_genesis.go).
	return keeper
}
