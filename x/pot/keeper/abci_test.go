package keeper_test

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/pot/keeper"
	"github.com/bitbadges/bitbadgeschain/x/pot/types"
)

// ---------------------------------------------------------------------------
// Mock ValidatorSetKeeper
// ---------------------------------------------------------------------------

type mockValidator struct {
	consAddr     sdk.ConsAddress
	operatorAddr string // account address (bb1...)
	power        int64
	jailed       bool
}

type MockValidatorSetKeeper struct {
	validators     map[string]*mockValidator // key: consAddr hex string
	disableCalls   []sdk.ConsAddress         // track disable calls for assertions
	enableCalls    []sdk.ConsAddress         // track enable calls for assertions
	canSafelyEnable map[string]bool          // consAddr hex -> override (default: true)
}

func NewMockValidatorSetKeeper() *MockValidatorSetKeeper {
	return &MockValidatorSetKeeper{
		validators:      make(map[string]*mockValidator),
		canSafelyEnable: make(map[string]bool),
	}
}

func (m *MockValidatorSetKeeper) AddValidator(consAddr sdk.ConsAddress, operatorAddr string, power int64, jailed bool) {
	m.validators[consAddr.String()] = &mockValidator{
		consAddr:     consAddr,
		operatorAddr: operatorAddr,
		power:        power,
		jailed:       jailed,
	}
}

func (m *MockValidatorSetKeeper) IterateActiveValidators(_ context.Context, fn func(val types.ValidatorInfo) bool) error {
	for _, v := range m.validators {
		if v.jailed || v.power <= 0 {
			continue
		}
		if fn(types.ValidatorInfo{
			ConsAddr:     v.consAddr,
			OperatorAddr: v.operatorAddr,
			Power:        v.power,
		}) {
			break
		}
	}
	return nil
}

func (m *MockValidatorSetKeeper) GetValidatorByConsAddr(_ context.Context, consAddr sdk.ConsAddress) (types.ValidatorInfo, error) {
	v, ok := m.validators[consAddr.String()]
	if !ok {
		return types.ValidatorInfo{}, fmt.Errorf("validator not found: %s", consAddr)
	}
	return types.ValidatorInfo{
		ConsAddr:     v.consAddr,
		OperatorAddr: v.operatorAddr,
		Power:        v.power,
	}, nil
}

func (m *MockValidatorSetKeeper) DisableValidator(_ context.Context, consAddr sdk.ConsAddress) error {
	m.disableCalls = append(m.disableCalls, consAddr)
	v, ok := m.validators[consAddr.String()]
	if !ok {
		return fmt.Errorf("validator not found: %s", consAddr)
	}
	v.jailed = true
	return nil
}

func (m *MockValidatorSetKeeper) EnableValidator(_ context.Context, consAddr sdk.ConsAddress) error {
	m.enableCalls = append(m.enableCalls, consAddr)
	v, ok := m.validators[consAddr.String()]
	if !ok {
		return fmt.Errorf("validator not found: %s", consAddr)
	}
	v.jailed = false
	return nil
}

func (m *MockValidatorSetKeeper) IsValidatorJailed(_ context.Context, consAddr sdk.ConsAddress) (bool, error) {
	v, ok := m.validators[consAddr.String()]
	if !ok {
		return false, fmt.Errorf("validator not found: %s", consAddr)
	}
	return v.jailed, nil
}

func (m *MockValidatorSetKeeper) CanSafelyEnable(_ context.Context, consAddr sdk.ConsAddress) bool {
	if val, ok := m.canSafelyEnable[consAddr.String()]; ok {
		return val
	}
	return true // default: safe to enable
}

// ---------------------------------------------------------------------------
// Mock TokenizationKeeper
// ---------------------------------------------------------------------------

type MockTokenizationKeeper struct {
	// key: "collectionId:tokenId:address" -> balance
	balances map[string]uint64
}

func NewMockTokenizationKeeper() *MockTokenizationKeeper {
	return &MockTokenizationKeeper{
		balances: make(map[string]uint64),
	}
}

func (m *MockTokenizationKeeper) SetBalance(collectionId, tokenId uint64, address string, balance uint64) {
	key := fmt.Sprintf("%d:%d:%s", collectionId, tokenId, address)
	m.balances[key] = balance
}

func (m *MockTokenizationKeeper) GetCredentialBalance(_ sdk.Context, collectionId, tokenId uint64, address string) (uint64, error) {
	key := fmt.Sprintf("%d:%d:%s", collectionId, tokenId, address)
	bal, ok := m.balances[key]
	if !ok {
		return 0, nil
	}
	return bal, nil
}

// ---------------------------------------------------------------------------
// Test helper: create a validator with a random consensus key
// ---------------------------------------------------------------------------

func makeTestValidator(t *testing.T) (sdk.ConsAddress, sdk.AccAddress) {
	t.Helper()
	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()
	consAddr := sdk.ConsAddress(pubKey.Address())
	accAddr := sdk.AccAddress(pubKey.Address())
	return consAddr, accAddr
}

// setupKeeperWithMocks creates a pot Keeper wired to the given mock keepers.
func setupKeeperWithMocks(t *testing.T, vs *MockValidatorSetKeeper, tk *MockTokenizationKeeper) (keeper.Keeper, sdk.Context) {
	t.Helper()

	storeKey := storetypes.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	appCodec := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)

	k := keeper.NewKeeper(
		appCodec,
		storeKey,
		log.NewNopLogger(),
		authority.String(),
		tk,
		vs,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())
	return k, ctx
}

// ---------------------------------------------------------------------------
// EndBlocker tests
// ---------------------------------------------------------------------------

func TestEndBlocker_DisablesValidatorWithoutCredential(t *testing.T) {
	vs := NewMockValidatorSetKeeper()
	tk := NewMockTokenizationKeeper()

	consAddr1, accAddr1 := makeTestValidator(t)
	consAddr2, accAddr2 := makeTestValidator(t)

	vs.AddValidator(consAddr1, accAddr1.String(), 100, false)
	vs.AddValidator(consAddr2, accAddr2.String(), 200, false)

	// Only validator 1 has a credential.
	tk.SetBalance(1, 1, accAddr1.String(), 1)

	k, ctx := setupKeeperWithMocks(t, vs, tk)

	require.NoError(t, k.SetParams(ctx, types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}))

	err := k.EndBlocker(ctx)
	require.NoError(t, err)

	// Validator 2 should have been disabled.
	require.Len(t, vs.disableCalls, 1)
	require.Equal(t, consAddr2, vs.disableCalls[0])

	// Validator 2 should be in compliance-jailed set.
	require.True(t, k.IsComplianceJailed(ctx, consAddr2))
	require.False(t, k.IsComplianceJailed(ctx, consAddr1))
}

func TestEndBlocker_AutoEnablesWhenCredentialRegained(t *testing.T) {
	vs := NewMockValidatorSetKeeper()
	tk := NewMockTokenizationKeeper()

	consAddr1, accAddr1 := makeTestValidator(t)
	consAddr2, accAddr2 := makeTestValidator(t)

	vs.AddValidator(consAddr1, accAddr1.String(), 100, false)
	vs.AddValidator(consAddr2, accAddr2.String(), 200, false)

	tk.SetBalance(1, 1, accAddr1.String(), 1)
	// Validator 2 initially has no credential.

	k, ctx := setupKeeperWithMocks(t, vs, tk)

	require.NoError(t, k.SetParams(ctx, types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}))

	// First EndBlocker: validator 2 gets disabled.
	err := k.EndBlocker(ctx)
	require.NoError(t, err)
	require.Len(t, vs.disableCalls, 1)
	require.True(t, k.IsComplianceJailed(ctx, consAddr2))

	// Give validator 2 a credential.
	tk.SetBalance(1, 1, accAddr2.String(), 1)

	// Second EndBlocker: validator 2 should be auto-enabled.
	vs.disableCalls = nil
	vs.enableCalls = nil
	err = k.EndBlocker(ctx)
	require.NoError(t, err)

	require.Len(t, vs.enableCalls, 1)
	require.Equal(t, consAddr2, vs.enableCalls[0])
	require.False(t, k.IsComplianceJailed(ctx, consAddr2))
}

func TestEndBlocker_DoesNotReDisableAlreadyComplianceJailed(t *testing.T) {
	vs := NewMockValidatorSetKeeper()
	tk := NewMockTokenizationKeeper()

	consAddr1, accAddr1 := makeTestValidator(t)
	consAddr2, accAddr2 := makeTestValidator(t)

	vs.AddValidator(consAddr1, accAddr1.String(), 100, false)
	vs.AddValidator(consAddr2, accAddr2.String(), 200, false)

	tk.SetBalance(1, 1, accAddr1.String(), 1)
	// Validator 2 has no credential.
	_ = accAddr2

	k, ctx := setupKeeperWithMocks(t, vs, tk)

	require.NoError(t, k.SetParams(ctx, types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}))

	// First EndBlocker: validator 2 gets disabled.
	err := k.EndBlocker(ctx)
	require.NoError(t, err)
	require.Len(t, vs.disableCalls, 1)

	// Second EndBlocker: validator 2 is already compliance-jailed, should NOT be disabled again.
	vs.disableCalls = nil
	err = k.EndBlocker(ctx)
	require.NoError(t, err)
	require.Empty(t, vs.disableCalls)
}

func TestEndBlocker_SafetyRefusesToDisableAllValidators(t *testing.T) {
	vs := NewMockValidatorSetKeeper()
	tk := NewMockTokenizationKeeper()

	consAddr1, accAddr1 := makeTestValidator(t)
	vs.AddValidator(consAddr1, accAddr1.String(), 100, false)
	// No credential for the only validator.
	_ = accAddr1

	k, ctx := setupKeeperWithMocks(t, vs, tk)

	require.NoError(t, k.SetParams(ctx, types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}))

	err := k.EndBlocker(ctx)
	require.NoError(t, err)

	// Should NOT have disabled anyone.
	require.Empty(t, vs.disableCalls)
	require.False(t, k.IsComplianceJailed(ctx, consAddr1))
}

func TestEndBlocker_ModuleDisabled_NoAction(t *testing.T) {
	vs := NewMockValidatorSetKeeper()
	tk := NewMockTokenizationKeeper()

	consAddr1, accAddr1 := makeTestValidator(t)
	vs.AddValidator(consAddr1, accAddr1.String(), 100, false)
	// No credential but module disabled.

	k, ctx := setupKeeperWithMocks(t, vs, tk)

	require.NoError(t, k.SetParams(ctx, types.Params{
		CredentialCollectionId: 0, // disabled
		CredentialTokenId:      0,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}))

	err := k.EndBlocker(ctx)
	require.NoError(t, err)
	require.Empty(t, vs.disableCalls)
}

func TestEndBlocker_AllValidatorsHaveCredentials_NoAction(t *testing.T) {
	vs := NewMockValidatorSetKeeper()
	tk := NewMockTokenizationKeeper()

	k, ctx := setupKeeperWithMocks(t, vs, tk)

	require.NoError(t, k.SetParams(ctx, types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}))

	for i := 0; i < 3; i++ {
		consAddr, accAddr := makeTestValidator(t)
		vs.AddValidator(consAddr, accAddr.String(), int64((i+1)*100), false)
		tk.SetBalance(1, 1, accAddr.String(), 1)
	}

	err := k.EndBlocker(ctx)
	require.NoError(t, err)
	require.Empty(t, vs.disableCalls)
	require.Empty(t, vs.enableCalls)
}

func TestEndBlocker_SafetyAllowsPartialDisabling(t *testing.T) {
	vs := NewMockValidatorSetKeeper()
	tk := NewMockTokenizationKeeper()

	k, ctx := setupKeeperWithMocks(t, vs, tk)

	require.NoError(t, k.SetParams(ctx, types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}))

	// 3 validators: only the first has a credential.
	var accAddrs []sdk.AccAddress
	for i := 0; i < 3; i++ {
		consAddr, accAddr := makeTestValidator(t)
		vs.AddValidator(consAddr, accAddr.String(), 100, false)
		accAddrs = append(accAddrs, accAddr)
	}
	tk.SetBalance(1, 1, accAddrs[0].String(), 1)

	err := k.EndBlocker(ctx)
	require.NoError(t, err)
	require.Len(t, vs.disableCalls, 2)
}

func TestEndBlocker_CanSafelyEnableFalse_DefersEnable(t *testing.T) {
	vs := NewMockValidatorSetKeeper()
	tk := NewMockTokenizationKeeper()

	consAddr1, accAddr1 := makeTestValidator(t)
	consAddr2, accAddr2 := makeTestValidator(t)

	vs.AddValidator(consAddr1, accAddr1.String(), 100, false)
	vs.AddValidator(consAddr2, accAddr2.String(), 200, false)

	tk.SetBalance(1, 1, accAddr1.String(), 1)
	// Validator 2 has no credential.

	k, ctx := setupKeeperWithMocks(t, vs, tk)

	require.NoError(t, k.SetParams(ctx, types.Params{
		CredentialCollectionId: 1,
		CredentialTokenId:      1,
		MinCredentialBalance:   1,
		Mode:                   types.ModeStakedMultiplier,
	}))

	// Disable validator 2.
	err := k.EndBlocker(ctx)
	require.NoError(t, err)
	require.True(t, k.IsComplianceJailed(ctx, consAddr2))

	// Give credential back, but mark CanSafelyEnable as false (e.g., slashing jail).
	tk.SetBalance(1, 1, accAddr2.String(), 1)
	vs.canSafelyEnable[consAddr2.String()] = false

	vs.disableCalls = nil
	vs.enableCalls = nil
	err = k.EndBlocker(ctx)
	require.NoError(t, err)

	// Should NOT have enabled — deferred.
	require.Empty(t, vs.enableCalls)
	// Should still be compliance-jailed.
	require.True(t, k.IsComplianceJailed(ctx, consAddr2))
}

// ---------------------------------------------------------------------------
// Compliance-jailed state tests
// ---------------------------------------------------------------------------

func TestComplianceJailedState(t *testing.T) {
	vs := NewMockValidatorSetKeeper()
	tk := NewMockTokenizationKeeper()
	k, ctx := setupKeeperWithMocks(t, vs, tk)

	consAddr := sdk.ConsAddress(ed25519.GenPrivKey().PubKey().Address())

	// Initially not compliance-jailed.
	require.False(t, k.IsComplianceJailed(ctx, consAddr))

	// Set compliance-jailed.
	k.SetComplianceJailed(ctx, consAddr)
	require.True(t, k.IsComplianceJailed(ctx, consAddr))

	// GetAll should include it.
	all := k.GetAllComplianceJailed(ctx)
	require.Len(t, all, 1)

	// Remove compliance-jailed.
	k.RemoveComplianceJailed(ctx, consAddr)
	require.False(t, k.IsComplianceJailed(ctx, consAddr))

	all = k.GetAllComplianceJailed(ctx)
	require.Empty(t, all)
}

// ---------------------------------------------------------------------------
// SavedPower (PowerStore) tests
// ---------------------------------------------------------------------------

func TestSavedPower(t *testing.T) {
	vs := NewMockValidatorSetKeeper()
	tk := NewMockTokenizationKeeper()
	k, ctx := setupKeeperWithMocks(t, vs, tk)

	consAddr := sdk.ConsAddress(ed25519.GenPrivKey().PubKey().Address())

	// Initially no saved power.
	power, found := k.GetSavedPower(ctx, consAddr)
	require.False(t, found)
	require.Equal(t, int64(0), power)

	// Set saved power.
	k.SetSavedPower(ctx, consAddr, 42)
	power, found = k.GetSavedPower(ctx, consAddr)
	require.True(t, found)
	require.Equal(t, int64(42), power)

	// Remove saved power.
	k.RemoveSavedPower(ctx, consAddr)
	power, found = k.GetSavedPower(ctx, consAddr)
	require.False(t, found)
	require.Equal(t, int64(0), power)
}
