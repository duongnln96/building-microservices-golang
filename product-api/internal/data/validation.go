package data

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// Validation contains
type ProductValidation struct {
	validator *validator.Validate
}

func NewValidation() *ProductValidation {
	return &ProductValidation{
		validator: validator.New(),
	}
}

// Validate the product request data
func (v *ProductValidation) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
