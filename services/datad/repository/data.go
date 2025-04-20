package data

import (
	"backend/services/datad/entity"
	"database/sql"
	"encoding/json"
)

type Repository struct {
	db *sql.DB
}

func NewDataRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateData(data *entity.Data) (string, error) {
	jsonData, err := json.Marshal(data.CompanyData)
	if err != nil {
		return "", err
	}

	query := `INSERT INTO data (data_id, company_data) VALUES (?, ?)`
	_, err = r.db.Exec(query, data.DataID, string(jsonData))
	if err != nil {
		return "", err
	}
	return data.DataID, nil
}

func (r *Repository) GetDataByID(id string) (*entity.Data, error) {
	query := `SELECT data_id, company_data FROM data WHERE data_id = ?`
	row := r.db.QueryRow(query, id)

	var dataID string
	var companyData string
	err := row.Scan(&dataID, &companyData)
	if err != nil {
		return nil, err
	}

	var data entity.Data
	err = json.Unmarshal([]byte(companyData), &data.CompanyData)
	if err != nil {
		return nil, err
	}
	data.DataID = dataID

	return &data, nil
}
