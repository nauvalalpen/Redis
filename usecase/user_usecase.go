package usecase

import "redis/domain"

type UserUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: repo}
}

func (u *UserUsecase) CreateUser(user domain.User) error {
	return u.userRepo.Save(user)
}