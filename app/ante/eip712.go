package ante

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/tidwall/gjson"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/crypto/ethsecp256k1"
	"github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/ethereum/eip712"
	ethereumtypes "github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/types"
	ethereum "github.com/bitbadges/bitbadgeschain/chain-handlers/ethereum/utils"
	solanatypes "github.com/bitbadges/bitbadgeschain/chain-handlers/solana/types"

	bitcointypes "github.com/bitbadges/bitbadgeschain/chain-handlers/bitcoin/types"

	solana "github.com/bitbadges/bitbadgeschain/chain-handlers/solana/utils"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var ethereumCodec codec.ProtoCodecMarshaler
var solanaCodec codec.ProtoCodecMarshaler

func init() {
	registry := codectypes.NewInterfaceRegistry()
	ethereum.RegisterInterfaces(registry)
	solana.RegisterInterfaces(registry)
	ethereumCodec = codec.NewProtoCodec(registry)
	solanaCodec = codec.NewProtoCodec(registry)
}

// Eip712SigVerificationDecorator Verify all signatures for a tx and return an error if any are invalid. Note,
// the Eip712SigVerificationDecorator decorator will not get executed on ReCheck.
//
// CONTRACT: Pubkeys are set in context for all signers before this decorator runs
// CONTRACT: Tx must implement SigVerifiableTx interface
type Eip712SigVerificationDecorator struct {
	ak              ante.AccountKeeper
	signModeHandler authsigning.SignModeHandler
	chain           string
}

// NewEip712SigVerificationDecorator creates a new Eip712SigVerificationDecorator
func NewEip712SigVerificationDecorator(ak ante.AccountKeeper, signModeHandler authsigning.SignModeHandler, chain string) Eip712SigVerificationDecorator {
	return Eip712SigVerificationDecorator{
		ak:              ak,
		signModeHandler: signModeHandler,
		chain:           chain,
	}
}

// AnteHandle handles validation of EIP712 signed cosmos txs.
// it is not run on RecheckTx
func (svd Eip712SigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// no need to verify signatures on recheck tx
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(types.ErrInvalidType, "tx %T doesn't implement authsigning.SigVerifiableTx", tx)
	}

	authSignTx, ok := tx.(authsigning.Tx)
	if !ok {
		return ctx, sdkerrors.Wrapf(types.ErrInvalidType, "tx %T doesn't implement the authsigning.Tx interface", tx)
	}

	// stdSigs contains the sequence number, account number, and signatures.
	// When simulating, this would just be a 0-length slice.
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return ctx, err
	}

	signerAddrs := sigTx.GetSigners()

	// EIP712 allows just one signature
	if len(sigs) != 1 {
		return ctx, sdkerrors.Wrapf(types.ErrUnauthorized, "invalid number of signers (%d);  EIP712 signatures allows just one signature", len(sigs))
	}

	// check that signer length and signature length are the same
	if len(sigs) != len(signerAddrs) {
		return ctx, sdkerrors.Wrapf(types.ErrUnauthorized, "invalid number of signer;  expected: %d, got %d", len(signerAddrs), len(sigs))
	}

	// EIP712 has just one signature, avoid looping here and only read index 0
	i := 0
	sig := sigs[i]

	acc, err := authante.GetSignerAcc(ctx, svd.ak, signerAddrs[i])
	if err != nil {
		return ctx, err
	}

	// retrieve pubkey
	pubKey := acc.GetPubKey()
	if !simulate && pubKey == nil {
		return ctx, sdkerrors.Wrap(types.ErrInvalidPubKey, "pubkey on account is not set")
	}

	// Check account sequence number.
	if sig.Sequence != acc.GetSequence() {
		return ctx, sdkerrors.Wrapf(
			types.ErrWrongSequence,
			"account sequence mismatch, expected %d, got %d", acc.GetSequence(), sig.Sequence,
		)
	}

	// retrieve signer data
	genesis := ctx.BlockHeight() == 0
	chainID := ctx.ChainID()

	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	signerData := authsigning.SignerData{
		ChainID:       chainID,
		AccountNumber: accNum,
		Sequence:      acc.GetSequence(),
	}

	if simulate {
		return next(ctx, tx, simulate)
	}

	if err := VerifySignature(pubKey, signerData, sig.Data, svd.signModeHandler, authSignTx, svd.chain); err != nil {
		errMsg := fmt.Errorf("signature verification failed; please verify account number (%d) and chain-id (%s): %w", accNum, chainID, err)
		return ctx, sdkerrors.Wrap(types.ErrUnauthorized, errMsg.Error())
	}

	return next(ctx, tx, simulate)
}

// VerifySignature verifies a transaction signature contained in SignatureData abstracting over different signing modes
// and single vs multi-signatures.
func VerifySignature(
	pubKey cryptotypes.PubKey,
	signerData authsigning.SignerData,
	sigData signing.SignatureData,
	_ authsigning.SignModeHandler,
	tx authsigning.Tx,
	chain string,
) error {
	switch data := sigData.(type) {
	case *signing.SingleSignatureData:
		if data.SignMode != signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON {
			return sdkerrors.Wrapf(types.ErrNotSupported, "unexpected SignatureData %T: wrong SignMode", sigData)
		}

		// Note: this prevents the user from sending trash data in the signature field
		if len(data.Signature) != 0 {
			return sdkerrors.Wrap(types.ErrTooManySignatures, "invalid signature value; EIP712 must have the cosmos transaction signature empty")
		}

		// @contract: this code is reached only when Msg has Web3Tx extension (so this custom Ante handler flow),
		// and the signature is SIGN_MODE_LEGACY_AMINO_JSON which is supported for EIP712 for now

		msgs := tx.GetMsgs()
		if len(msgs) == 0 {
			return sdkerrors.Wrap(types.ErrNoSignatures, "tx doesn't contain any msgs to verify signature")
		}

		txBytes := legacytx.StdSignBytes(
			signerData.ChainID,
			signerData.AccountNumber,
			signerData.Sequence,
			tx.GetTimeoutHeight(),
			legacytx.StdFee{
				Amount: tx.GetFee(),
				Gas:    tx.GetGas(),
			},
			msgs, tx.GetMemo(), tx.GetTip(),
		)

		signerChainID, err := ethereum.ParseChainID(signerData.ChainID)
		if err != nil {
			return sdkerrors.Wrapf(err, "failed to parse chain-id: %s", signerData.ChainID)
		}

		txWithExtensions, ok := tx.(authante.HasExtensionOptionsTx)
		if !ok {
			return sdkerrors.Wrap(types.ErrUnknownExtensionOptions, "tx doesnt contain any extensions")
		}

		opts := txWithExtensions.GetExtensionOptions()
		if len(opts) != 1 {
			return sdkerrors.Wrap(types.ErrUnknownExtensionOptions, "tx doesnt contain expected amount of extension options")
		}

		typedDataChainID := uint64(0)
		feePayerExt := ""
		feePayerSig := []byte{}

		//Get details from the extension option
		//TODO: Is feePayer really necessary in the extension? what is it used for?
		extOptEthereum, ok := opts[0].GetCachedValue().(*ethereumtypes.ExtensionOptionsWeb3Tx)
		if !ok && chain == "Ethereum" {
			return sdkerrors.Wrap(types.ErrUnknownExtensionOptions, "unknown extension option")
		} else if ok && chain == "Ethereum" {
			typedDataChainID = extOptEthereum.TypedDataChainID
			feePayerExt = extOptEthereum.FeePayer
			feePayerSig = extOptEthereum.FeePayerSig
		}

		extOptSolana, ok := opts[0].GetCachedValue().(*solanatypes.ExtensionOptionsWeb3TxSolana)
		if !ok && chain == "Solana" {
			return sdkerrors.Wrap(types.ErrUnknownExtensionOptions, "unknown extension option")
		} else if ok && chain == "Solana" {
			typedDataChainID = extOptSolana.TypedDataChainID
			feePayerExt = extOptSolana.FeePayer
			feePayerSig = extOptSolana.FeePayerSig
		}

		extOptBitcoin, ok := opts[0].GetCachedValue().(*bitcointypes.ExtensionOptionsWeb3TxBitcoin)
		if !ok && chain == "Bitcoin" {
			return sdkerrors.Wrap(types.ErrUnknownExtensionOptions, "unknown extension option")
		} else if ok && chain == "Bitcoin" {
			typedDataChainID = extOptBitcoin.TypedDataChainID
			feePayerExt = extOptBitcoin.FeePayer
			feePayerSig = extOptBitcoin.FeePayerSig
		}

		if typedDataChainID != signerChainID.Uint64() {
			return sdkerrors.Wrap(types.ErrInvalidChainID, "invalid chain-id")
		}

		if len(feePayerExt) == 0 {
			return sdkerrors.Wrap(types.ErrUnknownExtensionOptions, "no feePayer on ExtensionOptionsWeb3Tx")
		}

		feePayer, err := sdk.AccAddressFromBech32(feePayerExt)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to parse feePayer from ExtensionOptionsWeb3Tx")
		}

		recoveredFeePayerAcc := sdk.AccAddress(pubKey.Address().Bytes())
		if !recoveredFeePayerAcc.Equals(feePayer) {
			return sdkerrors.Wrapf(types.ErrorInvalidSigner, "failed to match fee payer in extension to the expected signer %s", recoveredFeePayerAcc)
		}

		//This uses the new EIP712 wrapper, not legacy
		//This also accounts for multiple msgs, not just one
		typedData, err := eip712.WrapTxToTypedData(typedDataChainID, txBytes, chain)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to pack tx data in EIP712 object")
		}

		//Within the wrapping proess, we added types for values that may not be in the message payload
		//but are expectd by the EIP712 typedData.  We need to add the empty values ("", false, 0) to the
		//message payload
		typedData, err = eip712.NormalizeEIP712TypedData(typedData)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to normalize EIP712 typed data for badges module")
		}

		//If chain is Solana, we need to use the Solana way to verify the signature (alphabetically sorted JSON keys)
		//Else, we use EIP712 typed signatures for Ethereum
		if chain == "Solana" || chain == "Bitcoin" {
			//We generate the Solana message payload using the code from generatin EIP712
			//We only use the message field (no types or domain)
			//Then, we sort the message field by alphabetizing the JSON keys

			//Creates the EIP712 message payload
			eip712Message := typedData.Message //map[string]interface{}
			// Marshal the map to JSON
			jsonData, err := json.Marshal(eip712Message)
			if err != nil {
				return sdkerrors.Wrap(err, "failed to marshal json")
			}

			// convert to gson and alphabetize the JSON
			jsonStr := gjson.ParseBytes(jsonData).String()
			sortedBytes, err := sdk.SortJSON([]byte(jsonStr))
			if err != nil {
				return sdkerrors.Wrap(err, "failed to sort JSON")
			}

			//Match address to pubkey to make sure it's equivalent and no random address is used
			//This is used for indexing purposes (to be able to map Solana addresses to Cosmos addresses,
			//you need to know the Solana address bc it takes a hash to convert, so you can't go the opposite way)
			//
			//Doesn't have any on-chain significance

			if chain == "Solana" {
				solanaAddress := extOptSolana.SolAddress
				addr := base58.Encode(pubKey.Bytes())
				if addr != solanaAddress {
					return sdkerrors.Wrap(types.ErrUnknownExtensionOptions, "provided solana address in extension does not match signer pubkey")
				}

				//Verify signature w/ ed25519
				valid := ed25519.Verify(pubKey.Bytes(), sortedBytes, feePayerSig)
				if !valid {
					return sdkerrors.Wrapf(types.ErrorInvalidSigner, "failed to verify delegated fee payer %s signature %s", recoveredFeePayerAcc, jsonStr)
				}
			} else if chain == "Bitcoin" {
				//verify bitcoin bip322 signature

				//TODO:
				
				return sdkerrors.Wrap(types.ErrorInvalidSigner, "unable to verify signer signature of Bitcoin signature: not implemented")
				
			}

			return nil
		} else {

			// return sdkerrors.Wrapf(types.ErrInvalidChainID, "%s", typedData)

			sigHash, _, err := apitypes.TypedDataAndHash(typedData)
			if err != nil {
				return sdkerrors.Wrapf(err, "failed to compute typed data hash")
			}

			if len(feePayerSig) != ethcrypto.SignatureLength {
				return sdkerrors.Wrap(types.ErrorInvalidSigner, "signature length doesn't match typical [R||S||V] signature 65 bytes")
			}

			// Remove the recovery offset if needed (ie. Metamask eip712 signature)
			if feePayerSig[ethcrypto.RecoveryIDOffset] == 27 || feePayerSig[ethcrypto.RecoveryIDOffset] == 28 {
				feePayerSig[ethcrypto.RecoveryIDOffset] -= 27
			}

			feePayerPubkey, err := secp256k1.RecoverPubkey(sigHash, feePayerSig)
			if err != nil {
				return sdkerrors.Wrap(err, "failed to recover delegated fee payer from sig")
			}

			ecPubKey, err := ethcrypto.UnmarshalPubkey(feePayerPubkey)
			if err != nil {
				return sdkerrors.Wrap(err, "failed to unmarshal recovered fee payer pubkey")
			}

			pk := &ethsecp256k1.PubKey{
				Key: ethcrypto.CompressPubkey(ecPubKey),
			}

			if !pubKey.Equals(pk) {
				return sdkerrors.Wrapf(types.ErrInvalidPubKey, "feePayer pubkey %s is different from transaction pubkey %s", pubKey, pk)
			}

			// VerifySignature of ethsecp256k1 accepts 64 byte signature [R||S]
			// WARNING! Under NO CIRCUMSTANCES try to use pubKey.VerifySignature there
			if !secp256k1.VerifySignature(pubKey.Bytes(), sigHash, feePayerSig[:len(feePayerSig)-1]) {
				return sdkerrors.Wrap(types.ErrorInvalidSigner, "unable to verify signer signature of EIP712 typed data")
			}
		}

		return nil
	default:
		return sdkerrors.Wrapf(types.ErrTooManySignatures, "unexpected SignatureData %T", sigData)
	}
}
