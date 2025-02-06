package utils

import "regexp"

func ValidateEmail(email string) bool {

	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if !regex.MatchString(email) {
		return false
	}

	return true
}

func ValidatePassword(password string) bool {

	regex := regexp.MustCompile(`^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$`)

	if !regex.MatchString(password) {
		return false
	}

	return true
}
