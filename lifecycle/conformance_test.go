package lifecycle

import (
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	"github.com/bitbadges/bitbadgeschain/app"
	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	signing "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFullTransactionConformance(t *testing.T) {
	runMu.Lock()
	defer runMu.Unlock()
	application := app.Setup(false)
	defer application.Close()
	ctx := application.NewContext(false)
	key := secp256k1.GenPrivKeyFromSecret([]byte("lifecycle-conformance-public-fixture"))
	addr := sdk.AccAddress(key.PubKey().Address())
	funds := sdk.NewCoins(sdk.NewCoin("ubadge", sdkmath.NewInt(1000000000)))
	require.NoError(t, application.BankKeeper.MintCoins(ctx, "mint", funds))
	require.NoError(t, application.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", addr, funds))
	config := application.GetTxConfig()
	fee := sdk.NewCoins(sdk.NewCoin("ubadge", sdkmath.NewInt(3000000)))
	build := func(sequence uint64, wrongSignature bool, msgs ...sdk.Msg) sdk.Tx {
		if len(msgs) == 0 {
			msgs = []sdk.Msg{&banktypes.MsgSend{FromAddress: addr.String(), ToAddress: bob, Amount: sdk.NewCoins(sdk.NewCoin("ubadge", sdkmath.NewInt(4)))}}
		}
		builder := config.NewTxBuilder()
		require.NoError(t, builder.SetMsgs(msgs...))
		builder.SetGasLimit(300000)
		builder.SetFeeAmount(fee)
		mode := signing.SignMode_SIGN_MODE_DIRECT
		require.NoError(t, builder.SetSignatures(signing.SignatureV2{PubKey: key.PubKey(), Data: &signing.SingleSignatureData{SignMode: mode}, Sequence: sequence}))
		data := authsigning.SignerData{Address: addr.String(), ChainID: ctx.ChainID(), AccountNumber: 0, Sequence: sequence, PubKey: key.PubKey()}
		sig, err := tx.SignWithPrivKey(ctx, mode, data, builder, key, config, sequence)
		require.NoError(t, err)
		require.NoError(t, builder.SetSignatures(sig))
		if wrongSignature {
			builder.SetMemo("tampered")
		}
		return builder.GetTx()
	}
	signed := build(0, false)
	gas, result, err := application.SimDeliver(config.TxEncoder(), signed)
	require.NoError(t, err)
	require.Positive(t, gas.GasUsed)
	require.NotEmpty(t, result.Events)
	require.Equal(t, "4", application.BankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(bob), "ubadge").Amount.String())
	require.Equal(t, "996999996", application.BankKeeper.GetBalance(ctx, addr, "ubadge").Amount.String())
	require.Equal(t, uint64(1), application.AccountKeeper.GetAccount(ctx, addr).GetSequence())
	_, _, err = application.SimDeliver(config.TxEncoder(), signed)
	require.Error(t, err, "replayed sequence rejected")
	_, _, err = application.SimDeliver(config.TxEncoder(), build(1, true))
	require.ErrorContains(t, err, "signature")
	require.Equal(t, "4", application.BankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(bob), "ubadge").Amount.String())
	m := artifact(t, "subscription-fixed-price")
	var value map[string]any
	require.NoError(t, json.Unmarshal(m.Value, &value))
	value["creator"] = addr.String()
	value["manager"] = addr.String()
	m.Value = raw(value)
	msg, err := decodeMessage(application, m)
	require.NoError(t, err)
	cached, _ := ctx.CacheContext()
	moduleResult, err := application.MsgServiceRouter().Handler(msg)(cached, msg)
	require.NoError(t, err)
	require.NotEmpty(t, moduleResult.Events)
	expected, found := application.TokenizationKeeper.GetCollectionFromStore(cached, sdkmath.NewUint(1))
	require.True(t, found)
	_, result, err = application.SimDeliver(config.TxEncoder(), build(1, false, msg))
	require.NoError(t, err)
	require.NotEmpty(t, result.Events)
	actual, found := application.TokenizationKeeper.GetCollectionFromStore(ctx, sdkmath.NewUint(1))
	require.True(t, found)
	require.Equal(t, expected, actual)
	require.Equal(t, uint64(2), application.AccountKeeper.GetAccount(ctx, addr).GetSequence())
	before := application.BankKeeper.GetBalance(ctx, addr, "ubadge").Amount
	missing := &tokenization.MsgUpdateCollection{Creator: addr.String(), CollectionId: sdkmath.NewUint(999), UpdateManager: true, Manager: bob}
	_, _, err = application.SimDeliver(config.TxEncoder(), build(2, false, missing))
	require.Error(t, err)
	require.Equal(t, before.Sub(sdkmath.NewInt(3000000)).String(), application.BankKeeper.GetBalance(ctx, addr, "ubadge").Amount.String())
	require.Equal(t, uint64(3), application.AccountKeeper.GetAccount(ctx, addr).GetSequence())
	fee = sdk.NewCoins(sdk.NewCoin("ubadge", sdkmath.NewInt(2000000000)))
	_, _, err = application.SimDeliver(config.TxEncoder(), build(3, false))
	require.ErrorContains(t, err, "insufficient funds")
	require.Equal(t, uint64(3), application.AccountKeeper.GetAccount(ctx, addr).GetSequence())

}
