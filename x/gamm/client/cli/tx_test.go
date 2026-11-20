package cli_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/app/params"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/third_party/osmoutils/osmocli"
	"github.com/bitbadges/bitbadgeschain/x/gamm/client/cli"
	"github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

func init() {
	params.SetAddressPrefixes()
}

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

func testFromAddr() string {
	return sdk.AccAddress([]byte("gamm-cli-test-addr01")).String()
}

func testAffiliateAddr() string {
	return sdk.AccAddress([]byte("gamm-cli-affil-addr01")).String()
}

func TestParseAffiliatesJSON(t *testing.T) {
	addr := testAffiliateAddr()

	tests := map[string]struct {
		input    string
		expected []poolmanagertypes.Affiliate
		wantErr  bool
	}{
		"empty string": {
			input:    "",
			expected: []poolmanagertypes.Affiliate{},
		},
		"empty array": {
			input:    "[]",
			expected: []poolmanagertypes.Affiliate{},
		},
		"null": {
			input:    "null",
			expected: []poolmanagertypes.Affiliate{},
		},
		"one affiliate": {
			input: fmt.Sprintf(`[{"address":"%s","basis_points_fee":"100"}]`, addr),
			expected: []poolmanagertypes.Affiliate{
				{Address: addr, BasisPointsFee: "100"},
			},
		},
		"invalid json": {
			input:   "{not-json",
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := cli.ParseAffiliatesJSON(tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, got)
		})
	}
}

func TestSwapExactAmountInMsg(t *testing.T) {
	from := testFromAddr()
	affiliate := testAffiliateAddr()
	desc, _ := cli.NewSwapExactAmountInCmd()

	tokenIn := sdk.NewCoin("ubadge", osmomath.NewInt(1000000000))
	routes := []poolmanagertypes.SwapAmountInRoute{{
		PoolId:        1,
		TokenOutDenom: "uusdc",
	}}

	osmocli.RunTxTestCases(t, desc, map[string]osmocli.TxCliTestCase[*types.MsgSwapExactAmountIn]{
		"without affiliates flag": {
			Cmd: fmt.Sprintf("1000000000ubadge 1 --swap-route-pool-ids=1 --swap-route-denoms=uusdc --from %s", from),
			ExpectedMsg: &types.MsgSwapExactAmountIn{
				Sender:            from,
				Routes:            routes,
				TokenIn:           tokenIn,
				TokenOutMinAmount: osmomath.NewInt(1),
				Affiliates:        []poolmanagertypes.Affiliate{},
			},
		},
		"empty affiliates json": {
			Cmd: fmt.Sprintf("1000000000ubadge 1 --swap-route-pool-ids=1 --swap-route-denoms=uusdc --affiliates=[] --from %s", from),
			ExpectedMsg: &types.MsgSwapExactAmountIn{
				Sender:            from,
				Routes:            routes,
				TokenIn:           tokenIn,
				TokenOutMinAmount: osmomath.NewInt(1),
				Affiliates:        []poolmanagertypes.Affiliate{},
			},
		},
		"with affiliates": {
			Cmd: fmt.Sprintf(`1000000000ubadge 1 --swap-route-pool-ids=1 --swap-route-denoms=uusdc --affiliates=[{"address":"%s","basis_points_fee":"100"}] --from %s`, affiliate, from),
			ExpectedMsg: &types.MsgSwapExactAmountIn{
				Sender:            from,
				Routes:            routes,
				TokenIn:           tokenIn,
				TokenOutMinAmount: osmomath.NewInt(1),
				Affiliates: []poolmanagertypes.Affiliate{
					{Address: affiliate, BasisPointsFee: "100"},
				},
			},
		},
		"invalid affiliates json": {
			Cmd:         fmt.Sprintf("1000000000ubadge 1 --swap-route-pool-ids=1 --swap-route-denoms=uusdc --affiliates=not-json --from %s", from),
			ExpectedErr: true,
		},
	})
}

func TestSwapExactAmountInWithIBCTransferMsg(t *testing.T) {
	from := testFromAddr()
	affiliate := testAffiliateAddr()
	desc, _ := cli.NewSwapExactAmountInWithIBCTransferCmd()

	ibcJSON := `{"source_channel":"channel-0","receiver":"cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a","memo":"","timeout_timestamp":1}`
	expectedIBC := types.IBCTransferInfo{
		SourceChannel:    "channel-0",
		Receiver:         "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a",
		Memo:             "",
		TimeoutTimestamp: 1,
	}

	osmocli.RunTxTestCases(t, desc, map[string]osmocli.TxCliTestCase[*types.MsgSwapExactAmountInWithIBCTransfer]{
		"without affiliates": {
			Cmd: fmt.Sprintf("1000000000ubadge 1 --swap-route-pool-ids=1 --swap-route-denoms=uusdc --ibc-transfer-info=%s --from %s", ibcJSON, from),
			ExpectedMsg: &types.MsgSwapExactAmountInWithIBCTransfer{
				Sender:            from,
				Routes:            []poolmanagertypes.SwapAmountInRoute{{PoolId: 1, TokenOutDenom: "uusdc"}},
				TokenIn:           sdk.NewCoin("ubadge", osmomath.NewInt(1000000000)),
				TokenOutMinAmount: osmomath.NewInt(1),
				IbcTransferInfo:   expectedIBC,
				Affiliates:        []poolmanagertypes.Affiliate{},
			},
		},
		"with affiliates": {
			Cmd: fmt.Sprintf(`1000000000ubadge 1 --swap-route-pool-ids=1 --swap-route-denoms=uusdc --ibc-transfer-info=%s --affiliates=[{"address":"%s","basis_points_fee":"50"}] --from %s`, ibcJSON, affiliate, from),
			ExpectedMsg: &types.MsgSwapExactAmountInWithIBCTransfer{
				Sender:            from,
				Routes:            []poolmanagertypes.SwapAmountInRoute{{PoolId: 1, TokenOutDenom: "uusdc"}},
				TokenIn:           sdk.NewCoin("ubadge", osmomath.NewInt(1000000000)),
				TokenOutMinAmount: osmomath.NewInt(1),
				IbcTransferInfo:   expectedIBC,
				Affiliates: []poolmanagertypes.Affiliate{
					{Address: affiliate, BasisPointsFee: "50"},
				},
			},
		},
	})
}

func TestTxCmdUseStrings(t *testing.T) {
	t.Run("swap-exact-amount-in", func(t *testing.T) {
		desc, _ := cli.NewSwapExactAmountInCmd()
		cmd := osmocli.BuildTxCli[*types.MsgSwapExactAmountIn](desc)
		require.Equal(t, "swap-exact-amount-in [token-in] [token-out-min-amount]", cmd.Use)
		require.Equal(t, 2, desc.NumArgs)
	})
	t.Run("swap-exact-amount-in-with-ibc-transfer", func(t *testing.T) {
		desc, _ := cli.NewSwapExactAmountInWithIBCTransferCmd()
		cmd := osmocli.BuildTxCli[*types.MsgSwapExactAmountInWithIBCTransfer](desc)
		require.Equal(t, "swap-exact-amount-in-with-ibc-transfer [token-in] [token-out-min-amount]", cmd.Use)
		require.Equal(t, 2, desc.NumArgs)
	})
	t.Run("join-pool", func(t *testing.T) {
		desc, _ := cli.NewJoinPoolCmd()
		cmd := osmocli.BuildTxCli[*types.MsgJoinPool](desc)
		require.Equal(t, "join-pool", cmd.Use)
		require.Equal(t, 0, desc.NumArgs)
	})
	t.Run("exit-pool", func(t *testing.T) {
		desc, _ := cli.NewExitPoolCmd()
		cmd := osmocli.BuildTxCli[*types.MsgExitPool](desc)
		require.Equal(t, "exit-pool", cmd.Use)
		require.Equal(t, 0, desc.NumArgs)
		require.NotContains(t, cmd.Flags().Lookup("share-amount-in").Usage, "TODO")
		require.NotContains(t, cmd.Flags().Lookup("min-amounts-out").Usage, "TODO")
	})
	t.Run("swap-exact-amount-out", func(t *testing.T) {
		desc, _ := cli.NewSwapExactAmountOutCmd()
		cmd := osmocli.BuildTxCli[*types.MsgSwapExactAmountOut](desc)
		require.Equal(t, "swap-exact-amount-out [token-out] [token-in-max-amount]", cmd.Use)
		require.Equal(t, 2, desc.NumArgs)
	})
}
