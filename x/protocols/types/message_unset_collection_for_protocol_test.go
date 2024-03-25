package types

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUnsetCollectionForProtocol_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUnsetCollectionForProtocol
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUnsetCollectionForProtocol{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUnsetCollectionForProtocol{
				Creator: sample.AccAddress(),
				Name:    "fdjhsajhksfd",
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
