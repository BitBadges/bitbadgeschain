package types_test

import (
	math "math"
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
				Creator:       "invalid_address",
				CollectionUri: "https://facebook.com",
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateUris{
				Creator:       sample.AccAddress(),
				CollectionUri: "https://facebook.com",
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
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
