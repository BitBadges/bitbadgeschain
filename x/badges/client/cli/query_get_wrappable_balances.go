package cli

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdGetWrappableBalances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-wrappable-balances [denom] [address]",
		Short: "Query GetWrappableBalances",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqDenom := args[0]
			if err != nil {
				return err
			}

			reqAddress, err := cast.ToStringE(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetWrappableBalancesRequest{
				Denom:   reqDenom,
				Address: reqAddress,
			}

			res, err := queryClient.GetWrappableBalances(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
