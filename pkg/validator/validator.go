package validator

import (
	"avito-internship/pkg/e"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"sync"
)

var (
	userIDRegex = regexp.MustCompile(`^u([1-9]|[1-9][0-9]|[1-9][0-9]{2})$`)
	prIDRegex   = regexp.MustCompile(`^pr-(100[1-9]|10[1-9][0-9]|1[1-9][0-9]{2}|[2-9][0-9]{3})$`)
	vOnce       sync.Once
)

// ValidateUserID валидация для user_id
func ValidateUserID(fl validator.FieldLevel) bool {
	return userIDRegex.MatchString(fl.Field().String())
}

// ValidatePullRequestID валидация для pull_request_id
func ValidatePullRequestID(fl validator.FieldLevel) bool {
	return prIDRegex.MatchString(fl.Field().String())
}

func RegisterValidators() error {
	const op = "validator.RegisterValidators"

	var result_err error
	vOnce.Do(func() {
		validator, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			result_err = e.Wrap(op+"binding.Validator.Engine() is not *validator.Validate", e.ErrValidatorFailed)
			return
		}

		if err := validator.RegisterValidation("userid", ValidateUserID); err != nil {
			result_err = e.Wrap(op, err)
			return
		}

		if err := validator.RegisterValidation("prid", ValidatePullRequestID); err != nil {
			result_err = e.Wrap(op, err)
			return
		}
	})

	return result_err
}
