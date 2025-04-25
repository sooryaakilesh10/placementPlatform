package data

import (
	"backend/services/datad/entity"
	"database/sql"
	"errors"
)

// ErrNotFound is returned when a requested entity is not found.
var ErrNotFound = errors.New("entity not found")

type Repository struct {
	db *sql.DB
}

func NewDataRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateCompany(company *entity.CompanyData) error {
	query := `INSERT INTO company_data (id, company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr_details) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, company.CompanyID, company.CompanyName, company.CompanyAddress, company.Drive, company.TypeOfDrive, company.FollowUp, company.IsContacted, company.Remarks, company.ContactDetails, company.HRDetails)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetCompany(id string) (*entity.CompanyData, error) {
	var company entity.CompanyData
	query := `
		SELECT 
			id, company_name, company_address, drive, type_of_drive, 
			follow_up, is_contacted, remarks, contact_details, hr_details 
		FROM company_data 
		WHERE id = ?
	`
	err := r.db.QueryRow(query, id).Scan(
		&company.CompanyID,
		&company.CompanyName,
		&company.CompanyAddress,
		&company.Drive,
		&company.TypeOfDrive,
		&company.FollowUp,
		&company.IsContacted,
		&company.Remarks,
		&company.ContactDetails,
		&company.HRDetails,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &company, nil
}

func (r *Repository) GetCompanyByName(name string) (*entity.CompanyData, error) {
	var company entity.CompanyData
	query := `
		SELECT 
			id, company_name, company_address, drive, type_of_drive, 
			follow_up, is_contacted, remarks, contact_details, hr_details 
		FROM company_data 
		WHERE name = ?
	`
	err := r.db.QueryRow(query, name).Scan(
		&company.CompanyID,
		&company.CompanyName,
		&company.CompanyAddress,
		&company.Drive,
		&company.TypeOfDrive,
		&company.FollowUp,
		&company.IsContacted,
		&company.Remarks,
		&company.ContactDetails,
		&company.HRDetails,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &company, nil
}

func (r *Repository) UpdateCompany(company *entity.CompanyData) error {
	query := `
		INSERT INTO company_data_approval 
		(id, company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr_details)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Handle optional HRDetails
	hrDetails := sql.NullString{
		String: company.HRDetails,
		Valid:  company.HRDetails != "",
	}

	_, err := r.db.Exec(query,
		company.CompanyID,
		company.CompanyName,
		company.CompanyAddress,
		company.Drive,
		company.TypeOfDrive,
		company.FollowUp,
		company.IsContacted,
		company.Remarks,
		company.ContactDetails,
		hrDetails,
	)

	return err
}

func (r *Repository) GetAwaitingApproval() ([]*entity.CompanyData, error) {
	var companies []*entity.CompanyData
	query := `
		SELECT 
			id, company_name, company_address, drive, type_of_drive, 
			follow_up, is_contacted, remarks, contact_details, hr_details 
		FROM company_data_approval 
		WHERE is_contacted = false
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var company entity.CompanyData
		err := rows.Scan(
			&company.CompanyID,
			&company.CompanyName,
			&company.CompanyAddress,
			&company.Drive,
			&company.TypeOfDrive,
			&company.FollowUp,
			&company.IsContacted,
			&company.Remarks,
			&company.ContactDetails,
			&company.HRDetails,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, &company)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return companies, nil
}