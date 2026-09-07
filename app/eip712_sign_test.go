package app

import (
	"math/big"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/evm/crypto/ethsecp256k1"
	"github.com/cosmos/evm/ethereum/eip712"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/cosmos/gogoproto/proto"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/stretchr/testify/require"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	tokenizationprecompile "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// eip712TestChainID is the Cosmos chain id the signed doc and the delivery
// context agree on. The test app's InitChain leaves it empty, so it is set on
// the header explicitly.
const eip712TestChainID = "bitbadges-eip712-test"

// TestEIP712SignedUniversalUpdateCollectionPassesAnte signs a
// MsgUniversalUpdateCollection the way an EVM wallet does — over the EIP-712
// typed-data hash rather than the raw sign bytes — and asserts the ante handler
// accepts it. The message carries an empty CollectionPermissions sub-struct so
// the typed-data tree contains an empty type, and the test pins the EIP-712
// encoding of that type to `Name()`, which is what eth-sig-util-based wallets
// hash on their side.
func TestEIP712SignedUniversalUpdateCollectionPassesAnte(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContextLegacy(false, cmtproto.Header{Height: 1, ChainID: eip712TestChainID})

	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	signer := sdk.AccAddress(priv.PubKey().Address())

	// Funding creates the account the ante handler looks up and pays the fee.
	funds := sdk.NewCoins(sdk.NewCoin(appparams.BaseCoinUnit, sdkmath.NewInt(1_000_000_000_000)))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", funds))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", signer, funds))
	acc := app.AccountKeeper.GetAccount(ctx, signer)
	require.NotNil(t, acc)

	msg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                     signer.String(),
		CollectionId:                sdkmath.NewUint(0),
		UpdateCollectionPermissions: true,
		CollectionPermissions:       &tokenizationtypes.CollectionPermissions{},
	}

	const gasLimit = 300_000
	fee := sdk.NewCoins(sdk.NewCoin(appparams.BaseCoinUnit, sdkmath.NewInt(10*gasLimit)))

	builder := app.txConfig.NewTxBuilder()
	require.NoError(t, builder.SetMsgs(msg))
	builder.SetFeeAmount(fee)
	builder.SetGasLimit(gasLimit)

	signMode := signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	sigV2 := signingtypes.SignatureV2{
		PubKey:   priv.PubKey(),
		Data:     &signingtypes.SingleSignatureData{SignMode: signMode},
		Sequence: acc.GetSequence(),
	}
	require.NoError(t, builder.SetSignatures(sigV2))

	signerData := authsigning.SignerData{
		Address:       signer.String(),
		ChainID:       ctx.ChainID(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
		PubKey:        priv.PubKey(),
	}
	signBytes, err := authsigning.GetSignBytesAdapter(ctx, app.txConfig.SignModeHandler(), signMode, signerData, builder.GetTx())
	require.NoError(t, err)

	// A wallet hashes the typed data, not the sign bytes.
	typedData, err := eip712.GetEIP712TypedDataForMsg(signBytes)
	require.NoError(t, err)

	var emptyTypes []string
	for name, fields := range typedData.Types {
		if len(fields) == 0 {
			emptyTypes = append(emptyTypes, name)
		}
	}
	require.NotEmpty(t, emptyTypes, "expected the empty CollectionPermissions to produce an empty EIP-712 type")

	encodedType := string(typedData.EncodeType(typedData.PrimaryType))
	for _, name := range emptyTypes {
		require.Contains(t, encodedType, name+"()", "empty type must encode as Name()")
	}

	sigHash, _, err := apitypes.TypedDataAndHash(typedData)
	require.NoError(t, err)
	sig, err := priv.Sign(sigHash)
	require.NoError(t, err)
	sig[ethcrypto.RecoveryIDOffset] += 27

	sigV2.Data = &signingtypes.SingleSignatureData{SignMode: signMode, Signature: sig}
	require.NoError(t, builder.SetSignatures(sigV2))

	_, err = app.AnteHandler()(ctx, builder.GetTx(), false)
	require.NoError(t, err)
}

// TestEIP712SignatureOverWrongPayloadIsRejected keeps the acceptance test
// honest: a signature over an EIP-712 hash for a different sequence must not
// pass the ante handler.
func TestEIP712SignatureOverWrongPayloadIsRejected(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContextLegacy(false, cmtproto.Header{Height: 1, ChainID: eip712TestChainID})

	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	signer := sdk.AccAddress(priv.PubKey().Address())

	funds := sdk.NewCoins(sdk.NewCoin(appparams.BaseCoinUnit, sdkmath.NewInt(1_000_000_000_000)))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", funds))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", signer, funds))
	acc := app.AccountKeeper.GetAccount(ctx, signer)
	require.NotNil(t, acc)

	msg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:               signer.String(),
		CollectionId:          sdkmath.NewUint(0),
		CollectionPermissions: &tokenizationtypes.CollectionPermissions{},
	}

	const gasLimit = 300_000
	builder := app.txConfig.NewTxBuilder()
	require.NoError(t, builder.SetMsgs(msg))
	builder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(appparams.BaseCoinUnit, sdkmath.NewInt(10*gasLimit))))
	builder.SetGasLimit(gasLimit)

	signMode := signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	sigV2 := signingtypes.SignatureV2{
		PubKey:   priv.PubKey(),
		Data:     &signingtypes.SingleSignatureData{SignMode: signMode},
		Sequence: acc.GetSequence(),
	}
	require.NoError(t, builder.SetSignatures(sigV2))

	// Sign bytes for the wrong sequence.
	signerData := authsigning.SignerData{
		Address:       signer.String(),
		ChainID:       ctx.ChainID(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence() + 1,
		PubKey:        priv.PubKey(),
	}
	signBytes, err := authsigning.GetSignBytesAdapter(ctx, app.txConfig.SignModeHandler(), signMode, signerData, builder.GetTx())
	require.NoError(t, err)
	typedData, err := eip712.GetEIP712TypedDataForMsg(signBytes)
	require.NoError(t, err)
	sigHash, _, err := apitypes.TypedDataAndHash(typedData)
	require.NoError(t, err)
	sig, err := priv.Sign(sigHash)
	require.NoError(t, err)
	sig[ethcrypto.RecoveryIDOffset] += 27

	sigV2.Data = &signingtypes.SingleSignatureData{SignMode: signMode, Signature: sig}
	require.NoError(t, builder.SetSignatures(sigV2))

	_, err = app.AnteHandler()(ctx, builder.GetTx(), false)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "signature verification failed"), err.Error())
}

// TestPrecompileEthTxDeliversThroughBaseApp is the regression guard for the
// common EVM path: a signed Ethereum transaction calling the tokenization
// precompile still executes through FinalizeBlock after the EIP-712 wiring.
func TestPrecompileEthTxDeliversThroughBaseApp(t *testing.T) {
	require.NoError(t, tokenizationprecompile.GetABILoadError())

	app := Setup(false)
	ctx := app.NewContext(false)
	key, from := newFundedEVMAccount(t, app, ctx)

	// Mainnet activates the precompile through an upgrade; the test genesis
	// does not, so a call to 0x1001 would otherwise hit an empty account.
	precompile := ethcommon.HexToAddress(tokenizationprecompile.TokenizationPrecompileAddress)
	require.NoError(t, app.EVMKeeper.EnableStaticPrecompiles(ctx, precompile))

	calldata, err := tokenizationprecompile.ABI.Pack("rangeContains", big.NewInt(1), big.NewInt(10), big.NewInt(5))
	require.NoError(t, err)

	tx, err := ethtypes.SignTx(
		ethtypes.NewTransaction(nonceOf(t, app, from), precompile, big.NewInt(0), 200_000, testGasPrice, calldata),
		ethtypes.NewEIP155Signer(evmChainID(t)),
		key,
	)
	require.NoError(t, err)

	res := deliverBlock(t, app, 1, from, tx)
	require.Equal(t, uint32(0), res[0].Code, "precompile call failed: %s", res[0].Log)

	var txMsgData sdk.TxMsgData
	require.NoError(t, proto.Unmarshal(res[0].Data, &txMsgData))
	require.Len(t, txMsgData.MsgResponses, 1)
	var response evmtypes.MsgEthereumTxResponse
	require.NoError(t, proto.Unmarshal(txMsgData.MsgResponses[0].Value, &response))
	require.Empty(t, response.VmError, "EVM reverted: %s", response.VmError)

	out, err := tokenizationprecompile.ABI.Unpack("rangeContains", response.Ret)
	require.NoError(t, err)
	require.Equal(t, []interface{}{true}, out)
}
