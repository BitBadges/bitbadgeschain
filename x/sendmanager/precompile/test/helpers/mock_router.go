package helpers

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/sendmanager/types"
)

// MockAliasDenomRouter is a mock implementation of AliasDenomRouter for testing
type MockAliasDenomRouter struct {
	prefix    string
	sendCalls []SendCall
}

type SendCall struct {
	From   string
	To     string
	Denom  string
	Amount sdkmath.Uint
}

// Ensure MockAliasDenomRouter implements AliasDenomRouter
var _ types.AliasDenomRouter = (*MockAliasDenomRouter)(nil)

// NewMockRouter creates a new mock router for testing
func NewMockRouter(prefix string) *MockAliasDenomRouter {
	return &MockAliasDenomRouter{
		prefix:    prefix,
		sendCalls: []SendCall{},
	}
}

func (m *MockAliasDenomRouter) CheckIsAliasDenom(ctx sdk.Context, denom string) bool {
	return len(denom) > len(m.prefix) && denom[:len(m.prefix)] == m.prefix
}

func (m *MockAliasDenomRouter) SendNativeTokensViaAliasDenom(ctx sdk.Context, recipientAddress, toAddress, denom string, amount sdkmath.Uint) error {
	m.sendCalls = append(m.sendCalls, SendCall{
		From:   recipientAddress,
		To:     toAddress,
		Denom:  denom,
		Amount: amount,
	})
	return nil
}

func (m *MockAliasDenomRouter) FundCommunityPoolViaAliasDenom(ctx sdk.Context, fromAddress, toAddress, denom string, amount sdkmath.Uint) error {
	return nil
}

func (m *MockAliasDenomRouter) SpendFromCommunityPoolViaAliasDenom(ctx sdk.Context, fromAddress, toAddress, denom string, amount sdkmath.Uint) error {
	return nil
}

func (m *MockAliasDenomRouter) SendFromModuleToAccountViaAliasDenom(ctx sdk.Context, moduleAddress, toAddress, denom string, amount sdkmath.Uint) error {
	return nil
}

func (m *MockAliasDenomRouter) SendFromAccountToModuleViaAliasDenom(ctx sdk.Context, fromAddress, moduleAddress, denom string, amount sdkmath.Uint) error {
	return nil
}

func (m *MockAliasDenomRouter) GetBalanceWithAliasRouting(ctx sdk.Context, address sdk.AccAddress, denom string) (sdk.Coin, error) {
	return sdk.NewCoin(denom, sdkmath.ZeroInt()), nil
}

// GetSendCalls returns all send calls made to this router
func (m *MockAliasDenomRouter) GetSendCalls() []SendCall {
	return m.sendCalls
}

// ResetSendCalls clears the send call history
func (m *MockAliasDenomRouter) ResetSendCalls() {
	m.sendCalls = []SendCall{}
}

