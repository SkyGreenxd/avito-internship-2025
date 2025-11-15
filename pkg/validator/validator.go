package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	userIDRegex = regexp.MustCompile(`^u([1-9]|[1-9][0-9]|[1-9][0-9]{2})$`)
	prIDRegex   = regexp.MustCompile(`^pr-(100[1-9]|10[1-9][0-9]|1[1-9][0-9]{2}|[2-9][0-9]{3})$`)
)

// ValidateUserID валидация для user_id
func ValidateUserID(fl validator.FieldLevel) bool {
	return userIDRegex.MatchString(fl.Field().String())
}

// ValidatePullRequestID валидация для pull_request_id
func ValidatePullRequestID(fl validator.FieldLevel) bool {
	return prIDRegex.MatchString(fl.Field().String())
}
