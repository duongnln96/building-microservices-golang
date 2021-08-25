package service

import (
	"github.com/duongnln96/building-microservices-golang/product-api/dto"
	"github.com/duongnln96/building-microservices-golang/product-api/entity"
	"github.com/duongnln96/building-microservices-golang/product-api/helper"
	"github.com/duongnln96/building-microservices-golang/product-api/repository"
	"go.uber.org/zap"
)

type AuthenServiceI interface {
	FindUserByEmail(email string) (entity.User, error)
	CreateUser(u *dto.UserRegisterDTO) error
	IsDuplicateEmail(email string) (bool, error)
	ValidateCredential(userCredential *dto.UserLogInDTO) (bool, entity.User)
}

type AuthenServiceDeps struct {
	Log      *zap.SugaredLogger
	UserRepo repository.UserRepoI
}

type authenService struct {
	log      *zap.SugaredLogger
	userRepo repository.UserRepoI
}

func NewAuthenService(deps AuthenServiceDeps) AuthenServiceI {
	return &authenService{
		log:      deps.Log,
		userRepo: deps.UserRepo,
	}
}

func (svc *authenService) FindUserByEmail(email string) (entity.User, error) {
	return svc.userRepo.FindUserByEmail(email)
}

func (svc *authenService) CreateUser(u *dto.UserRegisterDTO) error {
	user := entity.User{}
	user.Name = u.Name
	user.Email = u.Email
	user.Password = helper.Encrypt(u.Password)

	err := svc.userRepo.CreateUser(&user)
	if err != nil {
		return err
	}
	return nil
}

func (svc *authenService) IsDuplicateEmail(email string) (bool, error) {
	return svc.userRepo.CheckEmailExisted(email)
}

func (svc *authenService) ValidateCredential(userCredential *dto.UserLogInDTO) (bool, entity.User) {
	user, err := svc.userRepo.FindUserByEmail(userCredential.Email)
	if err != nil {
		return false, user
	}

	if user.Password != helper.Encrypt(userCredential.Password) {
		return false, user
	}

	return true, user
}
