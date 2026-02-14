package cli

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdIsAddressReservedProtocol() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-address-reserved-protocol [address]",
		Short: "Query if an address is a reserved protocol address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAddress, err := cast.ToStringE(args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryIsAddressReservedProtocolRequest{
				Address: reqAddress,
			}

			res, err := queryClient.IsAddressReservedProtocol(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
