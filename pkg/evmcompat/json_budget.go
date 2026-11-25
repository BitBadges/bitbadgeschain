package evmcompat

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultMaxJSONBytes        = 256 * 1024
	DefaultMaxJSONDepth        = 32
	DefaultMaxJSONTokens       = 10000
	DefaultMaxJSONObjectFields = 100
	DefaultMaxJSONArrayLength  = 100
	DefaultMaxJSONStringLength = 10000
	DefaultMaxJSONNumberLength = 128
)

// JSONBudgetOptions defines structural bounds for precompile JSON input validation.
type JSONBudgetOptions struct {
	MaxBytes        int
	MaxDepth        int
	MaxTokens       int
	MaxObjectFields int
	MaxArrayLength  int
	ArrayLimits     map[string]int
	MaxStringLength int
	MaxNumberLength int
	GasPerValue     uint64
}

// DefaultJSONBudgetOptions returns conservative structural bounds for precompile JSON inputs.
func DefaultJSONBudgetOptions() JSONBudgetOptions {
	return JSONBudgetOptions{
		MaxBytes:        DefaultMaxJSONBytes,
		MaxDepth:        DefaultMaxJSONDepth,
		MaxTokens:       DefaultMaxJSONTokens,
		MaxObjectFields: DefaultMaxJSONObjectFields,
		MaxArrayLength:  DefaultMaxJSONArrayLength,
		MaxStringLength: DefaultMaxJSONStringLength,
		MaxNumberLength: DefaultMaxJSONNumberLength,
		GasPerValue:     100,
	}
}

// ValidateJSONInput verifies structural limits on a JSON payload before decoding
// and bills per-value gas on ctx.GasMeter() when present.
func ValidateJSONInput(ctx sdk.Context, input string, opts JSONBudgetOptions) error {
	if opts.MaxBytes <= 0 {
		opts.MaxBytes = DefaultMaxJSONBytes
	}
	if opts.MaxDepth <= 0 {
		opts.MaxDepth = DefaultMaxJSONDepth
	}
	if opts.MaxTokens <= 0 {
		opts.MaxTokens = DefaultMaxJSONTokens
	}
	if opts.MaxObjectFields <= 0 {
		opts.MaxObjectFields = DefaultMaxJSONObjectFields
	}
	if opts.MaxArrayLength <= 0 {
		opts.MaxArrayLength = DefaultMaxJSONArrayLength
	}
	if opts.MaxStringLength <= 0 {
		opts.MaxStringLength = DefaultMaxJSONStringLength
	}
	if opts.MaxNumberLength <= 0 {
		opts.MaxNumberLength = DefaultMaxJSONNumberLength
	}

	if len(input) > opts.MaxBytes {
		return fmt.Errorf("JSON input size (%d) exceeds maximum allowed (%d)", len(input), opts.MaxBytes)
	}

	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()

	tokenCount := 0
	var consumeToken func() error
	consumeToken = func() error {
		tokenCount++
		if tokenCount > opts.MaxTokens {
			return fmt.Errorf("JSON token count exceeds maximum allowed (%d)", opts.MaxTokens)
		}
		if opts.GasPerValue > 0 && ctx.GasMeter() != nil {
			ctx.GasMeter().ConsumeGas(opts.GasPerValue, "precompile JSON value")
		}
		return nil
	}

	var value func(int, string) error
	value = func(depth int, field string) error {
		if depth > opts.MaxDepth {
			return fmt.Errorf("message nesting exceeds %d levels", opts.MaxDepth)
		}
		if err := consumeToken(); err != nil {
			return err
		}
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		switch tok := token.(type) {
		case string:
			if len(tok) > opts.MaxStringLength {
				return fmt.Errorf("%s length (%d) exceeds maximum allowed length (%d)", field, len(tok), opts.MaxStringLength)
			}
		case json.Number:
			if len(string(tok)) > opts.MaxNumberLength {
				return fmt.Errorf("%s numeric literal length (%d) exceeds maximum allowed (%d)", field, len(string(tok)), opts.MaxNumberLength)
			}
		case json.Delim:
			switch tok {
			case '[':
				limit := opts.MaxArrayLength
				if fieldLimit, ok := opts.ArrayLimits[field]; ok {
					limit = fieldLimit
				}
				count := 0
				for decoder.More() {
					count++
					if count > limit {
						return fmt.Errorf("%s size (%d) exceeds maximum allowed size (%d)", field, count, limit)
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
					if count > opts.MaxObjectFields {
						return fmt.Errorf("object exceeds %d fields", opts.MaxObjectFields)
					}
					key, err := decoder.Token()
					if err != nil {
						return err
					}
					name, ok := key.(string)
					if !ok {
						return fmt.Errorf("invalid object key")
					}
					if len(name) > opts.MaxStringLength {
						return fmt.Errorf("field name length (%d) exceeds maximum allowed length (%d)", len(name), opts.MaxStringLength)
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
