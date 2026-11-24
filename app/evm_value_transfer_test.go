package app

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
)

// TestEVMValueTransferMovesBaseDenom drives a signed EVM value transfer through
// FinalizeBlock on the 18-decimal chain and verifies exact 1:1 accounting
// between x/bank (abadge) and EVM StateDB/RPC.
func TestEVMValueTransferMovesIntegerDenom(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	evmDenom := evmtypes.GetEVMCoinDenom()
	extDenom := evmtypes.GetEVMCoinExtendedDenom()
	require.Equal(t, appparams.BaseCoinUnit, evmDenom, "precondition: EVM denom is the base denom")
	require.Equal(t, evmDenom, extDenom, "precondition: base and extended denom are identical at 18 decimals")

	key, from := newFundedEVMAccount(t, app, ctx)
	to := newEVMAddress(t)
	fromAcc := sdk.AccAddress(from.Bytes())
	toAcc := sdk.AccAddress(to.Bytes())

	fromBefore := app.BankKeeper.GetBalance(ctx, fromAcc, evmDenom).Amount
	require.Empty(t, app.BankKeeper.GetAllBalances(ctx, toAcc), "precondition: recipient starts empty")

	// A deliberately awkward 18-decimal amount (not a power of 10^9).
	value := new(big.Int).SetInt64(1_234_567_891_234_567)

	nonce := nonceOf(t, app, from)
	tx, err := ethtypes.SignTx(
		ethtypes.NewTransaction(nonce, to, value, 30_000, testGasPrice, nil),
		ethtypes.NewEIP155Signer(evmChainID(t)),
		key,
	)
	require.NoError(t, err)

	res := deliverBlock(t, app, 1, from, tx)
	require.Equal(t, uint32(0), res[0].Code, "tx failed: %s", res[0].Log)
	_ = ethLogs(t, res[0]) // asserts VmError is empty

	// Committed state.
	ctx = app.NewContext(true)

	feeWei := new(big.Int).Mul(testGasPrice, big.NewInt(res[0].GasUsed))

	wantFrom := fromBefore.Sub(sdkmath.NewIntFromBigInt(value)).Sub(sdkmath.NewIntFromBigInt(feeWei))
	gotFrom := app.BankKeeper.GetBalance(ctx, fromAcc, evmDenom).Amount
	require.Equal(t, wantFrom.String(), gotFrom.String(), "sender balance must drop by value + fee")

	gotTo := app.BankKeeper.GetBalance(ctx, toAcc, evmDenom).Amount
	require.Equal(t, value.String(), gotTo.String(), "recipient balance must equal the transferred value")

	// The EVM sees exactly what was sent.
	evmRes, err := app.EVMKeeper.Balance(ctx, &evmtypes.QueryBalanceRequest{Address: to.Hex()})
	require.NoError(t, err)
	require.Equal(t, value.String(), evmRes.Balance, "eth_getBalance of recipient must equal the transferred value")

	evmFromRes, err := app.EVMKeeper.Balance(ctx, &evmtypes.QueryBalanceRequest{Address: from.Hex()})
	require.NoError(t, err)
	require.Equal(t, wantFrom.String(), evmFromRes.Balance, "eth_getBalance of sender must match bank")
}
