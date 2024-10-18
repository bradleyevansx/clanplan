package userbus

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type QueryFilter struct {
	ID               *uuid.UUID    `bson:"_id"`
	Username         *string       `bson:"name"`
	Email            *mail.Address `bson:"email"`
	StartCreatedDate *time.Time    `bson:"date_created"`
	EndCreatedDate   *time.Time    `bson:"date_created"`
}
