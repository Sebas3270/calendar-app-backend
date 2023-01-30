package helpers

import (
	"github.com/Sebas3270/calendar-app-backend/models"
	"github.com/go-playground/validator"
)

var validate = validator.New()

type ErrorResponse struct {
	Field   string
	Message string
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return ""
}

func ValidateUserStruct(user models.User) []*ErrorResponse {

	var errors []*ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Message = msgForTag(err.Tag())
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateLogInStruct(logIn models.LogIn) []*ErrorResponse {

	var errors []*ErrorResponse
	err := validate.Struct(logIn)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Message = msgForTag(err.Tag())
			errors = append(errors, &element)
		}
	}
	return errors
}

func ValidateEventStruct(event models.Event) []*ErrorResponse {

	var errors []*ErrorResponse
	err := validate.Struct(event)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Message = msgForTag(err.Tag())
			errors = append(errors, &element)
		}
	}
	return errors
}
