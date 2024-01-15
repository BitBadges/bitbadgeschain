package custom_bindings

import (
	"bytes"
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/jsonpb"

	protocolTypes "github.com/bitbadges/bitbadgeschain/x/protocols/types"
)

// WASM handler for contracts calling into the badges module
func EncodeProtocolMessage() wasmKeeper.CustomEncoder {
	return func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {

		// Convert the RawMessage to a byte slice
		jsonData, err := msg.MarshalJSON()
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		//Create a reader from the byte slice
		reader := bytes.NewReader(jsonData)

		var badgeCustomMsg protocolTypes.ProtocolCustomMsgType
		err = jsonpb.Unmarshal(reader, &badgeCustomMsg)
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}
		switch {
			case badgeCustomMsg.CreateProtocolMsg != nil:
			badgeCustomMsg.CreateProtocolMsg.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.CreateProtocolMsg}, nil
		case badgeCustomMsg.UpdateProtocolMsg != nil:
			badgeCustomMsg.UpdateProtocolMsg.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.UpdateProtocolMsg}, nil
		case badgeCustomMsg.DeleteProtocolMsg != nil:
			badgeCustomMsg.DeleteProtocolMsg.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.DeleteProtocolMsg}, nil
		case badgeCustomMsg.SetCollectionForProtocolMsg != nil:
			badgeCustomMsg.SetCollectionForProtocolMsg.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.SetCollectionForProtocolMsg}, nil
		case badgeCustomMsg.UnsetCollectionForProtocolMsg != nil:
			badgeCustomMsg.UnsetCollectionForProtocolMsg.Creator = sender.String()
			return []sdk.Msg{badgeCustomMsg.UnsetCollectionForProtocolMsg}, nil
		default:
			return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Unknown custom badge message variant %s", badgeCustomMsg)
		}
	}
}
