package types_test

import (
	"testing"

	bitbadgesapp "github.com/bitbadges/bitbadgeschain/app"

	"github.com/stretchr/testify/suite"
)

// Bunch of weird config stuff to setup the app. Inherited most from Cosmos SDK tutorials and existing Cosmos SDK modules.
type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) SetupTest() {
	_ = bitbadgesapp.Setup(
		false,
	)
}

func TestBadgesTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
