package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// nil interface → true.
func TestIsBasicallyEmpty_NilInterface(t *testing.T) {
	require.True(t, IsBasicallyEmpty(nil))
}

// Typed nil pointer (interface holds nil but not typed nil) → true.
// gogoproto's Size() returns 0 for nil receiver, so this works even
// without an explicit reflect check.
func TestIsBasicallyEmpty_TypedNilPointer(t *testing.T) {
	var p *CosmosCoinBackedPath
	require.True(t, IsBasicallyEmpty(p))
}

// Zero-value struct (no fields populated) → true.
// This is the "after normalize-on-load" case: the field is non-nil but
// conceptually unset.
func TestIsBasicallyEmpty_ZeroValueStruct(t *testing.T) {
	p := &CosmosCoinBackedPath{}
	require.True(t, IsBasicallyEmpty(p), "fresh empty struct must be considered unset")
}

// Populated struct → false.
func TestIsBasicallyEmpty_PopulatedStruct(t *testing.T) {
	p := &CosmosCoinBackedPath{
		Address: "bb1xyz",
	}
	require.False(t, IsBasicallyEmpty(p), "struct with set field must NOT be considered empty")
}

// Partially populated still counts as not-empty (single field is enough).
func TestIsBasicallyEmpty_BoolFieldOnly(t *testing.T) {
	inv := &CollectionInvariants{
		DisablePoolCreation: true,
	}
	require.False(t, IsBasicallyEmpty(inv))
}

// Test against AddressChecks (one of our main empty struct types).
func TestIsBasicallyEmpty_AddressChecks(t *testing.T) {
	require.True(t, IsBasicallyEmpty((*AddressChecks)(nil)))
	require.True(t, IsBasicallyEmpty(&AddressChecks{}))
	require.False(t, IsBasicallyEmpty(&AddressChecks{MustBeEvmContract: true}))
}

// Test against AutoDeletionOptions.
func TestIsBasicallyEmpty_AutoDeletionOptions(t *testing.T) {
	require.True(t, IsBasicallyEmpty((*AutoDeletionOptions)(nil)))
	require.True(t, IsBasicallyEmpty(&AutoDeletionOptions{}))
	require.False(t, IsBasicallyEmpty(&AutoDeletionOptions{AfterOneUse: true}))
}

// Critical regression: after `keeper.NormalizeNilPointers` walks a
// fresh struct it leaves every sub-field as a non-nil pointer to a
// zero-value struct. `Size()` on the parent returns nonzero (it
// counts the embedded-message tag + length-0 marker), but the
// container is conceptually still unset. IsBasicallyEmpty must
// see through this.
func TestIsBasicallyEmpty_NormalizeFilledStillEmpty(t *testing.T) {
	// `&CosmosCoinBackedPath{Conversion: &Conversion{}}` — exactly the
	// shape after a normalize fill on a fresh struct that started nil.
	post := &CosmosCoinBackedPath{
		Conversion: &Conversion{},
	}
	require.True(t, IsBasicallyEmpty(post),
		"normalize-filled empty chain must be reported as 'unset' even though Size() > 0")

	// And a deeper nest still empty.
	deeper := &CosmosCoinBackedPath{
		Conversion: &Conversion{
			SideA: &ConversionSideAWithDenom{},
			SideB: nil,
		},
	}
	require.True(t, IsBasicallyEmpty(deeper))

	// Once any leaf has a real value, it's no longer empty.
	withValue := &CosmosCoinBackedPath{
		Conversion: &Conversion{
			SideA: &ConversionSideAWithDenom{Denom: "ucredit"},
		},
	}
	require.False(t, IsBasicallyEmpty(withValue))
}

// Verify the actual bootstrap-blocking scenario: a CollectionInvariants
// with CosmosCoinBackedPath set vs unset.
func TestIsBasicallyEmpty_BootstrapScenario(t *testing.T) {
	// Scenario A: user did not configure CosmosCoinBackedPath. Before
	// normalize-on-load, this would be `nil`. After normalize, this is
	// `&CosmosCoinBackedPath{}` (non-nil empty).
	invUnset := &CollectionInvariants{
		CosmosCoinBackedPath: &CosmosCoinBackedPath{}, // post-normalize empty
	}
	require.True(t, IsBasicallyEmpty(invUnset.CosmosCoinBackedPath),
		"post-normalize empty must be reported as 'unset' so Mint transfers aren't wrongly blocked")

	// Scenario B: user configured a real backed path.
	invSet := &CollectionInvariants{
		CosmosCoinBackedPath: &CosmosCoinBackedPath{Address: "bb1xyz"},
	}
	require.False(t, IsBasicallyEmpty(invSet.CosmosCoinBackedPath),
		"populated CosmosCoinBackedPath must be reported as set")
}
