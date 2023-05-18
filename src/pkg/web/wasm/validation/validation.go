package validation

import "go.uber.org/zap"

type Match string

var LOWERCASE = Match("abcdefghijklmnopqrstuvwxyz")
var UPPERCASE = Match("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var NUMBER = Match(`0123456789`)
var SPECIAL = Match(`$&+=?@#|<>^*()%!-`)

func (m *Match) Intersects(value string) bool {
	original := string(*m)

	for _, x := range original {
		for _, y := range value {
			if x == y {
				return true
			}
		}
	}

	return false
}

func CheckPasswordComplexity(password string) bool {
	if len(password) < 8 {
		zap.S().Debugf("password too short")
		return false
	}

	if len(password) > 50 {
		zap.S().Debugf("password too long")
		return false
	}

	if !LOWERCASE.Intersects(password) {
		zap.S().Debugf("password missing lowercase characters")
		return false
	}

	if !UPPERCASE.Intersects(password) {
		zap.S().Debugf("password missing uppercase characters")
		return false
	}

	if !NUMBER.Intersects(password) {
		zap.S().Debugf("password missing numeric characters")
		return false
	}

	if !SPECIAL.Intersects(password) {
		zap.S().Debugf("password missing special characters")
		return false
	}

	return true
}
