package attack_scenarios

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
	sendmanagerkeeper "github.com/bitbadges/bitbadgeschain/x/sendmanager/keeper"
)

// PrefixCollisionTestSuite tests prefix collision attack scenarios
type PrefixCollisionTestSuite struct {
	testutil.AITestSuite
}

func TestPrefixCollisionTestSuite(t *testing.T) {
	suite.Run(t, new(PrefixCollisionTestSuite))
}

// TestPrefixCollision_NonBadgeslpPrefixRejected tests that non-badgeslp prefixes are rejected
func (s *PrefixCollisionTestSuite) TestPrefixCollision_NonBadgeslpPrefixRejected() {
	router1 := testutil.GenerateMockRouter("a:")

	// Any prefix that is not "badgeslp:" should be rejected
	err := s.Keeper.RegisterRouter("a:", router1)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "only prefix")
}

// TestPrefixCollision_AnotherNonBadgeslpPrefixRejected tests that another custom prefix is rejected
func (s *PrefixCollisionTestSuite) TestPrefixCollision_AnotherNonBadgeslpPrefixRejected() {
	router := testutil.GenerateMockRouter("a:b:")

	err := s.Keeper.RegisterRouter("a:b:", router)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "only prefix")
}

// TestPrefixCollision_EmptyPrefix tests empty prefix handling
func (s *PrefixCollisionTestSuite) TestPrefixCollision_EmptyPrefix() {
	router := testutil.GenerateMockRouter("")

	err := s.Keeper.RegisterRouter("", router)
	s.Require().Error(err, "Empty prefix should be rejected")
	s.Require().Contains(err.Error(), "only prefix")
}

// TestPrefixCollision_DuplicateBadgeslpPrefix tests that re-registering badgeslp: overwrites successfully
func (s *PrefixCollisionTestSuite) TestPrefixCollision_DuplicateBadgeslpPrefix() {
	router1 := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)
	router2 := testutil.GenerateMockRouter(sendmanagerkeeper.AliasDenomPrefix)

	err := s.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router1)
	s.Require().NoError(err)

	// Re-registering badgeslp: should succeed (overwrites)
	err = s.Keeper.RegisterRouter(sendmanagerkeeper.AliasDenomPrefix, router2)
	s.Require().NoError(err)
}

// TestPrefixCollision_EmptyDenomRouting tests routing with empty denom
func (s *PrefixCollisionTestSuite) TestPrefixCollision_EmptyDenomRouting() {
	// Test that empty denom is rejected (validation prevents empty denoms)
	denom := ""

	testAddr := sdk.AccAddress("test-address-123456")

	// Empty denom should be rejected by validation
	_, err := s.Keeper.GetBalanceWithAliasRouting(s.Ctx, testAddr, denom)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "cannot be empty")
}
