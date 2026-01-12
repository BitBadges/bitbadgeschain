package types

const (
	// ModuleName defines the module name
	ModuleName = "badges"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_badges"

	// Version defines the current version the IBC module supports
	Version = "badges-1"

	// PortID is the default port id that module binds to
	PortID = "badges"
)

var ParamsKey = []byte("p_badges")

// PortKey defines the key to store the port ID in store
var PortKey = KeyPrefix("badges-port-")

func KeyPrefix(p string) []byte {
	return []byte(p)
}
