package keeper

import (
	"bytes"
	"reflect"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	"github.com/gogo/protobuf/proto"
)

// compareUintRanges compares two slices of UintRange for equality
func compareUintRanges(a, b []*types.UintRange) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Start.Equal(b[i].Start) || !a[i].End.Equal(b[i].End) {
			return false
		}
	}
	return true
}

// criteriaIsAbsent reports whether m carries no criteria at all.
//
// A plain `m == nil` is not enough and this is not a hypothetical: it panicked
// on mainnet. approvalCriteria is optional, so an approval that restricts
// nothing holds a nil *ApprovalCriteria. Assigning that to a proto.Message
// parameter produces an interface that is NOT nil — it is
// (type=*ApprovalCriteria, value=nil) — so `m == nil` is false, and
// proto.Marshal then calls MarshalToSizedBuffer on a nil receiver and
// dereferences it:
//
//	panic recovered in runTx err="recovered: runtime error: invalid memory
//	address or nil pointer dereference"
//	(*ApprovalCriteria).MarshalToSizedBuffer(0x0, ...)
//
// The interface cannot be avoided here — three distinct criteria types
// (ApprovalCriteria, OutgoingApprovalCriteria, IncomingApprovalCriteria) share
// this helper — so the nil has to be unwrapped by reflection instead.
func criteriaIsAbsent(m proto.Message) bool {
	if m == nil {
		return true
	}
	v := reflect.ValueOf(m)
	return v.Kind() == reflect.Ptr && v.IsNil()
}

// compareApprovalCriteria compares two protobuf messages for equality using
// canonical binary marshaling. proto.MarshalTextString is not deterministic
// across gogo/protobuf versions (whitespace, field ordering), so binary Marshal
// is used instead — it is the same encoding used for consensus state and is
// stable across nodes.
func compareApprovalCriteria(a, b proto.Message) bool {
	aAbsent, bAbsent := criteriaIsAbsent(a), criteriaIsAbsent(b)
	if aAbsent || bAbsent {
		// Absent on both sides is equal; absent on one side is not. An absent
		// criteria and a present-but-empty one are different states on the
		// wire, and changing between them is a real change to the approval.
		return aAbsent && bAbsent
	}
	aBytes, errA := proto.Marshal(a)
	bBytes, errB := proto.Marshal(b)
	if errA != nil || errB != nil {
		return false
	}
	return bytes.Equal(aBytes, bBytes)
}

// collectionApprovalEqual compares two CollectionApproval objects for equality,
// excluding the Version field (which is what we're trying to determine).
// Returns true if all fields except Version are equal.
func collectionApprovalEqual(a, b *types.CollectionApproval) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare all fields except Version
	return a.ApprovalId == b.ApprovalId &&
		a.FromListId == b.FromListId &&
		a.ToListId == b.ToListId &&
		a.InitiatedByListId == b.InitiatedByListId &&
		a.Uri == b.Uri &&
		a.CustomData == b.CustomData &&
		compareUintRanges(a.TransferTimes, b.TransferTimes) &&
		compareUintRanges(a.TokenIds, b.TokenIds) &&
		compareUintRanges(a.OwnershipTimes, b.OwnershipTimes) &&
		compareApprovalCriteria(a.ApprovalCriteria, b.ApprovalCriteria)
}

// userOutgoingApprovalEqual compares two UserOutgoingApproval objects for equality,
// excluding the Version field (which is what we're trying to determine).
// Returns true if all fields except Version are equal.
func userOutgoingApprovalEqual(a, b *types.UserOutgoingApproval) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare all fields except Version
	return a.ApprovalId == b.ApprovalId &&
		a.ToListId == b.ToListId &&
		a.InitiatedByListId == b.InitiatedByListId &&
		a.Uri == b.Uri &&
		a.CustomData == b.CustomData &&
		compareUintRanges(a.TransferTimes, b.TransferTimes) &&
		compareUintRanges(a.TokenIds, b.TokenIds) &&
		compareUintRanges(a.OwnershipTimes, b.OwnershipTimes) &&
		compareApprovalCriteria(a.ApprovalCriteria, b.ApprovalCriteria)
}

// userIncomingApprovalEqual compares two UserIncomingApproval objects for equality,
// excluding the Version field (which is what we're trying to determine).
// Returns true if all fields except Version are equal.
func userIncomingApprovalEqual(a, b *types.UserIncomingApproval) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare all fields except Version
	return a.ApprovalId == b.ApprovalId &&
		a.FromListId == b.FromListId &&
		a.InitiatedByListId == b.InitiatedByListId &&
		a.Uri == b.Uri &&
		a.CustomData == b.CustomData &&
		compareUintRanges(a.TransferTimes, b.TransferTimes) &&
		compareUintRanges(a.TokenIds, b.TokenIds) &&
		compareUintRanges(a.OwnershipTimes, b.OwnershipTimes) &&
		compareApprovalCriteria(a.ApprovalCriteria, b.ApprovalCriteria)
}
