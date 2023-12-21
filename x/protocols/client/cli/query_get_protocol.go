package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetProtocol() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-protocol [name]",
		Short: "Query getProtocol",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetProtocolRequest{
								Name: args[0],
            }

            

			res, err := queryClient.GetProtocol(cmd.Context(), params)
            if err != nil {
                return err
            }

            return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

    return cmd
}
