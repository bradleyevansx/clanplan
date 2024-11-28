package userapp

import (
	"clanplan/server/bus/domain/userbus"
	"encoding/json"
	"time"
)

type queryParams struct {
	ID               string
	Username         string
	Email            string
	StartCreatedDate string
	EndCreatedDate   string
}

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash []byte
	Enabled      bool
	DateCreated  string
	DateUpdated  string
}

// Encode implements the encoder interface.
func (app User) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppUser(bus userbus.User) User {
	return User{
		ID:           bus.ID.String(),
		Username:     bus.Username.String(),
		Email:        bus.Email.String(),
		PasswordHash: bus.PasswordHash,
		Enabled:      bus.Enabled,
		DateCreated:  bus.DateCreated.Format(time.RFC3339),
		DateUpdated:  bus.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []userbus.User) []User {
	app := make([]User, len(users))
	for i, usr := range users {
		app[i] = toAppUser(usr)
	}

	return app
}
