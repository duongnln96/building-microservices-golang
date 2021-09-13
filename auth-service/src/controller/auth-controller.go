package controller

import (
	"fmt"
	"net/http"

	"github.com/duongnln96/building-microservices-golang/auth-service/src/domain/dto"
	"github.com/duongnln96/building-microservices-golang/auth-service/src/service"
	"github.com/duongnln96/building-microservices-golang/auth-service/tools/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthControllerI interface {
	Login(c echo.Context) error
	Register(c echo.Context) error
}

type AuthControllerDeps struct {
	Log     *zap.SugaredLogger
	AuthSvc service.AuthServiceI
	JWTSvc  service.JWTServiceI
}

type authController struct {
	log     *zap.SugaredLogger
	authSvc service.AuthServiceI
	jwtSvc  service.JWTServiceI
}

func NewAuthController(deps AuthControllerDeps) AuthControllerI {
	return &authController{
		log:     deps.Log,
		authSvc: deps.AuthSvc,
		jwtSvc:  deps.JWTSvc,
	}
}

func (ac *authController) Register(c echo.Context) error {
	userRegDTO := dto.UserRegisterDTO{}

	if err := c.Bind(&userRegDTO); err != nil {
		res := utils.BuildErrorResponse("Fail", err.Error())
		return c.JSON(http.StatusBadRequest, res)
	}

	if err := c.Validate(userRegDTO); err != nil {
		res := utils.BuildErrorResponse("Fail", err.Error())
		return c.JSON(http.StatusBadRequest, res)
	}

	if err := ac.authSvc.CreateUser(&userRegDTO); err != nil {
		res := utils.BuildErrorResponse("Fail", err.Error())
		return c.JSON(http.StatusInternalServerError, res)
	}

	res := utils.BuildResponse("OK", nil)
	return c.JSON(http.StatusOK, res)
}

func (ac *authController) Login(c echo.Context) error {
	userLoginDto := dto.UserLoginDTO{}

	if err := c.Bind(&userLoginDto); err != nil {
		res := utils.BuildErrorResponse("Fail", err.Error())
		return c.JSON(http.StatusBadRequest, res)
	}

	if err := c.Validate(userLoginDto); err != nil {
		res := utils.BuildErrorResponse("Fail", err.Error())
		return c.JSON(http.StatusBadRequest, res)
	}

	if ok, user := ac.authSvc.ValidateCredential(&userLoginDto); ok {
		token := ac.jwtSvc.GenerateToken(user.Email)
		res := utils.BuildResponse("OK", token)
		return c.JSON(http.StatusAccepted, res)
	}

	res := utils.BuildErrorResponse("Fail", fmt.Errorf("Please check your credential").Error())
	return c.JSON(http.StatusUnauthorized, res)
}
