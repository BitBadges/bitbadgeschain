package types



func RemoveAddressMappingFromAddressMapping(mappingToRemove *AddressMapping, addressMapping *AddressMapping) (*AddressMapping, *AddressMapping) {
	//Each address mapping has a list of addresses and a boolean includeOnlySpecified.
	//Four cases (toRemove.IncludeOnlySpecified, addressMapping.IncludeOnlySpecified):
	// 1) (true, true) - Remove ABC from BCD
	//    Removed - duplicates from toRemove.Addresses and addressMapping.Addresses (BC)
	//    Remaining - non-duplicates from addressMapping.Addresses (D)
	// 2) (false, true) - Remove All but ABC from BCD
	//    Removed - non-duplicates from addressMapping.Addresses (D)
	//    Remaining - duplicates from toRemove.Addresses and addressMapping.Addresses (BC)
	// 3) (true, false) - Remove ABC from All but BCD
	//    Removed - non-duplicates from toRemove.Addresses (A)
	//		Remaining - everyone but combined list of toRemove.Addresses and addressMapping.Addresses (everyone but ABCD)
	// 4) (false, false) - Remove All but ABC from All but BCD
	//		Removed - everyone but combined list of toRemove.Addresses and addressMapping.Addresses (everyone but ABCD)
	//		Remaining - non-duplicates from toRemove.Addresses (A)

	duplicates := []string{}
	inToRemoveButNotMapping := []string{}
	inMappingButNotToRemove := []string{}


	for _, address := range mappingToRemove.Addresses {
		//Check if address is in addressMapping.Addresses
		found := false
		for _, address2 := range addressMapping.Addresses {
			if address == address2 {
				found = true
				break
			}
		}

		if found {
			duplicates = append(duplicates, address)
		} else {
			inToRemoveButNotMapping = append(inToRemoveButNotMapping, address)
		}
	}

	for _, address := range addressMapping.Addresses {
		//Check if address is in mappingToRemove.Addresses
		found := false
		for _, address2 := range mappingToRemove.Addresses {
			if address == address2 {
				found = true
				break
			}
		}

		if !found {
			inMappingButNotToRemove = append(inMappingButNotToRemove, address)
		}
	}

	removed := &AddressMapping{}
	remaining := &AddressMapping{}

	
	if mappingToRemove.IncludeOnlySpecified && addressMapping.IncludeOnlySpecified {
		//Case 1
		removed.IncludeOnlySpecified = true
		removed.Addresses = duplicates

		remaining.IncludeOnlySpecified = true
		remaining.Addresses = inMappingButNotToRemove
	} else if !mappingToRemove.IncludeOnlySpecified && addressMapping.IncludeOnlySpecified {
		//Case 2
		removed.IncludeOnlySpecified = true
		removed.Addresses = inMappingButNotToRemove

		remaining.IncludeOnlySpecified = true
		remaining.Addresses = duplicates
	} else if mappingToRemove.IncludeOnlySpecified && !addressMapping.IncludeOnlySpecified {
		//Case 3
		removed.IncludeOnlySpecified = true
		removed.Addresses = inToRemoveButNotMapping

		remaining.IncludeOnlySpecified = false
		remaining.Addresses = append(remaining.Addresses, inMappingButNotToRemove...)
		remaining.Addresses = append(remaining.Addresses, inToRemoveButNotMapping...)
		remaining.Addresses = append(remaining.Addresses, duplicates...)
	} else if !mappingToRemove.IncludeOnlySpecified && !addressMapping.IncludeOnlySpecified {
		//Case 4
		removed.IncludeOnlySpecified = false
		removed.Addresses = append(removed.Addresses, inMappingButNotToRemove...)
		removed.Addresses = append(removed.Addresses, inToRemoveButNotMapping...)
		removed.Addresses = append(removed.Addresses, duplicates...)

		remaining.IncludeOnlySpecified = true
		remaining.Addresses = inToRemoveButNotMapping
	}

	return remaining, removed
}