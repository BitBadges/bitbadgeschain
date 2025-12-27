package msg_handlers

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

type UpdateParamsTestSuite struct {
	testutil.AITestSuite
}

func TestUpdateParamsTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateParamsTestSuite))
}

func (suite *UpdateParamsTestSuite) TestUpdateParams_ValidMessage() {
	msg := &types.MsgUpdateParams{
		Authority: suite.Authority,
		Params:   types.DefaultParams(),
	}

	resp, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(resp)
}

func (suite *UpdateParamsTestSuite) TestUpdateParams_InvalidAuthority() {
	msg := &types.MsgUpdateParams{
		Authority: suite.Alice, // Not the authority
		Params:   types.DefaultParams(),
	}

	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
	// The error could be either "invalid authority address" (if address parsing fails) or "invalid authority" (if authority check fails)
	// Since Alice is a valid address but not the authority, it should fail the authority check
	errorMsg := err.Error()
	suite.Require().True(
		strings.Contains(errorMsg, "invalid authority") || strings.Contains(errorMsg, "invalid authority address"),
		"error should contain 'invalid authority' or 'invalid authority address', got: %s", errorMsg,
	)
}

func (suite *UpdateParamsTestSuite) TestUpdateParams_EmptyAuthority() {
	msg := &types.MsgUpdateParams{
		Authority: "",
		Params:   types.DefaultParams(),
	}

	_, err := suite.MsgServer.UpdateParams(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err)
}

