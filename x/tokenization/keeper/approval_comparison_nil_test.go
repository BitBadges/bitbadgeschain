package keeper

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// TestCompareApprovalCriteriaHandlesAbsentCriteria reproduces a panic seen on
// mainnet (v33, cosmos-sdk v0.53.6) while simulating MsgUniversalUpdateCollection:
//
//	panic recovered in runTx err="recovered: runtime error: invalid memory
//	address or nil pointer dereference"
//	(*ApprovalCriteria).MarshalToSizedBuffer(0x0, ...)
//	gogo/protobuf/proto.Marshal({0x5c4f110, 0x0})
//	keeper.compareApprovalCriteria(...) approval_comparison.go:35
//
// approvalCriteria is an optional field, so an approval that does not restrict
// transfers carries a nil pointer. compareApprovalCriteria takes proto.Message,
// and a nil *ApprovalCriteria stored in an interface is NOT an interface nil —
// it is (type=*ApprovalCriteria, value=nil). So both nil guards fall through
// and proto.Marshal dereferences the nil receiver.
//
// Each subtest passes a typed nil the way the three real call sites do
// (collectionApprovalEqual, userOutgoingApprovalEqual, userIncomingApprovalEqual),
// because the bug lives in the conversion to the interface, not in the values.
func TestCompareApprovalCriteriaHandlesAbsentCriteria(t *testing.T) {
	t.Run("collection", func(t *testing.T) {
		var absent *types.ApprovalCriteria
		if !compareApprovalCriteria(absent, absent) {
			t.Fatal("two absent collection criteria must compare equal")
		}
		if compareApprovalCriteria(absent, &types.ApprovalCriteria{}) {
			t.Fatal("an absent criteria must not equal a present empty one")
		}
	})

	t.Run("outgoing", func(t *testing.T) {
		var absent *types.OutgoingApprovalCriteria
		if !compareApprovalCriteria(absent, absent) {
			t.Fatal("two absent outgoing criteria must compare equal")
		}
	})

	t.Run("incoming", func(t *testing.T) {
		var absent *types.IncomingApprovalCriteria
		if !compareApprovalCriteria(absent, absent) {
			t.Fatal("two absent incoming criteria must compare equal")
		}
	})
}

// TestCollectionApprovalEqualWithAbsentCriteria drives the real call site rather
// than the helper, which is the path MsgUniversalUpdateCollection takes when it
// compares a submitted approval against the stored one to decide whether the
// version changed.
func TestCollectionApprovalEqualWithAbsentCriteria(t *testing.T) {
	approval := func() *types.CollectionApproval {
		return &types.CollectionApproval{
			ApprovalId:        "unrestricted",
			FromListId:        "All",
			ToListId:          "All",
			InitiatedByListId: "All",
			ApprovalCriteria:  nil,
		}
	}

	if !collectionApprovalEqual(approval(), approval()) {
		t.Fatal("an approval with no criteria must equal itself; a panic here is the mainnet bug")
	}
}
