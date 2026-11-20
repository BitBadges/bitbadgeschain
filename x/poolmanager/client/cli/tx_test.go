package cli_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/third_party/osmoutils/osmocli"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/client/cli"
	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

func TestParseCoinsNoSort(t *testing.T) {
	const (
		a = "aaa"
		b = "bbb"
		c = "ccc"
		d = "ddd"
	)

	var (
		ten = osmomath.NewInt(10)

		coinA = sdk.NewCoin(a, ten)
		coinB = sdk.NewCoin(b, ten)
		coinC = sdk.NewCoin(c, ten)
		coinD = sdk.NewCoin(d, ten)
	)

	tests := map[string]struct {
		coinsStr      string
		expectedCoins sdk.Coins
	}{
		"ascending": {
			coinsStr: "10aaa,10bbb,10ccc,10ddd",
			expectedCoins: sdk.Coins{
				coinA,
				coinB,
				coinC,
				coinD,
			},
		},
		"descending": {
			coinsStr: "10ddd,10ccc,10bbb,10aaa",
			expectedCoins: sdk.Coins{
				coinD,
				coinC,
				coinB,
				coinA,
			},
		},
		"mixed with different values.": {
			coinsStr: "100ddd,20bbb,300aaa,40ccc",
			expectedCoins: sdk.Coins{
				sdk.NewCoin(d, osmomath.NewInt(100)),
				sdk.NewCoin(b, osmomath.NewInt(20)),
				sdk.NewCoin(a, osmomath.NewInt(300)),
				sdk.NewCoin(c, osmomath.NewInt(40)),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			coins, err := cli.ParseCoinsNoSort(tc.coinsStr)

			require.NoError(t, err)
			require.Equal(t, tc.expectedCoins, coins)
		})
	}
}

func TestTxCmdUseStrings(t *testing.T) {
	t.Run("swap-exact-amount-in", func(t *testing.T) {
		desc, _ := cli.NewSwapExactAmountInCmd()
		cmd := osmocli.BuildTxCli[*types.MsgSwapExactAmountIn](desc)
		require.Equal(t, "swap-exact-amount-in [token-in] [token-out-min-amount]", cmd.Use)
		require.Equal(t, 2, desc.NumArgs)
	})
	t.Run("swap-exact-amount-out", func(t *testing.T) {
		desc, _ := cli.NewSwapExactAmountOutCmd()
		cmd := osmocli.BuildTxCli[*types.MsgSwapExactAmountOut](desc)
		require.Equal(t, "swap-exact-amount-out [token-out] [token-in-max-amount]", cmd.Use)
		require.Equal(t, 2, desc.NumArgs)
	})
}
