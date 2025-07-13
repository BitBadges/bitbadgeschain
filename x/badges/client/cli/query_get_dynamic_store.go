package cli

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = fmt.Sprintf

func CmdGetDynamicStore() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-dynamic-store [store-id]",
		Short: "Query getDynamicStore",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqStoreId := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetDynamicStoreRequest{

				StoreId: reqStoreId,
			}

			res, err := queryClient.GetDynamicStore(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
