// Helper for distinguishing "field unset" from "field set to zero value"
// after the chain normalizes nil pointers on load.
//
// Background: the proto schema marks nested struct fields
// `(gogoproto.nullable) = true` so the canonical Amino JSON omits
// them when the SDK's `dropEmptyProtoSubMessages` strips empty
// sub-messages from the proto wire. This is required for EIP-712
// typed-data parity with what MetaMask signs.
//
// To prevent nil-deref panics across ~170 keeper access sites,
// `keeper.NormalizeNilPointers` fills nil pointer-to-struct fields
// with fresh zero-value instances after storage load. After
// normalization, `field != nil` no longer reliably distinguishes
// "user set this field" from "user did not set this field" — both
// produce a non-nil pointer (one points to a zero-value struct,
// the other to a populated one).
//
// `IsBasicallyEmpty` recovers that distinction by checking the
// proto-encoded size: a populated message has wire size > 0, a
// zero-value message has size 0. Use this at "is field set"
// semantic check sites; keep `!= nil` for plain nil-safety guards
// (those continue to work after normalize).

package types

// SizedProto is the subset of proto messages we care about — every
// gogoproto-generated type implements `Size() int` and handles a nil
// receiver by returning 0.
type SizedProto interface {
	Size() int
}

// IsBasicallyEmpty reports whether `m` is conceptually unset:
// either a nil interface, a typed nil pointer (Size returns 0), or
// a non-nil pointer to an all-zero struct (Size returns 0).
//
// Pair with `keeper.NormalizeNilPointers` — that helper guarantees
// non-nil pointers everywhere downstream code reads, but this helper
// recovers the "user actually set this" semantic for security/permission
// checks like `if !IsBasicallyEmpty(invariants.CosmosCoinBackedPath)`.
func IsBasicallyEmpty(m SizedProto) bool {
	if m == nil {
		return true
	}
	return m.Size() == 0
}
