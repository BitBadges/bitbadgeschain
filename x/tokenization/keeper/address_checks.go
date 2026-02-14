package keeper

import (
	"github.com/ethereum/go-ethereum/common"

	approvalcriteria "github.com/bitbadges/bitbadgeschain/x/tokenization/approval_criteria"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CheckAddressChecks validates address checks for a given address
// Returns (deterministicErrorMsg, error) where deterministicErrorMsg is a deterministic error string
// This method is kept for backward compatibility with tests
func (k Keeper) CheckAddressChecks(ctx sdk.Context, addressChecks *types.AddressChecks, address string) (string, error) {
	if addressChecks == nil {
		return "", nil
	}

	// Create a dummy approval with the address checks for the checker
	approval := &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			SenderChecks: addressChecks, // We'll use sender checks as default
		},
	}

	// Use the new dynamic checker approach
	checkers := k.GetApprovalCriteriaCheckers(approval)
	for _, checker := range checkers {
		// Find the address checker (should be the sender one we just created)
		if _, ok := checker.(*approvalcriteria.AddressChecksChecker); ok {
			// Pass address as the "from" parameter since we're using "sender" check type
			// Pass nil for collection as this is a backward compatibility wrapper
			detErrMsg, err := checker.Check(ctx, approval, nil, "", address, "", "", "", nil, nil, "", false)
			if err != nil {
				return detErrMsg, err
			}
			return detErrMsg, nil
		}
	}
	return "", nil
}

// IsEVMContract checks if an address is an EVM contract (has code)
func (k Keeper) IsEVMContract(ctx sdk.Context, address string) (bool, error) {
	if k.evmKeeper == nil {
		// If EVMKeeper is not set, we can't check - return false
		// This allows the feature to work even if EVM module is not available
		return false, nil
	}

	// Convert bech32 address to Ethereum address
	// First try to parse as bech32 address
	accAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return false, err
	}

	// Convert Cosmos address to Ethereum address (20 bytes)
	ethAddr := common.BytesToAddress(accAddr.Bytes())

	// Check if the address has code (is a contract)
	return k.evmKeeper.IsContract(ctx, ethAddr), nil
}

// IsLiquidityPool checks if an address is a liquidity pool
func (k Keeper) IsLiquidityPool(ctx sdk.Context, address string) (bool, error) {
	if k.gammKeeper == nil {
		// If GammKeeper is not set, we can't check - return false
		// This allows the feature to work even if gamm module is not available
		return false, nil
	}

	// Validate address format
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return false, err
	}

	// The cache is populated when pools are created (see x/gamm/keeper/pool_service.go)
	_, found := k.GetPoolIdFromAddressCache(ctx, address)
	if found {
		return true, nil
	}

	return false, nil
}
