package user

import "backend/services/storaged/entity"

type Repository interface {
	Writer
	Reader
}

type Writer interface {
	CreateCompany(company *entity.Company) (string, error)
	UpdateCompany(company *entity.Company) error
	AssignCompanyToOfficer(companyID, officerID, assignedBy string) error
	UpdateApprovalStatus(companyID, status, notes string, updatedBy string) error
	ImportCompaniesFromCSV(companies []*entity.Company) ([]string, error)
}

type Reader interface {
	GetCompanyByID(id string) (*entity.Company, error)
	GetCompanyByEmail(email string) (*entity.Company, error)
	GetCompaniesByOfficerID(officerID string) ([]*entity.Company, error)
	GetAllCompanies() ([]*entity.Company, error)
}

type Usecase interface {
	// Company CRUD operations
	CreateCompany(company *entity.Company) (string, error)
	GetCompanyByID(id string) (*entity.Company, error)
	GetCompanyByEmail(email string) (*entity.Company, error)
	UpdateCompany(company *entity.Company) error

	// Assignment operations
	AssignCompanyToOfficer(companyID, officerID, assignedBy string) error
	GetCompaniesByOfficerID(officerID string) ([]*entity.Company, error)
	GetAllCompanies() ([]*entity.Company, error)

	// Bulk operations
	ImportCompaniesFromCSV(companies []*entity.Company) ([]string, error)

	// Approval operations
	UpdateApprovalStatus(companyID, status, notes string, updatedBy string) error
}
