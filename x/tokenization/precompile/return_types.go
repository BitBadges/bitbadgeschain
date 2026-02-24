package tokenization

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// convertEVMQueryChallengesToSolidity converts EvmQueryChallenges to the format expected by the Solidity CollectionInvariants struct.
// Each challenge is output as a 7-tuple: contractAddress, calldata, expectedResult, comparisonOperator, gasLimit, uri, customData.
// uri and customData are always included (empty string when not set) so invariants and approval criteria have full metadata.
func convertEVMQueryChallengesToSolidity(challenges []*tokenizationtypes.EVMQueryChallenge) []interface{} {
	if len(challenges) == 0 {
		return make([]interface{}, 0)
	}
	out := make([]interface{}, 0, len(challenges))
	for _, c := range challenges {
		if c == nil {
			continue
		}
		gasLimit := big.NewInt(0)
		if !c.GasLimit.IsNil() {
			gasLimit = c.GasLimit.BigInt()
		}
		out = append(out, []interface{}{
			c.ContractAddress,
			c.Calldata,
			c.ExpectedResult,
			c.ComparisonOperator,
			gasLimit,
			c.Uri,
			c.CustomData,
		})
	}
	return out
}

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

	// Convert collection approvals (full approvalCriteria)
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
		version := big.NewInt(0)
		if !app.Version.IsNil() {
			version = app.Version.BigInt()
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
			approvalCriteriaToSolidity(app.ApprovalCriteria),
			version,
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
			userOutgoingApprovalsToSolidity(collection.DefaultBalances.OutgoingApprovals),
			userIncomingApprovalsToSolidity(collection.DefaultBalances.IncomingApprovals),
			collection.DefaultBalances.AutoApproveSelfInitiatedOutgoingTransfers,
			collection.DefaultBalances.AutoApproveSelfInitiatedIncomingTransfers,
			collection.DefaultBalances.AutoApproveAllIncomingTransfers,
			userPermissionsToSolidity(collection.DefaultBalances.UserPermissions),
		}
	} else {
		emptySlice := make([]interface{}, 0)
		defaultBalances = []interface{}{
			emptySlice, // balances
			emptySlice, // outgoingApprovals
			emptySlice, // incomingApprovals
			false,
			false,
			false,
			userPermissionsToSolidity(nil),
		}
	}

	// Convert collection permissions
	collectionPermissions := collectionPermissionsToSolidity(collection.CollectionPermissions)

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
			conversionWithoutDenomToSolidity(path.Conversion),
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
			conversionWithoutDenomToSolidity(path.Conversion),
			path.Symbol,
			denomUnits,
			pathMetadata,
		}
	}

	// Convert invariants
	var invariants []interface{}
	if collection.Invariants != nil {
		// Get maxSupplyPerId from CollectionInvariants
		maxSupplyPerId := big.NewInt(0)
		if !collection.Invariants.MaxSupplyPerId.IsNil() {
			maxSupplyPerId = collection.Invariants.MaxSupplyPerId.BigInt()
		}
		evmQueryChallenges := convertEVMQueryChallengesToSolidity(collection.Invariants.EvmQueryChallenges)
		invariants = []interface{}{
			collection.Invariants.NoCustomOwnershipTimes,
			maxSupplyPerId,
			cosmosCoinBackedPathToSolidity(collection.Invariants.CosmosCoinBackedPath),
			collection.Invariants.NoForcefulPostMintTransfers,
			collection.Invariants.DisablePoolCreation,
			evmQueryChallenges,
		}
	} else {
		invariants = []interface{}{
			false,
			big.NewInt(0),
			cosmosCoinBackedPathToSolidity(nil),
			false,
			false,
			convertEVMQueryChallengesToSolidity(nil),
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

	// Convert outgoing approvals (full 10-element tuple with approvalCriteria and version)
	outgoingApprovals := userOutgoingApprovalsToSolidity(store.OutgoingApprovals)

	// Convert incoming approvals (full 10-element tuple with approvalCriteria and version)
	incomingApprovals := userIncomingApprovalsToSolidity(store.IncomingApprovals)

	// Convert user permissions
	userPermissions := userPermissionsToSolidity(store.UserPermissions)

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

// PackCollectionAsStruct packs a collection response as a Solidity struct tuple.
func PackCollectionAsStruct(method *abi.Method, collection *tokenizationtypes.TokenCollection) ([]byte, error) {
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

// PackUserBalanceStoreAsStruct packs a UserBalanceStore response as a Solidity struct tuple.
func PackUserBalanceStoreAsStruct(method *abi.Method, store *tokenizationtypes.UserBalanceStore) ([]byte, error) {
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

// PackApprovalTrackerAsStruct packs an ApprovalTracker response as a Solidity struct tuple.
func PackApprovalTrackerAsStruct(method *abi.Method, tracker *tokenizationtypes.ApprovalTracker) ([]byte, error) {
	structData, err := ConvertApprovalTrackerToSolidityStruct(tracker)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("convert approval tracker failed: %s", err))
	}
	packed, err := method.Outputs.Pack(structData...)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("pack approval tracker failed: %s", err))
	}
	return packed, nil
}

// PackDynamicStoreAsStruct packs a DynamicStore response as a Solidity struct tuple.
func PackDynamicStoreAsStruct(method *abi.Method, store *tokenizationtypes.DynamicStore) ([]byte, error) {
	structData, err := ConvertDynamicStoreToSolidityStruct(store)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("convert dynamic store failed: %s", err))
	}
	packed, err := method.Outputs.Pack(structData...)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("pack dynamic store failed: %s", err))
	}
	return packed, nil
}

// PackDynamicStoreValueAsStruct packs a DynamicStoreValue response as a Solidity struct tuple.
func PackDynamicStoreValueAsStruct(method *abi.Method, value *tokenizationtypes.DynamicStoreValue) ([]byte, error) {
	structData, err := ConvertDynamicStoreValueToSolidityStruct(value)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("convert dynamic store value failed: %s", err))
	}
	packed, err := method.Outputs.Pack(structData...)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("pack dynamic store value failed: %s", err))
	}
	return packed, nil
}

// PackVoteProofAsStruct packs a VoteProof response as a Solidity struct tuple.
func PackVoteProofAsStruct(method *abi.Method, proof *tokenizationtypes.VoteProof) ([]byte, error) {
	structData, err := ConvertVoteProofToSolidityStruct(proof)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("convert vote proof failed: %s", err))
	}
	packed, err := method.Outputs.Pack(structData...)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("pack vote proof failed: %s", err))
	}
	return packed, nil
}

// PackParamsAsStruct packs Params response as a Solidity struct tuple.
func PackParamsAsStruct(method *abi.Method, params *tokenizationtypes.Params) ([]byte, error) {
	structData, err := ConvertParamsToSolidityStruct(params)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("convert params failed: %s", err))
	}
	packed, err := method.Outputs.Pack(structData...)
	if err != nil {
		return nil, ErrInternalError(fmt.Sprintf("pack params failed: %s", err))
	}
	return packed, nil
}

// PackAddressListAsStruct packs an AddressList response as a Solidity struct tuple.
func PackAddressListAsStruct(method *abi.Method, list *tokenizationtypes.AddressList) ([]byte, error) {
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
