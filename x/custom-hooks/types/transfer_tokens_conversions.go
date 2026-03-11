package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// UnmarshalTransfersFromJSON takes the raw JSON "transfers" array from the memo (snake_case keys)
// and unmarshals it directly into protobuf Transfer types (which use camelCase JSON tags).
// This avoids needing intermediate types — just re-key and unmarshal.
func UnmarshalTransfersFromJSON(transfersJSON json.RawMessage) ([]*tokenizationtypes.Transfer, error) {
	// Parse into generic structure
	var raw []interface{}
	if err := json.Unmarshal(transfersJSON, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse transfers JSON: %w", err)
	}

	// Re-key all snake_case keys to camelCase recursively
	camelCased := make([]interface{}, len(raw))
	for i, item := range raw {
		camelCased[i] = snakeToCamelKeys(item)
	}

	// Marshal back to JSON with camelCase keys
	camelJSON, err := json.Marshal(camelCased)
	if err != nil {
		return nil, fmt.Errorf("failed to re-marshal transfers: %w", err)
	}

	// Unmarshal into protobuf types (which have camelCase json tags)
	var transfers []*tokenizationtypes.Transfer
	if err := json.Unmarshal(camelJSON, &transfers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into proto transfers: %w", err)
	}

	return transfers, nil
}

// snakeToCamelKeys recursively converts all map keys from snake_case to camelCase.
func snakeToCamelKeys(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(val))
		for k, v := range val {
			result[snakeToCamel(k)] = snakeToCamelKeys(v)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = snakeToCamelKeys(item)
		}
		return result
	default:
		return v
	}
}

// snakeToCamel converts a snake_case string to camelCase.
// e.g. "to_addresses" → "toAddresses", "only_check_prioritized_collection_approvals" → "onlyCheckPrioritizedCollectionApprovals"
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) <= 1 {
		return s
	}
	var b strings.Builder
	b.WriteString(parts[0])
	for _, part := range parts[1:] {
		if len(part) == 0 {
			continue
		}
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}
	return b.String()
}
