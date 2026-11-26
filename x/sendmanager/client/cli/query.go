package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

// GetQueryCmd returns the query commands for the sendmanager module.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdBalance())
	return cmd
}

// CmdBalance returns the sendmanager balance query command.
// It accepts the documented positional args and the v36 flags (--address, --denom).
func CmdBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance [address] [denom]",
		Short: "Query balance of a specific denom for an address with alias routing",
		Long: fmt.Sprintf(`Query the balance of a specific denomination for an address. Supports both standard coins and alias denoms (e.g., badgeslp:).

Two equivalent forms are accepted:

  Positional (as shown in Usage):
    %s query sendmanager balance <address> <denom>

  Flags (backward compatible):
    %s query sendmanager balance --address <address> --denom <denom>

Fields may be mixed (some positional, some flags). If the same field is given both as a positional argument and as a flag, the values must match; otherwise the command errors.`, "bitbadgeschaind", "bitbadgeschaind"),
		Example: fmt.Sprintf(`%s query sendmanager balance bb1... badgeslp:1:utoken
%s query sendmanager balance --address bb1... --denom badgeslp:1:utoken`, "bitbadgeschaind", "bitbadgeschaind"),
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			address, err := cmd.Flags().GetString(flagAddress)
			if err != nil {
				return err
			}
			denom, err := cmd.Flags().GetString(flagDenom)
			if err != nil {
				return err
			}

			input, err := ParseBalance(args, BalanceFlagValues{
				Address:  address,
				AddrSet:  cmd.Flags().Changed(flagAddress),
				Denom:    denom,
				DenomSet: cmd.Flags().Changed(flagDenom),
			})
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Balance(cmd.Context(), &types.QueryBalanceRequest{
				Address: input.Address,
				Denom:   input.Denom,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().String(flagAddress, "", "Address to query")
	cmd.Flags().String(flagDenom, "", "Denomination to query (standard coin or alias denom)")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
