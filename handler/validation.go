package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

func Validate(dto any) (bool, []string) {
	validate := validator.New()
	validate.RegisterValidation("notblank", validators.NotBlank)
	var fields []string;
	if err := validate.Struct(dto); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fields = append(fields, e.Field())
			fmt.Printf("Field: %s, Error: %s\n", e.Field(), e.Tag())
		}
		fmt.Println("Invalid struct!", err.Error())
		return false, fields
	}
	return true, fields
}