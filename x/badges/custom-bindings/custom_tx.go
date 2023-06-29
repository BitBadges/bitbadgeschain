package custom_bindings

import (
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// fromReflectRawMsg decodes msg.Data to an sdk.Msg using proto Any and json encoding.
// this needs to be registered on the Encoders
func EncodeBadgeMessage() wasmKeeper.CustomEncoder {
	return func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {

		// return nil, sdkerrors.Wrapf(types.ErrUnknownRequest, "unmarshaled to value: %s", msg)

		var badgeCustomMsg badgeCustomMsg
		err := json.Unmarshal(msg, &badgeCustomMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		switch {
		// case badgeCustomMsg.NewCollection != nil:
		// 	newCollectionMsg := badgeTypes.NewMsgNewCollection(
		// 		sender.String(),
		// 		badgeCustomMsg.NewCollection.Standard,
		// 		badgeCustomMsg.NewCollection.BadgesToCreate,
		// 		badgeCustomMsg.NewCollection.CollectionMetadata,
		// 		badgeCustomMsg.NewCollection.BadgeMetadata,
		// 		badgeCustomMsg.NewCollection.Permissions,
		// 		badgeCustomMsg.NewCollection.ApprovedTransfers,
		// 		badgeCustomMsg.NewCollection.ManagerApprovedTransfers,
		// 		badgeCustomMsg.NewCollection.Bytes,
		// 		badgeCustomMsg.NewCollection.Transfers,
		// 		badgeCustomMsg.NewCollection.Claims,
		// 		badgeCustomMsg.NewCollection.OffChainBalancesMetadata,
		// 	)
		// 	return []sdk.Msg{newCollectionMsg}, nil
		// case badgeCustomMsg.MintAndDistributeBadges != nil:
		// 	MintAndDistributeBadgesMsg := badgeTypes.NewMsgMintAndDistributeBadges(
		// 		sender.String(),
		// 		badgeCustomMsg.MintAndDistributeBadges.CollectionId,
		// 		badgeCustomMsg.MintAndDistributeBadges.BadgesToCreate,
		// 		badgeCustomMsg.MintAndDistributeBadges.Transfers,
		// 		badgeCustomMsg.MintAndDistributeBadges.Claims,
		// 		badgeCustomMsg.MintAndDistributeBadges.CollectionMetadata,
		// 		badgeCustomMsg.MintAndDistributeBadges.BadgeMetadata,
		// 		badgeCustomMsg.MintAndDistributeBadges.OffChainBalancesMetadata,
		// 	)
		// 	return []sdk.Msg{MintAndDistributeBadgesMsg}, nil
		// case badgeCustomMsg.ClaimBadge != nil:
		// 	claimBadgeMsg := badgeTypes.NewMsgClaimBadge(
		// 		sender.String(),
		// 		badgeCustomMsg.ClaimBadge.ClaimId,
		// 		badgeCustomMsg.ClaimBadge.CollectionId,
		// 		badgeCustomMsg.ClaimBadge.Solutions,
		// 	)
		// 	return []sdk.Msg{claimBadgeMsg}, nil
		// case badgeCustomMsg.DeleteCollection != nil:
		// 	deleteCollectionMsg := badgeTypes.NewMsgDeleteCollection(
		// 		sender.String(),
		// 		badgeCustomMsg.DeleteCollection.CollectionId,
		// 	)
		// 	return []sdk.Msg{deleteCollectionMsg}, nil
		// case badgeCustomMsg.RequestUpdateManager != nil:
		// 	requestUpdateManagerMsg := badgeTypes.NewMsgRequestUpdateManager(
		// 		sender.String(),
		// 		badgeCustomMsg.RequestUpdateManager.CollectionId,
		// 		badgeCustomMsg.RequestUpdateManager.AddRequest,
		// 	)
		// 	return []sdk.Msg{requestUpdateManagerMsg}, nil
		// case badgeCustomMsg.SetApproval != nil:
		// 	setApprovalMsg := badgeTypes.NewMsgSetApproval(
		// 		sender.String(),
		// 		badgeCustomMsg.SetApproval.CollectionId,
		// 		badgeCustomMsg.SetApproval.Address,
		// 		badgeCustomMsg.SetApproval.Balances,
		// 	)
		// 	return []sdk.Msg{setApprovalMsg}, nil
		// case badgeCustomMsg.TransferBadge != nil:
		// 	transferBadgeMsg := badgeTypes.NewMsgTransferBadge(
		// 		sender.String(),
		// 		badgeCustomMsg.TransferBadge.CollectionId,
		// 		badgeCustomMsg.TransferBadge.From,
		// 		badgeCustomMsg.TransferBadge.Transfers,
		// 	)
		// 	return []sdk.Msg{transferBadgeMsg}, nil
		// case badgeCustomMsg.UpdateManager != nil:
		// 	updateManagerMsg := badgeTypes.NewMsgUpdateManager(
		// 		sender.String(),
		// 		badgeCustomMsg.UpdateManager.CollectionId,
		// 		badgeCustomMsg.UpdateManager.Address,
		// 	)
		// 	return []sdk.Msg{updateManagerMsg}, nil
		// case badgeCustomMsg.UpdateCustomData != nil:
		// 	updateCustomDataMsg := badgeTypes.NewMsgUpdateCustomData(
		// 		sender.String(),
		// 		badgeCustomMsg.UpdateCustomData.CollectionId,
		// 		badgeCustomMsg.UpdateCustomData.Bytes,
		// 	)
		// 	return []sdk.Msg{updateCustomDataMsg}, nil
		// case badgeCustomMsg.UpdateCollectionApprovedTransfers != nil:
		// 	updateCollectionMetadataMsg := badgeTypes.NewMsgUpdateCollectionApprovedTransfers(
		// 		sender.String(),
		// 		badgeCustomMsg.UpdateCollectionApprovedTransfers.CollectionId,
		// 		badgeCustomMsg.NewCollection.ApprovedTransfers,
		// 	)
		// 	return []sdk.Msg{updateCollectionMetadataMsg}, nil
		// case badgeCustomMsg.UpdateCollectionPermissions != nil:
		// 	updateCollectionPermissionsMsg := badgeTypes.NewMsgUpdateCollectionPermissions(
		// 		sender.String(),
		// 		badgeCustomMsg.UpdateCollectionPermissions.CollectionId,
		// 		badgeCustomMsg.UpdateCollectionPermissions.Permissions,
		// 	)
		// 	return []sdk.Msg{updateCollectionPermissionsMsg}, nil
		// case badgeCustomMsg.UpdateMetadata != nil:
		// 	updateMetadataMsg := badgeTypes.NewMsgUpdateMetadata(
		// 		sender.String(),
		// 		badgeCustomMsg.UpdateMetadata.CollectionId,
		// 		badgeCustomMsg.UpdateMetadata.CollectionMetadata,
		// 		badgeCustomMsg.UpdateMetadata.BadgeMetadata,
		// 		badgeCustomMsg.UpdateMetadata.OffChainBalancesMetadata,
		// 	)
		// 	return []sdk.Msg{updateMetadataMsg}, nil
		default:
			return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Unknown custom badge message variant %s", badgeCustomMsg)
		}
	}
}

type badgeCustomMsg struct {
	// NewCollection                     *badgeTypes.MsgNewCollection                     `json:"newCollectionMsg,omitempty"`
	// MintAndDistributeBadges                         *badgeTypes.MsgMintAndDistributeBadges                         `json:"mintAndDistributeBadgesMsg,omitempty"`
	// ClaimBadge                        *badgeTypes.MsgClaimBadge                        `json:"claimBadgeMsg,omitempty"`
	// DeleteCollection                  *badgeTypes.MsgDeleteCollection                  `json:"deleteCollectionMsg,omitempty"`
	// RequestUpdateManager            *badgeTypes.MsgRequestUpdateManager            `json:"requestUpdateManagerMsg,omitempty"`
	// SetApproval                       *badgeTypes.MsgSetApproval                       `json:"setApprovalMsg,omitempty"`
	// TransferBadge                     *badgeTypes.MsgTransferBadge                     `json:"transferBadgeMsg,omitempty"`
	// UpdateManager                   *badgeTypes.MsgUpdateManager                   `json:"updateManagerMsg,omitempty"`
	// UpdateCollectionApprovedTransfers *badgeTypes.MsgUpdateCollectionApprovedTransfers `json:"UpdateCollectionApprovedTransfersMsg,omitempty"`
	// UpdateCollectionPermissions                 *badgeTypes.MsgUpdateCollectionPermissions                 `json:"updateCollectionPermissionsMsg,omitempty"`
	// UpdateMetadata                    *badgeTypes.MsgUpdateMetadata                    `json:"updateMetadataMsg,omitempty"`
}
