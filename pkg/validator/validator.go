package validator

import (
	"regexp"
	"strings"
)

func ValidateName(name string) bool {
	return len(strings.TrimSpace(name)) > 0
}

func ValidateSurname(surname string) bool {
	return len(strings.TrimSpace(surname)) > 0
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`/^([a-z0-9_\.-]+)@([\da-z\.-]+)\.([a-z\.]{2,6})$/g`)
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	return len(password) >= 6
}

func ValidateAdvertisement(title, description string, pictures []string, price float64) bool {
	return len(strings.TrimSpace(title)) > 0 &&
		len(strings.TrimSpace(description)) > 0 &&
		len(pictures) > 0 &&
		price > 0
}
