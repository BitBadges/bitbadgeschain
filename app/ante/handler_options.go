package ante

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	txsigning "cosmossdk.io/x/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ibcante "github.com/cosmos/ibc-go/v8/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	storetypes "cosmossdk.io/store/types"

	circuitante "cosmossdk.io/x/circuit/ante"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"

	// wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	// wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"

	corestoretypes "cosmossdk.io/core/store"
)

//TODO: We can play around with some more native ones
// circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),
// ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
// ante.NewUnorderedTxDecorator(unorderedtx.DefaultMaxUnOrderedTTL, options.TxManager, options.Environment, ante.DefaultSha256Cost),

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper, EVM Keeper and Fee Market Keeper.
type HandlerOptions struct {
	AccountKeeper   ante.AccountKeeper
	BankKeeper      authtypes.BankKeeper
	IBCKeeper       *ibckeeper.Keeper
	FeegrantKeeper  ante.FeegrantKeeper
	SignModeHandler *txsigning.HandlerMap
	SigGasConsumer  func(meter storetypes.GasMeter, sig signing.SignatureV2, params authtypes.Params) error
	TxFeeChecker    ante.TxFeeChecker
	CircuitKeeper   *circuitkeeper.Keeper
	// WasmConfig            *wasmTypes.WasmConfig
	// WasmKeeper            *wasmkeeper.Keeper
	TXCounterStoreService corestoretypes.KVStoreService
}

func (options HandlerOptions) Validate() error {
	if options.AccountKeeper == nil {
		return sdkerrors.Wrap(types.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return sdkerrors.Wrap(types.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return sdkerrors.Wrap(types.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.CircuitKeeper == nil {
		return sdkerrors.Wrap(types.ErrLogic, "circuit keeper is required for ante builder")
	}
	// if options.WasmConfig == nil {
	// 	return sdkerrors.Wrap(types.ErrLogic, "wasm config is required for ante builder")
	// }
	// if options.TXCounterStoreService == nil {
	// 	return sdkerrors.Wrap(types.ErrLogic, "wasm store service is required for ante builder")
	// }
	return nil
}

func newCosmosAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(),
		circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),

		// ante.NewRejectExtensionOptionsDecorator(),
		// ante.NewMempoolFeeDecorator(),

		// wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		// wasmkeeper.NewCountTXDecorator(options.TXCounterStoreService),
		// wasmkeeper.NewGasRegisterDecorator(options.WasmKeeper.GetGasRegister()),

		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	)
}

func newCosmosAnteHandlerEip712(options HandlerOptions, chain string) sdk.AnteHandler {
	if chain != "Ethereum" && chain != "Solana" && chain != "Bitcoin" {
		panic("chain must be either Ethereum or Solana")
	}

	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(),
		circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),
		// NOTE: extensions option decorator removed
		// ante.NewRejectExtensionOptionsDecorator(),

		// wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		// wasmkeeper.NewCountTXDecorator(options.TXCounterStoreService),
		// wasmkeeper.NewGasRegisterDecorator(options.WasmKeeper.GetGasRegister()),

		// ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		// SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewSetPubKeyDecorator(options.AccountKeeper),
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),

		NewCustomSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler, chain),

		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	)
}
