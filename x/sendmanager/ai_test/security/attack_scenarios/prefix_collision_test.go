package attack_scenarios

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/ai_test/testutil"
)

// PrefixCollisionTestSuite tests prefix collision attack scenarios
type PrefixCollisionTestSuite struct {
	testutil.AITestSuite
}

func TestPrefixCollisionTestSuite(t *testing.T) {
	suite.Run(t, new(PrefixCollisionTestSuite))
}

// TestPrefixCollision_OverlappingPrefixes tests that overlapping prefixes are detected
func (s *PrefixCollisionTestSuite) TestPrefixCollision_OverlappingPrefixes() {
	// Create mock routers
	router1 := testutil.GenerateMockRouter("badges:")
	router2 := testutil.GenerateMockRouter("badgeslp:")

	// Register first prefix
	err := s.Keeper.RegisterRouter("badges:", router1)
	s.Require().NoError(err)

	// Attempt to register overlapping prefix - should fail
	// "badgeslp:" starts with "badges:", so they overlap
	err = s.Keeper.RegisterRouter("badgeslp:", router2)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "overlaps")
}

// TestPrefixCollision_SubPrefixRegistration tests sub-prefix registration
func (s *PrefixCollisionTestSuite) TestPrefixCollision_SubPrefixRegistration() {
	router1 := testutil.GenerateMockRouter("badgeslp:")
	router2 := testutil.GenerateMockRouter("badges:")

	// Register longer prefix first
	err := s.Keeper.RegisterRouter("badgeslp:", router1)
	s.Require().NoError(err)

	// Register shorter prefix (sub-prefix) - should be prevented
	// "badgeslp:" starts with "badges:", so they overlap
	err = s.Keeper.RegisterRouter("badges:", router2)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "overlaps")
}

// TestPrefixCollision_EmptyPrefix tests empty prefix handling
func (s *PrefixCollisionTestSuite) TestPrefixCollision_EmptyPrefix() {
	router := testutil.GenerateMockRouter("")
	
	err := s.Keeper.RegisterRouter("", router)
	s.Require().Error(err, "Empty prefix should be rejected")
	s.Require().Contains(err.Error(), "cannot be empty")
}

// TestPrefixCollision_DuplicatePrefix tests duplicate prefix registration
func (s *PrefixCollisionTestSuite) TestPrefixCollision_DuplicatePrefix() {
	router1 := testutil.GenerateMockRouter("badges:")
	router2 := testutil.GenerateMockRouter("badges:")

	err := s.Keeper.RegisterRouter("badges:", router1)
	s.Require().NoError(err)

	// Attempt to register same prefix again
	err = s.Keeper.RegisterRouter("badges:", router2)
	s.Require().Error(err, "Duplicate prefix should be rejected")
	s.Require().Contains(err.Error(), "already registered")
}

// TestPrefixCollision_EmptyDenomRouting tests routing with empty denom
func (s *PrefixCollisionTestSuite) TestPrefixCollision_EmptyDenomRouting() {
	// Test that empty denom is rejected (validation prevents empty denoms)
	denom := ""
	
	// Create a test address directly instead of parsing Bech32 to avoid prefix issues
	// Use a simple address for testing
	testAddr := sdk.AccAddress("test-address-123456")
	
	// Empty denom should be rejected by validation
	_, err := s.Keeper.GetBalanceWithAliasRouting(s.Ctx, testAddr, denom)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "cannot be empty")
}

