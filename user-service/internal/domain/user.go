package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string    `json:"id" bson:"_id"`
	Email        string    `json:"email" bson:"email" validate:"required,email"`
	PasswordHash string    `json:"-" bson:"password_hash"`
	FirstName    string    `json:"first_name" bson:"first_name" validate:"required"`
	LastName     string    `json:"last_name" bson:"last_name" validate:"required"`
	Roles        []string  `json:"roles" bson:"roles"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

func NewUser(email, password, firstName, lastName string) (*User, error) {
	user := &User{
		ID:        generateID(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Roles:     []string{"customer"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func generateID() string {
	return uuid.New().String()
}
