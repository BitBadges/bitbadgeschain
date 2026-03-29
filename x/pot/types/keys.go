package types

const (
	// ModuleName defines the module name
	ModuleName = "pot"

	// StoreKey defines the primary module store key (for params and compliance-jailed set)
	StoreKey = ModuleName
)

var (
	// ParamsKey is the key for storing module params
	ParamsKey = []byte("p_pot")

	// ComplianceJailedPrefix is the prefix for the set of validator consensus addresses
	// that x/pot has jailed for missing credentials. This lets us distinguish
	// compliance-jailing from slashing-jailing so we can auto-unjail when
	// credentials are regained.
	// Each entry is ComplianceJailedPrefix | consAddr_bytes -> []byte{1}
	ComplianceJailedPrefix = []byte("cj/")

	// SavedPowerPrefix is the prefix for persisting a validator's voting power
	// before the PoA adapter disables them. When the validator is re-enabled,
	// the saved power is restored.
	// Each entry is SavedPowerPrefix | consAddr_bytes -> big-endian int64
	SavedPowerPrefix = []byte("sp/")
)

// ComplianceJailedKey returns the KV store key for tracking a compliance-jailed validator.
// Uses make+copy to avoid corrupting the package-level prefix slice.
func ComplianceJailedKey(consAddr []byte) []byte {
	key := make([]byte, len(ComplianceJailedPrefix)+len(consAddr))
	copy(key, ComplianceJailedPrefix)
	copy(key[len(ComplianceJailedPrefix):], consAddr)
	return key
}

// SavedPowerKey returns the KV store key for a validator's saved power.
func SavedPowerKey(consAddr []byte) []byte {
	key := make([]byte, len(SavedPowerPrefix)+len(consAddr))
	copy(key, SavedPowerPrefix)
	copy(key[len(SavedPowerPrefix):], consAddr)
	return key
}
