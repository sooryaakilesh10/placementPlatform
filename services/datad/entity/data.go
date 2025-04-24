package entity

import "github.com/google/uuid"

type Data struct {
	DataID      string
	CompanyData CompanyData
}

type CompanyData struct {
	CompanyID      string
	CompanyName    string
	CompanyAddress string
	Drive          string
	TypeOfDrive    string
	FollowUp       string
	IsContacted    bool
	Remarks        string
	ContactDetails string
	HRDetails      string
}

func NewCompany(companyName,
	CompanyAddress,
	Drive,
	TypeOfDrive,
	FollowUp,
	Remarks,
	ContactDetails,
	HRDetails string,
	isContacted bool,
) (*CompanyData, error) {
	return &CompanyData{
		CompanyID:      uuid.NewString(),
		CompanyName:    companyName,
		Drive:          Drive,
		TypeOfDrive:    TypeOfDrive,
		FollowUp:       FollowUp,
		IsContacted:    false,
		Remarks:        Remarks,
		ContactDetails: ContactDetails,
	}, nil
}

func (c *CompanyData) Validate() error {
	return nil
}
