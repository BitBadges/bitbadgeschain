package keeper

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
	// MaxTotalPostTransferInvariantGas limits total gas across all post-transfer invariants
	MaxTotalPostTransferInvariantGas = uint64(1000000)
	// DefaultPostTransferEVMQueryGasLimit is the default gas limit for post-transfer EVM queries
	DefaultPostTransferEVMQueryGasLimit = uint64(100000)
	// MaxPostTransferEVMQueryGasLimit is the maximum gas limit for a single post-transfer EVM query
	MaxPostTransferEVMQueryGasLimit = uint64(500000)
)

// CheckPostTransferInvariants runs all post-transfer invariants ONCE after all transfers complete
// Called after ALL SetBalanceForAddress() calls for all recipients in the message
func (k Keeper) CheckPostTransferInvariants(
	ctx sdk.Context,
	collection *types.TokenCollection,
	from string,
	toAddresses []string,
	initiatedBy string,
) error {
	if collection.Invariants == nil {
		return nil
	}

	invariants := collection.Invariants

	// Check maxSupplyPerId if set
	if !invariants.MaxSupplyPerId.IsNil() && !invariants.MaxSupplyPerId.IsZero() {
		if err := k.checkMaxSupplyPerIdInvariant(ctx, collection, invariants.MaxSupplyPerId); err != nil {
			return err
		}
	}

	// Check EVM query challenges
	if err := k.checkPostTransferEVMQueryInvariants(ctx, collection, from, toAddresses, initiatedBy, invariants.EvmQueryChallenges); err != nil {
		return err
	}

	return nil
}

// checkMaxSupplyPerIdInvariant verifies no token ID exceeds max supply
func (k Keeper) checkMaxSupplyPerIdInvariant(ctx sdk.Context, collection *types.TokenCollection, maxSupply types.Uint) error {
	balanceKey := ConstructBalanceKey(types.TotalAddress, collection.CollectionId)
	totalBalance, found := k.GetUserBalanceFromStore(ctx, balanceKey)
	if !found {
		return nil // No balances, invariant passes
	}

	for _, balance := range totalBalance.Balances {
		if balance.Amount.GT(maxSupply) {
			return sdkerrors.Wrapf(types.ErrInvalidRequest,
				"maxSupplyPerId invariant violated: amount %s exceeds max %s",
				balance.Amount.String(), maxSupply.String())
		}
	}
	return nil
}

// checkPostTransferEVMQueryInvariants runs all EVM query challenges for post-transfer invariants
func (k Keeper) checkPostTransferEVMQueryInvariants(
	ctx sdk.Context,
	collection *types.TokenCollection,
	from string,
	toAddresses []string,
	initiatedBy string,
	challenges []*types.EVMQueryChallenge,
) error {
	if len(challenges) == 0 {
		return nil
	}

	// Track total gas across all challenges (DoS protection)
	var totalGasUsed uint64

	for i, challenge := range challenges {
		if challenge == nil {
			continue
		}

		// Replace placeholders in calldata
		calldata, err := k.replacePostTransferInvariantPlaceholders(
			challenge.Calldata,
			from,
			toAddresses,
			initiatedBy,
			collection.CollectionId.String(),
		)
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidRequest,
				"post-transfer invariant %d placeholder replacement: %s", i, err.Error())
		}

		// Decode hex calldata
		calldataBytes, err := hex.DecodeString(strings.TrimPrefix(calldata, "0x"))
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidRequest,
				"invalid calldata hex in post-transfer invariant %d: %s", i, err.Error())
		}

		// Determine gas limit
		gasLimit := DefaultPostTransferEVMQueryGasLimit
		if !challenge.GasLimit.IsZero() {
			gasLimit = challenge.GasLimit.Uint64()
			if gasLimit > MaxPostTransferEVMQueryGasLimit {
				gasLimit = MaxPostTransferEVMQueryGasLimit
			}
		}

		// Check total gas limit
		if totalGasUsed+gasLimit > MaxTotalPostTransferInvariantGas {
			return sdkerrors.Wrapf(types.ErrInvalidRequest,
				"post-transfer invariants exceed total gas limit (%d)", MaxTotalPostTransferInvariantGas)
		}
		totalGasUsed += gasLimit

		// Execute the query with zero address as caller (from/to/initiator passed via calldata placeholders)
		result, err := k.ExecuteEVMQuery(ctx, challenge.ContractAddress, calldataBytes, gasLimit)
		if err != nil {
			return sdkerrors.Wrapf(types.ErrInvalidRequest,
				"post-transfer invariant %d failed: %s", i, err.Error())
		}

		// Compare result if expected result is provided
		if challenge.ExpectedResult != "" {
			expectedBytes, err := hex.DecodeString(strings.TrimPrefix(challenge.ExpectedResult, "0x"))
			if err != nil {
				return sdkerrors.Wrapf(types.ErrInvalidRequest,
					"invalid expected result hex in post-transfer invariant %d", i)
			}

			if !comparePostTransferInvariantResults(result, expectedBytes, challenge.ComparisonOperator) {
				return sdkerrors.Wrapf(types.ErrInvalidRequest,
					"post-transfer invariant %d violated: got %x, expected %x with op %s",
					i, result, expectedBytes, challenge.ComparisonOperator)
			}
		}
	}

	return nil
}

// replacePostTransferInvariantPlaceholders replaces placeholders in calldata
// $sender -> single address
// $recipients -> all recipients concatenated (N x 32 bytes)
// $recipient -> first recipient (for single-recipient convenience)
// $initiator -> initiator address
// $collectionId -> collection ID as uint256
func (k Keeper) replacePostTransferInvariantPlaceholders(
	calldata string,
	from string,
	toAddresses []string,
	initiatedBy string,
	collectionId string,
) (string, error) {
	senderHex, err := postTransferAddressToHexPadded(from)
	if err != nil {
		return "", sdkerrors.Wrapf(err, "sender address")
	}
	initiatorHex, err := postTransferAddressToHexPadded(initiatedBy)
	if err != nil {
		return "", sdkerrors.Wrapf(err, "initiator address")
	}
	calldata = strings.ReplaceAll(calldata, "$sender", senderHex)
	calldata = strings.ReplaceAll(calldata, "$initiator", initiatorHex)
	calldata = strings.ReplaceAll(calldata, "$collectionId", postTransferUint256ToHexPadded(collectionId))

	// $recipient = first recipient (convenience for single-recipient transfers)
	if len(toAddresses) > 0 {
		recipientHex, err := postTransferAddressToHexPadded(toAddresses[0])
		if err != nil {
			return "", sdkerrors.Wrapf(err, "recipient address")
		}
		calldata = strings.ReplaceAll(calldata, "$recipient", recipientHex)
	}

	// $recipients = concatenated addresses (for contracts that need all recipients)
	recipientsHex, err := encodePostTransferAddressArray(toAddresses)
	if err != nil {
		return "", sdkerrors.Wrapf(err, "recipients")
	}
	calldata = strings.ReplaceAll(calldata, "$recipients", recipientsHex)

	return calldata, nil
}

// encodePostTransferAddressArray concatenates addresses for inline replacement in calldata
func encodePostTransferAddressArray(addresses []string) (string, error) {
	var encoded strings.Builder
	for i, addr := range addresses {
		hexPadded, err := postTransferAddressToHexPadded(addr)
		if err != nil {
			return "", sdkerrors.Wrapf(err, "recipient %d", i)
		}
		encoded.WriteString(hexPadded)
	}
	return encoded.String(), nil
}

// postTransferAddressToHexPadded converts an address to 32-byte padded hex for ABI encoding.
// Accepts valid bech32, 0x + 40 hex chars, or protocol addresses (Mint, Total) which get a canonical encoding.
func postTransferAddressToHexPadded(address string) (string, error) {
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

// postTransferUint256ToHexPadded converts a uint string to 32-byte padded hex
func postTransferUint256ToHexPadded(value string) string {
	n := new(big.Int)
	n.SetString(value, 10)
	return fmt.Sprintf("%064x", n)
}

// comparePostTransferInvariantResults compares actual vs expected using the operator
func comparePostTransferInvariantResults(result, expected []byte, operator string) bool {
	switch operator {
	case "", "eq":
		return bytes.Equal(result, expected)
	case "ne":
		return !bytes.Equal(result, expected)
	case "gt", "gte", "lt", "lte":
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
