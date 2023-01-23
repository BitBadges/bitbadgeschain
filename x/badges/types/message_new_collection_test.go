package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestMsgNewBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgNewCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgNewCollection{
				Creator: "invalid_address",
				CollectionUri: "https://example.com",
				BadgeUri: "https://example.com/{id}",
				Permissions:         15,
				DisallowedTransfers: []*types.TransferMapping{},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgNewCollection{
				Creator: sample.AccAddress(),
				CollectionUri: "https://example.com",
				BadgeUri: "https://example.com/{id}",
				Permissions:         15,
				DisallowedTransfers: []*types.TransferMapping{},
			},
		}, {
			name: "invalid URI",
			msg: types.MsgNewCollection{
				Creator: sample.AccAddress(),
				CollectionUri: "",
				BadgeUri: "https://example.com/{id}",
				Permissions:         15,
				DisallowedTransfers: []*types.TransferMapping{},
			},

			err: types.ErrInvalidBadgeURI,
		},
		{
			name: "invalid Badge URI",
			msg: types.MsgNewCollection{
				Creator: sample.AccAddress(),
				CollectionUri: "https://example.com",
				BadgeUri: "",
				Permissions:         15,
				DisallowedTransfers: []*types.TransferMapping{},
			},
			err: types.ErrInvalidBadgeURI,
		},
		{
			name: "invalid Permissions",
			msg: types.MsgNewCollection{
				Creator: sample.AccAddress(),
				CollectionUri: "https://example.com",
				BadgeUri: "https://example.com/{id}",
				Permissions:         10000,
				DisallowedTransfers: []*types.TransferMapping{},
			},
			err: types.ErrInvalidPermissionsLeadingZeroes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
