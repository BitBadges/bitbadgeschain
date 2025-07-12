package types

const (
	// Note: Actual wasmx had "xwasm" to avoid prefix collision but we don't actually store anything here. This is just a wrapper module, so this should be fine
	// ModuleName defines the module name
	ModuleName = "wasmx"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_wasmx"

	// Version defines the current version the IBC module supports
	Version = "wasmx-1"

	// PortID is the default port id that module binds to
	PortID = "wasmx"
)

var (
	ParamsKey = []byte("p_wasmx")
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("wasmx-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
