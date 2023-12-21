package types

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateProtocol_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateProtocol
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateProtocol{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateProtocol{
				Creator: sample.AccAddress(),
				Name :"hjdsafkjal",
			},
		}, {
			name: "empty name",
			msg: MsgUpdateProtocol{
				Creator: sample.AccAddress(),
				Name :"",
			},
			err: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty"),
		}, {
			name: "invalid uri",
			msg: MsgUpdateProtocol{
				Creator: sample.AccAddress(),
				Name :"hjdsafkjal",
				Uri: "invalid_uri",
			},
			err: sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "uri cannot be invalid"),
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
