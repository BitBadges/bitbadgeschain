package custom_bindings

import (
	"encoding/json"

	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	badgeTypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// fromReflectRawMsg decodes msg.Data to an sdk.Msg using proto Any and json encoding.
// this needs to be registered on the Encoders
func EncodeBadgeMessage() wasmKeeper.CustomEncoder {
	return func(_sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {

		// return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unmarshaled to value: %s", msg)

		var badgeCustomMsg badgeCustomMsg
		err := json.Unmarshal(msg, &badgeCustomMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		// return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unmarshaled to value: %+v", badgeCustomMsg)

		switch {
		case badgeCustomMsg.RegisterAddresses != nil:
				registerAddressesMsg := badgeTypes.NewMsgRegisterAddresses(
					_sender.String(),
					badgeCustomMsg.RegisterAddresses.AddressesToRegister,
				)
				return []sdk.Msg{registerAddressesMsg}, nil
		case badgeCustomMsg.NewCollection != nil:
				newCollectionMsg := badgeTypes.NewMsgNewCollection(
					_sender.String(),
					badgeCustomMsg.NewCollection.Standard,
					badgeCustomMsg.NewCollection.BadgeSupplys,
					badgeCustomMsg.NewCollection.CollectionUri,
					badgeCustomMsg.NewCollection.BadgeUris,
					badgeCustomMsg.NewCollection.Permissions,
					badgeCustomMsg.NewCollection.DisallowedTransfers,
					badgeCustomMsg.NewCollection.ManagerApprovedTransfers,
					badgeCustomMsg.NewCollection.Bytes,
					badgeCustomMsg.NewCollection.Transfers,
					badgeCustomMsg.NewCollection.Claims,
				)
				return []sdk.Msg{newCollectionMsg}, nil
		case badgeCustomMsg.MintBadge != nil:
				mintBadgeMsg := badgeTypes.NewMsgMintBadge(
					_sender.String(),
					badgeCustomMsg.MintBadge.CollectionId,
					badgeCustomMsg.MintBadge.BadgeSupplys,
					badgeCustomMsg.MintBadge.Transfers,
					badgeCustomMsg.MintBadge.Claims,
					badgeCustomMsg.MintBadge.CollectionUri,
					badgeCustomMsg.MintBadge.BadgeUris,
				)
				return []sdk.Msg{mintBadgeMsg}, nil
		case badgeCustomMsg.ClaimBadge != nil:
				claimBadgeMsg := badgeTypes.NewMsgClaimBadge(
					_sender.String(),
					badgeCustomMsg.ClaimBadge.ClaimId,
					badgeCustomMsg.ClaimBadge.CollectionId,
					badgeCustomMsg.ClaimBadge.WhitelistProof,
					badgeCustomMsg.ClaimBadge.CodeProof,
				)
				return []sdk.Msg{claimBadgeMsg}, nil
		case badgeCustomMsg.DeleteCollection != nil:
				deleteCollectionMsg := badgeTypes.NewMsgDeleteCollection(
					_sender.String(),
					badgeCustomMsg.DeleteCollection.CollectionId,
				)
				return []sdk.Msg{deleteCollectionMsg}, nil
		case badgeCustomMsg.RequestTransferManager != nil:
				requestTransferManagerMsg := badgeTypes.NewMsgRequestTransferManager(
					_sender.String(),
					badgeCustomMsg.RequestTransferManager.CollectionId,
					badgeCustomMsg.RequestTransferManager.AddRequest,
				)
				return []sdk.Msg{requestTransferManagerMsg}, nil
		case badgeCustomMsg.SetApproval != nil:
				setApprovalMsg := badgeTypes.NewMsgSetApproval(
					_sender.String(),
					badgeCustomMsg.SetApproval.CollectionId,
					badgeCustomMsg.SetApproval.Address,
					badgeCustomMsg.SetApproval.Balances,
				)
				return []sdk.Msg{setApprovalMsg}, nil
		case badgeCustomMsg.TransferBadge != nil:
				transferBadgeMsg := badgeTypes.NewMsgTransferBadge(
					_sender.String(),
					badgeCustomMsg.TransferBadge.CollectionId,
					badgeCustomMsg.TransferBadge.From,
					badgeCustomMsg.TransferBadge.Transfers,
				)
				return []sdk.Msg{transferBadgeMsg}, nil
		case badgeCustomMsg.TransferManager != nil:
				transferManagerMsg := badgeTypes.NewMsgTransferManager(
					_sender.String(),
					badgeCustomMsg.TransferManager.CollectionId,
					badgeCustomMsg.TransferManager.Address,
				)
				return []sdk.Msg{transferManagerMsg}, nil
		case badgeCustomMsg.UpdateBytes != nil:
				updateBytesMsg := badgeTypes.NewMsgUpdateBytes(
					_sender.String(),
					badgeCustomMsg.UpdateBytes.CollectionId,
					badgeCustomMsg.UpdateBytes.NewBytes,
				)
				return []sdk.Msg{updateBytesMsg}, nil
		case badgeCustomMsg.UpdateDisallowedTransfers != nil:
				updateCollectionUriMsg := badgeTypes.NewMsgUpdateDisallowedTransfers(
					_sender.String(),
					badgeCustomMsg.UpdateDisallowedTransfers.CollectionId,
					badgeCustomMsg.NewCollection.DisallowedTransfers,
				)
				return []sdk.Msg{updateCollectionUriMsg}, nil
		case badgeCustomMsg.UpdatePermissions != nil:
				updatePermissionsMsg := badgeTypes.NewMsgUpdatePermissions(
					_sender.String(),
					badgeCustomMsg.UpdatePermissions.CollectionId,
					badgeCustomMsg.UpdatePermissions.Permissions,
				)
				return []sdk.Msg{updatePermissionsMsg}, nil
		case badgeCustomMsg.UpdateUris != nil:
				updateUrisMsg := badgeTypes.NewMsgUpdateUris(
					_sender.String(),
					badgeCustomMsg.UpdateUris.CollectionId,
					badgeCustomMsg.UpdateUris.CollectionUri,
					badgeCustomMsg.UpdateUris.BadgeUris,
				)
				return []sdk.Msg{updateUrisMsg}, nil
		default:
			return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Unknown custom badge message variant %s", badgeCustomMsg)
		}
	}
}

type badgeCustomMsg struct {
	RegisterAddresses *badgeTypes.MsgRegisterAddresses `json:"registerAddressesMsg,omitempty"`
	NewCollection    *badgeTypes.MsgNewCollection     `json:"newCollectionMsg,omitempty"`
	MintBadge		*badgeTypes.MsgMintBadge         `json:"mintBadgeMsg,omitempty"`
	ClaimBadge 	*badgeTypes.MsgClaimBadge        `json:"claimBadgeMsg,omitempty"`
	DeleteCollection *badgeTypes.MsgDeleteCollection  `json:"deleteCollectionMsg,omitempty"`
	RequestTransferManager *badgeTypes.MsgRequestTransferManager `json:"requestTransferManagerMsg,omitempty"`
	SetApproval *badgeTypes.MsgSetApproval `json:"setApprovalMsg,omitempty"`
	TransferBadge *badgeTypes.MsgTransferBadge `json:"transferBadgeMsg,omitempty"`
	TransferManager *badgeTypes.MsgTransferManager `json:"transferManagerMsg,omitempty"`
	UpdateBytes *badgeTypes.MsgUpdateBytes `json:"updateBytesMsg,omitempty"`
	UpdateDisallowedTransfers *badgeTypes.MsgUpdateDisallowedTransfers `json:"updateDisallowedTransfersMsg,omitempty"`
	UpdatePermissions *badgeTypes.MsgUpdatePermissions `json:"updatePermissionsMsg,omitempty"`
	UpdateUris *badgeTypes.MsgUpdateUris `json:"updateUrisMsg,omitempty"`
}