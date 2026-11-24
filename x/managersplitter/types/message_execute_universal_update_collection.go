package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgExecuteUniversalUpdateCollection) ValidateBasic() error {
	if msg == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "message cannot be nil")
	}
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

	// Validate that at least one update or create action is present
	tokenizationMsg := msg.UniversalUpdateCollectionMsg
	if tokenizationMsg.CollectionId.IsNil() {
		return sdkerrors.Wrap(ErrInvalidRequest, "collection ID cannot be nil in UniversalUpdateCollection message")
	}

	hasUpdateFlag := tokenizationMsg.CollectionId.IsZero() ||
		tokenizationMsg.UpdateValidTokenIds ||
		tokenizationMsg.UpdateCollectionPermissions ||
		tokenizationMsg.UpdateManager ||
		tokenizationMsg.UpdateCollectionMetadata ||
		tokenizationMsg.UpdateTokenMetadata ||
		tokenizationMsg.UpdateCustomData ||
		tokenizationMsg.UpdateCollectionApprovals ||
		tokenizationMsg.UpdateStandards ||
		tokenizationMsg.UpdateIsArchived ||
		len(tokenizationMsg.CosmosCoinWrapperPathsToAdd) > 0 ||
		len(tokenizationMsg.AliasPathsToAdd) > 0 ||
		len(tokenizationMsg.MintEscrowCoinsToTransfer) > 0

	if !hasUpdateFlag {
		return sdkerrors.Wrap(ErrInvalidRequest, "at least one update flag must be set in UniversalUpdateCollection message")
	}

	// Validate the nested UniversalUpdateCollection message
	// Note: We validate with the original creator, but it will be overwritten in the keeper
	// The creator field will be set to the manager splitter address in the keeper
	if err := tokenizationMsg.ValidateBasic(); err != nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid UniversalUpdateCollection message: %s", err)
	}

	if err := ValidateCanonicalAddresses(nil, msg.Executor, msg.ManagerSplitterAddress); err != nil {
		return err
	}

	return nil
}
