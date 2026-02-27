package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
)

// ensureBankDenomMetadata sets the ubadge denom metadata in the bank keeper.
// This MUST be called BEFORE module InitGenesis runs because the EVM module's
// InitGenesis requires bank denom metadata to be present.
//
// This is critical for ibc-go testing which creates bank genesis with empty metadata.
func (app *App) ensureBankDenomMetadata(ctx sdk.Context) {
	// Check if metadata already exists
	if _, found := app.BankKeeper.GetDenomMetaData(ctx, "ubadge"); found {
		return
	}

	// Set the ubadge denom metadata
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
	// CRITICAL: This must be set correctly or the ante handler will reject transactions with ubadge fees
	coinInfo := evmtypes.EvmCoinInfo{
		Denom:         "ubadge",
		ExtendedDenom: "abadge",
		DisplayDenom:  "BADGE",
		Decimals:      9,
	}
	if err := app.EVMKeeper.SetEvmCoinInfo(ctx, coinInfo); err != nil {
		ctx.Logger().Info("EVM coin info initialization skipped", "error", err)
	}

	// CRITICAL: Verify coin info is set correctly - if it's still "aatom", force set it
	// This is a safeguard to ensure coin info is always "ubadge" even if something resets it
	currentCoinInfo := app.EVMKeeper.GetEvmCoinInfo(ctx)
	if currentCoinInfo.Denom != "ubadge" {
		ctx.Logger().Warn("CRITICAL: EVM coin info denom is not ubadge, fixing",
			"current", currentCoinInfo.Denom,
			"expected", "ubadge",
			"evmParamsDenom", evmParams.EvmDenom)
		// Force set it - this is critical for transaction fees to work
		// Try multiple times to ensure it sticks
		for i := 0; i < 3; i++ {
			if err := app.EVMKeeper.SetEvmCoinInfo(ctx, coinInfo); err != nil {
				ctx.Logger().Error("CRITICAL: Failed to fix EVM coin info denom",
					"error", err,
					"attempt", i+1,
					"evmParamsDenom", evmParams.EvmDenom)
				if i == 2 {
					// Last attempt failed - this is critical
					return fmt.Errorf("CRITICAL: failed to set EVM coin info to ubadge after 3 attempts: current denom is %s, params denom is %s, error: %w",
						currentCoinInfo.Denom, evmParams.EvmDenom, err)
				}
				// Re-read params in case they changed
				evmParams = app.EVMKeeper.GetParams(ctx)
				continue
			}
			// Verify it was set
			currentCoinInfo = app.EVMKeeper.GetEvmCoinInfo(ctx)
			if currentCoinInfo.Denom == "ubadge" {
				ctx.Logger().Info("CRITICAL: EVM coin info denom fixed to ubadge", "attempt", i+1)
				break
			}
			if i == 2 {
				return fmt.Errorf("CRITICAL: EVM coin info still incorrect after 3 set attempts: got %s, expected ubadge, params denom is %s",
					currentCoinInfo.Denom, evmParams.EvmDenom)
			}
		}
	} else {
		ctx.Logger().Info("EVM coin info verified correctly", "denom", currentCoinInfo.Denom)
	}

	return nil
}
