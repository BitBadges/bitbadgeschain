//go:build !test
// +build !test

package app

import (
	evmkeeper "github.com/cosmos/evm/x/vm/keeper"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// configureEVMKeeper applies EVM keeper configuration for production builds.
// In production, we set a default EVM coin info as a fallback for RPC calls
// that happen before PreBlock initializes the store-based config.
func configureEVMKeeper(keeper *evmkeeper.Keeper) *evmkeeper.Keeper {
	return keeper.WithDefaultEvmCoinInfo(evmtypes.EvmCoinInfo{
		Denom:         "ubadge",
		ExtendedDenom: "ubadge",
		DisplayDenom:  "BADGE",
		Decimals:      9,
	})
}
