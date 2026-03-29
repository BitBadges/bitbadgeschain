package keeper

import (
	"strings"
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func TestCollectRejectionMessages_WithMessage(t *testing.T) {
	approvals := []*types.CollectionApproval{
		{
			ApprovalId: "approval-1",
			ApprovalCriteria: &types.ApprovalCriteria{
				RejectionMessage: "Only accredited investors may transfer",
			},
		},
	}
	result := collectRejectionMessages([]int{0}, approvals)
	if !strings.Contains(result, "Only accredited investors may transfer") {
		t.Errorf("expected rejection message in output, got: %s", result)
	}
	if !strings.Contains(result, "approval-1") {
		t.Errorf("expected approval ID in output, got: %s", result)
	}
}

func TestCollectRejectionMessages_EmptyMessage(t *testing.T) {
	approvals := []*types.CollectionApproval{
		{
			ApprovalId: "approval-1",
			ApprovalCriteria: &types.ApprovalCriteria{
				RejectionMessage: "",
			},
		},
	}
	result := collectRejectionMessages([]int{0}, approvals)
	if result != "" {
		t.Errorf("expected empty string for empty rejection message, got: %s", result)
	}
}

func TestCollectRejectionMessages_NilCriteria(t *testing.T) {
	approvals := []*types.CollectionApproval{
		{
			ApprovalId:       "approval-1",
			ApprovalCriteria: nil,
		},
	}
	result := collectRejectionMessages([]int{0}, approvals)
	if result != "" {
		t.Errorf("expected empty string for nil criteria, got: %s", result)
	}
}

func TestCollectRejectionMessages_MultipleApprovals(t *testing.T) {
	approvals := []*types.CollectionApproval{
		{
			ApprovalId: "approval-1",
			ApprovalCriteria: &types.ApprovalCriteria{
				RejectionMessage: "KYC required",
			},
		},
		{
			ApprovalId: "approval-2",
			ApprovalCriteria: &types.ApprovalCriteria{
				RejectionMessage: "Transfer locked until 2027",
			},
		},
		{
			ApprovalId: "approval-3",
			ApprovalCriteria: &types.ApprovalCriteria{
				RejectionMessage: "",
			},
		},
	}
	// Only indices 0 and 1 were checked
	result := collectRejectionMessages([]int{0, 1}, approvals)
	if !strings.Contains(result, "KYC required") {
		t.Errorf("expected first rejection message, got: %s", result)
	}
	if !strings.Contains(result, "Transfer locked until 2027") {
		t.Errorf("expected second rejection message, got: %s", result)
	}
	// approval-3 was not checked, so its (empty) message should not appear
}

func TestCollectRejectionMessages_OutOfBoundsIndex(t *testing.T) {
	approvals := []*types.CollectionApproval{
		{
			ApprovalId: "approval-1",
			ApprovalCriteria: &types.ApprovalCriteria{
				RejectionMessage: "Restricted",
			},
		},
	}
	// Index 5 is out of bounds - should not panic
	result := collectRejectionMessages([]int{0, 5}, approvals)
	if !strings.Contains(result, "Restricted") {
		t.Errorf("expected rejection message from valid index, got: %s", result)
	}
}

func TestBuildPotentialErrorsString_IncludesRejectionMessage(t *testing.T) {
	approvals := []*types.CollectionApproval{
		{
			ApprovalId: "test-approval",
			ApprovalCriteria: &types.ApprovalCriteria{
				RejectionMessage: "Transfers paused for maintenance",
			},
		},
	}
	errorsWithIdx := []ErrorWithIdx{
		{
			ErrorMsgs: []string{"addresses do not match"},
			Idx:       0,
		},
	}
	result := buildPotentialErrorsString(
		[]string{},       // no prioritized errors
		[]int{0},         // checked approval at index 0
		errorsWithIdx,
		approvals,
		approvals,
	)
	if !strings.Contains(result, "Transfers paused for maintenance") {
		t.Errorf("expected rejection message in error string, got: %s", result)
	}
}

func TestBuildPotentialErrorsString_NoRejectionMessage(t *testing.T) {
	approvals := []*types.CollectionApproval{
		{
			ApprovalId: "test-approval",
			ApprovalCriteria: &types.ApprovalCriteria{
				RejectionMessage: "",
			},
		},
	}
	errorsWithIdx := []ErrorWithIdx{
		{
			ErrorMsgs: []string{"addresses do not match"},
			Idx:       0,
		},
	}
	result := buildPotentialErrorsString(
		[]string{},
		[]int{0},
		errorsWithIdx,
		approvals,
		approvals,
	)
	// Should not contain " | " separator since there's no rejection message
	if strings.Contains(result, " | ") {
		t.Errorf("expected no rejection message separator, got: %s", result)
	}
}

func TestBuildPotentialErrorsString_PrioritizedWithRejectionMessage(t *testing.T) {
	approvals := []*types.CollectionApproval{
		{
			ApprovalId: "prio-approval",
			ApprovalCriteria: &types.ApprovalCriteria{
				RejectionMessage: "Only whitelisted addresses",
			},
		},
	}
	errorsWithIdx := []ErrorWithIdx{
		{
			ErrorMsgs: []string{"merkle challenge error"},
			Idx:       0,
		},
	}
	result := buildPotentialErrorsString(
		[]string{"merkle challenge error"}, // prioritized error
		[]int{0},
		errorsWithIdx,
		approvals,
		approvals,
	)
	if !strings.Contains(result, "Only whitelisted addresses") {
		t.Errorf("expected rejection message in prioritized error string, got: %s", result)
	}
	if !strings.Contains(result, "prioritized approvals") {
		t.Errorf("expected prioritized label in error string, got: %s", result)
	}
}
