package user

import (
	"backend/services/storaged/entity"
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

func (s *Service) CreateCompany(company *entity.Company) (string, error) {
	companyID, err := s.repo.CreateCompany(company)
	if err != nil {
		log.Printf("unable to create company in repository, err=%v", err)
		return "", err
	}
	return companyID, nil
}

func (s *Service) GetCompanyByID(id string) (*entity.Company, error) {
	company, err := s.repo.GetCompanyByID(id)
	if err != nil {
		log.Printf("unable to get company by id, err=%v", err)
		return nil, err
	}
	return company, nil
}

func (s *Service) GetCompanyByEmail(email string) (*entity.Company, error) {
	company, err := s.repo.GetCompanyByEmail(email)
	if err != nil {
		log.Printf("unable to get company by email, err=%v", err)
		return nil, err
	}
	return company, nil
}

func (s *Service) UpdateCompany(company *entity.Company) error {
	err := s.repo.UpdateCompany(company)
	if err != nil {
		log.Printf("unable to update company, err=%v", err)
		return err
	}
	return nil
}

func (s *Service) AssignCompanyToOfficer(companyID, officerID, assignedBy string) error {
	err := s.repo.AssignCompanyToOfficer(companyID, officerID, assignedBy)
	if err != nil {
		log.Printf("unable to assign company to officer, err=%v", err)
		return err
	}
	return nil
}

func (s *Service) GetCompaniesByOfficerID(officerID string) ([]*entity.Company, error) {
	companies, err := s.repo.GetCompaniesByOfficerID(officerID)
	if err != nil {
		log.Printf("unable to get companies by officer ID, err=%v", err)
		return nil, err
	}
	return companies, nil
}

func (s *Service) GetAllCompanies() ([]*entity.Company, error) {
	companies, err := s.repo.GetAllCompanies()
	if err != nil {
		log.Printf("unable to get all companies, err=%v", err)
		return nil, err
	}
	return companies, nil
}

func (s *Service) ImportCompaniesFromCSV(companies []*entity.Company) ([]string, error) {
	importedIDs, err := s.repo.ImportCompaniesFromCSV(companies)
	if err != nil {
		log.Printf("unable to import companies from CSV, err=%v", err)
		return nil, err
	}
	return importedIDs, nil
}

func (s *Service) UpdateApprovalStatus(companyID, status, notes string, updatedBy string) error {
	err := s.repo.UpdateApprovalStatus(companyID, status, notes, updatedBy)
	if err != nil {
		log.Printf("unable to update approval status, err=%v", err)
		return err
	}
	return nil
}
