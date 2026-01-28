package approval_criteria

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddressChecksChecker implements ApprovalCriteriaChecker for AddressChecks
// This checker can be used for sender, recipient, or initiator checks
type AddressChecksChecker struct {
	addressCheckService AddressCheckService
	addressChecks       *types.AddressChecks
	checkType           string // "sender", "recipient", or "initiator"
}

// NewAddressChecksChecker creates a new AddressChecksChecker for a specific check type
func NewAddressChecksChecker(addressCheckService AddressCheckService, addressChecks *types.AddressChecks, checkType string) *AddressChecksChecker {
	return &AddressChecksChecker{
		addressCheckService: addressCheckService,
		addressChecks:       addressChecks,
		checkType:           checkType,
	}
}

// Name returns the name of this checker
func (c *AddressChecksChecker) Name() string {
	return fmt.Sprintf("AddressChecks-%s", c.checkType)
}

// Check validates address checks for the appropriate address based on checkType
func (c *AddressChecksChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	if c.addressChecks == nil {
		return "", nil
	}

	// Determine which address to check based on checkType
	var addressToCheck string
	switch c.checkType {
	case "sender":
		addressToCheck = from
	case "recipient":
		addressToCheck = to
	case "initiator":
		addressToCheck = initiator
	default:
		return "", errorsmod.Wrapf(types.ErrInvalidCheckType, "checkType: %s", c.checkType)
	}

	// Check WASM contract requirements
	if c.addressChecks.MustBeWasmContract {
		detErrMsg := fmt.Sprintf("address %s must be a WASM contract", addressToCheck)
		isWasm, err := c.addressCheckService.IsWasmContract(ctx, addressToCheck)
		if err != nil {
			return detErrMsg, err
		}
		if !isWasm {
			return detErrMsg, errorsmod.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	if c.addressChecks.MustNotBeWasmContract {
		detErrMsg := fmt.Sprintf("address %s must not be a WASM contract", addressToCheck)
		isWasm, err := c.addressCheckService.IsWasmContract(ctx, addressToCheck)
		if err != nil {
			return detErrMsg, err
		}
		if isWasm {
			return detErrMsg, errorsmod.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	// Check liquidity pool requirements
	if c.addressChecks.MustBeLiquidityPool {
		detErrMsg := fmt.Sprintf("address %s must be a liquidity pool", addressToCheck)
		isPool, err := c.addressCheckService.IsLiquidityPool(ctx, addressToCheck)
		if err != nil {
			return detErrMsg, err
		}
		if !isPool {
			return detErrMsg, errorsmod.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	if c.addressChecks.MustNotBeLiquidityPool {
		detErrMsg := fmt.Sprintf("address %s must not be a liquidity pool", addressToCheck)
		isPool, err := c.addressCheckService.IsLiquidityPool(ctx, addressToCheck)
		if err != nil {
			return detErrMsg, err
		}
		if isPool {
			return detErrMsg, errorsmod.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
	}

	return "", nil
}
