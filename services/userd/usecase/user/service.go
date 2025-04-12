package user

import (
	"backend/services/userd/entity"
	"log"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(userName, email, pass, role string) (*entity.User, error) {
	user, err := entity.NewUser(userName, email, pass, role)
	if err != nil {
		log.Printf("unable to create user entity, err=%v", err)
		return nil, err
	}
	userID, err := s.repo.CreateUser(user)
	if err != nil {
		log.Printf("unable to create user in repository, err=%v", err)
		return nil, err
	}
	user.UserID = userID
	return user, nil
}

func (s *Service) GetUserByID(id string) (*entity.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		log.Printf("unable to get user by id, err=%v", err)
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUserByEmail(email string) (*entity.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		log.Printf("unable to get user by email, err=%v", err)
		return nil, err
	}
	return user, nil
}
