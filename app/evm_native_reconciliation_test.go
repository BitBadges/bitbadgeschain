package app

import (
	"bytes"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/x/vm/statedb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	"github.com/bitbadges/bitbadgeschain/pkg/evmcompat"
)

func TestNativeAtomicEventsFollowNestedCommit(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	sender := sdk.AccAddress(newEVMAddress(t).Bytes())
	receiver := sdk.AccAddress(newEVMAddress(t).Bytes())
	funds := sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 100))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", funds))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", sender, funds))
	db := statedb.New(ctx, app.EVMKeeper, statedb.TxConfig{})
	native, err := db.GetCacheContext()
	require.NoError(t, err)
	require.NoError(t, db.FlushToCacheCtx())
	handler := evmcompat.NewBalanceHandlerFactory(app.BankKeeper).NewBalanceHandler()
	handler.BeforeBalanceChange(native)
	outer := evmcompat.NewAtomicContext(native)
	inner := evmcompat.NewAtomicContext(outer.Ctx())
	require.NoError(t, app.BankKeeper.SendCoins(inner.Ctx(), sender, receiver, sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 10))))
	inner.Commit()
	require.NoError(t, db.FlushToCacheCtx())
	outer.Rollback()
	outer.Rollback()
	require.Empty(t, native.EventManager().Events(), "discarded nested writes must discard their events")
	accepted := evmcompat.NewAtomicContext(native)
	require.NoError(t, app.BankKeeper.SendCoins(accepted.Ctx(), sender, receiver, sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 3))))
	accepted.Commit()
	accepted.Commit()
	accepted.Rollback()
	require.NoError(t, handler.AfterBalanceChange(native, db))
	require.NoError(t, db.Commit())
	require.Equal(t, "97", app.BankKeeper.GetBalance(ctx, sender, appparams.BaseCoinUnit).Amount.String())
	require.Equal(t, "3", app.BankKeeper.GetBalance(ctx, receiver, appparams.BaseCoinUnit).Amount.String())
}

func TestNativeBalanceEventsPreserveFullAddress(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	sender := sdk.AccAddress(newEVMAddress(t).Bytes())
	receiver := sdk.AccAddress(append(bytes.Repeat([]byte{7}, 12), sender.Bytes()...))
	funds := sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 100))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", funds))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", sender, funds))
	supply := app.BankKeeper.GetSupply(ctx, appparams.BaseCoinUnit)
	db := statedb.New(ctx, app.EVMKeeper, statedb.TxConfig{})
	initial := db.GetBalance(common.BytesToAddress(sender)).Clone()
	native, err := db.GetCacheContext()
	require.NoError(t, err)
	require.NoError(t, db.FlushToCacheCtx())
	handler := evmcompat.NewBalanceHandlerFactory(app.BankKeeper).NewBalanceHandler()
	handler.BeforeBalanceChange(native)
	require.NoError(t, app.BankKeeper.SendCoins(native, sender, receiver, sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 10))))
	require.NoError(t, handler.AfterBalanceChange(native, db))
	require.NotEqual(t, initial.String(), db.GetBalance(common.BytesToAddress(sender)).String())
	require.NoError(t, db.Commit())
	require.Equal(t, "90", app.BankKeeper.GetBalance(ctx, sender, appparams.BaseCoinUnit).Amount.String())
	require.Equal(t, "10", app.BankKeeper.GetBalance(ctx, receiver, appparams.BaseCoinUnit).Amount.String())
	require.Equal(t, supply, app.BankKeeper.GetSupply(ctx, appparams.BaseCoinUnit))
	require.Equal(t, sdkmath.NewInt(100), app.BankKeeper.GetBalance(ctx, sender, appparams.BaseCoinUnit).Amount.Add(app.BankKeeper.GetBalance(ctx, receiver, appparams.BaseCoinUnit).Amount))

	db = statedb.New(ctx, app.EVMKeeper, statedb.TxConfig{})
	_ = db.GetBalance(common.BytesToAddress(sender))
	native, err = db.GetCacheContext()
	require.NoError(t, err)
	require.NoError(t, db.FlushToCacheCtx())
	handler = evmcompat.NewBalanceHandlerFactory(app.BankKeeper).NewBalanceHandler()
	handler.BeforeBalanceChange(native)
	require.NoError(t, app.BankKeeper.SendCoins(native, receiver, sender, sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 10))))
	require.NoError(t, handler.AfterBalanceChange(native, db))
	require.NoError(t, db.Commit())
	require.Equal(t, "100", app.BankKeeper.GetBalance(ctx, sender, appparams.BaseCoinUnit).Amount.String())
	require.True(t, app.BankKeeper.GetBalance(ctx, receiver, appparams.BaseCoinUnit).Amount.IsZero())
	require.Equal(t, supply, app.BankKeeper.GetSupply(ctx, appparams.BaseCoinUnit))
}
