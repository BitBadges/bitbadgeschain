package app

import (
	"crypto/ecdsa"
	"math/big"
	"strconv"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/gogoproto/proto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
)

// Block-level EVM test harness.
//
// Every other EVM test in this repo executes through evmKeeper.EthereumTx(ctx, msg),
// which calls the keeper directly and bypasses baseapp. That matters more than it
// sounds: evmtypes.PatchTxResponses — the thing that numbers EVM logs
// block-globally — is invoked by the tx runner during block execution, so it never
// runs in that path. A test written at the keeper level passes identically whether
// or not a runner is installed, which is exactly why the missing runner went
// unnoticed until it was nearly shipped.
//
// These helpers drive real, signed Ethereum transactions through app.FinalizeBlock
// so block-scoped behaviour is actually observable.

// logEmitterInitCode deploys a contract that emits two LOG0 records per call.
//
// Hand-assembled rather than compiled so the test carries no Solidity toolchain
// dependency:
//
//	deploy   600b 600c 6000 39 600b 6000 f3
//	         PUSH1 11, PUSH1 12, PUSH1 0, CODECOPY, PUSH1 11, PUSH1 0, RETURN
//	         -> copies the 11 runtime bytes at offset 12 and returns them
//
//	runtime  6000 6000 a0 6000 6000 a0 00
//	         (PUSH1 0, PUSH1 0, LOG0) x2, STOP
//	         -> two zero-length, zero-topic logs
//
// Two logs per transaction rather than one is deliberate: it distinguishes
// "indices increment within a transaction" from "indices continue across
// transactions", and only the latter is what the runner provides.
const logEmitterInitCode = "600b600c600039600b6000f3" + "60006000a060006000a000"

// logsPerCall must match the runtime above.
const logsPerCall = 2

// testGasPrice is priced comfortably above the fee market's base fee.
//
// Reading the base fee from a pre-block context does not work: the feemarket
// BeginBlocker recalculates it during FinalizeBlock, so a price sampled
// beforehand is stale and the tx is rejected with "max fee per gas less than
// block base fee". These tests are about log indexing, not fee dynamics, so
// they overpay rather than track it.
//
// The EVM prices gas in the 18-decimal extended denom, so this is ~2 units of
// the display denom per gas. Accounts are funded well above the resulting cost.
var testGasPrice = new(big.Int).SetUint64(2_000_000_000_000_000_000)

// maxTestGasLimit is the largest gas limit any transaction in this harness
// uses; txFundingHeadroom covers several such transactions per account.
const (
	maxTestGasLimit   = 500_000
	txFundingHeadroom = 100
)

func evmChainID(t *testing.T) *big.Int {
	t.Helper()
	id, err := strconv.ParseUint(appparams.GetEVMChainID(), 10, 64)
	require.NoError(t, err)
	return new(big.Int).SetUint64(id)
}

// newFundedEVMAccount creates an EVM key and gives it enough native balance to
// pay for gas.
func newFundedEVMAccount(t *testing.T, app *App, ctx sdk.Context) (*ecdsa.PrivateKey, ethcommon.Address) {
	t.Helper()

	key, err := ethcrypto.GenerateKey()
	require.NoError(t, err)
	addr := ethcrypto.PubkeyToAddress(key.PublicKey)

	// Funded from the actual worst-case cost rather than a magic constant: the
	// EVM prices gas in the 18-decimal extended denom, so the right amount
	// depends on the chain's decimals and a fixed number silently stops being
	// enough if those change.
	maxCost := new(big.Int).Mul(testGasPrice, big.NewInt(int64(maxTestGasLimit)))
	maxCost.Mul(maxCost, big.NewInt(txFundingHeadroom))
	funds := sdk.NewCoins(sdk.NewCoin(appparams.BaseCoinUnit, sdkmath.NewIntFromBigInt(maxCost)))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", funds))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", sdk.AccAddress(addr.Bytes()), funds))

	return key, addr
}

// encodeEthTx wraps a signed Ethereum transaction as a Cosmos tx.
//
// MsgEthereumTx.BuildTx is what attaches the ExtensionOptionsEthereumTx and the
// fee that the EVM ante handler expects; building the tx by hand without it is
// rejected before execution.
func encodeEthTx(t *testing.T, app *App, ethTx *ethtypes.Transaction, from ethcommon.Address) []byte {
	t.Helper()

	msg := &evmtypes.MsgEthereumTx{}
	msg.FromEthereumTx(ethTx)

	// FromEthereumTx only wraps the raw tx; From has to be set explicitly or
	// the ante handler rejects with "sender address is missing".
	msg.From = from.Bytes()

	signedTx, err := msg.BuildTx(app.txConfig.NewTxBuilder(), appparams.BaseCoinUnit)
	require.NoError(t, err)

	bz, err := app.txConfig.TxEncoder()(signedTx)
	require.NoError(t, err)
	return bz
}

// deliverBlock runs the given Ethereum transactions as a single block and
// returns their results in order.
func deliverBlock(t *testing.T, app *App, height int64, from ethcommon.Address, ethTxs ...*ethtypes.Transaction) []*abci.ExecTxResult {
	t.Helper()

	txBytes := make([][]byte, 0, len(ethTxs))
	for _, ethTx := range ethTxs {
		txBytes = append(txBytes, encodeEthTx(t, app, ethTx, from))
	}

	res, err := app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: height,
		Txs:    txBytes,
	})
	require.NoError(t, err)
	require.Len(t, res.TxResults, len(ethTxs))

	_, err = app.Commit()
	require.NoError(t, err)

	return res.TxResults
}

// ethLogs pulls the EVM logs out of a delivered transaction result. This reads
// exactly the field PatchTxResponses rewrites, so it observes the runner's
// effect rather than the StateDB's per-transaction numbering.
func ethLogs(t *testing.T, res *abci.ExecTxResult) []*evmtypes.Log {
	t.Helper()
	require.Equal(t, uint32(0), res.Code, "tx failed: %s", res.Log)

	var txMsgData sdk.TxMsgData
	require.NoError(t, proto.Unmarshal(res.Data, &txMsgData))

	var out []*evmtypes.Log
	for _, rsp := range txMsgData.MsgResponses {
		var response evmtypes.MsgEthereumTxResponse
		if rsp.TypeUrl != "/"+proto.MessageName(&response) {
			continue
		}
		require.NoError(t, proto.Unmarshal(rsp.Value, &response))
		require.Empty(t, response.VmError, "EVM reverted: %s", response.VmError)
		out = append(out, response.Logs...)
	}
	return out
}

// deployLogEmitter deploys the log-emitting contract in its own block and
// returns its address.
func deployLogEmitter(t *testing.T, app *App, key *ecdsa.PrivateKey, from ethcommon.Address, height int64) ethcommon.Address {
	t.Helper()

	chainID := evmChainID(t)
	nonce := nonceOf(t, app, from)

	deployTx, err := ethtypes.SignTx(
		ethtypes.NewContractCreation(nonce, big.NewInt(0), 500_000, testGasPrice, ethcommon.FromHex(logEmitterInitCode)),
		ethtypes.NewEIP155Signer(chainID),
		key,
	)
	require.NoError(t, err)

	res := deliverBlock(t, app, height, from, deployTx)
	require.Equal(t, uint32(0), res[0].Code, "deploy failed: %s", res[0].Log)

	return ethcrypto.CreateAddress(from, nonce)
}

func nonceOf(t *testing.T, app *App, addr ethcommon.Address) uint64 {
	t.Helper()
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: app.LastBlockHeight()})
	acc := app.AccountKeeper.GetAccount(ctx, sdk.AccAddress(addr.Bytes()))
	if acc == nil {
		return 0
	}
	return acc.GetSequence()
}

// callLogEmitter builds (but does not deliver) a call to the log-emitting
// contract at the given nonce.
func callLogEmitter(t *testing.T, key *ecdsa.PrivateKey, contract ethcommon.Address, nonce uint64) *ethtypes.Transaction {
	t.Helper()

	tx, err := ethtypes.SignTx(
		ethtypes.NewTransaction(nonce, contract, big.NewInt(0), 200_000, testGasPrice, nil),
		ethtypes.NewEIP155Signer(evmChainID(t)),
		key,
	)
	require.NoError(t, err)
	return tx
}
