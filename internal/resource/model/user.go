package model

import (
	"errors"
	"rest-api/internal/validator"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Name      string
	Email     string
	Password  password
	Phone     string
	Activated bool
	Cod       int
	Version   int
	Deleted   bool
}

type UserDTO struct {
	ID    uuid.UUID `json:"user_id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Phone string    `json:"phone"`
}

type UserSaveDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type password struct {
	Plaintext *string
	Hash      []byte
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}
}

func (u *UserDTO) ToModel() *User {
	return &User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}
}

func (u *UserSaveDTO) ToModel() (*User, error) {
	user := &User{
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}

	err := user.Password.Set(u.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(user.Phone != "", "phone", "must be provided")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
