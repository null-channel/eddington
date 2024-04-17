package types

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type CreateUserRequest struct {
	UserId int64      `json:"userId"`
	Traits UserTraits `json:"traits"`
}

type UserTraits struct {
	Email             string `json:"email"`
	Name              string `json:"name"`
	NewsLetterConsent bool   `json:"newsletterConsent"`
}

var validate *validator.Validate

func ValidateAgeGT(fl validator.FieldLevel) bool {
	ageLimiter, err := strconv.Atoi(fl.Param())
	if err != nil {
		fmt.Printf("Error using age parameter: %v\n", err)
		return false
	}
	dob := fl.Field().Interface().(time.Time)
	minDate := time.Now().AddDate(-1*ageLimiter, 0, 0)
	return dob.Before(minDate)
}

type NewUserRequest struct {
	ID                string    `json:"id"`
	Name              string    `json:"name" validate:"required"`
	Email             string    `json:"email" validate:"email"`
	NewsletterConsent bool      `json:"newsletterConsent" validate:"boolean"`
	DOB               time.Time `json:"dob" validate:"age=16"`
}

// ApiError represents an API error containing parameter and message details.
type ApiError struct {
	Param   string `json:"param"`
	Message string `json:"message"`
}

func ConstructErrorMeesages(err error) []ApiError {
	var apiErrors []ApiError
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.StructField()
		var message string

		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", fieldName)
		case "email":
			message = fmt.Sprintf("%s must be a valid email address", fieldName)
		case "age":
			message = fmt.Sprintf("%s must be greater than %s", fieldName, err.Param())
		default:
			message = fmt.Sprintf("%s is invalid", fieldName)
		}

		apiError := ApiError{
			Param:   fieldName,
			Message: message,
		}

		apiErrors = append(apiErrors, apiError)
	}
	return apiErrors
}
func (newUserReq *NewUserRequest) Validate() error {
	validate = validator.New()
	validate.RegisterValidation("age", ValidateAgeGT)
	err := validate.Struct(newUserReq)
	return err
}
