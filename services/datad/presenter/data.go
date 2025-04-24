package presenter

type CreateCompanyRequest struct {
	JWT            string `json:"jwt"`
	CompanyID      string `json:"companyID"`
	CompanyName    string `json:"companyName"`
	CompanyAddress string `json:"companyAddress"`
	Drive          string `json:"drive"`
	TypeOfDrive    string `json:"typeOfDrive"`
	FollowUp       string `json:"followUp"`
	IsContacted    bool   `json:"isContacted"`
	Remarks        string `json:"remarks"`
	ContactDetails string `json:"contactDetails"`
	HRDetails      string `json:"hrDetails"`
}

type GetCompanyResponse struct {
	CompanyID      string `json:"companyID"`
	CompanyName    string `json:"companyName"`
	CompanyAddress string `json:"companyAddress"`
	Drive          string `json:"drive"`
	TypeOfDrive    string `json:"typeOfDrive"`
	FollowUp       string `json:"followUp"`
	IsContacted    bool   `json:"isContacted"`
	Remarks        string `json:"remarks"`
	ContactDetails string `json:"contactDetails"`
	HRDetails      string `json:"hrDetails"`
}

type GetCompanyRequest struct {
	JWT string `json:"jwt"`
	ID  string `json:"id"`
}

type CreateCompanyResponse struct {
	CompanyID string `json:"companyID"`
}
