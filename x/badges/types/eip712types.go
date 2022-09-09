package types

import (
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

func GetMsgValueTypes(route string) []apitypes.Type {
	switch route {
	case TypeMsgNewBadge:
		return []apitypes.Type{ 
			{ Name: "creator", Type: "string" },
			{ Name: "uri", Type: "UriObject" },
			{ Name: "arbitraryBytes", Type: "bytes" },
			{ Name: "permissions", Type: "uint64" },
			{ Name: "defaultSubassetSupply", Type: "uint64" },
			{ Name: "freezeAddressRanges", Type: "IdRange[]" },
			{ Name: "standard", Type: "uint64" },
			{ Name: "subassetSupplys", Type: "uint64[]" },
			{ Name: "subassetAmountsToCreate", Type: "uint64[]" },
		}
	case TypeMsgNewSubBadge:
		return []apitypes.Type{ 
			{ Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "supplys", Type: "uint64[]" },
			{ Name: "amountsToCreate", Type: "uint64[]" },
		}
	case TypeMsgTransferBadge:
		return []apitypes.Type{	  { Name: "creator", Type: "string" },
			{ Name: "from", Type: "uint64" },
			{ Name: "toAddresses", Type: "uint64[]" },
			{ Name: "amounts", Type: "uint64[]" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "subbadgeRanges", Type: "IdRange[]" },
			{ Name: "expiration_time", Type: "uint64" },
			{ Name: "cantCancelBeforeTime", Type: "uint64" },
		}
	case TypeMsgRequestTransferBadge:
		return []apitypes.Type{{ Name: "creator", Type: "string" },
			{ Name: "from", Type: "uint64" },
			{ Name: "amount", Type: "uint64" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "subbadgeRanges", Type: "IdRange[]" },
			{ Name: "expiration_time", Type: "uint64" },
			{ Name: "cantCancelBeforeTime", Type: "uint64" },
		}
	case TypeMsgHandlePendingTransfer:
		return []apitypes.Type{{ Name: "creator", Type: "string" },
			{ Name: "accept", Type: "bool" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "nonceRanges", Type: "IdRange[]" },
			{ Name: "forcefulAccept", Type: "bool" },
		}
	case TypeMsgSetApproval:
		return []apitypes.Type{
			{ Name: "creator", Type: "string" },
			{ Name: "amount", Type: "uint64" },
			{ Name: "address", Type: "uint64" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "subbadgeRanges", Type: "IdRange[]" },
		}
	case TypeMsgRevokeBadge:
		return []apitypes.Type{  { Name: "creator", Type: "string" },
			{ Name: "addresses", Type: "uint64[]" },
			{ Name: "amounts", Type: "uint64[]" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "subbadgeRanges", Type: "IdRange[]" },
		}
	case TypeMsgFreezeAddress:
		return []apitypes.Type{
			{ Name: "creator", Type: "string" },
			{ Name: "addressRanges", Type: "IdRange[]" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "add", Type: "bool" },
		}
	case TypeMsgUpdateUris:
		return []apitypes.Type{	 { Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "uri", Type: "UriObject" },
		}
	case TypeMsgUpdatePermissions:
		return []apitypes.Type{{ Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "permissions", Type: "uint64" },
		}
	case TypeMsgUpdateBytes:
		return []apitypes.Type{  { Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "newBytes", Type: "bytes" },
		}
	case TypeMsgTransferManager:
		return []apitypes.Type{ { Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "address", Type: "uint64" },
		}
	case TypeMsgRequestTransferManager:
		return []apitypes.Type{ { Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
			{ Name: "add", Type: "bool" },
		}
	case TypeMsgSelfDestructBadge:
		return []apitypes.Type{ { Name: "creator", Type: "string" },
			{ Name: "badgeId", Type: "uint64" },
		}
	case TypeMsgPruneBalances:
		return []apitypes.Type{ { Name: "creator", Type: "string" },
			{ Name: "badgeIds", Type: "uint64[]" },
			{ Name: "addresses", Type: "uint64[]" },
		}
	case TypeMsgRegisterAddresses:
		return []apitypes.Type{{ Name: "creator", Type: "string" },
			{ Name: "addressesToRegister", Type: "string[]" },
		}
	default:
		return []apitypes.Type{}
	};
}
