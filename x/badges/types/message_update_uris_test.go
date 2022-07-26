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
				Uri: &types.UriObject{
					Uri:                    "example.com/",
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					BytesToInsert:          "",
					InsertIdIdx:            10,
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateUris{
				Creator: sample.AccAddress(),
				Uri: &types.UriObject{
					Uri:                    "example.com/",
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					BytesToInsert:          "",
					InsertIdIdx:            10,
				},
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
