package analyzer

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var ignoredEntityLemmas = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "by": true, "each": true,
	"for": true, "from": true, "in": true, "is": true, "of": true, "the": true,
	"to": true, "with": true,
	"amount": true, "date": true, "email": true, "name": true, "number": true,
	"phone": true, "title": true,
	"\u0430\u0442\u0440\u0438\u0431\u0443\u0442":                   true,
	"\u0430\u0442\u0440\u0438\u0431\u0443\u0442\u044b":             true,
	"\u0431\u0430\u0437\u0430":                                     true,
	"\u0431\u0434":                                                 true,
	"\u0434\u0430\u043d\u043d\u044b\u0435":                         true,
	"\u0434\u0430\u0442\u0430":                                     true,
	"\u0437\u0430\u043f\u0438\u0441\u044c":                         true,
	"\u0438\u043c\u044f":                                           true,
	"\u0438\u043d\u0444\u043e\u0440\u043c\u0430\u0446\u0438\u044f": true,
	"\u043a\u0430\u0436\u0434\u044b\u0439":                         true,
	"\u043a\u0430\u0436\u0434\u0430\u044f":                         true,
	"\u043a\u0430\u0436\u0434\u043e\u0435":                         true,
	"\u043d\u0435\u0441\u043a\u043e\u043b\u044c\u043a\u043e":       true,
	"\u043d\u043e\u043c\u0435\u0440":                               true,
	"\u043e\u043f\u0438\u0441\u0430\u043d\u0438\u0435":             true,
	"\u043f\u043e\u043b\u0435":                                     true,
	"\u043f\u043e\u043b\u044f":                                     true,
	"\u043f\u043e\u0447\u0442\u0430":                               true,
	"\u043f\u0440\u0435\u0434\u043c\u0435\u0442\u043d\u0430\u044f": true,
	"\u0441\u0438\u0441\u0442\u0435\u043c\u0430":                   true,
	"\u0441\u0442\u0440\u0443\u043a\u0442\u0443\u0440\u0430":       true,
	"\u0441\u0443\u043c\u043c\u0430":                               true,
	"\u0442\u0430\u0431\u043b\u0438\u0446\u0430":                   true,
	"\u0442\u0430\u0431\u043b\u0438\u0446\u044b":                   true,
	"\u0442\u0435\u043b\u0435\u0444\u043e\u043d":                   true,
	"\u0445\u0440\u0430\u043d\u0435\u043d\u0438\u0435":             true,
	"\u044d\u043a\u0437\u0435\u043c\u043f\u043b\u044f\u0440":       true,
}

var russianSingularReplacements = []struct {
	suffix string
	value  string
}{
	{"\u0438\u044f\u043c\u0438", "\u0438\u044f"},
	{"\u044f\u043c\u0438", "\u044f"},
	{"\u0430\u043c\u0438", "\u0430"},
	{"\u043e\u0433\u043e", ""},
	{"\u0435\u043c\u0443", ""},
	{"\u044b\u043c\u0438", ""},
	{"\u0438\u043c\u0438", ""},
	{"\u044b\u0435", ""},
	{"\u0438\u0435", "\u0438\u0435"},
	{"\u043e\u0432", ""},
	{"\u0435\u0432", ""},
	{"\u0435\u0439", "\u044c"},
	{"\u0430\u043c", "\u0430"},
	{"\u044f\u043c", "\u044f"},
	{"\u0430\u0445", "\u0430"},
	{"\u044f\u0445", "\u044f"},
	{"\u043e\u043c", ""},
	{"\u0435\u043c", "\u044c"},
	{"\u043e\u0439", "\u0430"},
	{"\u044b\u0445", ""},
	{"\u0438\u0445", ""},
	{"\u0443", ""},
	{"\u0435", ""},
	{"\u044b", ""},
	{"\u0438", ""},
}

var normalizedWordOverrides = map[string]string{
	"\u0434\u0430\u0442\u0443":       "\u0434\u0430\u0442\u0430",
	"\u0441\u0443\u043c\u043c\u0443": "\u0441\u0443\u043c\u043c\u0430",
}

func normalizeEntity(value string) string {
	value = normalizeWord(value)
	if value == "" {
		return value
	}

	if isASCII(value) && len(value) > 3 && strings.HasSuffix(value, "s") {
		value = strings.TrimSuffix(value, "s")
	}

	if !isASCII(value) {
		value = singularRussian(value)
	}

	first, size := utf8.DecodeRuneInString(value)
	return string(unicode.ToUpper(first)) + value[size:]
}

func normalizeAttribute(value string) string {
	return normalizeDomainWord(value)
}

func normalizeWord(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.Trim(value, ".,;:!?()[]{}\"'`")
	return value
}

func isIgnoredEntity(value string) bool {
	return ignoredEntityLemmas[normalizeDomainWord(value)]
}

func normalizeDomainWord(value string) string {
	value = normalizeWord(value)
	if override, exists := normalizedWordOverrides[value]; exists {
		return override
	}

	return value
}

func singularRussian(value string) string {
	for _, replacement := range russianSingularReplacements {
		if len([]rune(value)) > len([]rune(replacement.suffix))+2 && strings.HasSuffix(value, replacement.suffix) {
			return strings.TrimSuffix(value, replacement.suffix) + replacement.value
		}
	}
	return value
}

func isASCII(value string) bool {
	for _, r := range value {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func startsWithUpper(value string) bool {
	first, _ := utf8.DecodeRuneInString(value)
	return unicode.IsUpper(first)
}
