package cli

import (
	"encoding/json"
	"strconv"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

// message MsgNewCollection {
//     // See collections.proto for more details about these MsgNewBadge fields. Defines the badge details. Leave unneeded fields empty.
//     string creator = 1;
//     string collectionMetadata = 2;
//     repeated BadgeMetadata badgeMetadata = 3;

//     uint64 permissions = 4;
//     string bytes = 5;
//     repeated CollectionApprovedTransfer approvedTransfers = 6;
//     repeated CollectionApprovedTransfer managerApprovedTransfers = 7;
//     uint64 standard = 8;
//     //Badge supplys and amounts to create. For each idx, we create amounts[idx] badges each with a supply of supplys[idx].
//     //If supply[idx] == 0, we assume default supply. amountsToCreate[idx] can't equal 0.
//     repeated BadgeSupplyAndAmount badgesToCreate = 9;
//     repeated Transfers transfers = 10;
//     repeated Claim claims = 11;
// }

func CmdNewCollection() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new-collection [collectionMetadata] [badgeMetadata] [permissions] [bytes] [approvedTransfers] [managerApprovedTransfers] [standard] [supplys] [transfers] [claims] [offChainBalancesMetadata]",
		Short: "Broadcast message newCollection",
		Args:  cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argCollectionMetadata, err := cast.ToStringE(args[0])
			if err != nil {
				return err
			}

			var argBadgeMetadata []*types.BadgeMetadata
			err = json.Unmarshal([]byte(args[1]), &argBadgeMetadata)
			if err != nil {
				return err
			}

			argPermissions := types.NewUintFromString(args[2])
			if err != nil {
				return err
			}

			argBytesStr, err := cast.ToStringE(args[3])
			if err != nil {
				return err
			}

			var argApprovedTransfers []*types.CollectionApprovedTransfer
			err = json.Unmarshal([]byte(args[4]), &argApprovedTransfers)
			if err != nil {
				return err
			}

			var argManagerApprovedTransfers []*types.CollectionApprovedTransfer
			err = json.Unmarshal([]byte(args[5]), &argManagerApprovedTransfers)
			if err != nil {
				return err
			}

			argStandard := types.NewUintFromString(args[6])
			if err != nil {
				return err
			}

			var argBadgesToCreate []*types.BadgeSupplyAndAmount
			err = json.Unmarshal([]byte(args[7]), &argBadgesToCreate)
			if err != nil {
				return err
			}

			var argTransfers []*types.Transfer
			err = json.Unmarshal([]byte(args[8]), &argTransfers)
			if err != nil {
				return err
			}

			var argClaims []*types.Claim
			err = json.Unmarshal([]byte(args[9]), &argClaims)
			if err != nil {
				return err
			}

			argOffChainBalancesMetadata, err := cast.ToStringE(args[10])
			if err != nil {
				return err
			}

			msg := types.NewMsgNewCollection(
				clientCtx.GetFromAddress().String(),
				argStandard,
				argBadgesToCreate,
				argCollectionMetadata,
				argBadgeMetadata,
				argPermissions,
				argApprovedTransfers,
				argManagerApprovedTransfers,
				argBytesStr,
				argTransfers,
				argClaims,
				argOffChainBalancesMetadata,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
