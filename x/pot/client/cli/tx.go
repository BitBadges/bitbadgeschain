package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/pot/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/gogoproto/proto"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "pot transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(CmdUpdateParams())
	return cmd
}

func CmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [collection-id] [token-id] [min-balance] [mode]",
		Short: "Update x/pot params (must be sent from governance authority)",
		Long:  "Update x/pot params. Mode: staked_multiplier, equal, or credential_weighted. Set collection-id to 0 to disable.",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid collection-id: %w", err)
			}
			tokenId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid token-id: %w", err)
			}
			minBalance, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid min-balance: %w", err)
			}
			mode := args[3]

			params := types.Params{
				CredentialCollectionId: collectionId,
				CredentialTokenId:     tokenId,
				MinCredentialBalance:  minBalance,
				Mode:                  mode,
			}
			if err := params.Validate(); err != nil {
				return err
			}

			msg := &types.MsgUpdateParams{
				Authority: clientCtx.GetFromAddress().String(),
				Params:    params,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the pot module",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(CmdQueryParams())
	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query x/pot params",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Use ABCI store query since we have JSON-encoded params
			// (no protobuf gRPC gateway wired yet).
			resp, err := clientCtx.Client.ABCIQuery(cmd.Context(), fmt.Sprintf("/store/%s/key", types.StoreKey), types.ParamsKey)
			if err != nil {
				return fmt.Errorf("failed to query params: %w", err)
			}
			if resp.Response.Value == nil {
				fmt.Println("x/pot params: not set (module disabled)")
				return nil
			}
			var params types.Params
			if err := proto.Unmarshal(resp.Response.Value, &params); err != nil {
				// Try JSON fallback
				if err2 := json.Unmarshal(resp.Response.Value, &params); err2 != nil {
					return fmt.Errorf("failed to unmarshal params: proto: %w, json: %v", err, err2)
				}
			}
			out, _ := json.MarshalIndent(params, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
