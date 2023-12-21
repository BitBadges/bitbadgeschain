package types

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
)

func TestMsgSetCollectionForProtocol_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSetCollectionForProtocol
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSetCollectionForProtocol{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetCollectionForProtocol{
				Creator: sample.AccAddress(),
				Name: "non-empty",
				CollectionId: sdkmath.NewUintFromString("1"),
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
