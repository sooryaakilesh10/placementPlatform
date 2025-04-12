package repository

import (
	"backend/services/userd/entity"
	"database/sql"
	"errors"

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

func (r *Repository) CreateUser(user *entity.User) (string, error) {
	var existingUserID string

	err := r.db.QueryRow("SELECT user_id FROM users WHERE email = ?;", user.Email).Scan(&existingUserID)
	if err == nil {
		user.UserID = existingUserID
		return existingUserID, errors.New("user already present")
	} else if err != sql.ErrNoRows {
		return "", err
	}

	newUUID := uuid.NewString()
	user.UserID = newUUID

	query := `
		INSERT INTO users (user_id, user_name, email, pass, role)
		VALUES (?, ?, ?, ?, ?);
	`
	_, err = r.db.Exec(query, user.UserID, user.UserName, user.Email, user.Pass, user.Role)
	if err != nil {
		return "", err
	}

	return user.UserID, nil
}

func (r *Repository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow(`
	SELECT 
		user_id, user_name, email, pass, role 
		FROM users WHERE email = ?;
	`, email).Scan(&user.UserID,
		&user.UserName,
		&user.Email,
		&user.Pass,
		&user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow(`
		SELECT 
			user_id, user_name, email, pass, role 
			FROM users WHERE user_id = ?;
		`, id).Scan(&user.UserID,
		&user.UserName,
		&user.Email,
		&user.Pass,
		&user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
