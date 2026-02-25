package approval_criteria

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultEVMQueryGasLimit = uint64(100000)
	MaxEVMQueryGasLimit     = uint64(500000)
	// MaxTotalEVMQueryGas limits total gas across all challenges to prevent DoS
	// Even with 10 challenges, total gas is capped at 1M to prevent excessive computation
	MaxTotalEVMQueryGas = uint64(1000000)
)

// EVMQueryChallengesChecker implements ApprovalCriteriaChecker for EVMQueryChallenges
type EVMQueryChallengesChecker struct {
	evmQueryService EVMQueryService
}

// NewEVMQueryChallengesChecker creates a new EVMQueryChallengesChecker
func NewEVMQueryChallengesChecker(evmQueryService EVMQueryService) *EVMQueryChallengesChecker {
	return &EVMQueryChallengesChecker{
		evmQueryService: evmQueryService,
	}
}

// Name returns the name of this checker
func (c *EVMQueryChallengesChecker) Name() string {
	return "EVMQueryChallenges"
}

// Check validates EVM query challenges by executing read-only contract calls
func (c *EVMQueryChallengesChecker) Check(ctx sdk.Context, approval *types.CollectionApproval, collection *types.TokenCollection, to string, from string, initiator string, approvalLevel string, approverAddress string, merkleProofs []*types.MerkleProof, ethSignatureProofs []*types.ETHSignatureProof, memo string, isPrioritized bool) (string, error) {
	challenges := c.getChallenges(approval)
	if len(challenges) == 0 {
		return "", nil
	}

	// Track total gas used across all challenges to prevent DoS
	var totalGasUsed uint64

	for i, challenge := range challenges {
		if challenge == nil {
			continue // Skip nil challenges
		}

		// Replace placeholders in calldata
		calldata, err := c.replacePlaceholders(challenge.Calldata, initiator, from, to, collection.CollectionId.String())
		if err != nil {
			detErrMsg := fmt.Sprintf("EVM query challenge %d placeholder replacement: %s", i, err.Error())
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Decode hex calldata
		calldataBytes, err := hex.DecodeString(strings.TrimPrefix(calldata, "0x"))
		if err != nil {
			detErrMsg := fmt.Sprintf("invalid calldata hex in EVM query challenge %d: %s", i, err.Error())
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Determine gas limit
		gasLimit := DefaultEVMQueryGasLimit
		if !challenge.GasLimit.IsZero() {
			gasLimit = challenge.GasLimit.Uint64()
			if gasLimit > MaxEVMQueryGasLimit {
				gasLimit = MaxEVMQueryGasLimit
			}
		}

		// Check total gas limit across all challenges (DoS protection)
		if totalGasUsed+gasLimit > MaxTotalEVMQueryGas {
			detErrMsg := fmt.Sprintf("EVM query challenges exceed total gas limit (%d)", MaxTotalEVMQueryGas)
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
		totalGasUsed += gasLimit

		// Execute the query with zero address as caller (from/to/initiator passed via calldata placeholders)
		result, err := c.evmQueryService.ExecuteEVMQuery(ctx, "", challenge.ContractAddress, calldataBytes, gasLimit)
		if err != nil {
			detErrMsg := fmt.Sprintf("EVM query challenge %d failed: %s", i, err.Error())
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Compare result if expected result is provided
		if challenge.ExpectedResult != "" {
			expectedBytes, err := hex.DecodeString(strings.TrimPrefix(challenge.ExpectedResult, "0x"))
			if err != nil {
				detErrMsg := fmt.Sprintf("invalid expected result hex in EVM query challenge %d", i)
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}

			if !c.compareResults(result, expectedBytes, challenge.ComparisonOperator) {
				detErrMsg := fmt.Sprintf("EVM query challenge %d: result mismatch (got %x, expected %x with op %s)",
					i, result, expectedBytes, challenge.ComparisonOperator)
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
		}
	}

	return "", nil
}

// getChallenges retrieves EVM query challenges from the approval criteria
func (c *EVMQueryChallengesChecker) getChallenges(approval *types.CollectionApproval) []*types.EVMQueryChallenge {
	if approval.ApprovalCriteria != nil {
		return approval.ApprovalCriteria.EvmQueryChallenges
	}
	return nil
}

// replacePlaceholders replaces address placeholders in calldata with actual addresses
func (c *EVMQueryChallengesChecker) replacePlaceholders(calldata, initiator, from, to, collectionId string) (string, error) {
	initiatorHex, err := addressToHexPadded(initiator)
	if err != nil {
		return "", sdkerrors.Wrapf(err, "initiator address")
	}
	fromHex, err := addressToHexPadded(from)
	if err != nil {
		return "", sdkerrors.Wrapf(err, "sender address")
	}
	toHex, err := addressToHexPadded(to)
	if err != nil {
		return "", sdkerrors.Wrapf(err, "recipient address")
	}
	calldata = strings.ReplaceAll(calldata, "$initiator", initiatorHex)
	calldata = strings.ReplaceAll(calldata, "$sender", fromHex)
	calldata = strings.ReplaceAll(calldata, "$recipient", toHex)
	calldata = strings.ReplaceAll(calldata, "$collectionId", uint256ToHexPadded(collectionId))
	return calldata, nil
}

// uint256ToHexPadded converts a uint string to 32-byte padded hex
func uint256ToHexPadded(value string) string {
	n := new(big.Int)
	n.SetString(value, 10)
	return fmt.Sprintf("%064x", n)
}

// compareResults compares the actual result with the expected result using the specified operator
func (c *EVMQueryChallengesChecker) compareResults(result, expected []byte, operator string) bool {
	switch operator {
	case "", "eq":
		return bytes.Equal(result, expected)
	case "ne":
		return !bytes.Equal(result, expected)
	case "gt", "gte", "lt", "lte":
		// For numeric comparisons, interpret as big.Int
		resultInt := new(big.Int).SetBytes(result)
		expectedInt := new(big.Int).SetBytes(expected)
		cmp := resultInt.Cmp(expectedInt)
		switch operator {
		case "gt":
			return cmp > 0
		case "gte":
			return cmp >= 0
		case "lt":
			return cmp < 0
		case "lte":
			return cmp <= 0
		}
	}
	return false
}

// addressToHexPadded converts an address to a hex string left-padded to 32 bytes for ABI encoding.
// Accepts valid bech32, 0x + 40 hex chars, or protocol addresses (Mint, Total) which get a canonical encoding.
func addressToHexPadded(address string) (string, error) {
	// Protocol addresses: use canonical 32-byte hex so EVM calldata is valid
	if types.IsMintAddress(address) {
		return strings.Repeat("0", 64), nil
	}
	if types.IsTotalAddress(address) {
		return strings.Repeat("0", 63) + "1", nil
	}
	accAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		hexAddr := strings.TrimSpace(strings.TrimPrefix(address, "0x"))
		if len(hexAddr) == 40 {
			for _, c := range hexAddr {
				if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') {
					continue
				}
				return "", sdkerrors.Wrapf(types.ErrInvalidRequest, "invalid address for placeholder: not bech32 and not 40 hex chars")
			}
			return strings.Repeat("0", 24) + hexAddr, nil
		}
		return "", sdkerrors.Wrapf(types.ErrInvalidRequest, "invalid address for placeholder: must be bech32 or 0x + 40 hex chars")
	}
	return fmt.Sprintf("%064x", accAddr.Bytes()), nil
}
