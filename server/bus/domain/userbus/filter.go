package userbus

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type QueryFilter struct{
	ID               *uuid.UUID
	Username         *string
	Email            *mail.Address
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}