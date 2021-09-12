package middleware

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type UserValidationI interface {
	Validate(interface{})
}

type UserValidation struct {
	validator *validator.Validate
}

func NewValidation() *UserValidation {
	return &UserValidation{
		validator: validator.New(),
	}
}

func (v *UserValidation) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
