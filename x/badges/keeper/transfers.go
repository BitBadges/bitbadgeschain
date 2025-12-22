package keeper

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

type TransferMetadata struct {
	To              string
	From            string
	InitiatedBy     string
	ApproverAddress string
	ApprovalLevel   string
}

type ApprovalsUsed struct {
	ApprovalId      string
	ApprovalLevel   string
	ApproverAddress string
	Version         string
}

type CoinTransfers struct {
	From          string
	To            string
	Amount        string
	Denom         string
	IsProtocolFee bool
}

// EventTracking combines approvalsUsed and coinTransfersUsed for cleaner function signatures
type EventTracking struct {
	ApprovalsUsed *[]ApprovalsUsed
	CoinTransfers *[]CoinTransfers
}

func (k Keeper) HandleTransfers(ctx sdk.Context, collection *types.TokenCollection, transfers []*types.Transfer, initiatedBy string) error {
	var err error

	isArchived := types.GetIsArchived(ctx, collection)
	if isArchived {
		return customhookstypes.WrapErr(&ctx, ErrCollectionIsArchived, "collection is currently archived (read-only)")
	}

	// Validate transfers with invariants
	for _, transfer := range transfers {
		if err := types.ValidateTransferWithInvariants(ctx, transfer, true, collection); err != nil {
			// Create deterministic error message without using err.Error()
			detErrMsg := fmt.Sprintf("invariants validation failed for collection %s", collection.CollectionId.String())
			return customhookstypes.WrapErrSimple(&ctx, err, detErrMsg)
		}
	}

	for _, transfer := range transfers {
		fromUserBalance, _ := k.GetBalanceOrApplyDefault(ctx, collection, transfer.From)
		totalMinted := []*types.Balance{}

		for _, to := range transfer.ToAddresses {
			approvalsUsed := []ApprovalsUsed{}
			coinTransfers := []CoinTransfers{}
			eventTracking := &EventTracking{
				ApprovalsUsed: &approvalsUsed,
				CoinTransfers: &coinTransfers,
			}

			toUserBalance, _ := k.GetBalanceOrApplyDefault(ctx, collection, to)

			if transfer.PrecalculateBalancesFromApproval != nil && transfer.PrecalculateBalancesFromApproval.ApprovalId != "" {
				//Here, we precalculate balances from a specified approval
				transferMetadata := TransferMetadata{
					From:            transfer.From,
					To:              to,
					InitiatedBy:     initiatedBy,
					ApproverAddress: "",
					ApprovalLevel:   "collection",
				}
				approvals := collection.CollectionApprovals
				if transfer.PrecalculateBalancesFromApproval.ApprovalLevel == "collection" {
					if transfer.PrecalculateBalancesFromApproval.ApproverAddress != "" {
						return sdkerrors.Wrapf(ErrNotImplemented, "approver address must be blank for collection level approvals")
					}
				} else {
					if transfer.PrecalculateBalancesFromApproval.ApproverAddress != to && transfer.PrecalculateBalancesFromApproval.ApproverAddress != transfer.From {
						return sdkerrors.Wrapf(ErrNotImplemented, "approver address %s must match to or from address for user level precalculations", transfer.PrecalculateBalancesFromApproval.ApproverAddress)
					}

					handled := false
					if transfer.PrecalculateBalancesFromApproval.ApproverAddress == to && transfer.PrecalculateBalancesFromApproval.ApprovalLevel == "incoming" {
						userApprovals := toUserBalance.IncomingApprovals
						approvals = types.CastIncomingTransfersToCollectionTransfers(userApprovals, to)
						handled = true
					}

					if transfer.PrecalculateBalancesFromApproval.ApprovalLevel == "outgoing" && !handled && transfer.PrecalculateBalancesFromApproval.ApproverAddress == transfer.From {
						userApprovals := fromUserBalance.OutgoingApprovals
						approvals = types.CastOutgoingTransfersToCollectionTransfers(userApprovals, transfer.From)
						handled = true
					}

					if !handled {
						return sdkerrors.Wrapf(ErrNotImplemented, "could not determine approval to precalculate from %s", transfer.PrecalculateBalancesFromApproval.ApproverAddress)
					}
				}

				//Precaluclate the balances that will be transferred
				transfer.Balances, err = k.GetPredeterminedBalancesForPrecalculationId(
					ctx,
					collection,
					approvals,
					transfer,
					transferMetadata,
				)
				if err != nil {
					return err
				}

				//TODO: Deprecate this in favor of actually calculating the balances in indexer
				amountsJsonData, err := json.Marshal(transfer)
				if err != nil {
					return err
				}
				amountsStr := string(amountsJsonData)

				EmitMessageAndIndexerEvents(ctx,
					sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
					sdk.NewAttribute("creator", initiatedBy),
					sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
					sdk.NewAttribute("transfer", amountsStr),
				)
			}

			if types.IsMintAddress(transfer.From) {
				copiedBalances := types.DeepCopyBalances(transfer.Balances)
				totalMinted, err = types.AddBalances(ctx, totalMinted, copiedBalances)
				if err != nil {
					return err
				}
			}

			// Run global transfer checkers before HandleTransfer()
			// All must pass for the transfer to be allowed
			// In case of failure, whole transfer fails so it is rolled back
			for _, provider := range k.customGlobalTransferCheckerProviders {
				checkers := provider(ctx, transfer.From, to, initiatedBy, collection, transfer.Balances, transfer.Memo)
				for _, checker := range checkers {
					balances := types.DeepCopyBalances(transfer.Balances)
					detErrMsg, err := checker.Check(ctx, transfer.From, to, initiatedBy, collection, balances, transfer.Memo)
					if err != nil {
						if detErrMsg != "" {
							return sdkerrors.Wrapf(err, "%s: %s", checker.Name(), detErrMsg)
						}
						return sdkerrors.Wrapf(err, "%s: global transfer check failed", checker.Name())
					}
				}
			}

			fromUserBalance, toUserBalance, err = k.HandleTransfer(
				ctx,
				collection,
				transfer,
				fromUserBalance,
				toUserBalance,
				transfer.From,
				to,
				initiatedBy,
				eventTracking,
			)
			if err != nil {
				return err
			}

			if err := k.SetBalanceForAddress(ctx, collection, to, toUserBalance); err != nil {
				return err
			}

			// Calculate and distribute protocol fees
			protocolFeeTransfers, err := k.CalculateAndDistributeProtocolFees(ctx, coinTransfers, initiatedBy)
			if err != nil {
				return err
			}

			// Add protocol fee transfers to the main coinTransfers slice
			coinTransfers = append(coinTransfers, protocolFeeTransfers...)

			err = EmitUsedApprovalDetailsEvent(ctx, collection.CollectionId, transfer.From, to, initiatedBy, coinTransfers, approvalsUsed, transfer.Balances)
			if err != nil {
				return err
			}
		}

		if !types.IsMintAddress(transfer.From) {
			if err := k.SetBalanceForAddress(ctx, collection, transfer.From, fromUserBalance); err != nil {
				return err
			}
		} else {
			// Get current Total and increment it
			totalBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, types.TotalAddress)
			totalBalances.Balances, err = types.AddBalances(ctx, totalBalances.Balances, totalMinted)
			if err != nil {
				return err
			}

			if err := k.SetBalanceForAddress(ctx, collection, types.TotalAddress, totalBalances); err != nil {
				return err
			}
		}
	}

	return nil
}

// Step 1: Check if transfer is allowed on collection level (deducting collection approvals if needed). Will return what userApprovals we need to check.
// Step 2: Check necessary approvals on user level (deducting corresponding approvals if needed)
// Step 3: If all good, we can transfer the balances
func (k Keeper) HandleTransfer(
	ctx sdk.Context,
	collection *types.TokenCollection,
	transfer *types.Transfer,
	fromUserBalance *types.UserBalanceStore,
	toUserBalance *types.UserBalanceStore,
	from string,
	to string,
	initiatedBy string,
	eventTracking *EventTracking,
) (*types.UserBalanceStore, *types.UserBalanceStore, error) {
	if to == from {
		return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(ErrNotImplemented, "cannot send to self")
	}

	var err error

	// Check if sending from a special address (cosmos coin wrapper path) and set up one-time approval
	isSendingFromSpecialAddress := k.IsBackedOrWrappingPathAddress(ctx, collection, from)
	isSendingToSpecialAddress := k.IsBackedOrWrappingPathAddress(ctx, collection, to)

	// Little hacky way of doing this, but any transfer on the collection level that involves special addresses
	// should always specify the collection approvals
	if isSendingFromSpecialAddress || isSendingToSpecialAddress {
		transfer.OnlyCheckPrioritizedCollectionApprovals = true
	}

	if isSendingFromSpecialAddress {
		// Set up one-time outgoing approval for the special address to send tokens to the recipient
		// Similar to gamm pools
		currBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, from)

		// Add one-time outgoing approval
		oneTimeApproval := &types.UserOutgoingApproval{
			ToListId:          to,
			InitiatedByListId: initiatedBy,
			TransferTimes:     []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
			OwnershipTimes:    []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
			TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)}},
			Version:           sdkmath.NewUint(0),
			ApprovalId:        "one-time-outgoing",
		}

		// Set version for the approval
		oneTimeApproval.Version = k.IncrementApprovalVersion(ctx, collection.CollectionId, "outgoing", from, oneTimeApproval.ApprovalId)

		// Add to outgoing approvals
		currBalances.OutgoingApprovals = []*types.UserOutgoingApproval{oneTimeApproval}

		// Save the balance
		err = k.SetBalanceForAddress(ctx, collection, from, currBalances)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "failed to set one-time approval for special address")
		}

		// Update fromUserBalance to include the new approval
		fromUserBalance = currBalances
	}

	transferMetadata := TransferMetadata{
		From:            from,
		To:              to,
		InitiatedBy:     initiatedBy,
		ApproverAddress: "",
		ApprovalLevel:   "collection",
	}

	transferBalances := types.DeepCopyBalances(transfer.Balances)
	userApprovals, err := k.DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx, collection, transfer, transferMetadata, eventTracking, "collection")
	if err != nil {
		return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "collection transferability error")
	}

	if len(userApprovals) > 0 {
		for _, userApproval := range userApprovals {
			newTransfer := &types.Transfer{
				From:                                    from,
				ToAddresses:                             []string{to},
				Balances:                                userApproval.Balances,
				MerkleProofs:                            transfer.MerkleProofs,
				PrioritizedApprovals:                    transfer.PrioritizedApprovals,
				OnlyCheckPrioritizedCollectionApprovals: transfer.OnlyCheckPrioritizedCollectionApprovals,
				OnlyCheckPrioritizedIncomingApprovals:   transfer.OnlyCheckPrioritizedIncomingApprovals,
				OnlyCheckPrioritizedOutgoingApprovals:   transfer.OnlyCheckPrioritizedOutgoingApprovals,
				PrecalculationOptions:                   transfer.PrecalculationOptions,
			}

			if userApproval.Outgoing {
				transferMetadata.ApproverAddress = from
				transferMetadata.ApprovalLevel = "outgoing"

				err = k.DeductUserOutgoingApprovals(ctx, collection, transferBalances, newTransfer, transferMetadata, fromUserBalance, eventTracking, userApproval.UserRoyalties)
				if err != nil {
					return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "outgoing approvals for %s", from)
				}
			} else {
				transferMetadata.ApproverAddress = to
				transferMetadata.ApprovalLevel = "incoming"

				err = k.DeductUserIncomingApprovals(ctx, collection, transferBalances, newTransfer, transferMetadata, toUserBalance, eventTracking, userApproval.UserRoyalties)
				if err != nil {
					return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "incoming approvals for %s", to)
				}
			}
		}
	}

	for _, balance := range transferBalances {
		//Mint has unlimited balances
		//Backed addresses (wrapper and backed paths) have unlimited balances
		if !types.IsMintAddress(from) && !k.IsSpecialBackedAddress(ctx, collection, from) {
			fromUserBalance.Balances, err = types.SubtractBalance(ctx, fromUserBalance.Balances, balance, false)
			if err != nil {
				detErrMsg := fmt.Sprintf("inadequate balances for transfer from %s", from)
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, customhookstypes.WrapErrSimple(&ctx, err, detErrMsg)
			}
		}

		if !k.IsSpecialBackedAddress(ctx, collection, to) {
			toUserBalance.Balances, err = types.AddBalance(ctx, toUserBalance.Balances, balance)
			if err != nil {
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
			}
		}
	}

	// Handle special address wrapping/unwrapping to x/bank denominations
	err = k.HandleSpecialAddressWrapping(ctx, collection, transferBalances, from, to, initiatedBy)
	if err != nil {
		return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
	}

	// Handle special address backing/unbacking to x/bank denominations
	err = k.HandleSpecialAddressBacking(ctx, collection, transferBalances, from, to, initiatedBy)
	if err != nil {
		return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
	}

	// Handle auto-deletions for approvals if necessary
	fromUserBalance, toUserBalance, err = k.HandleAutoDeletions(ctx, collection, fromUserBalance, toUserBalance, *eventTracking.ApprovalsUsed)
	if err != nil {
		return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
	}

	// Clean up one-time approval if sending from special address
	if isSendingFromSpecialAddress {
		// Remove the one-time outgoing approval
		// This is needed as opposed to auto-deletion because technically the approval might not
		// be used if there is some forceful override (thus never deletes and we have a dangling approval)
		fromUserBalance.OutgoingApprovals = []*types.UserOutgoingApproval{}

		// Save the updated balance
		err = k.SetBalanceForAddress(ctx, collection, from, fromUserBalance)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "failed to clean up one-time approval for special address")
		}
	}

	return fromUserBalance, toUserBalance, nil
}
