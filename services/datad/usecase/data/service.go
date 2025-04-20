package data

import (
	"backend/pkg/common"
	"backend/services/datad/entity"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	repo      Repository
	JWTSecret string
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		JWTSecret: jwtSecret,
	}
}

// NewServiceWithAuth creates a new data service with JWT authentication
func NewServiceWithAuth(repo Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		JWTSecret: jwtSecret,
	}
}

func (s *Service) validateJWTAndRole(jwtString string) error {
	// Skip JWT validation if no JWT secret is set or jwt string is empty
	if s.JWTSecret == "" || jwtString == "" {
		log.Printf("Skipping JWT validation: secret empty=%v, token empty=%v", s.JWTSecret == "", jwtString == "")
		return nil
	}

	// Validate JWT and check if the user has the required role
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.JWTSecret), nil
	})

	if err != nil {
		log.Printf("unable to parse JWT, err=%v", err)
		return fmt.Errorf("invalid token: %w", err)
	}

	// Check if token is nil
	if token == nil {
		return fmt.Errorf("invalid token: token is nil")
	}

	// Try to extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid token: failed to process claims")
	}

	// Check for required role
	roleClaim, ok := claims["role"]
	if !ok {
		return fmt.Errorf("invalid token: role claim missing")
	}

	role, ok := roleClaim.(string)
	if !ok {
		return fmt.Errorf("invalid token: role claim type invalid")
	}

	// Check if ValidRolesToCreateData is nil or empty
	if common.ValidRolesToCreateData == nil || len(common.ValidRolesToCreateData) == 0 {
		return nil
	}

	// Check if the user's role is in the list of valid roles
	isAllowed := false
	for _, allowedRole := range common.ValidRolesToCreateData {
		if role == allowedRole {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		log.Printf("user role '%s' does not have permission to create data", role)
		return fmt.Errorf("permission denied: insufficient role")
	}

	return nil
}

func (s *Service) CreateData(jwtString string, companyData interface{}) (*entity.Data, error) {
	if err := s.validateJWTAndRole(jwtString); err != nil {
		return nil, err
	}

	// Check if companyData is nil
	if companyData == nil {
		return nil, fmt.Errorf("company data cannot be nil")
	}

	// Create entity and save to repository
	dataEn := entity.NewData(companyData)
	dataID, err := s.repo.CreateData(dataEn)
	if err != nil {
		log.Printf("unable to create data in repository, err=%v", err)
		return nil, err
	}

	dataEn.DataID = dataID
	return dataEn, nil
}

func (s *Service) GetDataByID(id string) (*entity.Data, error) {
	data, err := s.repo.GetDataByID(id)
	if err != nil {
		log.Printf("unable to get data by ID, err=%v", err)
		return nil, err
	}
	return data, nil
}
