package types_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

const idTestCreator = "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"

func tooLargeId() sdkmath.Uint {
	return sdkmath.NewUintFromBigInt(new(big.Int).Lsh(big.NewInt(1), 64))
}

// Every message that carries a collection or store id must reject nil and >uint64 ids
// with a validation error instead of panicking later in the keeper.
func TestMessageIdsAreValidated(t *testing.T) {
	build := map[string]func(id sdkmath.Uint) interface{ ValidateBasic() error }{
		"cast vote": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgCastVote{Creator: idTestCreator, CollectionId: id, ApprovalLevel: "collection", ApprovalId: "a", ProposalId: "p", YesWeight: sdkmath.NewUint(1)}
		},
		"delete collection": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgDeleteCollection{Creator: idTestCreator, CollectionId: id}
		},
		"delete dynamic store": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgDeleteDynamicStore{Creator: idTestCreator, StoreId: id}
		},
		"update dynamic store": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgUpdateDynamicStore{Creator: idTestCreator, StoreId: id}
		},
		"set dynamic store value": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgSetDynamicStoreValue{Creator: idTestCreator, StoreId: id, Address: idTestCreator}
		},
		"purge approvals": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgPurgeApprovals{Creator: idTestCreator, CollectionId: id, ApprovalsToPurge: []*types.ApprovalIdentifierDetails{{ApprovalId: "a", ApprovalLevel: "outgoing", ApproverAddress: idTestCreator, Version: sdkmath.NewUint(0)}}}
		},
		"universal update collection": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgUniversalUpdateCollection{Creator: idTestCreator, CollectionId: id}
		},
		"set manager": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgSetManager{Creator: idTestCreator, CollectionId: id, Manager: idTestCreator}
		},
		"transfer tokens": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgTransferTokens{Creator: idTestCreator, CollectionId: id}
		},
		"update user approvals": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgUpdateUserApprovals{Creator: idTestCreator, CollectionId: id}
		},
		"delete incoming approval": func(id sdkmath.Uint) interface{ ValidateBasic() error } {
			return &types.MsgDeleteIncomingApproval{Creator: idTestCreator, CollectionId: id, ApprovalId: "a"}
		},
	}

	for name, mk := range build {
		t.Run(name+" nil id", func(t *testing.T) {
			require.NotPanics(t, func() { require.Error(t, mk(sdkmath.Uint{}).ValidateBasic()) })
		})
		t.Run(name+" id above uint64", func(t *testing.T) {
			require.NotPanics(t, func() { require.Error(t, mk(tooLargeId()).ValidateBasic()) })
		})
	}
}

func TestNilElementsAreRejectedNotDereferenced(t *testing.T) {
	require.NotPanics(t, func() {
		require.Error(t, types.ValidateActionPermission([]*types.ActionPermission{nil}, false))
	})
	require.NotPanics(t, func() {
		require.Error(t, types.ValidateTokenMetadata([]*types.TokenMetadata{nil}, false))
	})
}

func TestSafeAddWithOverflowCheckReturnsError(t *testing.T) {
	max := sdkmath.NewUintFromBigInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)))
	require.NotPanics(t, func() {
		_, err := types.SafeAddWithOverflowCheck(max, sdkmath.OneUint())
		require.ErrorIs(t, err, types.ErrOverflow)
	})
	sum, err := types.SafeAddWithOverflowCheck(sdkmath.NewUint(2), sdkmath.NewUint(3))
	require.NoError(t, err)
	require.True(t, sum.Equal(sdkmath.NewUint(5)))
}
