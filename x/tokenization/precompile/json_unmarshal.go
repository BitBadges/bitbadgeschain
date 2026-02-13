package tokenization

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// convertEVMAddressToBech32 converts an EVM address (0x...) to bech32 format if needed
// If the address is already in bech32 format, it returns it unchanged
// If it's an EVM address, it converts it to bech32
func convertEVMAddressToBech32(addr string) string {
	if addr == "" {
		return addr
	}
	// Check if it's already a bech32 address (starts with prefix like "bb1", "cosmos1", etc.)
	// Bech32 addresses typically start with a prefix followed by "1"
	if _, err := sdk.AccAddressFromBech32(addr); err == nil {
		// Already valid bech32, return as-is
		return addr
	}
	// Check if it's an EVM address (0x...)
	if common.IsHexAddress(addr) {
		evmAddr := common.HexToAddress(addr)
		return sdk.AccAddress(evmAddr.Bytes()).String()
	}
	// If neither, return as-is (will fail validation later)
	return addr
}

// unmarshalMsgFromJSON unmarshals a JSON string into the appropriate Msg type based on method name
// and sets the Creator field from the contract caller for security.
func (p Precompile) unmarshalMsgFromJSON(methodName string, jsonStr string, contract *vm.Contract) (sdk.Msg, error) {
	// Get caller address
	caller := contract.Caller()
	if err := VerifyCaller(caller); err != nil {
		return nil, err
	}
	creatorCosmosAddr := sdk.AccAddress(caller.Bytes()).String()

	// Create the appropriate Msg type based on method name
	var msg sdk.Msg
	switch methodName {
	case TransferTokensMethod:
		msg = &tokenizationtypes.MsgTransferTokens{}
	case SetIncomingApprovalMethod:
		msg = &tokenizationtypes.MsgSetIncomingApproval{}
	case SetOutgoingApprovalMethod:
		msg = &tokenizationtypes.MsgSetOutgoingApproval{}
	case CreateCollectionMethod:
		msg = &tokenizationtypes.MsgCreateCollection{}
	case UpdateCollectionMethod:
		msg = &tokenizationtypes.MsgUpdateCollection{}
	case DeleteCollectionMethod:
		msg = &tokenizationtypes.MsgDeleteCollection{}
	case CreateAddressListsMethod:
		msg = &tokenizationtypes.MsgCreateAddressLists{}
	case UpdateUserApprovalsMethod:
		msg = &tokenizationtypes.MsgUpdateUserApprovals{}
	case DeleteIncomingApprovalMethod:
		msg = &tokenizationtypes.MsgDeleteIncomingApproval{}
	case DeleteOutgoingApprovalMethod:
		msg = &tokenizationtypes.MsgDeleteOutgoingApproval{}
	case PurgeApprovalsMethod:
		msg = &tokenizationtypes.MsgPurgeApprovals{}
	case CreateDynamicStoreMethod:
		msg = &tokenizationtypes.MsgCreateDynamicStore{}
	case UpdateDynamicStoreMethod:
		msg = &tokenizationtypes.MsgUpdateDynamicStore{}
	case DeleteDynamicStoreMethod:
		msg = &tokenizationtypes.MsgDeleteDynamicStore{}
	case SetDynamicStoreValueMethod:
		msg = &tokenizationtypes.MsgSetDynamicStoreValue{}
	case SetValidTokenIdsMethod:
		msg = &tokenizationtypes.MsgSetValidTokenIds{}
	case SetManagerMethod:
		msg = &tokenizationtypes.MsgSetManager{}
	case SetCollectionMetadataMethod:
		msg = &tokenizationtypes.MsgSetCollectionMetadata{}
	case SetTokenMetadataMethod:
		msg = &tokenizationtypes.MsgSetTokenMetadata{}
	case SetCustomDataMethod:
		msg = &tokenizationtypes.MsgSetCustomData{}
	case SetStandardsMethod:
		msg = &tokenizationtypes.MsgSetStandards{}
	case SetCollectionApprovalsMethod:
		msg = &tokenizationtypes.MsgSetCollectionApprovals{}
	case SetIsArchivedMethod:
		msg = &tokenizationtypes.MsgSetIsArchived{}
	case CastVoteMethod:
		msg = &tokenizationtypes.MsgCastVote{}
	case UniversalUpdateCollectionMethod:
		msg = &tokenizationtypes.MsgUniversalUpdateCollection{}
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unknown method: %s", methodName))
	}

	// Unmarshal JSON into the message
	// Try direct protobuf JSON unmarshaling first (ModuleCdc handles both standard JSON and protobuf JSON)
	// If that fails, try to provide better error messages
	if err := tokenizationtypes.ModuleCdc.UnmarshalJSON([]byte(jsonStr), msg); err != nil {
		// Validate JSON syntax separately for better error messages
		var jsonMap map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(jsonStr), &jsonMap); jsonErr != nil {
			// JSON syntax is invalid
			return nil, ErrInvalidInput(fmt.Sprintf("invalid JSON syntax: %v", jsonErr))
		}
		// JSON syntax is valid but protobuf unmarshaling failed
		// Provide more detailed error information
		return nil, ErrInvalidInput(fmt.Sprintf("failed to unmarshal JSON into %T: %v. JSON was: %s", msg, err, jsonStr))
	}

	// Set Creator field from contract caller (security: override any value in JSON)
	// Also convert EVM addresses to bech32 format for better UX
	switch m := msg.(type) {
	case *tokenizationtypes.MsgTransferTokens:
		m.Creator = creatorCosmosAddr
		// Convert ToAddresses in transfers
		for _, transfer := range m.Transfers {
			if transfer != nil {
				for i, toAddr := range transfer.ToAddresses {
					transfer.ToAddresses[i] = convertEVMAddressToBech32(toAddr)
				}
			}
		}
	case *tokenizationtypes.MsgSetIncomingApproval:
		m.Creator = creatorCosmosAddr
		// Convert addresses in approval criteria if present
		if m.Approval != nil && m.Approval.ApprovalCriteria != nil {
			convertAddressesInIncomingApprovalCriteria(m.Approval.ApprovalCriteria)
		}
	case *tokenizationtypes.MsgSetOutgoingApproval:
		m.Creator = creatorCosmosAddr
		// Convert addresses in approval criteria if present
		if m.Approval != nil && m.Approval.ApprovalCriteria != nil {
			convertAddressesInOutgoingApprovalCriteria(m.Approval.ApprovalCriteria)
		}
	case *tokenizationtypes.MsgCreateCollection:
		m.Creator = creatorCosmosAddr
		// Convert Manager if present
		if m.Manager != "" {
			m.Manager = convertEVMAddressToBech32(m.Manager)
		}
		// Convert addresses in collection approvals
		convertAddressesInCollectionApprovals(m.CollectionApprovals)
	case *tokenizationtypes.MsgUpdateCollection:
		m.Creator = creatorCosmosAddr
		// Convert Manager if present
		if m.Manager != "" {
			m.Manager = convertEVMAddressToBech32(m.Manager)
		}
		// Convert addresses in collection approvals
		convertAddressesInCollectionApprovals(m.CollectionApprovals)
	case *tokenizationtypes.MsgDeleteCollection:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgCreateAddressLists:
		m.Creator = creatorCosmosAddr
		// Convert addresses in address lists
		for _, list := range m.AddressLists {
			if list != nil {
				for i, addr := range list.Addresses {
					list.Addresses[i] = convertEVMAddressToBech32(addr)
				}
			}
		}
	case *tokenizationtypes.MsgUpdateUserApprovals:
		m.Creator = creatorCosmosAddr
		// Convert addresses in incoming approvals
		for _, approval := range m.IncomingApprovals {
			if approval != nil && approval.ApprovalCriteria != nil {
				convertAddressesInIncomingApprovalCriteria(approval.ApprovalCriteria)
			}
		}
		// Convert addresses in outgoing approvals
		for _, approval := range m.OutgoingApprovals {
			if approval != nil && approval.ApprovalCriteria != nil {
				convertAddressesInOutgoingApprovalCriteria(approval.ApprovalCriteria)
			}
		}
	case *tokenizationtypes.MsgDeleteIncomingApproval:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgDeleteOutgoingApproval:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgPurgeApprovals:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgCreateDynamicStore:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgUpdateDynamicStore:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgDeleteDynamicStore:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgSetDynamicStoreValue:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgSetValidTokenIds:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgSetManager:
		m.Creator = creatorCosmosAddr
		// Convert Manager if present
		if m.Manager != "" {
			m.Manager = convertEVMAddressToBech32(m.Manager)
		}
	case *tokenizationtypes.MsgSetCollectionMetadata:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgSetTokenMetadata:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgSetCustomData:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgSetStandards:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgSetCollectionApprovals:
		m.Creator = creatorCosmosAddr
		// Convert addresses in collection approvals
		convertAddressesInCollectionApprovals(m.CollectionApprovals)
	case *tokenizationtypes.MsgSetIsArchived:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgCastVote:
		m.Creator = creatorCosmosAddr
	case *tokenizationtypes.MsgUniversalUpdateCollection:
		m.Creator = creatorCosmosAddr
		// Convert Manager if present
		if m.Manager != "" {
			m.Manager = convertEVMAddressToBech32(m.Manager)
		}
		// Convert addresses in collection approvals
		convertAddressesInCollectionApprovals(m.CollectionApprovals)
	}

	// Validate message using ValidateBasic
	// Use panic recovery to handle cases where ValidateBasic might panic on nil/uninitialized fields
	if validator, ok := msg.(interface{ ValidateBasic() error }); ok {
		var validationErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Convert panic to error - this handles cases where ValidateBasic panics on nil fields
					// This is a safety measure for production readiness
					validationErr = fmt.Errorf("validation panic: %v", r)
				}
			}()
			validationErr = validator.ValidateBasic()
		}()
		if validationErr != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("message validation failed: %v", validationErr))
		}
	}

	return msg, nil
}

// convertAddressesInApprovalCriteria converts EVM addresses to bech32 in ApprovalCriteria
func convertAddressesInApprovalCriteria(criteria *tokenizationtypes.ApprovalCriteria) {
	if criteria == nil {
		return
	}
	// Convert addresses in CoinTransfers
	if criteria.CoinTransfers != nil {
		for _, coinTransfer := range criteria.CoinTransfers {
			if coinTransfer != nil && coinTransfer.To != "" {
				coinTransfer.To = convertEVMAddressToBech32(coinTransfer.To)
			}
		}
	}
	// Convert address in UserRoyalties (single struct, not slice)
	if criteria.UserRoyalties != nil && criteria.UserRoyalties.PayoutAddress != "" {
		criteria.UserRoyalties.PayoutAddress = convertEVMAddressToBech32(criteria.UserRoyalties.PayoutAddress)
	}
}

// convertAddressesInIncomingApprovalCriteria converts EVM addresses to bech32 in IncomingApprovalCriteria
func convertAddressesInIncomingApprovalCriteria(criteria *tokenizationtypes.IncomingApprovalCriteria) {
	if criteria == nil {
		return
	}
	// Convert addresses in CoinTransfers
	if criteria.CoinTransfers != nil {
		for _, coinTransfer := range criteria.CoinTransfers {
			if coinTransfer != nil && coinTransfer.To != "" {
				coinTransfer.To = convertEVMAddressToBech32(coinTransfer.To)
			}
		}
	}
}

// convertAddressesInOutgoingApprovalCriteria converts EVM addresses to bech32 in OutgoingApprovalCriteria
func convertAddressesInOutgoingApprovalCriteria(criteria *tokenizationtypes.OutgoingApprovalCriteria) {
	if criteria == nil {
		return
	}
	// Convert addresses in CoinTransfers
	if criteria.CoinTransfers != nil {
		for _, coinTransfer := range criteria.CoinTransfers {
			if coinTransfer != nil && coinTransfer.To != "" {
				coinTransfer.To = convertEVMAddressToBech32(coinTransfer.To)
			}
		}
	}
}

// convertAddressesInCollectionApprovals converts EVM addresses to bech32 in CollectionApproval slice
func convertAddressesInCollectionApprovals(approvals []*tokenizationtypes.CollectionApproval) {
	if approvals == nil {
		return
	}
	for _, approval := range approvals {
		if approval != nil && approval.ApprovalCriteria != nil {
			convertAddressesInApprovalCriteria(approval.ApprovalCriteria)
		}
	}
}

// unmarshalQueryFromJSON unmarshals a JSON string into the appropriate QueryRequest type
func (p Precompile) unmarshalQueryFromJSON(methodName string, jsonStr string) (interface{}, error) {
	var queryReq interface{}

	switch methodName {
	case GetCollectionMethod:
		queryReq = &tokenizationtypes.QueryGetCollectionRequest{}
	case GetBalanceMethod:
		queryReq = &tokenizationtypes.QueryGetBalanceRequest{}
	case GetAddressListMethod:
		queryReq = &tokenizationtypes.QueryGetAddressListRequest{}
	case GetApprovalTrackerMethod:
		queryReq = &tokenizationtypes.QueryGetApprovalTrackerRequest{}
	case GetChallengeTrackerMethod:
		queryReq = &tokenizationtypes.QueryGetChallengeTrackerRequest{}
	case GetETHSignatureTrackerMethod:
		queryReq = &tokenizationtypes.QueryGetETHSignatureTrackerRequest{}
	case GetDynamicStoreMethod:
		queryReq = &tokenizationtypes.QueryGetDynamicStoreRequest{}
	case GetDynamicStoreValueMethod:
		queryReq = &tokenizationtypes.QueryGetDynamicStoreValueRequest{}
	case GetWrappableBalancesMethod:
		queryReq = &tokenizationtypes.QueryGetWrappableBalancesRequest{}
	case IsAddressReservedProtocolMethod:
		queryReq = &tokenizationtypes.QueryIsAddressReservedProtocolRequest{}
	case GetAllReservedProtocolAddressesMethod:
		queryReq = &tokenizationtypes.QueryGetAllReservedProtocolAddressesRequest{}
	case GetVoteMethod:
		queryReq = &tokenizationtypes.QueryGetVoteRequest{}
	case GetVotesMethod:
		queryReq = &tokenizationtypes.QueryGetVotesRequest{}
	case ParamsMethod:
		queryReq = &tokenizationtypes.QueryParamsRequest{}
	// Note: GetBalanceAmount and GetTotalSupply are handled separately as they don't use standard query requests
	default:
		return nil, ErrInvalidInput(fmt.Sprintf("unknown query method: %s", methodName))
	}

	// Unmarshal JSON into the query request
	if err := json.Unmarshal([]byte(jsonStr), queryReq); err != nil {
		return nil, ErrInvalidInput(fmt.Sprintf("failed to unmarshal query JSON: %v", err))
	}

	// Validate query request using ValidateBasic if available
	// Use panic recovery to handle cases where ValidateBasic might panic on nil/uninitialized fields
	if validator, ok := queryReq.(interface{ ValidateBasic() error }); ok {
		var validationErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Convert panic to error - this handles cases where ValidateBasic panics on nil fields
					validationErr = fmt.Errorf("validation panic: %v", r)
				}
			}()
			validationErr = validator.ValidateBasic()
		}()
		if validationErr != nil {
			return nil, ErrInvalidInput(fmt.Sprintf("query validation failed: %v", validationErr))
		}
	}

	return queryReq, nil
}
