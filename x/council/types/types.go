package types

// Council represents a token-gated collective that can execute arbitrary sdk.Msgs.
// Membership is defined by ownership of an x/tokenization credential token.
type Council struct {
	Id                     uint64   `json:"id"`
	Creator                string   `json:"creator"`
	CredentialCollectionId uint64   `json:"credentialCollectionId"`
	CredentialTokenId      uint64   `json:"credentialTokenId"`
	VotingThreshold        uint64   `json:"votingThreshold"`        // percentage 0-100 of total credential supply needed
	ExecutionDelay         int64    `json:"executionDelay"`         // milliseconds after vote passes before execution allowed
	AllowedMsgTypes        []string `json:"allowedMsgTypes"`        // optional: restrict which msg type URLs can be proposed (empty = any)
	AccountAddress         string   `json:"accountAddress"`         // derived module account that executes messages
}

// Proposal represents a pending set of messages to be executed by a council.
type Proposal struct {
	CouncilId   uint64   `json:"councilId"`
	ProposalId  uint64   `json:"proposalId"`
	Proposer    string   `json:"proposer"`
	MsgTypeUrls []string `json:"msgTypeUrls"` // type URLs of the messages
	MsgBytes    [][]byte `json:"msgBytes"`    // proto/amino-encoded messages
	Status      string   `json:"status"`      // "pending" | "passed" | "executed" | "expired"
	YesWeight   uint64   `json:"yesWeight"`
	NoWeight    uint64   `json:"noWeight"`
	Deadline    int64    `json:"deadline"` // unix ms — voting deadline
	PassedAt    int64    `json:"passedAt"` // unix ms — when threshold was met
}

// Vote represents a single voter's vote on a proposal.
type Vote struct {
	CouncilId  uint64 `json:"councilId"`
	ProposalId uint64 `json:"proposalId"`
	Voter      string `json:"voter"`
	Weight     uint64 `json:"weight"`
	VoteYes    bool   `json:"voteYes"`
}

// Proposal statuses
const (
	ProposalStatusPending  = "pending"
	ProposalStatusPassed   = "passed"
	ProposalStatusExecuted = "executed"
	ProposalStatusExpired  = "expired"
)
