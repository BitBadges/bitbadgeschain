package gamm

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
)

// A nil *QueryPoolResponse handed to an interface is not an interface nil -- it
// is (type=*QueryPoolResponse, value=nil). Every type assertion in
// packQueryResponse therefore succeeds on it, and the branch then either
// dereferences the nil pointer or marshals it.
//
// This is the same defect that panicked mainnet in compareApprovalCriteria
// (PR #112). None of the gamm queriers returns (nil, nil) today, so it is
// latent rather than live -- it goes live the moment one does.

func TestPackQueryResponseRejectsTypedNil(t *testing.T) {
	// Both shapes: a branch that marshals via proto.Message, and a branch that
	// reads a field off the concrete type. Both assert true on a typed nil.
	cases := []struct {
		name   string
		method string
		resp   interface{}
	}{
		{"marshal branch", GetPoolMethod, (*gammtypes.QueryPoolResponse)(nil)},
		{"field-access branch", GetPoolTypeMethod, (*gammtypes.QueryPoolTypeResponse)(nil)},
		{"default marshal branch", GetTotalLiquidityMethod, (*gammtypes.QueryTotalLiquidityResponse)(nil)},
	}

	p := NewPrecompile(gammkeeper.Keeper{})
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, ok := p.Methods[tc.method]
			require.True(t, ok, "method %s missing from ABI", tc.method)

			// Must return an error. Before the guard this panicked.
			_, err := p.packQueryResponse(&m, tc.resp)
			require.Error(t, err)
		})
	}
}

// Pins the reason the guard is needed. gogoproto marshals a typed nil by
// dereferencing it; the SDK's ProtoCodec does not. If gogoproto ever stops
// panicking here, this fails and the guard can be reconsidered.
func TestGogoprotoMarshalPanicsOnTypedNil(t *testing.T) {
	var typedNil *gammtypes.QueryPoolResponse
	var asIface interface{} = typedNil

	// Deliberately a bare `==`, not require.NotNil: testify reflects into the
	// value and reports a typed nil as nil, which is the very distinction
	// under test here.
	require.False(t, asIface == nil, "a typed nil must not compare equal to interface nil")

	pm, ok := asIface.(proto.Message)
	require.True(t, ok, "a typed nil still satisfies proto.Message")

	require.Panics(t, func() {
		_, _ = proto.Marshal(pm) //nolint:errcheck // the panic is the assertion
	})
}

func TestPackQueryResponseAcceptsRealResponse(t *testing.T) {
	p := NewPrecompile(gammkeeper.Keeper{})
	m, ok := p.Methods[GetPoolTypeMethod]
	require.True(t, ok)

	out, err := p.packQueryResponse(&m, &gammtypes.QueryPoolTypeResponse{PoolType: "balancer"})
	require.NoError(t, err)
	require.NotEmpty(t, out)
}
