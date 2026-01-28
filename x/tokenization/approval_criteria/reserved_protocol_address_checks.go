package approval_criteria

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ReservedProtocolAddressChecker implements ApprovalCriteriaChecker for checking reserved protocol addresses
// This prevents forceful transfers from reserved protocol addresses when OverridesFromOutgoingApprovals is true
type ReservedProtocolAddressChecker struct {
	addressCheckService AddressCheckService
}

// NewReservedProtocolAddressChecker creates a new ReservedProtocolAddressChecker
func NewReservedProtocolAddressChecker(addressCheckService AddressCheckService) *ReservedProtocolAddressChecker {
	return &ReservedProtocolAddressChecker{
		addressCheckService: addressCheckService,
	}
}

// Name returns the name of this checker
func (c *ReservedProtocolAddressChecker) Name() string {
	return "ReservedProtocolAddress"
}

// Check validates that forceful transfers from reserved protocol addresses are disallowed
// This is an important check to prevent abuse of systems built on top of our standard
// Ex: For liquidity pools, we don't want to allow forceful revocations from manager or any addresses
//
//	because this would mess up the entire escrow system and could cause infinite liquidity glitches
//
// Bypass this check if the address is actually initiating it
func (c *ReservedProtocolAddressChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	approvalCriteria := approval.ApprovalCriteria
	if approvalCriteria == nil {
		return "", nil
	}

	// Only check if OverridesFromOutgoingApprovals is true and fromAddress is not the initiator
	if !approvalCriteria.OverridesFromOutgoingApprovals || from == initiator {
		return "", nil
	}

	// Check if the from address is a reserved protocol address
	if c.addressCheckService.IsAddressReservedProtocol(ctx, from) {
		detErrMsg := fmt.Sprintf("forceful transfers that bypass user approvals from reserved protocol addresses are globally disallowed (please use an approval that checks user-level approvals): %s", from)
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	return "", nil
}
