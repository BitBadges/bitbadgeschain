package types

func getDuplicatesAndNonDuplicates(list1 []string, list2 []string) ([]string, []string) {
	duplicates := []string{}
	inListOneButNotTwo := []string{}

	// Create a map for O(1) lookup instead of O(n) nested loop
	list2Map := make(map[string]bool, len(list2))
	for _, address := range list2 {
		list2Map[address] = true
	}

	for _, address := range list1 {
		// Check if address is in list2 using map lookup
		if list2Map[address] {
			duplicates = append(duplicates, address)
		} else {
			inListOneButNotTwo = append(inListOneButNotTwo, address)
		}
	}

	return duplicates, inListOneButNotTwo
}

// Each address list has a list of addresses and a boolean whitelist.
// Four cases (toRemove.Whitelist, addressList.Whitelist):
//  1. (true, true) - Remove ABC from BCD
//     Removed - duplicates from toRemove.Addresses and addressList.Addresses (BC)
//     Remaining - non-duplicates from addressList.Addresses (D)
//  2. (false, true) - Remove All but ABC from BCD
//     Removed - non-duplicates from addressList.Addresses (D)
//     Remaining - duplicates from toRemove.Addresses and addressList.Addresses (BC)
//  3. (true, false) - Remove ABC from All but BCD
//     Removed - non-duplicates from toRemove.Addresses (A)
//     Remaining - everyone but combined list of toRemove.Addresses and addressList.Addresses (everyone but ABCD)
//  4. (false, false) - Remove All but ABC from All but BCD
//     Removed - everyone but combined list of toRemove.Addresses and addressList.Addresses (everyone but ABCD)
//     Remaining - non-duplicates from toRemove.Addresses (A)
func RemoveAddressListFromAddressList(listToRemove *AddressList, addressList *AddressList) (*AddressList, *AddressList) {
	// Validate input parameters
	if listToRemove == nil {
		panic("listToRemove cannot be nil")
	}
	if addressList == nil {
		panic("addressList cannot be nil")
	}

	duplicates, inToRemoveButNotList := getDuplicatesAndNonDuplicates(listToRemove.Addresses, addressList.Addresses)
	_, inListButNotToRemove := getDuplicatesAndNonDuplicates(addressList.Addresses, listToRemove.Addresses)

	removed := &AddressList{}
	remaining := &AddressList{}

	switch {
	case listToRemove.Whitelist && addressList.Whitelist:
		handleCase1(removed, remaining, duplicates, inListButNotToRemove)
	case !listToRemove.Whitelist && addressList.Whitelist:
		handleCase2(removed, remaining, duplicates, inListButNotToRemove)
	case listToRemove.Whitelist && !addressList.Whitelist:
		handleCase3(removed, remaining, duplicates, inToRemoveButNotList, inListButNotToRemove)
	case !listToRemove.Whitelist && !addressList.Whitelist:
		handleCase4(removed, remaining, duplicates, inToRemoveButNotList, inListButNotToRemove)
	}

	return remaining, removed
}

// Case 1: (true, true) - Remove ABC from BCD
func handleCase1(removed, remaining *AddressList, duplicates, inListButNotToRemove []string) {
	removed.Whitelist = true
	removed.Addresses = duplicates

	remaining.Whitelist = true
	remaining.Addresses = inListButNotToRemove
}

// Case 2: (false, true) - Remove All but ABC from BCD
func handleCase2(removed, remaining *AddressList, duplicates, inListButNotToRemove []string) {
	removed.Whitelist = true
	removed.Addresses = inListButNotToRemove

	remaining.Whitelist = true
	remaining.Addresses = duplicates
}

// Case 3: (true, false) - Remove ABC from All but BCD
func handleCase3(removed, remaining *AddressList, duplicates, inToRemoveButNotList, inListButNotToRemove []string) {
	removed.Whitelist = true
	removed.Addresses = inToRemoveButNotList

	remaining.Whitelist = false
	remaining.Addresses = append(remaining.Addresses, inListButNotToRemove...)
	remaining.Addresses = append(remaining.Addresses, inToRemoveButNotList...)
	remaining.Addresses = append(remaining.Addresses, duplicates...)
}

// Case 4: (false, false) - Remove All but ABC from All but BCD
func handleCase4(removed, remaining *AddressList, duplicates, inToRemoveButNotList, inListButNotToRemove []string) {
	removed.Whitelist = false
	removed.Addresses = append(removed.Addresses, inListButNotToRemove...)
	removed.Addresses = append(removed.Addresses, inToRemoveButNotList...)
	removed.Addresses = append(removed.Addresses, duplicates...)

	remaining.Whitelist = true
	remaining.Addresses = inToRemoveButNotList
}
