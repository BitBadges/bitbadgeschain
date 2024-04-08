package custom_bindings

import (
	"bytes"
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/jsonpb"

	badgeTypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func EncodeBitBadgesModuleMessage() wasmKeeper.CustomEncoder {
	return func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		// Convert message and route to corresponding handler
		jsonData, err := msg.MarshalJSON()
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		reader := bytes.NewReader(jsonData)

		isBadgeModuleMsg := false
		var badgeCustomMsg badgeTypes.BadgeCustomMsgType
		err = jsonpb.Unmarshal(reader, &badgeCustomMsg)
		if err == nil {
			isBadgeModuleMsg = true
		}

		if isBadgeModuleMsg {
			reader = bytes.NewReader(jsonData)
			var badgeCustomMsg badgeTypes.BadgeCustomMsgType
			err = jsonpb.Unmarshal(reader, &badgeCustomMsg)
			if err != nil {
				return nil, sdkerrors.Wrap(err, err.Error())
			}

			switch {
			case badgeCustomMsg.UniversalUpdateCollectionMsg != nil:
				badgeCustomMsg.UniversalUpdateCollectionMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.UniversalUpdateCollectionMsg}, nil
			case badgeCustomMsg.CreateCollectionMsg != nil:
				badgeCustomMsg.CreateCollectionMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.CreateCollectionMsg}, nil
			case badgeCustomMsg.CreateAddressListsMsg != nil:
				badgeCustomMsg.CreateAddressListsMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.CreateAddressListsMsg}, nil
			case badgeCustomMsg.UpdateCollectionMsg != nil:
				badgeCustomMsg.UpdateCollectionMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.UpdateCollectionMsg}, nil
			case badgeCustomMsg.DeleteCollectionMsg != nil:
				badgeCustomMsg.DeleteCollectionMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.DeleteCollectionMsg}, nil
			case badgeCustomMsg.TransferBadgesMsg != nil:
				badgeCustomMsg.TransferBadgesMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.TransferBadgesMsg}, nil
			case badgeCustomMsg.UpdateUserApprovalsMsg != nil:
				badgeCustomMsg.UpdateUserApprovalsMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.UpdateUserApprovalsMsg}, nil
			default:
				return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Unknown custom badge message variant %s", badgeCustomMsg)
			}
		}

		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown custom badge message variant")
	}
}
