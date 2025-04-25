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

func (s *Service) CreateCompany(jwtString string,
	CompanyName,
	CompanyAddress,
	Drive,
	TypeOfDrive,
	FollowUp,
	Remarks,
	ContactDetails,
	HRDetails string,
	IsContacted bool) (string, error) {
	// if err := s.validateJWTAndRole(jwtString); err != nil {
	// 	return "", err
	// }

	companyData, err := entity.NewCompany(CompanyName,
		CompanyAddress,
		Drive,
		TypeOfDrive,
		FollowUp,
		Remarks,
		ContactDetails,
		HRDetails,
		IsContacted)
	if err != nil {
		log.Printf("unable to create company, err=%v", err)
		return "", err
	}

	err = s.repo.CreateCompany(companyData)
	if err != nil {
		log.Printf("unable to create company in repo, err=%v", err)
		return "", err
	}

	return companyData.CompanyID, nil
}

func (s *Service) GetCompany(jwtString string, id string) (*entity.CompanyData, error) {
	// if err := s.validateJWTAndRole(jwtString); err != nil {
	// 	return nil, err
	// }

	companyData, err := s.repo.GetCompany(id)
	if err != nil {
		log.Printf("unable to get company, err=%v", err)
		return nil, err
	}

	return companyData, nil
}

func (s *Service) GetCompanyByName(jwtString string, name string) (*entity.CompanyData, error) {
	// if err := s.validateJWTAndRole(jwtString); err != nil {
	// 	return nil, err
	// }

	companyData, err := s.repo.GetCompany(name)
	if err != nil {
		log.Printf("unable to get company, err=%v", err)
		return nil, err
	}

	return companyData, nil
}

func (s *Service) UpdateCompany(jwtString string,
	CompanyID,
	CompanyName,
	CompanyAddress,
	Drive,
	TypeOfDrive,
	FollowUp,
	Remarks,
	ContactDetails,
	HRDetails string,
	IsContacted bool) (*entity.CompanyData, error) {
	// if err := s.validateJWTAndRole(jwtString); err != nil {
	// 	return "", err
	// }

	companyData, err := entity.NewCompany(CompanyName,
		CompanyAddress,
		Drive,
		TypeOfDrive,
		FollowUp,
		Remarks,
		ContactDetails,
		HRDetails,
		IsContacted)
	if err != nil {
		log.Printf("unable to create company, err=%v", err)
		return nil, err
	}

	err = s.repo.UpdateCompany(companyData)
	if err != nil {
		log.Printf("unable to create company in repo, err=%v", err)
		return nil, err
	}

	return companyData, nil
}

// TODO check jwt token before calling this function
func (s *Service) GetAwaitingApproval() ([]*entity.CompanyData, error) {
	return s.repo.GetAwaitingApproval()
}

func (s *Service) SetAwaitingApproval(jwtString, companyID string, isApproved bool) (*entity.CompanyData, error) {
	// TODO: Add JWT validation/role check if necessary
	// if err := s.validateJWTAndRole(jwtString); err != nil { // Add appropriate role check later
	// 	return nil, err
	// }

	companyData, err := s.repo.SetAwaitingApproval(companyID, isApproved)
	if err != nil {
		log.Printf("unable to set awaiting approval status for company %s, err=%v", companyID, err)
		return nil, err
	}
	return companyData, nil
}
