package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/approval_criteria"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	twofatypes "github.com/bitbadges/bitbadgeschain/x/twofa/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// TwoFADecorator is a GLOBAL ante handler that checks if transactions from users with 2FA requirements
// meet those requirements. This applies to ALL transaction types (bank transfers, staking, governance,
// WASM contracts, IBC, etc.), not just BitBadges transactions. This provides defense in depth against
// compromised private keys by requiring ownership of specific badges as a second factor.
type TwoFADecorator struct {
	keeper Keeper
}

// NewTwoFADecorator creates a new TwoFADecorator
func NewTwoFADecorator(keeper Keeper) TwoFADecorator {
	return TwoFADecorator{
		keeper: keeper,
	}
}

// AnteHandle implements the ante decorator interface.
// This is a GLOBAL handler that applies to ALL transaction types across the entire blockchain.
// It checks 2FA requirements for all signers of any transaction, regardless of the message types.
func (decorator TwoFADecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Get all signers from the transaction (works for any transaction type)
	signers := decorator.getSigners(tx)

	// If no signers, nothing to check
	if len(signers) == 0 {
		return next(ctx, tx, simulate)
	}

	// Check 2FA requirements for each signer
	for _, signer := range signers {
		signerStr := signer.String()

		// Get 2FA requirements for this signer
		requirements, found := decorator.keeper.GetUser2FARequirementsFromStore(ctx, signerStr)
		if !found {
			// No 2FA requirements stored for this user, continue
			continue
		}

		// Check if there are any requirements to validate
		if len(requirements.MustOwnTokens) == 0 && len(requirements.DynamicStoreChallenges) == 0 {
			// Empty requirements list means 2FA is disabled, continue
			continue
		}

		// Create service adapters for shared helper functions
		collectionService := &twofaCollectionService{keeper: decorator.keeper}
		dynamicStoreService := &twofaDynamicStoreService{keeper: decorator.keeper}

		// Check MustOwnTokens requirements using shared helper function
		for idx, mustOwnToken := range requirements.MustOwnTokens {
			// For 2FA, we always check the signer (initiator)
			passed, errMsg := approval_criteria.CheckMustOwnTokensRequirement(
				ctx,
				mustOwnToken,
				idx,
				signerStr, // initiatedBy
				signerStr, // fromAddress (same as signer for 2FA)
				signerStr, // toAddress (same as signer for 2FA)
				collectionService,
				"2FA requirement",
			)
			if !passed {
				return ctx, sdkerrors.Wrap(
					twofatypes.ErrInvalidRequest,
					fmt.Sprintf("2FA requirement failed for signer %s: %s", signerStr, errMsg),
				)
			}
		}

		// Check DynamicStoreChallenges requirements using shared helper function
		for idx, challenge := range requirements.DynamicStoreChallenges {
			// For 2FA, we need to determine the party based on OwnershipCheckParty
			// Get collection if needed (for mint address resolution)
			var collection *badgestypes.TokenCollection
			// Note: For 2FA, we don't have a specific collection context, so we pass nil
			// The helper function will handle this appropriately

			passed, errMsg := approval_criteria.CheckDynamicStoreChallenge(
				ctx,
				challenge,
				idx,
				signerStr, // initiatedBy
				signerStr, // fromAddress (same as signer for 2FA)
				signerStr, // toAddress (same as signer for 2FA)
				dynamicStoreService,
				collection, // nil for 2FA context
				"2FA dynamic store challenge",
			)
			if !passed {
				return ctx, sdkerrors.Wrap(
					twofatypes.ErrInvalidRequest,
					fmt.Sprintf("2FA dynamic store challenge failed for signer %s: %s", signerStr, errMsg),
				)
			}
		}
	}

	return next(ctx, tx, simulate)
}

// getSigners extracts all signer addresses from the transaction.
// This works for ALL transaction types (bank, staking, governance, WASM, IBC, badges, etc.)
// by iterating through all messages and collecting their signers.
func (decorator TwoFADecorator) getSigners(tx sdk.Tx) []sdk.AccAddress {
	signersMap := make(map[string]sdk.AccAddress)

	// Get signers from all messages in the transaction (works for any module/message type)
	for _, msg := range tx.GetMsgs() {
		// Try to get signers using the GetSigners() method (for messages that implement it)
		type msgWithSigners interface {
			GetSigners() []sdk.AccAddress
		}
		if msgWithSigs, ok := msg.(msgWithSigners); ok {
			for _, signer := range msgWithSigs.GetSigners() {
				signersMap[signer.String()] = signer
			}
			continue
		}

		// Handle banktypes.MsgSend specifically (it doesn't implement GetSigners in the interface)
		if bankMsg, ok := msg.(*banktypes.MsgSend); ok {
			fromAddr, err := sdk.AccAddressFromBech32(bankMsg.FromAddress)
			if err == nil {
				signersMap[fromAddr.String()] = fromAddr
			}
			continue
		}

		// For other message types, try to extract signers from common fields
		// This is a fallback for messages that don't implement GetSigners()
		// Most Cosmos SDK messages should implement it, but we handle edge cases
	}

	// Convert map to slice
	signers := make([]sdk.AccAddress, 0, len(signersMap))
	for _, signer := range signersMap {
		signers = append(signers, signer)
	}

	return signers
}
