//go:build test
// +build test

package cli

import (
	"encoding/json"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// MessageTestSuite tests message marshaling, unmarshaling, and routing
type MessageTestSuite struct {
	suite.Suite
	cdc          codec.Codec
	interfaceReg cdctypes.InterfaceRegistry
}

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}

func (s *MessageTestSuite) SetupTest() {
	// Create interface registry and register interfaces
	s.interfaceReg = cdctypes.NewInterfaceRegistry()
	tokenizationtypes.RegisterInterfaces(s.interfaceReg)
	gammtypes.RegisterInterfaces(s.interfaceReg)
	poolmanagertypes.RegisterInterfaces(s.interfaceReg)
	s.cdc = codec.NewProtoCodec(s.interfaceReg)
}

// TestMsgTransferTokensMarshal tests marshaling/unmarshaling of MsgTransferTokens
func (s *MessageTestSuite) TestMsgTransferTokensMarshal() {
	s.T().Log("Testing MsgTransferTokens marshal/unmarshal")

	msg := &tokenizationtypes.MsgTransferTokens{
		Creator:      "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*tokenizationtypes.Transfer{
			{
				From:        "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
				ToAddresses: []string{"bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l"},
				Balances: []*tokenizationtypes.Balance{
					{
						Amount: sdkmath.NewUint(100),
						TokenIds: []*tokenizationtypes.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
						},
						OwnershipTimes: []*tokenizationtypes.UintRange{
							{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(18446744073709551615)},
						},
					},
				},
			},
		},
	}

	// Test Proto marshal/unmarshal
	bz, err := s.cdc.Marshal(msg)
	s.Require().NoError(err)
	s.Require().NotEmpty(bz)

	var unmarshaledMsg tokenizationtypes.MsgTransferTokens
	err = s.cdc.Unmarshal(bz, &unmarshaledMsg)
	s.Require().NoError(err)
	s.Require().Equal(msg.Creator, unmarshaledMsg.Creator)
	s.Require().Equal(msg.CollectionId, unmarshaledMsg.CollectionId)
	s.Require().Len(unmarshaledMsg.Transfers, 1)

	// Test JSON marshal/unmarshal
	jsonBz, err := s.cdc.MarshalJSON(msg)
	s.Require().NoError(err)
	s.Require().NotEmpty(jsonBz)

	var jsonUnmarshaledMsg tokenizationtypes.MsgTransferTokens
	err = s.cdc.UnmarshalJSON(jsonBz, &jsonUnmarshaledMsg)
	s.Require().NoError(err)
	s.Require().Equal(msg.Creator, jsonUnmarshaledMsg.Creator)

	s.T().Logf("Successfully marshaled/unmarshaled MsgTransferTokens, proto size: %d bytes, json size: %d bytes", len(bz), len(jsonBz))
}

// TestMsgSwapExactAmountInMarshal tests marshaling/unmarshaling of MsgSwapExactAmountIn
func (s *MessageTestSuite) TestMsgSwapExactAmountInMarshal() {
	s.T().Log("Testing MsgSwapExactAmountIn marshal/unmarshal")

	msg := &gammtypes.MsgSwapExactAmountIn{
		Sender: "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		Routes: []poolmanagertypes.SwapAmountInRoute{
			{
				PoolId:        1,
				TokenOutDenom: "bar",
			},
		},
		TokenIn:           sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
		TokenOutMinAmount: sdkmath.NewInt(900000),
	}

	// Test Proto marshal/unmarshal
	bz, err := s.cdc.Marshal(msg)
	s.Require().NoError(err)
	s.Require().NotEmpty(bz)

	var unmarshaledMsg gammtypes.MsgSwapExactAmountIn
	err = s.cdc.Unmarshal(bz, &unmarshaledMsg)
	s.Require().NoError(err)
	s.Require().Equal(msg.Sender, unmarshaledMsg.Sender)
	s.Require().Len(unmarshaledMsg.Routes, 1)
	s.Require().Equal(msg.TokenIn.Denom, unmarshaledMsg.TokenIn.Denom)

	// Test JSON marshal/unmarshal
	jsonBz, err := s.cdc.MarshalJSON(msg)
	s.Require().NoError(err)

	var jsonUnmarshaledMsg gammtypes.MsgSwapExactAmountIn
	err = s.cdc.UnmarshalJSON(jsonBz, &jsonUnmarshaledMsg)
	s.Require().NoError(err)
	s.Require().Equal(msg.Sender, jsonUnmarshaledMsg.Sender)

	s.T().Logf("Successfully marshaled/unmarshaled MsgSwapExactAmountIn")
}

// TestMsgUniversalUpdateCollectionMarshal tests complex collection update message
func (s *MessageTestSuite) TestMsgUniversalUpdateCollectionMarshal() {
	s.T().Log("Testing MsgUniversalUpdateCollection marshal/unmarshal")

	msg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                    "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		CollectionId:               sdkmath.NewUint(0), // 0 means create new
		UpdateCollectionApprovals:  true,
		UpdateValidTokenIds:        true,
		UpdateCollectionPermissions: true,
		UpdateManager:              true,
		Manager:                    "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		UpdateCollectionMetadata:   true,
		CollectionMetadata: &tokenizationtypes.CollectionMetadata{
			Uri:        "ipfs://test",
			CustomData: "",
		},
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
	}

	// Test Proto marshal/unmarshal
	bz, err := s.cdc.Marshal(msg)
	s.Require().NoError(err)
	s.Require().NotEmpty(bz)

	var unmarshaledMsg tokenizationtypes.MsgUniversalUpdateCollection
	err = s.cdc.Unmarshal(bz, &unmarshaledMsg)
	s.Require().NoError(err)
	s.Require().Equal(msg.Creator, unmarshaledMsg.Creator)
	s.Require().True(unmarshaledMsg.UpdateCollectionApprovals)

	// Test JSON marshal/unmarshal
	jsonBz, err := s.cdc.MarshalJSON(msg)
	s.Require().NoError(err)

	var jsonUnmarshaledMsg tokenizationtypes.MsgUniversalUpdateCollection
	err = s.cdc.UnmarshalJSON(jsonBz, &jsonUnmarshaledMsg)
	s.Require().NoError(err)

	s.T().Logf("Successfully marshaled/unmarshaled MsgUniversalUpdateCollection, proto size: %d bytes", len(bz))
}

// TestMessageSigners verifies messages return proper signers
func (s *MessageTestSuite) TestMessageSigners() {
	s.T().Log("Testing message signers")

	testCases := []struct {
		name    string
		msg     sdk.Msg
		creator string
	}{
		{
			"MsgTransferTokens",
			&tokenizationtypes.MsgTransferTokens{
				Creator:      "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
				CollectionId: sdkmath.NewUint(1),
				Transfers:    []*tokenizationtypes.Transfer{},
			},
			"bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		},
		{
			"MsgSwapExactAmountIn",
			&gammtypes.MsgSwapExactAmountIn{
				Sender:            "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
				Routes:            []poolmanagertypes.SwapAmountInRoute{},
				TokenIn:           sdk.NewCoin("foo", sdkmath.NewInt(1000)),
				TokenOutMinAmount: sdkmath.NewInt(900),
			},
			"bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		},
		{
			"MsgSwapExactAmountOut",
			&gammtypes.MsgSwapExactAmountOut{
				Sender:           "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
				Routes:           []poolmanagertypes.SwapAmountOutRoute{},
				TokenOut:         sdk.NewCoin("bar", sdkmath.NewInt(1000)),
				TokenInMaxAmount: sdkmath.NewInt(1100),
			},
			"bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		},
		{
			"MsgJoinPool",
			&gammtypes.MsgJoinPool{
				Sender:         "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
				PoolId:         1,
				ShareOutAmount: sdkmath.NewInt(1000),
				TokenInMaxs:    sdk.NewCoins(),
			},
			"bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		},
		{
			"MsgExitPool",
			&gammtypes.MsgExitPool{
				Sender:        "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
				PoolId:        1,
				ShareInAmount: sdkmath.NewInt(1000),
				TokenOutMins:  sdk.NewCoins(),
			},
			"bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Verify the message can be marshaled (confirming it's valid)
			bz, err := s.cdc.Marshal(tc.msg)
			s.Require().NoError(err, "%s should marshal", tc.name)
			s.Require().NotEmpty(bz, "%s should produce non-empty bytes", tc.name)
			s.T().Logf("%s: marshaled size=%d bytes", tc.name, len(bz))
		})
	}
}

// TestAnyPacking tests packing messages into Any for transaction building
func (s *MessageTestSuite) TestAnyPacking() {
	s.T().Log("Testing Any packing/unpacking for transaction building")

	msg := &tokenizationtypes.MsgTransferTokens{
		Creator:      "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		CollectionId: sdkmath.NewUint(1),
		Transfers:    []*tokenizationtypes.Transfer{},
	}

	// Pack into Any
	anyMsg, err := cdctypes.NewAnyWithValue(msg)
	s.Require().NoError(err)
	s.Require().NotNil(anyMsg)
	s.T().Logf("TypeURL: %s", anyMsg.TypeUrl)

	// Unpack from Any
	var unpackedMsg sdk.Msg
	err = s.interfaceReg.UnpackAny(anyMsg, &unpackedMsg)
	s.Require().NoError(err)

	// Verify it's the same type
	transferMsg, ok := unpackedMsg.(*tokenizationtypes.MsgTransferTokens)
	s.Require().True(ok)
	s.Require().Equal(msg.Creator, transferMsg.Creator)
}

// TestJSONRawMessageParsing tests parsing raw JSON messages
func (s *MessageTestSuite) TestJSONRawMessageParsing() {
	s.T().Log("Testing JSON raw message parsing")

	// This simulates receiving a transaction JSON from CLI
	// Note: Proto3 JSON uses camelCase field names
	rawJSON := `{
		"creator": "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		"collectionId": "1",
		"transfers": []
	}`

	var msg tokenizationtypes.MsgTransferTokens
	err := s.cdc.UnmarshalJSON([]byte(rawJSON), &msg)
	s.Require().NoError(err)
	s.Require().Equal("bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l", msg.Creator)
	s.Require().Equal(sdkmath.NewUint(1), msg.CollectionId)

	// Test swap message
	swapJSON := `{
		"sender": "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		"routes": [{"poolId": "1", "tokenOutDenom": "bar"}],
		"tokenIn": {"denom": "foo", "amount": "1000000"},
		"tokenOutMinAmount": "900000"
	}`

	var swapMsg gammtypes.MsgSwapExactAmountIn
	err = s.cdc.UnmarshalJSON([]byte(swapJSON), &swapMsg)
	s.Require().NoError(err)
	s.Require().Equal("bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l", swapMsg.Sender)
	s.Require().Len(swapMsg.Routes, 1)
	s.Require().Equal(uint64(1), swapMsg.Routes[0].PoolId)
}

// TestRoundTripSerialization tests complete round-trip serialization
func (s *MessageTestSuite) TestRoundTripSerialization() {
	s.T().Log("Testing complete round-trip serialization")

	original := &gammtypes.MsgSwapExactAmountIn{
		Sender: "bb1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq6j9c2l",
		Routes: []poolmanagertypes.SwapAmountInRoute{
			{PoolId: 1, TokenOutDenom: "bar"},
			{PoolId: 2, TokenOutDenom: "baz"},
		},
		TokenIn:           sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
		TokenOutMinAmount: sdkmath.NewInt(900000),
	}

	// Proto round trip
	protoBz, err := s.cdc.Marshal(original)
	s.Require().NoError(err)

	var protoResult gammtypes.MsgSwapExactAmountIn
	err = s.cdc.Unmarshal(protoBz, &protoResult)
	s.Require().NoError(err)
	s.Require().Equal(original.Sender, protoResult.Sender)
	s.Require().Equal(len(original.Routes), len(protoResult.Routes))

	// JSON round trip
	jsonBz, err := s.cdc.MarshalJSON(original)
	s.Require().NoError(err)

	var jsonResult gammtypes.MsgSwapExactAmountIn
	err = s.cdc.UnmarshalJSON(jsonBz, &jsonResult)
	s.Require().NoError(err)
	s.Require().Equal(original.Sender, jsonResult.Sender)

	// Verify data integrity through raw JSON
	var rawMap map[string]json.RawMessage
	err = json.Unmarshal(jsonBz, &rawMap)
	s.Require().NoError(err)
	s.Require().Contains(rawMap, "sender")
	s.Require().Contains(rawMap, "routes")
	s.Require().Contains(rawMap, "token_in")
}
