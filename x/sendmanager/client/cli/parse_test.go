package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func coins(t *testing.T, s string) sdk.Coins {
	t.Helper()
	c, err := sdk.ParseCoinsNormalized(s)
	require.NoError(t, err)
	return c
}

func TestParseSend_positionalOnly(t *testing.T) {
	t.Parallel()

	got, err := ParseSend(
		[]string{"alice", "bb1to", "1000ubadge", "7uqa"},
		SendFlagValues{},
	)
	require.NoError(t, err)
	require.Equal(t, "alice", got.FromAddress)
	require.Equal(t, "bb1to", got.ToAddress)
	require.True(t, got.Amount.Equal(coins(t, "1000ubadge,7uqa")))

	got, err = ParseSend(
		[]string{"alice", "bb1to", "1000ubadge,7uqa"},
		SendFlagValues{},
	)
	require.NoError(t, err)
	require.True(t, got.Amount.Equal(coins(t, "7uqa,1000ubadge")))
}

func TestParseSend_flagsOnly(t *testing.T) {
	t.Parallel()

	got, err := ParseSend(nil, SendFlagValues{
		From:      "alice",
		FromSet:   true,
		ToAddress: "bb1to",
		ToSet:     true,
		Amount:    []string{"1000ubadge", "7uqa"},
		AmountSet: true,
	})
	require.NoError(t, err)
	require.Equal(t, "alice", got.FromAddress)
	require.Equal(t, "bb1to", got.ToAddress)
	require.True(t, got.Amount.Equal(coins(t, "1000ubadge,7uqa")))
}

func TestParseSend_mixedFields(t *testing.T) {
	t.Parallel()

	got, err := ParseSend(
		[]string{"alice"},
		SendFlagValues{
			ToAddress: "bb1to",
			ToSet:     true,
			Amount:    []string{"5badgeslp:1:utoken"},
			AmountSet: true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "alice", got.FromAddress)
	require.Equal(t, "bb1to", got.ToAddress)
	require.True(t, got.Amount.Equal(coins(t, "5badgeslp:1:utoken")))
}

func TestParseSend_matchingBothForms(t *testing.T) {
	t.Parallel()

	got, err := ParseSend(
		[]string{"alice", "bb1to", "1000ubadge", "7uqa"},
		SendFlagValues{
			From:      "alice",
			FromSet:   true,
			ToAddress: "bb1to",
			ToSet:     true,
			Amount:    []string{"1000ubadge,7uqa"},
			AmountSet: true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "alice", got.FromAddress)
	require.Equal(t, "bb1to", got.ToAddress)
	require.True(t, got.Amount.Equal(coins(t, "1000ubadge,7uqa")))
}

func TestParseSend_conflicts(t *testing.T) {
	t.Parallel()

	_, err := ParseSend(
		[]string{"alice", "bb1to", "1000ubadge"},
		SendFlagValues{From: "bob", FromSet: true, ToAddress: "bb1to", ToSet: true, Amount: []string{"1000ubadge"}, AmountSet: true},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "conflicting from address")
	require.Contains(t, err.Error(), "alice")
	require.Contains(t, err.Error(), "bob")

	_, err = ParseSend(
		[]string{"alice", "bb1to", "1000ubadge"},
		SendFlagValues{ToAddress: "bb1other", ToSet: true},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "conflicting to address")
	require.Contains(t, err.Error(), "--to-address")

	_, err = ParseSend(
		[]string{"alice", "bb1to", "1000ubadge"},
		SendFlagValues{Amount: []string{"7uqa"}, AmountSet: true},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "conflicting amount")
}

func TestParseSend_missing(t *testing.T) {
	t.Parallel()

	_, err := ParseSend(nil, SendFlagValues{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing from address")

	_, err = ParseSend([]string{"alice"}, SendFlagValues{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing to address")

	_, err = ParseSend([]string{"alice", "bb1to"}, SendFlagValues{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing amount")

	_, err = ParseSend(nil, SendFlagValues{
		ToAddress: "bb1to", ToSet: true, Amount: []string{"1ubadge"}, AmountSet: true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing from address")
	require.Contains(t, err.Error(), "--from")

	_, err = ParseSend(nil, SendFlagValues{
		From: "alice", FromSet: true, Amount: []string{"1ubadge"}, AmountSet: true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing to address")
	require.Contains(t, err.Error(), "--to-address")

	_, err = ParseSend(nil, SendFlagValues{
		From: "alice", FromSet: true, ToAddress: "bb1to", ToSet: true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing amount")

	_, err = ParseSend(nil, SendFlagValues{
		From: "alice", FromSet: true, ToAddress: "bb1to", ToSet: true, AmountSet: true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing amount")
}

func TestParseSend_invalidAmount(t *testing.T) {
	t.Parallel()

	_, err := ParseSend([]string{"alice", "bb1to", "not-a-coin"}, SendFlagValues{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid amount")
}

func TestParseBalance_positionalOnly(t *testing.T) {
	t.Parallel()

	got, err := ParseBalance([]string{"bb1acct", "uqa"}, BalanceFlagValues{})
	require.NoError(t, err)
	require.Equal(t, BalanceInput{Address: "bb1acct", Denom: "uqa"}, got)
}

func TestParseBalance_flagsOnly(t *testing.T) {
	t.Parallel()

	got, err := ParseBalance(nil, BalanceFlagValues{
		Address:  "bb1acct",
		AddrSet:  true,
		Denom:    "badgeslp:1:utoken",
		DenomSet: true,
	})
	require.NoError(t, err)
	require.Equal(t, BalanceInput{Address: "bb1acct", Denom: "badgeslp:1:utoken"}, got)
}

func TestParseBalance_mixedAndMatching(t *testing.T) {
	t.Parallel()

	got, err := ParseBalance([]string{"bb1acct"}, BalanceFlagValues{Denom: "uqa", DenomSet: true})
	require.NoError(t, err)
	require.Equal(t, BalanceInput{Address: "bb1acct", Denom: "uqa"}, got)

	got, err = ParseBalance([]string{"bb1acct", "uqa"}, BalanceFlagValues{
		Address: "bb1acct", AddrSet: true, Denom: "uqa", DenomSet: true,
	})
	require.NoError(t, err)
	require.Equal(t, BalanceInput{Address: "bb1acct", Denom: "uqa"}, got)
}

func TestParseBalance_conflicts(t *testing.T) {
	t.Parallel()

	_, err := ParseBalance([]string{"bb1acct", "uqa"}, BalanceFlagValues{Address: "bb1other", AddrSet: true})
	require.Error(t, err)
	require.Contains(t, err.Error(), "conflicting address")
	require.Contains(t, err.Error(), "--address")

	_, err = ParseBalance([]string{"bb1acct", "uqa"}, BalanceFlagValues{Denom: "ubadge", DenomSet: true})
	require.Error(t, err)
	require.Contains(t, err.Error(), "conflicting denom")
	require.Contains(t, err.Error(), "--denom")
}

func TestParseBalance_missingAndExtra(t *testing.T) {
	t.Parallel()

	_, err := ParseBalance(nil, BalanceFlagValues{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing address")

	_, err = ParseBalance([]string{"bb1acct"}, BalanceFlagValues{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing denom")

	_, err = ParseBalance(nil, BalanceFlagValues{Address: "bb1acct", AddrSet: true})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing denom")
	require.Contains(t, err.Error(), "--denom")

	_, err = ParseBalance([]string{"bb1acct", "uqa", "extra"}, BalanceFlagValues{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "accepts at most 2 positional arg(s)")
}

func TestCommandHelpAndFlags(t *testing.T) {
	t.Parallel()

	txCmd := CmdSendWithAliasRouting()
	require.NotNil(t, txCmd.Flags().Lookup(flagToAddress))
	require.NotNil(t, txCmd.Flags().Lookup(flagAmount))
	require.Contains(t, txCmd.Use, "[from_key_or_address]")
	require.Contains(t, txCmd.Use, "[to_address]")
	require.Contains(t, txCmd.Use, "[amount]")
	require.True(t, strings.Contains(txCmd.Long, "--to-address") && strings.Contains(txCmd.Long, "--amount") && strings.Contains(txCmd.Long, "--from"))
	require.Contains(t, txCmd.Long, "Positional")
	require.Contains(t, txCmd.Long, "must match")

	qCmd := CmdBalance()
	require.NotNil(t, qCmd.Flags().Lookup(flagAddress))
	require.NotNil(t, qCmd.Flags().Lookup(flagDenom))
	require.Contains(t, qCmd.Use, "[address]")
	require.Contains(t, qCmd.Use, "[denom]")
	require.Contains(t, qCmd.Long, "--address")
	require.Contains(t, qCmd.Long, "--denom")
	require.Contains(t, qCmd.Long, "Positional")
	require.Contains(t, qCmd.Long, "must match")
}

func TestCobraParse_sendFlagsAndPositionals(t *testing.T) {
	t.Parallel()

	cmd := CmdSendWithAliasRouting()
	require.NoError(t, cmd.Flags().Parse([]string{
		"--from", "alice",
		"--to-address", "bb1to",
		"--amount", "1000ubadge",
		"--amount", "7uqa",
	}))
	from, err := cmd.Flags().GetString("from")
	require.NoError(t, err)
	to, err := cmd.Flags().GetString(flagToAddress)
	require.NoError(t, err)
	amt, err := cmd.Flags().GetStringArray(flagAmount)
	require.NoError(t, err)
	got, err := ParseSend(cmd.Flags().Args(), SendFlagValues{
		From: from, FromSet: cmd.Flags().Changed("from"),
		ToAddress: to, ToSet: cmd.Flags().Changed(flagToAddress),
		Amount: amt, AmountSet: cmd.Flags().Changed(flagAmount),
	})
	require.NoError(t, err)
	require.Equal(t, "alice", got.FromAddress)
	require.Equal(t, "bb1to", got.ToAddress)
	require.True(t, got.Amount.Equal(coins(t, "1000ubadge,7uqa")))

	cmd = CmdSendWithAliasRouting()
	require.NoError(t, cmd.Flags().Parse([]string{"alice", "bb1to", "1000ubadge", "7uqa"}))
	got, err = ParseSend(cmd.Flags().Args(), SendFlagValues{
		FromSet:   cmd.Flags().Changed("from"),
		ToSet:     cmd.Flags().Changed(flagToAddress),
		AmountSet: cmd.Flags().Changed(flagAmount),
	})
	require.NoError(t, err)
	require.Equal(t, "alice", got.FromAddress)
	require.Equal(t, "bb1to", got.ToAddress)
	require.True(t, got.Amount.Equal(coins(t, "1000ubadge,7uqa")))
}

func TestCobraParse_balanceFlagsAndPositionals(t *testing.T) {
	t.Parallel()

	cmd := CmdBalance()
	require.NoError(t, cmd.Flags().Parse([]string{"--address", "bb1acct", "--denom", "uqa"}))
	addr, err := cmd.Flags().GetString(flagAddress)
	require.NoError(t, err)
	denom, err := cmd.Flags().GetString(flagDenom)
	require.NoError(t, err)
	got, err := ParseBalance(cmd.Flags().Args(), BalanceFlagValues{
		Address: addr, AddrSet: cmd.Flags().Changed(flagAddress),
		Denom: denom, DenomSet: cmd.Flags().Changed(flagDenom),
	})
	require.NoError(t, err)
	require.Equal(t, BalanceInput{Address: "bb1acct", Denom: "uqa"}, got)

	cmd = CmdBalance()
	require.NoError(t, cmd.Flags().Parse([]string{"bb1acct", "uqa"}))
	got, err = ParseBalance(cmd.Flags().Args(), BalanceFlagValues{
		AddrSet:  cmd.Flags().Changed(flagAddress),
		DenomSet: cmd.Flags().Changed(flagDenom),
	})
	require.NoError(t, err)
	require.Equal(t, BalanceInput{Address: "bb1acct", Denom: "uqa"}, got)
}
