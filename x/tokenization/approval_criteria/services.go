package approval_criteria

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CollectionService provides methods to access collection data
type CollectionService interface {
	GetCollection(ctx sdk.Context, collectionId sdkmath.Uint) (*types.TokenCollection, bool)
	GetBalanceOrApplyDefault(ctx sdk.Context, collection *types.TokenCollection, userAddress string) (*types.UserBalanceStore, bool, error)
}

// AddressCheckService provides methods to check address types
type AddressCheckService interface {
	IsEVMContract(ctx sdk.Context, address string) (bool, error)
	IsLiquidityPool(ctx sdk.Context, address string) (bool, error)
	IsAddressReservedProtocol(ctx sdk.Context, address string) bool
}

// DynamicStoreService provides methods to access dynamic store data
type DynamicStoreService interface {
	GetDynamicStore(ctx sdk.Context, storeId sdkmath.Uint) (*types.DynamicStore, bool)
	GetDynamicStoreValue(ctx sdk.Context, storeId sdkmath.Uint, address string) (*types.DynamicStoreValue, bool)
}

// AddressListService provides methods to check addresses against address lists
type AddressListService interface {
	CheckAddresses(ctx sdk.Context, addressListId string, addressToCheck string) (bool, error)
}

// EVMQueryService provides methods to execute read-only EVM queries
type EVMQueryService interface {
	// ExecuteEVMQuery performs a read-only call to an EVM contract.
	// callerAddress is the address to use as the EVM caller (msg.sender in the contract); if empty, the implementation may use the zero address.
	// Returns the raw return data or an error.
	ExecuteEVMQuery(ctx sdk.Context, callerAddress string, contractAddress string, calldata []byte, gasLimit uint64) ([]byte, error)
}
