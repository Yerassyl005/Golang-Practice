package repository

import (
	"github.com/Yerassyl005/go-practice3/internal/repository/_postgres"
	"github.com/Yerassyl005/go-practice3/internal/repository/_postgres/users"
	"github.com/Yerassyl005/go-practice3/pkg/modules"
)

type UserRepository interface {
	GetUsers(limit, offset int) ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(user *modules.User) (int, error)
	UpdateUser(user *modules.User) error
	DeleteUser(id int) error
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}