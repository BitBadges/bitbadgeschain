package types

const (
	// ModuleName defines the module name
	ModuleName = "council"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	// CouncilKeyPrefix is the prefix for council entries: "c/" | council_id_bytes
	CouncilKeyPrefix = []byte("c/")

	// ProposalKeyPrefix is the prefix for proposal entries: "p/" | council_id_bytes | "/" | proposal_id_bytes
	ProposalKeyPrefix = []byte("p/")

	// VoteKeyPrefix is the prefix for vote entries: "v/" | council_id_bytes | "/" | proposal_id_bytes | "/" | voter_address
	VoteKeyPrefix = []byte("v/")

	// NextCouncilIdKey stores the next available council ID
	NextCouncilIdKey = []byte("next_council_id")

	// NextProposalIdKeyPrefix stores the next proposal ID per council: "np/" | council_id_bytes
	NextProposalIdKeyPrefix = []byte("np/")
)
