package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// initializeEVMCoinInfo initializes EVM coin info during InitChain.
// This ensures denom metadata exists and EVM params are configured correctly
// so that InitEvmCoinInfo can succeed. Needed for local dev (ignite serve)
// where upgrade handlers don't run.
//
// Note: EVM chain ID is set in evm.go during keeper initialization on every
// app startup based on the Cosmos chain ID from appOpts.
func (app *App) initializeEVMCoinInfo(ctx sdk.Context) error {
	// Set denom metadata if it doesn't exist
	if _, found := app.BankKeeper.GetDenomMetaData(ctx, "ubadge"); !found {
		app.BankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
			Description: "The native token of BitBadges Chain",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "ubadge", Exponent: 0},
				{Denom: "badge", Exponent: 9},
			},
			Base:    "ubadge",
			Display: "badge",
			Name:    "Badge",
			Symbol:  "BADGE",
		})
	}

	// Ensure EVM params are configured correctly
	evmParams := app.EVMKeeper.GetParams(ctx)
	if evmParams.EvmDenom != "ubadge" {
		evmParams.EvmDenom = "ubadge"
	}
	if evmParams.ExtendedDenomOptions == nil {
		evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: "abadge"}
	}
	if err := app.EVMKeeper.SetParams(ctx, evmParams); err != nil {
		return err
	}

	// Initialize EvmCoinInfo (may already be initialized)
	if err := app.EVMKeeper.SetEvmCoinInfo(ctx, evmtypes.EvmCoinInfo{
		Denom:         "ubadge",
		ExtendedDenom: "abadge",
		DisplayDenom:  "BADGE",
		Decimals:      9,
	}); err != nil {
		ctx.Logger().Info("EVM coin info initialization skipped", "error", err)
	}

	return nil
}

