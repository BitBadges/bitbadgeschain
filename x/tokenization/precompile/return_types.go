package tokenization

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// IMPORTANT: Return Type Simplifications
//
// Due to the complexity of nested structures in the tokenization types,
// some fields are returned as empty arrays in the Solidity struct representations.
// This is an intentional simplification for the EVM interface.
//
// Affected Fields (returned as empty arrays):
// - TokenCollection: approvalCriteria (in collectionApprovals), invariants.cosmosCoinBackedPath,
//   userPermissions (in defaultBalances)
// - UserBalanceStore: approvalCriteria (in outgoing/incomingApprovals)
// - CollectionApproval: approvalCriteria
//
// For full access to these fields, clients should:
// 1. Use the raw bytes returned from query methods
// 2. Decode using the appropriate protobuf codec
// 3. Or query the chain directly via gRPC/REST

// ConvertCollectionToSolidityStruct converts a proto TokenCollection to Solidity struct format
// Returns values that can be packed into ABI tuple
func ConvertCollectionToSolidityStruct(collection *tokenizationtypes.TokenCollection) ([]interface{}, error) {
	if collection == nil {
		return nil, fmt.Errorf("collection cannot be nil")
	}

	// Convert collection metadata
	var collectionMetadata []interface{}
	if collection.CollectionMetadata != nil {
		collectionMetadata = []interface{}{
			collection.CollectionMetadata.Uri,
			collection.CollectionMetadata.CustomData,
		}
	} else {
		collectionMetadata = []interface{}{"", ""}
	}

	// Convert token metadata array
	tokenMetadata := make([]interface{}, len(collection.TokenMetadata))
	for i, tm := range collection.TokenMetadata {
		tokenIds := make([]interface{}, len(tm.TokenIds))
		for j, tid := range tm.TokenIds {
			tokenIds[j] = []interface{}{
				tid.Start.BigInt(),
				tid.End.BigInt(),
			}
		}
		tokenMetadata[i] = []interface{}{
			tm.Uri,
			tm.CustomData,
			tokenIds,
		}
	}

	// Convert valid token IDs
	validTokenIds := make([]interface{}, len(collection.ValidTokenIds))
	for i, tid := range collection.ValidTokenIds {
		validTokenIds[i] = []interface{}{
			tid.Start.BigInt(),
			tid.End.BigInt(),
		}
	}

	// Convert collection approvals
	// NOTE: approvalCriteria is simplified to empty array - see package docs for details
	collectionApprovals := make([]interface{}, len(collection.CollectionApprovals))
	for i, app := range collection.CollectionApprovals {
		transferTimes := make([]interface{}, len(app.TransferTimes))
		for j, tt := range app.TransferTimes {
			transferTimes[j] = []interface{}{tt.Start.BigInt(), tt.End.BigInt()}
		}
		tokenIds := make([]interface{}, len(app.TokenIds))
		for j, tid := range app.TokenIds {
			tokenIds[j] = []interface{}{tid.Start.BigInt(), tid.End.BigInt()}
		}
		ownershipTimes := make([]interface{}, len(app.OwnershipTimes))
		for j, ot := range app.OwnershipTimes {
			ownershipTimes[j] = []interface{}{ot.Start.BigInt(), ot.End.BigInt()}
		}
		collectionApprovals[i] = []interface{}{
			app.FromListId,
			app.ToListId,
			app.InitiatedByListId,
			transferTimes,
			tokenIds,
			ownershipTimes,
			app.Uri,
			app.CustomData,
			app.ApprovalId,
			[]interface{}{}, // approvalCriteria - simplified
			app.Version.BigInt(),
		}
	}

	// Convert default balances (simplified)
	var defaultBalances []interface{}
	if collection.DefaultBalances != nil {
		balances := make([]interface{}, len(collection.DefaultBalances.Balances))
		for i, bal := range collection.DefaultBalances.Balances {
			ownershipTimes := make([]interface{}, len(bal.OwnershipTimes))
			for j, ot := range bal.OwnershipTimes {
				ownershipTimes[j] = []interface{}{ot.Start.BigInt(), ot.End.BigInt()}
			}
			tokenIds := make([]interface{}, len(bal.TokenIds))
			for j, tid := range bal.TokenIds {
				tokenIds[j] = []interface{}{tid.Start.BigInt(), tid.End.BigInt()}
			}
			balances[i] = []interface{}{
				bal.Amount.BigInt(),
				ownershipTimes,
				tokenIds,
			}
		}
		defaultBalances = []interface{}{
			balances,
			[]interface{}{}, // outgoingApprovals - simplified
			[]interface{}{}, // incomingApprovals - simplified
			collection.DefaultBalances.AutoApproveSelfInitiatedOutgoingTransfers,
			collection.DefaultBalances.AutoApproveSelfInitiatedIncomingTransfers,
			collection.DefaultBalances.AutoApproveAllIncomingTransfers,
			[]interface{}{}, // userPermissions - simplified
		}
	} else {
		defaultBalances = []interface{}{
			[]interface{}{},
			[]interface{}{},
			[]interface{}{},
			false,
			false,
			false,
			[]interface{}{},
		}
	}

	// Convert collection permissions (simplified)
	var collectionPermissions []interface{}
	if collection.CollectionPermissions != nil {
		collectionPermissions = []interface{}{
			[]interface{}{}, // canDeleteCollection - simplified
			[]interface{}{}, // canArchiveCollection - simplified
			[]interface{}{}, // canUpdateStandards - simplified
			[]interface{}{}, // canUpdateCustomData - simplified
			[]interface{}{}, // canUpdateManager - simplified
			[]interface{}{}, // canUpdateCollectionMetadata - simplified
			[]interface{}{}, // canUpdateValidTokenIds - simplified
			[]interface{}{}, // canUpdateTokenMetadata - simplified
			[]interface{}{}, // canUpdateCollectionApprovals - simplified
			[]interface{}{}, // canAddMoreAliasPaths - simplified
			[]interface{}{}, // canAddMoreCosmosCoinWrapperPaths - simplified
		}
	} else {
		collectionPermissions = []interface{}{
			[]interface{}{},
			[]interface{}{},
			[]interface{}{},
			[]interface{}{},
			[]interface{}{},
			[]interface{}{},
			[]interface{}{},
			[]interface{}{},
			[]interface{}{},
			[]interface{}{},
			[]interface{}{},
		}
	}

	// Convert cosmos coin wrapper paths (simplified)
	cosmosCoinWrapperPaths := make([]interface{}, len(collection.CosmosCoinWrapperPaths))
	for i, path := range collection.CosmosCoinWrapperPaths {
		denomUnits := make([]interface{}, len(path.DenomUnits))
		for j, du := range path.DenomUnits {
			var metadata []interface{}
			if du.Metadata != nil {
				metadata = []interface{}{du.Metadata.Uri, du.Metadata.CustomData}
			} else {
				metadata = []interface{}{"", ""}
			}
			denomUnits[j] = []interface{}{
				du.Decimals.BigInt(),
				du.Symbol,
				du.IsDefaultDisplay,
				metadata,
			}
		}
		var pathMetadata []interface{}
		if path.Metadata != nil {
			pathMetadata = []interface{}{path.Metadata.Uri, path.Metadata.CustomData}
		} else {
			pathMetadata = []interface{}{"", ""}
		}
		cosmosCoinWrapperPaths[i] = []interface{}{
			path.Denom,
			[]interface{}{}, // conversion - simplified
			path.Symbol,
			denomUnits,
			path.AllowOverrideWithAnyValidToken,
			pathMetadata,
		}
	}

	// Convert alias paths (simplified)
	aliasPaths := make([]interface{}, len(collection.AliasPaths))
	for i, path := range collection.AliasPaths {
		denomUnits := make([]interface{}, len(path.DenomUnits))
		for j, du := range path.DenomUnits {
			var metadata []interface{}
			if du.Metadata != nil {
				metadata = []interface{}{du.Metadata.Uri, du.Metadata.CustomData}
			} else {
				metadata = []interface{}{"", ""}
			}
			denomUnits[j] = []interface{}{
				du.Decimals.BigInt(),
				du.Symbol,
				du.IsDefaultDisplay,
				metadata,
			}
		}
		var pathMetadata []interface{}
		if path.Metadata != nil {
			pathMetadata = []interface{}{path.Metadata.Uri, path.Metadata.CustomData}
		} else {
			pathMetadata = []interface{}{"", ""}
		}
		aliasPaths[i] = []interface{}{
			path.Denom,
			[]interface{}{}, // conversion - simplified
			path.Symbol,
			denomUnits,
			pathMetadata,
		}
	}

	// Convert invariants (simplified)
	var invariants []interface{}
	if collection.Invariants != nil {
		invariants = []interface{}{
			collection.Invariants.NoCustomOwnershipTimes,
			collection.Invariants.MaxSupplyPerId.BigInt(),
			[]interface{}{}, // cosmosCoinBackedPath - simplified
			collection.Invariants.NoForcefulPostMintTransfers,
			collection.Invariants.DisablePoolCreation,
		}
	} else {
		invariants = []interface{}{
			false,
			big.NewInt(0),
			[]interface{}{},
			false,
			false,
		}
	}

	// Pack into tuple structure matching Solidity struct
	return []interface{}{
		collection.CollectionId.BigInt(),
		collectionMetadata,
		tokenMetadata,
		collection.CustomData,
		collection.Manager,
		collectionPermissions,
		collectionApprovals,
		collection.Standards,
		collection.IsArchived,
		defaultBalances,
		collection.CreatedBy,
		validTokenIds,
		collection.MintEscrowAddress,
		cosmosCoinWrapperPaths,
		invariants,
		aliasPaths,
	}, nil
}

// ConvertBalanceToSolidityStruct converts a proto Balance to Solidity struct format
func ConvertBalanceToSolidityStruct(balance *tokenizationtypes.Balance) ([]interface{}, error) {
	if balance == nil {
		return nil, fmt.Errorf("balance cannot be nil")
	}

	ownershipTimes := make([]interface{}, len(balance.OwnershipTimes))
	for i, ot := range balance.OwnershipTimes {
		ownershipTimes[i] = []interface{}{
			ot.Start.BigInt(),
			ot.End.BigInt(),
		}
	}

	tokenIds := make([]interface{}, len(balance.TokenIds))
	for i, tid := range balance.TokenIds {
		tokenIds[i] = []interface{}{
			tid.Start.BigInt(),
			tid.End.BigInt(),
		}
	}

	return []interface{}{
		balance.Amount.BigInt(),
		ownershipTimes,
		tokenIds,
	}, nil
}

// ConvertUserBalanceStoreToSolidityStruct converts a proto UserBalanceStore to Solidity struct format
func ConvertUserBalanceStoreToSolidityStruct(store *tokenizationtypes.UserBalanceStore) ([]interface{}, error) {
	if store == nil {
		return nil, fmt.Errorf("store cannot be nil")
	}

	balances := make([]interface{}, len(store.Balances))
	for i, bal := range store.Balances {
		converted, err := ConvertBalanceToSolidityStruct(bal)
		if err != nil {
			return nil, fmt.Errorf("balance[%d]: %w", i, err)
		}
		balances[i] = converted
	}

	// Convert outgoing approvals (simplified)
	outgoingApprovals := make([]interface{}, len(store.OutgoingApprovals))
	for i, app := range store.OutgoingApprovals {
		transferTimes := make([]interface{}, len(app.TransferTimes))
		for j, tt := range app.TransferTimes {
			transferTimes[j] = []interface{}{tt.Start.BigInt(), tt.End.BigInt()}
		}
		tokenIds := make([]interface{}, len(app.TokenIds))
		for j, tid := range app.TokenIds {
			tokenIds[j] = []interface{}{tid.Start.BigInt(), tid.End.BigInt()}
		}
		ownershipTimes := make([]interface{}, len(app.OwnershipTimes))
		for j, ot := range app.OwnershipTimes {
			ownershipTimes[j] = []interface{}{ot.Start.BigInt(), ot.End.BigInt()}
		}
		outgoingApprovals[i] = []interface{}{
			app.ApprovalId,
			app.ToListId,
			app.InitiatedByListId,
			transferTimes,
			tokenIds,
			ownershipTimes,
			app.Uri,
			app.CustomData,
		}
	}

	// Convert incoming approvals (simplified)
	incomingApprovals := make([]interface{}, len(store.IncomingApprovals))
	for i, app := range store.IncomingApprovals {
		transferTimes := make([]interface{}, len(app.TransferTimes))
		for j, tt := range app.TransferTimes {
			transferTimes[j] = []interface{}{tt.Start.BigInt(), tt.End.BigInt()}
		}
		tokenIds := make([]interface{}, len(app.TokenIds))
		for j, tid := range app.TokenIds {
			tokenIds[j] = []interface{}{tid.Start.BigInt(), tid.End.BigInt()}
		}
		ownershipTimes := make([]interface{}, len(app.OwnershipTimes))
		for j, ot := range app.OwnershipTimes {
			ownershipTimes[j] = []interface{}{ot.Start.BigInt(), ot.End.BigInt()}
		}
		incomingApprovals[i] = []interface{}{
			app.ApprovalId,
			app.FromListId,
			app.InitiatedByListId,
			transferTimes,
			tokenIds,
			ownershipTimes,
			app.Uri,
			app.CustomData,
		}
	}

	// Convert user permissions (simplified)
	var userPermissions []interface{}
	if store.UserPermissions != nil {
		userPermissions = []interface{}{
			[]interface{}{}, // canUpdateOutgoingApprovals - simplified
			[]interface{}{}, // canUpdateIncomingApprovals - simplified
			[]interface{}{}, // canUpdateAutoApproveSelfInitiatedOutgoingTransfers - simplified
			[]interface{}{}, // canUpdateAutoApproveSelfInitiatedIncomingTransfers - simplified
			[]interface{}{}, // canUpdateAutoApproveAllIncomingTransfers - simplified
		}
	} else {
		userPermissions = []interface{}{
			[]interface{}{}, []interface{}{}, []interface{}{}, []interface{}{}, []interface{}{},
		}
	}

	return []interface{}{
		balances,
		outgoingApprovals,
		incomingApprovals,
		store.AutoApproveSelfInitiatedOutgoingTransfers,
		store.AutoApproveSelfInitiatedIncomingTransfers,
		store.AutoApproveAllIncomingTransfers,
		userPermissions,
	}, nil
}

// PackCollectionAsStruct packs a collection response as a Solidity struct tuple
// If ABI expects bytes, marshal to bytes. Otherwise pack as struct tuple.
func PackCollectionAsStruct(method *abi.Method, collection *tokenizationtypes.TokenCollection) ([]byte, error) {
	// Check if ABI expects bytes (legacy) or struct tuple
	if len(method.Outputs) == 1 && method.Outputs[0].Type.T == abi.BytesTy {
		// Legacy: marshal to bytes
		bz, err := tokenizationtypes.ModuleCdc.Marshal(collection)
		if err != nil {
			return nil, ErrInternalError(fmt.Sprintf("marshal collection failed: %s", err))
		}
		return bz, nil
	}
	// New: pack as struct tuple
	structData, err := ConvertCollectionToSolidityStruct(collection)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("convert collection failed: %s", err))
	}
	packed, err := method.Outputs.Pack(structData...)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("pack collection failed: %s", err))
	}
	return packed, nil
}

// PackBalanceAsStruct packs a balance response as a Solidity struct tuple
func PackBalanceAsStruct(method *abi.Method, balance *tokenizationtypes.Balance) ([]byte, error) {
	structData, err := ConvertBalanceToSolidityStruct(balance)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("convert balance failed: %s", err))
	}
	packed, err := method.Outputs.Pack(structData...)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("pack balance failed: %s", err))
	}
	return packed, nil
}

// PackUserBalanceStoreAsStruct packs a UserBalanceStore response as a Solidity struct tuple
// If ABI expects bytes, marshal to bytes. Otherwise pack as struct tuple.
func PackUserBalanceStoreAsStruct(method *abi.Method, store *tokenizationtypes.UserBalanceStore) ([]byte, error) {
	// Check if ABI expects bytes (legacy) or struct tuple
	if len(method.Outputs) == 1 && method.Outputs[0].Type.T == abi.BytesTy {
		// Legacy: marshal to bytes
		bz, err := tokenizationtypes.ModuleCdc.Marshal(store)
		if err != nil {
			return nil, ErrInternalError(fmt.Sprintf("marshal user balance store failed: %s", err))
		}
		return bz, nil
	}
	// New: pack as struct tuple
	structData, err := ConvertUserBalanceStoreToSolidityStruct(store)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("convert user balance store failed: %s", err))
	}
	packed, err := method.Outputs.Pack(structData...)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("pack user balance store failed: %s", err))
	}
	return packed, nil
}

// ConvertAddressListToSolidityStruct converts a proto AddressList to Solidity struct format
func ConvertAddressListToSolidityStruct(list *tokenizationtypes.AddressList) ([]interface{}, error) {
	if list == nil {
		return nil, ErrInvalidInput("list cannot be nil")
	}

	// Convert addresses (EVM addresses need to be converted back from Cosmos format)
	addresses := make([]interface{}, len(list.Addresses))
	for i, addr := range list.Addresses {
		// Try to convert from Cosmos address format to EVM address
		cosmosAddr, err := sdk.AccAddressFromBech32(addr)
		if err == nil {
			// Convert to EVM address format
			evmAddr := common.BytesToAddress(cosmosAddr.Bytes())
			addresses[i] = evmAddr
		} else {
			// Keep as string if not a valid Cosmos address
			addresses[i] = addr
		}
	}

	return []interface{}{
		list.ListId,
		addresses,
		list.Whitelist,
		list.Uri,
		list.CustomData,
		list.CreatedBy,
	}, nil
}

// PackAddressListAsStruct packs an AddressList response as a Solidity struct tuple
// If ABI expects bytes, marshal to bytes. Otherwise pack as struct tuple.
func PackAddressListAsStruct(method *abi.Method, list *tokenizationtypes.AddressList) ([]byte, error) {
	// Check if ABI expects bytes (legacy) or struct tuple
	if len(method.Outputs) == 1 && method.Outputs[0].Type.T == abi.BytesTy {
		// Legacy: marshal to bytes
		bz, err := tokenizationtypes.ModuleCdc.Marshal(list)
		if err != nil {
			return nil, ErrInternalError(fmt.Sprintf("marshal address list failed: %s", err))
		}
		return bz, nil
	}
	// New: pack as struct tuple
	structData, err := ConvertAddressListToSolidityStruct(list)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("convert address list failed: %s", err))
	}
	packed, err := method.Outputs.Pack(structData...)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("pack address list failed: %s", err))
	}
	return packed, nil
}
