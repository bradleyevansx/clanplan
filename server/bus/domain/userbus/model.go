package userbus

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        mail.Address
	PasswordHash []byte
	Enabled      bool
	DateCreated  time.Time
	DateUpdated  time.Time
}
type NewUser struct {
	Username        string
	Email           mail.Address
	Password        string
	PasswordConfirm string
}
type UpdateUser struct {
	Username *string
	Email    *mail.Address
	Password *string
	Enabled  *bool
}
