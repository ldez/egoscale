package helpers

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	abbr "github.com/BluntSporks/abbreviation"
)

var uppercaseAcronym = sync.Map{}

// ConfigureAcronym allows you to add additional words which will be considered acronyms
func ConfigureAcronym(key, val string) {
	uppercaseAcronym.Store(key, val)
}

// RenderReference renders OpenAPI reference from path to go style.
func RenderReference(referencePath string) string {
	return ToCamel(filepath.Base(referencePath))
}

// Header retruns header file for generated go source files.
func Header(packageName, version string) []byte {
	return []byte(fmt.Sprintf(`// Package %s provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/egoscale/v3/generator version %s DO NOT EDIT.
`,
		packageName, version,
	))
}

// ToLowerCamel converts a string to lowerCamelCase
func ToLowerCamel(s string) string {
	return toInitialCamel(s, true)
}

// ToCamel converts a string to CamelCase
func ToCamel(s string) string {
	return toInitialCamel(s, false)
}

// toInitialCamel got inspiration from https://github.com/iancoleman/strcase
// with improvement on acronym conversion.
func toInitialCamel(s string, lower bool) string {
	if s == "" {
		return ""
	}

	pattern := `[-_./\s]`
	s = trimBySeparators(s, pattern)
	words := splitBySeparators(s, pattern)

	for i, w := range words {
		if w == "" {
			continue
		}

		_, ok := abbr.Acronyms[strings.ToUpper(w)]
		if ok {
			words[i] = strings.ToUpper(w)
		}

		a, hasAcronym := uppercaseAcronym.Load(w)
		if hasAcronym {
			words[i] = a.(string)
		}

		if i == 0 && lower {
			words[i] = strings.ToLower(words[i])
			continue
		}

		bytes := []byte(words[i])
		v := bytes[0]
		vIsLow := v >= 'a' && v <= 'z'
		if vIsLow {
			bytes[0] += 'A'
			bytes[0] -= 'a'
		}

		words[i] = string(bytes)
	}

	return strings.Join(words, "")
}

func trimBySeparators(s, separators string) string {
	trimPattern := fmt.Sprintf("^%s+|%s+$", separators, separators)
	re := regexp.MustCompile(trimPattern)
	return re.ReplaceAllString(s, "")
}

func splitBySeparators(s, separators string) []string {
	re := regexp.MustCompile(separators)
	parts := re.Split(s, -1)
	var result []string
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

// RenderDoc returns proper go doc comment from
// an OpenAPI spec field documentation.
func RenderDoc(doc string) string {
	if doc == "null" {
		return ""
	}

	docs := strings.Split(doc, "\n")
	r := []string{}
	for i, d := range docs {
		if d == "" {
			docs = append(docs[:i], docs[i+1:]...)
			continue
		}
		r = append(r, "// "+strings.TrimSpace(d))
	}

	if len(r) == 0 {
		return ""
	}

	return strings.Join(r, "\n")
}
