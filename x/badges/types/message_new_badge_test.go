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
		msg  types.MsgNewBadge
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgNewBadge{
				Creator: "invalid_address",
				Uri: &types.UriObject{
					Uri:                    []byte("example.com/"),
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					BytesToInsert:          []byte{},
					InsertIdIdx:            10,
				},
				Permissions:         15,
				FreezeAddressRanges: []*types.IdRange{},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgNewBadge{
				Creator: sample.AccAddress(),
				Uri: &types.UriObject{
					Uri:                    []byte("example.com/"),
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					BytesToInsert:          []byte{},
					InsertIdIdx:            10,
				},
				Permissions:         15,
				FreezeAddressRanges: []*types.IdRange{},
			},
		}, {
			name: "invalid URI",
			msg: types.MsgNewBadge{
				Creator: sample.AccAddress(),
				Uri: &types.UriObject{
					Uri:                    []byte(""),
					Scheme:                 0,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					BytesToInsert:          []byte{},
					InsertIdIdx:            10,
				},
				Permissions:         15,
				FreezeAddressRanges: []*types.IdRange{},
			},

			err: types.ErrInvalidBadgeURI,
		},
		{
			name: "invalid Subasset URI",
			msg: types.MsgNewBadge{
				Creator: sample.AccAddress(),
				Uri: &types.UriObject{
					Uri:                    []byte("example.com/"),
					Scheme:                 0,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					BytesToInsert:          []byte("  "),
					InsertIdIdx:            10,
				},
				Permissions:         15,
				FreezeAddressRanges: []*types.IdRange{},
			},
			err: types.ErrInvalidBadgeURI,
		},
		{
			name: "invalid Permissions",
			msg: types.MsgNewBadge{
				Creator: sample.AccAddress(),
				Uri: &types.UriObject{
					Uri:                    []byte("example.com/"),
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					BytesToInsert:          []byte{},
					InsertIdIdx:            10,
				},
				Permissions:         10000,
				FreezeAddressRanges: []*types.IdRange{},
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
