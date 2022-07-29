package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

var _ = strconv.Itoa(0)

func CmdGetBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-balance [badge-id] [subbadge-id] [address]",
		Short: "Query getBalance",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqBadgeId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			reqSubbadgeId, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}
			reqAddress, err := cast.ToUint64E(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetBalanceRequest{

				BadgeId:    reqBadgeId,
				SubbadgeId: reqSubbadgeId,
				Address:    reqAddress,
			}

			res, err := queryClient.GetBalance(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
