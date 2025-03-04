package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserID string

type User struct {
	ID        UserID
	Name      string
	Email     string
	Phone     string
	Password  string `json:"-"` // exclude password from json responses
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func NewUserID() UserID {
	return UserID(uuid.New().String())
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
