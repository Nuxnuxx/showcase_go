package services

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type (
	CustomValidator struct {
		Validator *validator.Validate
	}

	HumanErrors struct {
		Value 	 string
		Error 	 string
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func CreateHumanErrors(err error) map[string]HumanErrors {
	errors := make(map[string]HumanErrors)

	for _, v := range err.(validator.ValidationErrors) {
		error := strings.Builder{}
		error.WriteString(fmt.Sprintf("%s should be %s %s",
			strings.Split(v.Namespace(), ".")[1],
			v.Param(),
			v.Tag(),
		))

		errors[strings.ToLower(v.Field())] = HumanErrors{
			Value: v.Value().(string),
			Error: error.String(),
		}
	}

	return errors
}
