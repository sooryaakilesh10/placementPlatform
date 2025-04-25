package data

import "backend/services/datad/entity"

type Repository interface {
	Writer
	Reader
}

type Writer interface {
	CreateCompany(companyData *entity.CompanyData) error
	UpdateCompany(companyData *entity.CompanyData) error
	SetAwaitingApproval(companyID string, isApproved bool) (*entity.CompanyData, error)
}

type Reader interface {
	GetCompany(id string) (*entity.CompanyData, error)
	GetAwaitingApproval() ([]*entity.CompanyData, error)
}

type Usecase interface {
	CreateCompany(jwt,
		companyName,
		CompanyAddress,
		Drive,
		TypeOfDrive,
		FollowUp,
		Remarks,
		ContactDetails,
		HRDetails string,
		isContacted bool) (string, error)
	GetCompany(jwtString, id string) (*entity.CompanyData, error)
	GetCompanyByName(jwtString, name string) (*entity.CompanyData, error)
	UpdateCompany(jwt,
		companyID,
		companyName,
		CompanyAddress,
		Drive,
		TypeOfDrive,
		FollowUp,
		Remarks,
		ContactDetails,
		HRDetails string,
		isContacted bool) (*entity.CompanyData, error)
	GetAwaitingApproval() ([]*entity.CompanyData, error)
	SetAwaitingApproval(jwtString, companyID string, isApproved bool) (*entity.CompanyData, error)
}
