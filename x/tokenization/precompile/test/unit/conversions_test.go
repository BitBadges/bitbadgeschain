package tokenization_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ConversionsTestSuite is a placeholder for conversion tests
// Note: Conversion functions have been removed as part of the JSON-based API migration
// These tests are kept for reference but are skipped
type ConversionsTestSuite struct {
	suite.Suite
}

func TestConversionsTestSuite(t *testing.T) {
	suite.Run(t, new(ConversionsTestSuite))
}

// All conversion tests are skipped - conversion functions removed in favor of JSON unmarshaling
// The precompile now accepts JSON strings directly, so conversion from EVM structs is no longer needed
func (suite *ConversionsTestSuite) TestConversions_Skipped() {
	suite.T().Skip("Conversion functions removed - using JSON unmarshaling instead. " +
		"All precompile methods now accept JSON strings matching the Cosmos SDK message structure.")
}
