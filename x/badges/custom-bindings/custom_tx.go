package custom_bindings

import (
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	badgeTypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// fromReflectRawMsg decodes msg.Data to an sdk.Msg using proto Any and json encoding.
// this needs to be registered on the Encoders
func EncodeBadgeMessage() wasmKeeper.CustomEncoder {
	return func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		var badgeCustomMsg badgeCustomMsg
		err := json.Unmarshal(msg, &badgeCustomMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		switch {
		case badgeCustomMsg.NewCollection != nil:
			badgeCustomMsg.NewCollection.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.NewCollection}, nil
		case badgeCustomMsg.MintAndDistributeBadges != nil:
			badgeCustomMsg.MintAndDistributeBadges.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.MintAndDistributeBadges}, nil
		case badgeCustomMsg.DeleteCollection != nil:
			badgeCustomMsg.DeleteCollection.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.DeleteCollection}, nil
		case badgeCustomMsg.TransferBadge != nil:
			badgeCustomMsg.TransferBadge.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.TransferBadge}, nil
		case badgeCustomMsg.UpdateManager != nil:
			badgeCustomMsg.UpdateManager.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.UpdateManager}, nil
		case badgeCustomMsg.UpdateCollectionApprovedTransfers != nil:
			badgeCustomMsg.UpdateCollectionApprovedTransfers.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.UpdateCollectionApprovedTransfers}, nil
		case badgeCustomMsg.UpdateCollectionPermissions != nil:
			badgeCustomMsg.UpdateCollectionPermissions.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.UpdateCollectionPermissions}, nil
		case badgeCustomMsg.UpdateMetadata != nil:
			badgeCustomMsg.UpdateMetadata.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.UpdateMetadata}, nil
		case badgeCustomMsg.ArchiveCollection != nil:
			badgeCustomMsg.ArchiveCollection.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.ArchiveCollection}, nil
		case badgeCustomMsg.UpdateUserApprovedTransfers != nil:
			badgeCustomMsg.UpdateUserApprovedTransfers.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.UpdateUserApprovedTransfers}, nil
		case badgeCustomMsg.UpdateUserPermissions != nil:
			badgeCustomMsg.UpdateUserPermissions.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.UpdateUserPermissions}, nil
		default:
			return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Unknown custom badge message variant %s", badgeCustomMsg)
		}
	}
}

type badgeCustomMsg struct {
	NewCollection                     *badgeTypes.MsgNewCollection                     `json:"newCollectionMsg,omitempty"`
	MintAndDistributeBadges           *badgeTypes.MsgMintAndDistributeBadges                         `json:"mintAndDistributeBadgesMsg,omitempty"`
	DeleteCollection                  *badgeTypes.MsgDeleteCollection                  `json:"deleteCollectionMsg,omitempty"`
	TransferBadge                     *badgeTypes.MsgTransferBadge                     `json:"transferBadgeMsg,omitempty"`
	UpdateManager                   *badgeTypes.MsgUpdateManager                   `json:"updateManagerMsg,omitempty"`
	UpdateCollectionApprovedTransfers *badgeTypes.MsgUpdateCollectionApprovedTransfers `json:"UpdateCollectionApprovedTransfersMsg,omitempty"`
	UpdateCollectionPermissions                 *badgeTypes.MsgUpdateCollectionPermissions                 `json:"updateCollectionPermissionsMsg,omitempty"`
	UpdateMetadata                    *badgeTypes.MsgUpdateMetadata                    `json:"updateMetadataMsg,omitempty"`
	ArchiveCollection 							 *badgeTypes.MsgArchiveCollection                 `json:"archiveCollectionMsg,omitempty"`
	UpdateUserApprovedTransfers      *badgeTypes.MsgUpdateUserApprovedTransfers      `json:"updateUserApprovedTransfersMsg,omitempty"`
	UpdateUserPermissions 					*badgeTypes.MsgUpdateUserPermissions            `json:"updateUserPermissionsMsg,omitempty"`
}
