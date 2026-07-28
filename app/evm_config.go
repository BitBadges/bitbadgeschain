//go:build !test
// +build !test

package app

import (
	evmmodule "github.com/cosmos/evm/x/vm"
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
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
	evmmodule.SetGlobalConfigVariables(evmtypes.EvmCoinInfo{
		Denom:         "ubadge", // Base 9-decimal denomination
		ExtendedDenom: "abadge", // Extended 18-decimal denomination used by the EVM
		DisplayDenom:  "BADGE",
		Decimals:      9, // Base decimals
	})
	return keeper
}
