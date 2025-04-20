package user

import "backend/services/userd/entity"

type Repository interface {
	Writer
	Reader
}

type Writer interface {
	CreateUser(user *entity.User) (string, error)
}

type Reader interface {
	GetUserByID(id string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
}

type Usecase interface {
	CreateUser(userName, email, pass, role string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	Login(email, pass string) (string, error)
}
