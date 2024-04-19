package types

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/null-channel/eddington/api/app/models"
	userTypes "github.com/null-channel/eddington/api/users/types"
)

var validate *validator.Validate

type Application struct {
	Name          string           `json:"name" validate:"required"`
	Image         string           `json:"image"`
	GitRepo       string           `json:"gitRepo" validate:"required"`
	RepoType      models.BuildType `json:"repoType" validate:"required,repoType"`
	ResourceGroup string           `json:"resourceGroup"`
	Directory     string           `json:"directory" validate:"required"`
}

func ConstructErrorMessages(err error) []userTypes.ApiError {
	var apiErrors []userTypes.ApiError
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.StructField()
		var message string

		switch err.Tag() {
		case "required":
			message = fmt.Sprintf("%s is required", fieldName)
		case "repoType":
			message = fmt.Sprintf("%s unkown repo type.Please choose on of the following ['NodeJS', 'Go', 'Rust', 'Python', 'ContainerImage'].", err.Value())
		default:
			message = fmt.Sprintf("%s is invalid", fieldName)
		}
		apiError := userTypes.ApiError{
			Param:   fieldName,
			Message: message,
		}

		apiErrors = append(apiErrors, apiError)
	}
	return apiErrors
}
func validateRepoType(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(models.BuildType)
	return value.String() != "unknown"
}

func (application *Application) Validate() error {
	validate = validator.New()
	validate.RegisterValidation("repoType", validateRepoType)
	err := validate.Struct(application)
	return err
}
