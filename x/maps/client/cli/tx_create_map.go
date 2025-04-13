package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/spf13/cobra"

	"github.com/bitbadges/bitbadgeschain/x/maps/types"

	"github.com/cosmos/cosmos-sdk/client"
)

var _ = strconv.Itoa(0)

func CmdCreateMap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-map [tx-json]",
		Short: "Broadcast message createMap",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txJSON := args[0]

			var txData types.MsgCreateMap
			if err := jsonpb.UnmarshalString(txJSON, &txData); err != nil {
				return err
			}

			if err := txData.ValidateBasic(); err != nil {
				return err
			}

			txData.Creator = clientCtx.GetFromAddress().String()

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &txData)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSetValue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-value [tx-json]",
		Short: "Broadcast message setValue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txJSON := args[0]

			var txData types.MsgSetValue
			if err := jsonpb.UnmarshalString(txJSON, &txData); err != nil {
				return err
			}

			if err := txData.ValidateBasic(); err != nil {
				return err
			}

			txData.Creator = clientCtx.GetFromAddress().String()
			if txData.Options == nil {
				txData.Options = &types.SetOptions{}
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &txData)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateMap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-map [tx-json]",
		Short: "Broadcast message updateMap",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txJSON := args[0]

			var txData types.MsgUpdateMap
			if err := jsonpb.UnmarshalString(txJSON, &txData); err != nil {
				return err
			}

			if err := txData.ValidateBasic(); err != nil {
				return err
			}

			txData.Creator = clientCtx.GetFromAddress().String()

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &txData)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteMap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-map [tx-json]",
		Short: "Broadcast message deleteMap",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txJSON := args[0]

			var txData types.MsgDeleteMap
			if err := jsonpb.UnmarshalString(txJSON, &txData); err != nil {
				return err
			}

			if err := txData.ValidateBasic(); err != nil {
				return err
			}

			txData.Creator = clientCtx.GetFromAddress().String()

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &txData)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
