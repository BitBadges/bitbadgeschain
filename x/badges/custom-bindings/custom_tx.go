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
		case badgeCustomMsg.CreateAddressMappings != nil:
			badgeCustomMsg.CreateAddressMappings.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.CreateAddressMappings}, nil
		case badgeCustomMsg.UpdateCollection != nil:
			badgeCustomMsg.UpdateCollection.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.UpdateCollection}, nil
		case badgeCustomMsg.DeleteCollection != nil:
			badgeCustomMsg.DeleteCollection.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.DeleteCollection}, nil
		case badgeCustomMsg.TransferBadges != nil:
			badgeCustomMsg.TransferBadges.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.TransferBadges}, nil
		case badgeCustomMsg.UpdateUserApprovedTransfers != nil:
			badgeCustomMsg.UpdateUserApprovedTransfers.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.UpdateUserApprovedTransfers}, nil
		default:
			return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Unknown custom badge message variant %s", badgeCustomMsg)
		}
	}
}

type badgeCustomMsg struct {
	CreateAddressMappings 					 *badgeTypes.MsgCreateAddressMappings             `json:"createAddressMappingsMsg,omitempty"`
	UpdateCollection   						*badgeTypes.MsgUpdateCollection                  `json:"updateCollectionMsg,omitempty"`
	DeleteCollection                  *badgeTypes.MsgDeleteCollection                  `json:"deleteCollectionMsg,omitempty"`
	TransferBadges                     *badgeTypes.MsgTransferBadges                     `json:"transferBadgesMsg,omitempty"`
	UpdateUserApprovedTransfers       *badgeTypes.MsgUpdateUserApprovedTransfers       `json:"updateUserApprovedTransfersMsg,omitempty"`
}
