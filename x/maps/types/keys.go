package types

const (
	// ModuleName defines the module name
	ModuleName = "maps"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_maps"

	// Version defines the current version the IBC module supports
	Version = "maps-1"

	// PortID is the default port id that module binds to
	PortID = "maps"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("maps-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
