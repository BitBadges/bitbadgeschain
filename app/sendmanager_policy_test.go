package app

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	appparams "github.com/bitbadges/bitbadgeschain/app/params"
	sendkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
	sendtypes "github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
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
			// Internal module operations retain their own authorization policy.
			require.NoError(t, app.SendmanagerKeeper.SendCoinsWithAliasRouting(ctx, sender, receiver, funds))
		})
	}
}
