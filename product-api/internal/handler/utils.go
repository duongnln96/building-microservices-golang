package handler

import (
	"net/http"
	"strconv"

	protos "github.com/duongnln96/building-microservices-golang/currency/protos/currency"
	"github.com/duongnln96/building-microservices-golang/product-api/internal/utils"
	"github.com/labstack/echo/v4"
)

// getProductIDParam returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func (p *productHandler) getProductIDParam(c echo.Context) int {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		p.log.Panicf("Cannot convert param string %+v", err)
	}

	return id
}

func (p *productHandler) getProductQuery(c echo.Context) string {
	currency := c.QueryParam("currency")
	return currency
}

func (p *productHandler) getRate(currency string) (float64, error) {
	rateRequest := protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[currency]),
	}

	rate, err := p.cc.GetRate(p.ctx, &rateRequest)
	return float64(rate.Rate), err
}

func (p *productHandler) responseData(c echo.Context, i interface{}) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return utils.ToJSON(i, c.Response())
}

func (p *productHandler) responseStatusOK(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (p *productHandler) responseErrCode(c echo.Context, statuscode int, msg string) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(statuscode)
	return utils.ToJSON(&GenericError{Message: msg}, c.Response())
}
