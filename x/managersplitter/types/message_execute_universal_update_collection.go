package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgExecuteUniversalUpdateCollection) ValidateBasic() error {
	// Validate executor address format
	executorAddr, err := sdk.AccAddressFromBech32(msg.Executor)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid executor address (%s)", err)
	}
	if executorAddr.Empty() {
		return sdkerrors.Wrap(ErrInvalidAddress, "executor address cannot be empty")
	}

	// Validate manager splitter address format
	if msg.ManagerSplitterAddress == "" {
		return sdkerrors.Wrap(ErrInvalidAddress, "manager splitter address cannot be empty")
	}
	managerSplitterAddr, err := sdk.AccAddressFromBech32(msg.ManagerSplitterAddress)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid manager splitter address format (%s)", err)
	}
	if managerSplitterAddr.Empty() {
		return sdkerrors.Wrap(ErrInvalidAddress, "manager splitter address cannot be empty")
	}

	// Validate that executor and manager splitter are different addresses
	// (executor is the signer, manager splitter is the contract address)
	if msg.Executor == msg.ManagerSplitterAddress {
		return sdkerrors.Wrap(ErrInvalidRequest, "executor and manager splitter address cannot be the same")
	}

	// Validate that UniversalUpdateCollectionMsg is not nil
	if msg.UniversalUpdateCollectionMsg == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "universal update collection message cannot be nil")
	}

	// Validate that at least one update flag is set
	badgesMsg := msg.UniversalUpdateCollectionMsg
	hasUpdateFlag := badgesMsg.UpdateValidTokenIds ||
		badgesMsg.UpdateCollectionPermissions ||
		badgesMsg.UpdateManager ||
		badgesMsg.UpdateCollectionMetadata ||
		badgesMsg.UpdateTokenMetadata ||
		badgesMsg.UpdateCustomData ||
		badgesMsg.UpdateCollectionApprovals ||
		badgesMsg.UpdateStandards ||
		badgesMsg.UpdateIsArchived

	if !hasUpdateFlag {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one update flag must be set in UniversalUpdateCollection message")
	}

	// Validate the nested UniversalUpdateCollection message
	// Note: We validate with the original creator, but it will be overwritten in the keeper
	// The creator field will be set to the manager splitter address in the keeper
	if err := badgesMsg.ValidateBasic(); err != nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid UniversalUpdateCollection message: %s", err)
	}

	// Additional validation: ensure collection ID is valid (not nil)
	if badgesMsg.CollectionId.IsNil() {
		return sdkerrors.Wrap(ErrInvalidRequest, "collection ID cannot be nil in UniversalUpdateCollection message")
	}

	return nil
}
