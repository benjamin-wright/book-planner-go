package validation

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
		return false
	}

	if len(password) > 50 {
		return false
	}

	if !LOWERCASE.Intersects(password) {
		return false
	}

	if !UPPERCASE.Intersects(password) {
		return false
	}

	if !NUMBER.Intersects(password) {
		return false
	}

	if !SPECIAL.Intersects(password) {
		return false
	}

	return true
}
