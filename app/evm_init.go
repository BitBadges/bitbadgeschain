package app

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
)

// ensureBankDenomMetadata sets the base-denom metadata in the bank keeper.
// This MUST be called BEFORE module InitGenesis runs because the EVM module's
// InitGenesis requires bank denom metadata to be present.
//
// This is critical for ibc-go testing which creates bank genesis with empty metadata.
func (app *App) ensureBankDenomMetadata(ctx sdk.Context) {
	// Check if metadata already exists
	if _, found := app.BankKeeper.GetDenomMetaData(ctx, appparams.BaseCoinUnit); found {
		return
	}

	// Set the base denom metadata
	app.BankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
		Description: "The native token of BitBadges Chain",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: appparams.BaseCoinUnit, Exponent: 0},
			{Denom: appparams.DisplayCoinUnit, Exponent: appparams.BaseCoinDecimals},
		},
		Base:    appparams.BaseCoinUnit,
		Display: appparams.DisplayCoinUnit,
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
	if _, found := app.BankKeeper.GetDenomMetaData(ctx, appparams.BaseCoinUnit); !found {
		app.BankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
			Description: "The native token of BitBadges Chain",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: appparams.BaseCoinUnit, Exponent: 0},
				{Denom: appparams.DisplayCoinUnit, Exponent: appparams.BaseCoinDecimals},
			},
			Base:    appparams.BaseCoinUnit,
			Display: appparams.DisplayCoinUnit,
			Name:    "Badge",
			Symbol:  "BADGE",
		})
	}

	// Ensure EVM params are configured correctly
	evmParams := app.EVMKeeper.GetParams(ctx)
	if evmParams.EvmDenom != appparams.BaseCoinUnit {
		evmParams.EvmDenom = appparams.BaseCoinUnit
	}
	// Set unconditionally, not only when nil. Upstream's DefaultParams seeds
	// ExtendedDenomOptions with sdk.DefaultBondDenom ("stake"), so a nil-guard
	// leaves a wrong-but-non-nil value in place — and the EVM then reads
	// balances in a denom nobody holds, reporting every account as empty. That
	// is a silent failure: no error, just zeros.
	evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{
		ExtendedDenom: appparams.ExtendedCoinUnit,
	}
	if err := app.EVMKeeper.SetParams(ctx, evmParams); err != nil {
		return err
	}

	// Initialize EvmCoinInfo (may already be initialized)
	// CRITICAL: This must be set correctly or the ante handler will reject transactions with abadge fees
	coinInfo := evmtypes.EvmCoinInfo{
		Denom:         appparams.BaseCoinUnit,
		ExtendedDenom: appparams.ExtendedCoinUnit,
		DisplayDenom:  "BADGE",
		Decimals:      appparams.BaseCoinDecimals,
	}
	if err := app.EVMKeeper.SetEvmCoinInfo(ctx, coinInfo); err != nil {
		ctx.Logger().Info("EVM coin info initialization skipped", "error", err)
	}

	// CRITICAL: Verify coin info is set correctly - if it's still "aatom", force set it
	// This is a safeguard to ensure coin info is always "abadge" even if something resets it
	currentCoinInfo := app.EVMKeeper.GetEvmCoinInfo(ctx)
	if currentCoinInfo.Denom != appparams.BaseCoinUnit {
		ctx.Logger().Warn("CRITICAL: EVM coin info denom is not the base denom, fixing",
			"current", currentCoinInfo.Denom,
			"expected", appparams.BaseCoinUnit,
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
					return fmt.Errorf("CRITICAL: failed to set EVM coin info to base denom after 3 attempts: current denom is %s, params denom is %s, error: %w",
						currentCoinInfo.Denom, evmParams.EvmDenom, err)
				}
				// Re-read params in case they changed
				evmParams = app.EVMKeeper.GetParams(ctx)
				continue
			}
			// Verify it was set
			currentCoinInfo = app.EVMKeeper.GetEvmCoinInfo(ctx)
			if currentCoinInfo.Denom == appparams.BaseCoinUnit {
				ctx.Logger().Info("CRITICAL: EVM coin info denom fixed to base denom", "attempt", i+1)
				break
			}
			if i == 2 {
				return fmt.Errorf("CRITICAL: EVM coin info still incorrect after 3 set attempts: got %s, expected base denom, params denom is %s",
					currentCoinInfo.Denom, evmParams.EvmDenom)
			}
		}
	} else {
		ctx.Logger().Info("EVM coin info verified correctly", "denom", currentCoinInfo.Denom)
	}

	return nil
}
