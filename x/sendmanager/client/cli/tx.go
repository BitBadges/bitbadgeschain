package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

// GetTxCmd returns the transaction commands for the sendmanager module.
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdSendWithAliasRouting())
	return cmd
}

// CmdSendWithAliasRouting returns the send-with-alias-routing tx command.
// It accepts the documented positional args and the v36 flags (--from,
// --to-address, repeated --amount).
func CmdSendWithAliasRouting() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-with-alias-routing [from_key_or_address] [to_address] [amount]",
		Short: "Send coins with alias denom routing",
		Long: fmt.Sprintf(`Send coins from one address to another. Supports both standard coins and alias denoms (e.g., badgeslp:). This mirrors cosmos bank MsgSend but routes through sendmanager.

Two equivalent forms are accepted:

  Positional (as shown in Usage):
    %s tx sendmanager send-with-alias-routing <from_key_or_address> <to_address> <amount> [...]

  Flags (backward compatible):
    %s tx sendmanager send-with-alias-routing --from <from_key_or_address> --to-address <to_address> --amount <coin> [--amount <coin> ...]

<amount> may be a comma-separated coins string (1000ubadge,7uqa) or space-separated coins (1000ubadge 7uqa). Repeat --amount once per coin. --from is also the signing key.

Fields may be mixed (some positional, some flags). If the same field is given both as a positional argument and as a flag, the values must match (amounts are compared after normalization); otherwise the command errors.`, "bitbadgeschaind", "bitbadgeschaind"),
		Example: fmt.Sprintf(`%s tx sendmanager send-with-alias-routing alice bb1... 1000ubadge 5badgeslp:1:utoken
%s tx sendmanager send-with-alias-routing --from alice --to-address bb1... --amount 1000ubadge --amount 5badgeslp:1:utoken`, "bitbadgeschaind", "bitbadgeschaind"),
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			toAddress, err := cmd.Flags().GetString(flagToAddress)
			if err != nil {
				return err
			}
			amountVals, err := cmd.Flags().GetStringArray(flagAmount)
			if err != nil {
				return err
			}
			from, err := cmd.Flags().GetString(flags.FlagFrom)
			if err != nil {
				return err
			}

			input, err := ParseSend(args, SendFlagValues{
				From:      from,
				FromSet:   cmd.Flags().Changed(flags.FlagFrom),
				ToAddress: toAddress,
				ToSet:     cmd.Flags().Changed(flagToAddress),
				Amount:    amountVals,
				AmountSet: cmd.Flags().Changed(flagAmount),
			})
			if err != nil {
				return err
			}

			if !cmd.Flags().Changed(flags.FlagFrom) {
				if err := cmd.Flags().Set(flags.FlagFrom, input.FromAddress); err != nil {
					return err
				}
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			toAddr, err := resolveAddressOrKeyName(clientCtx, input.ToAddress)
			if err != nil {
				return err
			}

			msg := &types.MsgSendWithAliasRouting{
				FromAddress: clientCtx.GetFromAddress().String(),
				ToAddress:   toAddr,
				Amount:      input.Amount,
			}

			// If --from was already set and positional from was also given,
			// ParseSend required the raw strings to match, so the signer is
			// the positional from_key_or_address as well.
			if msg.FromAddress == "" {
				msg.FromAddress = input.FromAddress
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagToAddress, "", "Recipient address or key name")
	cmd.Flags().StringArray(flagAmount, nil, "Amount to send (repeat the flag for multiple coins, e.g. --amount 1000ubadge --amount 7uqa)")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// resolveAddressOrKeyName accepts a bech32 address or a local key name, matching
// v36 autocli's "--to-address account address or key name" flag.
func resolveAddressOrKeyName(clientCtx client.Context, nameOrAddr string) (string, error) {
	if _, err := sdk.AccAddressFromBech32(nameOrAddr); err == nil {
		return nameOrAddr, nil
	}
	if clientCtx.Keyring == nil {
		return nameOrAddr, nil
	}
	rec, err := clientCtx.Keyring.Key(nameOrAddr)
	if err != nil {
		return nameOrAddr, nil
	}
	addr, err := rec.GetAddress()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}
