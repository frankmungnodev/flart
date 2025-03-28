package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToPascalCase(str string) string {
	// Split the string by non-alphanumeric characters
	words := strings.FieldsFunc(str, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	// Capitalize first letter of each word
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[0:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, "")
}
