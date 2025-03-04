package user

import (
	"time"

	"github.com/google/uuid"
)

type UserID string

type User struct {
	ID        UserID
	Name      string
	Email     string
	Phone     string
	Password  string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func NewUserID() UserID {
	return UserID(uuid.New().String())
}
