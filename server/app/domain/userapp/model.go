package userapp

import (
	"clanplan/server/bus/domain/userbus"
	"clanplan/server/bus/types/name"
	"encoding/json"
	"fmt"
	"net/mail"
	"time"
)

type queryParams struct {
	Page             string
	Row              string
	Order            string
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

type NewUser struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func (app *NewUser) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

func toNewBusUser(app NewUser) (userbus.NewUser, error) {
	un, err := name.Parse(app.Password)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: ", err)
	}

	em, err := mail.ParseAddress(app.Email)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: ", err)
	}

	return userbus.NewUser{
		Username: un,
		Email:    *em,
		Password: app.Password,
	}, nil
}
