package keeper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ErrorWithIdx tracks error messages by approval index
type ErrorWithIdx struct {
	ErrorMsgs []string
	Idx       int
}

// addPartialSuccessErrors adds error messages for approvals that had partial success
// (i.e., approvals that were used but didn't fully satisfy all balance requirements)
func addPartialSuccessErrors(
	errorsWithIdx []ErrorWithIdx,
	approvalsUsed []ApprovalsUsed,
	approvals []*types.CollectionApproval,
) []ErrorWithIdx {
	for _, approvalUsed := range approvalsUsed {
		matchingIdx := -1
		for i, approval := range approvals {
			if approvalUsed.ApprovalId == approval.ApprovalId {
				matchingIdx = i
				break
			}
		}

		// Skip if we can't find the approval in our approvals array
		if matchingIdx == -1 {
			continue
		}

		matchingErrorElementIdx := -1
		for i, errorWithIdx := range errorsWithIdx {
			if errorWithIdx.Idx == matchingIdx {
				matchingErrorElementIdx = i
				break
			}
		}

		if matchingErrorElementIdx == -1 {
			errorsWithIdx = append(errorsWithIdx, ErrorWithIdx{ErrorMsgs: []string{}, Idx: matchingIdx})
			matchingErrorElementIdx = len(errorsWithIdx) - 1
		}

		if len(errorsWithIdx[matchingErrorElementIdx].ErrorMsgs) == 0 {
			errorsWithIdx[matchingErrorElementIdx].ErrorMsgs = append(errorsWithIdx[matchingErrorElementIdx].ErrorMsgs, "approval had partial success but not full success for all balances")
		}
	}

	return errorsWithIdx
}

// buildPotentialErrorsString builds the potential errors string based on whether there are prioritized errors
// or auto-scan errors
func buildPotentialErrorsString(
	potentialErrors []string,
	approvalIdxsChecked []int,
	errorsWithIdx []ErrorWithIdx,
) string {
	if len(potentialErrors) > 0 {
		return " - errors w/ prioritized approvals: " + strings.Join(potentialErrors, ", ")
	}

	approvalsWithPluralConditional := "approval"
	if len(approvalIdxsChecked) > 1 {
		approvalsWithPluralConditional = "approvals"
	}
	approvalIdxsCheckedStr := []string{}
	for _, idx := range approvalIdxsChecked {
		approvalIdxsCheckedStr = append(approvalIdxsCheckedStr, strconv.Itoa(idx))
	}
	potentialErrorsStr := " - auto-scan failed (checked " + strconv.Itoa(len(approvalIdxsChecked)) + " potentially matching " + approvalsWithPluralConditional
	if len(approvalIdxsChecked) > 0 {
		potentialErrorsStr += ", idxs: " + strings.Join(approvalIdxsCheckedStr, ", ")
	}
	potentialErrorsStr += ")"

	// Try to be smart about error logs. If we only checked one approval via auto-scan, we can just log the errors for that approval
	if len(approvalIdxsChecked) == 1 {
		idxChecked := approvalIdxsChecked[0]
		errorsForIdx := []string{}
		for _, errorWithIdx := range errorsWithIdx {
			if errorWithIdx.Idx == idxChecked {
				errorsForIdx = errorWithIdx.ErrorMsgs
				break
			}
		}

		potentialErrorsStr = ": errors for only matching approval idx " + strconv.Itoa(idxChecked) + ": " + strings.Join(errorsForIdx, ", ")
	}

	return potentialErrorsStr
}

// buildTransferString builds a descriptive string for the transfer attempt
func buildTransferString(
	remainingBalances []*types.Balance,
	fromAddress string,
	toAddress string,
	initiatedBy string,
) string {
	return fmt.Sprintf("attempting to transfer ID %s from %s to %s initiated by %s",
		remainingBalances[0].TokenIds[0].Start.String(), fromAddress, toAddress, initiatedBy)
}

// buildApprovalFailureError builds the final error message for failed approval checks
func buildApprovalFailureError(
	ctx sdk.Context,
	approvalLevel string,
	transferStr string,
	potentialErrorsStr string,
) error {
	detErrMsg := fmt.Sprintf("%s approvals not satisfied: %s%s", approvalLevel, transferStr, potentialErrorsStr)
	return customhookstypes.WrapErrSimple(&ctx, ErrInadequateApprovals, detErrMsg)
}

// buildMultipleRoyaltiesError builds an error message when multiple different user royalty percentages are found
func buildMultipleRoyaltiesError(ctx sdk.Context) error {
	detErrMsg := "multiple user-level royalties found - please split your transfer up to use one collection approval w/ royalty per transfer"
	return customhookstypes.WrapErrSimple(&ctx, ErrInadequateApprovals, detErrMsg)
}
