package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetCollectionIdForProtocol() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-collection-id-for-protocol [name] [address]",
		Short: "Query getCollectionIdForProtocol",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetCollectionIdForProtocolRequest{
								Name: args[0],
								Address: args[1],	
			}

            

			res, err := queryClient.GetCollectionIdForProtocol(cmd.Context(), params)
			if err != nil {

					return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

    return cmd
}
