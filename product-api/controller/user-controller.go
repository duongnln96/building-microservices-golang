package controller

type UserControllerI interface {
}

type UserControllerDeps struct {
}

type userController struct {
}

func NewUserController(deps UserControllerI) UserControllerI {
	return &userController{}
}
