package presenter

import "time"

// HRData contains contact information for company HR
type HRData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Position    string `json:"position"`
	LinkedInURL string `json:"linkedin_url,omitempty"`
}

// DriveDetailsData contains information about a recruitment drive
type DriveDetailsData struct {
	Status           string    `json:"status"`
	ScheduledDate    time.Time `json:"scheduled_date,omitempty"`
	NumberOfOffers   int       `json:"number_of_offers,omitempty"`
	NumberHired      int       `json:"number_hired,omitempty"`
	RolesOffered     []string  `json:"roles_offered,omitempty"`
	MinCGPA          float32   `json:"min_cgpa,omitempty"`
	EligibleBranches []string  `json:"eligible_branches,omitempty"`
	Notes            string    `json:"notes,omitempty"`
}

// CompanyDetailsData contains general information about a company
type CompanyDetailsData struct {
	Website     string `json:"website"`
	Industry    string `json:"industry"`
	FoundedYear int    `json:"founded_year,omitempty"`
	Size        string `json:"size,omitempty"`
	Description string `json:"description,omitempty"`
	Logo        string `json:"logo_url,omitempty"`
}

// CompanyRequest is used for creating a new company
type CompanyRequest struct {
	ID             uint64             `json:"id"`
	CompanyName    string             `json:"company_name"`
	LastContacted  bool               `json:"last_contacted"`
	FollowUp       time.Time          `json:"follow_up"`
	Packages       []float32          `json:"packages"`
	Remarks        string             `json:"remarks"`
	HR             HRData             `json:"hr"`
	TargetBranch   string             `json:"target_branch"`
	IsValidation   bool               `json:"is_validation"`
	Approved       bool               `json:"approved"`
	DriveDetails   DriveDetailsData   `json:"drive_details"`
	CompanyDetails CompanyDetailsData `json:"company_details"`
	Location       string             `json:"location"`
	CreatedBy      string             `json:"created_by"`
}

// CompanyUpdateRequest is used for updating an existing company
type CompanyUpdateRequest struct {
	ID             string             `json:"id"`
	CompanyName    string             `json:"company_name,omitempty"`
	LastContacted  bool               `json:"last_contacted,omitempty"`
	FollowUp       time.Time          `json:"follow_up,omitempty"`
	Packages       []float32          `json:"packages,omitempty"`
	Remarks        string             `json:"remarks,omitempty"`
	HR             HRData             `json:"hr,omitempty"`
	TargetBranch   string             `json:"target_branch,omitempty"`
	IsValidation   bool               `json:"is_validation,omitempty"`
	DriveDetails   DriveDetailsData   `json:"drive_details,omitempty"`
	CompanyDetails CompanyDetailsData `json:"company_details,omitempty"`
	Location       string             `json:"location,omitempty"`
	UpdatedBy      string             `json:"updated_by"`
}

// CompanyResponse is used for sending company data back to the client
type CompanyResponse struct {
	ID              string             `json:"id"`
	CompanyName     string             `json:"company_name"`
	LastContacted   bool               `json:"last_contacted"`
	FollowUp        time.Time          `json:"follow_up"`
	Packages        []float32          `json:"packages"`
	Remarks         string             `json:"remarks"`
	HR              HRData             `json:"hr"`
	TargetBranch    string             `json:"target_branch"`
	IsValidation    bool               `json:"is_validation"`
	Approved        bool               `json:"approved"`
	DriveDetails    DriveDetailsData   `json:"drive_details"`
	CompanyDetails  CompanyDetailsData `json:"company_details"`
	Location        string             `json:"location"`
	AssignedTo      string             `json:"assigned_to,omitempty"`
	CreatedBy       string             `json:"created_by,omitempty"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty"`
	IsDataValidated bool               `json:"is_data_validated"`
	ApprovalStatus  string             `json:"approval_status"`
	ApprovalNotes   string             `json:"approval_notes,omitempty"`
}

// AssignmentRequest is used for assigning a company to a placement officer
type AssignmentRequest struct {
	CompanyID  string `json:"company_id"`
	OfficerID  string `json:"officer_id"`
	AssignedBy string `json:"assigned_by"`
}

// ApprovalRequest is used for updating a company's approval status
type ApprovalRequest struct {
	CompanyID string `json:"company_id"`
	Status    string `json:"status"` // APPROVED or REJECTED
	Notes     string `json:"notes,omitempty"`
	UpdatedBy string `json:"updated_by"`
}

// SystemSettingsRequest is used for updating system settings
type SystemSettingsRequest struct {
	ApprovalMode string `json:"approval_mode"` // auto or manual
	UpdatedBy    string `json:"updated_by"`
}
