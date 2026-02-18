package types

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

const (
	// MaxBulkQueries is the maximum number of queries allowed in a bulk request
	MaxBulkQueries = 100
)

// NewOwnershipQueryPacket creates a new OwnershipQueryPacket
func NewOwnershipQueryPacket(
	queryId string,
	address string,
	collectionId string,
	tokenId string,
	ownershipTime string,
) *OwnershipQueryPacket {
	return &OwnershipQueryPacket{
		QueryId:       queryId,
		Address:       address,
		CollectionId:  collectionId,
		TokenId:       tokenId,
		OwnershipTime: ownershipTime,
	}
}

// ValidateBasic performs basic validation of the OwnershipQueryPacket
func (p *OwnershipQueryPacket) ValidateBasic() error {
	if p.QueryId == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "query_id cannot be empty")
	}

	if p.Address == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "address cannot be empty")
	}

	if p.CollectionId == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "collection_id cannot be empty")
	}

	// Validate collection_id is a valid uint
	_, err := sdkmath.ParseUint(p.CollectionId)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid collection_id: %s", err.Error())
	}

	// Validate token_id is a valid uint
	if p.TokenId == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "token_id cannot be empty")
	}
	_, err = sdkmath.ParseUint(p.TokenId)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid token_id: %s", err.Error())
	}

	// Validate ownership_time is a valid uint
	if p.OwnershipTime == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "ownership_time cannot be empty")
	}
	_, err = sdkmath.ParseUint(p.OwnershipTime)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid ownership_time: %s", err.Error())
	}

	return nil
}

// GetBytes returns the packet data as bytes for IBC transmission
func (p *OwnershipQueryPacket) GetBytes() []byte {
	packetData := TokenizationPacketData{
		Packet: &TokenizationPacketData_OwnershipQuery{
			OwnershipQuery: p,
		},
	}
	return ModuleCdc.MustMarshal(&packetData)
}

// NewOwnershipQueryResponsePacket creates a new OwnershipQueryResponsePacket
func NewOwnershipQueryResponsePacket(
	queryId string,
	ownsTokens bool,
	totalAmount sdkmath.Uint,
	proofHeight uint64,
	errorMsg string,
) *OwnershipQueryResponsePacket {
	return &OwnershipQueryResponsePacket{
		QueryId:     queryId,
		OwnsTokens:  ownsTokens,
		TotalAmount: Uint(totalAmount),
		ProofHeight: proofHeight,
		Error:       errorMsg,
	}
}

// NewErrorOwnershipQueryResponsePacket creates an error response packet
func NewErrorOwnershipQueryResponsePacket(queryId string, err error) *OwnershipQueryResponsePacket {
	return &OwnershipQueryResponsePacket{
		QueryId:     queryId,
		OwnsTokens:  false,
		TotalAmount: Uint(sdkmath.ZeroUint()),
		Error:       err.Error(),
	}
}

// GetBytes returns the packet data as bytes for IBC transmission
func (p *OwnershipQueryResponsePacket) GetBytes() []byte {
	packetData := TokenizationPacketData{
		Packet: &TokenizationPacketData_OwnershipQueryResponse{
			OwnershipQueryResponse: p,
		},
	}
	return ModuleCdc.MustMarshal(&packetData)
}

// NewBulkOwnershipQueryPacket creates a new BulkOwnershipQueryPacket
func NewBulkOwnershipQueryPacket(
	queryId string,
	queries []*OwnershipQueryPacket,
) *BulkOwnershipQueryPacket {
	return &BulkOwnershipQueryPacket{
		QueryId: queryId,
		Queries: queries,
	}
}

// ValidateBasic performs basic validation of the BulkOwnershipQueryPacket
func (p *BulkOwnershipQueryPacket) ValidateBasic() error {
	if p.QueryId == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "query_id cannot be empty")
	}

	if len(p.Queries) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "queries cannot be empty")
	}

	if len(p.Queries) > MaxBulkQueries {
		return sdkerrors.Wrapf(ErrInvalidRequest, "too many queries: %d (max %d)", len(p.Queries), MaxBulkQueries)
	}

	// Validate each individual query
	for i, q := range p.Queries {
		if q == nil {
			return sdkerrors.Wrapf(ErrInvalidRequest, "queries[%d] is nil", i)
		}
		if err := q.ValidateBasic(); err != nil {
			return sdkerrors.Wrapf(err, "queries[%d] validation failed", i)
		}
	}

	return nil
}

// GetBytes returns the packet data as bytes for IBC transmission
func (p *BulkOwnershipQueryPacket) GetBytes() []byte {
	packetData := TokenizationPacketData{
		Packet: &TokenizationPacketData_BulkOwnershipQuery{
			BulkOwnershipQuery: p,
		},
	}
	return ModuleCdc.MustMarshal(&packetData)
}

// NewBulkOwnershipQueryResponsePacket creates a new BulkOwnershipQueryResponsePacket
func NewBulkOwnershipQueryResponsePacket(
	queryId string,
	responses []*OwnershipQueryResponsePacket,
) *BulkOwnershipQueryResponsePacket {
	return &BulkOwnershipQueryResponsePacket{
		QueryId:   queryId,
		Responses: responses,
	}
}

// GetBytes returns the packet data as bytes for IBC transmission
func (p *BulkOwnershipQueryResponsePacket) GetBytes() []byte {
	packetData := TokenizationPacketData{
		Packet: &TokenizationPacketData_BulkOwnershipQueryResponse{
			BulkOwnershipQueryResponse: p,
		},
	}
	return ModuleCdc.MustMarshal(&packetData)
}

// NewFullBalanceQueryPacket creates a new FullBalanceQueryPacket
func NewFullBalanceQueryPacket(
	queryId string,
	address string,
	collectionId string,
) *FullBalanceQueryPacket {
	return &FullBalanceQueryPacket{
		QueryId:      queryId,
		Address:      address,
		CollectionId: collectionId,
	}
}

// ValidateBasic performs basic validation of the FullBalanceQueryPacket
func (p *FullBalanceQueryPacket) ValidateBasic() error {
	if p.QueryId == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "query_id cannot be empty")
	}

	if p.Address == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "address cannot be empty")
	}

	if p.CollectionId == "" {
		return sdkerrors.Wrap(ErrInvalidRequest, "collection_id cannot be empty")
	}

	// Validate collection_id is a valid uint
	_, err := sdkmath.ParseUint(p.CollectionId)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid collection_id: %s", err.Error())
	}

	return nil
}

// GetBytes returns the packet data as bytes for IBC transmission
func (p *FullBalanceQueryPacket) GetBytes() []byte {
	packetData := TokenizationPacketData{
		Packet: &TokenizationPacketData_FullBalanceQuery{
			FullBalanceQuery: p,
		},
	}
	return ModuleCdc.MustMarshal(&packetData)
}

// NewFullBalanceQueryResponsePacket creates a new FullBalanceQueryResponsePacket
func NewFullBalanceQueryResponsePacket(
	queryId string,
	balanceStore []byte,
	proofHeight uint64,
	errorMsg string,
) *FullBalanceQueryResponsePacket {
	return &FullBalanceQueryResponsePacket{
		QueryId:      queryId,
		BalanceStore: balanceStore,
		ProofHeight:  proofHeight,
		Error:        errorMsg,
	}
}

// NewErrorFullBalanceQueryResponsePacket creates an error response packet for FullBalanceQuery
func NewErrorFullBalanceQueryResponsePacket(queryId string, err error) *FullBalanceQueryResponsePacket {
	return &FullBalanceQueryResponsePacket{
		QueryId: queryId,
		Error:   err.Error(),
	}
}

// GetBytes returns the packet data as bytes for IBC transmission
func (p *FullBalanceQueryResponsePacket) GetBytes() []byte {
	packetData := TokenizationPacketData{
		Packet: &TokenizationPacketData_FullBalanceQueryResponse{
			FullBalanceQueryResponse: p,
		},
	}
	return ModuleCdc.MustMarshal(&packetData)
}

// ICQPacketType represents the type of ICQ packet
type ICQPacketType string

const (
	ICQPacketTypeOwnershipQuery             ICQPacketType = "ownership_query"
	ICQPacketTypeOwnershipQueryResponse     ICQPacketType = "ownership_query_response"
	ICQPacketTypeBulkOwnershipQuery         ICQPacketType = "bulk_ownership_query"
	ICQPacketTypeBulkOwnershipQueryResponse ICQPacketType = "bulk_ownership_query_response"
	ICQPacketTypeFullBalanceQuery           ICQPacketType = "full_balance_query"
	ICQPacketTypeFullBalanceQueryResponse   ICQPacketType = "full_balance_query_response"
	ICQPacketTypeUnknown                    ICQPacketType = "unknown"
)

// GetICQPacketType returns the type of ICQ packet from TokenizationPacketData
func GetICQPacketType(packetData *TokenizationPacketData) ICQPacketType {
	switch packetData.Packet.(type) {
	case *TokenizationPacketData_OwnershipQuery:
		return ICQPacketTypeOwnershipQuery
	case *TokenizationPacketData_OwnershipQueryResponse:
		return ICQPacketTypeOwnershipQueryResponse
	case *TokenizationPacketData_BulkOwnershipQuery:
		return ICQPacketTypeBulkOwnershipQuery
	case *TokenizationPacketData_BulkOwnershipQueryResponse:
		return ICQPacketTypeBulkOwnershipQueryResponse
	case *TokenizationPacketData_FullBalanceQuery:
		return ICQPacketTypeFullBalanceQuery
	case *TokenizationPacketData_FullBalanceQueryResponse:
		return ICQPacketTypeFullBalanceQueryResponse
	default:
		return ICQPacketTypeUnknown
	}
}

// String returns the string representation of the packet type
func (t ICQPacketType) String() string {
	return string(t)
}

// ICQEvent names for event emission
const (
	EventTypeICQRequest  = "icq_request"
	EventTypeICQResponse = "icq_response"

	AttributeKeyQueryId      = "query_id"
	AttributeKeyAddress      = "address"
	AttributeKeyCollectionId = "collection_id"
	AttributeKeyOwnsTokens   = "owns_tokens"
	AttributeKeyTotalAmount  = "total_amount"
	AttributeKeyError        = "error"
	AttributeKeyBulkCount    = "bulk_count"
)

// FormatOwnershipQueryEvent formats an ownership query event
func FormatOwnershipQueryEvent(query *OwnershipQueryPacket) string {
	return fmt.Sprintf("ICQ ownership query: id=%s, address=%s, collection=%s",
		query.QueryId, query.Address, query.CollectionId)
}

// FormatOwnershipQueryResponseEvent formats an ownership query response event
func FormatOwnershipQueryResponseEvent(response *OwnershipQueryResponsePacket) string {
	if response.Error != "" {
		return fmt.Sprintf("ICQ ownership response: id=%s, error=%s",
			response.QueryId, response.Error)
	}
	return fmt.Sprintf("ICQ ownership response: id=%s, owns=%t, amount=%s",
		response.QueryId, response.OwnsTokens, response.TotalAmount.String())
}

// FormatFullBalanceQueryEvent formats a full balance query event
func FormatFullBalanceQueryEvent(query *FullBalanceQueryPacket) string {
	return fmt.Sprintf("ICQ full balance query: id=%s, address=%s, collection=%s",
		query.QueryId, query.Address, query.CollectionId)
}

// FormatFullBalanceQueryResponseEvent formats a full balance query response event
func FormatFullBalanceQueryResponseEvent(response *FullBalanceQueryResponsePacket) string {
	if response.Error != "" {
		return fmt.Sprintf("ICQ full balance response: id=%s, error=%s",
			response.QueryId, response.Error)
	}
	return fmt.Sprintf("ICQ full balance response: id=%s, balance_store_size=%d bytes",
		response.QueryId, len(response.BalanceStore))
}
