package validation

import (
	"net/url"
	"regexp"
)

func GetMissingFields(form url.Values, fields []string) []string {
	missing := []string{}

	for _, field := range fields {
		if !form.Has(field) {
			missing = append(missing, field)
		}
	}

	return missing
}

var MATCH_LOWERCASE = regexp.MustCompile(`[a-z]`)
var MATCH_UPPERCASE = regexp.MustCompile(`[A-Z]`)
var MATCH_NUMBER = regexp.MustCompile(`[0-9]`)
var MATCH_SPECIAL = regexp.MustCompile(`[$&+=?@#|<>^*()%!-]`)

func CheckPasswordComplexity(password string) bool {
	if len(password) < 8 {
		return false
	}

	if len(password) > 50 {
		return false
	}

	bytes := []byte(password)

	if !MATCH_LOWERCASE.Match(bytes) {
		return false
	}

	if !MATCH_UPPERCASE.Match(bytes) {
		return false
	}

	if !MATCH_NUMBER.Match(bytes) {
		return false
	}

	if !MATCH_SPECIAL.Match(bytes) {
		return false
	}

	return true
}
