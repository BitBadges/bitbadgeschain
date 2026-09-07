package tokenization

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bound every nested wire field before protobuf decoding and semantic validation.
func meterJSONInput(ctx sdk.Context, input string) error {
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()
	var value func(int, string) error
	value = func(depth int, field string) error {
		if depth > 32 {
			return fmt.Errorf("message nesting exceeds 32 levels")
		}
		ctx.GasMeter().ConsumeGas(GasPerApprovalField, "precompile JSON value")
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		switch token := token.(type) {
		case string:
			return ValidateMetadataLength(token, field)
		case json.Delim:
			switch token {
			case '[':
				limit := MaxApprovalRanges
				if field == "addresses" {
					limit = MaxAddressListEntries
				}
				count := 0
				for decoder.More() {
					count++
					if err := ValidateArraySizeAllowEmpty(count, limit, field); err != nil {
						return err
					}
					if err := value(depth+1, field); err != nil {
						return err
					}
				}
				end, err := decoder.Token()
				if err != nil {
					return err
				}
				if end != json.Delim(']') {
					return fmt.Errorf("invalid array")
				}
			case '{':
				count := 0
				for decoder.More() {
					count++
					if count > 100 {
						return fmt.Errorf("object exceeds 100 fields")
					}
					key, err := decoder.Token()
					if err != nil {
						return err
					}
					name, ok := key.(string)
					if !ok {
						return fmt.Errorf("invalid object key")
					}
					if err := ValidateMetadataLength(name, "field name"); err != nil {
						return err
					}
					if err := value(depth+1, name); err != nil {
						return err
					}
				}
				end, err := decoder.Token()
				if err != nil {
					return err
				}
				if end != json.Delim('}') {
					return fmt.Errorf("invalid object")
				}
			default:
				return fmt.Errorf("unexpected JSON delimiter")
			}
		}
		return nil
	}
	if err := value(0, "message"); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return fmt.Errorf("unexpected trailing JSON input")
	}
	return nil
}
