package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetBalanceForToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance-for-token [collection-id] [address] [token-id] [time]",
		Short: "Query the balance amount for a specific token ID at a specific time",
		Long: `Query the balance amount for a specific token ID at a specific time.
The time parameter is optional and defaults to the current block time.
Time should be specified in milliseconds since epoch.` + "\n" + QueryHelpLinks("balance-for-token"),
		Args: cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqCollectionId := args[0]
			reqAddress := args[1]
			reqTokenId := args[2]

			reqTime := ""
			if len(args) > 3 {
				reqTime = args[3]
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetBalanceForTokenRequest{
				CollectionId: reqCollectionId,
				Address:      reqAddress,
				TokenId:      reqTokenId,
				Time:         reqTime,
			}

			res, err := queryClient.GetBalanceForToken(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
