package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/trevormil/bitbadgeschain/testutil/sample"

	"github.com/trevormil/bitbadgeschain/x/badges/types"
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
				Creator:      "invalid_address",
				Uri:          "https://bitbadge.com/badge.svg",
				SubassetUris: "https://bitbadge.com/badge.svg",
				Permissions:  15,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: types.MsgNewBadge{
				Creator:      sample.AccAddress(),
				Uri:          "https://bitbadge.com/badge.svg",
				SubassetUris: "https://bitbadge.com/badge.svg",
				Permissions:  15,
			},
		}, {
			name: "invalid URI",
			msg: types.MsgNewBadge{
				Creator:      sample.AccAddress(),
				Uri:          "ht",
				SubassetUris: "https://bitbadge.com/badge.svg",
				Permissions:  15,
			},

			err: types.ErrInvalidBadgeURI,
		},
		{
			name: "invalid Subasset URI",
			msg: types.MsgNewBadge{
				Creator:      sample.AccAddress(),
				Uri:          "http://x.com",
				SubassetUris: "sfd",
				Permissions:  15,
			},
			err: types.ErrInvalidBadgeURI,
		},
		{
			name: "invalid Permissions",
			msg: types.MsgNewBadge{
				Creator:      sample.AccAddress(),
				Uri:          "http://x.com",
				SubassetUris: "http://x.com",
				Permissions:  10000,
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
