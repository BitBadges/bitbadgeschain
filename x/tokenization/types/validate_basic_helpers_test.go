package types_test

import (
	math "math"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

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

func GetValidCollectionMetadata() *types.CollectionMetadata {
	return &types.CollectionMetadata{
		Uri: "https://example.com/{id}",
	}
}

func GetValidTokenMetadata() []*types.TokenMetadata {
	return []*types.TokenMetadata{
		{
			Uri: "https://example.com/{id}",
			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
			},
		},
	}
}
