package keeper

import (
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

// compareApprovalCriteria compares two ApprovalCriteria values for
// semantic equality under the nullable=true schema. Old storage bytes
// (written under the pre-upgrade nullable=false binary) carry explicit
// `&Empty{}` sub-messages with length-0 tags; new submissions via the
// SDK omit those tags. types.ProtoEqualNullableAware collapses both
// to a canonical form before comparing, so a no-op resubmission of an
// existing approval doesn't falsely register as "edited."
func compareApprovalCriteria(a, b proto.Message) bool {
	return types.ProtoEqualNullableAware(a, b)
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
