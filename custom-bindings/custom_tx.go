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
	balancertypes "github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	managersplittertypes "github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	mapstypes "github.com/bitbadges/bitbadgeschain/x/maps/types"
	tokenTypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func EncodeBitBadgesModuleMessage() wasmKeeper.CustomEncoder {
	return func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
		// Convert message and route to corresponding handler
		jsonData, err := msg.MarshalJSON()
		if err != nil {
			return nil, sdkerrors.Wrap(err, err.Error())
		}

		reader := bytes.NewReader(jsonData)

		isTokenizationModuleMsg := false
		var tokenizationCustomMsg tokenTypes.TokenizationCustomMsgType
		err = jsonpb.Unmarshal(reader, &tokenizationCustomMsg)
		if err == nil {
			isTokenizationModuleMsg = true
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

		isManagersplitterModuleMsg := false
		var managersplitterCustomMsg managersplittertypes.ManagersplitterCustomMsgType
		err = jsonpb.Unmarshal(reader, &managersplitterCustomMsg)
		if err == nil {
			isManagersplitterModuleMsg = true
		}

		if isTokenizationModuleMsg {
			reader = bytes.NewReader(jsonData)
			var tokenizationCustomMsg tokenTypes.TokenizationCustomMsgType
			err = jsonpb.Unmarshal(reader, &tokenizationCustomMsg)
			if err != nil {
				return nil, sdkerrors.Wrap(err, err.Error())
			}

			switch {
			case tokenizationCustomMsg.UniversalUpdateCollectionMsg != nil:
				tokenizationCustomMsg.UniversalUpdateCollectionMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.UniversalUpdateCollectionMsg}, nil
			case tokenizationCustomMsg.CreateCollectionMsg != nil:
				tokenizationCustomMsg.CreateCollectionMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.CreateCollectionMsg}, nil
			case tokenizationCustomMsg.CreateAddressListsMsg != nil:
				tokenizationCustomMsg.CreateAddressListsMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.CreateAddressListsMsg}, nil
			case tokenizationCustomMsg.UpdateCollectionMsg != nil:
				tokenizationCustomMsg.UpdateCollectionMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.UpdateCollectionMsg}, nil
			case tokenizationCustomMsg.DeleteCollectionMsg != nil:
				tokenizationCustomMsg.DeleteCollectionMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.DeleteCollectionMsg}, nil
			case tokenizationCustomMsg.TransferTokensMsg != nil:
				tokenizationCustomMsg.TransferTokensMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.TransferTokensMsg}, nil
			case tokenizationCustomMsg.UpdateUserApprovalsMsg != nil:
				tokenizationCustomMsg.UpdateUserApprovalsMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.UpdateUserApprovalsMsg}, nil
			case tokenizationCustomMsg.CreateDynamicStoreMsg != nil:
				tokenizationCustomMsg.CreateDynamicStoreMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.CreateDynamicStoreMsg}, nil
			case tokenizationCustomMsg.UpdateDynamicStoreMsg != nil:
				tokenizationCustomMsg.UpdateDynamicStoreMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.UpdateDynamicStoreMsg}, nil
			case tokenizationCustomMsg.DeleteDynamicStoreMsg != nil:
				tokenizationCustomMsg.DeleteDynamicStoreMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.DeleteDynamicStoreMsg}, nil
			case tokenizationCustomMsg.SetDynamicStoreValueMsg != nil:
				tokenizationCustomMsg.SetDynamicStoreValueMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetDynamicStoreValueMsg}, nil
			case tokenizationCustomMsg.SetIncomingApprovalMsg != nil:
				tokenizationCustomMsg.SetIncomingApprovalMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetIncomingApprovalMsg}, nil
			case tokenizationCustomMsg.DeleteIncomingApprovalMsg != nil:
				tokenizationCustomMsg.DeleteIncomingApprovalMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.DeleteIncomingApprovalMsg}, nil
			case tokenizationCustomMsg.SetOutgoingApprovalMsg != nil:
				tokenizationCustomMsg.SetOutgoingApprovalMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetOutgoingApprovalMsg}, nil
			case tokenizationCustomMsg.DeleteOutgoingApprovalMsg != nil:
				tokenizationCustomMsg.DeleteOutgoingApprovalMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.DeleteOutgoingApprovalMsg}, nil
			case tokenizationCustomMsg.PurgeApprovalsMsg != nil:
				tokenizationCustomMsg.PurgeApprovalsMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.PurgeApprovalsMsg}, nil
			case tokenizationCustomMsg.SetValidTokenIdsMsg != nil:
				tokenizationCustomMsg.SetValidTokenIdsMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetValidTokenIdsMsg}, nil
			case tokenizationCustomMsg.SetManagerMsg != nil:
				tokenizationCustomMsg.SetManagerMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetManagerMsg}, nil
			case tokenizationCustomMsg.SetCollectionMetadataMsg != nil:
				tokenizationCustomMsg.SetCollectionMetadataMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetCollectionMetadataMsg}, nil
			case tokenizationCustomMsg.SetTokenMetadataMsg != nil:
				tokenizationCustomMsg.SetTokenMetadataMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetTokenMetadataMsg}, nil
			case tokenizationCustomMsg.SetCustomDataMsg != nil:
				tokenizationCustomMsg.SetCustomDataMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetCustomDataMsg}, nil
			case tokenizationCustomMsg.SetStandardsMsg != nil:
				tokenizationCustomMsg.SetStandardsMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetStandardsMsg}, nil
			case tokenizationCustomMsg.SetCollectionApprovalsMsg != nil:
				tokenizationCustomMsg.SetCollectionApprovalsMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetCollectionApprovalsMsg}, nil
			case tokenizationCustomMsg.SetIsArchivedMsg != nil:
				tokenizationCustomMsg.SetIsArchivedMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.SetIsArchivedMsg}, nil
			case tokenizationCustomMsg.CastVoteMsg != nil:
				tokenizationCustomMsg.CastVoteMsg.Creator = sender.String()
				return []sdk.Msg{tokenizationCustomMsg.CastVoteMsg}, nil
			default:
				return nil, sdkerrors.Wrapf(types.ErrInvalidMsg, "Unknown custom tokenization message variant %s", tokenizationCustomMsg)
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
		} else if isManagersplitterModuleMsg {
			reader = bytes.NewReader(jsonData)
			var managersplitterCustomMsg managersplittertypes.ManagersplitterCustomMsgType
			err = jsonpb.Unmarshal(reader, &managersplitterCustomMsg)
			if err != nil {
				return nil, sdkerrors.Wrap(err, err.Error())
			}

			switch {
			case managersplitterCustomMsg.CreateManagerSplitterMsg != nil:
				managersplitterCustomMsg.CreateManagerSplitterMsg.Admin = sender.String()
				return []sdk.Msg{managersplitterCustomMsg.CreateManagerSplitterMsg}, nil
			case managersplitterCustomMsg.UpdateManagerSplitterMsg != nil:
				managersplitterCustomMsg.UpdateManagerSplitterMsg.Admin = sender.String()
				return []sdk.Msg{managersplitterCustomMsg.UpdateManagerSplitterMsg}, nil
			case managersplitterCustomMsg.DeleteManagerSplitterMsg != nil:
				managersplitterCustomMsg.DeleteManagerSplitterMsg.Admin = sender.String()
				return []sdk.Msg{managersplitterCustomMsg.DeleteManagerSplitterMsg}, nil
			case managersplitterCustomMsg.ExecuteUniversalUpdateCollectionMsg != nil:
				managersplitterCustomMsg.ExecuteUniversalUpdateCollectionMsg.Executor = sender.String()
				return []sdk.Msg{managersplitterCustomMsg.ExecuteUniversalUpdateCollectionMsg}, nil
			}
		}

		return nil, sdkerrors.Wrap(types.ErrInvalidMsg, "Unknown custom badge message variant")
	}
}
