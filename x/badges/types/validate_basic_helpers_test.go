package types_test

import (
	math "math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
)

func GetValidUserPermissions() *types.UserPermissions {
	return &types.UserPermissions{
		CanUpdateOutgoingApprovals: []*types.UserOutgoingApprovalPermission{
			{
				PermanentlyPermittedTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
				PermanentlyForbiddenTimes: []*types.UintRange{{Start: sdkmath.NewUint(6), End: sdkmath.NewUint(8)}},
			},
		},
	}
}

func GetValidCollectionMetadataTimeline() []*types.CollectionMetadataTimeline {
	return []*types.CollectionMetadataTimeline{
		{
			CollectionMetadata: &types.CollectionMetadata{
				Uri: "https://example.com/{id}",
			},
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
			},
		},
	}
}

func GetValidTokenMetadataTimeline() []*types.TokenMetadataTimeline {
	return []*types.TokenMetadataTimeline{
		{
			TokenMetadata: []*types.TokenMetadata{
				{
					Uri: "https://example.com/{id}",
					TokenIds: []*types.UintRange{
						{
							Start: sdkmath.NewUint(1),

							End: sdkmath.NewUint(math.MaxUint64),
						},
					},
				},
			},
			TimelineTimes: []*types.UintRange{
				{
					Start: sdkmath.NewUint(0),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
			},
		},
	}
}
