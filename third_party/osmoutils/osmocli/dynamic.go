package osmocli

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/cosmos/gogoproto/proto"

	"github.com/bitbadges/bitbadgeschain/third_party/osmoutils"
)

type Descriptor interface {
	GetCustomFlagOverrides() map[string]string
	GetCustomFieldParsers() map[string]CustomFieldParserFn
	GetUse() string
	AttachToUse(str string)
}

// fields that are not provided as arguments
var nonAttachableFields []string = []string{"sender", "pagination", "owner", "admin"}

// attachFieldsToUse extracts fields from reqP proto message and dynamically appends them into Use field.
// Fields supplied via CustomFlagOverrides or CustomFieldParsers (typically FlagOnlyParser) are omitted
// so Usage matches the positional args cobra actually accepts.
// If Use already contains a `[` placeholder, it is treated as an explicit Use string and left unchanged.
func attachFieldsToUse[reqP proto.Message](desc Descriptor) {
	if strings.Contains(desc.GetUse(), "[") {
		return
	}

	req := osmoutils.MakeNew[reqP]()
	v := reflect.ValueOf(req).Type().Elem() // get underlying non-pointer struct
	var useField string
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Field(i).Name
		fn := pascalToKebab(fieldName)

		if osmoutils.Contains(nonAttachableFields, fn) {
			continue
		}
		if isFlagProvidedField(desc, fieldName) {
			continue
		}

		useField += fmt.Sprintf(" [%s]", fn)
	}

	desc.AttachToUse(useField)
}

// isFlagProvidedField reports whether fieldName is populated from a flag rather than a positional arg.
// CustomFlagOverrides keys are matched case-insensitively as both the raw struct field name and kebab-case,
// matching ParseField's lowercased lookup after Sanitize().
func isFlagProvidedField(desc Descriptor, fieldName string) bool {
	lower := strings.ToLower(fieldName)
	kebab := pascalToKebab(fieldName)
	for k, v := range desc.GetCustomFlagOverrides() {
		if v == "" {
			continue
		}
		kl := strings.ToLower(k)
		if kl == lower || k == kebab || kl == kebab {
			return true
		}
	}
	for k := range desc.GetCustomFieldParsers() {
		if strings.ToLower(k) == lower {
			return true
		}
	}
	return false
}

// pascalToKebab converts PascalCase string to kebab-case string
func pascalToKebab(s string) string {
	reg := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	s = reg.ReplaceAllString(s, `${1}-${2}`)

	// Convert everything to lowercase
	return strings.ToLower(s)
}
