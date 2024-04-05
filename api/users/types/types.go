package types

import (
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

var validate = validator.New()

type DOB struct {
	time.Time
}

func (t *DOB) UnmarshalJSON(b []byte) error {
	parseTime, err := time.Parse(`"2006-01-02T15:04:05.000z"`, string(b))
	if err != nil {
		return err
	}
	t.Time = parseTime
	return nil
}
func ageValidation(fl validator.FieldLevel) bool {
	ageLimiter := fl.Param()
	limit, err := time.ParseDuration("-" + ageLimiter + "y")
	if err != nil {
		return false
	}
	dob := fl.Field().Interface().(DOB)
	minDate := time.Now().Add(limit)
	return dob.Before(minDate)
}

type NewUserRequest struct {
	ID                string `json:"id" validator:"required"`
	Name              string `json:"name" validator:"required"`
	Email             string `json:"email" validator:"email"`
	NewsletterConsent bool   `json:"newsletterConsent" validator:"bool"`
	DOB               DOB    `json:"dob" validator:"age=16"`
}

func (newUserReq *NewUserRequest) Validate() error {
	validate.RegisterValidation("age", ageValidation)
	err := validate.Struct(newUserReq)
	return err
}
