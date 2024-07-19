package cli

import (
	"strconv"

	"bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

// message QueryGetChallengeTrackerRequest {
//   string collectionId = 1 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
//   string approvalLevel = 2; //"collection" or "incoming" or "outgoing"
//   string approverAddress = 3; //if approvalLevel is "collection", leave blank
//   string challengeTrackerId = 4;
//   string leafIndex = 5 [(gogoproto.customtype) = "Uint", (gogoproto.nullable) = false];
// }

func CmdGetChallengeTracker() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-num-used-for-merkle-challenge [collectionId] [approvalLevel] [approverAddress] [approvalId] [challengeTrackerId] [leafIndex]",
		Short: "Query getChallengeTracker",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetChallengeTrackerRequest{
				CollectionId:       types.NewUintFromString(args[0]),
				ApprovalLevel:      args[1],
				ApproverAddress:    args[2],
				ApprovalId:         args[3],
				ChallengeTrackerId: args[4],
				LeafIndex:          types.NewUintFromString(args[5]),
			}

			res, err := queryClient.GetChallengeTracker(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
