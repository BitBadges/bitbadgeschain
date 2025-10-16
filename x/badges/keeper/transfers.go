package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

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

func (k Keeper) HandleTransfers(ctx sdk.Context, collection *types.BadgeCollection, transfers []*types.Transfer, initiatedBy string) error {
	var err error

	isArchived := types.GetIsArchived(ctx, collection)
	if isArchived {
		return ErrCollectionIsArchived
	}

	// Validate transfers with invariants
	for _, transfer := range transfers {
		if err := types.ValidateTransferWithInvariants(ctx, transfer, true, collection); err != nil {
			return err
		}
	}

	for _, transfer := range transfers {
		numAttempts := sdkmath.NewUint(1)
		if !transfer.NumAttempts.IsNil() {
			numAttempts = transfer.NumAttempts
		}

		// Convert to uint64 for more efficient loop iteration
		// Add bounds checking to prevent potential issues with very large values
		if numAttempts.GT(sdkmath.NewUint(10000)) {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "numAttempts cannot exceed 10000, got %s", numAttempts.String())
		}

		numAttemptsUint64 := numAttempts.Uint64()
		for i := uint64(0); i < numAttemptsUint64; i++ {
			if i > 0 {
				ctx.GasMeter().ConsumeGas(1000, "HandleTransfers: Gas consumed for each attempt")
			}

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
				protocolFeeTransfers, err := k.CalculateAndDistributeProtocolFees(ctx, coinTransfers, initiatedBy, transfer.AffiliateAddress)
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
	}

	return nil
}

// Step 1: Check if transfer is allowed on collection level (deducting collection approvals if needed). Will return what userApprovals we need to check.
// Step 2: Check necessary approvals on user level (deducting corresponding approvals if needed)
// Step 3: If all good, we can transfer the balances
func (k Keeper) HandleTransfer(
	ctx sdk.Context,
	collection *types.BadgeCollection,
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
	transferMetadata := TransferMetadata{
		From:            from,
		To:              to,
		InitiatedBy:     initiatedBy,
		ApproverAddress: "",
		ApprovalLevel:   "collection",
	}

	transferBalances := types.DeepCopyBalances(transfer.Balances)
	userApprovals, err := k.DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx, collection, transfer, transferMetadata, eventTracking)
	if err != nil {
		return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "collection approvals not satisfied")
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
				AffiliateAddress:                        transfer.AffiliateAddress,
				NumAttempts:                             transfer.NumAttempts,
			}

			if userApproval.Outgoing {
				transferMetadata.ApproverAddress = from
				transferMetadata.ApprovalLevel = "outgoing"

				err = k.DeductUserOutgoingApprovals(ctx, collection, transferBalances, newTransfer, transferMetadata, fromUserBalance, eventTracking, userApproval.UserRoyalties)
				if err != nil {
					return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "outgoing approvals for %s not satisfied", from)
				}
			} else {
				transferMetadata.ApproverAddress = to
				transferMetadata.ApprovalLevel = "incoming"

				err = k.DeductUserIncomingApprovals(ctx, collection, transferBalances, newTransfer, transferMetadata, toUserBalance, eventTracking, userApproval.UserRoyalties)
				if err != nil {
					return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "incoming approvals for %s not satisfied", to)
				}
			}
		}
	}

	for _, balance := range transferBalances {
		//Mint has unlimited balances
		if !types.IsMintAddress(from) {
			fromUserBalance.Balances, err = types.SubtractBalance(ctx, fromUserBalance.Balances, balance, false)
			if err != nil {
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "inadequate balances for transfer from %s", from)
			}
		}

		toUserBalance.Balances, err = types.AddBalance(ctx, toUserBalance.Balances, balance)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
		}
	}

	// Handle special address wrapping/unwrapping to x/bank denominations
	err = k.HandleSpecialAddressWrapping(ctx, collection, transferBalances, from, to)
	if err != nil {
		return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
	}

	// Handle auto-deletions for approvals if necessary
	fromUserBalance, toUserBalance, err = k.HandleAutoDeletions(ctx, collection, fromUserBalance, toUserBalance, *eventTracking.ApprovalsUsed)
	if err != nil {
		return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
	}

	return fromUserBalance, toUserBalance, nil
}
