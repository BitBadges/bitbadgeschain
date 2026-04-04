package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/pot/types"
)

func TestDefaultParams(t *testing.T) {
	p := types.DefaultParams()
	require.Equal(t, uint64(0), p.CredentialCollectionId)
	require.Equal(t, uint64(0), p.CredentialTokenId)
	require.Equal(t, uint64(1), p.MinCredentialBalance)
	require.Equal(t, types.ModeStakedMultiplier, p.Mode)
}

func TestParamsValidate_StakedMultiplier(t *testing.T) {
	p := types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}
	require.NoError(t, p.Validate())
}

func TestParamsValidate_Equal_Rejected(t *testing.T) {
	p := types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeEqual,
	}
	err := p.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "only \"staked_multiplier\" is currently supported")
}

func TestParamsValidate_CredentialWeighted_Rejected(t *testing.T) {
	p := types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeCredentialWeighted,
	}
	err := p.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "only \"staked_multiplier\" is currently supported")
}

func TestParamsValidate_InvalidMode(t *testing.T) {
	p := types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   "invalid_mode",
	}
	err := p.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid mode")
}

func TestParamsValidate_MinCredentialBalanceZero(t *testing.T) {
	p := types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   0,
		Mode:                   types.ModeStakedMultiplier,
	}
	err := p.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "min_credential_balance must be > 0")
}

func TestIsEnabled_CollectionIdZero(t *testing.T) {
	p := types.Params{
		CredentialCollectionId: 0,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}
	require.False(t, p.IsEnabled())
}

func TestIsEnabled_CollectionIdPositive(t *testing.T) {
	p := types.Params{
		CredentialCollectionId: 42,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}
	require.True(t, p.IsEnabled())
}

func TestParamsValidate_DefaultParamsAreValid(t *testing.T) {
	// DefaultParams has MinCredentialBalance=1, Mode=staked_multiplier, so it should be valid.
	p := types.DefaultParams()
	require.NoError(t, p.Validate())
}

func TestParamsValidate_EmptyMode(t *testing.T) {
	p := types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   "",
	}
	err := p.Validate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid mode")
}
