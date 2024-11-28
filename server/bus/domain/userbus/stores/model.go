package userdb

import (
	"clanplan/server/bus/domain/userbus"
	"clanplan/server/bus/types/name"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type user struct {
	ID           uuid.UUID `bson:"_id"`
	Username     string    `bson:"name"`
	Email        string    `bson:"email"`
	PasswordHash []byte    `bson:"password_hash"`
	Enabled      bool      `bson:"enabled"`
	DateCreated  time.Time `bson:"date_created"`
	DateUpdated  time.Time `bson:"date_updated"`
}

func toDbUser(u userbus.User) user {
	db := user{
		ID:           u.ID,
		Username:     u.Username.String(),
		Email:        u.Email.Address,
		PasswordHash: u.PasswordHash,
		Enabled:      u.Enabled,
		DateCreated:  u.DateCreated,
		DateUpdated:  u.DateUpdated,
	}
	return db
}

func toBusUser(u user) userbus.User {
	name, err := name.Parse(u.Username)
	if err != nil {
		panic(err)
	}
	bus := userbus.User{
		ID:           u.ID,
		Username:     name,
		Email:        mail.Address{Address: u.Email},
		PasswordHash: u.PasswordHash,
		Enabled:      u.Enabled,
		DateCreated:  u.DateCreated,
		DateUpdated:  u.DateUpdated,
	}
	return bus
}

func toBusUsers(users []user) []userbus.User {
	busUsers := make([]userbus.User, 0, len(users))
	for _, u := range users {
		busUsers = append(busUsers, toBusUser(u))
	}
	return busUsers
}
