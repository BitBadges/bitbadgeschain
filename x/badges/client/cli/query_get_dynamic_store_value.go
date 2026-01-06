package cli

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = fmt.Sprintf

func CmdGetDynamicStoreValue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dynamic-store-value [store-id] [address]",
		Short: "Query dynamic store value",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqStoreId := args[0]
			reqAddress := args[1]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetDynamicStoreValueRequest{

				StoreId: reqStoreId,
				Address: reqAddress,
			}

			res, err := queryClient.GetDynamicStoreValue(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
} 