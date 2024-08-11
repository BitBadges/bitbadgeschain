package types

const (
	// ModuleName defines the module name
	ModuleName = "offers"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_offers"

	// Version defines the current version the IBC module supports
	Version = "offers-1"

	// PortID is the default port id that module binds to
	PortID = "offers"
)

var (
	ParamsKey = []byte("p_offers")
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("offers-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
