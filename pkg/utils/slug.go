package utils

import (
	"fmt"
	"strings"
	"unicode"
)

func Slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	lastDash := false

	for _, r := range value {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			builder.WriteRune(r)
			lastDash = false
		case !lastDash:
			builder.WriteByte('-')
			lastDash = true
		}
	}

	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return "item"
	}
	return result
}

func UniqueSlug(base string, exists func(string) (bool, error)) (string, error) {
	base = Slugify(base)
	for suffix := 1; suffix <= 100; suffix++ {
		candidate := base
		if suffix > 1 {
			candidate = fmt.Sprintf("%s-%d", base, suffix)
		}
		found, err := exists(candidate)
		if err != nil {
			return "", err
		}
		if !found {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not generate a unique slug")
}
