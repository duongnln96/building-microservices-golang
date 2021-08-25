package service

import (
	"github.com/duongnln96/building-microservices-golang/product-api/dto"
	"github.com/duongnln96/building-microservices-golang/product-api/entity"
	"github.com/duongnln96/building-microservices-golang/product-api/helper"
	"github.com/duongnln96/building-microservices-golang/product-api/repository"
	"go.uber.org/zap"
)

type UserServiceI interface {
	UpdateUser(*dto.UserUpdateDTO) error
}

type UserServiceDeps struct {
	Log  *zap.SugaredLogger
	Repo repository.UserRepoI
}

type userService struct {
	log  *zap.SugaredLogger
	repo repository.UserRepoI
}

func NewUserService(deps UserServiceDeps) UserServiceI {
	return &userService{
		log:  deps.Log,
		repo: deps.Repo,
	}
}

func (svc *userService) UpdateUser(u *dto.UserUpdateDTO) error {
	oldData, err := svc.repo.FintUserByID(u.ID)
	if err != nil {
		svc.log.Debugf("User service cannot find user want to update %+v", err)
		return err
	}

	user := entity.User{}
	user.Id = u.ID

	if u.Name != oldData.Name {
		user.Name = u.Name
	} else {
		user.Name = oldData.Name
	}

	if u.Email != oldData.Email {
		user.Email = u.Email
	} else {
		user.Email = oldData.Email
	}

	if u.Password != "" {
		user.Password = helper.Encrypt(u.Password)
	}

	err = svc.repo.UpdateUser(&user)
	if err != nil {
		svc.log.Debugf("User service cannot update user %+v", err)
		return err
	}

	return nil
}
