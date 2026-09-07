package app

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	sendmanagerprecompile "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
)

// TestSendManagerPrecompileSendPersistsAfterCommit moves bank coins through
// the sendmanager precompile inside a signed EVM transaction and checks the
// balances after the block is committed.
//
// At commit x/vm writes every dirty StateDB account's balance back to the
// bank. A plain call leaves the sender's object clean, but as soon as the same
// transaction also moves EVM value (here: msg.value on the precompile call,
// as any payable contract in the path would) the sender is dirty and its
// StateDB balance wins. A bank send performed by a precompile therefore has
// to be mirrored into the StateDB, or the commit hands the sender its coins
// back.
func TestSendManagerPrecompileSendPersistsAfterCommit(t *testing.T) {
	for _, valueUnits := range []int64{0, 1} {
		t.Run(fmt.Sprintf("value=%d", valueUnits), func(t *testing.T) {
			app := Setup(false)
			ctx := app.NewContext(false)
			require.NoError(t, app.EnableAllPrecompiles(ctx))

			denom := evmtypes.GetEVMCoinDenom()
			key, from := newFundedEVMAccount(t, app, ctx)
			to := newEVMAddress(t)
			fromAcc := sdk.AccAddress(from.Bytes())
			toAcc := sdk.AccAddress(to.Bytes())
			fromBefore := app.BankKeeper.GetBalance(ctx, fromAcc, denom).Amount

			const amount = int64(1_000)
			msgJSON := fmt.Sprintf(`{"to_address":%q,"amount":[{"denom":%q,"amount":"%d"}]}`, to.Hex(), denom, amount)
			input, err := sendmanagerprecompile.ABI.Pack(sendmanagerprecompile.SendMethod, msgJSON)
			require.NoError(t, err)

			precompileAddr := ethcommon.HexToAddress(sendmanagerprecompile.SendManagerPrecompileAddress)
			value := new(big.Int).Mul(big.NewInt(valueUnits), ubadgeToWeiFactor)
			nonce := nonceOf(t, app, from)
			tx, err := ethtypes.SignTx(
				ethtypes.NewTransaction(nonce, precompileAddr, value, maxTestGasLimit, testGasPrice, input),
				ethtypes.NewEIP155Signer(evmChainID(t)),
				key,
			)
			require.NoError(t, err)

			res := deliverBlock(t, app, 1, from, tx)
			require.Equal(t, uint32(0), res[0].Code, "tx failed: %s", res[0].Log)
			_ = ethLogs(t, res[0])

			ctx = app.NewContext(true)
			feeUnits := new(big.Int).Div(new(big.Int).Mul(testGasPrice, big.NewInt(res[0].GasUsed)), ubadgeToWeiFactor)

			wantFrom := fromBefore.Sub(sdkmath.NewInt(amount + valueUnits)).Sub(sdkmath.NewIntFromBigInt(feeUnits))
			require.Equal(t, wantFrom.String(), app.BankKeeper.GetBalance(ctx, fromAcc, denom).Amount.String(),
				"sender must be poorer by exactly the sent amount, the value and the fee after commit")
			require.Equal(t, sdkmath.NewInt(amount).String(), app.BankKeeper.GetBalance(ctx, toAcc, denom).Amount.String(),
				"recipient must hold exactly the sent amount after commit")

			// The EVM's view agrees with x/bank on both sides.
			wantToWei := new(big.Int).Mul(big.NewInt(amount), ubadgeToWeiFactor)
			evmRes, err := app.EVMKeeper.Balance(ctx, &evmtypes.QueryBalanceRequest{Address: to.Hex()})
			require.NoError(t, err)
			require.Equal(t, wantToWei.String(), evmRes.Balance)

			wantFromWei := new(big.Int).Mul(wantFrom.BigInt(), ubadgeToWeiFactor)
			evmRes, err = app.EVMKeeper.Balance(ctx, &evmtypes.QueryBalanceRequest{Address: from.Hex()})
			require.NoError(t, err)
			require.Equal(t, wantFromWei.String(), evmRes.Balance)
		})
	}
}
