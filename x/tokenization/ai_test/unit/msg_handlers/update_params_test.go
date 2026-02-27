package msg_handlers_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type UpdateParamsTestSuite struct {
	testutil.AITestSuite
}

func TestUpdateParamsSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(UpdateParamsTestSuite))
}

// TestUpdateParams_GovernanceAuthorityRequired tests that governance authority is required
func (suite *UpdateParamsTestSuite) TestUpdateParams_GovernanceAuthorityRequired() {
	authority := suite.Keeper.GetAuthority()

	params := types.Params{
		AllowedDenoms:       []string{"ubadge", "stake"},
		AffiliatePercentage: sdkmath.NewUint(10),
	}

	msg := &types.MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}

	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "governance authority should be able to update params")

	// Verify params were updated
	updatedParams := suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Equal([]string{"ubadge", "stake"}, updatedParams.AllowedDenoms)
	suite.Require().Equal(sdkmath.NewUint(10), updatedParams.AffiliatePercentage)
}

// TestUpdateParams_NonGovernanceRejected tests that non-governance authority is rejected
func (suite *UpdateParamsTestSuite) TestUpdateParams_NonGovernanceRejected() {
	params := types.Params{
		AllowedDenoms:       []string{"ubadge"},
		AffiliatePercentage: sdkmath.NewUint(5),
	}

	// Try using Alice (not governance authority)
	msg := &types.MsgUpdateParams{
		Authority: suite.Alice,
		Params:    params,
	}

	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-governance authority should not be able to update params")
	suite.Require().Contains(err.Error(), "invalid authority", "error should mention invalid authority")
}

// TestUpdateParams_AllowedDenomsUpdate tests updating allowed_denoms
func (suite *UpdateParamsTestSuite) TestUpdateParams_AllowedDenomsUpdate() {
	authority := suite.Keeper.GetAuthority()

	// Set initial params with some denoms
	initialParams := types.Params{
		AllowedDenoms:       []string{"ubadge"},
		AffiliatePercentage: sdkmath.NewUint(0),
	}
	initialMsg := &types.MsgUpdateParams{
		Authority: authority,
		Params:    initialParams,
	}
	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), initialMsg)
	suite.Require().NoError(err)

	// Verify initial params
	params := suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Equal([]string{"ubadge"}, params.AllowedDenoms)

	// Update to different denoms
	updatedParams := types.Params{
		AllowedDenoms:       []string{"ubadge", "uatom", "uosmo"},
		AffiliatePercentage: sdkmath.NewUint(0),
	}
	updateMsg := &types.MsgUpdateParams{
		Authority: authority,
		Params:    updatedParams,
	}
	_, err = suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "updating allowed denoms should succeed")

	// Verify updated params
	params = suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Equal([]string{"ubadge", "uatom", "uosmo"}, params.AllowedDenoms)
}

// TestUpdateParams_EmptyAllowedDenoms tests setting empty allowed_denoms
func (suite *UpdateParamsTestSuite) TestUpdateParams_EmptyAllowedDenoms() {
	authority := suite.Keeper.GetAuthority()

	params := types.Params{
		AllowedDenoms:       []string{}, // Empty
		AffiliatePercentage: sdkmath.NewUint(0),
	}

	msg := &types.MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}

	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting empty allowed denoms should succeed")

	// Verify params - use Len check since empty slice and nil are equivalent
	updatedParams := suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Len(updatedParams.AllowedDenoms, 0, "allowed denoms should be empty")
}

// TestUpdateParams_ParamsPersistedCorrectly tests that params are persisted correctly
func (suite *UpdateParamsTestSuite) TestUpdateParams_ParamsPersistedCorrectly() {
	authority := suite.Keeper.GetAuthority()

	// Set params
	params := types.Params{
		AllowedDenoms:       []string{"ubadge", "ustake", "utest"},
		AffiliatePercentage: sdkmath.NewUint(25),
	}

	msg := &types.MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}

	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Verify all params are persisted correctly
	updatedParams := suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Equal([]string{"ubadge", "ustake", "utest"}, updatedParams.AllowedDenoms, "allowed denoms should be persisted")
	suite.Require().Equal(sdkmath.NewUint(25), updatedParams.AffiliatePercentage, "affiliate percentage should be persisted")
}

// TestUpdateParams_InvalidAuthority tests various invalid authority scenarios
func (suite *UpdateParamsTestSuite) TestUpdateParams_InvalidAuthority() {
	params := types.DefaultParams()

	invalidAuthorities := []string{
		"",                      // Empty
		"invalid",               // Invalid bech32
		suite.Bob,               // Valid address but not authority
		suite.Manager,           // Manager is not governance
		"bb1invalidaddress1234", // Invalid checksum
	}

	for _, invalidAuth := range invalidAuthorities {
		msg := &types.MsgUpdateParams{
			Authority: invalidAuth,
			Params:    params,
		}

		_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().Error(err, "invalid authority %s should be rejected", invalidAuth)
	}
}

// TestUpdateParams_MultipleDenomTypes tests various denom configurations
func (suite *UpdateParamsTestSuite) TestUpdateParams_MultipleDenomTypes() {
	authority := suite.Keeper.GetAuthority()

	testCases := []struct {
		name   string
		denoms []string
	}{
		{
			name:   "single denom",
			denoms: []string{"ubadge"},
		},
		{
			name:   "multiple denoms",
			denoms: []string{"ubadge", "uatom", "uosmo", "ustake"},
		},
		{
			name:   "ibc denom",
			denoms: []string{"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"},
		},
		{
			name:   "mixed denoms",
			denoms: []string{"ubadge", "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			params := types.Params{
				AllowedDenoms:       tc.denoms,
				AffiliatePercentage: sdkmath.NewUint(0),
			}

			msg := &types.MsgUpdateParams{
				Authority: authority,
				Params:    params,
			}

			_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
			suite.Require().NoError(err, "setting denoms %v should succeed", tc.denoms)

			// Verify
			updatedParams := suite.Keeper.GetParams(suite.Ctx)
			suite.Require().Equal(tc.denoms, updatedParams.AllowedDenoms)
		})
	}
}

// TestUpdateParams_AffiliatePercentage tests various affiliate percentage values
func (suite *UpdateParamsTestSuite) TestUpdateParams_AffiliatePercentage() {
	authority := suite.Keeper.GetAuthority()

	testCases := []struct {
		name       string
		percentage sdkmath.Uint
	}{
		{
			name:       "zero percentage",
			percentage: sdkmath.NewUint(0),
		},
		{
			name:       "small percentage",
			percentage: sdkmath.NewUint(5),
		},
		{
			name:       "medium percentage",
			percentage: sdkmath.NewUint(50),
		},
		{
			name:       "high percentage",
			percentage: sdkmath.NewUint(100),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			params := types.Params{
				AllowedDenoms:       []string{"ubadge"},
				AffiliatePercentage: tc.percentage,
			}

			msg := &types.MsgUpdateParams{
				Authority: authority,
				Params:    params,
			}

			_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
			suite.Require().NoError(err, "setting affiliate percentage %s should succeed", tc.percentage.String())

			// Verify
			updatedParams := suite.Keeper.GetParams(suite.Ctx)
			suite.Require().Equal(tc.percentage, updatedParams.AffiliatePercentage)
		})
	}
}

// TestUpdateParams_SequentialUpdates tests multiple sequential param updates
func (suite *UpdateParamsTestSuite) TestUpdateParams_SequentialUpdates() {
	authority := suite.Keeper.GetAuthority()

	// First update
	params1 := types.Params{
		AllowedDenoms:       []string{"ubadge"},
		AffiliatePercentage: sdkmath.NewUint(10),
	}
	msg1 := &types.MsgUpdateParams{
		Authority: authority,
		Params:    params1,
	}
	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Verify first update
	updatedParams := suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Equal([]string{"ubadge"}, updatedParams.AllowedDenoms)
	suite.Require().Equal(sdkmath.NewUint(10), updatedParams.AffiliatePercentage)

	// Second update
	params2 := types.Params{
		AllowedDenoms:       []string{"ubadge", "uatom"},
		AffiliatePercentage: sdkmath.NewUint(20),
	}
	msg2 := &types.MsgUpdateParams{
		Authority: authority,
		Params:    params2,
	}
	_, err = suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Verify second update
	updatedParams = suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Equal([]string{"ubadge", "uatom"}, updatedParams.AllowedDenoms)
	suite.Require().Equal(sdkmath.NewUint(20), updatedParams.AffiliatePercentage)

	// Third update - completely different
	params3 := types.Params{
		AllowedDenoms:       []string{"ustake"},
		AffiliatePercentage: sdkmath.NewUint(5),
	}
	msg3 := &types.MsgUpdateParams{
		Authority: authority,
		Params:    params3,
	}
	_, err = suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().NoError(err)

	// Verify third update - previous values should be replaced
	updatedParams = suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Equal([]string{"ustake"}, updatedParams.AllowedDenoms)
	suite.Require().Equal(sdkmath.NewUint(5), updatedParams.AffiliatePercentage)
}

// TestUpdateParams_DefaultParams tests updating with default params
func (suite *UpdateParamsTestSuite) TestUpdateParams_DefaultParams() {
	authority := suite.Keeper.GetAuthority()

	defaultParams := types.DefaultParams()

	msg := &types.MsgUpdateParams{
		Authority: authority,
		Params:    defaultParams,
	}

	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating with default params should succeed")

	// Verify params match default
	updatedParams := suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Equal(defaultParams.AllowedDenoms, updatedParams.AllowedDenoms)
	suite.Require().Equal(defaultParams.AffiliatePercentage, updatedParams.AffiliatePercentage)
}

// TestUpdateParams_PartialUpdate tests that all fields must be provided
func (suite *UpdateParamsTestSuite) TestUpdateParams_PartialUpdate() {
	authority := suite.Keeper.GetAuthority()

	// First set known values
	initialParams := types.Params{
		AllowedDenoms:       []string{"ubadge", "uatom"},
		AffiliatePercentage: sdkmath.NewUint(50),
	}
	initialMsg := &types.MsgUpdateParams{
		Authority: authority,
		Params:    initialParams,
	}
	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), initialMsg)
	suite.Require().NoError(err)

	// Update with different allowed denoms but same affiliate percentage
	updateParams := types.Params{
		AllowedDenoms:       []string{"ustake"}, // Changed
		AffiliatePercentage: sdkmath.NewUint(50), // Same
	}
	updateMsg := &types.MsgUpdateParams{
		Authority: authority,
		Params:    updateParams,
	}
	_, err = suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Verify update
	updatedParams := suite.Keeper.GetParams(suite.Ctx)
	suite.Require().Equal([]string{"ustake"}, updatedParams.AllowedDenoms, "allowed denoms should be updated")
	suite.Require().Equal(sdkmath.NewUint(50), updatedParams.AffiliatePercentage, "affiliate percentage should remain same")
}

// TestUpdateParams_ZeroAffiliatePercentage tests zero affiliate percentage is valid
func (suite *UpdateParamsTestSuite) TestUpdateParams_ZeroAffiliatePercentage() {
	authority := suite.Keeper.GetAuthority()

	params := types.Params{
		AllowedDenoms:       []string{"ubadge"},
		AffiliatePercentage: sdkmath.NewUint(0),
	}

	msg := &types.MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}

	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "zero affiliate percentage should be valid")

	updatedParams := suite.Keeper.GetParams(suite.Ctx)
	suite.Require().True(updatedParams.AffiliatePercentage.IsZero(), "affiliate percentage should be zero")
}
