package cli

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [collection-id] [approval-level] [approver-address] [approval-id] [proposal-id] [voter-address]",
		Short: "Query vote",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqCollectionId := args[0]
			reqApprovalLevel := args[1]
			reqApproverAddress := args[2]
			reqApprovalId := args[3]
			reqProposalId := args[4]
			reqVoterAddress := args[5]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetVoteRequest{
				CollectionId:    reqCollectionId,
				ApprovalLevel:   reqApprovalLevel,
				ApproverAddress: reqApproverAddress,
				ApprovalId:      reqApprovalId,
				ProposalId:      reqProposalId,
				VoterAddress:    reqVoterAddress,
			}

			res, err := queryClient.GetVote(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
