package types

func getDuplicatesAndNonDuplicates(list1 []string, list2 []string) ([]string, []string) {
	duplicates := []string{}
	inListOneButNotTwo := []string{}

	for _, address := range list1 {
		//Check if address is in list2
		found := false
		for _, address2 := range list2 {
			if address == address2 {
				found = true
				break
			}
		}

		if found {
			duplicates = append(duplicates, address)
		} else {
			inListOneButNotTwo = append(inListOneButNotTwo, address)
		}
	}

	return duplicates, inListOneButNotTwo
}

//Each address list has a list of addresses and a boolean whitelist.
//Four cases (toRemove.Whitelist, addressList.Whitelist):
// 1) (true, true) - Remove ABC from BCD
//    Removed - duplicates from toRemove.Addresses and addressList.Addresses (BC)
//    Remaining - non-duplicates from addressList.Addresses (D)
// 2) (false, true) - Remove All but ABC from BCD
//    Removed - non-duplicates from addressList.Addresses (D)
//    Remaining - duplicates from toRemove.Addresses and addressList.Addresses (BC)
// 3) (true, false) - Remove ABC from All but BCD
//    Removed - non-duplicates from toRemove.Addresses (A)
//		Remaining - everyone but combined list of toRemove.Addresses and addressList.Addresses (everyone but ABCD)
// 4) (false, false) - Remove All but ABC from All but BCD
//		Removed - everyone but combined list of toRemove.Addresses and addressList.Addresses (everyone but ABCD)
//		Remaining - non-duplicates from toRemove.Addresses (A)
func RemoveAddressListFromAddressList(listToRemove *AddressList, addressList *AddressList) (*AddressList, *AddressList) {
	duplicates, inToRemoveButNotList := getDuplicatesAndNonDuplicates(listToRemove.Addresses, addressList.Addresses)
	_, inListButNotToRemove := getDuplicatesAndNonDuplicates(addressList.Addresses, listToRemove.Addresses)

	removed := &AddressList{}
	remaining := &AddressList{}

	if listToRemove.Whitelist && addressList.Whitelist {
		//Case 1
		removed.Whitelist = true
		removed.Addresses = duplicates

		remaining.Whitelist = true
		remaining.Addresses = inListButNotToRemove
	} else if !listToRemove.Whitelist && addressList.Whitelist {
		//Case 2
		removed.Whitelist = true
		removed.Addresses = inListButNotToRemove

		remaining.Whitelist = true
		remaining.Addresses = duplicates
	} else if listToRemove.Whitelist && !addressList.Whitelist {
		//Case 3
		removed.Whitelist = true
		removed.Addresses = inToRemoveButNotList

		remaining.Whitelist = false
		remaining.Addresses = append(remaining.Addresses, inListButNotToRemove...)
		remaining.Addresses = append(remaining.Addresses, inToRemoveButNotList...)
		remaining.Addresses = append(remaining.Addresses, duplicates...)
	} else if !listToRemove.Whitelist && !addressList.Whitelist {
		//Case 4
		removed.Whitelist = false
		removed.Addresses = append(removed.Addresses, inListButNotToRemove...)
		removed.Addresses = append(removed.Addresses, inToRemoveButNotList...)
		removed.Addresses = append(removed.Addresses, duplicates...)

		remaining.Whitelist = true
		remaining.Addresses = inToRemoveButNotList
	}

	return remaining, removed
}
