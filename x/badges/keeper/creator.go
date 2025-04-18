package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetCreator(ctx sdk.Context, creator string, creatorOverride string) (string, error) {

	// Check if creator is a contract address or normal address
	// Contract addresses are 32 bytes and normal addresses are 20 bytes
	accAddress, err := sdk.AccAddressFromBech32(creator)
	if err != nil {
		return "", err
	}

	creatorIsContract := len(accAddress.Bytes()) == 32
	isApprovedContract := false
	for _, address := range k.ApprovedContractAddresses {
		if address == creator {
			isApprovedContract = true
		}
	}

	if creatorIsContract && !isApprovedContract {
		return "", fmt.Errorf("the only entrypoint for modules and contracts is via an approved contract address")
	}

	if creator == creatorOverride {
		return creator, nil
	}

	// If creatorOverride is set, we need to verify actual creator is an approved contract address
	// IMPORTANT: Approved contract addresses should never be allowed to specify alternate creators other than the initial signer themselves
	// This is to prevent malicious contracts from overriding the creator and bypassing all permissions
	if creatorOverride != "" {
		if isApprovedContract {
			return creatorOverride, nil
		}
		return "", fmt.Errorf("the only entrypoint for modules and contracts is via an approved contract address")
	}

	// Otherwise, use the original creator
	return creator, nil
}
