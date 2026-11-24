package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/sample"

	"github.com/stretchr/testify/require"
)

func TestMsgUpdateCollection_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateCollection
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateCollection{
				Creator:      "invalid_address",
				CollectionId: sdkmath.NewUint(1),
			},
			err: ErrInvalidAddress,
		}, {
			name: "zero collection id",
			msg: MsgUpdateCollection{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.ZeroUint(),
			},
			err: ErrInvalidCollectionID,
		}, {
			name: "valid address",
			msg: MsgUpdateCollection{
				Creator:      sample.AccAddress(),
				CollectionId: sdkmath.NewUint(1),
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
