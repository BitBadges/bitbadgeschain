package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"bitbadgeschain/x/anchor/types"

	sdkmath "cosmossdk.io/math"
)

func CmdQueryData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getAnchorData [location-id]",
		Short: "get anchor data at a location",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			locationId := sdkmath.NewUintFromString(args[0])

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.GetValueAtLocation(cmd.Context(), &types.QueryGetValueAtLocationRequest{
				LocationId: locationId,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
