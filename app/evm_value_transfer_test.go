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

// TestEVMValueTransferMovesIntegerDenom drives a signed EVM value transfer
// through FinalizeBlock and checks the money moved in x/bank's integer denom.
//
// The chain is 9-decimal: x/vm's StateDB works in the 18-decimal extended
// denom and, at commit, writes each dirty account's balance back through the
// bank keeper it was given. That write has to land in the integer denom (plus
// the precisebank fractional store) — an x/bank entry keyed by the extended
// denom is a balance nothing reads and nothing backs.
func TestEVMValueTransferMovesIntegerDenom(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	evmDenom := evmtypes.GetEVMCoinDenom()
	extDenom := evmtypes.GetEVMCoinExtendedDenom()
	require.Equal(t, appparams.BaseCoinUnit, evmDenom, "precondition: EVM denom is the integer denom")
	require.NotEqual(t, evmDenom, extDenom, "precondition: extended denom is distinct, so the split is observable")

	key, from := newFundedEVMAccount(t, app, ctx)
	to := newEVMAddress(t)
	fromAcc := sdk.AccAddress(from.Bytes())
	toAcc := sdk.AccAddress(to.Bytes())

	fromBefore := app.BankKeeper.GetBalance(ctx, fromAcc, evmDenom).Amount
	extSupplyBefore := app.BankKeeper.GetSupply(ctx, extDenom).Amount
	require.Empty(t, app.BankKeeper.GetAllBalances(ctx, toAcc), "precondition: recipient starts empty")
	require.True(t, app.BankKeeper.GetAllBalances(ctx, fromAcc).AmountOf(extDenom).IsZero(),
		"precondition: sender holds no extended-denom coin")

	// A deliberately awkward integer amount, so a pass cannot come from a
	// power-of-ten coincidence.
	const valueUnits = int64(1_234_567)
	value := new(big.Int).Mul(big.NewInt(valueUnits), ubadgeToWeiFactor)

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

	// testGasPrice is an exact multiple of the 9→18 decimal factor, so the fee
	// is a whole number of integer-denom units.
	feeWei := new(big.Int).Mul(testGasPrice, big.NewInt(res[0].GasUsed))
	feeUnits := new(big.Int).Div(feeWei, ubadgeToWeiFactor)
	require.Zero(t, new(big.Int).Mod(feeWei, ubadgeToWeiFactor).Sign(), "fee must be whole integer-denom units")

	wantFrom := fromBefore.Sub(sdkmath.NewInt(valueUnits)).Sub(sdkmath.NewIntFromBigInt(feeUnits))
	gotFrom := app.BankKeeper.GetBalance(ctx, fromAcc, evmDenom).Amount
	require.Equal(t, wantFrom.String(), gotFrom.String(), "sender integer-denom balance must drop by value + fee")

	gotTo := app.BankKeeper.GetBalance(ctx, toAcc, evmDenom).Amount
	require.Equal(t, sdkmath.NewInt(valueUnits).String(), gotTo.String(), "recipient integer-denom balance must equal the value")

	// The EVM sees exactly what was sent.
	evmRes, err := app.EVMKeeper.Balance(ctx, &evmtypes.QueryBalanceRequest{Address: to.Hex()})
	require.NoError(t, err)
	require.Equal(t, value.String(), evmRes.Balance, "eth_getBalance of recipient must equal the transferred value")

	// No account carries an x/bank entry in the extended denom, and its
	// supply is untouched.
	for _, acc := range []sdk.AccAddress{fromAcc, toAcc} {
		all := app.BankKeeper.GetAllBalances(ctx, acc)
		require.True(t, all.AmountOf(extDenom).IsZero(), "%s must hold no %s in x/bank, got %s", acc, extDenom, all)
	}
	require.Equal(t, extSupplyBefore.String(), app.BankKeeper.GetSupply(ctx, extDenom).Amount.String(),
		"extended-denom supply must be unchanged")
}

// TestEVMValueTransferFractionalUnitsUsePreciseBank sends a value that is not
// a whole number of integer-denom units and checks the remainder lands in
// precisebank's fractional store on both sides, with the integer denom
// carrying the rest.
func TestEVMValueTransferFractionalUnitsUsePreciseBank(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	evmDenom := evmtypes.GetEVMCoinDenom()
	extDenom := evmtypes.GetEVMCoinExtendedDenom()
	require.NotEqual(t, evmDenom, extDenom)

	key, from := newFundedEVMAccount(t, app, ctx)
	to := newEVMAddress(t)
	fromAcc := sdk.AccAddress(from.Bytes())
	toAcc := sdk.AccAddress(to.Bytes())

	fromBefore := app.BankKeeper.GetBalance(ctx, fromAcc, evmDenom).Amount
	require.True(t, app.PreciseBankKeeper.GetFractionalBalance(ctx, fromAcc).IsZero(), "precondition: sender has no fractional balance")

	// 7 whole units plus 0.25 of a unit.
	const wholeUnits = int64(7)
	fraction := new(big.Int).Div(ubadgeToWeiFactor, big.NewInt(4))
	value := new(big.Int).Mul(big.NewInt(wholeUnits), ubadgeToWeiFactor)
	value.Add(value, fraction)

	nonce := nonceOf(t, app, from)
	tx, err := ethtypes.SignTx(
		ethtypes.NewTransaction(nonce, to, value, 30_000, testGasPrice, nil),
		ethtypes.NewEIP155Signer(evmChainID(t)),
		key,
	)
	require.NoError(t, err)

	res := deliverBlock(t, app, 1, from, tx)
	_ = ethLogs(t, res[0])
	ctx = app.NewContext(true)

	feeUnits := new(big.Int).Div(new(big.Int).Mul(testGasPrice, big.NewInt(res[0].GasUsed)), ubadgeToWeiFactor)

	// Recipient: 7 whole units in x/bank, a quarter unit in the fractional store.
	require.Equal(t, sdkmath.NewInt(wholeUnits).String(),
		app.BankKeeper.GetBalance(ctx, toAcc, evmDenom).Amount.String())
	require.Equal(t, fraction.String(),
		app.PreciseBankKeeper.GetFractionalBalance(ctx, toAcc).String())

	// Sender: the quarter unit is borrowed from a whole one, so the integer
	// balance drops by 8 (+ fee) and three quarters of a unit is left over.
	wantFromInt := fromBefore.Sub(sdkmath.NewInt(wholeUnits + 1)).Sub(sdkmath.NewIntFromBigInt(feeUnits))
	require.Equal(t, wantFromInt.String(),
		app.BankKeeper.GetBalance(ctx, fromAcc, evmDenom).Amount.String())
	require.Equal(t, new(big.Int).Sub(ubadgeToWeiFactor, fraction).String(),
		app.PreciseBankKeeper.GetFractionalBalance(ctx, fromAcc).String())

	// The EVM's view of both sides is exact.
	evmRes, err := app.EVMKeeper.Balance(ctx, &evmtypes.QueryBalanceRequest{Address: to.Hex()})
	require.NoError(t, err)
	require.Equal(t, value.String(), evmRes.Balance)

	wantFromWei := new(big.Int).Mul(fromBefore.BigInt(), ubadgeToWeiFactor)
	wantFromWei.Sub(wantFromWei, value)
	wantFromWei.Sub(wantFromWei, new(big.Int).Mul(testGasPrice, big.NewInt(res[0].GasUsed)))
	evmRes, err = app.EVMKeeper.Balance(ctx, &evmtypes.QueryBalanceRequest{Address: from.Hex()})
	require.NoError(t, err)
	require.Equal(t, wantFromWei.String(), evmRes.Balance)

	for _, acc := range []sdk.AccAddress{fromAcc, toAcc} {
		require.True(t, app.BankKeeper.GetAllBalances(ctx, acc).AmountOf(extDenom).IsZero(),
			"%s must hold no %s in x/bank", acc, extDenom)
	}
}
