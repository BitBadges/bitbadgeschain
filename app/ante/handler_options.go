package ante

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	twofakeeper "github.com/bitbadges/bitbadgeschain/x/twofa/keeper"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	txsigning "cosmossdk.io/x/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ibcante "github.com/cosmos/ibc-go/v10/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	storetypes "cosmossdk.io/store/types"

	circuitante "cosmossdk.io/x/circuit/ante"
	circuitkeeper "cosmossdk.io/x/circuit/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	corestoretypes "cosmossdk.io/core/store"
)

// NewAnteHandler returns an ante handler responsible for performing
// transaction-level processing (e.g. fee payment, signature verification) before
// being passed onto it's respective handler.
//
// Decorators are executed in the order listed below, combining standard SDK decorators
// with custom decorators for Circuit Breaker, WASM, and IBC functionality.
func NewAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		// Setup and context
		ante.NewSetUpContextDecorator(),

		// Custom: Circuit Breaker (must be early to protect against problematic txs)
		circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),

		// Custom: WASM decorators (after setup context to enforce limits early)
		wasmkeeper.NewLimitSimulationGasDecorator(options.WasmConfig.SimulationGasLimit),
		wasmkeeper.NewCountTXDecorator(options.TXCounterStoreService),
		wasmkeeper.NewGasRegisterDecorator(options.WasmKeeper.GetGasRegister()),

		// Standard SDK decorators
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

		// Custom: Global 2FA decorator (applies to ALL transaction types, not just BitBadges)
		// This provides defense in depth by requiring badge ownership as a second factor
		// for any transaction from users who have set 2FA requirements.
		// Must be after signature verification to ensure signers are known.
		twofakeeper.NewTwoFADecorator(options.TwoFAKeeper),

		// Custom: IBC decorator (at the end)
		ibcante.NewRedundantRelayDecorator(options.IBCKeeper),
	)
}

// TODO: We can play around with some more native ones
// circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),
// ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
// ante.NewUnorderedTxDecorator(unorderedtx.DefaultMaxUnOrderedTTL, options.TxManager, options.Environment, ante.DefaultSha256Cost),

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper, Circuit Breaker keeper, and WASM keepers.
type HandlerOptions struct {
	AccountKeeper         ante.AccountKeeper
	BankKeeper            authtypes.BankKeeper
	IBCKeeper             *ibckeeper.Keeper
	FeegrantKeeper        ante.FeegrantKeeper
	SignModeHandler       *txsigning.HandlerMap
	SigGasConsumer        func(meter storetypes.GasMeter, sig signing.SignatureV2, params authtypes.Params) error // defaults to authante.DefaultSigVerificationGasConsumer
	TxFeeChecker          ante.TxFeeChecker
	CircuitKeeper         *circuitkeeper.Keeper
	WasmConfig            wasmtypes.NodeConfig
	WasmKeeper            *wasmkeeper.Keeper
	TXCounterStoreService corestoretypes.KVStoreService
	TwoFAKeeper           twofakeeper.Keeper
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
	// WasmConfig is now a value type, not a pointer - validation not needed
	if options.TXCounterStoreService == nil {
		return sdkerrors.Wrap(types.ErrLogic, "wasm store service is required for ante builder")
	}
	return nil
}
