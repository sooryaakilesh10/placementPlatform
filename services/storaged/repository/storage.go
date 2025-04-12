package repository

import (
	"backend/services/storaged/entity"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateCompany(company *entity.Company) (string, error) {
	// Generate a new UUID for the company
	companyID := uuid.NewString()
	company.ID = companyID

	// Marshal JSON arrays
	packagesJSON, err := json.Marshal(company.Packages)
	if err != nil {
		return "", err
	}

	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert into companies table
	_, err = tx.Exec(`
		INSERT INTO companies (
			company_id, company_name, last_contacted, follow_up, 
			packages, remarks, target_branch, is_validation, 
			approved, location, hr_name, hr_email, hr_phone, 
			hr_position, hr_linkedin_url, website, industry, 
			founded_year, company_size, description, logo_url, 
			assigned_to, created_by, is_data_validated, 
			approval_status, approval_notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		company.ID,
		company.CompanyName,
		company.LastContacted,
		company.FollowUp,
		packagesJSON,
		company.Remarks,
		company.TargetBranch,
		company.IsValidation,
		company.Approved,
		company.Location,
		company.HR.Name,
		company.HR.Email,
		company.HR.Phone,
		company.HR.Position,
		company.HR.LinkedInURL,
		company.CompanyDetails.Website,
		company.CompanyDetails.Industry,
		company.CompanyDetails.FoundedYear,
		company.CompanyDetails.Size,
		company.CompanyDetails.Description,
		company.CompanyDetails.Logo,
		company.AssignedTo,
		company.CreatedBy,
		company.IsDataValidated,
		company.ApprovalStatus,
		company.ApprovalNotes,
	)

	if err != nil {
		return "", err
	}

	// Create recruitment drive if provided
	driveID := uuid.NewString()
	if company.DriveDetails.Status != "" {
		rolesJSON, err := json.Marshal(company.DriveDetails.RolesOffered)
		if err != nil {
			return "", err
		}

		branchesJSON, err := json.Marshal(company.DriveDetails.EligibleBranches)
		if err != nil {
			return "", err
		}

		_, err = tx.Exec(`
			INSERT INTO recruitment_drives (
				drive_id, company_id, status, scheduled_date, 
				number_of_offers, number_hired, roles_offered, 
				min_cgpa, eligible_branches, notes, created_by
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			driveID,
			company.ID,
			company.DriveDetails.Status,
			company.DriveDetails.ScheduledDate,
			company.DriveDetails.NumberOfOffers,
			company.DriveDetails.NumberHired,
			rolesJSON,
			company.DriveDetails.MinCGPA,
			branchesJSON,
			company.DriveDetails.Notes,
			company.CreatedBy,
		)

		if err != nil {
			return "", err
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return "", err
	}

	return company.ID, nil
}

func (r *Repository) GetCompanyByID(id string) (*entity.Company, error) {
	var company entity.Company
	var packagesJSON, rolesJSON, branchesJSON []byte
	var hrName, hrEmail, hrPhone, hrPosition, hrLinkedInURL string
	var website, industry, size, description, logoURL string
	var foundedYear int
	var driveStatus, driveNotes sql.NullString
	var scheduledDate sql.NullTime
	var numOffers, numHired sql.NullInt32
	var minCGPA sql.NullFloat64

	// Query company data
	err := r.db.QueryRow(`
		SELECT 
			c.company_id, c.company_name, c.last_contacted, c.follow_up, 
			c.packages, c.remarks, c.target_branch, c.is_validation, 
			c.approved, c.location, c.hr_name, c.hr_email, c.hr_phone, 
			c.hr_position, c.hr_linkedin_url, c.website, c.industry, 
			c.founded_year, c.company_size, c.description, c.logo_url, 
			c.assigned_to, c.created_by, c.created_at, c.updated_at, 
			c.is_data_validated, c.approval_status, c.approval_notes,
			rd.status, rd.scheduled_date, rd.number_of_offers, rd.number_hired,
			rd.roles_offered, rd.min_cgpa, rd.eligible_branches, rd.notes
		FROM companies c
		LEFT JOIN recruitment_drives rd ON c.company_id = rd.company_id
		WHERE c.company_id = ?
	`, id).Scan(
		&company.ID,
		&company.CompanyName,
		&company.LastContacted,
		&company.FollowUp,
		&packagesJSON,
		&company.Remarks,
		&company.TargetBranch,
		&company.IsValidation,
		&company.Approved,
		&company.Location,
		&hrName,
		&hrEmail,
		&hrPhone,
		&hrPosition,
		&hrLinkedInURL,
		&website,
		&industry,
		&foundedYear,
		&size,
		&description,
		&logoURL,
		&company.AssignedTo,
		&company.CreatedBy,
		&company.CreatedAt,
		&company.UpdatedAt,
		&company.IsDataValidated,
		&company.ApprovalStatus,
		&company.ApprovalNotes,
		&driveStatus,
		&scheduledDate,
		&numOffers,
		&numHired,
		&rolesJSON,
		&minCGPA,
		&branchesJSON,
		&driveNotes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("company not found")
		}
		return nil, err
	}

	// Set HR data
	company.HR = entity.HR{
		Name:        hrName,
		Email:       hrEmail,
		Phone:       hrPhone,
		Position:    hrPosition,
		LinkedInURL: hrLinkedInURL,
	}

	// Set company details
	company.CompanyDetails = entity.CompanyDetails{
		Website:     website,
		Industry:    industry,
		FoundedYear: foundedYear,
		Size:        size,
		Description: description,
		Logo:        logoURL,
	}

	// Unmarshal packages
	if err = json.Unmarshal(packagesJSON, &company.Packages); err != nil {
		return nil, err
	}

	// Set drive details if available
	if driveStatus.Valid {
		company.DriveDetails.Status = entity.DriveStatus(driveStatus.String)

		if scheduledDate.Valid {
			company.DriveDetails.ScheduledDate = scheduledDate.Time
		}

		if numOffers.Valid {
			company.DriveDetails.NumberOfOffers = int(numOffers.Int32)
		}

		if numHired.Valid {
			company.DriveDetails.NumberHired = int(numHired.Int32)
		}

		if minCGPA.Valid {
			company.DriveDetails.MinCGPA = float32(minCGPA.Float64)
		}

		if driveNotes.Valid {
			company.DriveDetails.Notes = driveNotes.String
		}

		// Unmarshal roles and branches if present
		if len(rolesJSON) > 0 {
			if err = json.Unmarshal(rolesJSON, &company.DriveDetails.RolesOffered); err != nil {
				log.Printf("Error unmarshaling roles: %v", err)
			}
		}

		if len(branchesJSON) > 0 {
			if err = json.Unmarshal(branchesJSON, &company.DriveDetails.EligibleBranches); err != nil {
				log.Printf("Error unmarshaling branches: %v", err)
			}
		}
	}

	return &company, nil
}

func (r *Repository) GetCompanyByEmail(email string) (*entity.Company, error) {
	// Search by HR email
	var companyID string
	err := r.db.QueryRow(`
		SELECT company_id FROM companies WHERE hr_email = ?
	`, email).Scan(&companyID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("company not found")
		}
		return nil, err
	}

	return r.GetCompanyByID(companyID)
}

func (r *Repository) UpdateCompany(company *entity.Company) error {
	company.UpdatedAt = time.Now()

	// Marshal JSON arrays
	packagesJSON, err := json.Marshal(company.Packages)
	if err != nil {
		return err
	}

	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update company record
	_, err = tx.Exec(`
		UPDATE companies SET
			company_name = ?, last_contacted = ?, follow_up = ?,
			packages = ?, remarks = ?, target_branch = ?, is_validation = ?,
			approved = ?, location = ?, hr_name = ?, hr_email = ?, hr_phone = ?,
			hr_position = ?, hr_linkedin_url = ?, website = ?, industry = ?,
			founded_year = ?, company_size = ?, description = ?, logo_url = ?,
			assigned_to = ?, is_data_validated = ?, approval_status = ?,
			approval_notes = ?, updated_at = ?
		WHERE company_id = ?
	`,
		company.CompanyName,
		company.LastContacted,
		company.FollowUp,
		packagesJSON,
		company.Remarks,
		company.TargetBranch,
		company.IsValidation,
		company.Approved,
		company.Location,
		company.HR.Name,
		company.HR.Email,
		company.HR.Phone,
		company.HR.Position,
		company.HR.LinkedInURL,
		company.CompanyDetails.Website,
		company.CompanyDetails.Industry,
		company.CompanyDetails.FoundedYear,
		company.CompanyDetails.Size,
		company.CompanyDetails.Description,
		company.CompanyDetails.Logo,
		company.AssignedTo,
		company.IsDataValidated,
		company.ApprovalStatus,
		company.ApprovalNotes,
		company.UpdatedAt,
		company.ID,
	)

	if err != nil {
		return err
	}

	// Update drive information if present
	if company.DriveDetails.Status != "" {
		rolesJSON, err := json.Marshal(company.DriveDetails.RolesOffered)
		if err != nil {
			return err
		}

		branchesJSON, err := json.Marshal(company.DriveDetails.EligibleBranches)
		if err != nil {
			return err
		}

		// Check if a drive record exists for this company
		var driveID string
		err = tx.QueryRow("SELECT drive_id FROM recruitment_drives WHERE company_id = ?", company.ID).Scan(&driveID)

		if err == nil {
			// Drive exists, update it
			_, err = tx.Exec(`
				UPDATE recruitment_drives SET
					status = ?, scheduled_date = ?, number_of_offers = ?,
					number_hired = ?, roles_offered = ?, min_cgpa = ?,
					eligible_branches = ?, notes = ?, updated_at = ?
				WHERE drive_id = ?
			`,
				company.DriveDetails.Status,
				company.DriveDetails.ScheduledDate,
				company.DriveDetails.NumberOfOffers,
				company.DriveDetails.NumberHired,
				rolesJSON,
				company.DriveDetails.MinCGPA,
				branchesJSON,
				company.DriveDetails.Notes,
				time.Now(),
				driveID,
			)
		} else if err == sql.ErrNoRows {
			// Drive doesn't exist, create it
			driveID = uuid.NewString()
			_, err = tx.Exec(`
				INSERT INTO recruitment_drives (
					drive_id, company_id, status, scheduled_date,
					number_of_offers, number_hired, roles_offered,
					min_cgpa, eligible_branches, notes, created_by
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
				driveID,
				company.ID,
				company.DriveDetails.Status,
				company.DriveDetails.ScheduledDate,
				company.DriveDetails.NumberOfOffers,
				company.DriveDetails.NumberHired,
				rolesJSON,
				company.DriveDetails.MinCGPA,
				branchesJSON,
				company.DriveDetails.Notes,
				company.CreatedBy,
			)
		}

		if err != nil {
			return err
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetCompaniesByOfficerID(officerID string) ([]*entity.Company, error) {
	rows, err := r.db.Query(`
		SELECT company_id FROM companies WHERE assigned_to = ?
	`, officerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*entity.Company
	for rows.Next() {
		var companyID string
		if err := rows.Scan(&companyID); err != nil {
			return nil, err
		}

		company, err := r.GetCompanyByID(companyID)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}

func (r *Repository) AssignCompanyToOfficer(companyID, officerID, assignedBy string) error {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update company assignment
	_, err = tx.Exec(`
		UPDATE companies SET assigned_to = ?, updated_at = ? WHERE company_id = ?
	`, officerID, time.Now(), companyID)
	if err != nil {
		return err
	}

	// Record assignment in company_assignments table
	assignmentID := uuid.NewString()
	_, err = tx.Exec(`
		INSERT INTO company_assignments (
			assignment_id, company_id, officer_id, assigned_by
		) VALUES (?, ?, ?, ?)
	`, assignmentID, companyID, officerID, assignedBy)
	if err != nil {
		return err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllCompanies() ([]*entity.Company, error) {
	rows, err := r.db.Query(`
		SELECT company_id FROM companies
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*entity.Company
	for rows.Next() {
		var companyID string
		if err := rows.Scan(&companyID); err != nil {
			return nil, err
		}

		company, err := r.GetCompanyByID(companyID)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}

func (r *Repository) ImportCompaniesFromCSV(companies []*entity.Company) ([]string, error) {
	var importedIDs []string

	for _, company := range companies {
		id, err := r.CreateCompany(company)
		if err != nil {
			log.Printf("Error importing company %s: %v", company.CompanyName, err)
			continue
		}
		importedIDs = append(importedIDs, id)
	}

	return importedIDs, nil
}

func (r *Repository) UpdateApprovalStatus(companyID, status, notes string, updatedBy string) error {
	_, err := r.db.Exec(`
		UPDATE companies SET 
			approval_status = ?, 
			approval_notes = ?, 
			updated_at = ?
		WHERE company_id = ?
	`, status, notes, time.Now(), companyID)

	return err
}
