package ante

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	protov2 "google.golang.org/protobuf/proto"

	signingv1beta1 "cosmossdk.io/api/cosmos/tx/signing/v1beta1"
	sdkerrors "cosmossdk.io/errors"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bech32cosmos "github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/tidwall/gjson"
	"github.com/unisat-wallet/libbrc20-indexer/utils/bip322"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"bitbadgeschain/chain-handlers/ethereum/crypto/ethsecp256k1"
	"bitbadgeschain/chain-handlers/ethereum/ethereum/eip712"
	ethereumtypes "bitbadgeschain/chain-handlers/ethereum/types"
	ethereum "bitbadgeschain/chain-handlers/ethereum/utils"
	solanatypes "bitbadgeschain/chain-handlers/solana/types"

	bitcointypes "bitbadgeschain/chain-handlers/bitcoin/types"

	solana "bitbadgeschain/chain-handlers/solana/utils"

	"bitbadgeschain/x/badges/types"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/storyicon/sigverify"

	bitcoin "bitbadgeschain/chain-handlers/bitcoin/utils"
	ethereumcodec "bitbadgeschain/chain-handlers/ethereum/crypto/codec"

	txsigning "cosmossdk.io/x/tx/signing"
)

var ethereumCodec codec.ProtoCodecMarshaler
var solanaCodec codec.ProtoCodecMarshaler

func init() {
	registry := codectypes.NewInterfaceRegistry()
	ethereumcodec.RegisterInterfaces(registry)
	bitcoin.RegisterInterfaces(registry)
	ethereum.RegisterInterfaces(registry)
	solana.RegisterInterfaces(registry)
	ethereumCodec = codec.NewProtoCodec(registry)
	solanaCodec = codec.NewProtoCodec(registry)
}

type CustomSigVerificationDecorator struct {
	ak              ante.AccountKeeper
	signModeHandler *txsigning.HandlerMap
	chain           string
}

func NewCustomSigVerificationDecorator(ak ante.AccountKeeper,
	signModeHandler *txsigning.HandlerMap, chain string) CustomSigVerificationDecorator {
	return CustomSigVerificationDecorator{
		ak:              ak,
		signModeHandler: signModeHandler,
		chain:           chain,
	}
}

func (svd CustomSigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
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

	signerAddrs, err := sigTx.GetSigners()
	if err != nil {
		return ctx, err
	}

	if len(sigs) != 1 {
		return ctx, sdkerrors.Wrapf(types.ErrUnauthorized, "invalid number of signers (%d);  EIP712 signatures allows just one signature", len(sigs))
	}

	// check that signer length and signature length are the same
	if len(sigs) != len(signerAddrs) {
		return ctx, sdkerrors.Wrapf(types.ErrUnauthorized, "invalid number of signer;  expected: %d, got %d", len(signerAddrs), len(sigs))
	}

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

	anyPk, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return ctx, err
	}

	txSigningData := txsigning.SignerData{
		ChainID:       chainID,
		AccountNumber: accNum,
		Sequence:      acc.GetSequence(),
		Address:       sdk.AccAddress(acc.GetAddress()).String(),
	}

	if pubKey != nil {
		txSigningData.PubKey = &anypb.Any{
			TypeUrl: anyPk.TypeUrl,
			Value:   anyPk.Value,
		}
	}

	if simulate {
		return next(ctx, tx, simulate)
	}

	if err := VerifySignature(pubKey, signerData, sig.Data, svd.signModeHandler, authSignTx, svd.chain, ctx, txSigningData); err != nil {
		errMsg := fmt.Errorf("signature verification failed; please verify account number (%d) and chain-id (%s): %w", accNum, chainID, err)
		return ctx, sdkerrors.Wrap(types.ErrUnauthorized, errMsg.Error())
	}

	return next(ctx, tx, simulate)
}

func deepCopyAnySlice(input []*anypb.Any) ([]*anypb.Any, error) {
	output := make([]*anypb.Any, len(input))
	for i, v := range input {
		if v != nil {
			// Use proto.Clone to deep copy each element
			cloned := proto.Clone(v).(*anypb.Any)
			output[i] = cloned
		}
	}
	return output, nil
}

// VerifySignature verifies a transaction signature contained in SignatureData abstracting over different signing modes
// and single vs multi-signatures.
func VerifySignature(
	pubKey cryptotypes.PubKey,
	signerData authsigning.SignerData,
	sigData signing.SignatureData,
	signModeMap *txsigning.HandlerMap,
	tx authsigning.Tx,
	chain string,
	ctx sdk.Context,
	txSigningData txsigning.SignerData,
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

		msgs := tx.GetMsgs()
		if len(msgs) == 0 {
			return sdkerrors.Wrap(types.ErrNoSignatures, "tx doesn't contain any msgs to verify signature")
		}

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

		adaptableTx := tx.(authsigning.V2AdaptableTx)
		txData := adaptableTx.GetSigningTxData()

		deepCopiedExtOpts, err := deepCopyAnySlice(txData.Body.ExtensionOptions)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to deep copy extension options")
		}

		deepCopiedNonCritExtOpts, err := deepCopyAnySlice(txData.Body.NonCriticalExtensionOptions)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to deep copy non-critical extension options")
		}

		//HACK: Little hacky. We only use txBytes for generating tx JSON (extensions not included)
		//			And, the GetSignBytes function throws for extensions with Amino. So, we remove and readd.
		//			IDK if all of this is necessary
		txData.Body.ExtensionOptions = nil
		txData.Body.NonCriticalExtensionOptions = nil
		bodyBz, err := protov2.Marshal(txData.Body)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to marshal tx body")
		}
		txData.BodyBytes = bodyBz
		txBytes, err := signModeMap.GetSignBytes(ctx, signingv1beta1.SignMode_SIGN_MODE_LEGACY_AMINO_JSON, txSigningData, txData)
		if err != nil {
			return sdkerrors.Wrap(err, "failed to get sign bytes")
		}

		txData.Body.ExtensionOptions = deepCopiedExtOpts
		txData.Body.NonCriticalExtensionOptions = deepCopiedNonCritExtOpts

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

		//We generate the Solana message payload using the code from generatin EIP712
		//We only use the message field (no types or domain)
		//Then, we sort the message field by alphabetizing the JSON keys

		//Creates the EIP712 message payload
		// Marshal the map to JSON
		eip712Message := typedData.Message //map[string]interface{}
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

		sha256JsonHash := sha256.Sum256(sortedBytes)
		jsonHashHexStr := hex.EncodeToString(sha256JsonHash[:])

		humanReadableStr := "This is a BitBadges transaction with the content hash: " + jsonHashHexStr

		//If chain is Solana, we need to use the Solana way to verify the signature (alphabetically sorted JSON keys)
		//Else, we use EIP712 typed signatures for Ethereum
		if chain == "Solana" || chain == "Bitcoin" {

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

				//Verify signature w/ ed25519 or the SHA256 hash of the sorted JSON (as a workaround for the 1000 byte limit)
				standardMsgSigValid := ed25519.Verify(pubKey.Bytes(), sortedBytes, feePayerSig)
				hashedMsgSigValid := ed25519.Verify(pubKey.Bytes(), []byte(jsonHashHexStr), feePayerSig)
				humanReadableMsgSigValid := ed25519.Verify(pubKey.Bytes(), []byte(humanReadableStr), feePayerSig)
				if !standardMsgSigValid && !hashedMsgSigValid && !humanReadableMsgSigValid {
					return sdkerrors.Wrapf(types.ErrorInvalidSigner, "failed to verify delegated fee payer %s signature %s %s", recoveredFeePayerAcc, jsonStr, jsonHashHexStr)
				}

			} else if chain == "Bitcoin" {
				//verify bitcoin bip322 signature
				message := string(sortedBytes)
				// signature := hex.EncodeToString(feePayerSig)
				cosmosAddress := feePayer.String()
				signature := base64.StdEncoding.EncodeToString(feePayerSig)
				_, base256Bytes, err := bech32cosmos.DecodeAndConvert(cosmosAddress)
				if err != nil {
					return sdkerrors.Wrap(err, "failed to decode and convert signer address")
				}

				base32Bytes, err := bech32.ConvertBits(base256Bytes, 8, 5, true)
				if err != nil {
					return sdkerrors.Wrap(err, "failed to convert signer address to base32")
				}

				btcAddressBytes := []byte{} //witness version byte
				btcAddressBytes = append(btcAddressBytes, 0)
				btcAddressBytes = append(btcAddressBytes, base32Bytes...)

				signerAddress, err := bech32.Encode("bc", btcAddressBytes)
				if err != nil {
					return sdkerrors.Wrap(err, "failed to convert and encode signer address")
				}

				standardSigIsValid, err := VerifyBIP322Signature(signerAddress, signature, message)
				if err != nil {
					return sdkerrors.Wrapf(err, "failed to verify bitcoin signature %s %s %s", message, signature, signerAddress)
				}

				humanReadableSigIsValid, err := VerifyBIP322Signature(signerAddress, signature, humanReadableStr)
				if err != nil {
					return sdkerrors.Wrapf(err, "failed to verify bitcoin signature %s %s %s", humanReadableStr, signature, signerAddress)
				}

				if !standardSigIsValid && !humanReadableSigIsValid {
					return sdkerrors.Wrap(types.ErrorInvalidSigner, "unable to verify signer signature of Bitcoin signature")
				}
			}

			return nil
		} else {

			if len(feePayerSig) != ethcrypto.SignatureLength {
				return sdkerrors.Wrap(types.ErrorInvalidSigner, "signature length doesn't match typical [R||S||V] signature 65 bytes")
			}

			byteStr := []byte(jsonStr)
			hexSig := hex.EncodeToString(feePayerSig)
			pubKeyBytes := pubKey.Address().Bytes()

			//Standard JSON signature verification
			isHashSigValid, hashSigErr := sigverify.VerifyEllipticCurveHexSignatureEx(
				ethcommon.Address(pubKeyBytes),
				// []byte(jsonHashHexStr),
				byteStr,
				"0x"+hexSig,
			)
			if isHashSigValid && hashSigErr == nil {
				return nil
			}

			isHumanReadableStrSigValid, humanReadableStrSigErr := sigverify.VerifyEllipticCurveHexSignatureEx(
				ethcommon.Address(pubKey.Address().Bytes()),
				[]byte(humanReadableStr),
				"0x"+hex.EncodeToString(feePayerSig),
			)
			if isHumanReadableStrSigValid && humanReadableStrSigErr == nil {
				return nil
			}

			//If we do not pass the above check, we will try to EIP-712 sign the message

			//TODO: Make this num nested structs to account for the nested structs in the EIP712 message?
			if len(jsonStr) > 1000 {
				return sdkerrors.Wrapf(types.ErrInvalidChainID, "could not verify standard JSON signature and tx is tooqq expensive for Ethereum EIP712 signature verification")
			}

			sigHash, _, err := apitypes.TypedDataAndHash(typedData)
			if err != nil {
				return sdkerrors.Wrapf(err, "failed to compute typed data hash")
			}

			// Remove the recovery offset if needed (ie. Metamask eip712 signature)
			if feePayerSig[ethcrypto.RecoveryIDOffset] == 27 || feePayerSig[ethcrypto.RecoveryIDOffset] == 28 {
				feePayerSig[ethcrypto.RecoveryIDOffset] -= 27
			}

			if CheckFeePayerPubKey(pubKey, sigHash, feePayerSig) != nil {
				return sdkerrors.Wrap(types.ErrorInvalidSigner, "unable to verify Ethereum signature")
			}

			standardMsgSigValid := secp256k1.VerifySignature(pubKey.Bytes(), sigHash, feePayerSig[:len(feePayerSig)-1])
			if !standardMsgSigValid {
				return sdkerrors.Wrap(types.ErrorInvalidSigner, "unable to verify Ethereum signature")
			}
		}

		return nil
	default:
		return sdkerrors.Wrapf(types.ErrTooManySignatures, "unexpected SignatureData %T", sigData)
	}
}

func VerifyBIP322Signature(signerAddress string, signature string, message string) (bool, error) {
	//These will be in the following format:
	//[0x02 or 0x03] [LENGTH_BYTE, ...(LENGTH bytes that go into witness[0])] [0x21, ...(33 byte (len = 0x21) public key that was used to sign the message)]
	//- First byte is either a 0x02 or 0x03 (not sure exactly why but apparently it is - https://github.com/ACken2/bip322-js/blob/f0f9373b3a1da19e017c518b891522eaa4bcccdd/src/helpers/Address.ts#L85)
	//- Next part is the details that go into witness[0] (length of this part is determined by the second byte) - again, not exactly sure what this is
	//- Last part is the public key that was used to sign the message (33 bytes) - There is a length byte 0x21 before this part to denote the length of the public key (0x21 = 33)
	encodedSigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	//Convert the address to a public key script (I believe this is the denotation for a pay-to-this-address script but not 100% sure)
	pkScript := convertAddressToScriptPubkey(signerAddress)

	//Decode the length of the witness[0] part
	encodedSigLenByte := encodedSigBytes[1]
	encodedSigLen := int(encodedSigLenByte)

	//Extract witness[0] and witness[1] from the encoded signature
	PUB_KEY_LEN := 33
	encodedSig := encodedSigBytes[len(encodedSigBytes)-PUB_KEY_LEN-1-encodedSigLen : len(encodedSigBytes)-PUB_KEY_LEN-1]
	pubKey := encodedSigBytes[len(encodedSigBytes)-33:]

	//Recreate the witness to pass into the VerifySignature function
	witness := wire.TxWitness{
		encodedSig,
		pubKey,
	}

	//This is the BIP322 verification function from the Unisat wallet
	verified := bip322.VerifySignature(witness, pkScript, message)
	return verified, nil
}

func convertAddressToScriptPubkey(address string) []byte {
	btcAddress, err := btcutil.DecodeAddress(address, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}

	outputScript, err := txscript.PayToAddrScript(btcAddress)
	if err != nil {
		panic(err)
	}

	return outputScript
}

func CheckFeePayerPubKey(pubKey cryptotypes.PubKey, message []byte, feePayerSig []byte) error {
	feePayerPubkey, err := secp256k1.RecoverPubkey(message, feePayerSig)
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

	return nil
}

func GetNumNestedStructs(m map[string]interface{}) int {
	numNestedStructs := 0
	for _, v := range m {
		switch v.(type) {
		case map[string]interface{}:
			numNestedStructs += GetNumNestedStructs(v.(map[string]interface{}))

		case []interface{}:
			for _, v := range v.([]interface{}) {
				switch v.(type) {
				case map[string]interface{}:
					numNestedStructs += GetNumNestedStructs(v.(map[string]interface{}))
				}
			}
		}

	}
	return numNestedStructs + 1
}
