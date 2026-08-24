package types

// Canonical IBC denominations recognised by x/tokenization.
//
// These are hashes of the full IBC denom trace, so they encode the exact route
// a token took to reach this chain. Two routes for the same underlying asset
// produce two different denoms, which is why USDC appears twice below.
const (
	// USDCDenom is the canonical USDC on BitBadges, routed through Injective.
	//
	//	trace: transfer/channel-40/transfer/channel-148/uusdc
	//	       channel-40  BitBadges -> Injective
	//	       channel-148 Injective -> Noble
	//
	// Everything new should be priced and quoted in this denom.
	USDCDenom = "ibc/0E485657AEF4C39D551E7D53463734E4C445A96E6C814DC4C2FF0031470B40BB"

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
