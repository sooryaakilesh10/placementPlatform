package user

import (
	"backend/services/userd/entity"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	repo      Repository
	JWTSecret string
}

func NewService(repo Repository, jwt string) *Service {
	return &Service{
		repo:      repo,
		JWTSecret: jwt,
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

func (s *Service) Login(email, pass string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		log.Printf("unable to get user by email, err=%v", err)
		return "", err
	}

	if user.Pass != pass {
		log.Printf("invalid password for user %s", user.UserName)
		return "", err
	}

	return s.generateJWT(user.UserID, user.Email, user.UserName, user.Role)
}

func (s *Service) generateJWT(ID, name, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": ID,
		"user_name": name,
		"email":  email,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.JWTSecret))
}
