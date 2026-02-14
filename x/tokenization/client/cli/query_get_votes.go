package cli

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetVotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [collection-id] [approval-level] [approver-address] [approval-id] [proposal-id]",
		Short: "Query votes",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqCollectionId := args[0]
			reqApprovalLevel := args[1]
			reqApproverAddress := args[2]
			reqApprovalId := args[3]
			reqProposalId := args[4]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetVotesRequest{
				CollectionId:    reqCollectionId,
				ApprovalLevel:   reqApprovalLevel,
				ApproverAddress: reqApproverAddress,
				ApprovalId:      reqApprovalId,
				ProposalId:      reqProposalId,
			}

			res, err := queryClient.GetVotes(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
