package types

// Canonical IBC denominations recognised by x/tokenization.
//
// These are hashes of the full IBC denom trace, so they encode the exact route
// a token took to reach this chain. Two routes for the same underlying asset
// produce two different denoms, which is why USDC appears twice below.
const (
	// USDCDenom is the canonical USDC on BitBadges: Circle's NATIVE USDC on
	// Injective (USDC.inj, erc20 contract
	// 0xa00C59fF5a080D2b954d0c75e46E22a0c371235a, held as an Injective bank
	// denom), sent ONE hop over the existing BitBadges<->Injective channel.
	//
	//	trace: transfer/channel-40/erc20:0xa00C59fF5a080D2b954d0c75e46E22a0c371235a
	//	       channel-40  BitBadges -> Injective
	//
	// The erc20 address inside the trace is CHECKSUMMED, exactly as Injective's
	// bank module spells the denom — the hash is case-sensitive, so the
	// lowercase spelling would be a different (wrong) denom.
	//
	// Deliberately NOT the Noble voucher held on Injective
	// (transfer/channel-148/uusdc there): forwarding that voucher would mint a
	// 2-hop denom here. Canonical USDC is the single-hop native asset.
	// Everything new should be priced and quoted in this denom.
	USDCDenom = "ibc/E1116484B327AEE59CDC3DA73D319834781A13DB2A7DFC1F38A30CD45ABF58B8"

	// USDCNobleDenom is the legacy Noble-direct USDC.
	//
	//	trace: transfer/channel-2/uusdc
	//
	// Deprecated for new use, but deliberately still allowed. It remains the
	// backing asset for collections that declared a backed path against it —
	// the backed-path escrow address is derived from the denom string itself
	// (see generatePathAddress in x/tokenization/keeper), so those collections
	// cannot be repointed without moving their escrow. They stay on this denom.
	//
	// Displayed to users as "USDC.noble" to distinguish it from the canonical
	// USDC above.
	USDCNobleDenom = "ibc/F082B65C88E4B6D5EF1DB243CDA1D331D002759E938A0F5CD3FFDC5D53B3E349"

	// ATOMDenom is Cosmos Hub ATOM (trace: transfer/channel-3/uatom).
	ATOMDenom = "ibc/A4DB47A9D3CF9A068D454513891B526702455D3EF08FB9EB558C561F9DC2B701"

	// OSMODenom is Osmosis OSMO (trace: transfer/channel-0/uosmo).
	OSMODenom = "ibc/ED07A3391A112B175915CD8FAF43A2DA8E4790EDE12566649D0C2F97716B8518"

	// NativeDenom is the chain's own staking and gas token.
	NativeDenom = "ubadge"
)
