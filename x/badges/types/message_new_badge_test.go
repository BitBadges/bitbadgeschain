package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/trevormil/bitbadgeschain/testutil/sample"
)

func TestMsgNewBadge_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgNewBadge
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgNewBadge{
				Creator: 				"invalid_address",
				Uri:    				"https://bitbadge.com/badge.svg",
				SubassetUris: 			"https://bitbadge.com/badge.svg",
				Permissions:			15,
				FreezeAddressesDigest: 	"",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid state",
			msg: MsgNewBadge{
				Creator: sample.AccAddress(),
				Uri:    				"https://bitbadge.com/badge.svg",
				SubassetUris: 			"https://bitbadge.com/badge.svg",
				Permissions:			15,
				FreezeAddressesDigest: 	"",
			},
		}, {
			name: "invalid URI",
			msg: MsgNewBadge{
				Creator: sample.AccAddress(),
				Uri:    				"ht",
				SubassetUris: 			"https://bitbadge.com/badge.svg",
				Permissions:			15,
				FreezeAddressesDigest: 	"",
			},
			
			err: ErrInvalidBadgeURI,
		},
		{
			name: "invalid Subasset URI",
			msg: MsgNewBadge{
				Creator: sample.AccAddress(),
				Uri:    				"http://x.com",
				SubassetUris: 			"sfd",
				Permissions:			15,
				FreezeAddressesDigest: 	"",
			},
			err: ErrInvalidBadgeURI,
		},
		{
			name: "invalid Permissions",
			msg: MsgNewBadge{
				Creator: sample.AccAddress(),
				Uri:    				"http://x.com",
				SubassetUris: 			"http://x.com",
				Permissions:			10000,
				FreezeAddressesDigest: 	"",
			},
			err: ErrInvalidPermissionsLeadingZeroes,
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
