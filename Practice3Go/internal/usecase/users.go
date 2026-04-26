package usecase

import (
	"github.com/Yerassyl005/go-practice3/internal/repository"
	"github.com/Yerassyl005/go-practice3/pkg/modules"
)

type UserUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

//////////////////////
// GET ALL USERS
//////////////////////

func (u *UserUsecase) GetUsers(limit, offset int) ([]modules.User, error) {
	return u.repo.GetUsers(limit, offset)
}

//////////////////////
// GET USER BY ID
//////////////////////

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
	return u.repo.GetUserByID(id)
}

//////////////////////
// CREATE USER
//////////////////////

func (u *UserUsecase) CreateUser(user *modules.User) (int, error) {
	return u.repo.CreateUser(user)
}

//////////////////////
// UPDATE USER
//////////////////////

func (u *UserUsecase) UpdateUser(user *modules.User) error {
	return u.repo.UpdateUser(user)
}

//////////////////////
// DELETE USER
//////////////////////

func (u *UserUsecase) DeleteUser(id int) error {
	return u.repo.DeleteUser(id)
}