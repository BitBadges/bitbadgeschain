package osmocli

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

type dummyDesc struct {
	use       string
	overrides map[string]string
	parsers   map[string]CustomFieldParserFn
}

func (d dummyDesc) GetCustomFlagOverrides() map[string]string { return d.overrides }
func (d dummyDesc) GetCustomFieldParsers() map[string]CustomFieldParserFn {
	return d.parsers
}
func (d dummyDesc) GetUse() string          { return d.use }
func (d *dummyDesc) AttachToUse(str string) { d.use += str }

type dummyReq struct {
	Sender         string
	PoolId         uint64
	ShareOutAmount string
	TokenInMaxs    string
	TokenIn        string
}

func (dummyReq) Reset()         {}
func (dummyReq) String() string { return "" }
func (dummyReq) ProtoMessage()  {}

func noopParser(_ string, _ *pflag.FlagSet) (any, FieldReadLocation, error) {
	return nil, UsedFlag, nil
}

func TestAttachFieldsToUse(t *testing.T) {
	t.Run("skips sender, flag overrides, and custom parsers", func(t *testing.T) {
		desc := &dummyDesc{
			use: "join-pool",
			overrides: map[string]string{
				"poolid":         "pool-id",
				"ShareOutAmount": "share-amount-out",
			},
			parsers: map[string]CustomFieldParserFn{
				"TokenInMaxs": noopParser,
			},
		}
		attachFieldsToUse[*dummyReq](desc)
		require.Equal(t, "join-pool [token-in]", desc.use)
	})

	t.Run("leaves explicit Use with placeholders unchanged", func(t *testing.T) {
		desc := &dummyDesc{
			use: "swap-exact-amount-out [token-out] [token-in-max-amount]",
		}
		attachFieldsToUse[*dummyReq](desc)
		require.Equal(t, "swap-exact-amount-out [token-out] [token-in-max-amount]", desc.use)
	})
}
