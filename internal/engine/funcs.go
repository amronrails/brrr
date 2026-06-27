package engine

import (
	"strings"
	"text/template"
	"unicode"
)

// FuncMap returns the template helper functions made available to every
// template rendered by the engine. They are intentionally small and
// dependency-free so the generated naming is predictable across Go, SQL and
// TypeScript output.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"pascal":   Pascal,
		"camel":    Camel,
		"snake":    Snake,
		"kebab":    Kebab,
		"screamingSnake": func(s string) string { return strings.ToUpper(Snake(s)) },
		"lower":    strings.ToLower,
		"upper":    strings.ToUpper,
		"title":    Title,
		"plural":   Plural,
		"singular": Singular,
		"add":      func(a, b int) int { return a + b },
		// table/route helpers built from the above
		"tableName": func(s string) string { return Plural(Snake(s)) },
		"routePath": func(s string) string { return Plural(Kebab(s)) },
	}
}

// words splits an identifier written in any common case style into its lower
// cased component words. It is the basis for every other case conversion.
func words(s string) []string {
	var out []string
	var cur strings.Builder
	runes := []rune(s)
	flush := func() {
		if cur.Len() > 0 {
			out = append(out, strings.ToLower(cur.String()))
			cur.Reset()
		}
	}
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case r == '_' || r == '-' || r == ' ' || r == '.' || r == '/':
			flush()
		case unicode.IsUpper(r):
			// boundary before an uppercase run start, and before the last
			// uppercase of an acronym that precedes a lowercase (e.g. "HTTPServer").
			prevLower := i > 0 && unicode.IsLower(runes[i-1])
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if prevLower || (cur.Len() > 0 && nextLower) {
				flush()
			}
			cur.WriteRune(r)
		default:
			cur.WriteRune(r)
		}
	}
	flush()
	return out
}

// Pascal converts an identifier to PascalCase ("user_profile" -> "UserProfile").
func Pascal(s string) string {
	var b strings.Builder
	for _, w := range words(s) {
		b.WriteString(capitalize(w))
	}
	return b.String()
}

// Camel converts an identifier to camelCase ("user_profile" -> "userProfile").
func Camel(s string) string {
	ws := words(s)
	var b strings.Builder
	for i, w := range ws {
		if i == 0 {
			b.WriteString(w)
			continue
		}
		b.WriteString(capitalize(w))
	}
	return b.String()
}

// Snake converts an identifier to snake_case.
func Snake(s string) string { return strings.Join(words(s), "_") }

// Kebab converts an identifier to kebab-case.
func Kebab(s string) string { return strings.Join(words(s), "-") }

// Title converts an identifier to space separated Title Case ("user_profile" ->
// "User Profile"), useful for UI labels.
func Title(s string) string {
	ws := words(s)
	for i, w := range ws {
		ws[i] = capitalize(w)
	}
	return strings.Join(ws, " ")
}

func capitalize(w string) string {
	if w == "" {
		return ""
	}
	// Preserve common all-caps acronyms.
	if acro, ok := acronyms[w]; ok {
		return acro
	}
	r := []rune(w)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

var acronyms = map[string]string{
	"id": "ID", "url": "URL", "api": "API", "http": "HTTP", "json": "JSON",
	"jwt": "JWT", "uuid": "UUID", "sql": "SQL", "html": "HTML", "ui": "UI",
}

// Plural returns a naive English plural. It covers the cases that matter for
// table and route naming; irregular nouns can be overridden by callers if
// needed later.
func Plural(s string) string {
	if s == "" {
		return s
	}
	lower := strings.ToLower(s)
	for sing, plur := range irregulars {
		if lower == sing {
			return matchCase(s, plur)
		}
	}
	switch {
	case endsWithAny(lower, "s", "x", "z", "ch", "sh"):
		return s + "es"
	case strings.HasSuffix(lower, "y") && !isVowel(beforeLast(lower)):
		return s[:len(s)-1] + "ies"
	default:
		return s + "s"
	}
}

// Singular returns a naive English singular form.
func Singular(s string) string {
	if s == "" {
		return s
	}
	lower := strings.ToLower(s)
	for sing, plur := range irregulars {
		if lower == plur {
			return matchCase(s, sing)
		}
	}
	switch {
	case strings.HasSuffix(lower, "ies") && len(s) > 3:
		return s[:len(s)-3] + "y"
	case strings.HasSuffix(lower, "ses") || strings.HasSuffix(lower, "xes") ||
		strings.HasSuffix(lower, "zes") || strings.HasSuffix(lower, "ches") ||
		strings.HasSuffix(lower, "shes"):
		return s[:len(s)-2]
	case strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss"):
		return s[:len(s)-1]
	default:
		return s
	}
}

var irregulars = map[string]string{
	"person": "people", "man": "men", "woman": "women", "child": "children",
	"tooth": "teeth", "foot": "feet", "mouse": "mice", "goose": "geese",
}

func endsWithAny(s string, suffixes ...string) bool {
	for _, suf := range suffixes {
		if strings.HasSuffix(s, suf) {
			return true
		}
	}
	return false
}

func isVowel(r rune) bool { return strings.ContainsRune("aeiou", r) }

func beforeLast(s string) rune {
	r := []rune(s)
	if len(r) < 2 {
		return 0
	}
	return r[len(r)-2]
}

// matchCase applies the leading-capitalisation of src to dst.
func matchCase(src, dst string) string {
	if src == "" {
		return dst
	}
	if unicode.IsUpper([]rune(src)[0]) {
		return capitalize(dst)
	}
	return dst
}
