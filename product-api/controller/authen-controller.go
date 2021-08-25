package controller

import (
	"fmt"
	"net/http"

	"github.com/duongnln96/building-microservices-golang/product-api/dto"
	"github.com/duongnln96/building-microservices-golang/product-api/helper"
	"github.com/duongnln96/building-microservices-golang/product-api/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthenControllerI interface {
	Login(c echo.Context) error
	Register(c echo.Context) error
}

type AuthenControllerDeps struct {
	Log           *zap.SugaredLogger
	AuthenService service.AuthenServiceI
	JwtService    service.JWTSeriveI
}

type authenController struct {
	log           *zap.SugaredLogger
	authenService service.AuthenServiceI
	jwtService    service.JWTSeriveI
}

func NewAuthorizeController(deps AuthenControllerDeps) AuthenControllerI {
	return &authenController{
		log:           deps.Log,
		authenService: deps.AuthenService,
		jwtService:    deps.JwtService,
	}
}

func (ac *authenController) Login(c echo.Context) error {
	userLoginDto := dto.UserLogInDTO{}

	if err := c.Bind(&userLoginDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err, nil)
		return c.JSON(http.StatusBadRequest, res)
	}

	if err := c.Validate(userLoginDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err, nil)
		return c.JSON(http.StatusBadRequest, res)
	}

	if ok, user := ac.authenService.ValidateCredential(&userLoginDto); ok {
		token := ac.jwtService.GenerateToken(user.Email)
		res := helper.BuildResponse("OK", token)
		return c.JSON(http.StatusAccepted, res)
	}

	res := helper.BuildErrorResponse("Fail", fmt.Errorf("Please check your credential"), nil)
	return c.JSON(http.StatusUnauthorized, res)
}

func (ac *authenController) Register(c echo.Context) error {
	userRegDto := dto.UserRegisterDTO{}

	if err := c.Bind(&userRegDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err, nil)
		return c.JSON(http.StatusBadRequest, res)
	}

	if err := c.Validate(userRegDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err, nil)
		return c.JSON(http.StatusBadRequest, res)
	}

	if err := ac.authenService.CreateUser(&userRegDto); err != nil {
		res := helper.BuildErrorResponse("Fail", err, nil)
		return c.JSON(http.StatusInternalServerError, res)
	}

	res := helper.BuildResponse("OK", nil)
	return c.JSON(http.StatusOK, res)
}
