package api

import (
	"strings"
	"unicode"

	validator "github.com/darnellsylvain/auth52/internal"
)


func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(v.Matches(email, validator.EmailRX), "email", "must be a valid email")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 characters long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 characters long")

	v.Check(strings.IndexFunc(password, unicode.IsUpper) >= 0, "password", "must contain an uppercase letter")
	v.Check(strings.IndexFunc(password, unicode.IsLower) >= 0, "password", "must contain a lowercase letter")
	v.Check(strings.IndexFunc(password, unicode.IsDigit) >= 0, "password", "must contain a digit")
	v.Check(strings.ContainsAny(password, "!@#$%^&*"), "password", "must contain a special character (!@#$%^&*)")
}

func ValidateSignupParams(v *validator.Validator, params *SignupParams) {
	ValidateEmail(v, params.Email)
	ValidatePassword(v, params.Password)
}