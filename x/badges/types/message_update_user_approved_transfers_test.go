package types_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/sample"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateUserApprovals_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateUserApprovals
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateUserApprovals{
				Creator: "invalid_address",
			},
			err: types.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgUpdateUserApprovals{
				Creator: sample.AccAddress(),
			},
		},
		{
			name: "tracker ID = ID of another approval",
			msg: types.MsgUpdateUserApprovals{
				Creator: sample.AccAddress(),
				OutgoingApprovals: []*types.UserOutgoingApproval{
					{
						ToListId: "All",
						InitiatedByListId: "All",
						ApprovalId: "approval_id",
						AmountTrackerId: "other_approval_id",
						ChallengeTrackerId: "challenge_id",
					},
					{
						ToListId: "All",
						InitiatedByListId: "All",
						ApprovalId: "other_approval_id",
						AmountTrackerId: "other_approval_id",
						ChallengeTrackerId: "challenge_id",
					},
				},
			},
			err: types.ErrAmountTrackerIdIsNil,
		},
		{
			name: "tracker ID = ID of another approval",
			msg: types.MsgUpdateUserApprovals{
				Creator: sample.AccAddress(),
				OutgoingApprovals: []*types.UserOutgoingApproval{
					{
						ToListId: "All",
						InitiatedByListId: "All",
						ApprovalId: "approval_id",
						AmountTrackerId: "other_approval_id",
						ChallengeTrackerId: "challenge_id",
					},
					{
						ToListId: "All",
						InitiatedByListId: "All",
						ApprovalId: "dfsdgffds_approval_id",
						AmountTrackerId: "other_approval_id",
						ChallengeTrackerId: "challenge_id",
					},
				},
			},
		},
		{
			name: "tracker ID = ID of another approval",
			msg: types.MsgUpdateUserApprovals{
				Creator: sample.AccAddress(),
				OutgoingApprovals: []*types.UserOutgoingApproval{
					{
						ToListId: "All",
						InitiatedByListId: "All",
						ApprovalId: "approval_id",
						AmountTrackerId: "sdafaf",
						ChallengeTrackerId: "other_approval_id",
					},
					{
						ToListId: "All",
						InitiatedByListId: "All",
						ApprovalId: "other_approval_id",
						AmountTrackerId: "afdsa",
						ChallengeTrackerId: "asfdadsf",
					},
				},
			},
			err: types.ErrAmountTrackerIdIsNil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.Error(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
