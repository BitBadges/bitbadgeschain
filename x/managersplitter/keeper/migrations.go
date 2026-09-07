package keeper

import (
	"bytes"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/pkg/storewalk"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) MigrateV35CanonicalAddresses(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	canonical := func(value string) string {
		address, err := sdk.AccAddressFromBech32(value)
		if err != nil {
			return value
		}
		return address.String()
	}
	return storewalk.Prefix(ctx, store, ManagerSplitterKey, func(key, value []byte) error {
		var splitter types.ManagerSplitter
		if err := k.cdc.Unmarshal(value, &splitter); err != nil {
			return err
		}
		splitter.Admin = canonical(splitter.Admin)
		splitter.Address = canonical(splitter.Address)
		for _, permission := range types.PermissionFields(splitter.Permissions) {
			if permission == nil {
				continue
			}
			seen := map[string]bool{}
			addresses := []string{}
			for _, address := range permission.ApprovedAddresses {
				address = canonical(address)
				if !seen[address] {
					addresses = append(addresses, address)
					seen[address] = true
				}
			}
			permission.ApprovedAddresses = addresses
		}
		target := managerSplitterStoreKey(splitter.Address)
		if !bytes.Equal(key, target) && store.Has(target) {
			return fmt.Errorf("conflicting manager splitter identities: %s", splitter.Address)
		}
		encoded, err := k.cdc.Marshal(&splitter)
		if err != nil {
			return err
		}
		if !bytes.Equal(key, target) || !bytes.Equal(value, encoded) {
			store.Set(target, encoded)
		}
		if !bytes.Equal(key, target) {
			store.Delete(key)
		}
		return nil
	})
}
