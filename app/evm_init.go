package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// initializeEVMCoinInfo initializes EVM coin info for local development.
// This ensures that denom metadata exists and EVM params are configured correctly
// so that InitEvmCoinInfo can succeed. This is needed for local dev (ignite serve)
// where upgrade handlers don't run automatically.
func (app *App) initializeEVMCoinInfo(ctx sdk.Context) error {
	// Check if denom metadata for "ubadge" exists
	_, found := app.BankKeeper.GetDenomMetaData(ctx, "ubadge")
	if !found {
		// Set denom metadata if it doesn't exist
		denomMetadata := banktypes.Metadata{
			Description: "The native token of BitBadges Chain",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "ubadge",
					Exponent: 0,
					Aliases:  nil,
				},
				{
					Denom:    "badge",
					Exponent: 9,
					Aliases:  nil,
				},
			},
			Base:    "ubadge",
			Display: "badge",
			Name:    "Badge",
			Symbol:  "BADGE",
		}
		app.BankKeeper.SetDenomMetaData(ctx, denomMetadata)
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

	// Initialize EvmCoinInfo if it hasn't been initialized yet
	// This will only succeed if denom metadata exists and params are set correctly
	if err := app.EVMKeeper.InitEvmCoinInfo(ctx); err != nil {
		// Log the error but don't fail - this might already be initialized
		ctx.Logger().Info("EVM coin info initialization skipped (may already be initialized)", "error", err)
	}

	return nil
}

