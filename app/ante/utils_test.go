package ante_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	types2 "github.com/cosmos/cosmos-sdk/x/bank/types"
	types3 "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/ethereum/eip712"
	"github.com/evmos/ethermint/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	bitbadgesapp "github.com/bitbadges/bitbadgeschain/app"
	ante "github.com/bitbadges/bitbadgeschain/app/ante"
	"github.com/bitbadges/bitbadgeschain/encoding"
	"github.com/evmos/ethermint/tests"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

type AnteTestSuite struct {
	suite.Suite

	ctx             sdk.Context
	app             *bitbadgesapp.App
	clientCtx       client.Context
	anteHandler     sdk.AnteHandler
	ethSigner       ethtypes.Signer
	enableFeemarket bool
	enableLondonHF  bool
	evmParamsOption func(*evmtypes.Params)
}

var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

const TestGasLimit uint64 = 100000

func (suite *AnteTestSuite) SetupTest() {
	checkTx := false
	simapp.FlagEnabledValue = true
	simapp.FlagCommitValue = true

	encodingConfig := encoding.MakeConfig(bitbadgesapp.ModuleBasics)
	// We're using TestMsg amino encoding in some tests, so register it here.
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)

	app := bitbadgesapp.Setup(
		false, nil,
	)

	suite.app = app

	suite.ctx = suite.app.BaseApp.NewContext(checkTx, tmproto.Header{Height: 2, ChainID: "ethermint_9000-1", Time: time.Now().UTC()})
	suite.ctx = suite.ctx.WithMinGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(evmtypes.DefaultEVMDenom, sdk.OneInt())))
	suite.ctx = suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1000000000000000000))

	infCtx := suite.ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	suite.app.AccountKeeper.SetParams(infCtx, authtypes.DefaultParams())

	suite.clientCtx = client.Context{}.WithTxConfig(encodingConfig.TxConfig)

	options := ante.HandlerOptions{
		AccountKeeper:   suite.app.AccountKeeper,
		BankKeeper:      suite.app.BankKeeper,
		FeegrantKeeper:  suite.app.FeeGrantKeeper,
		IBCKeeper:       suite.app.IBCKeeper,
		SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
	}

	suite.Require().NoError(options.Validate())

	suite.anteHandler = ante.NewAnteHandler(options)
	suite.ethSigner = ethtypes.LatestSignerForChainID(big.NewInt(9000))
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, &AnteTestSuite{
		enableLondonHF: true,
	})
}

func (suite *AnteTestSuite) CreateTestCosmosTxBuilder(gasPrice sdk.Int, denom string, msgs ...sdk.Msg) client.TxBuilder {
	txBuilder := suite.clientCtx.TxConfig.NewTxBuilder()

	txBuilder.SetGasLimit(TestGasLimit)
	fees := &sdk.Coins{{Denom: denom, Amount: gasPrice.MulRaw(int64(TestGasLimit))}}
	txBuilder.SetFeeAmount(*fees)
	err := txBuilder.SetMsgs(msgs...)
	suite.Require().NoError(err)
	return txBuilder
}

func (suite *AnteTestSuite) CreateTestEIP712TxBuilderMsgSend(from sdk.AccAddress, priv cryptotypes.PrivKey, chainId string, gas uint64, gasAmount sdk.Coins) client.TxBuilder {
	// Build MsgSend
	recipient := sdk.AccAddress(common.Address{}.Bytes())
	msgSend := types2.NewMsgSend(from, recipient, sdk.NewCoins(sdk.NewCoin(evmtypes.DefaultEVMDenom, sdk.NewInt(1))))
	return suite.CreateTestEIP712CosmosTxBuilder(from, priv, chainId, gas, gasAmount, msgSend)
}

func (suite *AnteTestSuite) CreateTestEIP712TxBuilderMsgDelegate(from sdk.AccAddress, priv cryptotypes.PrivKey, chainId string, gas uint64, gasAmount sdk.Coins) client.TxBuilder {
	// Build MsgSend
	valEthAddr := tests.GenerateAddress()
	valAddr := sdk.ValAddress(valEthAddr.Bytes())
	msgSend := types3.NewMsgDelegate(from, valAddr, sdk.NewCoin(evmtypes.DefaultEVMDenom, sdk.NewInt(20)))
	return suite.CreateTestEIP712CosmosTxBuilder(from, priv, chainId, gas, gasAmount, msgSend)
}

func (suite *AnteTestSuite) CreateTestEIP712CosmosTxBuilder(
	from sdk.AccAddress, priv cryptotypes.PrivKey, chainId string, gas uint64, gasAmount sdk.Coins, msg sdk.Msg,
) client.TxBuilder {
	var err error

	nonce, err := suite.app.AccountKeeper.GetSequence(suite.ctx, from)
	suite.Require().NoError(err)

	pc, err := types.ParseChainID(chainId)
	suite.Require().NoError(err)
	ethChainId := pc.Uint64()

	// GenerateTypedData TypedData
	var ethermintCodec codec.ProtoCodecMarshaler
	fee := legacytx.NewStdFee(gas, gasAmount)
	accNumber := suite.app.AccountKeeper.GetAccount(suite.ctx, from).GetAccountNumber()

	data := legacytx.StdSignBytes(chainId, accNumber, nonce, 0, fee, []sdk.Msg{msg}, "", nil)
	typedData, err := eip712.WrapTxToTypedData(ethermintCodec, ethChainId, msg, data, &eip712.FeeDelegationOptions{
		FeePayer: from,
	})
	suite.Require().NoError(err)

	sigHash, _, err := apitypes.TypedDataAndHash(typedData)
	suite.Require().NoError(err)

	// Sign typedData
	keyringSigner := tests.NewSigner(priv)
	signature, pubKey, err := keyringSigner.SignByAddress(from, sigHash)
	suite.Require().NoError(err)
	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper

	// Add ExtensionOptionsWeb3Tx extension
	var option *codectypes.Any
	option, err = codectypes.NewAnyWithValue(&types.ExtensionOptionsWeb3Tx{
		FeePayer:         from.String(),
		TypedDataChainID: ethChainId,
		FeePayerSig:      signature,
	})
	suite.Require().NoError(err)

	suite.clientCtx.TxConfig.SignModeHandler()
	txBuilder := suite.clientCtx.TxConfig.NewTxBuilder()
	builder, ok := txBuilder.(authtx.ExtensionOptionsTxBuilder)
	suite.Require().True(ok)

	builder.SetExtensionOptions(option)
	builder.SetFeeAmount(gasAmount)
	builder.SetGasLimit(gas)

	sigsV2 := signing.SignatureV2{
		PubKey: pubKey,
		Data: &signing.SingleSignatureData{
			SignMode: signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
		},
		Sequence: nonce,
	}

	err = builder.SetSignatures(sigsV2)
	suite.Require().NoError(err)

	err = builder.SetMsgs(msg)
	suite.Require().NoError(err)

	return builder
}

func NextFn(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
	return ctx, nil
}

var _ sdk.Tx = &invalidTx{}

type invalidTx struct{}

func (invalidTx) GetMsgs() []sdk.Msg   { return []sdk.Msg{nil} }
func (invalidTx) ValidateBasic() error { return nil }
