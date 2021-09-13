package middleware

import (
	"fmt"
	"net/http"

	"github.com/duongnln96/building-microservices-golang/product-api/config"
	"github.com/duongnln96/building-microservices-golang/product-api/src/service"
	"github.com/duongnln96/building-microservices-golang/product-api/tools/helper"
	"github.com/labstack/echo/v4"
)

func AuthorizeJWT(svc service.JWTSeriveI) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !config.GetConfig().Server.UseJwt {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				res := helper.BuildErrorResponse("Fail", fmt.Errorf("No token found").Error(), nil)
				return c.JSON(http.StatusBadRequest, res)
			}

			token, err := svc.ValidateToken(authHeader)
			if token.Valid {
				return next(c)
			}
			res := helper.BuildErrorResponse("Fail", fmt.Errorf("No token found %+v", err).Error(), nil)
			return c.JSON(http.StatusUnauthorized, res)
		}
	}
}
