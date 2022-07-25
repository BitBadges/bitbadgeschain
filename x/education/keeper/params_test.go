package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/trevormil/bitbadgeschain/testutil/keeper"
	"github.com/trevormil/bitbadgeschain/x/education/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.EducationKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
