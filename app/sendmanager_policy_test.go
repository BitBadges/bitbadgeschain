package app

import (
	"testing"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	sendkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	sendtypes "github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	transfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"
)

func TestSendmanagerPublicBankPolicy(t *testing.T) {
	for _, blocked := range []bool{false, true} {
		name := "disabled denom"
		if blocked {
			name = "restricted recipient"
		}
		t.Run(name, func(t *testing.T) {
			app := Setup(false)
			ctx := app.NewContext(false)
			sender := sdk.AccAddress(newEVMAddress(t).Bytes())
			receiver := sdk.AccAddress(newEVMAddress(t).Bytes())
			funds := sdk.NewCoins(sdk.NewInt64Coin(appparams.BaseCoinUnit, 10))
			require.NoError(t, app.BankKeeper.MintCoins(ctx, "mint", funds))
			require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", sender, funds))
			if blocked {
				receiver = authtypes.NewModuleAddress("mint")
				require.True(t, app.BankKeeper.BlockedAddr(receiver))
			} else {
				app.BankKeeper.SetSendEnabled(ctx, appparams.BaseCoinUnit, false)
			}
			_, err := sendkeeper.NewMsgServerImpl(app.SendmanagerKeeper).SendWithAliasRouting(ctx, &sendtypes.MsgSendWithAliasRouting{
				FromAddress: sender.String(), ToAddress: receiver.String(), Amount: funds,
			})
			require.Error(t, err)
			require.Equal(t, funds, app.BankKeeper.GetAllBalances(ctx, sender))
			require.True(t, app.BankKeeper.GetBalance(ctx, receiver, appparams.BaseCoinUnit).IsZero())
			require.Error(t, app.SendmanagerKeeper.SendCoinsWithAliasRouting(ctx, sender, receiver, funds))
			require.Error(t, app.SendmanagerKeeper.SendCoinWithAliasRouting(ctx, sender, receiver, &funds[0]))
			require.Equal(t, funds, app.BankKeeper.GetAllBalances(ctx, sender))
		})
	}
}

func TestIBCReceiveRejectsBlockedModuleAccounts(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)
	for _, name := range []string{"gamm", "tokenization", "evm", "precisebank"} {
		receiver := authtypes.NewModuleAddress(name)
		err := app.TransferKeeper.OnRecvPacket(ctx, transfertypes.InternalTransferRepresentation{
			Token:  transfertypes.Token{Denom: transfertypes.Denom{Base: "uatom"}, Amount: "1"},
			Sender: "remote-sender", Receiver: receiver.String(),
		}, "transfer", "channel-0", "transfer", "channel-1")
		require.ErrorContains(t, err, "not allowed to receive funds", name)
		require.Empty(t, app.BankKeeper.GetAllBalances(ctx, receiver))
	}
}
