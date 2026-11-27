package cli

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	flagToAddress = "to-address"
	flagAmount    = "amount"
	flagAddress   = "address"
	flagDenom     = "denom"
)

// SendFlagValues is the flag-form input for send-with-alias-routing.
// From is --from (the v36 source of from_address; there is no --from-address flag).
type SendFlagValues struct {
	From      string
	FromSet   bool
	ToAddress string
	ToSet     bool
	Amount    []string
	AmountSet bool
}

// SendInput is the resolved send-with-alias-routing message fields.
type SendInput struct {
	FromAddress string
	ToAddress   string
	Amount      sdk.Coins
}

// BalanceFlagValues is the flag-form input for query sendmanager balance.
type BalanceFlagValues struct {
	Address  string
	AddrSet  bool
	Denom    string
	DenomSet bool
}

// BalanceInput is the resolved balance query fields.
type BalanceInput struct {
	Address string
	Denom   string
}

// ParseSend combines positional args [from_key_or_address] [to_address] [amount...]
// with the v36 flags --from, --to-address, and repeated --amount.
//
// Each field may come from positionals, flags, or both. If both are given,
// amounts must be equal after normalization and string fields must be
// identical; otherwise a conflict error is returned. A field given in
// neither form is a missing-field error.
func ParseSend(args []string, flags SendFlagValues) (SendInput, error) {
	var fromPos, toPos string
	var amountPos []string
	if len(args) >= 1 {
		fromPos = args[0]
	}
	if len(args) >= 2 {
		toPos = args[1]
	}
	if len(args) >= 3 {
		amountPos = args[2:]
	}

	from, err := pickString("from address", "from", fromPos, flags.From, flags.FromSet)
	if err != nil {
		return SendInput{}, err
	}
	to, err := pickString("to address", flagToAddress, toPos, flags.ToAddress, flags.ToSet)
	if err != nil {
		return SendInput{}, err
	}

	amount, err := pickAmount(amountPos, flags.Amount, flags.AmountSet)
	if err != nil {
		return SendInput{}, err
	}

	return SendInput{FromAddress: from, ToAddress: to, Amount: amount}, nil
}

// ParseBalance combines positional args [address] [denom] with the v36
// flags --address and --denom. Same per-field precedence as ParseSend:
// either form is enough; both must match; neither is an error.
func ParseBalance(args []string, flags BalanceFlagValues) (BalanceInput, error) {
	if len(args) > 2 {
		return BalanceInput{}, fmt.Errorf("accepts at most 2 positional arg(s), received %d", len(args))
	}

	var addrPos, denomPos string
	if len(args) >= 1 {
		addrPos = args[0]
	}
	if len(args) >= 2 {
		denomPos = args[1]
	}

	addr, err := pickString("address", flagAddress, addrPos, flags.Address, flags.AddrSet)
	if err != nil {
		return BalanceInput{}, err
	}
	denom, err := pickString("denom", flagDenom, denomPos, flags.Denom, flags.DenomSet)
	if err != nil {
		return BalanceInput{}, err
	}

	return BalanceInput{Address: addr, Denom: denom}, nil
}

func pickString(field, flagName, positional, flag string, flagSet bool) (string, error) {
	hasPos := positional != ""
	if hasPos && flagSet {
		if positional == flag {
			return positional, nil
		}
		return "", fmt.Errorf("conflicting %s: positional %q vs --%s %q", field, positional, flagName, flag)
	}
	if hasPos {
		return positional, nil
	}
	if flagSet && flag != "" {
		return flag, nil
	}
	return "", fmt.Errorf("missing %s: provide it as a positional argument or --%s", field, flagName)
}

func pickAmount(positional, flagVals []string, flagSet bool) (sdk.Coins, error) {
	hasPos := len(positional) > 0
	if hasPos && flagSet {
		posCoins, err := parseAmountValues(positional)
		if err != nil {
			return nil, err
		}
		flagCoins, err := parseAmountValues(flagVals)
		if err != nil {
			return nil, err
		}
		if !posCoins.Equal(flagCoins) {
			return nil, fmt.Errorf("conflicting amount: positional %q vs --%s %q", posCoins, flagAmount, flagCoins)
		}
		return posCoins, nil
	}

	var vals []string
	switch {
	case hasPos:
		vals = positional
	case flagSet:
		vals = flagVals
	default:
		return nil, fmt.Errorf("missing amount: provide it as a positional argument or --%s", flagAmount)
	}

	coins, err := parseAmountValues(vals)
	if err != nil {
		return nil, err
	}
	return coins, nil
}

// parseAmountValues accepts one or more coin strings. Each value may be a
// single coin (1000ubadge) or a comma-separated coins string (1000ubadge,7uqa),
// matching bank send and the repeated --amount flag form.
func parseAmountValues(vals []string) (sdk.Coins, error) {
	if len(vals) == 0 {
		return nil, fmt.Errorf("missing amount: provide it as a positional argument or --%s", flagAmount)
	}

	var coins sdk.Coins
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if v == "" {
			return nil, fmt.Errorf("invalid amount %q", v)
		}
		parsed, err := sdk.ParseCoinsNormalized(v)
		if err != nil {
			return nil, fmt.Errorf("invalid amount %q: %w", v, err)
		}
		if parsed.IsZero() {
			return nil, fmt.Errorf("invalid amount %q: coins cannot be empty", v)
		}
		coins = coins.Add(parsed...)
	}
	if coins.IsZero() || !coins.IsValid() {
		return nil, fmt.Errorf("invalid amount %q", strings.Join(vals, " "))
	}
	return coins, nil
}
