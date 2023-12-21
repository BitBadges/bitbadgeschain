package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetApprovalsTrackers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-approvals-trackers [collectionId] [approvalLevel] [approverAddress] [amountTrackerId] [trackerType] [approvedAddress]",
		Short: "Query getApprovalsTrackers",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetApprovalsTrackerRequest{
				CollectionId:    types.NewUintFromString(args[0]),
				ApprovalLevel:   args[1],
				ApproverAddress: args[2],
				AmountTrackerId: args[3],
				TrackerType:     args[4],
				ApprovedAddress: args[5],
			}

			res, err := queryClient.GetApprovalsTracker(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
