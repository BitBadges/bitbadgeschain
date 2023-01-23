package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func TestMsgUpdateUris_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateUris
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateUris{
				Creator: "invalid_address",
				CollectionUri: "https://facebook.com",
				BadgeUri: "https://facebook.com",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateUris{
				Creator: sample.AccAddress(),
				CollectionUri: "https://facebook.com",
				BadgeUri: "https://facebook.com/{id}",
			},
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
