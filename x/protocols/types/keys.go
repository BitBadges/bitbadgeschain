package types

const (
	// ModuleName defines the module name
	ModuleName = "protocols"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_protocols"

	// Version defines the current version the IBC module supports
	Version = "protocols-1"

	// PortID is the default port id that module binds to
	PortID = "protocols"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("protocols-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
