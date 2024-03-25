package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	//Note these are alphanumerically sorted (needed for approvals test)
	alice   = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
	bob     = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	charlie = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
)

func TestRemoveOverlaps(t *testing.T) {
	remaining, _ := types.UniversalRemoveOverlaps(sdk.Context{}, &types.UniversalPermissionDetails{
		BadgeId: &types.UintRange{
			Start: sdkmath.NewUint(5),
			End:   sdkmath.NewUint(5),
		},
		TimelineTime: &types.UintRange{
			Start: sdkmath.NewUint(5),
			End:   sdkmath.NewUint(5),
		},
		TransferTime: &types.UintRange{
			Start: sdkmath.NewUint(5),
			End:   sdkmath.NewUint(5),
		},
		OwnershipTime: &types.UintRange{
			Start: sdkmath.NewUint(5),
			End:   sdkmath.NewUint(5),
		},

		ToList: &types.AddressList{
			Addresses: []string{alice, bob, charlie},
			Whitelist: true,
		},
		FromList: &types.AddressList{
			Addresses: []string{alice, bob, charlie},
			Whitelist: true,
		},
		InitiatedByList: &types.AddressList{
			Addresses: []string{alice, bob, charlie},
			Whitelist: true,
		},

		AmountTrackerIdList: &types.AddressList{},

		ChallengeTrackerIdList: &types.AddressList{},
		ApprovalIdList:         &types.AddressList{},
	}, &types.UniversalPermissionDetails{
		BadgeId: &types.UintRange{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(10),
		},
		TimelineTime: &types.UintRange{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(10),
		},
		TransferTime: &types.UintRange{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(10),
		},
		OwnershipTime: &types.UintRange{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(10),
		},
		ToList: &types.AddressList{
			Addresses: []string{alice, bob, charlie},
			Whitelist: true,
		},
		FromList: &types.AddressList{
			Addresses: []string{alice, bob, charlie},
			Whitelist: true,
		},
		InitiatedByList: &types.AddressList{
			Addresses: []string{alice, bob, charlie},
			Whitelist: true,
		},

		AmountTrackerIdList: &types.AddressList{},

		ChallengeTrackerIdList: &types.AddressList{},
		ApprovalIdList:         &types.AddressList{},
	})
	expected := []*types.UniversalPermissionDetails{
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(4),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList:    &types.AddressList{},
			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(6),
				End:   sdkmath.NewUint(10),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(4),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(6),
				End:   sdkmath.NewUint(10),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(4),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},

		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(6),
				End:   sdkmath.NewUint(10),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(4),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},

		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(6),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
	}

	for idx, r := range remaining {
		t.Log(r)
		t.Log(expected[idx])
	}

	for idx, r := range expected {
		x := remaining[idx]
		require.Equal(t, r.BadgeId.End, x.BadgeId.End)
		require.Equal(t, r.BadgeId.Start, x.BadgeId.Start)

		require.Equal(t, r.TimelineTime.End, x.TimelineTime.End)
		require.Equal(t, r.TimelineTime.Start, x.TimelineTime.Start)

		require.Equal(t, r.TransferTime.End, x.TransferTime.End)
		require.Equal(t, r.TransferTime.Start, x.TransferTime.Start)

		for _, addr := range r.ToList.Addresses {
			require.Contains(t, x.ToList.Addresses, addr)
		}

		for _, addr := range r.FromList.Addresses {
			require.Contains(t, x.FromList.Addresses, addr)
		}

		for _, addr := range r.InitiatedByList.Addresses {
			require.Contains(t, x.InitiatedByList.Addresses, addr)
		}

		require.Equal(t, r.ToList.Whitelist, x.ToList.Whitelist)
		require.Equal(t, r.FromList.Whitelist, x.FromList.Whitelist)
		require.Equal(t, r.InitiatedByList.Whitelist, x.InitiatedByList.Whitelist)
	}

	require.Equal(t, expected, remaining)
}

func TestRemoveAddresses(t *testing.T) {
	remaining, _ := types.UniversalRemoveOverlaps(sdk.Context{}, &types.UniversalPermissionDetails{
		BadgeId: &types.UintRange{
			Start: sdkmath.NewUint(5),
			End:   sdkmath.NewUint(5),
		},
		TimelineTime: &types.UintRange{
			Start: sdkmath.NewUint(5),
			End:   sdkmath.NewUint(5),
		},
		TransferTime: &types.UintRange{
			Start: sdkmath.NewUint(5),
			End:   sdkmath.NewUint(5),
		},
		OwnershipTime: &types.UintRange{
			Start: sdkmath.NewUint(5),
			End:   sdkmath.NewUint(5),
		},
		ToList: &types.AddressList{
			Addresses: []string{alice},
			Whitelist: true,
		},
		FromList: &types.AddressList{
			Addresses: []string{alice},
			Whitelist: true,
		},
		InitiatedByList: &types.AddressList{
			Addresses: []string{alice},
			Whitelist: true,
		},

		AmountTrackerIdList: &types.AddressList{},

		ChallengeTrackerIdList: &types.AddressList{},
		ApprovalIdList:         &types.AddressList{},
	}, &types.UniversalPermissionDetails{
		BadgeId: &types.UintRange{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(10),
		},
		TimelineTime: &types.UintRange{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(10),
		},
		TransferTime: &types.UintRange{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(10),
		},
		OwnershipTime: &types.UintRange{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(10),
		},
		ToList: &types.AddressList{
			Addresses: []string{alice, bob, charlie},
			Whitelist: true,
		},
		FromList: &types.AddressList{
			Addresses: []string{alice, bob, charlie},
			Whitelist: true,
		},
		InitiatedByList: &types.AddressList{
			Addresses: []string{alice, bob, charlie},
			Whitelist: true,
		},

		AmountTrackerIdList: &types.AddressList{},

		ChallengeTrackerIdList: &types.AddressList{},
		ApprovalIdList:         &types.AddressList{},
	})
	expected := []*types.UniversalPermissionDetails{
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(4),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(6),
				End:   sdkmath.NewUint(10),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(4),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(6),
				End:   sdkmath.NewUint(10),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(4),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},

		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(6),
				End:   sdkmath.NewUint(10),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(4),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},

		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(6),
				End:   sdkmath.NewUint(10),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},

		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			ToList: &types.AddressList{
				Addresses: []string{bob, charlie},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{bob, charlie},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{alice, bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
		{
			TimelineTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			BadgeId: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			TransferTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			OwnershipTime: &types.UintRange{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(5),
			},
			ToList: &types.AddressList{
				Addresses: []string{alice},
				Whitelist: true,
			},
			FromList: &types.AddressList{
				Addresses: []string{alice},
				Whitelist: true,
			},
			InitiatedByList: &types.AddressList{
				Addresses: []string{bob, charlie},
				Whitelist: true,
			},

			AmountTrackerIdList: &types.AddressList{},

			ChallengeTrackerIdList: &types.AddressList{},
			ApprovalIdList:         &types.AddressList{},
		},
	}

	// for idx, r := range remaining {
	// 	t.Log(r)
	// 	t.Log(expected[idx])
	// }

	for idx, r := range expected {
		t.Log(idx)

		x := remaining[idx]
		require.Equal(t, r.BadgeId.End, x.BadgeId.End)
		require.Equal(t, r.BadgeId.Start, x.BadgeId.Start)

		require.Equal(t, r.TimelineTime.End, x.TimelineTime.End)
		require.Equal(t, r.TimelineTime.Start, x.TimelineTime.Start)

		require.Equal(t, r.TransferTime.End, x.TransferTime.End)
		require.Equal(t, r.TransferTime.Start, x.TransferTime.Start)

		for _, addr := range r.ToList.Addresses {
			require.Contains(t, x.ToList.Addresses, addr)
		}

		for _, addr := range r.FromList.Addresses {
			require.Contains(t, x.FromList.Addresses, addr)
		}

		for _, addr := range r.InitiatedByList.Addresses {
			require.Contains(t, x.InitiatedByList.Addresses, addr)
		}

		require.Equal(t, r.ToList.Whitelist, x.ToList.Whitelist)
		require.Equal(t, r.FromList.Whitelist, x.FromList.Whitelist)
		require.Equal(t, r.InitiatedByList.Whitelist, x.InitiatedByList.Whitelist)
	}

	require.Equal(t, expected, remaining)
}
