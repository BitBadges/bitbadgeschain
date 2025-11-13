package custom_bindings

import (
	"bytes"
	"encoding/json"

	sdkerrors "cosmossdk.io/errors"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/jsonpb"

	anchortypes "github.com/bitbadges/bitbadgeschain/x/anchor/types"
	tokenTypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	balancertypes "github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	mapstypes "github.com/bitbadges/bitbadgeschain/x/maps/types"
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
		var badgeCustomMsg tokenTypes.BadgeCustomMsgType
		err = jsonpb.Unmarshal(reader, &badgeCustomMsg)
		if err == nil {
			isBadgeModuleMsg = true
		}

		isAnchorModuleMsg := false
		var anchorCustomMsg anchortypes.AnchorCustomMsgType
		err = jsonpb.Unmarshal(reader, &anchorCustomMsg)
		if err == nil {
			isAnchorModuleMsg = true
		}

		isMapsModuleMsg := false
		var mapsCustomMsg mapstypes.MapCustomMsgType
		err = jsonpb.Unmarshal(reader, &mapsCustomMsg)
		if err == nil {
			isMapsModuleMsg = true
		}

		isGammModuleMsg := false
		var gammCustomMsg gammtypes.GammCustomMsgType
		err = jsonpb.Unmarshal(reader, &gammCustomMsg)
		if err == nil {
			isGammModuleMsg = true
		}

		isBalancerModuleMsg := false
		var balancerCustomMsg balancertypes.BalancerCustomMsgType
		err = jsonpb.Unmarshal(reader, &balancerCustomMsg)
		if err == nil {
			isBalancerModuleMsg = true
		}

		if isBadgeModuleMsg {
			reader = bytes.NewReader(jsonData)
			var badgeCustomMsg tokenTypes.BadgeCustomMsgType
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
			case badgeCustomMsg.TransferTokensMsg != nil:
				badgeCustomMsg.TransferTokensMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.TransferTokensMsg}, nil
			case badgeCustomMsg.UpdateUserApprovalsMsg != nil:
				badgeCustomMsg.UpdateUserApprovalsMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.UpdateUserApprovalsMsg}, nil
			case badgeCustomMsg.CreateDynamicStoreMsg != nil:
				badgeCustomMsg.CreateDynamicStoreMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.CreateDynamicStoreMsg}, nil
			case badgeCustomMsg.UpdateDynamicStoreMsg != nil:
				badgeCustomMsg.UpdateDynamicStoreMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.UpdateDynamicStoreMsg}, nil
			case badgeCustomMsg.DeleteDynamicStoreMsg != nil:
				badgeCustomMsg.DeleteDynamicStoreMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.DeleteDynamicStoreMsg}, nil
			case badgeCustomMsg.SetDynamicStoreValueMsg != nil:
				badgeCustomMsg.SetDynamicStoreValueMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetDynamicStoreValueMsg}, nil
			case badgeCustomMsg.IncrementStoreValueMsg != nil:
				badgeCustomMsg.IncrementStoreValueMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.IncrementStoreValueMsg}, nil
			case badgeCustomMsg.DecrementStoreValueMsg != nil:
				badgeCustomMsg.DecrementStoreValueMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.DecrementStoreValueMsg}, nil
			case badgeCustomMsg.SetIncomingApprovalMsg != nil:
				badgeCustomMsg.SetIncomingApprovalMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetIncomingApprovalMsg}, nil
			case badgeCustomMsg.DeleteIncomingApprovalMsg != nil:
				badgeCustomMsg.DeleteIncomingApprovalMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.DeleteIncomingApprovalMsg}, nil
			case badgeCustomMsg.SetOutgoingApprovalMsg != nil:
				badgeCustomMsg.SetOutgoingApprovalMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetOutgoingApprovalMsg}, nil
			case badgeCustomMsg.DeleteOutgoingApprovalMsg != nil:
				badgeCustomMsg.DeleteOutgoingApprovalMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.DeleteOutgoingApprovalMsg}, nil
			case badgeCustomMsg.PurgeApprovalsMsg != nil:
				badgeCustomMsg.PurgeApprovalsMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.PurgeApprovalsMsg}, nil
			case badgeCustomMsg.SetValidTokenIdsMsg != nil:
				badgeCustomMsg.SetValidTokenIdsMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetValidTokenIdsMsg}, nil
			case badgeCustomMsg.SetManagerMsg != nil:
				badgeCustomMsg.SetManagerMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetManagerMsg}, nil
			case badgeCustomMsg.SetCollectionMetadataMsg != nil:
				badgeCustomMsg.SetCollectionMetadataMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetCollectionMetadataMsg}, nil
			case badgeCustomMsg.SetTokenMetadataMsg != nil:
				badgeCustomMsg.SetTokenMetadataMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetTokenMetadataMsg}, nil
			case badgeCustomMsg.SetCustomDataMsg != nil:
				badgeCustomMsg.SetCustomDataMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetCustomDataMsg}, nil
			case badgeCustomMsg.SetStandardsMsg != nil:
				badgeCustomMsg.SetStandardsMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetStandardsMsg}, nil
			case badgeCustomMsg.SetCollectionApprovalsMsg != nil:
				badgeCustomMsg.SetCollectionApprovalsMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetCollectionApprovalsMsg}, nil
			case badgeCustomMsg.SetIsArchivedMsg != nil:
				badgeCustomMsg.SetIsArchivedMsg.Creator = sender.String()
				return []sdk.Msg{badgeCustomMsg.SetIsArchivedMsg}, nil
			default:
				return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Unknown custom badge message variant %s", badgeCustomMsg)
			}
		} else if isAnchorModuleMsg {
			reader = bytes.NewReader(jsonData)
			var anchorCustomMsg anchortypes.AnchorCustomMsgType
			err = jsonpb.Unmarshal(reader, &anchorCustomMsg)
			if err != nil {
				return nil, sdkerrors.Wrap(err, err.Error())
			}

			switch {
			case anchorCustomMsg.AddCustomDataMsg != nil:
				anchorCustomMsg.AddCustomDataMsg.Creator = sender.String()
				return []sdk.Msg{anchorCustomMsg.AddCustomDataMsg}, nil
			}
		} else if isMapsModuleMsg {
			reader = bytes.NewReader(jsonData)
			var mapsCustomMsg mapstypes.MapCustomMsgType
			err = jsonpb.Unmarshal(reader, &mapsCustomMsg)
			if err != nil {
				return nil, sdkerrors.Wrap(err, err.Error())
			}

			switch {
			case mapsCustomMsg.CreateMapMsg != nil:
				mapsCustomMsg.CreateMapMsg.Creator = sender.String()
				return []sdk.Msg{mapsCustomMsg.CreateMapMsg}, nil
			case mapsCustomMsg.UpdateMapMsg != nil:
				mapsCustomMsg.UpdateMapMsg.Creator = sender.String()
				return []sdk.Msg{mapsCustomMsg.UpdateMapMsg}, nil
			case mapsCustomMsg.DeleteMapMsg != nil:
				mapsCustomMsg.DeleteMapMsg.Creator = sender.String()
				return []sdk.Msg{mapsCustomMsg.DeleteMapMsg}, nil
			case mapsCustomMsg.SetValueMsg != nil:
				mapsCustomMsg.SetValueMsg.Creator = sender.String()
				return []sdk.Msg{mapsCustomMsg.SetValueMsg}, nil
			}
		} else if isGammModuleMsg {
			reader = bytes.NewReader(jsonData)
			var gammCustomMsg gammtypes.GammCustomMsgType
			err = jsonpb.Unmarshal(reader, &gammCustomMsg)
			if err != nil {
				return nil, sdkerrors.Wrap(err, err.Error())
			}

			switch {
			case gammCustomMsg.JoinPoolMsg != nil:
				gammCustomMsg.JoinPoolMsg.Sender = sender.String()
				return []sdk.Msg{gammCustomMsg.JoinPoolMsg}, nil
			case gammCustomMsg.ExitPoolMsg != nil:
				gammCustomMsg.ExitPoolMsg.Sender = sender.String()
				return []sdk.Msg{gammCustomMsg.ExitPoolMsg}, nil
			case gammCustomMsg.SwapExactAmountInMsg != nil:
				gammCustomMsg.SwapExactAmountInMsg.Sender = sender.String()
				return []sdk.Msg{gammCustomMsg.SwapExactAmountInMsg}, nil
			case gammCustomMsg.SwapExactAmountOutMsg != nil:
				gammCustomMsg.SwapExactAmountOutMsg.Sender = sender.String()
				return []sdk.Msg{gammCustomMsg.SwapExactAmountOutMsg}, nil
			case gammCustomMsg.JoinSwapExternAmountInMsg != nil:
				gammCustomMsg.JoinSwapExternAmountInMsg.Sender = sender.String()
				return []sdk.Msg{gammCustomMsg.JoinSwapExternAmountInMsg}, nil
			case gammCustomMsg.JoinSwapShareAmountOutMsg != nil:
				gammCustomMsg.JoinSwapShareAmountOutMsg.Sender = sender.String()
				return []sdk.Msg{gammCustomMsg.JoinSwapShareAmountOutMsg}, nil
			case gammCustomMsg.ExitSwapShareAmountInMsg != nil:
				gammCustomMsg.ExitSwapShareAmountInMsg.Sender = sender.String()
				return []sdk.Msg{gammCustomMsg.ExitSwapShareAmountInMsg}, nil
			case gammCustomMsg.ExitSwapExternAmountOutMsg != nil:
				gammCustomMsg.ExitSwapExternAmountOutMsg.Sender = sender.String()
				return []sdk.Msg{gammCustomMsg.ExitSwapExternAmountOutMsg}, nil
			case gammCustomMsg.SwapExactAmountInWithIBCTransferMsg != nil:
				gammCustomMsg.SwapExactAmountInWithIBCTransferMsg.Sender = sender.String()
				return []sdk.Msg{gammCustomMsg.SwapExactAmountInWithIBCTransferMsg}, nil
			}
		} else if isBalancerModuleMsg {
			reader = bytes.NewReader(jsonData)
			var balancerCustomMsg balancertypes.BalancerCustomMsgType
			err = jsonpb.Unmarshal(reader, &balancerCustomMsg)
			if err != nil {
				return nil, sdkerrors.Wrap(err, err.Error())
			}

			switch {
			case balancerCustomMsg.CreateBalancerPoolMsg != nil:
				balancerCustomMsg.CreateBalancerPoolMsg.Sender = sender.String()
				return []sdk.Msg{balancerCustomMsg.CreateBalancerPoolMsg}, nil
			}
		}

		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown custom badge message variant")
	}
}
