package cli

import (
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

// message QueryGetETHSignatureTrackerRequest {
//   string collectionId = 1;
//   string approvalLevel = 2; //"collection" or "incoming" or "outgoing"
//   string approverAddress = 3; //if approvalLevel is "collection", leave blank
//   string approvalId = 4;
//   string challengeTrackerId = 5;
//   string signature = 6;
// }

func CmdGetETHSignatureTracker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "num-used-for-eth-signature-challenge [collectionId] [approvalLevel] [approverAddress] [approvalId] [challengeTrackerId] [signature]",
		Short: "Query ETH signature tracker",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetETHSignatureTrackerRequest{
				CollectionId:       args[0],
				ApprovalLevel:      args[1],
				ApproverAddress:    args[2],
				ApprovalId:         args[3],
				ChallengeTrackerId: args[4],
				Signature:          args[5],
			}

			res, err := queryClient.GetETHSignatureTracker(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
