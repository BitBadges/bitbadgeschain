package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"
)

const (
	//Note these are alphanumerically sorted (needed for approvals test)
	alice   = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
	bob     = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	charlie = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
)

func TestRemoveOverlaps(t *testing.T) {
	remaining, _ := types.UniversalRemoveOverlaps(&types.UniversalPermissionDetails{
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
		ToMapping: &types.AddressMapping{
			Addresses:        []string{alice, bob, charlie},
			IncludeAddresses: true,
		},
		FromMapping: &types.AddressMapping{
			Addresses:        []string{alice, bob, charlie},
			IncludeAddresses: true,
		},
		InitiatedByMapping: &types.AddressMapping{
			Addresses:        []string{alice, bob, charlie},
			IncludeAddresses: true,
		},
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
		ToMapping: &types.AddressMapping{
			Addresses:        []string{alice, bob, charlie},
			IncludeAddresses: true,
		},
		FromMapping: &types.AddressMapping{
			Addresses:        []string{alice, bob, charlie},
			IncludeAddresses: true,
		},
		InitiatedByMapping: &types.AddressMapping{
			Addresses:        []string{alice, bob, charlie},
			IncludeAddresses: true,
		},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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

		for _, addr := range r.ToMapping.Addresses {
			require.Contains(t, x.ToMapping.Addresses, addr)
		}

		for _, addr := range r.FromMapping.Addresses {
			require.Contains(t, x.FromMapping.Addresses, addr)
		}

		for _, addr := range r.InitiatedByMapping.Addresses {
			require.Contains(t, x.InitiatedByMapping.Addresses, addr)
		}

		require.Equal(t, r.ToMapping.IncludeAddresses, x.ToMapping.IncludeAddresses)
		require.Equal(t, r.FromMapping.IncludeAddresses, x.FromMapping.IncludeAddresses)
		require.Equal(t, r.InitiatedByMapping.IncludeAddresses, x.InitiatedByMapping.IncludeAddresses)
	}

	require.Equal(t, expected, remaining)
}

func TestRemoveAddresses(t *testing.T) {
	remaining, _ := types.UniversalRemoveOverlaps(&types.UniversalPermissionDetails{
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
		ToMapping: &types.AddressMapping{
			Addresses:        []string{alice},
			IncludeAddresses: true,
		},
		FromMapping: &types.AddressMapping{
			Addresses:        []string{alice},
			IncludeAddresses: true,
		},
		InitiatedByMapping: &types.AddressMapping{
			Addresses:        []string{alice},
			IncludeAddresses: true,
		},
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
		ToMapping: &types.AddressMapping{
			Addresses:        []string{alice, bob, charlie},
			IncludeAddresses: true,
		},
		FromMapping: &types.AddressMapping{
			Addresses:        []string{alice, bob, charlie},
			IncludeAddresses: true,
		},
		InitiatedByMapping: &types.AddressMapping{
			Addresses:        []string{alice, bob, charlie},
			IncludeAddresses: true,
		},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{alice, bob, charlie},
				IncludeAddresses: true,
			},
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
			ToMapping: &types.AddressMapping{
				Addresses:        []string{bob, charlie},
				IncludeAddresses: true,
			},
			FromMapping: &types.AddressMapping{
				Addresses:        []string{bob, charlie},
				IncludeAddresses: true,
			},
			InitiatedByMapping: &types.AddressMapping{
				Addresses:        []string{bob, charlie},
				IncludeAddresses: true,
			},
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

		for _, addr := range r.ToMapping.Addresses {
			require.Contains(t, x.ToMapping.Addresses, addr)
		}

		for _, addr := range r.FromMapping.Addresses {
			require.Contains(t, x.FromMapping.Addresses, addr)
		}

		for _, addr := range r.InitiatedByMapping.Addresses {
			require.Contains(t, x.InitiatedByMapping.Addresses, addr)
		}

		require.Equal(t, r.ToMapping.IncludeAddresses, x.ToMapping.IncludeAddresses)
		require.Equal(t, r.FromMapping.IncludeAddresses, x.FromMapping.IncludeAddresses)
		require.Equal(t, r.InitiatedByMapping.IncludeAddresses, x.InitiatedByMapping.IncludeAddresses)
	}

	require.Equal(t, expected, remaining)
}
