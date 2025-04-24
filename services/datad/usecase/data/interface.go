package data

import "backend/services/datad/entity"

type Repository interface {
	Writer
	Reader
}

type Writer interface {
	CreateCompany(companyData *entity.CompanyData) error
}

type Reader interface {
	GetCompany(id string) (*entity.CompanyData, error)
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
}
