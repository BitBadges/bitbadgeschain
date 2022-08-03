package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/trevormil/bitbadgeschain/testutil/sample"
)

func TestMsgUpdateUris_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateUris
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateUris{
				Creator: "invalid_address",
				Uri:    "https://example.com",
				SubassetUri: "https://example.com",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateUris{
				Creator: sample.AccAddress(),
				Uri:    "https://example.com",
				SubassetUri: "https://example.com",
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
