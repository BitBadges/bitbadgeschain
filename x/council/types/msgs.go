package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------------------------------------------------------------------------
// MsgCreateCouncil
// ---------------------------------------------------------------------------

type MsgCreateCouncil struct {
	Creator                string   `json:"creator"`
	CredentialCollectionId uint64   `json:"credentialCollectionId"`
	CredentialTokenId      uint64   `json:"credentialTokenId"`
	VotingThreshold        uint64   `json:"votingThreshold"`
	ExecutionDelay         int64    `json:"executionDelay"`
	AllowedMsgTypes        []string `json:"allowedMsgTypes"`
}

type MsgCreateCouncilResponse struct {
	CouncilId uint64 `json:"councilId"`
}

func (msg *MsgCreateCouncil) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return ErrInvalidAddress
	}
	if msg.VotingThreshold == 0 || msg.VotingThreshold > 100 {
		return ErrInvalidThreshold
	}
	if msg.ExecutionDelay < 0 {
		return ErrInvalidExecutionDelay
	}
	return nil
}

func (msg *MsgCreateCouncil) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{addr}
}

func (msg *MsgCreateCouncil) ProtoMessage()             {}
func (msg *MsgCreateCouncil) Reset()                    {}
func (msg *MsgCreateCouncil) String() string            { return "" }
func (msg *MsgCreateCouncilResponse) ProtoMessage()     {}
func (msg *MsgCreateCouncilResponse) Reset()            {}
func (msg *MsgCreateCouncilResponse) String() string    { return "" }

// ---------------------------------------------------------------------------
// MsgPropose
// ---------------------------------------------------------------------------

type MsgPropose struct {
	Proposer    string   `json:"proposer"`
	CouncilId   uint64   `json:"councilId"`
	MsgTypeUrls []string `json:"msgTypeUrls"`
	MsgBytes    [][]byte `json:"msgBytes"`
	Deadline    int64    `json:"deadline"` // unix ms
}

type MsgProposeResponse struct {
	ProposalId uint64 `json:"proposalId"`
}

func (msg *MsgPropose) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Proposer); err != nil {
		return ErrInvalidAddress
	}
	if len(msg.MsgTypeUrls) == 0 {
		return ErrNoMessages
	}
	if len(msg.MsgTypeUrls) != len(msg.MsgBytes) {
		return ErrInvalidMsgCount
	}
	return nil
}

func (msg *MsgPropose) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Proposer)
	return []sdk.AccAddress{addr}
}

func (msg *MsgPropose) ProtoMessage()          {}
func (msg *MsgPropose) Reset()                 {}
func (msg *MsgPropose) String() string         { return "" }
func (msg *MsgProposeResponse) ProtoMessage()  {}
func (msg *MsgProposeResponse) Reset()         {}
func (msg *MsgProposeResponse) String() string { return "" }

// ---------------------------------------------------------------------------
// MsgVote
// ---------------------------------------------------------------------------

type MsgVote struct {
	Voter      string `json:"voter"`
	CouncilId  uint64 `json:"councilId"`
	ProposalId uint64 `json:"proposalId"`
	VoteYes    bool   `json:"voteYes"`
}

type MsgVoteResponse struct{}

func (msg *MsgVote) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Voter); err != nil {
		return ErrInvalidAddress
	}
	return nil
}

func (msg *MsgVote) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Voter)
	return []sdk.AccAddress{addr}
}

func (msg *MsgVote) ProtoMessage()          {}
func (msg *MsgVote) Reset()                 {}
func (msg *MsgVote) String() string         { return "" }
func (msg *MsgVoteResponse) ProtoMessage()  {}
func (msg *MsgVoteResponse) Reset()         {}
func (msg *MsgVoteResponse) String() string { return "" }

// ---------------------------------------------------------------------------
// MsgExecuteProposal
// ---------------------------------------------------------------------------

type MsgExecuteProposal struct {
	Sender     string `json:"sender"`
	CouncilId  uint64 `json:"councilId"`
	ProposalId uint64 `json:"proposalId"`
}

type MsgExecuteProposalResponse struct{}

func (msg *MsgExecuteProposal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return ErrInvalidAddress
	}
	return nil
}

func (msg *MsgExecuteProposal) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{addr}
}

func (msg *MsgExecuteProposal) ProtoMessage()          {}
func (msg *MsgExecuteProposal) Reset()                 {}
func (msg *MsgExecuteProposal) String() string         { return "" }
func (msg *MsgExecuteProposalResponse) ProtoMessage()  {}
func (msg *MsgExecuteProposalResponse) Reset()         {}
func (msg *MsgExecuteProposalResponse) String() string { return "" }
