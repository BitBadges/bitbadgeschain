//go:build test
// +build test

package app

import (
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
)

// configureEVMKeeper applies EVM keeper configuration for test builds.
// In tests, we don't call WithDefaultEvmCoinInfo because it conflicts with
// the testing-specific global state in cosmos/evm. The genesis initialization
// will properly configure the EVM coin info through InitGenesis.
func configureEVMKeeper(keeper *evmkeeper.Keeper) *evmkeeper.Keeper {
	// Don't set default coin info in tests - it conflicts with the
	// EVM module's test build which uses a single global variable for both
	// default and configured coin info.
	return keeper
}
