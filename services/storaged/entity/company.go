package entity

import (
	"errors"
	"time"
)

// DriveStatus represents the status of a recruitment drive
type DriveStatus string

const (
	DriveStatusScheduled DriveStatus = "SCHEDULED"
	DriveStatusCompleted DriveStatus = "COMPLETED"
	DriveStatusNoHiring  DriveStatus = "NO_HIRING"
)

// HR contains contact information for company HR
type HR struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Position    string `json:"position"`
	LinkedInURL string `json:"linkedin_url,omitempty"`
}

// DriveDetails contains information about a recruitment drive
type DriveDetails struct {
	Status           DriveStatus `json:"status"`
	ScheduledDate    time.Time   `json:"scheduled_date,omitempty"`
	NumberOfOffers   int         `json:"number_of_offers,omitempty"`
	NumberHired      int         `json:"number_hired,omitempty"`
	RolesOffered     []string    `json:"roles_offered,omitempty"`
	MinCGPA          float32     `json:"min_cgpa,omitempty"`
	EligibleBranches []string    `json:"eligible_branches,omitempty"`
	Notes            string      `json:"notes,omitempty"`
}

// CompanyDetails contains general information about a company
type CompanyDetails struct {
	Website     string `json:"website"`
	Industry    string `json:"industry"`
	FoundedYear int    `json:"founded_year,omitempty"`
	Size        string `json:"size,omitempty"`
	Description string `json:"description,omitempty"`
	Logo        string `json:"logo_url,omitempty"`
}

// Company represents a company in the placement system
type Company struct {
	ID              string
	CompanyName     string
	LastContacted   bool
	FollowUp        time.Time
	Packages        []float32
	Remarks         string
	HR              HR
	TargetBranch    string
	IsValidation    bool
	Approved        bool
	DriveDetails    DriveDetails
	CompanyDetails  CompanyDetails
	Location        string
	AssignedTo      string // ID of the placement officer assigned
	CreatedBy       string // ID of the user who created this company
	CreatedAt       time.Time
	UpdatedAt       time.Time
	IsDataValidated bool   // Flag to track if company data is validated
	ApprovalStatus  string // PENDING, APPROVED, REJECTED
	ApprovalNotes   string // Notes about approval
}

func NewCompany(
	id uint64,
	companyName string,
	lastContacted bool,
	followUp time.Time,
	remarks string,
	targetBranch string,
	packages []float32,
	hr HR,
	driveDetails DriveDetails,
	companyDetails CompanyDetails,
	approved bool,
	isValidated bool,
	location string,
) (*Company, error) {
	company := &Company{
		ID:              "", // Will be set by the repository
		CompanyName:     companyName,
		LastContacted:   lastContacted,
		FollowUp:        followUp,
		Remarks:         remarks,
		TargetBranch:    targetBranch,
		Packages:        packages,
		HR:              hr,
		DriveDetails:    driveDetails,
		CompanyDetails:  companyDetails,
		Approved:        approved,
		IsValidation:    isValidated,
		Location:        location,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ApprovalStatus:  "PENDING",
		IsDataValidated: false,
	}

	if err := company.validate(); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *Company) validate() error {
	if s.CompanyName == "" {
		return errors.New("company name cannot be empty")
	}

	if s.TargetBranch == "" {
		return errors.New("target branch cannot be empty")
	}

	if len(s.Packages) == 0 {
		return errors.New("packages cannot be empty")
	}

	if s.Location == "" {
		return errors.New("location cannot be empty")
	}

	return nil
}

// ValidateData marks the company data as validated
func (s *Company) ValidateData() {
	s.IsDataValidated = true
	s.UpdatedAt = time.Now()
}

// AssignToOfficer assigns the company to a placement officer
func (s *Company) AssignToOfficer(officerID string) {
	s.AssignedTo = officerID
	s.UpdatedAt = time.Now()
}

// UpdateDriveStatus updates the status of the recruitment drive
func (s *Company) UpdateDriveStatus(status DriveStatus, details DriveDetails) {
	s.DriveDetails = details
	s.DriveDetails.Status = status
	s.UpdatedAt = time.Now()
	// When drive is updated, we reset approval to pending
	s.ApprovalStatus = "PENDING"
}

// SetApprovalStatus sets the approval status of company updates
func (s *Company) SetApprovalStatus(status string, notes string) {
	s.ApprovalStatus = status
	s.ApprovalNotes = notes
	s.UpdatedAt = time.Now()
}
