package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetApprovalTrackers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-approvals-trackers [collectionId] [approvalLevel] [approverAddress] [approvalId] [amountTrackerId] [trackerType] [approvedAddress]",
		Short: "Query getApprovalTrackers",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetApprovalTrackerRequest{
				CollectionId:    args[0],
				ApprovalLevel:   args[1],
				ApproverAddress: args[2],
				ApprovalId:      args[3],
				AmountTrackerId: args[4],
				TrackerType:     args[5],
				ApprovedAddress: args[6],
			}

			res, err := queryClient.GetApprovalTracker(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
