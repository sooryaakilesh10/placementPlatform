package entity

import (
	"backend/pkg/common"
	"errors"
	"regexp"
	"slices"
	"time"
)

// helper
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type User struct {
	UserID    string
	UserName  string
	Email     string
	Pass      string
	Role      string
	CreatedAt time.Time
}

func NewUser(userName, email, pass, role string) (*User, error) {
	user := &User{
		UserName:  userName,
		Email:     email,
		Pass:      pass,
		Role:      role,
		CreatedAt: time.Now(),
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) validate() error {
	if u.UserName == "" {
		return errors.New("user name cannot be empty")
	}
	if u.Email == "" {
		return errors.New("email cannot be empty")
	} else if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}
	if u.Pass == "" {
		return errors.New("password cannot be empty")
	}
	if !slices.Contains(common.ValidRoles, u.Role) {
		return errors.New("invalid role")
	}
	return nil
}
