package ante

import (
	"fmt"
	"runtime/debug"

	"bitbadgeschain/x/anchor/types"

	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	storetypes "cosmossdk.io/store/types"

	"bitbadgeschain/chain-handlers/ethereum/crypto/ethsecp256k1"

	ed25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
)

const (
	secp256k1VerifyCost uint64 = 21000
	ed25519VerifyCost   uint64 = 21000
)

// NewAnteHandler returns an ante handler responsible for attempting to route an
// Ethereum or SDK transaction to an internal ante handler for performing
// transaction-level processing (e.g. fee payment, signature verification) before
// being passed onto it's respective handler.
func NewAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		defer Recover(ctx.Logger(), &err)

		txWithExtensions, ok := tx.(authante.HasExtensionOptionsTx)
		if ok {
			opts := txWithExtensions.GetExtensionOptions()
			if len(opts) > 0 {
				switch typeURL := opts[0].GetTypeUrl(); typeURL {
				case "/ethereum.ExtensionOptionsWeb3Tx":
					// handle as Cosmos SDK tx, except signature is checked for EIP712 representation
					anteHandler = newCosmosAnteHandlerEip712(options, "Ethereum")
				case "/solana.ExtensionOptionsWeb3TxSolana":
					// handle as Cosmos SDK tx, except signature is checked for JSON representation intended to be signed w/ solana (Phantom wallet)
					anteHandler = newCosmosAnteHandlerEip712(options, "Solana")
				case "/bitcoin.ExtensionOptionsWeb3TxBitcoin":
					// handle as Cosmos SDK tx, except signature is checked for JSON representation intended to be signed w/ bitcoin (Ledger wallet)
					anteHandler = newCosmosAnteHandlerEip712(options, "Bitcoin")
				default:
					return ctx, sdkerrors.Wrapf(
						types.ErrUnknownRequest,
						"rejecting tx with unsupported extension option: %s", typeURL,
					)
				}

				return anteHandler(ctx, tx, sim)
			}
		}

		// handle as totally normal Cosmos SDK tx
		switch tx.(type) {
		case sdk.Tx:
			anteHandler = newCosmosAnteHandler(options)
		default:
			return ctx, sdkerrors.Wrapf(types.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}

		return anteHandler(ctx, tx, sim)
	}
}

func Recover(logger log.Logger, err *error) {
	if r := recover(); r != nil {
		*err = sdkerrors.Wrapf(sdkerrors.ErrPanic, "%v", r)

		if e, ok := r.(error); ok {
			logger.Error(
				"ante handler panicked",
				"error", e,
				"stack trace", string(debug.Stack()),
			)
		} else {
			logger.Error(
				"ante handler panicked",
				"recover", fmt.Sprintf("%v", r),
			)
		}
	}
}

var _ authante.SignatureVerificationGasConsumer = DefaultSigVerificationGasConsumer

// DefaultSigVerificationGasConsumer is the default implementation of SignatureVerificationGasConsumer. It consumes gas
// for signature verification based upon the public key type. The cost is fetched from the given params and is matched
// by the concrete type.
func DefaultSigVerificationGasConsumer(
	meter storetypes.GasMeter, sig signing.SignatureV2, params authtypes.Params,
) error {
	// support for ethereum ECDSA secp256k1 keys
	_, ok := sig.PubKey.(*ethsecp256k1.PubKey)
	if ok {
		meter.ConsumeGas(secp256k1VerifyCost, "ante verify: eth_secp256k1")
		return nil
	}

	_, ok = sig.PubKey.(*ed25519.PubKey)
	if ok {
		meter.ConsumeGas(ed25519VerifyCost, "ante verify: solana ed25519")
		return nil
	}

	return authante.DefaultSigVerificationGasConsumer(meter, sig, params)
}
