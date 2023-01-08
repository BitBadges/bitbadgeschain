package keeper

var (
	// ReservedAddressAccountNumber is the account number for the reserved address
	ReservedAddresses = []uint64{
		
	}
)

func IsReservedAddress(address uint64) bool {
	for _, reserved := range ReservedAddresses {
		if reserved == address {
			return true
		}
	}
	return false
}