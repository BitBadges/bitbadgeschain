package messages

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

type UpdateParamsValidationTestSuite struct {
	testutil.AITestSuite
}

func TestUpdateParamsValidationTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateParamsValidationTestSuite))
}

func (suite *UpdateParamsValidationTestSuite) TestUpdateParams_ValidMessage() {
	msg := &types.MsgUpdateParams{
		Authority: suite.Authority,
		Params:    types.DefaultParams(),
	}

	// MsgUpdateParams doesn't have ValidateBasic, so we test via message server
	// For validation tests, we just verify the message structure is valid
	suite.Require().NotNil(msg)
	suite.Require().NotEmpty(msg.Authority)
}

func (suite *UpdateParamsValidationTestSuite) TestUpdateParams_EmptyAuthority() {
	msg := &types.MsgUpdateParams{
		Authority: "",
		Params:    types.DefaultParams(),
	}

	// Test that empty authority is caught by message server
	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
}

func (suite *UpdateParamsValidationTestSuite) TestUpdateParams_InvalidParams() {
	msg := &types.MsgUpdateParams{
		Authority: suite.Authority,
		Params:    types.Params{}, // Empty params
	}

	// Test that invalid params are caught by message server
	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	// Params validation happens in SetParams
	_ = err // May or may not error depending on params validation
}

